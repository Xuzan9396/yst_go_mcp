package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		"0.0.3",
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
				mcp.DefaultNumber(360),
				mcp.Description("登录超时时间（秒），默认 360 秒（6 分钟）"),
			),
		),
		handleBrowserLogin,
	)

	// 2. collect_reports 工具
	//s.AddTool(
	//	mcp.NewTool("collect_reports",
	//		mcp.WithDescription("采集指定月份范围的日报数据"),
	//		mcp.WithString("start_month",
	//			mcp.Required(),
	//			mcp.Description("起始月份，格式 YYYY-MM (例如: 2025-01)"),
	//		),
	//		mcp.WithString("end_month",
	//			mcp.Required(),
	//			mcp.Description("结束月份，格式 YYYY-MM (例如: 2025-03)"),
	//		),
	//		mcp.WithString("output_file",
	//			mcp.Description("输出文件路径（可选，默认 ~/.yst_go_mcp/output/new.md）"),
	//		),
	//	),
	//	handleCollectReports,
	//)

	// 3. clear_saved_cookies 工具
	s.AddTool(
		mcp.NewTool("clear_saved_cookies",
			mcp.WithDescription("清除已保存的 Cookie 和浏览器数据"),
		),
		handleClearCookies,
	)

	// 4. auto_collect_reports 工具（自动化采集）
	s.AddTool(
		mcp.NewTool("auto_collect_reports",
			mcp.WithDescription("自动采集日报数据（如果未登录会自动启动浏览器登录）"),
			mcp.WithString("start_month",
				mcp.Required(),
				mcp.Description("起始月份，格式 YYYY-MM (例如: 2025-01)"),
			),
			mcp.WithString("end_month",
				mcp.Required(),
				mcp.Description("结束月份，格式 YYYY-MM (例如: 2025-03)"),
			),
			mcp.WithString("output_file",
				mcp.Description("输出文件路径（可选，默认,mac是下载目录 ~/Downloads/x月日报.md，windows是保留桌面 C:\\Users\\用户名\\Desktop\\x月日报.md）"),
			),
			mcp.WithNumber("login_timeout",
				mcp.DefaultNumber(360),
				mcp.Description("登录超时时间（秒），默认 360 秒（6 分钟）"),
			),
		),
		handleAutoCollectReports,
	)

	// 5. generate_summary_csv 工具（读取日报 MD 文件，输出内容供 AI 整理成 CSV）
	s.AddTool(
		mcp.NewTool("generate_summary_csv",
			mcp.WithDescription("读取日报详情 MD 文件内容，返回给 AI 模型整理生成 CSV 汇总表格"),
			mcp.WithString("md_file_path",
				mcp.Required(),
				mcp.Description("日报详情 MD 文件的完整路径"),
			),
		),
		handleGenerateSummaryCSV,
	)
}

