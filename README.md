# Clean Invisible Text for Alfred

A planned Alfred Workflow for reviewing and cleaning dangerous invisible Unicode characters locally.

> Status: specification and roadmap. The Workflow is not implemented yet.

The Workflow is a thin macOS frontend for [go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text). It does not implement Unicode cleaning rules independently.

## Planned actions

- Check selected text or clipboard content.
- Reveal findings without exposing the full text in logs.
- Clean text and copy the result back to the clipboard.
- Report counts, code points, and applied actions.

Work in progress belongs in GitHub Issues and Projects. Stable Alfred behavior belongs in this repository's specifications.

## License

MIT
