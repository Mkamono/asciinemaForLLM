# asciinemaForLLM

asciinema の録画ファイルを LLM が読みやすい形式に変換するツールです。

## 機能

- asciinema の .cast ファイルを構造化テキストやCSV形式に変換
- コマンドの入力と出力を明確に分離
- 実行時間、開始・終了タイムスタンプの表示
- **コマンドの終了ステータス（exit code）を取得・表示**
- **作業ディレクトリ情報を抽出・表示**
- 端末エスケープシーケンスの自動除去
- 録画から変換まで一括処理
- LLM向けに最適化されたCSV出力形式

## 前提

### asciinemaのインストール

https://docs.asciinema.org/getting-started/

miseでの入れ方を紹介

```bash
mise use -g python
pip install --user pipx
mise use -g pipx:asciinema
```

## インストール

### mise から
```bash
mise use -g go
mise use -g go:github.com/Mkamono/asciinemaForLLM
```

### go installから
```bash
go install github.com/Mkamono/asciinemaForLLM@latest
```

### ソースからビルド
```bash
git clone https://github.com/Mkamono/asciinemaForLLM.git
cd asciinemaForLLM
go build -o asciinemaForLLM
```

## クイックスタート

### 最も実用的な使い方（推奨）

```bash
# 録画開始 → CSV出力 → 元ファイル削除まで一括処理
asciinemaForLLM record my_session --output=csv --cleanup

# ターミナルでコマンドを実行
# 例: ls, pwd, echo "Hello World"

# 録画終了（Ctrl+D または exit）
```

終了すると `my_session_formatted.csv` ファイルのみが残り、元の .cast ファイルは自動削除されます。(cleanupオプション)

### その他の使い方

#### 既存の .cast ファイルをCSV変換

```bash
# CSV形式で変換 + 元ファイル削除
asciinemaForLLM file existing_session.cast --output=csv --cleanup
```

#### 構造化テキスト形式が必要な場合

```bash
# 人間が読みやすい形式で出力
asciinemaForLLM record my_session.cast --cleanup

# 既存ファイルを構造化テキストに変換
asciinemaForLLM file existing_session.cast --cleanup
```

#### パイプで使用する場合（従来の方法）

```bash
# CSV形式
cat session.cast | asciinemaForLLM format --output=csv

# 構造化テキスト形式
cat session.cast | asciinemaForLLM format
```

## 使い方

### 基本的な使用方法

#### 1. 標準入力からの変換（従来の方法）

```bash
# パイプでファイルを渡す
cat demo.cast | asciinemaForLLM

# format サブコマンドを明示的に指定
cat demo.cast | asciinemaForLLM format
```

#### 2. 録画から変換まで一括処理

```bash
# 録画開始（ファイル名は自動生成）
asciinemaForLLM record

# ファイル名を指定して録画
asciinemaForLLM record my_session.cast

# 録画後に元の .cast ファイルを削除
asciinemaForLLM record my_session.cast --cleanup
```

#### 3. 既存ファイルの変換

```bash
# 既存の .cast ファイルを変換
asciinemaForLLM file demo.cast

# 出力ファイル名を指定
asciinemaForLLM file demo.cast output.md

# 変換後に元ファイルを削除
asciinemaForLLM file demo.cast output.md --cleanup
```

### コマンド一覧

| コマンド                | 説明                                                               |
| ----------------------- | ------------------------------------------------------------------ |
| `format`                | 標準入力から .cast データを読み取り、構造化テキストまたはCSVを出力 |
| `record [filename]`     | asciinema 録画を開始し、終了後に自動でフォーマット                 |
| `file <input> [output]` | 既存の .cast ファイルをフォーマット                                |
| `--help`, `-h`          | ヘルプメッセージを表示                                             |

### オプション

| オプション        | 説明                                                                 |
| ----------------- | -------------------------------------------------------------------- |
| `--cleanup`       | 処理後に元の .cast ファイルを削除（record、file コマンドで使用可能） |
| `--output=FORMAT` | 出力形式を指定（structured\|csv、デフォルト: structured）            |

## 出力例

### 構造化テキスト形式
```
Terminal Session (fish shell, 148x35)
Recorded: 2025-07-08 14:14:24
Working Directory: /Users/kamonomakoto/Documents/repo/asciinemaForLLM

COMMAND: echo "Hello, world"
START TIME: 3.433s
DURATION: 2.119s
EXIT CODE: 0
OUTPUT: Hello, world

COMMAND: exit
START TIME: 5.552s
DURATION: 0.002s
EXIT CODE: 0
OUTPUT: (no output)
```

### CSV形式（LLM向け）
```csv
shell,width,height,recorded,working_dir,command,start_time,duration,exit_code,output
fish,148,35,2025-07-08 14:14:24,/Users/kamonomakoto/Documents/repo/asciinemaForLLM,"echo ""Hello, world""",3.433,2.119,0,"Hello, world"
fish,148,35,2025-07-08 14:14:24,/Users/kamonomakoto/Documents/repo/asciinemaForLLM,exit,5.552,0.002,0,(no output)
```

## 開発

### テストの実行

```bash
# 全テストケースを実行
mise run test

# 個別のテストケース
cat test/echo/demo.cast | go run main.go | diff - test/echo/expectation
```

### 新しいテストケースの追加

1. `test/新しいテスト名/` ディレクトリを作成
2. `demo.cast` ファイル（asciinema 録画ファイル）を配置
3. `expectation` ファイル（期待される出力）を配置
4. `mise run test` でテスト実行

## ライセンス

MIT License
