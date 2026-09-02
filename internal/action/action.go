// Package action implements the Check/Reveal/Clean/Copy report actions from
// docs/specification.md, orchestrating internal/clipboard,
// internal/tempinput, and internal/cliinvoke into Alfred Script Filter
// responses.
//
// Modifier keys on a result row (docs/specification.md Accessibility and
// keyboard flow — alternate actions live on the existing row, not a
// secondary menu):
//   - plain Enter: copy a report (excludes the original text)
//   - cmd+Enter:   copy a report that includes the original text
//   - shift+Enter (Clean only, Warning state only): re-run Clean with
//     --keep-warnings before writing the clipboard
package action

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/y-marui/alfred-clean-invisible-text/internal/cliinvoke"
	"github.com/y-marui/alfred-clean-invisible-text/internal/clipboard"
	"github.com/y-marui/alfred-clean-invisible-text/internal/scriptfilter"
	"github.com/y-marui/alfred-clean-invisible-text/internal/tempinput"
)

// Source is where the input text for an action comes from, per
// docs/specification.md Entry points.
type Source string

const (
	SourceClipboard Source = "clipboard"
	SourceSelection Source = "selection"
)

// Request identifies the input for Check/Reveal/Clean: either the selection
// captured at invocation time (Text, when Source is SourceSelection) or a
// fresh clipboard read at execution time (Source is SourceClipboard, Text
// ignored).
type Request struct {
	Source Source
	Text   string
}

// resolve returns the text to run an action against, or an Alfred item
// describing why it couldn't (an empty or unsupported clipboard) — this is
// never a Go error, since an empty/unsupported clipboard is an expected,
// user-facing condition, not an exceptional one.
func (r Request) resolve() (string, *scriptfilter.Item) {
	if r.Source == SourceSelection {
		return r.Text, nil
	}
	text, err := clipboard.ReadPlainText()
	switch {
	case err == nil:
		return text, nil
	case errors.Is(err, clipboard.ErrEmpty):
		return "", &scriptfilter.Item{Title: "Error", Subtitle: "Clipboard is empty", Valid: scriptfilter.BoolPtr(false)}
	case errors.Is(err, clipboard.ErrUnsupported):
		return "", &scriptfilter.Item{Title: "Error", Subtitle: "Clipboard doesn't contain text", Valid: scriptfilter.BoolPtr(false)}
	default:
		return "", &scriptfilter.Item{Title: "Error", Subtitle: "Could not read the clipboard", Valid: scriptfilter.BoolPtr(false)}
	}
}

// List returns the top-level Check/Reveal/Clean chooser (Alfred result
// rows). For SourceSelection, text (the selection) is carried forward via
// each item's variables so the next step doesn't need to re-resolve it; for
// SourceClipboard, no text is carried — the clipboard is read fresh at
// execution time.
func List(source Source, text string) scriptfilter.Response {
	baseVars := map[string]string{"source": string(source)}
	if source == SourceSelection {
		baseVars["text"] = text
	}

	row := func(uid, title, subtitle, arg string) scriptfilter.Item {
		vars := make(map[string]string, len(baseVars)+1)
		for k, v := range baseVars {
			vars[k] = v
		}
		vars["action"] = arg
		return scriptfilter.Item{
			UID:       uid,
			Title:     title,
			Subtitle:  subtitle,
			Arg:       arg,
			Valid:     scriptfilter.BoolPtr(true),
			Variables: vars,
		}
	}

	return scriptfilter.Response{Items: []scriptfilter.Item{
		row("check", "Check", "Inspect for invisible Unicode characters", "check"),
		row("reveal", "Reveal", "Show every finding without changing anything", "reveal"),
		row("clean", "Clean", "Remove invisible characters and update the clipboard", "clean"),
	}}
}

