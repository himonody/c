#!/bin/bash

# Habit Admin 开发环境启动脚本

echo "🚀 启动 Habit Admin 开发环境..."

# 检查 Node.js 版本
node_version=$(node -v | cut -d'v' -f2)
required_version="16"

if [ "$(printf '%s\n' "$required_version" "$node_version" | sort -V | head -n1)" != "$required_version" ]; then
    echo "❌ 需要 Node.js >= $required_version，当前版本: $node_version"
    exit 1
fi

# 检查 pnpm 是否安装
if ! command -v pnpm &> /dev/null; then
    echo "📦 安装 pnpm..."
    npm install -g pnpm
fi

# 检查依赖是否安装
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖..."
    pnpm install
fi

# 检查环境变量文件
if [ ! -f ".env" ]; then
    echo "📝 创建环境变量文件..."
    cp .env.example .env
    echo "⚠️  请编辑 .env 文件配置正确的 API 地址"
fi

# 启动开发服务器
echo "🌟 启动开发服务器..."
pnpm dev
