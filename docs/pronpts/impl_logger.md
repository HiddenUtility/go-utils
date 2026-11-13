# Go Text Logger 実装仕様書の作成・実装指示

## 1. 目的

Go で簡易的なテキストログライブラリを実装する。

実装する機能は以下のとおり。

- ログレベル管理
- ログレコード生成
- ログファイルへの書き込み
- ログのバッファリング
- ファイルサイズに基づくログローテーション
- 指定世代数を超えた古いログファイルの削除
- Logger の Close 時の自動書き込み

本仕様書に記載された要件を優先すること。

**重要：独自の機能追加、仕様変更、過剰な抽象化を行わないこと。**

仕様に記載されていない事項については、既存プロジェクトの設計・規約に従うこと。判断できない場合は、勝手に実装せずユーザーに質問すること。

---

## 2. 実装前の確認事項

実装を開始する前に、以下を確認する。

1. リポジトリのディレクトリ構成
2. 既存の `go.mod`
3. 既存の Go ソースコード
4. 既存のコーディング規約
5. 既存のテスト構成
6. 既存の依存パッケージ

既存プロジェクトがある場合は、その構成に合わせて実装する。

既存のコードを無断で削除・変更しないこと。

原則として Go 標準ライブラリのみを使用する。

不必要な外部依存を追加しない。

---

## 3. 基本仕様

### 3.1. 使用言語

- Go
- Go の標準的な設計・命名規則に従う。
- 不必要な外部依存を追加しない。
- 不必要な goroutine を使用しない。
- 非同期処理は実装しない。

### 3.2. ログファイル

出力先の例：

`logs/app.log`

ログファイルの親ディレクトリが存在しない場合は、自動的に作成する。

パス操作には `path/filepath` を使用する。

### 3.3. デフォルト設定

| 設定 | デフォルト値 |
|---|---|
| MaxFileNum | 7 |
| MaxFileSize | 204800 Bytes |
| Realtime | false |

MaxFileSize は 200 KiB（200 × 1024 Bytes）とする。

MaxFileNum は現在のログファイルを含む最大ファイル数とする。

---

## 4. LogLevel

### 4.1. 型定義

Go には enum がないため、独自の整数型と定数を使用する。

```go
type LogLevel int

const (
    CRITICAL LogLevel = iota
    ERROR
    WARNING
    INFO
    DEBUG
)
```

指定された定数名を維持する。

### 4.2. 文字列変換

LogLevel に String メソッドを実装する。

```go
func (level LogLevel) String() string
```

変換結果：

| 定数 | 文字列 |
|---|---|
| CRITICAL | CRITICAL |
| ERROR | ERROR |
| WARNING | WARNING |
| INFO | INFO |
| DEBUG | DEBUG |

以下のように利用できること。

```go
level := INFO
fmt.Println(level.String())
```

不正な LogLevel が渡された場合は `"UNKNOWN"` を返す。

---

## 5. LogRecode

### 5.1. 構造体定義

ログ1件分の情報を保持する構造体。

```go
type LogRecode struct {
    Timestamp string
    Level     LogLevel
    Msg       string
    PID       int
}
```

各フィールドの意味：

| フィールド | 型 | 説明 |
|---|---|---|
| Timestamp | string | ログ生成時刻 |
| Level | LogLevel | ログレベル |
| Msg | string | ログメッセージ |
| PID | int | プロセスID |

### 5.2. Timestamp

ログレコードの生成時に現在時刻を取得する。

以下の形式で保存する。

`2026-11-11T12:13:45.123456`

要件：

- ISO 8601 に準拠した形式とする。
- ローカル時刻を使用する。
- タイムゾーン情報は付与しない。
- マイクロ秒精度とする。
- 形式は `YYYY-MM-DDTHH:MM:SS.ffffff` とする。
- Go の `time.Now()` を使用する。
- `time.Local` に基づくローカル時刻とする。

実装例：

```go
timestamp := time.Now().Format("2006-01-02T15:04:05.000000")
```

UTC に変換しないこと。

### 5.3. PID

ログレコード生成時に、実行中プロセスの PID を取得する。

`os.Getpid()` を使用する。

### 5.4. ログ生成メソッド

以下の5メソッドを実装する。

```go
func NewCritical(msg string) LogRecode
func NewError(msg string) LogRecode
func NewWarning(msg string) LogRecode
func NewInfo(msg string) LogRecode
func NewDebug(msg string) LogRecode
```

各メソッドの処理：

1. 引数 msg を受け取る。
2. 現在時刻を取得する。
3. PID を取得する。
4. メソッドに対応する LogLevel を設定する。
5. LogRecode インスタンスを生成して返す。

対応関係：

