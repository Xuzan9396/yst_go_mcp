package browser

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Xuzan9396/yst_go_mcp/internal/cookie"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	LoginURL  = "https://kpi.drojian.dev/site/login"
	TargetURL = "https://kpi.drojian.dev/report/report-daily/my-list"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36"
)

// Login 浏览器登录管理器
type Login struct {
	cookieManager *cookie.Manager
}

// NewLogin 创建浏览器登录管理器
func NewLogin() *Login {
	return &Login{
		cookieManager: cookie.NewManager(),
	}
}

// LaunchBrowserLogin 启动浏览器进行登录
func (l *Login) LaunchBrowserLogin(ctx context.Context, timeout int) error {
	log.Println("正在启动浏览器...")
	log.Printf("超时时间: %d 秒", timeout)

	// 设置 chromedp 选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent(UserAgent),
		chromedp.WindowSize(1920, 1080),
		chromedp.UserDataDir(l.cookieManager.GetBrowserProfileDir()),
	)

	// 创建上下文
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	browserCtx, browserCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer browserCancel()

	// 设置超时
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, time.Duration(timeout)*time.Second)
	defer timeoutCancel()

	// 导航到目标页面
	log.Printf("正在打开页面: %s", TargetURL)
	if err := chromedp.Run(timeoutCtx, chromedp.Navigate(TargetURL)); err != nil {
		log.Printf("首次访问出错（可能需要登录）: %v", err)
	}

	// 等待登录完成
	log.Println("⏳ 等待登录完成...")
	log.Println("提示：登录成功后，页面会跳转到日报列表页面")

	if err := l.waitForLoginSuccess(timeoutCtx, timeout); err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	log.Println("✓ 登录成功！正在提取 Cookie...")

	// 提取并保存 Cookies
	if err := chromedp.Run(timeoutCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		// 获取所有 Cookies
		cookiesData, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}

		// 转换 Cookie 格式
		var cookieList []cookie.Cookie
		for _, c := range cookiesData {
			cookieList = append(cookieList, cookie.Cookie{
				Name:   c.Name,
				Value:  c.Value,
				Domain: c.Domain,
				Path:   c.Path,
			})
		}

		// 保存 Cookie
		if err := l.cookieManager.SaveCookies(cookieList); err != nil {
			return fmt.Errorf("保存 Cookie 失败: %w", err)
		}

		log.Println("✓ Cookie 已保存")
		return nil
	})); err != nil {
		return fmt.Errorf("提取 Cookie 失败: %w", err)
	}

	log.Println("🎉 登录流程完成！现在可以使用 collect_reports 采集数据了")

	// 等待 3 秒让用户看到结果
	time.Sleep(3 * time.Second)

	return nil
}

// waitForLoginSuccess 等待登录成功
func (l *Login) waitForLoginSuccess(ctx context.Context, timeout int) error {
	checkInterval := 7 * time.Second
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	log.Printf("等待登录成功，超时时间: %d 秒，检查间隔: %v", timeout, checkInterval)

	for time.Now().Before(deadline) {
		var currentURL string
		if err := chromedp.Run(ctx, chromedp.Location(&currentURL)); err != nil {
			log.Printf("获取 URL 失败: %v", err)
			time.Sleep(checkInterval)
			continue
		}

		elapsed := int(time.Since(time.Now().Add(-time.Duration(timeout) * time.Second)).Seconds())
		log.Printf("[%ds] 检查中...", elapsed)
		log.Printf("  当前URL: %s", currentURL)

		// 检查是否已登录到系统
		if strings.Contains(currentURL, "kpi.drojian.dev") && !strings.Contains(currentURL, "accounts.google.com") {
			log.Println("  ✓ 已登录到系统！")

			// 如果不在目标页面，尝试跳转
			if !strings.Contains(currentURL, "my-list") && !strings.Contains(currentURL, "report-daily") {
				log.Println("  → 尝试跳转到日报页面...")
				if err := chromedp.Run(ctx, chromedp.Navigate(TargetURL)); err != nil {
					log.Printf("  ⚠ 跳转失败: %v", err)
				} else {
					time.Sleep(2 * time.Second)
					if err := chromedp.Run(ctx, chromedp.Location(&currentURL)); err == nil {
						log.Printf("  → 跳转后URL: %s", currentURL)
					}
				}
			}

			log.Printf("[%ds] ✓✓✓ 登录成功（已在系统内）！✓✓✓", elapsed)
			time.Sleep(2 * time.Second) // 等待 Cookie 完全保存
			return nil
		}

		log.Println("  ⏳ 等待跳转到目标页面...")
		log.Println("  （需要URL包含: my-list 或 report-daily）")
		time.Sleep(checkInterval)
	}

	return fmt.Errorf("登录超时（%d 秒）", timeout)
}
