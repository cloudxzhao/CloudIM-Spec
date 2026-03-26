#!/bin/bash

# 前端构建脚本

set -e

echo "Building CloudIM Client..."

cd "$(dirname "$0")"

# 安装依赖
npm install

# 构建
npm run build

echo "Build completed successfully!"
echo "Output: apps/client/dist/"
