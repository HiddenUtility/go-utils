# Codex Implementation Specification: Go JSON Utility Library

## 1. Project Overview

Go言語を使用して、JSONとGo構造体の相互変換、およびJSONファイルの読み書きを行うユーティリティライブラリを実装してください。

JSONのシリアライズ・デシリアライズにはGo標準パッケージの `encoding/json` を使用してください。

本ライブラリの主な目的は以下のとおりです。

1. 任意のGo構造体をJSON文字列に変換する。
2. JSON文字列を任意のGo構造体に変換する。
3. JSONファイルの読み込みと書き込みを行う。
4. インデント付きの整形されたJSONを出力する。
5. JSON入出力用の構造体（DTO）と、アプリケーション内部で使用するドメイン構造体を分離する。
6. DTOからドメイン構造体への変換時に、値の検証とデフォルト値の補完を行う。

特に重要なのは、**JSONの不完全なデータやnullを許容しつつ、ドメイン構造体ではポインタを極力使用せず、常に利用可能な値を保持すること**です。

## 2. Development Requirements

### 2.1 Programming Language

- Goを使用する。
- Go 1.18以上を対象とする。
- Goの標準的な設計・命名規則に従う。
- 安全性と可読性を重視する。
- 不必要な依存関係を追加しない。

### 2.2 Dependencies

Go標準パッケージのみを使用する。

| Package | Purpose |
|---|---|
| encoding/json | JSON Serialize / Deserialize |
| os | JSONファイルの読み書き |
| fmt | エラーメッセージの生成 |
| errors | エラー判定 |
| strings | 文字列処理 |
| time | 日時の検証 |
| net/mail | メールアドレスの形式検証 |
| path/filepath | ファイルパス処理 |
| testing | テスト |

外部ライブラリは原則として追加しないこと。

### 2.3 Dependency License

外部ライブラリを追加する必要が生じた場合は、MIT、Apache-2.0などの一般的に利用可能なライセンスであることを確認する。

既存の標準パッケージで実装できる場合は、新しい依存関係を追加しないこと。

## 3. Project Structure

以下の構成を基本として実装すること。

```text
json-utility/
├── go.mod
├── json.go
├── error.go
├── validation.go
├── models/
│   ├── dto.go
│   └── domain.go
├── examples/
│   └── user_example.go
└── tests/
    ├── json_test.go
    ├── file_test.go
    └── conversion_test.go
```

各ファイルの役割は以下のとおり。

| File | Responsibility |
|---|---|
| go.mod | Go Moduleの定義 |
| json.go | JSON Serialize / Deserialize、ファイルIO |
| error.go | エラー型の定義 |
| models/dto.go | JSON入出力用構造体 |
| models/domain.go | ドメイン構造体 |
| validation.go | DTOからDomainへの変換時の検証 |
| tests/ | 各機能のテスト |
| examples/ | 使用例 |

必要に応じてファイルを追加してもよいが、責務を明確に分離すること。

Goのパッケージ設計に従い、循環インポートが発生しないようにすること。

## 4. JSON Utility Functions

任意の型に対応するジェネリック関数を実装すること。

### 4.1 Serialize

任意のGo構造体をJSON文字列に変換する。

Function:

```go
func ToJSON[T any](value T) (string, error)
```

Requirements:

- ジェネリック型を使用する。
- 任意のJSON変換可能な型を受け取る。
- JSON文字列を返す。
- `json.Marshal` を使用する。
- 変換に失敗した場合はエラーを返す。
- パニックを発生させない。

### 4.2 Deserialize

JSON文字列を任意の型に変換する。

Function:

```go
func FromJSON[T any](data string) (T, error)
```

Requirements:

- ジェネリック型を使用する。
- 呼び出し側が変換先の型を指定できる。
- `json.Unmarshal` を使用する。
- JSON形式が不正な場合はエラーを返す。
- 型が一致しない場合はエラーを返す。
- エラー発生時にはゼロ値とエラーを返す。

使用例:

```go
user, err := FromJSON[UserDTO](jsonString)
if err != nil {
    return err
}
```

### 4.3 Pretty Serialize

インデント付きのJSON文字列を生成する。

Function:

```go
func ToJSONPretty[T any](value T) (string, error)
```

Requirements:

- `json.MarshalIndent` を使用する。
- JSONの階層構造が分かるようにインデントする。
- インデントには2スペースを使用する。
- Serializeと同様にエラーを返す。

