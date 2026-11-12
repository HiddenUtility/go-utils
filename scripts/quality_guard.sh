#!/bin/bash

# scripts/quality_guard.sh
# 品質ガードスクリプト: go fmt, golangci-lint, go test を実行します。

set -e

# 書き込み権限のあるディレクトリをキャッシュとして使用
export GOTMPDIR=/tmp
export GOCACHE=/tmp/go-cache
mkdir -p $GOCACHE

echo "=== Quality Guard Start ==="

echo "[1/3] Running go fmt..."
go fmt ./...

echo "[2/3] Running golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
    echo "Error: golangci-lint is not installed."
    exit 1
fi
golangci-lint run -v

echo "[3/3] Running go test..."
go test -v ./...

echo "=== Quality Guard Passed! ==="
