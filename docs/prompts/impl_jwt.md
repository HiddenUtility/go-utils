# Codex Implementation Specification: JWT Utility Library (Go)

## 1. Project Overview

Go言語（Golang）でJWT（JSON Web Token）を生成・検証・読み込みするユーティリティライブラリを実装してください。

本ライブラリは、JWTの署名付きトークンの生成、秘密鍵の生成、Payloadのシリアライズ・デシリアライズ、および有効期限の検証機能を提供します。

### 1.1 基本仕様

| 項目 | 仕様 |
|---|---|
| 言語 | Go |
| JWT署名アルゴリズム | HS256 |
| UUID | UUID v4 |
| Payload形式 | JSON |
| シリアライズ | encoding/json |
| デシリアライズ | encoding/json |
| 時刻表現 | Unix Timestamp（秒） |
| 秘密鍵 | 暗号学的に安全な乱数から生成 |
| エラー処理 | error |
| 非同期処理 | 使用しない |

### 1.2 実装方針

以下の方針を厳守してください。

1. 仕様に記載されていない機能を独自に追加しないこと。
2. 不必要な外部パッケージを追加しないこと。
3. エラーの種類を明確に区別すること。
4. JWTの署名検証にはHS256のみを使用すること。
5. 秘密鍵の生成には暗号学的に安全な乱数生成器を使用すること。
6. UUIDはUUID v4を使用すること。
7. Payloadは標準の `encoding/json` によってシリアライズ・デシリアライズできること。
8. 既存のJWT標準仕様と互換性を持たせること。
9. Goの標準的なエラー処理方法に従うこと。
10. コンパイル可能で、テストが実行できる状態まで実装すること。

---

## 2. Dependencies

以下のパッケージを使用してください。

| Package | 用途 |
|---|---|
| github.com/golang-jwt/jwt/v5 | JWTの生成・署名・検証 |
| github.com/google/uuid | UUID v4の生成 |
| crypto/rand | 暗号学的に安全な乱数生成 |
| encoding/hex | バイナリデータの16進数変換 |
| encoding/json | JSONのシリアライズ・デシリアライズ |
| errors | エラー判定 |
| fmt | エラーメッセージ |
| time | 時刻の取得 |

`crypto/rand`、`encoding/hex`、`encoding/json`、`errors`、`fmt`、`time` はGo標準ライブラリを使用してください。

### 2.1 go.mod

必要な外部パッケージを定義してください。

- `github.com/golang-jwt/jwt/v5`
- `github.com/google/uuid`

バージョンは実装時点で利用可能な安定版を使用してください。

不要な外部パッケージを追加しないでください。

---

## 3. JWT Secret Key

JWTの署名・検証に使用する秘密鍵を生成する関数を実装してください。

### 3.1 Secret Key Generation

関数名：

`GenerateSecretKey`

シグネチャ：

```go
func GenerateSecretKey() (string, error)
```

仕様：

1. `crypto/rand` を使用する。
2. 32バイト（256ビット）以上のランダムデータを生成する。
3. 生成したバイナリデータを16進数文字列に変換する。
4. 生成した文字列を秘密鍵として返す。
5. 乱数生成に失敗した場合はエラーを返す。
6. 秘密鍵は毎回異なる値になること。

生成する秘密鍵の形式：

- 64文字以上の16進数文字列。
- 大文字・小文字は問わない。
- 暗号学的に安全な乱数から生成すること。

注意：

- UUIDを単純に秘密鍵として使用しないこと。
- 固定文字列や時刻などから秘密鍵を生成しないこと。
- 生成した秘密鍵をログに出力しないこと。
- `crypto/rand` のエラーを握りつぶさないこと。

---

## 4. JWT Payload

JWTのPayloadを表現する構造体を実装してください。

構造体名：

`JwtPayload`

### 4.1 Struct Definition