### 4.4 Generic Type Requirements

JSONユーティリティ関数は特定の構造体に依存してはならない。

例えば、以下のような異なる型で利用できること。

```go
userJSON, err := ToJSON(user)
configJSON, err := ToJSON(config)
settingsJSON, err := ToJSON(settings)
```

## 5. JSON File I/O

JSONファイルの読み込みと書き込みを実装する。

### 5.1 Read JSON File

Function:

```go
func ReadJSONFile[T any](path string) (T, error)
```

Requirements:

- JSONファイルを読み込む。
- 読み込んだJSONを指定された型に変換する。
- `os.ReadFile` を使用する。
- ファイルが存在しない場合はエラーを返す。
- ファイルアクセスに失敗した場合はエラーを返す。
- JSON形式が不正な場合はエラーを返す。
- 型が一致しない場合はエラーを返す。
- エラー発生時にはゼロ値とエラーを返す。

### 5.2 Write JSON File

Function:

```go
func WriteJSONFile[T any](path string, value T) error
```

Requirements:

- 指定されたパスにJSONファイルを書き込む。
- JSONはインデントなしのコンパクト形式とする。
- `json.Marshal` を使用する。
- `os.WriteFile` を使用する。
- 既存ファイルがある場合は上書きする。
- ファイルが存在しない場合は新規作成する。
- ファイルのパーミッションは0644を基本とする。
- 書き込みに失敗した場合はエラーを返す。

### 5.3 Write Pretty JSON File

Function:

```go
func WriteJSONFilePretty[T any](path string, value T) error
```

Requirements:

- 指定されたパスにJSONファイルを書き込む。
- `json.MarshalIndent` を使用する。
- インデントには2スペースを使用する。
- 既存ファイルがある場合は上書きする。
- ファイルが存在しない場合は新規作成する。
- ファイルのパーミッションは0644を基本とする。
- 書き込みに失敗した場合はエラーを返す。

### 5.4 File I/O Restrictions

- ファイルの読み書きには `os` を使用する。
- ファイルパスはstringで受け取る。
- エラーを握りつぶさない。
- ファイルIOのエラーとJSON変換エラーを区別できるようにする。
- 暗黙的に親ディレクトリを作成する必要はない。
- ファイルの文字コードはUTF-8を前提とする。

## 6. DTO and Domain Model Design

JSON入出力用の構造体と、アプリケーション内部で使用するドメイン構造体を分離する。

### 6.1 DTO (Data Transfer Object)

DTOはJSONとの変換を目的とした構造体である。

Requirements:

- JSONのフィールド名を明示的に指定するため、`json` タグを使用する。
- JSONにnullが含まれる可能性を考慮する。
- JSONにフィールドが存在しない可能性を考慮する。
- 必要に応じてポインタ型を使用する。
- JSONの入力形式に適した構造にする。

JSONのnullとフィールド未指定を区別する必要がない場合は、両方とも同じデフォルト値に変換できる設計にすること。

### 6.2 Domain Model

Domain Modelはアプリケーション内部で使用する構造体である。

Requirements:

- 原則としてポインタ型を使用しない。
- すべてのフィールドに利用可能な値を保持する。
- DTOから変換された段階でデフォルト値の補完が完了していること。
- 値の検証が完了していること。
- JSONの入力形式に依存しない構造にする。

重要：

DTOとDomain Modelを同一の構造体にしてはならない。

DTOからDomain Modelへの変換を明示的に実装すること。

## 7. Sample Domain Models

動作確認と実装例としてUserドメインを定義する。

入れ子構造、配列、null、複数のデータ型を扱えるようにすること。

### 7.1 User DTO

```go
package models

type UserDTO struct {
    ID        *uint64       `json:"id"`
    Name      *string       `json:"name"`
    Age       *int64        `json:"age"`
    Email     *string       `json:"email"`
    Active    *bool         `json:"active"`
    Tags      []string      `json:"tags"`
    Address   *AddressDTO   `json:"address"`
    Settings  *SettingsDTO  `json:"settings"`
    CreatedAt *string       `json:"created_at"`
}
```

### 7.2 Address DTO

```go
package models

type AddressDTO struct {
    Country    *string `json:"country"`
    PostalCode *string `json:"postal_code"`
    Prefecture *string `json:"prefecture"`
    City       *string `json:"city"`
    Street     *string `json:"street"`
}
```

### 7.3 Settings DTO

