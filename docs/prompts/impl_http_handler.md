# Codex Implementation Specification: AsyncHttpClient (Go)

## 1. プロジェクト概要

Go 言語を使用して、簡易的な HTTP Client ライブラリを実装してください。

本ライブラリは、JSON を使用する REST API との通信を主な用途とする。

HTTP 通信には Go 標準ライブラリの `net/http` を使用する。

GET / POST / PUT / DELETE の4種類の HTTP メソッドを提供し、API 通信に必要な共通設定を Client インスタンスに保持する。

レスポンスボディは Go のジェネリック機能を使用して、任意の構造体に変換できるようにする。

### 1.1 基本方針

- Go で実装する。
- HTTP 通信には標準ライブラリ `net/http` を使用する。
- JSON のシリアライズ・デシリアライズには `encoding/json` を使用する。
- 非同期処理は Go 標準の goroutine / context の仕組みを利用する。
- HTTP Client はインスタンスとして保持し、複数のリクエストで再利用する。
- レスポンスはジェネリック型を使用して任意の構造体に変換する。
- エラーは独自の `HttpClientError` に統一する。
- HTTP リトライ機能は実装しない。
- 必要以上の機能を追加しない。

## 2. 重要な設計方針

### 2.1 URL の扱い

すべての HTTP メソッドで、完全な URL を引数として受け取る。

例：

- `https://api.example.com/users`
- `https://api.example.com/users/1`
- `https://api.example.com/users?page=1&limit=10`

Client インスタンスには Base URL を保持しない。

URL の組み立て機能は HTTP Client の実装範囲外とする。

以下の機能は実装しない。

- Base URL の保持
- パスの結合
- クエリパラメータの組み立て
- URL の自動補完

HTTP Client は呼び出し側から渡された完全な URL を使用して通信する。

### 2.2 Client インスタンス

HTTP 通信に必要な共通設定は `HttpClient` インスタンスに保持する。

保持する設定：

- `http.Client`
- Authorization ヘッダー
- 任意の共通 HTTP ヘッダー
- タイムアウト設定

一度生成した Client インスタンスを複数の API 通信で再利用できるようにする。

### 2.3 HTTP メソッド

以下の4つを独立した関数として実装する。

- GET
- POST
- PUT
- DELETE

### 2.4 Go のジェネリック仕様に関する重要事項

Go では、通常の構造体メソッドに独自の型パラメータを定義できない。

そのため、以下の形式で実装する。

```go
user, err := Get[User](ctx, client, url)
```

HTTP 通信の共通設定は `HttpClient` インスタンスに保持する。

HTTP 通信関数はジェネリックなパッケージ関数として実装する。

以下のようなメソッド形式は実装しない。

```go
// Go ではこのようなジェネリックメソッドは定義できない
func (c *HttpClient) Get[T any](ctx context.Context, url string) (T, error)
```

この制約を考慮して、公開 API を設計すること。

## 3. 使用するライブラリ

Go の標準ライブラリのみを使用する。

| Package | 用途 |
|---|---|
| net/http | HTTP 通信 |
| context | リクエストのキャンセル・タイムアウト管理 |
| encoding/json | JSON シリアライズ・デシリアライズ |
| errors | エラー判定 |
| fmt | エラーメッセージ |
| time | タイムアウト設定 |
| io | レスポンスボディの読み込み |

原則として外部ライブラリを追加しない。

## 4. HttpClient の仕様

### 4.1 構造体

以下の構成を基本とする。

```go
type HttpClient struct {
    client  *http.Client
    headers http.Header
}
```

必要に応じて非公開フィールドを追加してもよい。

ただし、不要なフィールドは追加しないこと。

### 4.2 コンストラクタ

以下のコンストラクタを実装する。

#### New

```go
func New() (*HttpClient, error)
```

デフォルト設定で Client を生成する。

デフォルト設定：

- タイムアウト：30秒
- Authorization：未設定
- その他の共通ヘッダー：未設定

#### NewWithTimeout

```go
func NewWithTimeout(timeout time.Duration) (*HttpClient, error)
```

指定されたタイムアウトで Client を生成する。

### 4.3 NewWithOptions

共通ヘッダーなどをまとめて設定できるコンストラクタを実装する。

```go
func NewWithOptions(
    timeout time.Duration,
    headers http.Header,
) (*HttpClient, error)
```

仕様：