```go
package jwtutil

type JwtPayload struct {
    // 発行者
    Iss string `json:"iss,omitempty"`

    // 主体（ユーザーIDなど）
    Sub string `json:"sub,omitempty"`

    // 対象者
    Aud string `json:"aud,omitempty"`

    // 有効期限（Unix Timestamp 秒）
    Exp int64 `json:"exp"`

    // 有効開始日時（Unix Timestamp 秒）
    Nbf int64 `json:"nbf"`

    // 発行日時（Unix Timestamp 秒）
    Iat int64 `json:"iat"`

    // JWT ID（UUID v4）
    Jti string `json:"jti"`

    // アプリケーション独自のデータ
    Data any `json:"data,omitempty"`
}
```

### 4.2 Payload Attributes

JWT標準クレームとして以下を定義してください。

| フィールド | JSON名 | 意味 | 必須 |
|---|---|---|---|
| Iss | iss | 発行者（Issuer） | 任意 |
| Sub | sub | 主体（Subject） | 任意 |
| Aud | aud | 対象者（Audience） | 任意 |
| Exp | exp | 有効期限（Expiration Time） | 必須 |
| Nbf | nbf | 有効開始日時（Not Before） | 必須 |
| Iat | iat | 発行日時（Issued At） | 必須 |
| Jti | jti | JWT ID | 必須 |
| Data | data | アプリケーション独自データ | 任意 |

時刻はすべてUnix Timestampの秒単位で扱ってください。

`Exp`、`Nbf`、`Iat` は `int64` を使用してください。

### 4.3 jti

`Jti` はJWTを一意に識別するためのIDです。

以下の仕様にしてください。

- UUID v4を使用する。
- Payload生成時に自動生成する。
- 呼び出し側から指定する必要はない。
- 生成したUUIDを文字列として保持する。

UUIDの生成には `github.com/google/uuid` を使用してください。

UUID生成に失敗した場合はエラーを返してください。

### 4.4 data

`Data` はアプリケーション固有の任意のJSONデータを格納するためのフィールドです。

Goの `any` 型を使用してください。

以下のようなデータを格納できること。

```json
{
    "user_id": 123,
    "username": "test_user",
    "role": "admin",
    "permissions": ["read", "write"]
}
```

以下の条件を満たしてください。

1. 任意のJSONオブジェクトを格納できる。
2. JSON配列を格納できる。
3. 文字列、数値、真偽値なども格納できる。
4. `Data` が `nil` の場合はJSONからフィールド自体を省略する。
5. `Data` が `nil` でない場合は通常のJSONとして出力する。

`json:"data,omitempty"` を使用してください。

注意：

Goの `any` 型にJSONをデシリアライズした場合、数値は通常 `float64` として扱われます。この動作は許容します。

---

## 5. Payload Constructor

Payloadを生成するコンストラクタを実装してください。

関数名：

`NewJwtPayload`

シグネチャ：

```go
func NewJwtPayload(
    exp int64,
    nbf int64,
    data any,
) (*JwtPayload, error)
```

仕様：

1. `exp` を有効期限として設定する。
2. `nbf` を有効開始日時として設定する。
3. `iat` は現在時刻から自動生成する。
4. `jti` はUUID v4から自動生成する。
5. `iss` は空文字列に設定する。
6. `sub` は空文字列に設定する。
7. `aud` は空文字列に設定する。
8. `data` は引数の値を設定する。

現在時刻は以下の形式で取得してください。

```go
time.Now().Unix()
```

時刻はUnix Timestampの秒単位とします。

UUID生成に失敗した場合はエラーを返してください。

コンストラクタ内で `panic` を発生させないでください。

---

## 6. JWT Serialization

PayloadをJWT形式の文字列にシリアライズする関数を実装してください。

関数名：

`EncodeJWT`

シグネチャ：

```go
func EncodeJWT(
    payload *JwtPayload,
    secretKey string,
) (string, error)
```

### 6.1 処理内容

1. `payload` がnilでないことを確認する。
2. `secretKey` が空文字列でないことを確認する。
3. `JwtPayload` をJSONにシリアライズする。
4. JWT標準の形式に従ってトークンを生成する。
5. 署名アルゴリズムにはHS256を使用する。
6. `secretKey` を署名鍵として使用する。
7. 生成したJWT文字列を返す。
8. シリアライズや署名処理に失敗した場合はエラーを返す。