```go
package models

type SettingsDTO struct {
    Theme                 *string  `json:"theme"`
    NotificationsEnabled  *bool    `json:"notifications_enabled"`
    Language              *string  `json:"language"`
    DisplayOptions        []string `json:"display_options"`
}
```

注意：

- ポインタ型のフィールドはJSONのnullと未指定の両方でnilになる。
- スライス型のフィールドはJSONのnullと未指定の両方でnilになる。
- どちらもDomainへの変換時にはデフォルト値を設定する。
- JSONの型が不一致の場合はDeserialize Errorとする。

DTOには必要に応じてDefault相当の初期化関数を実装してもよい。

### 7.4 Domain User

```go
package models

type User struct {
    ID        uint64
    Name      string
    Age       uint32
    Email     string
    Active    bool
    Tags      []string
    Address   Address
    Settings  Settings
    CreatedAt string
}
```

### 7.5 Domain Address

```go
package models

type Address struct {
    Country    string
    PostalCode string
    Prefecture string
    City       string
    Street     string
}
```

### 7.6 Domain Settings

```go
package models

type Settings struct {
    Theme                string
    NotificationsEnabled bool
    Language             string
    DisplayOptions       []string
}
```

Domain Modelには原則としてポインタ型を使用しないこと。

## 8. Default Values

DTOからDomain Modelへ変換する際に、nullまたはフィールド未指定の値をデフォルト値に置き換える。

以下のデフォルト値を使用する。

### 8.1 User Defaults

| Field | Default |
|---|---|
| ID | 0 |
| Name | "" |
| Age | 0 |
| Email | "" |
| Active | false |
| Tags | 空のスライス |
| Address | Address{} |
| Settings | DefaultSettings() |
| CreatedAt | "" |

### 8.2 Address Defaults

| Field | Default |
|---|---|
| Country | "" |
| PostalCode | "" |
| Prefecture | "" |
| City | "" |
| Street | "" |

### 8.3 Settings Defaults

| Field | Default |
|---|---|
| Theme | "light" |
| NotificationsEnabled | false |
| Language | "ja" |
| DisplayOptions | 空のスライス |

### 8.4 Default Implementation

GoにはRustのDefaultトレイトに相当する仕組みがないため、明示的な初期化関数を実装すること。

以下の関数を実装する。

```go
func DefaultUser() User {
    return User{
        ID:        0,
        Name:      "",
        Age:       0,
        Email:     "",
        Active:    false,
        Tags:      []string{},
        Address:   DefaultAddress(),
        Settings:  DefaultSettings(),
        CreatedAt: "",
    }
}

func DefaultAddress() Address {
    return Address{
        Country:    "",
        PostalCode: "",
        Prefecture: "",
        City:       "",
        Street:     "",
    }
}

func DefaultSettings() Settings {
    return Settings{
        Theme:                "light",
        NotificationsEnabled: false,
        Language:             "ja",
        DisplayOptions:       []string{},
    }
}
```

デフォルト値は一元管理し、DTOからDomainへの変換処理で個別にハードコードしないこと。

## 9. DTO to Domain Conversion

DTOからDomain Modelへの変換を実装する。

変換処理では、以下の2つを実施する。

1. デフォルト値の補完
2. 値のバリデーション

### 9.1 Conversion Functions

以下のような変換関数を実装する。

```go
func (dto UserDTO) ToDomain() (User, error)
```

AddressDTO、SettingsDTOについても必要に応じて同様の変換関数を実装する。

```go
func (dto AddressDTO) ToDomain() (Address, error)

func (dto SettingsDTO) ToDomain() (Settings, error)
```

### 9.2 Null Handling

DTOのポインタ型フィールドがnilの場合、Domain Modelのデフォルト値を設定する。

例：

```go
func stringOrDefault(value *string, defaultValue string) string {
    if value == nil {
        return defaultValue
    }
    return *value
}
```

同様に、uint64、int64、boolなどの型についても必要に応じて共通の処理を実装すること。

ただし、Goのジェネリクスを使用して複雑にしすぎないこと。

### 9.3 Slice Handling

スライス型のフィールドがnilの場合は空のスライスを設定する。

例：

```go
tags := dto.Tags
if tags == nil {
    tags = []string{}
}
```

Domain Modelではnilスライスではなく空のスライスを基本とすること。

### 9.4 Nested Structure Handling

入れ子構造のDTOがnilの場合も、Domain Modelのデフォルト値を設定する。

例：

