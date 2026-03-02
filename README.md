# cos-event-tools

イベント撮影などで撮影したときの作業を効率化するために作成したツールなどを保管しておくリポジトリ

## ツール一覧

### cos-mkdir

CSV を読み込み、イベント日付・イベント名・名前からディレクトリを作成する CLI です。

```bash
go run ./cmd/cos-mkdir events.csv
```

カレントディレクトリに CSV ファイルが 1 つだけある場合は、引数を省略して実行できます。

```bash
go run ./cmd/cos-mkdir
```

### cos-zip

指定したディレクトリ直下の各ディレクトリをそれぞれ zip 化します。上書きの挙動は次のとおりです。

- 既存の zip はデフォルトで y/n を確認し、`n` または何も入力せずに Enter を押すとそのファイルはスキップします。
- `-f` を付けると確認せずに上書きします。

```bash
# デフォルト（既存 zip は y/n 確認）
go run ./cmd/cos-zip <root-directory>

# 既存 zip を確認なしで上書きする場合
go run ./cmd/cos-zip -f <root-directory>
```

## 一括インストール（`~/bin`）

`cmd/` 配下の `package main` を自動検出して、各CLIをまとめてビルドし `~/bin` に配置できます。
新しいCLIを `cmd/<tool>/main.go` として追加した場合も、スクリプト側のメンテナンスは不要です。

```bash
# デフォルト: ~/bin に配置
./scripts/install-all-bins.sh

# 配置先を指定する場合
./scripts/install-all-bins.sh -d ~/bin
```