JWT Headerのアルゴリズムは必ずHS256にしてください。

以下のライブラリを使用してください。

`github.com/golang-jwt/jwt/v5`

JWTの生成には `jwt.NewWithClaims` を使用するか、同等の安全な方法を採用してください。

署名アルゴリズムは `jwt.SigningMethodHS256` を明示してください。

注意：

JWTの生成処理では有効期限の自動更新やPayloadの書き換えを行わないでください。

---

## 7. JWT Deserialization and Verification

JWT文字列からPayloadを読み込む関数を実装してください。

関数名：

`DecodeJWT`

シグネチャ：

```go
func DecodeJWT(
    token string,
    secretKey string,
) (*JwtPayload, error)
```

### 7.1 処理内容

以下の順序で処理してください。

1. 秘密鍵が空文字列でないことを確認する。
2. JWT文字列が空でないことを確認する。
3. JWTの形式を確認する。
4. JWT Headerの形式を確認する。
5. JWT Headerの署名アルゴリズムがHS256であることを確認する。
6. 秘密鍵を使用して署名を検証する。
7. 署名検証に成功した場合、Payloadをデシリアライズする。
8. Payloadの有効期限を検証する。
9. Payloadの有効開始日時を検証する。
10. すべての検証に成功した場合、`JwtPayload` を返す。

### 7.2 Signature Verification

署名の検証には `github.com/golang-jwt/jwt/v5` を使用してください。

以下の場合は署名検証エラーとして扱ってください。

- 秘密鍵が間違っている。
- JWTの署名が改ざんされている。
- 署名が正しくない。

秘密鍵の不一致と署名の改ざんは、どちらも `ErrInvalidSignature` として扱って構いません。

### 7.3 Algorithm Restriction

JWT Headerに指定されているアルゴリズムを無条件に信用してはいけません。

必ずHS256のみを許可してください。

`jwt.ParseWithClaims` のオプションに以下を指定してください。

```go
jwt.WithValidMethods([]string{
    jwt.SigningMethodHS256.Alg(),
})
```

これによりHS256以外のアルゴリズムを拒否してください。

### 7.4 Error Distinction

以下のエラーを区別できるようにしてください。

- JWTの形式が不正。
- JWT Headerのパース失敗。
- JWT Payloadのパース失敗。
- 署名検証失敗。
- 有効期限切れ。
- 有効開始日時前。
- 秘密鍵の不正。

特に、秘密鍵の間違いによる署名検証失敗とJSONのパース失敗を同じエラーにしないでください。

`github.com/golang-jwt/jwt/v5` のエラーを `errors.Is` などで判定し、独自のエラー型に変換してください。

JWT HeaderのパースとPayloadのパースは区別して処理してください。

必要に応じて、JWTの各セグメントを事前に形式確認する処理を実装して構いません。

ただし、署名検証に成功する前にPayloadの内容を信頼してはいけません。

---

## 8. Payload Validation

Payloadの有効期限および有効開始日時を検証する関数を実装してください。

関数名：

`ValidateTime`

シグネチャ：

```go
func ValidateTime(
    payload *JwtPayload,
) error
```

### 8.1 Expiration Validation

`Exp` は有効期限を表します。

以下の場合は有効期限切れと判定してください。

```go
now >= payload.Exp
```

### 8.2 Not Before Validation

`Nbf` は有効開始日時を表します。

以下の場合は有効開始日時前と判定してください。

```go
now < payload.Nbf
```

### 8.3 Valid Condition

以下の条件をすべて満たした場合に有効とします。

```go
payload.Nbf <= now && now < payload.Exp
```

### 8.4 Time Handling

- 現在時刻は `time.Now().Unix()` で取得する。
- 時刻の比較は秒単位で行う。
- タイムゾーンによる影響を受けないこと。
- 有効期限の判定に猶予時間（leeway）を設けないこと。
- `payload` がnilの場合はエラーを返す。

