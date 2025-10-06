package collector

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Xuzan9396/yst_go_mcp/internal/cookie"
)

const (
	BaseURL       = "https://kpi.drojian.dev"
	ReportListURL = BaseURL + "/report/report-daily/my-list"
	UserAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36"
)

// Collector 日报采集器
type Collector struct {
	client        *http.Client
	cookieManager *cookie.Manager
}

// Report 日报信息
type Report struct {
	Text string
	Link string
}

// NewCollector 创建日报采集器
func NewCollector() *Collector {
	jar, _ := cookiejar.New(nil)
	return &Collector{
		client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		cookieManager: cookie.NewManager(),
	}
}

// LoadSavedCookies 加载保存的 Cookies
func (c *Collector) LoadSavedCookies() error {
	cookies, err := c.cookieManager.LoadCookies()
	if err != nil {
		return fmt.Errorf("加载 Cookie 失败: %w", err)
	}

	if len(cookies) == 0 {
		return fmt.Errorf("没有保存的 Cookie")
	}

	// 转换为 http.Cookie 格式
	baseURL, _ := url.Parse(BaseURL)
	var httpCookies []*http.Cookie
	for _, ck := range cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:   ck.Name,
			Value:  ck.Value,
			Domain: ck.Domain,
			Path:   ck.Path,
		})
	}

	c.client.Jar.SetCookies(baseURL, httpCookies)
	return nil
}

// CheckLoginStatus 检查登录状态
func (c *Collector) CheckLoginStatus() bool {
	resp, err := c.client.Get(ReportListURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 如果返回 200 且不是重定向到登录页，说明已登录
	return resp.StatusCode == http.StatusOK && !strings.Contains(resp.Request.URL.String(), "login")
}

// FetchMonthReports 获取指定月份的日报列表
func (c *Collector) FetchMonthReports(month string) ([]Report, error) {
	reportURL := fmt.Sprintf("%s?month=%s", ReportListURL, month)

	req, err := http.NewRequest("GET", reportURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh-TW;q=0.9,zh;q=0.8,en;q=0.7")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	var reports []Report
	doc.Find("#report_list li").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		link, _ := s.Find("a").Attr("href")

		if text != "" {
			reports = append(reports, Report{
				Text: text,
				Link: link,
			})
		}
	})

	return reports, nil
}

// GenerateMonthRange 生成月份范围列表
func (c *Collector) GenerateMonthRange(startMonth, endMonth string) ([]string, error) {
	start, err := time.Parse("2006-01", startMonth)
	if err != nil {
		return nil, fmt.Errorf("起始月份格式错误: %w", err)
	}

	end, err := time.Parse("2006-01", endMonth)
	if err != nil {
		return nil, fmt.Errorf("结束月份格式错误: %w", err)
	}

	var months []string
	current := start
	for !current.After(end) {
		months = append(months, current.Format("2006-01"))
		current = current.AddDate(0, 1, 0)
	}

	return months, nil
}

