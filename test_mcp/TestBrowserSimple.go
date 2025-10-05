package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("=== 简单浏览器测试 ===\n")
	fmt.Println("正在启动 Chrome 浏览器...")
	fmt.Println("⏰ 浏览器将保持打开 30 秒，请观察是否有窗口弹出\n")

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	// 使用最简单的配置
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false), // 明确设置为非无头模式
		chromedp.Flag("disable-gpu", false),
		chromedp.WindowSize(1280, 800),
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer taskCancel()

	fmt.Println("✓ 上下文创建成功")
	fmt.Println("→ 正在导航到测试页面...")

	var title string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.Sleep(2*time.Second), // 等待页面加载
		chromedp.Title(&title),
	)

	if err != nil {
		log.Fatalf("❌ 执行失败: %v", err)
	}

	fmt.Printf("\n✅ 成功！页面标题: %s\n", title)
	fmt.Println("\n请检查是否有 Chrome 窗口弹出...")
	fmt.Println("如果看到窗口，说明浏览器工作正常！")
	fmt.Println("\n等待 30 秒后自动关闭...")

	// 保持浏览器打开
	time.Sleep(30 * time.Second)

	fmt.Println("\n✓ 测试完成，浏览器即将关闭")
}