- `timeout` に指定された時間をタイムアウトとして使用する。
- `headers` に指定された HTTP ヘッダーを共通ヘッダーとして保持する。
- Authorization ヘッダーを含め、任意の HTTP ヘッダーを設定できる。
- `headers` が nil の場合はヘッダー未設定として扱う。
- 渡された HeaderMap はコピーして保持する。
- 呼び出し側が元の HeaderMap を変更しても Client の設定に影響しないようにする。

タイムアウトは正の値を指定すること。

0以下の場合は設定エラーを返す。

### 4.4 Authorization ヘッダー

Authorization ヘッダーに認証情報を設定できるようにする。

Authorization は必須ではない。

以下の形式に対応する。

```http
Authorization: Bearer YOUR_TOKEN
```

ただし、Bearer 形式に限定しない。

呼び出し側が Authorization ヘッダーの値を自由に指定できるようにする。

例：

```go
headers := make(http.Header)

headers.Set(
    "Authorization",
    "Bearer YOUR_TOKEN",
)
```

Authorization 以外の任意の HTTP ヘッダーも設定できるようにする。

例：

```go
headers.Set("Accept", "application/json")
headers.Set("X-API-Version", "1")
```

### 4.5 Client の再利用

`http.Client` は Client インスタンス生成時に一度だけ生成する。

各 HTTP リクエストで `http.Client` を再生成してはならない。

Go の `http.Client` は複数の goroutine から安全に利用できるため、同じ `HttpClient` インスタンスを並行利用できる設計とする。

共通ヘッダーは生成後に変更しない。

## 5. HTTP 通信 API

### 5.1 共通仕様

すべての HTTP 通信関数はジェネリック関数として実装する。

共通の戻り値は以下とする。

```go
(T, error)
```

成功時は指定された型 `T` に変換した値を返す。

失敗時はゼロ値の `T` と `HttpClientError` を返す。

すべての関数で以下を受け取る。

- `ctx context.Context`
- `client *HttpClient`
- `url string`

POST / PUT では追加で JSON リクエストボディを受け取る。

### 5.2 GET

```go
func Get[T any](
    ctx context.Context,
    client *HttpClient,
    url string,
) (T, error)
```

仕様：

1. 指定された URL に GET リクエストを送信する。
2. `ctx` をリクエストに設定する。
3. Client に設定された共通ヘッダーを付与する。
4. タイムアウト設定を適用する。
5. HTTP ステータスコードを確認する。
6. 2xx 以外の場合は `HttpClientError` を返す。
7. 2xx の場合はレスポンスボディを JSON として解析する。
8. 指定されたジェネリック型 `T` に変換する。

GET のクエリパラメータは URL に含めて渡す。

GET 専用のクエリパラメータ構造体は実装しない。

### 5.3 POST

```go
func Post[T any](
    ctx context.Context,
    client *HttpClient,
    url string,
    body any,
) (T, error)
```

仕様：

1. 指定された URL に POST リクエストを送信する。
2. リクエストボディを JSON にシリアライズする。
3. Content-Type に `application/json` を設定する。
4. Client に設定された共通ヘッダーを付与する。
5. `ctx` をリクエストに設定する。
6. タイムアウト設定を適用する。
7. HTTP ステータスコードを確認する。
8. 2xx 以外の場合は `HttpClientError` を返す。
9. 2xx の場合はレスポンスボディを JSON として解析する。
10. 指定されたジェネリック型 `T` に変換する。

リクエストボディは `encoding/json` でシリアライズ可能な値を受け取る。

### 5.4 PUT

```go
func Put[T any](
    ctx context.Context,
    client *HttpClient,
    url string,
    body any,
) (T, error)
```

仕様：

POST と同じ JSON 処理を行い、HTTP メソッドを PUT とする。

### 5.5 DELETE

```go
func Delete[T any](
    ctx context.Context,
    client *HttpClient,
    url string,
) (T, error)
```

仕様：

1. 指定された URL に DELETE リクエストを送信する。
2. Client に設定された共通ヘッダーを付与する。
3. `ctx` をリクエストに設定する。
4. タイムアウト設定を適用する。
5. HTTP ステータスコードを確認する。
6. 2xx 以外の場合は `HttpClientError` を返す。
7. 2xx の場合はレスポンスボディを JSON として解析する。
8. 指定されたジェネリック型 `T` に変換する。

DELETE のリクエストボディは実装しない。

