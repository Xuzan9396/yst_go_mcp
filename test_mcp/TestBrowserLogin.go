package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	fmt.Println("=== 测试浏览器登录模块 ===\n")

	serverPath, err := filepath.Abs("../yst-go-mcp")
	if err != nil {
		log.Fatalf("获取服务器路径失败: %v", err)
	}

	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		log.Fatalf("服务器不存在: %s\n请先编译: go build -o yst-go-mcp ./cmd/yst-go-mcp", serverPath)
	}

	fmt.Printf("连接到服务器: %s\n\n", serverPath)

	// 创建 STDIO 客户端
	c, err := client.NewStdioMCPClient(serverPath)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Second) // 6分钟超时
	defer cancel()

	// 初始化连接
	fmt.Println("📡 初始化 MCP 连接...")
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "browser-login-test",
		Version: "1.0.0",
	}

	serverInfo, err := c.Initialize(ctx, initReq)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Printf("✓ 连接成功：%s v%s\n\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 测试浏览器登录
	fmt.Println("🌐 测试浏览器登录功能")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("⚠️  注意事项：")
	fmt.Println("  1. 此测试会打开浏览器窗口")
	fmt.Println("  2. 需要手动完成 Google OAuth 登录")
	fmt.Println("  3. 登录超时时间：300 秒（5分钟）")
	fmt.Println("  4. 登录成功后会自动保存 Cookie")
	fmt.Println()
	fmt.Print("按 Enter 键开始测试（或 Ctrl+C 取消）...")
	fmt.Scanln()

	fmt.Println("\n🚀 调用 browser_login 工具...")
	fmt.Println("提示：请在弹出的浏览器中完成登录操作\n")

	// 调用 browser_login 工具
	callReq := mcp.CallToolRequest{}
	callReq.Params.Name = "browser_login"
	callReq.Params.Arguments = map[string]interface{}{
		"timeout": 300, // 5分钟超时
	}

	startTime := time.Now()
	result, err := c.CallTool(ctx, callReq)
	elapsed := time.Since(startTime)

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 测试结果")
	fmt.Println(strings.Repeat("=", 60))

	if err != nil {
		fmt.Printf("❌ 测试失败: %v\n", err)
		fmt.Printf("⏱  耗时: %.1f 秒\n", elapsed.Seconds())
		os.Exit(1)
	}

	fmt.Printf("✅ 测试成功！\n")
	fmt.Printf("⏱  耗时: %.1f 秒\n", elapsed.Seconds())
	fmt.Println()

	// 显示返回内容
	if len(result.Content) > 0 {
		fmt.Println("📄 返回内容:")
		for _, content := range result.Content {
			if textContent, ok := mcp.AsTextContent(content); ok {
				fmt.Printf("   %s\n", textContent.Text)
			}
		}
	}

	if result.IsError {
		fmt.Println("\n⚠️  工具返回了错误状态")
	}

	// 验证 Cookie 是否已保存
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("🔍 验证 Cookie 保存状态")
	fmt.Println(strings.Repeat("-", 60))

	cookiePaths := []string{
		"../data/cookies.json",
		filepath.Join(os.Getenv("HOME"), ".yst_go_mcp/data/cookies.json"),
	}

	cookieFound := false
	for _, path := range cookiePaths {
		if _, err := os.Stat(path); err == nil {
			info, _ := os.Stat(path)
			fmt.Printf("✓ Cookie 文件存在: %s\n", path)
			fmt.Printf("  大小: %d 字节\n", info.Size())
			fmt.Printf("  修改时间: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
			cookieFound = true
			break
		}
	}

	if !cookieFound {
		fmt.Println("⚠️  未找到 Cookie 文件")
	}

	// 测试采集功能（验证登录状态）
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("🔍 验证登录状态（尝试采集数据）")
	fmt.Println(strings.Repeat("-", 60))

	collectReq := mcp.CallToolRequest{}
	collectReq.Params.Name = "collect_reports"
	collectReq.Params.Arguments = map[string]interface{}{
		"start_month": "2025-01",
		"end_month":   "2025-01",
		"output_file": "test_login_verify.md",
	}

	collectResult, err := c.CallTool(ctx, collectReq)
	if err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
	} else {
		if collectResult.IsError {
			fmt.Println("⚠️  登录可能未成功（采集失败）")
			if len(collectResult.Content) > 0 {
				for _, content := range collectResult.Content {
					if textContent, ok := mcp.AsTextContent(content); ok {
						fmt.Printf("   错误: %s\n", textContent.Text)
					}
				}
			}
		} else {
			fmt.Println("✅ 登录验证成功！可以正常采集数据")
			if len(collectResult.Content) > 0 {
				for _, content := range collectResult.Content {
					if textContent, ok := mcp.AsTextContent(content); ok {
						fmt.Printf("   %s\n", textContent.Text)
					}
				}
			}
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 浏览器登录测试完成！")
	fmt.Println(strings.Repeat("=", 60))
}
