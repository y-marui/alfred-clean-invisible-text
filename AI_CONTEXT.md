## Reference Order

AI reads the following in order at the start of a task:

1. `README.md` (overview, setup)
2. `docs/specification.md` (interaction model, states, accessibility — the normative source for Workflow behavior)

Read as needed (any order):
- `CONTRIBUTING.md` (PR/Issue rules)
- `docs/dependency-policy.md` (how the CLI is pinned and verified)
- `docs/decisions/` (ADRs — architecture decisions and their rationale)
- `docs/architecture.md` (module/component structure, once implementation exists)
- `docs/file-map.md` (file-level dependencies; explore and append if stale or missing)
- `docs/ui-design.md` (not applicable — Alfred's native Script Filter/Universal Action UI is used as-is; see `docs/specification.md` Accessibility and keyboard flow)

`DEVELOPING.md` does not exist yet — there is no implementation to build/test. Add it (build, test, naming conventions) when Go source lands (issue #4).

## Project Overview

A planned Alfred Workflow that is a thin macOS frontend for the
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text)
CLI. It does not implement Unicode cleaning rules itself — see
[ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md). Currently at the
specification stage: `docs/specification.md` and ADR 0001–0002 are settled;
no Go source exists in this repository yet (tracked by the v0.1 Workflow
milestone, issues #2–#5).

### Technology Stack

- Go (the Workflow itself will be implemented in Go, matching the CLI it
  wraps, once implementation starts — see issue #4)
- No source code yet; this repository is currently docs-only

### Main Directories

| Path | Role |
|---|---|
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
