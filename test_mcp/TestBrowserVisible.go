package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("=== 强制可见浏览器测试 ===\n")

	// 先激活 Chrome（macOS 专用）
	fmt.Println("→ 尝试激活 Chrome 窗口...")
	exec.Command("osascript", "-e", `tell application "System Events" to set frontmost of first process whose name is "Google Chrome" to true`).Run()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 使用用户的正常 Chrome 配置
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", false),
		chromedp.Flag("disable-dev-shm-usage", false),
		chromedp.WindowSize(1400, 900),
		// 尝试使用新窗口
		chromedp.Flag("new-window", true),
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	fmt.Println("→ 正在启动浏览器...")
	fmt.Println("→ 请注意观察屏幕上是否有新的 Chrome 窗口")
	fmt.Println()

	err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://kpi.drojian.dev/report/report-daily/my-list"),
		chromedp.Sleep(3*time.Second),
	)

	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Println("✅ 浏览器已启动并导航到登录页面")
	fmt.Println("📍 URL: https://kpi.drojian.dev/report/report-daily/my-list")
	fmt.Println()
	fmt.Println("⏰ 窗口将保持打开 60 秒")
	fmt.Println("   请检查:")
	fmt.Println("   1. Chrome 是否在屏幕上可见")
	fmt.Println("   2. 是否显示登录页面")
	fmt.Println("   3. 是否可以手动登录")
	fmt.Println()
	fmt.Println("等待中...")

	// 每 10 秒尝试激活一次窗口
	for i := 0; i < 6; i++ {
		time.Sleep(10 * time.Second)
		exec.Command("osascript", "-e", `tell application "Google Chrome" to activate`).Run()
		fmt.Printf("  %d/60 秒...\n", (i+1)*10)
	}

	fmt.Println("\n✓ 测试完成")
}
