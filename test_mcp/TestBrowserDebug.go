package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("=== Chrome 浏览器启动调试 ===\n")

	// 测试 1: 检查 Chrome 是否可用
	fmt.Println("测试 1: 检查系统 Chrome")
	fmt.Println("-----------------------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 尝试使用系统 Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	// 启用调试日志
	browserCtx, browserCancel := chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Printf))
	defer browserCancel()

	fmt.Println("正在启动 Chrome 浏览器...")
	fmt.Println("如果浏览器成功启动，你应该能看到一个窗口弹出")
	fmt.Println()

	// 尝试导航到百度（简单测试）
	var title string
	err := chromedp.Run(browserCtx,
		chromedp.Navigate("https://www.baidu.com"),
		chromedp.Title(&title),
	)

	if err != nil {
		fmt.Printf("❌ 启动失败: %v\n\n", err)

		// 尝试无头模式
		fmt.Println("测试 2: 尝试无头模式")
		fmt.Println("-----------------------------------")

		opts2 := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
		)

		allocCtx2, allocCancel2 := chromedp.NewExecAllocator(context.Background(), opts2...)
		defer allocCancel2()

		browserCtx2, browserCancel2 := chromedp.NewContext(allocCtx2)
		defer browserCancel2()

		err2 := chromedp.Run(browserCtx2,
			chromedp.Navigate("https://www.baidu.com"),
			chromedp.Title(&title),
		)

		if err2 != nil {
			fmt.Printf("❌ 无头模式也失败: %v\n", err2)
		} else {
			fmt.Printf("✓ 无头模式成功！页面标题: %s\n", title)
			fmt.Println("\n⚠️ 问题：有头模式失败，但无头模式成功")
			fmt.Println("   可能原因：")
			fmt.Println("   1. 显示服务器配置问题（如果是 SSH 连接）")
			fmt.Println("   2. Chrome 版本或权限问题")
		}
		return
	}

	fmt.Printf("✅ 浏览器启动成功！\n")
	fmt.Printf("   页面标题: %s\n", title)
	fmt.Println("\n请检查是否有浏览器窗口弹出...")
	fmt.Println("窗口会在 5 秒后关闭")

	time.Sleep(5 * time.Second)
}