## 6. 共通 HTTP リクエスト処理

GET / POST / PUT / DELETE で共通する処理は、内部関数にまとめる。

以下の処理を共通化する。

1. HTTP リクエストの生成
2. Context の設定
3. 共通ヘッダーの付与
4. HTTP リクエストの送信
5. HTTP ステータスコードの確認
6. エラー時のレスポンスボディ取得
7. 成功時のレスポンスボディ取得
8. タイムアウト・通信エラーの分類

内部関数は非公開とする。

各 HTTP メソッドで同じ処理を重複して実装しないこと。

共通処理を使用しても、公開される GET / POST / PUT / DELETE は独立した関数として維持すること。

### 6.1 リクエストヘッダー

すべてのリクエストに Client の共通ヘッダーを適用する。

POST / PUT では以下を設定する。

```http
Content-Type: application/json
```

JSON リクエストの場合は Content-Type を application/json とする。

共通ヘッダーに Content-Type が指定されていても、JSON リクエストでは application/json を優先する。

### 6.2 HTTP ステータスコード

HTTP ステータスコードが 200〜299 の場合のみ成功とする。

それ以外は HttpStatus エラーとする。

### 6.3 レスポンスボディ

成功時のレスポンスボディを `encoding/json` でデシリアライズする。

レスポンスボディを読み取った後は必ず `Close()` すること。

### 6.4 HTTP エラー時のレスポンスボディ

HTTP ステータスコードが2xx以外の場合、レスポンスボディを文字列として取得する。

取得した文字列を `HttpClientError` に保存する。

レスポンスボディの取得に失敗した場合でも、HTTP ステータスコードの情報は失わないこと。

## 7. HttpClientError の仕様

エラーの種類を呼び出し側で判別できるようにする。

独自の `HttpClientError` 構造体とエラー種別を定義する。

### 7.1 エラー種別

以下のエラー種別を定義する。

```go
type HttpClientErrorType int

const (
    ErrorTypeHttpStatus HttpClientErrorType = iota
    ErrorTypeConnection
    ErrorTypeTimeout
    ErrorTypeJsonSerialization
    ErrorTypeJsonDeserialization
    ErrorTypeRequest
    ErrorTypeConfiguration
)
```

各エラー種別の意味：

| エラー種別 | 説明 |
|---|---|
| ErrorTypeHttpStatus | HTTP ステータスコードが2xx以外 |
| ErrorTypeConnection | 接続関連のエラー |
| ErrorTypeTimeout | タイムアウト |
| ErrorTypeJsonSerialization | リクエスト JSON の変換失敗 |
| ErrorTypeJsonDeserialization | レスポンス JSON の変換失敗 |
| ErrorTypeRequest | その他の HTTP リクエストエラー |
| ErrorTypeConfiguration | Client の設定エラー |

### 7.2 HttpClientError 構造体

以下を基本構成とする。

```go
type HttpClientError struct {
    Type       HttpClientErrorType
    StatusCode int
    Body       string
    Err        error
}
```

各フィールドの用途：

- `Type`: エラー種別
- `StatusCode`: HTTP ステータスコード
- `Body`: HTTP エラー時のレスポンスボディ
- `Err`: 元のエラー

HTTP ステータスエラー以外では `StatusCode` を0とする。

HTTP ステータスエラー以外では `Body` を空文字列とする。

### 7.3 Error メソッド

`error` インターフェースを実装する。

```go
func (e *HttpClientError) Error() string
```

エラー種別に応じた分かりやすいエラーメッセージを返す。

HTTP ステータスエラーではステータスコードを含める。

元のエラー情報がある場合は、その内容も含める。

### 7.4 Unwrap メソッド

元のエラーを取得できるようにする。

```go
func (e *HttpClientError) Unwrap() error
```

`errors.Is()` および `errors.As()` で元のエラーを判定できるようにする。

### 7.5 エラー判定関数

`HttpClientError` に以下のメソッドを実装する。

```go
func (e *HttpClientError) IsTimeout() bool

func (e *HttpClientError) IsConnectionError() bool

func (e *HttpClientError) IsJsonError() bool

func (e *HttpClientError) IsHttpError() bool

func (e *HttpClientError) IsAuthError() bool

func (e *HttpClientError) RequiresReauthentication() bool

func (e *HttpClientError) StatusCodeValue() int
```

各関数の仕様：

#### IsTimeout

