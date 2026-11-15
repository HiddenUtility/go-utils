# Codex Implementation Specification: SqlLogger (Golang)

## 1. 目的

Go言語でSQLiteを利用したシンプルなログ記録ライブラリ `SqlLogger` を実装する。

本ライブラリは、ログレコードをメモリ上に一時的に蓄積し、SQLiteデータベースへ保存する機能を提供する。

基本方針：

- Goで一般的かつ自然な設計を採用する。
- Goの標準ライブラリを優先して利用する。
- SQLiteドライバにはPure Goで利用できるライブラリを採用する。
- Goで実現可能な範囲で、安定性・保守性・安全性を重視する。
- 不必要に複雑な設計や過剰な最適化を行わない。
- エラーを適切に呼び出し元へ返す。
- 複数のgoroutineから利用される可能性を考慮し、データ競合を防止する。

この仕様書に記載されていない機能を独自判断で追加しないこと。

---

## 2. 開発環境・依存ライブラリ

### 2.1. 言語

- Go
- 安定版のGoを使用する。
- Go Modulesを利用する。

### 2.2. SQLiteドライバ

SQLiteドライバには `modernc.org/sqlite` を採用する。

理由：

- Pure Goで利用できる。
- CGOを必要としない。
- `database/sql` の標準インターフェースを利用できる。
- SQLiteの基本機能を利用できる。

SQLiteへのアクセスにはGo標準の `database/sql` を使用する。

### 2.3. その他のライブラリ

原則としてGo標準ライブラリを利用する。

主に以下を使用する。

- `database/sql`
- `encoding/json`
- `errors`
- `fmt`
- `os`
- `path/filepath`
- `sync`
- `time`

不要な外部ライブラリを追加しないこと。

---

## 3. ログレベル

### 3.1. SqlLogLevel

ログレベルを表す型を定義する。

```go
type SqlLogLevel int
```

以下の定数を定義する。

```go
const (
    SqlLogDebug    SqlLogLevel = 1
    SqlLogInfo     SqlLogLevel = 2
    SqlLogWarning  SqlLogLevel = 3
    SqlLogError    SqlLogLevel = 4
    SqlLogCritical SqlLogLevel = 5
)
```

ログレベルの優先順位は以下とする。

| レベル | 数値 | 優先順位 |
|---|---:|---|
| DEBUG | 1 | 低 |
| INFO | 2 | ↓ |
| WARNING | 3 | ↓ |
| ERROR | 4 | ↓ |
| CRITICAL | 5 | 高 |

**数値が大きいほど重要度が高い。**

ログレベルの数値は明示的に指定し、`iota` による暗黙的な採番には依存しないこと。

### 3.2. String()

`String()` メソッドを実装する。

```go
func (l SqlLogLevel) String() string
```

各ログレベルに対応する文字列を返す。

| 値 | 戻り値 |
|---|---|
| SqlLogDebug | DEBUG |
| SqlLogInfo | INFO |
| SqlLogWarning | WARNING |
| SqlLogError | ERROR |
| SqlLogCritical | CRITICAL |

定義されていない値が指定された場合もpanicを発生させず、`UNKNOWN` などの識別可能な文字列を返す。

### 3.3. JSON変換

`encoding/json` でログレベルを文字列として扱えるようにする。

以下のメソッドを実装する。

```go
func (l SqlLogLevel) MarshalJSON() ([]byte, error)

func (l *SqlLogLevel) UnmarshalJSON(data []byte) error
```

JSON上では以下のように表現する。

```json
"INFO"
```

不正なログレベル文字列が指定された場合はエラーを返す。

---

## 4. SqlLogRecode

ログレコードを表す構造体を定義する。

構造体名は `SqlLogRecode` とする。

```go
type SqlLogRecode struct {
    Timestamp float64     `json:"timestamp"`
    Level     SqlLogLevel `json:"level"`
    Name      string      `json:"name"`
    Msg       string      `json:"msg"`
    PID       int         `json:"pid"`
}
```

### 4.1. 各フィールド

