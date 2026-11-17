# Codex Implementation Specification: Case Converter (Go)

## 1. Project Overview

Go言語で文字列のケーススタイルを相互変換するユーティリティライブラリを実装してください。

本ライブラリは、入力された文字列から単語を抽出し、指定されたケーススタイルに変換する機能を提供します。

### 対応するケーススタイル

以下の4種類に対応してください。

1. `camelCase`
2. `snake_case`
3. `PascalCase`
4. `kebab-case`

### 実装方針

- Go言語で実装すること。
- Goの標準ライブラリを優先して使用すること。
- 外部パッケージへの依存は追加しないこと。
- 文字列の解析とケース変換を明確に分離すること。
- すべての変換関数は決定的な結果を返すこと。
- 入力文字列を破壊的に変更しないこと。
- 単体テストを必ず実装すること。
- Unicode文字の処理にはGoの標準ライブラリを使用すること。
- Goの標準的なパッケージ構成と命名規則に従うこと。

---

## 2. Functional Requirements

### 2.1. Case Conversion Functions

以下の4つの公開関数を実装してください。

| Function | Description |
|---|---|
| `ToCamelCase` | camelCaseへ変換 |
| `ToSnakeCase` | snake_caseへ変換 |
| `ToPascalCase` | PascalCaseへ変換 |
| `ToKebabCase` | kebab-caseへ変換 |

#### Function Specifications

すべての関数について、以下の仕様を守ってください。

- 引数は `string` とする。
- 戻り値は `string` とする。
- 入力文字列が空の場合は空文字列を返す。
- 入力文字列に含まれる単語を解析してから変換する。
- 変換時に元の入力文字列を変更しない。
- 変換処理でパニックを発生させないこと。
- 不正なUTF-8を含む文字列についても、パニックせず安全に処理すること。

関数シグネチャ:

```go
func ToCamelCase(input string) string

func ToSnakeCase(input string) string

func ToPascalCase(input string) string

func ToKebabCase(input string) string
```

これらの関数は外部パッケージから利用できるように公開してください。

---

## 3. Word Parsing Specification

### 3.1. Overview

ケース変換に先立ち、入力文字列を単語単位に分割してください。

すべての変換関数で共通の単語解析処理を使用してください。

解析処理と変換処理を個別に実装することは禁止します。

例えば、`HTTPServer_test-Example` は以下の単語配列に変換されます。

```text
["HTTP", "Server", "test", "Example"]
```

### 3.2. Delimiter Characters

以下の文字は単語の区切り文字として扱ってください。

- アンダースコア `_`
- ハイフン `-`
- 半角スペース
- タブ
- 改行
- その他のUnicode空白文字（`unicode.IsSpace` がtrueを返す文字）

区切り文字が連続している場合は、1つの区切りとして扱ってください。

先頭および末尾の区切り文字は無視してください。

### 3.3. Uppercase Boundary

英小文字の直後に英大文字が出現した場合、単語を分割してください。

例:

```text
helloWorld
```

解析結果:

```text
["hello", "World"]
```

ASCII英字の大文字・小文字判定には `unicode` パッケージを使用してください。

ただし、単語境界の基本判定はASCII英字を中心に実装してください。

### 3.4. Consecutive Uppercase Boundary

連続する大文字は、次の文字が小文字の場合、その直前で分割してください。

ただし、最後の大文字は次の単語に含めてください。

例:

```text
HTTPServer
```

解析結果:

```text
["HTTP", "Server"]
```

```text
XMLParser
```

解析結果:

```text
["XML", "Parser"]
```

```text
ABC
```

解析結果:

```text
["ABC"]
```

実装では、現在の文字と前後の文字を参照し、連続する大文字の末尾を正しく判定してください。

### 3.5. Numeric Boundary

英字と数字の境界では単語を分割してください。

以下の2パターンを対象とします。

- 英字の直後に数字が出現した場合
- 数字の直後に英字が出現した場合

例:

```text
version2Value
```

解析結果:

```text
["version", "2", "Value"]
```

```text
123ABC
```

解析結果:

```text
["123", "ABC"]
```

```text
IPv4Address
```

解析結果:

```text
["IP", "v", "4", "Address"]
```