// Check runs `check` against req's input and returns a one-row summary.
func Check(binaryPath string, req Request) (scriptfilter.Response, error) {
	input, errItem := req.resolve()
	if errItem != nil {
		return scriptfilter.Response{Items: []scriptfilter.Item{*errItem}}, nil
	}

	result, err := runCLI(binaryPath, cliinvoke.Check, input, false)
	if err != nil {
		return scriptfilter.Response{}, err
	}
	if result.State() == cliinvoke.StateError {
		return scriptfilter.Response{Items: []scriptfilter.Item{errorItem(result)}}, nil
	}

	n := len(result.File.Findings)
	subtitle := "No findings"
	if n > 0 {
		subtitle = fmt.Sprintf("%d finding(s)", n)
	}
	item := scriptfilter.Item{Title: result.State().String(), Subtitle: subtitle}
	attachReportMods(&item, result.File.Findings, result.State(), input)
	return scriptfilter.Response{Items: []scriptfilter.Item{item}}, nil
}

// Reveal runs `explain` against req's input and returns one row per finding
// (or a single "no findings" row). Every row carries the same
// complete-report Copy report action — docs/specification.md treats Copy
// report as a follow-up on the Reveal result as a whole, not on an
// individual finding.
func Reveal(binaryPath string, req Request) (scriptfilter.Response, error) {
	input, errItem := req.resolve()
	if errItem != nil {
		return scriptfilter.Response{Items: []scriptfilter.Item{*errItem}}, nil
	}

	result, err := runCLI(binaryPath, cliinvoke.Explain, input, false)
	if err != nil {
		return scriptfilter.Response{}, err
	}
	if result.State() == cliinvoke.StateError {
		return scriptfilter.Response{Items: []scriptfilter.Item{errorItem(result)}}, nil
	}

	if len(result.File.Findings) == 0 {
		item := scriptfilter.Item{Title: "No findings", Subtitle: "This text has no invisible-character findings"}
		attachReportMods(&item, nil, result.State(), input)
		return scriptfilter.Response{Items: []scriptfilter.Item{item}}, nil
	}

	items := make([]scriptfilter.Item, 0, len(result.File.Findings))
	for _, f := range result.File.Findings {
		item := findingItem(f)
		attachReportMods(&item, result.File.Findings, result.State(), input)
		items = append(items, item)
	}
	return scriptfilter.Response{Items: items}, nil
}

// Clean runs `fix` against req's input and, only on success (State Clean,
// Cleaned, or Warning — never Error), replaces the clipboard with the
// cleaned result. Its result row carries the Copy report modifier actions,
// plus — only in the Warning state, and only when this run didn't already
// keep them — a shift-modifier action that re-runs Clean with
// --keep-warnings before writing the clipboard.
func Clean(binaryPath string, req Request, keepWarnings bool) (scriptfilter.Response, error) {
	input, errItem := req.resolve()
	if errItem != nil {
		return scriptfilter.Response{Items: []scriptfilter.Item{*errItem}}, nil
	}

	result, cleaned, err := runFix(binaryPath, input, keepWarnings)
	if err != nil {
		return scriptfilter.Response{}, err
	}
	if result.State() == cliinvoke.StateError {
		return scriptfilter.Response{Items: []scriptfilter.Item{errorItem(result)}}, nil
	}

	if err := clipboard.WritePlainText(cleaned); err != nil {
		return scriptfilter.Response{}, fmt.Errorf("action: writing cleaned text to clipboard: %w", err)
	}

	item := scriptfilter.Item{Title: result.State().String(), Subtitle: cleanSubtitle(result)}
	attachReportMods(&item, result.File.Findings, result.State(), input)
	if result.State() == cliinvoke.StateWarning && !keepWarnings {
		shiftVars := map[string]string{
			"action": "clean-keep-warnings",
			"source": string(req.Source),
		}
		if req.Source == SourceSelection {
			shiftVars["text"] = req.Text
		}
		item.Mods["shift"] = scriptfilter.Mod{
			Subtitle:  "Re-run keeping unclassified characters",
			Arg:       "clean-keep-warnings",
			Variables: shiftVars,
		}
	}
	return scriptfilter.Response{Items: []scriptfilter.Item{item}}, nil
}