### 8.5 Validation Order

以下の順序で検証してください。

1. 有効期限（Exp）
2. 有効開始日時（Nbf）

両方に問題がある場合は、有効期限切れを優先して返してください。

### 8.6 decode_jwtとの連携

`DecodeJWT` の内部でも必ず有効期限と有効開始日時を検証してください。

個別の `ValidateTime` 関数も公開し、すでに取得済みのPayloadに対して再検証できるようにしてください。

JWTの読み込み時に検証を省略するオプションは設けないでください。

`jwt.ParseWithClaims` のオプションでは以下を明示してください。

- アルゴリズム：HS256のみ。
- 有効期限検証：有効。
- 有効開始日時検証：有効。
- Leeway：0。
- 有効期限（Exp）：必須。
- 有効開始日時（Nbf）：必須。

`jwt/v5` の `WithExpirationRequired` を使用してください。

`Nbf` の必須性については、独自のPayload検証でも確認してください。

また、ライブラリの標準検証に加え、独自の `ValidateTime` を使用してエラーの種類と判定条件を統一してください。

---

## 9. Error Handling

独自のエラー型を実装してください。

エラーの判定には `errors.Is` を利用できるようにしてください。

### 9.1 Error Definitions

以下のエラーを定義してください。

```go
package jwtutil

import "errors"

var (
    // 秘密鍵が不正、または空
    ErrInvalidKey = errors.New("invalid secret key")

    // JWTの形式が不正
    ErrInvalidFormat = errors.New("invalid jwt format")

    // JWT Headerのパース失敗
    ErrInvalidHeader = errors.New("invalid jwt header")

    // JWT Payloadのパース失敗
    ErrInvalidPayload = errors.New("invalid jwt payload")

    // 署名検証失敗
    ErrInvalidSignature = errors.New("invalid jwt signature")

    // 有効期限切れ
    ErrExpired = errors.New("jwt expired")

    // 有効開始日時前
    ErrNotYetValid = errors.New("jwt is not yet valid")

    // 乱数生成失敗
    ErrRandomGeneration = errors.New("random generation failed")

    // UUID生成失敗
    ErrUUIDGeneration = errors.New("uuid generation failed")

    // その他のJWT関連エラー
    ErrOther = errors.New("jwt processing failed")
)
```

実際の実装では、必要に応じて元のエラーを保持してください。

### 9.2 Error Mapping

`github.com/golang-jwt/jwt/v5` のエラーを以下のように変換してください。

| JWTライブラリ側のエラー | 独自エラー |
|---|---|
| ErrTokenSignatureInvalid | ErrInvalidSignature |
| ErrTokenExpired | ErrExpired |
| ErrTokenNotValidYet | ErrNotYetValid |
| JWT Headerのパース失敗 | ErrInvalidHeader |
| PayloadのJSONパース失敗 | ErrInvalidPayload |
| JWT形式の不正 | ErrInvalidFormat |
| その他のエラー | ErrOther |

注意：

JWTライブラリのエラー判定には `errors.Is` を使用してください。

### 9.3 Error Wrapping

元のエラー情報を失わないようにしてください。

例えば以下の形式を使用してください。

```go
return fmt.Errorf("%w: %v", ErrInvalidSignature, err)
```

呼び出し側では以下のようにエラーを判定できること。

```go
if errors.Is(err, ErrInvalidSignature) {
    // 署名検証失敗
}
```

秘密鍵が空の場合は `ErrInvalidKey` を返してください。

重要：

- 署名検証エラーを `ErrInvalidFormat` に変換しないこと。
- JSONパースエラーを `ErrInvalidSignature` に変換しないこと。
- エラーを握りつぶさないこと。
- `panic` を通常のエラー処理に使用しないこと。
- 秘密鍵そのものやJWTの秘密情報をエラーメッセージに含めないこと。

---

## 10. Public API

外部から使用する公開APIは以下を基本としてください。

