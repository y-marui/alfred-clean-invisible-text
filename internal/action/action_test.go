package action

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/y-marui/alfred-clean-invisible-text/internal/cliinvoke"
	"github.com/y-marui/alfred-clean-invisible-text/internal/clipboard"
)

func fakeCLIPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("testdata/fakecli.sh")
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	return abs
}

func requireMacOS(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("clipboard-touching actions are macOS-only")
	}
	for _, bin := range []string{"pbcopy", "pbpaste", "osascript"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not found", bin)
		}
	}
}

func saveAndRestoreClipboard(t *testing.T) {
	t.Helper()
	original, err := clipboard.ReadPlainText()
	t.Cleanup(func() {
		if err == nil || err == clipboard.ErrEmpty {
			_ = clipboard.WritePlainText(original)
		}
	})
}

func TestList(t *testing.T) {
	resp := List(SourceClipboard, "")
	if len(resp.Items) != 3 {
		t.Fatalf("List(clipboard) items = %d, want 3", len(resp.Items))
	}
	for _, item := range resp.Items {
		if item.Variables["source"] != "clipboard" {
			t.Errorf("item %q source = %q, want clipboard", item.Title, item.Variables["source"])
		}
		if _, ok := item.Variables["text"]; ok {
			t.Errorf("item %q unexpectedly carries text for clipboard source", item.Title)
		}
	}

	resp = List(SourceSelection, "hello")
	for _, item := range resp.Items {
		if item.Variables["text"] != "hello" {
			t.Errorf("item %q text = %q, want %q", item.Title, item.Variables["text"], "hello")
		}
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		scenario    string
		wantTitle   string
		wantHasMods bool
		wantValid   bool
	}{
		{"clean", "Clean", true, true},
		{"cleaned", "Cleaned", true, true},
		{"warning", "Warning", true, true},
		{"file-error", "Error", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Setenv("FAKECLI_SCENARIO", tc.scenario)
			resp, err := Check(fakeCLIPath(t), Request{Source: SourceSelection, Text: "hello world"})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			if len(resp.Items) != 1 {
				t.Fatalf("items = %d, want 1", len(resp.Items))
			}
			item := resp.Items[0]
			if item.Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", item.Title, tc.wantTitle)
			}
			hasMods := item.Mods != nil && item.Mods["cmd"].Arg == "copy-report"
			if hasMods != tc.wantHasMods {
				t.Errorf("has copy-report cmd mod = %v, want %v", hasMods, tc.wantHasMods)
			}
			valid := item.Valid == nil || *item.Valid
			if valid != tc.wantValid {
				t.Errorf("Valid = %v, want %v (nil counts as true)", valid, tc.wantValid)
			}
		})
	}
}

func TestCheck_CopyReportVariants(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "cleaned")
	resp, err := Check(fakeCLIPath(t), Request{Source: SourceSelection, Text: "secret input text"})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	item := resp.Items[0]

	if strings.Contains(item.Variables["report"], "secret input text") {
		t.Error("default report unexpectedly includes the original text")
	}
	cmdReport := item.Mods["cmd"].Variables["report"]
	if !strings.Contains(cmdReport, "secret input text") {
		t.Error("cmd-modifier report does not include the original text")
	}
}