| フィールド | 型 | 説明 |
|---|---|---|
| Timestamp | float64 | Unix time（秒） |
| Level | SqlLogLevel | ログレベル |
| Name | string | ログの発生元 |
| Msg | string | ログメッセージ |
| PID | int | プロセスID |

`Timestamp` はUnix timeの秒数とする。

小数部分には秒未満の精度を保持する。

### 4.2. ログレコード生成関数

以下の関数をパッケージレベルに定義する。

```go
func Critical(name string, msg string) SqlLogRecode

func Error(name string, msg string) SqlLogRecode

func Warning(name string, msg string) SqlLogRecode

func Info(name string, msg string) SqlLogRecode

func Debug(name string, msg string) SqlLogRecode
```

各関数は以下を設定した `SqlLogRecode` を返す。

- `Timestamp`: 関数呼び出し時点のUnix time
- `Level`: 関数名に対応するログレベル
- `Name`: 引数 `name`
- `Msg`: 引数 `msg`
- `PID`: `os.Getpid()` で取得したプロセスID

### 4.3. String()

ログレコードを文字列として表示できるようにする。

```go
func (r SqlLogRecode) String() string
```

表示するタイムスタンプはローカル時刻に変換する。

タイムゾーン情報は表示しない。

フォーマットは以下とする。

```text
2006-01-02T15:04:05.000000
```

マイクロ秒まで表示する。

出力形式の例：

```text
2026-11-11T12:13:45.123456 [INFO] [example] message (PID:1234)
```

`String()` は表示用のメソッドとし、構造体内部の `Timestamp` を変更してはならない。

### 4.4. JSON変換

`encoding/json` を使用して `SqlLogRecode` をJSONへ変換できるようにする。

フィールド名は構造体タグに従い、以下の形式とする。

```json
{
  "timestamp": 1794370425.123456,
  "level": "INFO",
  "name": "example",
  "msg": "message",
  "pid": 1234
}
```

`SqlLogLevel` は前述のJSON変換メソッドを利用すること。

---

## 5. SqlLogger

SQLiteデータベースへのログ保存を管理する構造体を定義する。

```go
type SqlLogger struct {
    dbPath    string
    recodes   []SqlLogRecode
    realtime  bool
    logLevel  SqlLogLevel
    nowRowNum uint64
    maxRecNum uint64

    db *sql.DB

    mu sync.Mutex
}
```

各フィールドは外部から直接変更できないようにする。

### 5.1. フィールドの説明

| フィールド | 説明 |
|---|---|
| dbPath | SQLiteデータベースのファイルパス |
| recodes | DBへ保存する前のログレコード |
| realtime | ログ追加時に即時保存するか |
| logLevel | 記録対象とする最低ログレベル |
| nowRowNum | DB内の現在のレコード数 |
| maxRecNum | DB内の最大レコード数の基準 |
| db | SQLite接続 |
| mu | 排他制御用Mutex |

`db` は `database/sql` の `*sql.DB` を使用する。

`mu` は複数goroutineからの同時アクセスによるデータ競合を防ぐために使用する。

---

## 6. コンストラクタ

以下の2つを実装する。

### 6.1. New

```go
func New(dbPath string) (*SqlLogger, error)
```

デフォルト設定でSqlLoggerを生成する。

デフォルト値：

| 項目 | 値 |
|---|---|
| dbPath | 引数で指定 |
| realtime | false |
| logLevel | SqlLogInfo |
| nowRowNum | DB内の実際のレコード数 |
| maxRecNum | 10,000,000 |

### 6.2. NewWith

```go
func NewWith(
    dbPath string,
    realtime bool,
    logLevel SqlLogLevel,
    maxRecNum uint64,
) (*SqlLogger, error)
```

すべてのオプションを指定してSqlLoggerを生成する。

`logLevel` には DEBUG、INFO、WARNING、ERROR、CRITICAL のいずれかを指定する。

定義されていない値が指定された場合はエラーを返す。

### 6.3. 初期化処理

コンストラクタでは以下を実行する。

