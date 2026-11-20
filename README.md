
###  `golangci-lint`（v2）を使うコマンド


#### 1. インストール（Go 1.26 環境で v2 を入れる場合）

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

```

（※ `$(go env GOPATH)/bin` に `PATH` が通っているか、`golangci-lint version` で `2.x` 系になったかを確認してください。）

#### 2. lint を実行するコマンド

リポジトリのルートディレクトリで以下を実行します。

```bash
# CI と同じく設定ファイル（.golangci.yml）を読み込んで全パッケージをチェック
golangci-lint run ./...

```

## よく使う補助コマンド

```bash
# 設定ファイルの妥当性チェック
golangci-lint config verify

# 特定のファイルやディレクトリだけ絞って確認
golangci-lint run ./pkg/io/...

# 自動修正に対応している指摘があれば修正（errcheck の多くは手動対応が必要）
golangci-lint run --fix ./...

```