```go
var address Address

if dto.Address != nil {
    converted, err := dto.Address.ToDomain()
    if err != nil {
        return User{}, err
    }
    address = converted
} else {
    address = DefaultAddress()
}
```

Settingsについても同様に処理する。

## 10. Validation Rules

DTOからDomain Modelへ変換する際に、以下の検証を行う。

### 10.1 General Rules

| Input | Behavior |
|---|---|
| 正常な値 | そのまま変換 |
| null | デフォルト値を設定 |
| フィールド未指定 | デフォルト値を設定 |
| 型不一致 | JSON Deserialize Error |
| 値が制約に違反 | Validation Error |

型不一致はJSONデシリアライズ時点でエラーとし、バリデーションエラーとは区別すること。

nullやフィールド未指定はエラーにしない。

### 10.2 User Validation

| Field | Validation |
|---|---|
| ID | 0以上の値 |
| Name | 文字数が100文字以内 |
| Age | 0以上150以下 |
| Email | メールアドレス形式 |
| Active | 追加制約なし |
| Tags | 各要素が100文字以内 |
| CreatedAt | 空文字列を許容。空でない場合は日時形式を検証 |

注意：

- AgeはDTOではint64、Domainではuint32とする。
- 負数の場合はValidation Errorとする。
- 150を超える場合もValidation Errorとする。
- IDはuint64として扱うため、負数はJSON Deserialize Errorとする。
- メールアドレスの検証には `net/mail` を使用する。
- メールアドレスの検証では `mail.ParseAddress` の結果が入力文字列そのものと一致することも確認する。表示名付きアドレスは許可しない。
- 日時形式はRFC 3339形式を基本とする。
- 日時の検証には `time.Parse(time.RFC3339, value)` を使用する。
- CreatedAtが空文字列の場合はデフォルト値として許可する。

### 10.3 Address Validation

| Field | Validation |
|---|---|
| Country | 100文字以内 |
| PostalCode | 20文字以内 |
| Prefecture | 100文字以内 |
| City | 100文字以内 |
| Street | 200文字以内 |

空文字列は許容する。

### 10.4 Settings Validation

| Field | Validation |
|---|---|
| Theme | "light" または "dark" |
| NotificationsEnabled | 追加制約なし |
| Language | 10文字以内 |
| DisplayOptions | 各要素が50文字以内 |

Themeについては、nullまたは未指定の場合は "light" とする。

不正なThemeが入力された場合はValidation Errorとする。

### 10.5 String Length Validation

GoではUTF-8の文字列をバイト数で数えることに注意すること。

今回の文字数制限は、原則としてUnicodeコードポイント数を基準とする。

必要に応じて `unicode/utf8` の `utf8.RuneCountInString` を使用すること。

### 10.6 Validation Error Handling

検証に失敗した場合は、どのフィールドが原因なのか分かるエラーを返すこと。

例：

```text
Validation error: age must be between 0 and 150
Validation error: email has invalid format
Validation error: settings.theme must be "light" or "dark"
```

ValidationErrorはフィールド名とエラー内容を保持できる構造体とする。