ErrorTypeTimeout の場合に true を返す。

#### IsConnectionError

ErrorTypeConnection の場合に true を返す。

#### IsJsonError

以下の場合に true を返す。

- ErrorTypeJsonSerialization
- ErrorTypeJsonDeserialization

#### IsHttpError

ErrorTypeHttpStatus の場合に true を返す。

#### IsAuthError

認証・認可に関連する HTTP ステータスの場合に true を返す。

対象：

- 401 Unauthorized
- 403 Forbidden

#### RequiresReauthentication

再認証が必要な HTTP ステータスの場合に true を返す。

対象：

- 401 Unauthorized

403 Forbidden は権限不足として扱い、再認証が必要とは判定しない。

#### StatusCodeValue

ErrorTypeHttpStatus の場合は HTTP ステータスコードを返す。

それ以外の場合は0を返す。

### 7.6 エラー分類

HTTP 通信時のエラーは、以下のルールで分類する。

#### タイムアウト

以下を使用して判定する。

```go
var netErr net.Error

if errors.As(err, &netErr) && netErr.Timeout() {
    // Timeout
}
```

Context の DeadlineExceeded もタイムアウトとして扱う。

#### 接続エラー

接続先への接続に失敗した場合は Connection とする。

例：

- 接続先サーバーが停止している
- TCP 接続に失敗した
- DNS 解決に失敗した

`net.OpError` などのエラー情報を利用して判定する。

タイムアウトを Connection として分類しないこと。

#### JSON シリアライズエラー

リクエストボディを JSON に変換できない場合は JsonSerialization とする。

例：

- JSON に変換できない値
- 循環参照
- サポートされていない値

#### JSON デシリアライズエラー

レスポンス JSON を指定された型に変換できない場合は JsonDeserialization とする。

#### その他のリクエストエラー

上記に該当しない HTTP 通信関連のエラーは Request とする。

### 7.7 エラー分類の注意事項

エラー分類は内部の共通処理に集約する。

各 HTTP メソッドで個別に異なるエラー分類を実装しないこと。

HTTP ステータスエラーと通信エラーを混同しないこと。

エラーの元情報を失わないようにすること。

## 8. レスポンス処理

### 8.1 ジェネリック型

レスポンス JSON を任意の型に変換できるようにする。

Go のジェネリック機能を使用する。

例：

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

user, err := Get[User](
    ctx,
    client,
    "https://api.example.com/users/1",
)
```

複数のデータを含むレスポンスにも対応する。

```go
users, err := Get[[]User](
    ctx,
    client,
    "https://api.example.com/users",
)
```

### 8.2 JSON デシリアライズ

`json.Unmarshal()` を使用する。

JSON の変換に失敗した場合は JsonDeserialization エラーを返す。

### 8.3 空のレスポンスボディ

HTTP 204 No Content など、成功ステータスでレスポンスボディが空の場合がある。

今回の基本 API は JSON レスポンスを返す API を対象とする。

そのため、空のレスポンスボディは JsonDeserialization エラーとして扱う。

呼び出し側が空レスポンスを扱う必要がある場合に備え、将来的に別の API を追加できる設計にする。

現段階では空レスポンス専用 API を実装しない。

## 9. リトライ処理

リトライ処理は一切実装しない。

以下の場合も再試行しない。

- 接続エラー
- タイムアウト
- HTTP ステータスエラー
- JSON エラー

エラーが発生した場合は、そのまま呼び出し側に返す。

## 10. タイムアウト

### 10.1 デフォルト

デフォルトタイムアウトは30秒とする。

### 10.2 カスタマイズ

`NewWithTimeout()` または `NewWithOptions()` で変更できるようにする。

`http.Client.Timeout` を使用する。

タイムアウトにはリクエスト送信からレスポンスボディ読み込みまでの時間を含める。

Context によるキャンセル・タイムアウトも尊重する。

タイムアウトが発生した場合は ErrorTypeTimeout として返す。

## 11. Context の扱い

すべての HTTP 通信関数で `context.Context` を受け取る。

以下のように使用する。

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    10*time.Second,
)

defer cancel()

user, err := Get[User](
    ctx,
    client,
    "https://api.example.com/users/1",
)
```

仕様：

- リクエストには `http.NewRequestWithContext()` を使用する。
- 呼び出し側によるキャンセルを尊重する。
- Context の Deadline を尊重する。
- Context のキャンセルによって通信を停止できるようにする。
- Context の DeadlineExceeded は Timeout として分類する。
- 明示的なキャンセルは Request として分類する。