| メソッド | LogLevel |
|---|---|
| NewCritical | CRITICAL |
| NewError | ERROR |
| NewWarning | WARNING |
| NewInfo | INFO |
| NewDebug | DEBUG |

例：

```go
record := NewInfo("Application started")
```

### 5.5. ToLog

ログレコードをテキストに変換するメソッドを実装する。

```go
func (record LogRecode) ToLog() string
```

出力形式：

`timestamp [level] msg`

例：

```text
2026-11-11T12:13:45.123456 [INFO] Application started
```

要件：

- Timestamp を先頭にする。
- LogLevel を `[]` で囲む。
- Msg を末尾にする。
- PID はログ文字列に含めない。
- 戻り値は string とする。
- 改行文字は含めない。改行はファイル書き込み時に追加する。

---

## 6. Logger

### 6.1. 構造体定義

Logger 本体を以下のフィールドで構成する。

```go
type Logger struct {
    LogPath     string
    Recodes     []LogRecode
    MaxFileNum  uint32
    MaxFileSize uint32
    Realtime    bool
    NowSize     uint32
}
```

各フィールドの意味：

| フィールド | 型 | 説明 |
|---|---|---|
| LogPath | string | ログ出力先 |
| Recodes | []LogRecode | 未出力ログ |
| MaxFileNum | uint32 | 最大ログファイル数 |
| MaxFileSize | uint32 | 最大ログファイルサイズ（Bytes） |
| Realtime | bool | 即時書き込み設定 |
| NowSize | uint32 | 現在のログファイルサイズ（Bytes） |

### 6.2. 初期化メソッド

以下の2種類を実装する。

#### New

```go
func New(logPath string) (*Logger, error)
```

ファイルパスだけを指定する初期化メソッド。

以下のデフォルト設定を使用する。

```go
MaxFileNum = 7
MaxFileSize = 200 * 1024
Realtime = false
```

#### NewWith

```go
func NewWith(
    logPath string,
    maxFileNum uint32,
    maxFileSize uint32,
    realtime bool,
) (*Logger, error)
```

全オプションを指定して初期化する。

#### 初期化処理

両メソッドは同じ初期化処理を利用する。

1. LogPath を設定する。
2. LogPath の親ディレクトリを取得する。
3. 親ディレクトリが存在しない場合は `os.MkdirAll()` で作成する。
4. ログファイルが存在するか確認する。
5. 存在する場合は `os.Stat()` でファイルサイズを取得し、NowSize に設定する。
6. 存在しない場合は NowSize を 0 にする。
7. Recodes を空のスライスで初期化する。
8. 指定された設定値を各フィールドに保存する。

ファイルサイズは実際のファイルサイズを使用すること。

ログファイル自体が存在しない場合は、初期化時点ではファイルを作成せず、空のバッファ状態として扱う。

MaxFileNum = 0 は無効な設定として扱い、error を返す。

---

## 7. Logger のログメソッド

Logger に以下の5メソッドを実装する。

```go
func (logger *Logger) Critical(msg string) error
func (logger *Logger) Error(msg string) error
func (logger *Logger) Warning(msg string) error
func (logger *Logger) Info(msg string) error
func (logger *Logger) Debug(msg string) error
```

### 7.1. 共通処理

各メソッドは以下の処理を行う。

1. 対応する LogRecode を生成する。
2. Recodes に追加する。
3. Realtime が true の場合は Out() を実行する。
4. Realtime が false の場合は書き込みを行わず、バッファに保持する。

### 7.2. Realtime の動作

| Realtime | 動作 |
|---|---|
| true | ログ追加のたびに Out() を実行 |
| false | ログをバッファに蓄積 |

Out() がエラーになった場合は、そのエラーを呼び出し元に返す。

---

## 8. Out メソッド

### 8.1. メソッド定義

```go
func (logger *Logger) Out() error
```

Recodes に保持されているログをログファイルに書き出す。

### 8.2. 基本処理

1. Recodes が空の場合は何もせず nil を返す。
2. Recodes の各 LogRecode を ToLog() で文字列に変換する。
3. 各ログの末尾に改行を追加する。
4. UTF-8 でログファイルに追記する。
5. 書き込んだバイト数を NowSize に加算する。
6. 全レコードの書き込みが完了したら Recodes を空にする。
7. ファイルサイズが MaxFileSize を超えた場合はローテーションを実行する。

### 8.3. 書き込み方式

既存ファイルを上書きしてはならない。

必ず追記モードで書き込む。

以下のフラグを使用する。

```go
os.O_APPEND | os.O_CREATE | os.O_WRONLY
```

新しいログファイルが存在しない場合は作成する。

### 8.4. バッファの保持

