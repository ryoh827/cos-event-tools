# Repository Guidelines

このリポジトリは README に記載されている通り、イベント撮影後の事務作業を効率化するためのツール群を集約します。Go製CLIを複数保持する前提で、すべての変更は以下の指針に従ってください。README と本ガイドは常に同期させ、情報が重複する場合も最新の方針を優先します。

## Project Structure & Module Organization
- `go.mod` はモジュール名 `github.com/ryoh827/cos-event-tools` を宣言します。新規パッケージを追加する際はこのモジュール配下に配置してください。
- CLI は `cmd/<tool>/main.go` へ置き、共有ロジックが必要なら `internal/` や `pkg/` に切り出します。既存例として `cmd/cos-mkdir/` があり、構造の参考にできます。
- ビルド生成物や検証用のCSV・JSONなどは git 管理から除外し、`git status -sb` で確認してからコミットします。

## Build, Test, and Development Commands
- `go build ./cmd/<tool>`：個別ツールのビルド。`cos-mkdir` を確認する場合は `go build ./cmd/cos-mkdir` を実行します。`GOOS`/`GOARCH` を指定して配布バイナリを検証することも可能です。
- `go run ./cmd/<tool> [args]`：CSV変換やファイル生成など目的別の検証を素早く行えます。例: `go run ./cmd/cos-mkdir sample/event.csv`。
- `go test ./...`：すべてのパッケージを対象にユニットテストを実行します。必要に応じ `-run` で絞り込み、`-v` で詳細ログを出力します。
- `go install ./cmd/<tool>`：頻繁に使うツールを `$GOBIN` へ配置し、日常のオペレーションフローに組み込みます。
- `./scripts/install-all-bins.sh`：`cmd/` 配下の `package main` を自動検出し、各CLIをまとめてビルドして `~/bin` へ配置します。配置先を変更する場合は `./scripts/install-all-bins.sh -d ~/bin` のように `-d` を指定します。

## Coding Style & Naming Conventions
- Goの標準スタイルに従い、`gofmt` と `goimports` を必ず適用します。Lintを追加する場合は設定ファイルを共有し、CI 導入時には README で告知します。
- 外向けAPI以外はローワーキャメルケース（例: `parseDate`）。コマンド名・ディレクトリ名はハイフンなしの短い英語を推奨します。
- CLI では `run` 関数が引数処理とエラー制御を担当し、個別ロジックは小さなヘルパーに分割してテストしやすくします。

## Testing Guidelines
- 各パッケージに `_test.go` を置き、テーブル駆動テストでフォーマット変換・ファイル生成・エラーパスを網羅します。`t.TempDir()` を活用して実ファイルを安全に扱ってください。
- 異常系（欠損列、無効な日付、既存ディレクトリなど）を必ず含め、`testdata/` 配下にサンプル入力を保存すると再現性が高まります。
- CLI の追加機能は `go run` を用いた手動検証結果も PR に記録し、README に使用例が必要か検討します。必要なら `README.md` の「cos-event-tools」紹介文に新しいユースケースを追記してください。

## Commit & Pull Request Guidelines
- コミットメッセージは命令形・短文（例: `feat: add csv sanitizing helper`）。Issue があれば本文に参照を追加し、実行したコマンドや確認結果を列挙します。
- PR では目的・アーキテクチャ判断・テスト結果（`go build`, `go test`, `go run` 等）を記述し、挙動が分かるログやスクリーンショットを添付してください。
- ツール追加や仕様変更時は README と AGENTS.md の更新要否を確認し、他のツールにも影響する場合は共通ルールを優先して調整します。レビュー前に `git diff README.md AGENTS.md` で記載内容が齟齬なく揃っているかチェックするのがおすすめです。

## Development Workflow & Security Notes
- Issue → branch（`feature/dir-generator` など）→ PR の順で進め、作業前後に `git status -sb` と `go env GOPATH` を確認してキャッシュや生成物の差分を把握します。
- CSV に含まれる個人情報はローカルから持ち出さず、共有時はダミーデータを使用します。ディレクトリ権限やファイル名規則を変更する際はガイドへ理由と手順を追記してください。端末を共有環境で使う場合は `GOENV` などに機微情報を残さないよう注意します。
