# cos-event-tools

イベント撮影などで撮影したときの作業を効率化するために作成したツールなどを保管しておくリポジトリ

## ツール一覧

### cos-mkdir

CSV を読み込み、イベント日付・イベント名・名前からディレクトリを作成する CLI です。

```bash
go run ./cmd/cos-mkdir events.csv
```

### cos-zip

指定したディレクトリ直下の各ディレクトリをそれぞれ zip 化します。既存の zip はデフォルトで上書き前に y/n を確認し、`-f` で無確認の上書きができます。

```bash
# デフォルト（既存 zip は y/n 確認）
go run ./cmd/cos-zip <root-directory>

# 既存 zip を確認なしで上書きする場合
go run ./cmd/cos-zip -f <root-directory>
```
