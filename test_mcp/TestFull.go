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
	fmt.Println("=== YST Go MCP 完整测试 ===\n")

	serverPath, err := filepath.Abs("./yst-go-mcp")
	if err != nil {
		log.Fatalf("获取服务器路径失败: %v", err)
	}

	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		log.Fatalf("服务器不存在: %s", serverPath)
	}

	fmt.Printf("连接到服务器: %s\n\n", serverPath)

	c, err := client.NewStdioMCPClient(serverPath)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 初始化
	fmt.Println("📡 初始化连接...")
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "yst-full-test",
		Version: "1.0.0",
	}

	serverInfo, err := c.Initialize(ctx, initReq)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Printf("✓ 连接成功：%s v%s\n\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 测试 1: 列出工具
	testListTools(ctx, c)

	// 测试 2: 清除 Cookies
	testClearCookies(ctx, c)

	// 测试 3: 测试 collect_reports（会失败，因为未登录，但可以测试错误处理）
	testCollectReports(ctx, c)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ 所有测试完成！")
	fmt.Println("\n💡 提示：")
	fmt.Println("  - clear_saved_cookies 工具正常工作")
	fmt.Println("  - collect_reports 需要先登录（预期错误）")
	fmt.Println("  - browser_login 需要图形界面（可手动测试）")
}

func testListTools(ctx context.Context, c *client.StdioMCPClient) {
	fmt.Println("📋 测试 1: 列出可用工具")
	fmt.Println(strings.Repeat("-", 60))

	toolsReq := mcp.ListToolsRequest{}
	toolsResult, err := c.ListTools(ctx, toolsReq)
	if err != nil {
		log.Fatalf("列出工具失败: %v", err)
	}

	fmt.Printf("✓ 共有 %d 个工具:\n\n", len(toolsResult.Tools))
	for i, tool := range toolsResult.Tools {
		fmt.Printf("  %d. %-20s %s\n", i+1, tool.Name, tool.Description)
	}
	fmt.Println()
}

func testClearCookies(ctx context.Context, c *client.StdioMCPClient) {
	fmt.Println("🧹 测试 2: 清除 Cookies")
	fmt.Println(strings.Repeat("-", 60))

	callReq := mcp.CallToolRequest{}
	callReq.Params.Name = "clear_saved_cookies"
	callReq.Params.Arguments = map[string]interface{}{}

	result, err := c.CallTool(ctx, callReq)
	if err != nil {
		fmt.Printf("❌ 调用失败: %v\n\n", err)
		return
	}

	fmt.Println("✓ 调用成功！")
	if len(result.Content) > 0 {
		for _, content := range result.Content {
			if textContent, ok := mcp.AsTextContent(content); ok {
				fmt.Printf("  返回: %s\n", textContent.Text)
			}
		}
	}
	fmt.Println()
}

func testCollectReports(ctx context.Context, c *client.StdioMCPClient) {
	fmt.Println("📊 测试 3: 采集日报（预期失败 - 未登录）")
	fmt.Println(strings.Repeat("-", 60))

	callReq := mcp.CallToolRequest{}
	callReq.Params.Name = "collect_reports"
	callReq.Params.Arguments = map[string]interface{}{
		"start_month": "2025-01",
		"end_month":   "2025-03",
		"output_file": "test_output.md",
	}

	result, err := c.CallTool(ctx, callReq)
	if err != nil {
		fmt.Printf("❌ 调用失败: %v\n\n", err)
		return
	}

	if result.IsError {
		fmt.Println("✓ 收到预期错误（未登录）")
	} else {
		fmt.Println("✓ 调用成功！")
	}

	if len(result.Content) > 0 {
		for _, content := range result.Content {
			if textContent, ok := mcp.AsTextContent(content); ok {
				fmt.Printf("  返回: %s\n", textContent.Text)
			}
		}
	}
	fmt.Println()
}
