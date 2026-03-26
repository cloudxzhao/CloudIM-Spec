#!/bin/bash

# 后端构建脚本

set -e

echo "Building CloudIM Server..."

cd "$(dirname "$0")"

# 构建
go build -o bin/server ./cmd/server

echo "Build completed successfully!"
echo "Binary: apps/server/bin/server"