数字が連続する場合は、1つの単語として扱ってください。

例:

```text
item123
```

解析結果:

```text
["item", "123"]
```

数字の判定には `unicode.IsDigit` を使用してください。

ただし、境界判定の基本対象はASCII英字と数字とし、Unicodeの複雑な文字分類による単語境界解析は対象外とします。

### 3.6. Other Characters

英字以外の文字について、以下のルールを適用してください。

- 記号類は単語区切りとして扱う。
- Unicode文字については可能な範囲で保持する。
- UTF-8文字列をバイト単位で分割しない。
- 文字列の走査には `[]rune` または適切なUTF-8文字境界を使用する。
- Unicode文字の大文字・小文字変換にはGo標準ライブラリの機能を使用する。

Unicodeの複雑な単語境界解析や言語依存の形態素解析は対象外とします。

なお、不正なUTF-8を含む文字列を `[]rune` に変換した場合、無効なバイト列はUnicodeの置換文字 `utf8.RuneError` に変換されます。

この場合もパニックを発生させず、安全に処理してください。

---

## 4. Case Conversion Rules

単語解析後、以下の規則で変換してください。

### 4.1. camelCase

ルール:

1. 最初の単語はすべて小文字にする。
2. 2番目以降の単語は先頭を大文字にする。
3. 2番目以降の単語の残りは小文字にする。
4. 単語間に区切り文字を挿入しない。

例:

| Input | Output |
|---|---|
| `hello_world` | `helloWorld` |
| `HelloWorld` | `helloWorld` |
| `HTTPServer` | `httpServer` |
| `user-id` | `userId` |
| `version2Value` | `version2Value` |
| `FOO_BAR` | `fooBar` |

### 4.2. snake_case

ルール:

1. すべての単語を小文字にする。
2. 単語間に `_` を挿入する。

例:

| Input | Output |
|---|---|
| `helloWorld` | `hello_world` |
| `HelloWorld` | `hello_world` |
| `HTTPServer` | `http_server` |
| `user-id` | `user_id` |
| `version2Value` | `version_2_value` |
| `FOO_BAR` | `foo_bar` |

### 4.3. PascalCase

ルール:

1. すべての単語の先頭を大文字にする。
2. 各単語の残りを小文字にする。
3. 単語間に区切り文字を挿入しない。

例:

| Input | Output |
|---|---|
| `hello_world` | `HelloWorld` |
| `helloWorld` | `HelloWorld` |
| `HTTPServer` | `HttpServer` |
| `user-id` | `UserId` |
| `version2Value` | `Version2Value` |
| `FOO_BAR` | `FooBar` |

### 4.4. kebab-case

ルール:

1. すべての単語を小文字にする。
2. 単語間に `-` を挿入する。

例:

| Input | Output |
|---|---|
| `helloWorld` | `hello-world` |
| `HelloWorld` | `hello-world` |
| `HTTPServer` | `http-server` |
| `user_id` | `user-id` |
| `version2Value` | `version-2-value` |
| `FOO_BAR` | `foo-bar` |

---

## 5. Character Conversion Details

### 5.1. Lowercase Conversion

小文字への変換には、Go標準ライブラリの `strings.ToLower` を使用してください。

Unicode文字の小文字変換にも対応してください。

### 5.2. Uppercase Conversion

先頭文字の大文字化には、Go標準ライブラリの `unicode.ToUpper` を使用してください。

Unicode文字の大文字変換にも対応してください。

### 5.3. First Character Handling

単語の先頭文字だけを大文字にし、残りを小文字にする処理を共通化してください。

例えば、以下の関数を内部実装として用意してください。

```go
func lowercaseWord(word string) string

func capitalizeWord(word string) string
```

`capitalizeWord` は以下の処理を行ってください。

1. 入力単語を `[]rune` に変換する。
2. 空の単語であれば空文字列を返す。
3. 先頭のruneを `unicode.ToUpper` で大文字にする。
4. 2文字目以降を小文字にする。
5. 結果をstringとして返す。

空文字列の単語を受け取った場合もパニックを発生させないこと。

注意:

Go標準ライブラリの `unicode.ToUpper` は1文字を1文字に変換する機能です。Unicodeの大文字変換によって複数文字に展開されるケースはサポート対象外とします。

