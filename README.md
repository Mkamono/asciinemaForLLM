# asciinemaForLLM

asciinema の録画ファイルを LLM が読みやすい形式に変換するツールです。

## 機能

- asciinema の .cast ファイルを構造化された Markdown 形式に変換
- コマンドの入力と出力を明確に分離
- 実行時間、開始・終了タイムスタンプの表示
- 端末エスケープシーケンスの自動除去
- 録画から変換まで一括処理

## インストール

### GitHubリリースから（推奨）
```bash
# リリースページからバイナリをダウンロード
# https://github.com/Mkamono/asciinemaForLLM/releases
```

### go installから
```bash
go install github.com/Mkamono/asciinemaForLLM@latest
```

### ソースからビルド
```bash
git clone https://github.com/Mkamono/asciinemaForLLM.git
cd asciinemaForLLM
go build -o asciinema-for-llm
```

## クイックスタート

### 最も実用的な使い方（推奨）

```bash
# 録画開始 → CSV出力 → 元ファイル削除まで一括処理
asciinema-for-llm record my_session --output=csv --cleanup

# ターミナルでコマンドを実行
# 例: ls, pwd, echo "Hello World"

# 録画終了（Ctrl+D または exit）
```

終了すると `my_session_formatted.csv` ファイルのみが残り、元の .cast ファイルは自動削除されます。

### その他の使い方

#### 既存の .cast ファイルをCSV変換

```bash
# CSV形式で変換 + 元ファイル削除
asciinema-for-llm file existing_session.cast --output=csv --cleanup
```

#### 構造化テキスト形式が必要な場合

```bash
# 人間が読みやすい形式で出力
asciinema-for-llm record my_session.cast --cleanup

# 既存ファイルを構造化テキストに変換
asciinema-for-llm file existing_session.cast --cleanup
```

#### パイプで使用する場合（従来の方法）

```bash
# CSV形式
cat session.cast | asciinema-for-llm format --output=csv

# 構造化テキスト形式
cat session.cast | asciinema-for-llm format
```

## 使い方

### 基本的な使用方法

#### 1. 標準入力からの変換（従来の方法）

```bash
# パイプでファイルを渡す
cat demo.cast | asciinema-for-llm

# format サブコマンドを明示的に指定
cat demo.cast | asciinema-for-llm format
```

#### 2. 録画から変換まで一括処理

```bash
# 録画開始（ファイル名は自動生成）
asciinema-for-llm record

# ファイル名を指定して録画
asciinema-for-llm record my_session.cast

# 録画後に元の .cast ファイルを削除
asciinema-for-llm record my_session.cast --cleanup
```

#### 3. 既存ファイルの変換

```bash
# 既存の .cast ファイルを変換
asciinema-for-llm file demo.cast

# 出力ファイル名を指定
asciinema-for-llm file demo.cast output.md

# 変換後に元ファイルを削除
asciinema-for-llm file demo.cast output.md --cleanup
```

### コマンド一覧

| コマンド                | 説明                                                                  |
| ----------------------- | --------------------------------------------------------------------- |
| `format`                | 標準入力から .cast データを読み取り、フォーマット済み Markdown を出力 |
| `record [filename]`     | asciinema 録画を開始し、終了後に自動でフォーマット                    |
| `file <input> [output]` | 既存の .cast ファイルをフォーマット                                   |
| `--help`, `-h`          | ヘルプメッセージを表示                                                |

### オプション

| オプション  | 説明                                                                 |
| ----------- | -------------------------------------------------------------------- |
| `--cleanup` | 処理後に元の .cast ファイルを削除（record、file コマンドで使用可能） |

## 出力例

```markdown
# Terminal Session Analysis
Recorded at: 2025-01-01 12:00:00
Terminal: 80x24
Shell: /bin/bash

## Command 1
**Command:** `ls -la`
**Start Time:** 1.234s
**End Time:** 2.567s
**Duration:** 1.333s
**Output:**
```
total 16
drwxr-xr-x  4 user user  128 Jan  1 12:00 .
drwxr-xr-x  3 user user   96 Jan  1 11:59 ..
-rw-r--r--  1 user user 1234 Jan  1 12:00 file.txt
```

## Command 2
**Command:** `exit`
**Start Time:** 5.000s
**End Time:** 5.001s
**Duration:** 0.001s
**Output:**
```
(no output)
```
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