```go
package jwtutil

func GenerateSecretKey() (string, error)

func NewJwtPayload(
    exp int64,
    nbf int64,
    data any,
) (*JwtPayload, error)

func EncodeJWT(
    payload *JwtPayload,
    secretKey string,
) (string, error)

func DecodeJWT(
    token string,
    secretKey string,
) (*JwtPayload, error)

func ValidateTime(
    payload *JwtPayload,
) error
```

また、以下も外部から利用できるように公開してください。

- `JwtPayload`
- `ErrInvalidKey`
- `ErrInvalidFormat`
- `ErrInvalidHeader`
- `ErrInvalidPayload`
- `ErrInvalidSignature`
- `ErrExpired`
- `ErrNotYetValid`
- `ErrRandomGeneration`
- `ErrUUIDGeneration`
- `ErrOther`

公開APIを必要以上に増やさないでください。

---

## 11. Module Structure

シンプルなライブラリ構成にしてください。

推奨構成：

```text
jwtutil/
├── go.mod
├── payload.go
├── error.go
├── jwt.go
└── jwt_test.go
```

各ファイルの責務：

| ファイル | 責務 |
|---|---|
| go.mod | Go Moduleと依存関係 |
| payload.go | JwtPayload構造体、コンストラクタ |
| error.go | エラー定義 |
| jwt.go | 秘密鍵生成、JWT生成・読み込み、有効期限検証 |
| jwt_test.go | 単体テスト・結合テスト |

この構成は推奨ですが、同等の責務分離ができる場合は多少変更しても構いません。

---

## 12. Unit Tests

以下のテストを必ず実装してください。

### 12.1 Secret Key Tests

1. 秘密鍵が正常に生成できる。
2. 生成した秘密鍵が64文字以上である。
3. 生成した秘密鍵が16進数文字列である。
4. 連続して生成した秘密鍵が異なる。

### 12.2 Payload Tests

1. Payloadを正常に生成できる。
2. `Iat` が現在時刻に設定されている。
3. `Jti` がUUID v4形式である。
4. `Iss`、`Sub`、`Aud` が初期状態で空文字列である。
5. `Data` に任意のJSONを設定できる。
6. PayloadをJSONにシリアライズできる。
7. JSONからPayloadにデシリアライズできる。
8. `Data` がnilの場合、JSONにdataフィールドが出力されない。
9. `Data` にJSONオブジェクトを設定できる。
10. `Data` にJSON配列を設定できる。
11. `Data` に文字列・数値・真偽値を設定できる。

### 12.3 JWT Encoding Tests

1. 正しい秘密鍵でJWTを生成できる。
2. 生成したJWTが3つのセグメントで構成されている。
3. JWT HeaderのアルゴリズムがHS256である。
4. Payloadの内容が正しく保持されている。
5. 空の秘密鍵ではエラーになる。
6. nilのPayloadではエラーになる。

### 12.4 JWT Decoding Tests

1. 正しい秘密鍵でJWTを読み込める。
2. 読み込んだPayloadが元のPayloadと一致する。
3. 間違った秘密鍵では `ErrInvalidSignature` になる。
4. 改ざんされたJWTでは `ErrInvalidSignature` になる。
5. 不正なJWT形式では `ErrInvalidFormat` または適切な形式エラーになる。
6. 不正なHeaderでは `ErrInvalidHeader` になる。
7. 不正なPayload JSONでは `ErrInvalidPayload` になる。
8. 空の秘密鍵では `ErrInvalidKey` になる。
9. HS256以外のアルゴリズムを使用したJWTは拒否される。
10. Expが存在しないJWTは拒否される。
11. Nbfが存在しないJWTは拒否される。

### 12.5 Time Validation Tests

1. 現在時刻が有効期間内なら成功する。
2. 有効期限を過ぎたPayloadでは `ErrExpired` になる。
3. 有効期限と現在時刻が一致した場合は `ErrExpired` になる。
4. 有効開始日時が未来の場合は `ErrNotYetValid` になる。
5. 有効開始日時と現在時刻が一致した場合は成功する。
6. 有効期限と有効開始日時の両方に問題がある場合は `ErrExpired` になる。
7. nilのPayloadではエラーになる。

