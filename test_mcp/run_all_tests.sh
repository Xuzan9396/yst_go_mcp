#!/bin/bash

echo "======================================"
echo "YST Go MCP 测试套件"
echo "======================================"
echo

# 检查服务器是否存在
if [ ! -f "../yst-go-mcp" ]; then
    echo "❌ 服务器不存在，正在编译..."
    cd ..
    go build -o yst-go-mcp ./cmd/yst-go-mcp
    if [ $? -ne 0 ]; then
        echo "❌ 编译失败"
        exit 1
    fi
    cd test_mcp
    echo "✓ 编译成功"
    echo
fi

echo "测试 1: MCP 客户端连接测试"
echo "--------------------------------------"
go run TestClient.go
if [ $? -ne 0 ]; then
    echo "❌ 测试 1 失败"
    exit 1
fi
echo

echo "测试 2: STDIO 协议测试"
echo "--------------------------------------"
go run TestStdio.go
if [ $? -ne 0 ]; then
    echo "❌ 测试 2 失败"
    exit 1
fi
echo

echo "测试 3: 完整功能测试"
echo "--------------------------------------"
go run TestFull.go
if [ $? -ne 0 ]; then
    echo "❌ 测试 3 失败"
    exit 1
fi
echo

echo "======================================"
echo "✅ 所有自动化测试通过！"
echo "======================================"
echo
echo "💡 提示："
echo "   如需测试浏览器登录，请运行："
echo "   go run TestBrowserLogin.go"
