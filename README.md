# cos-event-tools

イベント撮影などで撮影したときの作業を効率化するために作成したツールなどを保管しておくリポジトリ

## ツール一覧

### cos-mkdir

CSV を読み込み、イベント日付・イベント名・名前からディレクトリを作成する CLI です。

```bash
go run ./cmd/cos-mkdir events.csv
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

### cos-edrm

Lightroom や Photoshop で加工されたファイル名に付与される `-Edit` サフィックスを取り除き、元のファイル名にリネームする CLI です。指定したディレクトリ以下を再帰的に探索し、引数を省略するとカレントディレクトリを対象にします。同名ファイルが既に存在する場合は安全のためエラーで停止します。

```bash
# カレントディレクトリ配下を探索
go run ./cmd/cos-edrm

# 任意のディレクトリを指定
go run ./cmd/cos-edrm /path/to/exported/photos
```