1. 引数の妥当性を確認する。
2. データベースファイルの親ディレクトリが存在しない場合は作成する。
3. SQLiteデータベースへ接続する。
4. SQLiteのテーブルが存在しない場合は作成する。
5. 必要なインデックスが存在しない場合は作成する。
6. DB内のレコード数を取得する。
7. `nowRowNum` に取得した件数を設定する。

データベースファイルが既に存在する場合、既存データを削除してはならない。

データベースのパスには相対パス・絶対パスの両方を使用できるようにする。

### 6.4. SQLite接続設定

`database/sql` の接続プールはSQLiteの利用に適した設定とする。

基本設定：

```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(0)
```

接続はSqlLoggerの利用期間中維持する。

不要なタイミングで接続を開閉しないこと。

---

## 7. SQLiteデータベース設計

### 7.1. テーブル

テーブル名は `sql_logs` とする。

以下のカラムを作成する。

| カラム | SQLite型 | 説明 |
|---|---|---|
| id | INTEGER | 主キー、自動採番 |
| timestamp | REAL | Unix time |
| log_level | TEXT | ログレベル |
| name | TEXT | ログ発生元 |
| pid | INTEGER | プロセスID |
| msg | TEXT | ログメッセージ |

`id` はDB内部でレコードを一意に識別するために使用する。

`SqlLogRecode` には含めない。

### 7.2. テーブル作成SQL

概念上、以下の構造とする。

```sql
CREATE TABLE IF NOT EXISTS sql_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp REAL NOT NULL,
    log_level TEXT NOT NULL,
    name TEXT NOT NULL,
    pid INTEGER NOT NULL,
    msg TEXT NOT NULL
);
```

### 7.3. インデックス

検索性能を考慮して以下のインデックスを作成する。

```sql
CREATE INDEX IF NOT EXISTS idx_sql_logs_timestamp
ON sql_logs(timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_sql_logs_log_level
ON sql_logs(log_level);

CREATE INDEX IF NOT EXISTS idx_sql_logs_name
ON sql_logs(name);

CREATE INDEX IF NOT EXISTS idx_sql_logs_log_level_timestamp
ON sql_logs(log_level, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_sql_logs_name_timestamp
ON sql_logs(name, timestamp DESC);
```

既存のインデックスがある場合は再作成しない。

---

## 8. ログレベルによるフィルタリング

SqlLoggerは `logLevel` によって記録するログレベルを制御する。

設定されたログレベルより重要度の低いログは記録しない。

**数値が大きいほど重要度が高いため、判定条件は以下とする。**

```go
recordLevel >= logger.logLevel
```

例：

| 設定 logLevel | DEBUG | INFO | WARNING | ERROR | CRITICAL |
|---|---|---|---|---|---|
| DEBUG (1) | 記録 | 記録 | 記録 | 記録 | 記録 |
| INFO (2) | 無視 | 記録 | 記録 | 記録 | 記録 |
| WARNING (3) | 無視 | 無視 | 記録 | 記録 | 記録 |
| ERROR (4) | 無視 | 無視 | 無視 | 記録 | 記録 |
| CRITICAL (5) | 無視 | 無視 | 無視 | 無視 | 記録 |

フィルタリングで無視されたログは、`recodes` に追加してはならない。

無視されたログについてDBへの書き込みを行ってはならない。

---

## 9. ログ追加メソッド

SqlLoggerに以下のメソッドを実装する。

```go
func (l *SqlLogger) Critical(msg string) error

func (l *SqlLogger) Error(msg string) error

func (l *SqlLogger) Warning(msg string) error

func (l *SqlLogger) Info(msg string) error

func (l *SqlLogger) Debug(msg string) error
```

### 9.1. 共通処理

各メソッドは以下の処理を行う。

1. SqlLoggerが利用可能な状態であることを確認する。
2. 設定された `logLevel` によって記録対象か判定する。
3. 対象外であれば何もせず `nil` を返す。
4. 対象であればログレコードを生成する。
5. `Name` は固定文字列 `"logger"` とする。
6. `Msg` には引数の文字列を設定する。
7. `recodes` に追加する。
8. `realtime == true` の場合はDBへの保存処理を直ちに実行する。

