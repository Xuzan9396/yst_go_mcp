package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Xuzan9396/yst_go_mcp/internal/browser"
	"github.com/Xuzan9396/yst_go_mcp/internal/collector"
	"github.com/Xuzan9396/yst_go_mcp/internal/cookie"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建 MCP Server
	mcpServer := server.NewMCPServer(
		"YST Go MCP",
		"0.1.0",
	)

	// 注册工具
	registerTools(mcpServer)

	// 启动 STDIO Server
	log.Println("YST Go MCP Server 启动中...")
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// registerTools 注册所有工具
func registerTools(s *server.MCPServer) {
	// 1. browser_login 工具
	s.AddTool(
		mcp.NewTool("browser_login",
			mcp.WithDescription("启动浏览器进行登录"),
			mcp.WithNumber("timeout",
				mcp.DefaultNumber(300),
				mcp.Description("登录超时时间（秒），默认 300 秒"),
			),
		),
		handleBrowserLogin,
	)

	// 2. collect_reports 工具
	s.AddTool(
		mcp.NewTool("collect_reports",
			mcp.WithDescription("采集指定月份范围的日报数据"),
			mcp.WithString("start_month",
				mcp.Required(),
				mcp.Description("起始月份，格式 YYYY-MM (例如: 2025-01)"),
			),
			mcp.WithString("end_month",
				mcp.Required(),
				mcp.Description("结束月份，格式 YYYY-MM (例如: 2025-03)"),
			),
			mcp.WithString("output_file",
				mcp.Description("输出文件路径（可选，默认 ~/.yst_go_mcp/output/new.md）"),
			),
		),
		handleCollectReports,
	)

	// 3. clear_saved_cookies 工具
	s.AddTool(
		mcp.NewTool("clear_saved_cookies",
			mcp.WithDescription("清除已保存的 Cookie 和浏览器数据"),
		),
		handleClearCookies,
	)
}

// handleBrowserLogin 处理浏览器登录
func handleBrowserLogin(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	timeout := 300
	if val, ok := arguments["timeout"].(float64); ok {
		timeout = int(val)
	}

	log.Printf("browser_login 工具被调用，timeout=%d", timeout)

	loginManager := browser.NewLogin()
	if err := loginManager.LaunchBrowserLogin(context.Background(), timeout); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("登录失败: %v", err)), nil
	}

	return mcp.NewToolResultText("✅ 登录成功！Cookie 已保存，现在可以使用 collect_reports 采集数据了"), nil
}

// handleCollectReports 处理日报采集
func handleCollectReports(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	startMonth, ok := arguments["start_month"].(string)
	if !ok || startMonth == "" {
		return mcp.NewToolResultError("start_month 参数必须提供"), nil
	}

	endMonth, ok := arguments["end_month"].(string)
	if !ok || endMonth == "" {
		return mcp.NewToolResultError("end_month 参数必须提供"), nil
	}

	outputFile, _ := arguments["output_file"].(string)

	log.Printf("collect_reports 工具被调用: %s 到 %s, 输出: %s", startMonth, endMonth, outputFile)

	c := collector.NewCollector()
	result, err := c.Collect(startMonth, endMonth, outputFile)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("采集失败: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// handleClearCookies 处理清除 Cookies
func handleClearCookies(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	log.Println("clear_saved_cookies 工具被调用")

	manager := cookie.NewManager()
	if err := manager.ClearCookies(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("清除失败: %v", err)), nil
	}

	return mcp.NewToolResultText("✓ Cookie 和浏览器数据已清除"), nil
}
