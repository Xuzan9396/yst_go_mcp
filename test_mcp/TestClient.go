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
	fmt.Println("=== YST Go MCP 客户端测试 ===\n")

	// 获取服务器二进制路径
	serverPath, err := filepath.Abs("./yst-go-mcp")
	if err != nil {
		log.Fatalf("获取服务器路径失败: %v", err)
	}

	// 检查服务器是否存在
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 初始化连接
	fmt.Println("📡 初始化连接...")
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "yst-test-client",
		Version: "1.0.0",
	}
	initReq.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initReq)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Printf("✓ 连接成功！\n")
	fmt.Printf("  服务器名称: %s\n", serverInfo.ServerInfo.Name)
	fmt.Printf("  服务器版本: %s\n", serverInfo.ServerInfo.Version)
	fmt.Printf("  协议版本: %s\n\n", serverInfo.ProtocolVersion)

	// 列出可用工具
	fmt.Println("📋 列出可用工具...")
	toolsReq := mcp.ListToolsRequest{}
	toolsResult, err := c.ListTools(ctx, toolsReq)
	if err != nil {
		log.Fatalf("列出工具失败: %v", err)
	}

	fmt.Printf("✓ 共有 %d 个工具:\n\n", len(toolsResult.Tools))
	for i, tool := range toolsResult.Tools {
		fmt.Printf("%d. %s\n", i+1, tool.Name)
		fmt.Printf("   描述: %s\n", tool.Description)
		if tool.InputSchema.Properties != nil {
			fmt.Printf("   参数: %v\n", getPropertyNames(tool.InputSchema.Properties))
		}
		fmt.Println()
	}

	// 测试调用工具 - clear_saved_cookies（不需要参数）
	fmt.Println("🧪 测试调用 clear_saved_cookies 工具...")
	callReq := mcp.CallToolRequest{}
	callReq.Params.Name = "clear_saved_cookies"
	callReq.Params.Arguments = map[string]interface{}{}

	result, err := c.CallTool(ctx, callReq)
	if err != nil {
		log.Printf("❌ 工具调用失败: %v", err)
	} else {
		fmt.Println("✓ 工具调用成功！")
		if len(result.Content) > 0 {
			for _, content := range result.Content {
				if textContent, ok := mcp.AsTextContent(content); ok {
					fmt.Printf("   返回: %s\n", textContent.Text)
				}
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("测试完成！\n")
	fmt.Println("你可以尝试的其他测试：")
	fmt.Println("1. browser_login - 打开浏览器登录（需要浏览器）")
	fmt.Println("2. collect_reports - 采集日报（需要先登录）")
	fmt.Println("   示例参数:")
	fmt.Println("   {")
	fmt.Println("     \"start_month\": \"2025-01\",")
	fmt.Println("     \"end_month\": \"2025-03\",")
	fmt.Println("     \"output_file\": \"test_output.md\"")
	fmt.Println("   }")
}

// getPropertyNames 获取属性名列表
func getPropertyNames(props map[string]interface{}) []string {
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	return names
}