これはGo標準ライブラリの仕様による制約です。

---

## 6. Internal Architecture

以下の構造を推奨します。

### 6.1. Word Parser

```go
func splitWords(input string) []string
```

責務:

- 入力文字列の走査
- 区切り文字の判定
- 大文字・小文字境界の判定
- 数字境界の判定
- 単語配列の生成
- 不正なUTF-8を含む入力の安全な処理

実装上の注意:

- `[]rune` を利用して文字単位で処理すること。
- 必要に応じて前後のruneを参照すること。
- 空の単語は生成しないこと。
- 解析処理はすべての変換関数で共通化すること。

### 6.2. Word Normalization

```go
func lowercaseWord(word string) string

func capitalizeWord(word string) string
```

責務:

- 単語の小文字化
- 単語の先頭文字の大文字化
- Unicode文字への対応

### 6.3. Public Conversion Functions

以下の公開関数を実装してください。

```go
func ToCamelCase(input string) string

func ToSnakeCase(input string) string

func ToPascalCase(input string) string

func ToKebabCase(input string) string
```

各公開関数は、共通の `splitWords` を呼び出し、変換規則に従って結果を組み立ててください。

変換処理に重複した単語解析ロジックを実装しないこと。

文字列の組み立てには、必要に応じて `strings.Builder` を使用してください。

---

## 7. Error Handling

本ライブラリは文字列変換ユーティリティであるため、通常の入力に対してエラーを返す必要はありません。

- Goのstringは任意のバイト列を保持できるため、不正なUTF-8が含まれる可能性がある。
- 不正なUTF-8を含む場合もパニックを発生させないこと。
- 空文字列は正常入力として扱う。
- 区切り文字だけの入力も正常入力として扱う。
- 空の単語が生成された場合は結果から除外する。
- UTF-8文字列をバイト単位で不正に切断しないこと。
- 不正なUTF-8は `[]rune` への変換時に置換文字へ変換されることを許容する。
- パニックを発生させないこと。

---

## 8. Unit Test Specification

Go標準の `testing` パッケージを使用して単体テストを実装してください。

テストファイル名は以下としてください。

```text
case_converter_test.go
```

すべての公開関数について、以下のテストを実施してください。

### 8.1. Basic Tests

以下の入力と期待値をテストしてください。

| Input | camelCase | snake_case | PascalCase | kebab-case |
|---|---|---|---|---|
| `hello_world` | `helloWorld` | `hello_world` | `HelloWorld` | `hello-world` |
| `helloWorld` | `helloWorld` | `hello_world` | `HelloWorld` | `hello-world` |
| `HelloWorld` | `helloWorld` | `hello_world` | `HelloWorld` | `hello-world` |
| `HTTPServer` | `httpServer` | `http_server` | `HttpServer` | `http-server` |
| `XMLParser` | `xmlParser` | `xml_parser` | `XmlParser` | `xml-parser` |
| `user-id` | `userId` | `user_id` | `UserId` | `user-id` |
| `FOO_BAR` | `fooBar` | `foo_bar` | `FooBar` | `foo-bar` |
| `version2Value` | `version2Value` | `version_2_value` | `Version2Value` | `version-2-value` |
| `123ABC` | `123Abc` | `123_abc` | `123Abc` | `123-abc` |

テストはテーブル駆動テスト形式で実装することを推奨します。

### 8.2. Empty Input Tests

以下の入力をテストしてください。

```text
""
" "
"___"
"---"
"   "
```

すべての変換関数で空文字列が返ること。

### 8.3. Consecutive Delimiter Tests

以下の入力をテストしてください。

```text
"hello__world"
"hello--world"
"hello _ world"
"__hello__world__"
```

連続する区切り文字によって空の単語が生成されないこと。

### 8.4. Boundary Tests

以下の入力をテストしてください。

```text
"aB"
"ABc"
"ABC"
"HTTPServer"
"XMLHttpRequest"
"foo123bar"
"123foo456"
"version2Value"
```

期待する単語分割が正しく行われること。

内部関数 `splitWords` に対しても、単語配列を比較する単体テストを実装してください。

期待する解析結果:

