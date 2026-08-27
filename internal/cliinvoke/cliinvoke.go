// Package cliinvoke runs the pinned clean-invisible-text CLI's
// check/explain/fix --json against a single file and classifies the result
// into the Clean/Cleaned/Warning/Error state model from
// docs/specification.md States.
package cliinvoke

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// Command names accepted by Run, matching the CLI's subcommands.
const (
	Check   = "check"
	Explain = "explain"
	Fix     = "fix"
)

// Finding is one entry of the CLI's --json findings array
// (docs/cli.md in go-clean-invisible-text).
type Finding struct {
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Offset      int    `json:"offset"`
	Rune        string `json:"rune"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Action      string `json:"action"` // "remove", "replace", or "warn"
	Replacement string `json:"replacement"`
}

// FileResult is the CLI's --json report for a single file.
type FileResult struct {
	Path     string    `json:"path"`
	Findings []Finding `json:"findings"`
	Changed  bool      `json:"changed"`
	Error    *string   `json:"error"`
}

// Result is the outcome of one CLI invocation against a single temp file.
type Result struct {
	ExitCode int
	File     FileResult
}

// State is the Clean/Cleaned/Warning/Error model from
// docs/specification.md States.
type State int

const (
	StateClean State = iota
	StateCleaned
	StateWarning
	StateError
)

func (s State) String() string {
	switch s {
	case StateClean:
		return "Clean"
	case StateCleaned:
		return "Cleaned"
	case StateWarning:
		return "Warning"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// State classifies Result per docs/specification.md States: exit 0 with no
// error is Clean; exit 1 with no error and no action=="warn" finding is
// Cleaned; exit 1 with no error and at least one action=="warn" finding is
// Warning; anything else is Error.
func (r Result) State() State {
	if r.File.Error != nil {
		return StateError
	}
	switch r.ExitCode {
	case 0:
		return StateClean
	case 1:
		for _, f := range r.File.Findings {
			if f.Action == "warn" {
				return StateWarning
			}
		}
		return StateCleaned
	default:
		return StateError
	}
}

// Run invokes binaryPath (a pinned clean-invisible-text binary, see
// internal/cliasset) as `command --json [--keep-warnings] path` and parses
// its single-file JSON report. It never includes file content in an error —
// only the command and path.
func Run(binaryPath, command, path string, keepWarnings bool) (Result, error) {
	args := []string{command, "--json"}
	if keepWarnings {
		args = append(args, "--keep-warnings")
	}
	args = append(args, path)

	cmd := exec.Command(binaryPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	runErr := cmd.Run()
	exitCode := 0
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			return Result{}, fmt.Errorf("cliinvoke: running %s %s: %w", binaryPath, command, runErr)
		}
	}

	var results []FileResult
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		return Result{}, fmt.Errorf("cliinvoke: parsing %s output: %w", command, err)
	}
	if len(results) != 1 {
		return Result{}, fmt.Errorf("cliinvoke: expected exactly 1 file result from %s, got %d", command, len(results))
	}

	return Result{ExitCode: exitCode, File: results[0]}, nil
}
