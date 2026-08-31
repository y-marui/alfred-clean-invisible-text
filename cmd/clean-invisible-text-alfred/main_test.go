package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles this command once per test run and returns its path,
// so subprocess tests exercise the exact code path Alfred would invoke.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	out := filepath.Join(dir, "clean-invisible-text-alfred")
	cmd := exec.Command("go", "build", "-o", out, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, output)
	}
	return out
}

func runBinary(t *testing.T, bin string, env map[string]string, args ...string) (stdout string, exitCode int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode()
		}
		t.Fatalf("running binary: %v", err)
	}
	return string(out), 0
}

func TestList_Clipboard(t *testing.T) {
	bin := buildBinary(t)
	stdout, exitCode := runBinary(t, bin, nil, "list", "clipboard")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}

	var resp struct {
		Items []struct {
			Title     string            `json:"title"`
			Arg       string            `json:"arg"`
			Variables map[string]string `json:"variables"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, stdout)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(resp.Items))
	}
	for _, item := range resp.Items {
		if item.Variables["source"] != "clipboard" {
			t.Errorf("item %q source = %q, want clipboard", item.Title, item.Variables["source"])
		}
		if _, ok := item.Variables["text"]; ok {
			t.Errorf("item %q unexpectedly carries text", item.Title)
		}
	}
}

func TestList_Selection(t *testing.T) {
	bin := buildBinary(t)
	stdout, exitCode := runBinary(t, bin, nil, "list", "selection", "hello world")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}

	var resp struct {
		Items []struct {
			Variables map[string]string `json:"variables"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, stdout)
	}
	for _, item := range resp.Items {
		if item.Variables["text"] != "hello world" {
			t.Errorf("text = %q, want %q", item.Variables["text"], "hello world")
		}
	}
}

func TestRun_UnknownAction(t *testing.T) {
	bin := buildBinary(t)
	stdout, exitCode := runBinary(t, bin, map[string]string{"source": "selection", "text": "hi"}, "run", "bogus-action")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0 (should degrade gracefully)", exitCode)
	}
	if !jsonContains(t, stdout, "Nothing to do") {
		t.Errorf("stdout = %s, want a graceful \"Nothing to do\" item", stdout)
	}
}

func TestRun_MissingPinnedBinary(t *testing.T) {
	bin := buildBinary(t)
	env := map[string]string{
		"ALFRED_CLEAN_ASSETS_DIR": t.TempDir(), // empty: no pinned binary staged
		"source":                  "selection",
		"text":                    "hi",
	}
	stdout, exitCode := runBinary(t, bin, env, "run", "check")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0 (errors must render as a JSON item, not a crash)", exitCode)
	}
	if !jsonContains(t, stdout, "Error") {
		t.Errorf("stdout = %s, want an Error item", stdout)
	}
}

func TestRun_EndToEnd_RealPinnedCLI(t *testing.T) {
	assetsDir, err := filepath.Abs("../../assets/bin")
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	if _, err := os.Stat(assetsDir); err != nil {
		t.Skipf("pinned CLI binary not staged (run `make fetch-cli`): %v", err)
	}

	bin := buildBinary(t)
	env := map[string]string{
		"ALFRED_CLEAN_ASSETS_DIR": assetsDir,
		"source":                  "selection",
		"text":                    "hello​world", // ZERO WIDTH SPACE
	}
	stdout, exitCode := runBinary(t, bin, env, "run", "check")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !jsonContains(t, stdout, "1 finding") {
		t.Errorf("stdout = %s, want it to report 1 finding", stdout)
	}
}

func jsonContains(t *testing.T, stdout, substr string) bool {
	t.Helper()
	var resp struct {
		Items []struct {
			Title    string `json:"title"`
			Subtitle string `json:"subtitle"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, stdout)
	}
	for _, item := range resp.Items {
		if strings.Contains(item.Title, substr) || strings.Contains(item.Subtitle, substr) {
			return true
		}
	}
	return false
}