| Input | Expected Words |
|---|---|
| `aB` | `["a", "B"]` |
| `ABc` | `["A", "Bc"]` |
| `ABC` | `["ABC"]` |
| `HTTPServer` | `["HTTP", "Server"]` |
| `XMLHttpRequest` | `["XML", "Http", "Request"]` |
| `foo123bar` | `["foo", "123", "bar"]` |
| `123foo456` | `["123", "foo", "456"]` |
| `version2Value` | `["version", "2", "Value"]` |

### 8.5. Unicode Tests

以下の入力をテストしてください。

```text
"こんにちは_world"
"café_name"
"überValue"
```

確認事項:

- UTF-8文字列の処理でパニックしないこと。
- Unicode文字を不正に切断しないこと。
- Go標準ライブラリの大文字・小文字変換が適用されること。
- 日本語など、ASCII英字以外の文字が可能な範囲で保持されること。

Unicodeの単語境界判定は本仕様の対象外であるため、ASCII英字と数字を中心とした境界判定を行ってください。

### 8.6. Invalid UTF-8 Tests

Goのstringは不正なUTF-8を保持できるため、不正なUTF-8を含む入力についてもテストしてください。

例:

```go
input := string([]byte{0xff, 0xfe, 'a', 'b', 'c'})
```

確認事項:

- 4種類すべての変換関数でパニックしないこと。
- 不正なUTF-8を含む入力を安全に処理できること。
- 置換文字への変換を許容すること。

---

## 9. Coding Standards

以下のルールを遵守してください。

- Goの標準的な命名規則に従う。
- 公開関数名は `PascalCase` とする。
- 非公開関数名は `camelCase` とする。
- パッケージ名は短く、簡潔で、適切な名前にする。
- 公開関数にはGoDoc形式のドキュメントコメントを付ける。
- 不要なメモリアロケーションを避ける。
- 不要な外部パッケージを追加しない。
- 不必要な `panic` を使用しない。
- 不必要な `unsafe` を使用しない。
- Go標準のフォーマット規則に従う。
- `gofmt` によるフォーマットを適用する。
- `go vet` の警告を可能な限り解消する。
- Go標準の `testing` パッケージを使用する。

---

## 10. Implementation Restrictions

以下の実装は禁止します。

1. 変換関数ごとに異なる単語解析ロジックを実装すること。
2. UTF-8文字列をバイト単位で不正に切断すること。
3. 空文字列を考慮しない実装。
4. 連続する区切り文字によって空の単語を生成する実装。
5. 正規表現パッケージなど、不要な外部依存を追加すること。
6. 入力文字列を破壊的に変更すること。
7. 指定された変換規則を独自に変更すること。
8. Unicodeの大文字・小文字変換に対応するためだけに外部ライブラリを追加すること。
9. `splitWords` の処理を各公開変換関数に重複実装すること。

---

## 11. Completion Criteria

実装完了時には、以下をすべて満たしてください。

- [ ] 4種類の公開変換関数が実装されている。
- [ ] 共通の単語解析処理が実装されている。
- [ ] 大文字・小文字境界が正しく判定される。
- [ ] 連続する大文字の境界が正しく判定される。
- [ ] 数字と英字の境界が正しく判定される。
- [ ] 区切り文字が正しく処理される。
- [ ] Unicode文字を不正に切断しない。
- [ ] 不正なUTF-8を含む入力でパニックしない。
- [ ] 空文字列や区切り文字だけの入力に対応している。
- [ ] すべての単体テストが成功する。
- [ ] `gofmt -d .` でフォーマット差分が発生しない。
- [ ] `go vet ./...` が成功する。
- [ ] `go test ./...` が成功する。

---

## 12. Final Report

実装完了後、以下の内容を報告してください。

1. 実装した機能の一覧
2. 作成・変更したファイルの一覧
3. 実装した単語解析ロジックの概要
4. 実施したテストの一覧
5. 各コマンドの実行結果
   - `gofmt -d .`
   - `go vet ./...`
   - `go test ./...`
6. 未対応事項や既知の制限事項

テストが失敗している場合、成功したと報告してはいけません。

仕様に記載されていない機能を独自判断で追加しないでください。

以上の仕様に従って実装してください。