## 12. 共通ヘッダー

Client 生成時に任意の HTTP ヘッダーを設定できるようにする。

例：

```go
headers := make(http.Header)

headers.Set(
    "Authorization",
    "Bearer YOUR_TOKEN",
)

headers.Set(
    "Accept",
    "application/json",
)

client, err := NewWithOptions(
    30*time.Second,
    headers,
)
```

共通ヘッダーはすべての HTTP リクエストに適用する。

POST / PUT の JSON リクエストでは Content-Type を application/json とする。

ヘッダーの設定はリクエストごとにコピーして適用する。

Client が保持する共通ヘッダーを直接変更しないこと。

これにより複数 goroutine から同じ Client を安全に利用できるようにする。

## 13. 使用例

### 13.1 Client の生成

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func main() {
    headers := make(http.Header)

    headers.Set(
        "Authorization",
        "Bearer YOUR_TOKEN",
    )

    headers.Set(
        "Accept",
        "application/json",
    )

    client, err := NewWithOptions(
        30*time.Second,
        headers,
    )

    if err != nil {
        panic(err)
    }

    _ = client
    _ = context.Background()
    _ = fmt.Println
}
```

### 13.2 GET

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

user, err := Get[User](
    context.Background(),
    client,
    "https://api.example.com/users/1",
)

if err != nil {
    // エラー処理
    return
}

fmt.Println(user.Name)
```

### 13.3 POST

```go
type CreateUser struct {
    Name string `json:"name"`
}

type CreateUserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

body := CreateUser{
    Name: "test",
}

response, err := Post[CreateUserResponse](
    context.Background(),
    client,
    "https://api.example.com/users",
    body,
)

if err != nil {
    // エラー処理
    return
}

fmt.Println(response.ID)
```

### 13.4 PUT

```go
response, err := Put[CreateUserResponse](
    context.Background(),
    client,
    "https://api.example.com/users/1",
    body,
)
```

### 13.5 DELETE

```go
response, err := Delete[map[string]any](
    context.Background(),
    client,
    "https://api.example.com/users/1",
)
```

### 13.6 エラー判定

```go
user, err := Get[User](
    context.Background(),
    client,
    "https://api.example.com/users/1",
)

if err != nil {
    var httpErr *HttpClientError

    if errors.As(err, &httpErr) {

        if httpErr.RequiresReauthentication() {
            fmt.Println("再認証が必要です")

        } else if httpErr.IsTimeout() {
            fmt.Println("通信がタイムアウトしました")

        } else if httpErr.IsConnectionError() {
            fmt.Println("接続エラーが発生しました")

        } else if httpErr.IsJsonError() {
            fmt.Println("JSON の処理に失敗しました")

        } else if httpErr.IsHttpError() {
            fmt.Println(
                "HTTP エラー:",
                httpErr.StatusCodeValue(),
            )

        } else {
            fmt.Println("その他のエラー:", err)
        }
    }

    return
}

fmt.Println(user.Name)
```

## 14. テスト仕様

外部 API に依存しないテストを実装する。

Go 標準ライブラリの `net/http/httptest` を使用する。

### 14.1 Client 生成

以下をテストする。

- デフォルト設定で生成できること。
- デフォルトタイムアウトが30秒であること。
- タイムアウトを指定して生成できること。
- 共通ヘッダーを指定して生成できること。
- Authorization ヘッダーを指定できること。
- Authorization 未設定でも生成できること。
- nil の HeaderMap でも生成できること。
- 0以下のタイムアウトを指定した場合にエラーになること。
- 元の HeaderMap を変更しても Client の設定が変わらないこと。

### 14.2 GET

以下をテストする。

- GET リクエストが送信されること。
- HTTP メソッドが GET であること。
- JSON レスポンスが指定した構造体に変換されること。
- 共通ヘッダーが送信されること。
- クエリパラメータ付き URL を使用できること。
- 2xx 以外の場合に HttpStatus が返ること。

### 14.3 POST

以下をテストする。

- POST リクエストが送信されること。
- HTTP メソッドが POST であること。
- リクエストボディが正しく JSON に変換されること。
- Content-Type が application/json であること。
- 共通ヘッダーが送信されること。
- JSON レスポンスが指定した構造体に変換されること。
- JSON シリアライズエラーが正しく分類されること。
- 2xx 以外の場合に HttpStatus が返ること。