// CopyReport replaces the clipboard with reportText, a report already built
// by Check/Reveal/Clean (via BuildReport) and threaded through as an Alfred
// variable — it never re-scans (docs/specification.md Copy report).
func CopyReport(reportText string) error {
	return clipboard.WritePlainText(reportText)
}

// BuildReport renders findings and state as a plain-text report: only
// code points, names, categories, locations, actions, and counts. If
// includeText is set, the original input is appended as a clearly labeled,
// separate section — the explicit opt-in variant from
// docs/specification.md Copy report.
func BuildReport(findings []cliinvoke.Finding, state cliinvoke.State, includeText bool, input string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s — %d finding(s)\n", state, len(findings))
	for _, f := range findings {
		fmt.Fprintf(&b, "%d:%d: %s %s [%s] -> %s\n", f.Line, f.Column, f.Rune, f.Name, f.Category, f.Action)
	}
	if includeText {
		b.WriteString("\n--- Original text ---\n")
		b.WriteString(input)
	}
	return b.String()
}

// attachReportMods sets item.Arg to the default (text-excluded) Copy report
// action and adds a cmd-modifier action for the text-included variant.
func attachReportMods(item *scriptfilter.Item, findings []cliinvoke.Finding, state cliinvoke.State, input string) {
	item.Valid = scriptfilter.BoolPtr(true)
	item.Arg = "copy-report"
	item.Variables = map[string]string{
		"action": "copy-report",
		"report": BuildReport(findings, state, false, input),
	}
	item.Mods = map[string]scriptfilter.Mod{
		"cmd": {
			Subtitle: "Copy report, including the original text",
			Arg:      "copy-report",
			Variables: map[string]string{
				"action": "copy-report",
				"report": BuildReport(findings, state, true, input),
			},
		},
	}
}

func cleanSubtitle(result cliinvoke.Result) string {
	switch result.State() {
	case cliinvoke.StateClean:
		return "Already clean — clipboard unchanged"
	case cliinvoke.StateCleaned:
		return fmt.Sprintf("Removed %d finding(s) and updated the clipboard", len(result.File.Findings))
	case cliinvoke.StateWarning:
		return fmt.Sprintf("Removed %d finding(s), including unclassified characters — review the report", len(result.File.Findings))
	default:
		return ""
	}
}

func findingItem(f cliinvoke.Finding) scriptfilter.Item {
	return scriptfilter.Item{
		Title:    fmt.Sprintf("%s %s", f.Rune, f.Name),
		Subtitle: fmt.Sprintf("%s · line %d, col %d · %s", f.Category, f.Line, f.Column, f.Action),
	}
}

func errorItem(result cliinvoke.Result) scriptfilter.Item {
	msg := "The CLI could not process this input"
	if result.File.Error != nil {
		msg = *result.File.Error
	}
	return scriptfilter.Item{Title: "Error", Subtitle: msg, Valid: scriptfilter.BoolPtr(false)}
}

func runCLI(binaryPath, command, input string, keepWarnings bool) (cliinvoke.Result, error) {
	var result cliinvoke.Result
	err := tempinput.WithTempFile(input, func(path string) error {
		var runErr error
		result, runErr = cliinvoke.Run(binaryPath, command, path, keepWarnings)
		return runErr
	})
	return result, err
}

func runFix(binaryPath, input string, keepWarnings bool) (cliinvoke.Result, string, error) {
	var result cliinvoke.Result
	var cleaned string
	err := tempinput.WithTempFile(input, func(path string) error {
		var runErr error
		result, runErr = cliinvoke.Run(binaryPath, cliinvoke.Fix, path, keepWarnings)
		if runErr != nil {
			return runErr
		}
		if result.State() == cliinvoke.StateError {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		cleaned = string(data)
		return nil
	})
	return result, cleaned, err
}