重要なルール：

- 書き込みに成功した場合のみ Recodes を空にする。
- 書き込みに失敗した場合は Recodes を保持する。
- エラーを握りつぶさず、呼び出し元に返す。
- 書き込み失敗時にログデータを勝手に破棄しない。

`os.File.Write()` の戻り値を確認すること。

書き込みバイト数が要求バイト数に満たない場合もエラーとして扱う。

部分書き込みが発生した場合に、次回 Out() で同じデータが重複出力されないよう、書き込み済みバイト数を管理すること。

ただし、今回の基本設計では過剰な状態管理を追加しないこと。必要な最小限の実装とする。

### 8.5. 書き込み完了とローテーション

ログデータの書き込みとローテーションを区別して管理する。

ログデータの書き込みが完了した後にローテーションが失敗した場合、書き込み済みデータを Recodes に残してはならない。

この場合はローテーションエラーを返し、書き込み済みログを再出力しないこと。

---

## 9. ログローテーション

### 9.1. 発動条件

Out() によってバッファを書き込んだ後、以下を判定する。

```go
logger.NowSize > logger.MaxFileSize
```

最大サイズと等しい場合はローテーションしない。

最大サイズを超えた場合にのみ実行する。

### 9.2. ファイル構成

例えば MaxFileNum = 3 の場合：

ローテーション前：

```text
logs/
    app.log
    app.log.1
    app.log.2
```

ローテーション後：

```text
logs/
    app.log
    app.log.1
    app.log.2
```

世代の意味：

| ファイル | 意味 |
|---|---|
| app.log | 現在のログ |
| app.log.1 | 1世代前 |
| app.log.2 | 2世代前 |

最も古いファイルは削除する。

### 9.3. ローテーション手順

MaxFileNum = 3 の場合を例にする。

1. app.log.2 が存在すれば削除する。
2. app.log.1 が存在すれば app.log.2 にリネームする。
3. app.log が存在すれば app.log.1 にリネームする。
4. 新しい app.log を作成する。
5. NowSize を 0 に設定する。

### 9.4. 一般化した処理

MaxFileNum の値に応じて同じ処理を行う。

重要：

- 現行ファイルを含めた総数を MaxFileNum 以下にする。
- 古いファイルから削除する。
- 世代番号の大きいファイルから順番に処理する。
- ファイルのリネーム順序を間違えない。
- 存在しないファイルをリネームしようとしない。
- MaxFileNum = 1 の場合は旧 app.log を削除し、新しい app.log を作成する。

ローテーションに失敗した場合はエラーを返す。

### 9.5. ローテーション後

ローテーション完了後：

```go
logger.NowSize = 0
```

新しいログファイルへの書き込みに備える。

---

## 10. Close メソッド

Go には Rust の Drop に相当するデストラクタが存在しない。

GC によるメモリ解放時にログ出力を実行することは保証できない。

そのため、Logger に Close() メソッドを実装する。

```go
func (logger *Logger) Close() error
```

### 10.1. 動作

1. Out() を実行する。
2. Out() が成功した場合は nil を返す。
3. Out() が失敗した場合はエラーを返す。

Close() は未出力のログをファイルへ書き出すための明示的な終了処理とする。

### 10.2. 使用方法

以下のように使用できること。

```go
logger, err := New("logs/app.log")
if err != nil {
    return err
}
defer logger.Close()
```

Close() を複数回呼び出しても安全であること。

一度正常に書き込みが完了している場合、2回目の呼び出しで同じログが重複出力されないこと。

### 10.3. 重要事項

Logger は Close() を呼び出すことを前提に使用する。

Go の GC による自動解放にログ出力を依存させないこと。

Close() のエラーを呼び出し元が確認できるようにすること。

---

## 11. エラー処理

ファイル操作に関するエラーは Go 標準の error インターフェースで返す。

対象：

- ディレクトリ作成失敗
- ファイル情報取得失敗
- ファイルオープン失敗
- ファイル書き込み失敗
- ファイル削除失敗
- ファイルリネーム失敗
- ローテーション失敗
- 無効な初期化設定

エラーを不必要に握りつぶさない。

エラーが発生した場合は、呼び出し元が原因を特定できるように、操作内容を付加して返す。

例：

```go
return fmt.Errorf("failed to write log file: %w", err)
```

ただし、ログ出力の失敗によってアプリケーションの本来の処理を妨害しないよう、利用側でエラーを無視できる設計にする。

---

## 12. 必須テスト

テストには `testing` パッケージを使用する。

一時ディレクトリには `t.TempDir()` を使用する。

実際の logs ディレクトリやユーザーのファイルを変更しないこと。

### 12.1. LogLevel