### 14.4 PUT

以下をテストする。

- PUT リクエストが送信されること。
- HTTP メソッドが PUT であること。
- リクエストボディが正しく JSON に変換されること。
- JSON レスポンスが指定した構造体に変換されること。
- 2xx 以外の場合に HttpStatus が返ること。

### 14.5 DELETE

以下をテストする。

- DELETE リクエストが送信されること。
- HTTP メソッドが DELETE であること。
- JSON レスポンスが指定した構造体に変換されること。
- 2xx 以外の場合に HttpStatus が返ること。

### 14.6 エラー分類

以下をそれぞれテストする。

- HTTP 400
- HTTP 401
- HTTP 403
- HTTP 404
- HTTP 500
- 接続エラー
- タイムアウト
- Context キャンセル
- JSON シリアライズエラー
- JSON デシリアライズエラー

特に以下の判定が正しいこと。

| 状態 | IsAuthError | RequiresReauthentication |
|---|---|---|
| 401 | true | true |
| 403 | true | false |
| 404 | false | false |
| 500 | false | false |

### 14.7 並行利用

同一の HttpClient インスタンスを複数の goroutine から利用できることを確認する。

`go test -race ./...` でもデータ競合が発生しないこと。

## 15. 実装時の制約

以下を必ず守ること。

1. GET / POST / PUT / DELETE はそれぞれ独立した公開関数として実装する。
2. すべての HTTP メソッドで完全な URL を受け取る。
3. Base URL は保持しない。
4. Authorization は任意設定とする。
5. 任意の共通 HTTP ヘッダーを設定できるようにする。
6. レスポンスはジェネリック型で受け取る。
7. すべてのエラーを HttpClientError に統一する。
8. エラー種別を判定するメソッドを実装する。
9. HTTP 401 と HTTP 403 を区別する。
10. 通信エラーとタイムアウトを区別する。
11. JSON シリアライズとデシリアライズのエラーを区別する。
12. デフォルトタイムアウトは30秒とする。
13. Context によるキャンセルをサポートする。
14. リトライ機能を実装しない。
15. URL 生成機能を実装しない。
16. HTTP Client をリクエストごとに生成しない。
17. 同一 Client インスタンスの並行利用を安全にする。
18. 指定されていない機能を独断で追加しない。
19. 不必要な外部ライブラリを追加しない。
20. Go の言語仕様に適合した公開 API にする。

## 16. 実装手順

以下の順番で実装すること。

1. Go Module を初期化する。
2. HttpClientErrorType を実装する。
3. HttpClientError を実装する。
4. HttpClient の構造体とコンストラクタを実装する。
5. 共通ヘッダーとタイムアウトの設定を実装する。
6. 内部共通 HTTP リクエスト処理を実装する。
7. GET を実装する。
8. POST を実装する。
9. PUT を実装する。
10. DELETE を実装する。
11. エラー判定メソッドを実装する。
12. 単体テスト・統合テストを実装する。
13. gofmt を実行する。
14. go vet ./... を実行する。
15. go test ./... を実行する。
16. go test -race ./... を実行する。
17. コンパイルエラー・警告・テスト失敗を修正する。

## 17. 完了条件

以下をすべて満たした場合に実装完了とする。

- GET / POST / PUT / DELETE が正常に動作する。
- 完全な URL を使用して通信できる。
- Authorization ヘッダーを設定できる。
- 任意の共通 HTTP ヘッダーを設定できる。
- JSON リクエストを送信できる。
- JSON レスポンスを任意の構造体に変換できる。
- HTTP ステータスエラーを判定できる。
- 401 と 403 を区別できる。
- 通信エラーを判定できる。
- JSON シリアライズエラーを判定できる。
- JSON デシリアライズエラーを判定できる。
- タイムアウトを判定できる。
- Context によるキャンセルが動作する。
- デフォルトタイムアウトが30秒である。
- 同一 Client インスタンスを並行利用できる。
- すべてのテストが成功する。
- gofmt が成功する。
- go vet ./... が成功する。
- go test -race ./... が成功する。

実装完了後は、以下を報告すること。

- 実装した機能
- 作成・変更したファイル
- 公開 API の使用方法
- テスト結果
- go vet の結果
- Race Detector の結果
- 未実装事項や制約

以上の仕様に従って実装してください。