ログレコードの生成には、前述のログレコード生成関数を利用してもよい。

### 9.2. エラー処理

ログ追加時に発生したエラーは呼び出し元へ返す。

ただし、ログレベルによって無視された場合はエラーとしない。

### 9.3. realtime

`realtime` の設定によって保存タイミングを制御する。

| realtime | 動作 |
|---|---|
| false | メモリ上に蓄積し、明示的なOutまたはCloseで保存 |
| true | ログ追加後、直ちにDBへ保存 |

---

## 10. Out

メモリ上に蓄積したログレコードをDBへ保存する。

```go
func (l *SqlLogger) Out() error
```

### 10.1. 基本動作

1. `recodes` が空の場合は何もせず `nil` を返す。
2. DBへトランザクションを開始する。
3. INSERT用のSQLを準備する。
4. `recodes` 内のすべてのレコードをDBへ追加する。
5. すべて成功したらCOMMITする。
6. COMMIT成功後に `recodes` をクリアする。
7. `nowRowNum` を保存したレコード数だけ増加させる。

### 10.2. トランザクション

複数のログレコードを保存する場合、1回のトランザクションで処理する。

レコードごとにトランザクションを開始してはならない。

可能な限りPrepared Statementを再利用する。

### 10.3. エラー時の動作

以下の場合、トランザクションをROLLBACKする。

- トランザクション開始に失敗した場合
- SQLの準備に失敗した場合
- INSERTに失敗した場合
- COMMITに失敗した場合

失敗した場合は以下を守る。

- `recodes` をクリアしない。
- `nowRowNum` を変更しない。
- エラーを呼び出し元へ返す。

保存が成功していないレコードを失わないこと。

### 10.4. 排他制御

`Out()` はMutexによって排他制御する。

ログ追加処理から `Out()` を呼び出す場合のデッドロックを防ぐため、内部処理を分離する。

例：

```go
func (l *SqlLogger) Out() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    return l.outLocked()
}
```

`outLocked()` はMutexを取得せず、呼び出し元がロックを保持していることを前提とする。

---

## 11. GetRecs

最新のログレコードを取得する。

```go
func (l *SqlLogger) GetRecs() ([]SqlLogRecode, error)
```

### 11.1. 取得条件

最新の1000件を取得する。

SQLの並び順は以下とする。

```sql
ORDER BY timestamp DESC, id DESC
```

同一タイムスタンプのレコードがある場合は `id` の降順で並べる。

### 11.2. 戻り値

取得したログレコードの配列を返す。

取得結果が0件の場合は空のスライスを返す。

`nil` スライスではなく空のスライスを返すことを基本とする。

### 11.3. エラー処理

DBアクセスやレコードの読み込みに失敗した場合はエラーを返す。

---

## 12. GetRecsWithLimit

取得件数を指定してログレコードを取得する。

```go
func (l *SqlLogger) GetRecsWithLimit(
    limit uint64,
) ([]SqlLogRecode, error)
```

### 12.1. 取得条件

指定された件数を上限として取得する。

並び順は `GetRecs()` と同様にする。

```sql
ORDER BY timestamp DESC, id DESC
```

### 12.2. limit

`limit == 0` の場合は空のスライスを返す。

SQLに渡す際は、SQLiteが扱える整数範囲を考慮し、不正な値や変換不能な値を適切にエラーとして扱う。

---

## 13. GetRecsWith

条件を指定してログレコードを取得する。

```go
func (l *SqlLogger) GetRecsWith(
    limit uint64,
    name *string,
    logLevel *SqlLogLevel,
) ([]SqlLogRecode, error)
```

### 13.1. 引数

| 引数 | 説明 |
|---|---|
| limit | 最大取得件数 |
| name | ログ発生元による絞り込み。nilの場合は無効 |
| logLevel | ログレベルによる絞り込み。nilの場合は無効 |