- 各定数が正しい文字列になること。
- String() が正常に動作すること。
- 不正な値が UNKNOWN になること。

### 12.2. LogRecode

- NewCritical が CRITICAL を設定すること。
- NewError が ERROR を設定すること。
- NewWarning が WARNING を設定すること。
- NewInfo が INFO を設定すること。
- NewDebug が DEBUG を設定すること。
- Msg が正しく保存されること。
- PID が現在のプロセス ID と一致すること。
- Timestamp が指定形式になっていること。
- ToLog() が正しい形式を返すこと。

### 12.3. Logger 初期化

- New() がデフォルト値を設定すること。
- NewWith() が指定値を設定すること。
- 親ディレクトリが存在しない場合に作成されること。
- 既存ログファイルのサイズを NowSize に設定すること。
- 未作成のログファイルの場合 NowSize が 0 になること。
- MaxFileNum = 0 を拒否すること。

### 12.4. バッファリング

- Realtime = false ではログメソッド呼び出し時にファイルへ書き込まれないこと。
- Recodes にログが追加されること。
- Out() でログが書き込まれること。
- 書き込み成功後 Recodes が空になること。
- Out() を複数回呼び出しても不要な重複書き込みが発生しないこと。

### 12.5. Realtime

- Realtime = true ではログメソッド呼び出し時にファイルへ書き込まれること。
- 書き込み後 Recodes が空になること。

### 12.6. ローテーション

- MaxFileSize 以下ではローテーションされないこと。
- MaxFileSize を超えるとローテーションされること。
- app.log が app.log.1 になること。
- app.log.1 が app.log.2 になること。
- 最大世代数を超えたファイルが削除されること。
- MaxFileNum = 1 が正常に動作すること。
- ローテーション後 NowSize が 0 になること。

### 12.7. Close

- Close() 呼び出し時にバッファが書き込まれること。
- Close() を複数回呼び出しても重複書き込みされないこと。
- バッファが空の場合でも安全に呼び出せること。
- 書き込み失敗時にエラーが返されること。

### 12.8. エラー処理

- 無効なパスを指定した場合にエラーが返されること。
- 書き込み失敗時にバッファが失われないこと。
- ローテーション失敗時にエラーが返されること。
- 書き込み済みログが再出力されないこと。

---

## 13. 実装上の注意事項

以下の事項を必ず守ること。

1. 仕様にない機能を勝手に追加しない。
2. ログ形式を変更しない。
3. Timestamp にタイムゾーンを追加しない。
4. PID をログ文字列に追加しない。
5. MaxFileSize の単位を変更しない。
6. MaxFileNum の数え方を変更しない。
7. ログを上書きしない。追記する。
8. 書き込み失敗時にバッファを破棄しない。
9. GC による自動解放をログ出力のトリガーにしない。
10. Close() を実装する。
11. 不必要な外部パッケージを追加しない。
12. 不必要な goroutine を使用しない。
13. 既存プロジェクトのコードを無断で破壊しない。

---

## 14. 実装完了条件

以下をすべて満たした場合のみ実装完了とする。

- [ ] LogLevel が実装されている。
- [ ] LogLevel の String() が実装されている。
- [ ] LogRecode が実装されている。
- [ ] 5種類のログ生成関数が実装されている。
- [ ] Timestamp がローカル時刻・マイクロ秒精度である。
- [ ] PID が取得される。
- [ ] ToLog() が実装されている。
- [ ] Logger が実装されている。
- [ ] New() が実装されている。
- [ ] NewWith() が実装されている。
- [ ] 5種類の Logger ログメソッドが実装されている。
- [ ] Out() が実装されている。
- [ ] Realtime が動作する。
- [ ] ログファイルの追記が動作する。
- [ ] ログローテーションが動作する。
- [ ] 古いログファイルが削除される。
- [ ] Close() によるバッファ出力が動作する。
- [ ] エラー処理が仕様どおりである。
- [ ] 必須テストがすべて実装されている。

---

## 15. 品質確認

実装後、以下を実行する。

```bash
gofmt -w .
go vet ./...
go test ./...
```

すべて成功すること。

失敗した場合は原因を調査して修正し、再実行する。

テストを通すためだけにテストを削除したり、不適切に警告を無視したりしてはならない。

---

## 16. 作業完了時の報告

実装完了後、以下を報告する。

1. 作成・変更したファイル
2. 実装した機能
3. 追加した依存パッケージとその理由
4. 実施したテスト
5. gofmt / go vet / go test の実行結果
6. 未対応事項
7. 仕様から変更した箇所がある場合は、その理由

仕様どおりに実装できなかった項目がある場合、完了したように報告してはならない。

必ず未完了項目と理由を明記する。
