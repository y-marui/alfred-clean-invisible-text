# Clean Invisible Text for Alfred

> **このファイルは正本(日本語版)です。**
> 英語版(参照)は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

危険な不可視 Unicode 文字をローカルでレビュー・クリーニングするための Alfred Workflow。

> **Status:** [リリース済み](https://github.com/y-marui/alfred-clean-invisible-text/releases/latest)
> (署名・notarize済み)。Apple Silicon上で両方のエントリーポイントの動作を確認済み
> — 残作業(Intel実機・Galleryスクリーンショット)は
> [ロードマップ](https://github.com/y-marui/alfred-clean-invisible-text/issues/1) 参照。

この Workflow は
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text)
の薄い macOS フロントエンドであり、Unicode クリーニングのルール自体は独自実装しない
([ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md) 参照)。

## Requirements

- Alfred 5 以降
- macOS 13 (Ventura) 以降、Intel または Apple Silicon — 固定しているCLI自体の
  下限に合わせている。詳細は [ADR 0003](docs/decisions/0003-v1-compatibility-and-upgrade-policy.md) 参照

## Setup

署名済みの[最新リリース](https://github.com/y-marui/alfred-clean-invisible-text/releases/latest)の
`.alfredworkflow` をダウンロードするか、最新のCLIピンを使うためにソースからビルドする:

```bash
git clone https://github.com/y-marui/alfred-clean-invisible-text
cd alfred-clean-invisible-text
make fetch-cli       # 固定されたCLIバイナリをダウンロード・検証する
make build-workflow  # → dist/*.alfredworkflow
```

`.alfredworkflow` をダブルクリックしてAlfredに読み込む。キーワード
(`cit`、クリップボード対象)とUniversal Action(選択テキスト対象、Alfredの
Universal Actionsパレットから)のどちらも、追加のセットアップなしでそのまま
動作する。

## Usage

クリップボードに対する `cit` キーワード、または選択テキストへのUniversal Action
でトリガーする([docs/specification.md](docs/specification.md) に完全な
相互作用モデル・状態・アクセシビリティに関する記述がある)。その後、以下のいずれかを選ぶ:

| アクション | 説明 |
|---|---|
| Check | テキストを検査し、findingsの有無を1行で要約表示する |
| Reveal | 変更を書き込まずに、すべてのfinding(コードポイント・名前・カテゴリ・位置)を表示する |
| Clean | CLIのクリーナーを実行し、成功後にクリップボードのプレーンテキストを置き換える |
| Copy report | デフォルトでは元のテキストを除外した、findingsの構造化レポートをコピーする |

Check/Reveal/Clean の結果行では: **Enter** でレポートをコピー(元のテキストは除外)、
**⌘+Enter** で元のテキストを含むレポートをコピーする。Clean の結果がWarning状態の
場合のみ、**⇧+Enter** で未分類の文字を保持したまま Clean を再実行する。

## Documentation

- [docs/specification.md](docs/specification.md) — Alfred Workflow の仕様
- [docs/dependency-policy.md](docs/dependency-policy.md) — CLIの固定・検証方法
- [docs/release-process.md](docs/release-process.md) — Workflowリリースの作成・公開手順
- [docs/alfred-gallery-readiness.md](docs/alfred-gallery-readiness.md) — Alfred Gallery掲載チェックリスト
- [docs/decisions/](docs/decisions/) — アーキテクチャ決定記録 (ADR)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md) を参照。

## License

[MIT](LICENSE)

---
*この文書には英語版 [README.md](README.md) があります。編集時は同一コミットで更新してください。*