### 13.2. 検索条件

各引数が指定された場合、対応する条件を追加する。

- `name != nil`: `name` が一致するレコードを取得する。
- `logLevel != nil`: `log_level` が一致するレコードを取得する。
- 両方指定された場合: AND条件で検索する。
- 両方nilの場合: 条件なしで取得する。

検索に指定する `logLevel` は完全一致とする。指定したレベル以上のログを取得する検索ではない。

SQLの並び順は以下とする。

```sql
ORDER BY timestamp DESC, id DESC
```

### 13.3. SQLの安全性

検索条件の値をSQL文字列に直接連結してはならない。

必ずプレースホルダーを使用する。

SQLインジェクションを防止すること。

### 13.4. 戻り値

取得したログレコードの配列とエラーを返す。

---

## 14. Reduce

DB内のログレコード数が設定上限を超えた場合に、古いログを削除する。

```go
func (l *SqlLogger) Reduce(days uint) error
```

### 14.1. 実行条件

以下の場合のみ削除処理を実行する。

```go
l.nowRowNum > l.maxRecNum
```

条件を満たさない場合は何もせず `nil` を返す。

### 14.2. 削除条件

引数 `days` で指定された日数より古いレコードを削除する。

現在のUnix timeを `now` とした場合、削除基準は以下とする。

```text
cutoff = now - days * 86400
```

以下の条件に該当するレコードを削除する。

```sql
timestamp < cutoff
```

### 14.3. 削除処理

削除処理はトランザクション内で実行する。

削除成功後、DB内のレコード数を再取得する。

```sql
SELECT COUNT(*) FROM sql_logs;
```

取得した件数を `nowRowNum` に設定する。

### 14.4. 注意事項

- Reduceは自動実行しない。
- ログ追加時やOut実行時に勝手にReduceを呼び出してはならない。
- 利用者が明示的に呼び出した場合のみ実行する。
- 削除に失敗した場合は `nowRowNum` を変更しない。
- DBの件数と `nowRowNum` が一致するようにする。

---

## 15. Initialize

DBを初期化する。

```go
func (l *SqlLogger) Initialize() error
```

### 15.1. 基本動作

以下を実行する。

1. メモリ上に蓄積されているログレコードを破棄する。
2. DB接続を閉じる。
3. DBファイルを削除する。
4. 必要に応じてSQLiteのWAL・SHM関連ファイルも削除する。
5. DBを再作成する。
6. テーブルを作成する。
7. インデックスを作成する。
8. DBへ再接続する。
9. `nowRowNum` を0にする。

### 15.2. 注意事項

InitializeはDB内のデータをすべて削除する破壊的な操作である。

この動作を勝手に変更しないこと。

DBファイル削除時にエラーが発生した場合は、呼び出し元へエラーを返す。

SQLite接続が残った状態でDBファイルを削除しないこと。

初期化処理はMutexで排他制御する。

---

## 16. Close

Goではオブジェクトの破棄タイミングが保証されないため、明示的な終了処理を提供する。

```go
func (l *SqlLogger) Close() error
```

`io.Closer` として利用できるようにする。

### 16.1. 基本動作

1. 未保存のログレコードを `Out()` で保存する。
2. 保存に成功した場合、DB接続を閉じる。
3. Close済みであることを記録する。

使用例：

```go
logger, err := New("logs/logs.db")
if err != nil {
    return err
}
defer logger.Close()
```

### 16.2. エラー処理

Outに失敗した場合はDB接続を閉じず、エラーを返す。

この場合、バッファを保持し、呼び出し元が再度Closeを試行できるようにする。

DB接続のCloseに失敗した場合もエラーを返す。

### 16.3. 複数回のClose

Closeは複数回呼び出されても安全に動作するようにする。

既に正常終了している場合は `nil` を返す。

### 16.4. ガベージコレクション

`runtime.SetFinalizer` による自動保存・自動終了処理は実装しない。

利用者は `Close()` を明示的に呼び出すこと。

---

## 17. エラー設計