例：

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf(
        "validation error: %s: %s",
        e.Field,
        e.Message,
    )
}
```

複数の検証エラーがある場合は、最初のエラーで終了してよい。

ただし、検証順序は安定させること。

## 11. Error Design

エラーの種類を区別できる設計にすること。

最低限、以下のエラーを扱う。

- JSON Serialize Error
- JSON Deserialize Error
- File I/O Error
- Validation Error

必要に応じて独自のエラー型を定義すること。

例：

```go
var (
    ErrSerialize   = errors.New("json serialize error")
    ErrDeserialize = errors.New("json deserialize error")
    ErrFileIO      = errors.New("file io error")
)
```

エラーを返す場合は、原因となったエラーを保持すること。

例：

```go
return "", fmt.Errorf("%w: %w", ErrSerialize, err)
```

ValidationErrorは独立した構造体として定義する。

Requirements:

- エラーの原因を確認できること。
- `errors.Is` によるエラー判定を可能にすること。
- `errors.As` によるValidationErrorの取得を可能にすること。
- 元のエラー情報を失わないこと。
- エラーを握りつぶさないこと。
- panicを通常のエラー処理に使用しないこと。

## 12. Example JSON

以下のようなJSONを扱えるようにする。

### 12.1 Complete JSON

```json
{
  "id": 100,
  "name": "Taro",
  "age": 30,
  "email": "taro@example.com",
  "active": true,
  "tags": ["go", "json", "developer"],
  "address": {
    "country": "Japan",
    "postal_code": "100-0001",
    "prefecture": "Tokyo",
    "city": "Chiyoda",
    "street": "Marunouchi 1-1"
  },
  "settings": {
    "theme": "dark",
    "notifications_enabled": true,
    "language": "ja",
    "display_options": ["compact", "show_email"]
  },
  "created_at": "2026-01-01T12:00:00+09:00"
}
```

### 12.2 JSON with Null Values

```json
{
  "id": 100,
  "name": null,
  "age": null,
  "email": null,
  "active": null,
  "tags": null,
  "address": null,
  "settings": null,
  "created_at": null
}
```

このJSONは正常にDomain Modelへ変換できること。

各フィールドには定義済みのデフォルト値が設定されること。

### 12.3 JSON with Missing Fields

```json
{
  "id": 100,
  "name": "Taro"
}
```

このJSONも正常にDomain Modelへ変換できること。

存在しないフィールドはデフォルト値で補完すること。

### 12.4 JSON with Invalid Values

```json
{
  "id": 100,
  "name": "Taro",
  "age": -1,
  "email": "invalid-email"
}
```

このJSONはDTOへのDeserializeには成功するが、Domain Modelへの変換時にValidation Errorを返すこと。

## 13. Public API

以下のAPIを外部から利用できるようにする。

### JSON Conversion

```go
func ToJSON[T any](value T) (string, error)

func FromJSON[T any](data string) (T, error)

func ToJSONPretty[T any](value T) (string, error)
```

### JSON File I/O

```go
func ReadJSONFile[T any](path string) (T, error)

func WriteJSONFile[T any](path string, value T) error

func WriteJSONFilePretty[T any](path string, value T) error
```

### Model Conversion

```go
func (dto UserDTO) ToDomain() (User, error)

func (dto AddressDTO) ToDomain() (Address, error)

func (dto SettingsDTO) ToDomain() (Settings, error)
```

### Default Functions

```go
func DefaultUser() User

func DefaultAddress() Address

func DefaultSettings() Settings
```

User、Address、Settings、各DTOおよび必要なエラー型も外部から利用できるようにする。

Goの公開識別子は大文字で始めること。

## 14. Usage Example

以下のような利用方法を提供すること。

### 14.1 JSON String to Domain

```go
package main

import (
    "fmt"

    "your-module/models"
    "your-module"
)

func main() {
    jsonData := `{
        "id": 100,
        "name": "Taro",
        "age": 30,
        "email": "taro@example.com"
    }`

    dto, err := jsonutility.FromJSON[models.UserDTO](jsonData)
    if err != nil {
        panic(err)
    }

    user, err := dto.ToDomain()
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", user)
}
```

### 14.2 Domain to JSON

```go
user := models.DefaultUser()

jsonData, err := jsonutility.ToJSONPretty(user)
if err != nil {
    return err
}

fmt.Println(jsonData)
```

### 14.3 JSON File to Domain

```go
dto, err := jsonutility.ReadJSONFile[models.UserDTO]("user.json")
if err != nil {
    return err
}

