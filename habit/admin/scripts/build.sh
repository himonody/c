#!/bin/bash

# Habit Admin 生产环境构建脚本

echo "🏗️  构建 Habit Admin 生产版本..."

# 清理旧的构建文件
echo "🧹 清理旧文件..."
rm -rf dist

# 检查 Node.js 版本
node_version=$(node -v | cut -d'v' -f2)
required_version="16"

if [ "$(printf '%s\n' "$required_version" "$node_version" | sort -V | head -n1)" != "$required_version" ]; then
    echo "❌ 需要 Node.js >= $required_version，当前版本: $node_version"
    exit 1
fi

# 安装依赖
echo "📦 安装依赖..."
pnpm install --frozen-lockfile

# 运行类型检查
echo "🔍 运行类型检查..."
pnpm type-check

# 运行代码检查
echo "🔧 运行代码检查..."
pnpm lint

# 构建应用
echo "🏗️  构建应用..."
pnpm build:prod

# 检查构建结果
if [ -d "dist" ] && [ -f "dist/index.html" ]; then
    echo "✅ 构建成功！"
    echo "📊 构建统计:"
    du -sh dist
    echo ""
    echo "🚀 部署说明:"
    echo "1. 本地预览: pnpm preview"
    echo "2. Docker 部署: docker build -t habit-admin ."
    echo "3. 直接部署: 将 dist 目录内容部署到 Web 服务器"
else
    echo "❌ 构建失败！"
    exit 1
fi