Go標準の `error` を利用する。

必要に応じて以下のエラーを定義する。

```go
var (
    ErrLoggerClosed = errors.New("sql logger is closed")
    ErrInvalidLevel = errors.New("invalid log level")
)
```

エラーを返す際は、可能な限り `fmt.Errorf` と `%w` を利用して元のエラーを保持する。

例：

```go
return fmt.Errorf("failed to insert log record: %w", err)
```

呼び出し元が `errors.Is` でエラーを判定できるようにする。

---

## 18. 排他制御・並行処理

SqlLoggerは複数goroutineから利用される可能性を考慮する。

以下の操作はMutexで排他制御する。

- ログ追加
- Out
- GetRecs
- GetRecsWithLimit
- GetRecsWith
- Reduce
- Initialize
- Close

### 18.1. Mutex

`sync.Mutex` を使用する。

同一インスタンスに対して複数goroutineから操作された場合でも、データ競合や不整合が発生しないようにする。

### 18.2. デッドロック防止

Mutexを取得するメソッドから、同じMutexを取得する別のメソッドを直接呼び出してはならない。

必要に応じてロック済みの内部処理を別メソッドに分離する。

例：

- `Out()` はロックを取得して `outLocked()` を呼ぶ。
- ログ追加メソッドはロックを取得した状態で `outLocked()` を呼ぶ。

### 18.3. DB接続

`database/sql` の接続設定とMutexを組み合わせ、SqlLogger内部でのDB操作が安全に行われるようにする。

---

## 19. テスト仕様

実装完了後、以下のテストを実施する。

テストには一時ディレクトリを使用し、実際の利用者データを変更しないこと。

### 19.1. SqlLogLevel

- DEBUGが数値1である。
- INFOが数値2である。
- WARNINGが数値3である。
- ERRORが数値4である。
- CRITICALが数値5である。
- 各レベルのString()が正しい文字列を返す。
- 数値が大きいほど重要度が高い。
- 不正な値でpanicしない。
- JSON変換が正しく行われる。
- 不正なJSONレベル文字列でエラーになる。

### 19.2. SqlLogRecode

- 各生成関数が正しいレベルを設定する。
- TimestampがUnix timeとして生成される。
- NameとMsgが正しく設定される。
- PIDが現在のプロセスIDになる。
- String()が指定フォーマットで表示される。
- JSON変換が正しく行われる。

### 19.3. コンストラクタ

- Newで新規DBが作成される。
- NewWithでオプションが正しく設定される。
- 親ディレクトリが存在しない場合に作成される。
- 既存DBのデータが保持される。
- nowRowNumが実際のDB件数になる。
- 不正なlogLevelを指定するとエラーになる。
- DB初期化に失敗した場合にエラーを返す。

### 19.4. ログレベルフィルタリング

- DEBUG設定ですべてのログが記録される。
- INFO設定でDEBUGが無視される。
- WARNING設定でDEBUGとINFOが無視される。
- ERROR設定でDEBUG、INFO、WARNINGが無視される。
- CRITICAL設定でCRITICALのみ記録される。
- 設定レベル以上の重要度を持つログが記録される。
- 無視されたログがDBに保存されない。
- 無視されたログがバッファにも追加されない。

### 19.5. ログ追加・Out

- 各ログメソッドで正しいログレベルが設定される。
- Nameが `"logger"` になる。
- realtime=falseでは即時保存されない。
- realtime=trueでは即時保存される。
- Outで蓄積されたログが保存される。
- Out成功後にバッファがクリアされる。
- Out失敗時にバッファが保持される。
- Out失敗時にnowRowNumが変更されない。
- 複数レコードが1トランザクションで保存される。

### 19.6. 検索

- GetRecsで最新1000件を取得できる。
- GetRecsWithLimitで指定件数を取得できる。
- GetRecsWithでnameによる絞り込みができる。
- GetRecsWithでlogLevelによる完全一致検索ができる。
- nameとlogLevelのAND検索ができる。
- 検索結果が新しい順になる。
- 同一タイムスタンプの場合にid順で安定して取得できる。
- 0件の場合に空スライスを返す。
- SQLインジェクションが発生しない。