### 12.6 Automatic Validation Tests

1. 有効なJWTは `DecodeJWT` で正常に読み込める。
2. 有効期限切れのJWTは `DecodeJWT` で `ErrExpired` になる。
3. 有効開始日時前のJWTは `DecodeJWT` で `ErrNotYetValid` になる。

時刻に依存するテストは、可能な限り固定時刻を使用する内部検証関数を用意して、テスト結果が不安定にならないようにしてください。

---

## 13. Security Requirements

JWTを扱うライブラリとして、以下のセキュリティ要件を守ってください。

1. 署名アルゴリズムはHS256に固定する。
2. 秘密鍵は暗号学的に安全な乱数から生成する。
3. 秘密鍵をソースコードにハードコードしない。
4. 秘密鍵をログに出力しない。
5. 署名検証を省略しない。
6. JWTの署名検証に成功する前にPayloadを信頼しない。
7. JWT Headerに指定されたアルゴリズムを無条件に信用しない。
8. HS256以外のアルゴリズムを許可しない。
9. 有効期限の検証を省略しない。
10. 有効開始日時の検証を省略しない。
11. 不正なJWTを受け取った場合でもpanicを発生させない。
12. JWTの検証に失敗した場合はPayloadを返さない。

---

## 14. Implementation Restrictions

以下の制約を守ってください。

- 非同期処理を導入しない。
- 不要なインターフェースや抽象化を追加しない。
- 不要な設定項目を追加しない。
- RSAなどの公開鍵方式を追加しない。
- HS384、HS512など他の署名アルゴリズムを追加しない。
- JWTのリフレッシュ機能を追加しない。
- JWTのブラックリスト機能を追加しない。
- JWTの保存機能を追加しない。
- 独自の暗号化処理を実装しない。
- Payloadの独自データには `any` 型を使用する。
- 時刻はUnix Timestampの秒単位で統一する。
- Go標準のエラー処理方法に従う。
- 不要な外部パッケージを追加しない。

---

## 15. Completion Criteria

実装完了前に以下を確認してください。

1. 仕様に記載されたすべての関数が実装されている。
2. すべての公開APIが外部パッケージから利用できる。
3. JWTの生成・読み込みが正常に動作する。
4. HS256による署名検証が正常に動作する。
5. UUID v4によるJtiの自動生成が動作する。
6. 任意のJSONデータをPayloadに格納できる。
7. 有効期限と有効開始日時を正しく検証できる。
8. 秘密鍵の不一致とJSONパースエラーを区別できる。
9. エラーがerror型で適切に返される。
10. すべてのテストが成功する。
11. `gofmt` が適用されている。
12. `go vet ./...` が成功する。
13. `go test ./...` が成功する。

---

## 16. Final Instructions for Codex

実装を開始する前に、既存プロジェクトの構成とgo.modを確認してください。

既存の実装や設定がある場合は、むやみに変更・削除しないでください。

仕様に矛盾や実装上の重大な問題がある場合は、実装を開始する前にユーザーへ確認してください。

特に問題がなければ、以下の順序で実装してください。

1. go.modの依存関係を確認・追加する。
2. エラー定義を実装する。
3. JwtPayloadを実装する。
4. 秘密鍵生成関数を実装する。
5. JWT生成関数を実装する。
6. JWT読み込み・署名検証関数を実装する。
7. 有効期限検証関数を実装する。
8. 単体テスト・結合テストを実装する。
9. gofmtを実行する。
10. go vet ./... を実行する。
11. go test ./... を実行する。
12. 実装内容とテスト結果をユーザーに報告する。

実装が完了したら、以下を報告してください。

- 実装した機能
- 追加・変更したファイル
- 使用した外部パッケージ
- テストの実行結果
- 未対応事項や既知の問題

仕様に記載されていない機能は実装しないでください。