user, err := dto.ToDomain()
if err != nil {
    return err
}
```

### 14.4 Domain to JSON File

```go
err := jsonutility.WriteJSONFilePretty("user.json", user)
if err != nil {
    return err
}
```

## 15. Unit Tests

必ずユニットテストおよび統合テストを実装すること。

Go標準の `testing` パッケージを使用する。

テストはテーブル駆動テストを基本とする。

### 15.1 JSON Conversion Tests

以下をテストする。

- 正常な構造体をJSON文字列に変換できる。
- JSON文字列を指定した構造体に変換できる。
- インデント付きJSONを生成できる。
- 不正なJSONでエラーが返る。
- 型が一致しない場合にエラーが返る。
- 複数の異なる構造体でジェネリック関数を利用できる。

### 15.2 File I/O Tests

以下をテストする。

- JSONファイルを読み込める。
- JSONファイルを書き込める。
- インデント付きJSONファイルを書き込める。
- 書き込んだJSONを再度読み込める。
- 存在しないファイルを読み込んだ場合にエラーになる。
- 書き込みに失敗した場合にエラーになる。
- 不正なJSONファイルを読み込んだ場合にエラーになる。

テストファイルは `t.TempDir()` を使用してテストごとに独立した一時ディレクトリを作成すること。

テスト終了後の後片付けは `testing` パッケージに任せること。

### 15.3 DTO to Domain Tests

以下を必ずテストする。

1. すべてのフィールドが正常なJSONから正しく変換できる。
2. すべてのフィールドがnullの場合にデフォルト値が設定される。
3. フィールドが省略された場合にデフォルト値が設定される。
4. 入れ子構造がnullの場合にデフォルト値が設定される。
5. 配列がnullの場合に空のスライスが設定される。
6. 正常な配列が正しく変換される。
7. Ageが負数の場合にエラーになる。
8. Ageが150を超える場合にエラーになる。
9. Emailが不正な場合にエラーになる。
10. Nameが100文字を超える場合にエラーになる。
11. Themeがlight/dark以外の場合にエラーになる。
12. CreatedAtが不正な日時の場合にエラーになる。
13. Domain Modelの各フィールドがポインタ型ではないこと。
14. 空文字列や空配列が仕様どおりに扱われること。

### 15.4 Default Tests

User、Address、Settingsそれぞれのデフォルト値が仕様どおりであることをテストする。

## 16. Implementation Constraints

以下の事項を必ず守ること。

1. JSON変換関数はジェネリック型で実装する。
2. JSON入出力用DTOとDomain Modelを分離する。
3. Domain Modelには原則としてポインタ型を使用しない。
4. nullとフィールド未指定はデフォルト値で補完する。
5. 型不一致はDeserialize Errorとして扱う。
6. 値の制約違反はValidation Errorとして扱う。
7. エラーを握りつぶさない。
8. 不必要な依存関係を追加しない。
9. panicを通常のエラー処理に使用しない。
10. サンプルコードだけでなく、実際にコンパイル可能なライブラリとして完成させる。
11. すべての公開APIに適切なGoDocコメントを記述する。
12. `gofmt` を実行する。
13. `go vet ./...` を実行し、問題を解消する。
14. `go test ./...` を実行し、すべてのテストを成功させる。

## 17. Implementation Workflow

以下の順番で実装すること。

1. go.modの作成とGoバージョンの確認。
2. プロジェクト構成を作成する。
3. JsonErrorとValidationErrorを実装する。
4. DTO構造体を実装する。
5. Domain Modelとデフォルト値関数を実装する。
6. DTOからDomain Modelへの変換を実装する。
7. Validation処理を実装する。
8. JSON Serialize / Deserializeを実装する。
9. JSONファイルの読み書きを実装する。
10. サンプルコードを作成する。
11. ユニットテスト・統合テストを実装する。
12. gofmtを実行する。
13. go vet ./... を実行する。
14. go test ./... を実行する。
15. すべてのテスト結果と実装内容を確認する。

各工程でコンパイルエラーやテストエラーが発生した場合は、原因を調査して修正すること。

## 18. Completion Criteria

以下をすべて満たした場合に実装完了とする。

- [ ] JSON Serializeが動作する。
- [ ] JSON Deserializeが動作する。
- [ ] Pretty Serializeが動作する。
- [ ] JSONファイルの読み込みが動作する。
- [ ] JSONファイルの書き込みが動作する。
- [ ] Pretty JSONファイルの書き込みが動作する。
- [ ] ジェネリック型による任意の構造体への変換が動作する。
- [ ] DTOとDomain Modelが分離されている。
- [ ] DTOからDomain Modelへの変換が動作する。
- [ ] nullがデフォルト値に補完される。
- [ ] フィールド未指定がデフォルト値に補完される。
- [ ] 入れ子構造が正しく変換される。
- [ ] 配列が正しく変換される。
- [ ] バリデーションが動作する。
- [ ] エラーが適切に処理される。
- [ ] すべてのテストが成功する。
- [ ] gofmtが成功する。
- [ ] go vet ./... が成功する。

## 19. Final Report

実装完了後、以下を報告すること。

1. 実装した機能の一覧
2. 作成・変更したファイルの一覧
3. 公開APIの一覧
4. DTOとDomain Modelの変換仕様
5. デフォルト値とバリデーションの仕様
6. 追加した依存関係
7. go test ./... の実行結果
8. go vet ./... の実行結果
9. 未実装項目や既知の問題

実行していないテストやコマンドを成功したと報告してはならない。

**重要：仕様にない機能を独自判断で追加しないこと。実装上の矛盾や不明点が見つかった場合は、実装を開始する前、または該当箇所の実装前にユーザーへ質問すること。**