// Collect 采集指定月份范围的日报并保存
func (c *Collector) Collect(startMonth, endMonth, outputFile string) (string, error) {
	// 处理输出文件路径
	if outputFile == "" {
		outputFile = c.getDefaultOutputFile()
	} else if !filepath.IsAbs(outputFile) {
		// 相对路径转换为绝对路径
		outputFile = filepath.Join(c.getDefaultOutputDir(), outputFile)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 加载已保存的 Cookie
	if c.cookieManager.HasCookies() {
		if err := c.LoadSavedCookies(); err != nil {
			return "", fmt.Errorf("加载 Cookie 失败: %w", err)
		}
	}

	// 检查登录状态
	if !c.CheckLoginStatus() {
		return "", fmt.Errorf("未登录或登录已过期，请先使用 browser_login 工具登录")
	}

	// 生成月份范围
	months, err := c.GenerateMonthRange(startMonth, endMonth)
	if err != nil {
		return "", fmt.Errorf("生成月份范围失败: %w", err)
	}

	// 采集所有月份的数据
	allReports := make(map[string][]Report)
	for _, month := range months {
		log.Printf("正在采集 %s 月份日报...", month)
		reports, err := c.FetchMonthReports(month)
		if err != nil {
			log.Printf("采集 %s 月份失败: %v", month, err)
			continue
		}
		allReports[month] = reports
		log.Printf("  ✓ 采集到 %d 条日报", len(reports))
	}

	// 生成 Markdown 文件
	if err := c.generateMarkdown(allReports, outputFile); err != nil {
		return "", fmt.Errorf("生成 Markdown 失败: %w", err)
	}

	totalCount := 0
	for _, reports := range allReports {
		totalCount += len(reports)
	}

	return fmt.Sprintf("✓ 采集完成！共采集 %d 个月份，%d 条日报，已保存到 %s",
		len(months), totalCount, outputFile), nil
}

// generateMarkdown 生成 Markdown 文件
func (c *Collector) generateMarkdown(allReports map[string][]Report, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	// 写入标题
	fmt.Fprintln(f, "# YST 日报整理")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "生成时间：%s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 按月份排序输出
	var months []string
	for month := range allReports {
		months = append(months, month)
	}
	// 简单排序
	for i := 0; i < len(months); i++ {
		for j := i + 1; j < len(months); j++ {
			if months[i] > months[j] {
				months[i], months[j] = months[j], months[i]
			}
		}
	}

	for _, month := range months {
		reports := allReports[month]
		fmt.Fprintf(f, "## %s 月份日报 (%d 条)\n\n", month, len(reports))

		if len(reports) == 0 {
			fmt.Fprintln(f, "*暂无数据*\n")
			continue
		}

		for i, report := range reports {
			fmt.Fprintf(f, "### %d. %s\n\n", i+1, report.Text)
			if report.Link != "" {
				fmt.Fprintf(f, "链接：%s\n\n", report.Link)
			}
			fmt.Fprintln(f, "---\n")
		}
	}

	return nil
}

// getDefaultOutputDir 获取默认输出目录
func (c *Collector) getDefaultOutputDir() string {
	// 优先使用当前工作目录（AI 客户端的项目目录）
	if cwd, err := os.Getwd(); err == nil {
		// 检查当前目录是否可写
		testFile := filepath.Join(cwd, ".yst_test_write")
		if f, err := os.Create(testFile); err == nil {
			f.Close()
			os.Remove(testFile)
			log.Printf("使用当前工作目录作为输出目录: %s", cwd)
			return cwd
		}
	}

	// 如果当前目录不可写，使用用户主目录
	homeDir, err := os.UserHomeDir()
	if err == nil {
		outputDir := filepath.Join(homeDir, ".yst_go_mcp", "output")
		os.MkdirAll(outputDir, 0755)
		log.Printf("使用用户主目录作为输出目录: %s", outputDir)
		return outputDir
	}

	// 最后的备选方案
	log.Println("使用相对路径 output 作为输出目录")
	return "output"
}

// getDefaultOutputFile 获取默认输出文件
func (c *Collector) getDefaultOutputFile() string {
	return filepath.Join(c.getDefaultOutputDir(), "日报详情.md")
}

// ReadMarkdownForSummary 读取日报 MD 文件，返回内容和 CSV 保存路径
func (c *Collector) ReadMarkdownForSummary(mdFilePath string) (string, string, error) {
	// 读取 MD 文件内容
	content, err := os.ReadFile(mdFilePath)
	if err != nil {
		return "", "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 生成 CSV 文件路径（同目录）
	dir := filepath.Dir(mdFilePath)
	baseName := filepath.Base(mdFilePath)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// 提取月份信息（假设文件名包含月份，如 "1月日报.md"）
	month := strings.Split(baseName, "月")[0]
	csvPath := filepath.Join(dir, month+"月汇总总结.csv")

	return string(content), csvPath, nil
}
