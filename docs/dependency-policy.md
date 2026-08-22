# CLI Dependency Policy

The Workflow embeds released macOS binaries from go-clean-invisible-text.

- The CLI version is pinned for each Workflow release.
- SHA-256 checksums are verified during packaging.
- Intel and Apple Silicon binaries come from the same upstream CLI release.
- Runtime downloads are not performed during ordinary text processing.
- The packaged CLI version is included in diagnostics.
- Unicode rules are never copied into Workflow scripts.