// handleBrowserLogin 处理浏览器登录
func handleBrowserLogin(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	timeout := 360
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

// handleAutoCollectReports 处理自动采集（自动登录+采集）
func handleAutoCollectReports(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	// 解析参数
	startMonth, ok := arguments["start_month"].(string)
	if !ok || startMonth == "" {
		return mcp.NewToolResultError("start_month 参数必须提供"), nil
	}

	endMonth, ok := arguments["end_month"].(string)
	if !ok || endMonth == "" {
		return mcp.NewToolResultError("end_month 参数必须提供"), nil
	}

	outputFile, _ := arguments["output_file"].(string)

	loginTimeout := 360
	if val, ok := arguments["login_timeout"].(float64); ok {
		loginTimeout = int(val)
	}

	log.Printf("auto_collect_reports 工具被调用: %s 到 %s, 超时: %d 秒", startMonth, endMonth, loginTimeout)

	// 创建 cookie 管理器
	cookieManager := cookie.NewManager()
	c := collector.NewCollector()

	// 检查 cookie 是否存在且有效
	needLogin := false
	if !cookieManager.HasCookies() {
		log.Println("未找到 Cookie 文件，需要登录")
		needLogin = true
	} else {
		// 尝试加载 cookie 并检查登录状态
		if err := c.LoadSavedCookies(); err != nil {
			log.Printf("加载 Cookie 失败: %v，需要重新登录", err)
			needLogin = true
		} else if !c.CheckLoginStatus() {
			log.Println("Cookie 已过期，需要重新登录")
			needLogin = true
		}
	}

	// 如果需要登录
	if needLogin {
		log.Println("🔐 开始自动登录流程...")

		// 使用 channel 接收登录结果
		loginResult := make(chan error, 1)

		// 启动浏览器登录（异步）
		go func() {
			loginManager := browser.NewLogin()
			err := loginManager.LaunchBrowserLogin(context.Background(), loginTimeout)
			loginResult <- err
		}()

		// 定时检测登录状态
		checkInterval := 3 * time.Second
		deadline := time.Now().Add(time.Duration(loginTimeout) * time.Second)

		log.Printf("⏳ 等待登录完成（超时: %d 秒）...", loginTimeout)
		log.Println("💡 提示：请在浏览器中完成 Google 登录")

		for {
			select {
			case err := <-loginResult:
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("登录失败: %v", err)), nil
				}
				log.Println("✓ 登录成功！")
				goto LOGIN_SUCCESS

			case <-time.After(checkInterval):
				// 检查 cookie 文件是否已创建
				if cookieManager.HasCookies() {
					log.Println("✓ 检测到 Cookie 文件已创建")

					// 尝试验证登录状态
					newCollector := collector.NewCollector()
					if err := newCollector.LoadSavedCookies(); err == nil {
						if newCollector.CheckLoginStatus() {
							log.Println("✓ 登录状态验证成功！")
							c = newCollector
							goto LOGIN_SUCCESS
						}
					}
					log.Println("⏳ Cookie 文件存在但登录状态未就绪，继续等待...")
				}

				// 检查超时
				if time.Now().After(deadline) {
					return mcp.NewToolResultError(fmt.Sprintf("登录超时（%d 秒）", loginTimeout)), nil
				}

				elapsed := int(time.Since(deadline.Add(-time.Duration(loginTimeout) * time.Second)).Seconds())
				log.Printf("⏳ [%ds/%ds] 等待登录中...", elapsed, loginTimeout)
			}
		}

	LOGIN_SUCCESS:
		log.Println("🎉 登录流程完成！")

		// 如果之前没有加载过 cookie，现在加载
		if err := c.LoadSavedCookies(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("加载 Cookie 失败: %v", err)), nil
		}
	} else {
		log.Println("✓ Cookie 有效，跳过登录")
	}

	// 开始采集数据
	log.Printf("📊 开始采集日报数据: %s 到 %s", startMonth, endMonth)
	result, err := c.Collect(startMonth, endMonth, outputFile)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("采集失败: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// handleGenerateSummaryCSV 处理读取日报 MD 并输出内容供 AI 整理
func handleGenerateSummaryCSV(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	mdFilePath, ok := arguments["md_file_path"].(string)
	if !ok || mdFilePath == "" {
		return mcp.NewToolResultError("md_file_path 参数必须提供"), nil
	}

	log.Printf("generate_summary_csv 工具被调用: %s", mdFilePath)

	c := collector.NewCollector()
	content, csvPath, err := c.ReadMarkdownForSummary(mdFilePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("读取 MD 文件失败: %v", err)), nil
	}

	result := fmt.Sprintf(`📄 已读取日报详情文件: %s

请根据以下日报内容，整理生成 CSV 格式的月度汇总表格，包含以下列：
- 序号
- 主要工作任务
- 权重
- 任务成果情况

生成的 CSV 文件应保存到: %s

日报内容如下：
---
%s
---

请分析日报内容，提取主要工作任务，并生成符合格式的 CSV 文件。`, mdFilePath, csvPath, content)

	return mcp.NewToolResultText(result), nil
}