func TestReveal(t *testing.T) {
	req := Request{Source: SourceSelection, Text: "hello"}

	t.Setenv("FAKECLI_SCENARIO", "clean")
	resp, err := Reveal(fakeCLIPath(t), req)
	if err != nil {
		t.Fatalf("Reveal: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Title != "No findings" {
		t.Errorf("Reveal(clean) = %+v, want a single \"No findings\" item", resp.Items)
	}

	t.Setenv("FAKECLI_SCENARIO", "cleaned")
	resp, err = Reveal(fakeCLIPath(t), req)
	if err != nil {
		t.Fatalf("Reveal: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("items = %d, want 1 finding row", len(resp.Items))
	}
	if !strings.Contains(resp.Items[0].Title, "ZERO WIDTH SPACE") {
		t.Errorf("Title = %q, want it to mention the finding", resp.Items[0].Title)
	}

	t.Setenv("FAKECLI_SCENARIO", "file-error")
	resp, err = Reveal(fakeCLIPath(t), req)
	if err != nil {
		t.Fatalf("Reveal: %v", err)
	}
	if resp.Items[0].Title != "Error" {
		t.Errorf("Title = %q, want %q", resp.Items[0].Title, "Error")
	}
	if v := resp.Items[0].Valid; v == nil || *v {
		t.Errorf("Valid = %v, want false so Alfred won't accept Enter on an Error row", v)
	}
}

func TestClean_WritesClipboardOnSuccess(t *testing.T) {
	requireMacOS(t)
	saveAndRestoreClipboard(t)

	t.Setenv("FAKECLI_SCENARIO", "cleaned")
	if err := clipboard.WritePlainText("stale clipboard content"); err != nil {
		t.Fatalf("seeding clipboard: %v", err)
	}

	resp, err := Clean(fakeCLIPath(t), Request{Source: SourceSelection, Text: "hello"}, false)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Title != "Cleaned" {
		t.Fatalf("items = %+v, want a single Cleaned item", resp.Items)
	}

	got, err := clipboard.ReadPlainText()
	if err != nil {
		t.Fatalf("ReadPlainText after Clean: %v", err)
	}
	if got == "stale clipboard content" {
		t.Error("clipboard was not replaced")
	}
}

func TestClean_DoesNotWriteClipboardOnError(t *testing.T) {
	requireMacOS(t)
	saveAndRestoreClipboard(t)

	t.Setenv("FAKECLI_SCENARIO", "file-error")
	const sentinel = "original clipboard must survive"
	if err := clipboard.WritePlainText(sentinel); err != nil {
		t.Fatalf("seeding clipboard: %v", err)
	}

	resp, err := Clean(fakeCLIPath(t), Request{Source: SourceSelection, Text: "hello"}, false)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	if resp.Items[0].Title != "Error" {
		t.Fatalf("Title = %q, want %q", resp.Items[0].Title, "Error")
	}
	if v := resp.Items[0].Valid; v == nil || *v {
		t.Errorf("Valid = %v, want false so Alfred won't accept Enter on an Error row", v)
	}

	got, err := clipboard.ReadPlainText()
	if err != nil {
		t.Fatalf("ReadPlainText: %v", err)
	}
	if got != sentinel {
		t.Errorf("clipboard = %q, want unchanged %q", got, sentinel)
	}
}

func TestClean_WarningOffersKeepWarningsModifier(t *testing.T) {
	requireMacOS(t)
	saveAndRestoreClipboard(t)

	req := Request{Source: SourceSelection, Text: "hello"}

	t.Setenv("FAKECLI_SCENARIO", "warning")
	resp, err := Clean(fakeCLIPath(t), req, false)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	mod, ok := resp.Items[0].Mods["shift"]
	if !ok {
		t.Fatal("Warning state result has no shift modifier")
	}
	if mod.Arg != "clean-keep-warnings" {
		t.Errorf("shift mod Arg = %q, want %q", mod.Arg, "clean-keep-warnings")
	}
	if mod.Variables["source"] != "selection" || mod.Variables["text"] != "hello" {
		t.Errorf("shift mod Variables = %v, want source/text carried forward for the loop-back", mod.Variables)
	}

	t.Setenv("FAKECLI_SCENARIO", "cleaned")
	resp, err = Clean(fakeCLIPath(t), req, false)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	if _, ok := resp.Items[0].Mods["shift"]; ok {
		t.Error("Cleaned (non-Warning) state unexpectedly has a shift modifier")
	}
}

func TestCheck_EmptyClipboard(t *testing.T) {
	requireMacOS(t)
	saveAndRestoreClipboard(t)

	if err := clipboard.WritePlainText(""); err != nil {
		t.Fatalf("seeding clipboard: %v", err)
	}

	resp, err := Check(fakeCLIPath(t), Request{Source: SourceClipboard})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Title != "Error" {
		t.Fatalf("items = %+v, want a single Error item", resp.Items)
	}
	if resp.Items[0].Subtitle != "Clipboard is empty" {
		t.Errorf("Subtitle = %q, want %q", resp.Items[0].Subtitle, "Clipboard is empty")
	}
	if v := resp.Items[0].Valid; v == nil || *v {
		t.Errorf("Valid = %v, want false so Alfred won't accept Enter on an Error row", v)
	}
}

func TestBuildReport(t *testing.T) {
	findings := []cliinvoke.Finding{
		{Line: 1, Column: 2, Rune: "U+200B", Name: "ZERO WIDTH SPACE", Category: "zwsp", Action: "remove"},
	}

	report := BuildReport(findings, cliinvoke.StateCleaned, false, "the original text")
	if strings.Contains(report, "the original text") {
		t.Error("report without includeText unexpectedly contains the original text")
	}
	if !strings.Contains(report, "ZERO WIDTH SPACE") {
		t.Error("report does not mention the finding")
	}

	report = BuildReport(findings, cliinvoke.StateCleaned, true, "the original text")
	if !strings.Contains(report, "the original text") {
		t.Error("report with includeText does not contain the original text")
	}
}
