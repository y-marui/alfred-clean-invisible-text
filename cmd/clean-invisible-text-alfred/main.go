// Command clean-invisible-text-alfred is the binary the packaged Alfred
// Workflow invokes (see workflow/info.plist). Alfred always runs scripts
// with the working directory set to the workflow's own bundle, so the
// pinned CLI binaries at assets/bin/ resolve as a relative path — see
// internal/cliasset and docs/dependency-policy.md.
//
// Subcommands, each corresponding to one workflow/info.plist node's script:
//
//	list <source> [text]   — the Check/Reveal/Clean chooser (Script Filter)
//	run [action]            — runs the chosen action (Script Filter); action
//	                          falls back to $action if omitted; source/text
//	                          come from the $source/$text env vars Alfred
//	                          exports from the chosen item's variables
//	copy-report              — writes $report to the clipboard (Run Script)
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/y-marui/alfred-clean-invisible-text/internal/action"
	"github.com/y-marui/alfred-clean-invisible-text/internal/cliasset"
	"github.com/y-marui/alfred-clean-invisible-text/internal/scriptfilter"
)

func main() {
	if len(os.Args) < 2 {
		writeErrorResponse(fmt.Errorf("missing subcommand"))
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		runList()
	case "run":
		runAction()
	case "copy-report":
		runCopyReport()
	default:
		writeErrorResponse(fmt.Errorf("unknown subcommand %q", os.Args[1]))
		os.Exit(1)
	}
}

func runList() {
	if len(os.Args) < 3 {
		writeErrorResponse(fmt.Errorf("list: missing source"))
		return
	}
	source := action.Source(os.Args[2])
	text := ""
	if len(os.Args) > 3 {
		text = os.Args[3]
	}
	writeResponse(action.List(source, text))
}

func runAction() {
	act := os.Getenv("action")
	if len(os.Args) > 2 && os.Args[2] != "" {
		act = os.Args[2]
	}

	switch act {
	case "check", "reveal", "clean", "clean-keep-warnings":
		// handled below, once the pinned CLI binary is resolved
	default:
		writeResponse(scriptfilter.Response{Items: []scriptfilter.Item{{Title: "Nothing to do"}}})
		return
	}

	binaryPath, err := cliasset.ResolvePath(assetsDir())
	if err != nil {
		writeErrorResponse(err)
		return
	}

	req := action.Request{
		Source: action.Source(os.Getenv("source")),
		Text:   os.Getenv("text"),
	}

	var resp scriptfilter.Response
	switch act {
	case "check":
		resp, err = action.Check(binaryPath, req)
	case "reveal":
		resp, err = action.Reveal(binaryPath, req)
	case "clean":
		resp, err = action.Clean(binaryPath, req, false)
	case "clean-keep-warnings":
		resp, err = action.Clean(binaryPath, req, true)
	}
	if err != nil {
		writeErrorResponse(err)
		return
	}
	writeResponse(resp)
}

func runCopyReport() {
	report := os.Getenv("report")
	if err := action.CopyReport(report); err != nil {
		notify("Clean Invisible Text", "Could not copy the report")
		os.Exit(1)
	}
	notify("Clean Invisible Text", "Report copied to clipboard")
}

func assetsDir() string {
	if dir := os.Getenv("ALFRED_CLEAN_ASSETS_DIR"); dir != "" {
		return dir
	}
	return "assets/bin"
}

func writeResponse(resp scriptfilter.Response) {
	if err := resp.Write(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "clean-invisible-text-alfred: writing response:", err)
		os.Exit(1)
	}
}

// writeErrorResponse always emits valid Script Filter JSON — even on an
// unexpected internal failure — so Alfred shows a readable error row
// instead of an empty/unparseable result.
func writeErrorResponse(err error) {
	resp := scriptfilter.Response{Items: []scriptfilter.Item{{Title: "Error", Subtitle: err.Error()}}}
	_ = resp.Write(os.Stdout)
}

// notify shows a macOS notification for terminal (Run Script) actions that
// have no Alfred result row of their own to convey their outcome in text —
// docs/specification.md Accessibility and keyboard flow. Best-effort: a
// notification failure must not be treated as the underlying action
// failing.
func notify(title, message string) {
	script := fmt.Sprintf("display notification %q with title %q", message, title)
	_ = exec.Command("osascript", "-e", script).Run()
}