### 19.7. Reduce

- nowRowNumがmaxRecNum以下なら削除しない。
- nowRowNumがmaxRecNumを超えた場合に削除する。
- 指定日数より古いレコードだけ削除する。
- 削除後にnowRowNumが正しく更新される。
- 削除失敗時にnowRowNumが変更されない。
- 自動実行されない。

### 19.8. Initialize

- DB内のレコードがすべて削除される。
- テーブルが再作成される。
- インデックスが再作成される。
- nowRowNumが0になる。
- recodesがクリアされる。
- 再初期化後もログを保存できる。

### 19.9. Close

- Closeで未保存ログが保存される。
- 保存失敗時にDB接続が維持される。
- 保存失敗時にバッファが保持される。
- Close後の操作が適切なエラーになる。
- 複数回Closeしても問題が発生しない。

### 19.10. 並行処理

複数goroutineから同時にログ追加・検索・保存を実行してもデータ競合が発生しないこと。

以下を実行する。

```bash
go test -race ./...
```

---

## 20. 実装手順

以下の順番で実装する。

1. Go Modulesとパッケージ構成を作成する。
2. SQLiteドライバを追加する。
3. SqlLogLevelを実装する。
4. SqlLogRecodeを実装する。
5. SQLiteの接続・テーブル作成処理を実装する。
6. SqlLoggerのコンストラクタを実装する。
7. ログレベルフィルタリングを実装する。
8. ログ追加メソッドを実装する。
9. Outを実装する。
10. 検索メソッドを実装する。
11. Reduceを実装する。
12. Initializeを実装する。
13. Closeを実装する。
14. 単体テストを作成する。
15. 並行処理テストを作成する。
16. 全テストを実行し、不具合を修正する。

各段階で必要なテストを実行し、問題がないことを確認してから次へ進む。

---

## 21. 品質確認

以下のコマンドを実行し、すべて成功すること。

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

### 21.1. コード品質

- Goの一般的な命名規則に従う。
- 公開する型・関数・メソッドには適切なコメントを付ける。
- 不要な抽象化を行わない。
- 同じ処理を不必要に重複させない。
- エラーを握りつぶさない。
- リソースを適切に解放する。
- Mutexのロック・アンロックを適切に管理する。
- SQLの値には必ずプレースホルダーを使用する。

---

## 22. 禁止事項

以下の実装は禁止する。

- ログレベルによるフィルタリングを省略すること。
- ログレベルの数値と重要度の関係を逆転させること。
- Outの失敗時にバッファを破棄すること。
- Outの失敗時にnowRowNumを更新すること。
- Reduceを自動実行すること。
- SQLへ値を直接文字列連結すること。
- 不必要にDB接続を開閉すること。
- Mutexの二重ロックによるデッドロックを発生させること。
- Closeの代わりにruntime.SetFinalizerを使用すること。
- エラーを無条件に握りつぶすこと。
- 仕様にない機能を独自判断で追加すること。
- Goで不要に複雑な実装を行うこと。

---

## 23. 完了条件

以下をすべて満たした時点で実装完了とする。

- SqlLogLevelがDEBUG=1、INFO=2、WARNING=3、ERROR=4、CRITICAL=5で実装されている。
- SqlLogRecodeが実装されている。
- SqlLoggerが実装されている。
- SQLiteへの保存・検索が正常に動作する。
- ログレベルによるフィルタリングが正常に動作する。
- realtimeの動作が正常である。
- Outのトランザクション処理が正常である。
- Reduceが仕様どおり動作する。
- Initializeが正常に動作する。
- Closeが正常に動作する。
- エラー処理が適切に実装されている。
- 並行処理でデータ競合が発生しない。
- すべてのテストが成功する。
- `go vet ./...` が成功する。
- `go test -race ./...` が成功する。

実装完了後は、作成・変更したファイル、実装した機能、実行したテストとその結果を報告すること。
