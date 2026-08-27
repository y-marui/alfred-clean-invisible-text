# Clean Invisible Text for Alfred

> **このファイルは正本(日本語版)です。**
> 英語版(参照)は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

危険な不可視 Unicode 文字をローカルでレビュー・クリーニングするための、計画中の Alfred Workflow。

> **Status:** 仕様策定・ロードマップの段階。Workflow 本体は未実装
> ([ロードマップ](https://github.com/y-marui/alfred-clean-invisible-text/issues/1) 参照)。

この Workflow は
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text)
の薄い macOS フロントエンドであり、Unicode クリーニングのルール自体は独自実装しない
([ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md) 参照)。

## Setup

まだ利用できない。最初のリリース後は、署名済みの `.alfredworkflow` を
ダブルクリックして Alfred に読み込む形になる予定
(このアーティファクトのビルド・検証は issue #4 で追跡)。

## Usage

選択テキストへの Universal Action、またはクリップボードに対するキーワードで
トリガーする、計画中のアクション
([docs/specification.md](docs/specification.md) に完全な相互作用モデル・状態・
アクセシビリティに関する記述がある):

| アクション | 説明 |
|---|---|
| Check | テキストを検査し、findingsの有無を1行で要約表示する |
| Reveal | 変更を書き込まずに、すべてのfinding(コードポイント・名前・カテゴリ・位置)を表示する |
| Clean | CLIのクリーナーを実行し、成功後にクリップボードのプレーンテキストを置き換える |
| Copy report | デフォルトでは元のテキストを除外した、findingsの構造化レポートをコピーする |

## Documentation

- [docs/specification.md](docs/specification.md) — Alfred Workflow の仕様
- [docs/dependency-policy.md](docs/dependency-policy.md) — CLIの固定・検証方法
- [docs/decisions/](docs/decisions/) — アーキテクチャ決定記録 (ADR)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md) を参照。

## License

[MIT](LICENSE)

---
*この文書には英語版 [README.md](README.md) があります。編集時は同一コミットで更新してください。*
