## Reference Order

AI reads the following in order at the start of a task:

1. `README.md` (overview, setup)
2. `DEVELOPING.md` (build, test, implementation conventions)

Read as needed (any order):
- `CONTRIBUTING.md` (PR/Issue rules)
- `docs/specification.md` (interaction model, states, accessibility — the normative source for Workflow behavior)
- `docs/dependency-policy.md` (how the CLI is pinned and verified)
- `docs/release-process.md` (how a Workflow release is cut and published)
- `docs/alfred-gallery-readiness.md` (Alfred Gallery submission checklist and open blockers)
- `docs/decisions/` (ADRs — architecture decisions and their rationale)
- `docs/architecture.md` (module/component structure)
- `docs/alfred-workflow-notes/workflow-object-schema.md` (reverse-engineered `info.plist` object schema — Alfred doesn't document this; read before touching `workflow/info.plist`)
- `docs/file-map.md` (file-level dependencies; explore and append if stale or missing)
- `docs/ui-design.md` (not applicable — Alfred's native Script Filter/Universal Action UI is used as-is; see `docs/specification.md` Accessibility and keyboard flow)

## Project Overview

An Alfred Workflow that is a thin macOS frontend for the
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text)
CLI. It does not implement Unicode cleaning rules itself — see
[ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md). Specification
(`docs/specification.md`, ADR 0001–0002) is settled; the privacy-safe
clipboard/temp-file primitives (`internal/clipboard`, `internal/tempinput`,
issue #3), the pinned/verified CLI binary selection (`internal/cliasset`,
issue #5), the Check/Reveal/Clean/Copy report action logic (`internal/action`,
`internal/cliinvoke`, `internal/scriptfilter`, `cmd/clean-invisible-text-alfred`),
and the Alfred wiring itself (`workflow/info.plist`, `scripts/build-workflow.sh`
→ `dist/*.alfredworkflow`) all exist. Both entry points — the Keyword
(`cit`, clipboard) and the Universal Action (selected text) — are fully
wired in `workflow/info.plist` with no manual setup, and both have been
verified working via Alfred's own Workflow debugger on Apple Silicon; see
[docs/alfred-workflow-notes/workflow-object-schema.md](docs/alfred-workflow-notes/workflow-object-schema.md) for the
reverse-engineered plist schema this relies on. Verification on real Intel
hardware remains open but is optional/best-effort, not a blocker (issue
#4) — Intel Macs are increasingly rare, and the universal binary itself is
already verified via `lipo` at build time.

### Technology Stack

- Go (see `go.mod` for the toolchain version)
- No third-party Go modules — `internal/clipboard` shells out to the macOS
  system binaries `pbcopy`/`pbpaste`/`osascript`

### Main Directories

| Path | Role |
|---|---|
| `cmd/clean-invisible-text-alfred/` | The binary Alfred invokes |
| `workflow/` | `info.plist` (the Alfred object graph), `icon.png` |
| `internal/action/` | Check/Reveal/Clean/Copy report orchestration |
| `internal/cliinvoke/` | Pinned CLI invocation and Clean/Cleaned/Warning/Error state classification |
| `internal/scriptfilter/` | Alfred Script Filter JSON types |
| `internal/clipboard/` | macOS pasteboard plain-text read/write |
| `internal/tempinput/` | Private, single-use temp file for `check`/`explain`/`fix` input |
| `internal/cliasset/` | Pinned CLI version/checksums, runtime binary selection and verification |
| `docs/` | Specification, dependency policy, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`, see below) |

## Applied Charter Principles

- Charter reference: use `docs/dev-charter/CHARTER_INDEX.md` to find the relevant topic, then read only that file
- YAGNI, minimal diff scope, reuse existing patterns before adding new ones — `docs/dev-charter/PRINCIPLES.md`
- Secrets and pre-commit security gates — `docs/dev-charter/SECURITY_POLICY.md`
- Public-facing text (README, CLI/error output, commit/PR text) is English; internal Japanese is fine — `docs/dev-charter/LANGUAGE_POLICY.md`
- Multi-step, sub-issue-tracked work happens on an `epic/<name>` branch off `main`, reported to the parent issue on creation — `docs/dev-charter/PROJECT_LIFECYCLE.md`
- Do not directly edit files under `docs/dev-charter/`; changes go through an Issue in the dev-charter repository and `git subtree pull`

## Document Sync Rule

When a spec, rule, or structural change happens, update the related documentation
in the same piece of work. This includes files under `docs/` as well as root files
such as `AI_CONTEXT.md` and `README.md`.

## Project-Specific Rules

- Unicode detection/cleaning behavior is never reimplemented here — it is defined exclusively by `go-clean-invisible-text` (`CONTRIBUTING.md`, `docs/dependency-policy.md`)
- A change that alters observable Workflow behavior must update `docs/specification.md` or add an ADR under `docs/decisions/` — closed issues are not the source of truth
- Roadmap and task tracking live in GitHub Issues/Milestones under this repository (issue #1 is the roadmap), not in Markdown files in this repo
- A PR into a non-default branch (e.g. an `epic/<name>` branch) does not auto-close the issue it references (`Refs #N`, not `Closes #N`); close the issue manually once its epic-branch PR merges

## AI Tool Assignments

- **Tools in use**: Claude Code, Codex, GitHub Copilot, Gemini CLI, local LLM (Ollama)
- **Canonical responsibilities**: `docs/dev-charter/AI_COLLABORATION_RULES.md`, "AI Tool Responsibilities" and "Rules for Multi-AI Usage"
- **Project-specific overrides**: none

## Prohibited Actions

- Reimplementing or duplicating Unicode category/cleaning rules from `go-clean-invisible-text`
- Adding telemetry or network transmission outside explicit update/release-download operations (`SECURITY.md`)
- Committing secrets or credentials
- Direct edits under `docs/dev-charter/`
