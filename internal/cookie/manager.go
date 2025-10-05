package cookie

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Cookie 表示浏览器 Cookie
type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

// Manager Cookie 管理器
type Manager struct {
	cookieFile string
}

// NewManager 创建 Cookie 管理器
func NewManager() *Manager {
	dataDir := getDataDir()
	return &Manager{
		cookieFile: filepath.Join(dataDir, "cookies.json"),
	}
}

// getDataDir 获取数据目录
func getDataDir() string {
	// 检查是否是打包后的可执行文件
	exePath, err := os.Executable()
	if err == nil {
		// 打包后使用用户主目录
		homeDir, err := os.UserHomeDir()
		if err == nil {
			dataDir := filepath.Join(homeDir, ".yst_go_mcp", "data")
			// 检查是否在标准安装目录下运行
			if !filepath.HasPrefix(exePath, "/tmp") && !filepath.HasPrefix(exePath, os.TempDir()) {
				os.MkdirAll(dataDir, 0755)
				return dataDir
			}
		}
	}

	// 开发模式：使用项目目录
	return "data"
}

// SaveCookies 保存 Cookies 到文件
func (m *Manager) SaveCookies(cookies []Cookie) error {
	// 确保目录存在
	dir := filepath.Dir(m.cookieFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 Cookie 失败: %w", err)
	}

	if err := os.WriteFile(m.cookieFile, data, 0600); err != nil {
		return fmt.Errorf("保存 Cookie 文件失败: %w", err)
	}

	return nil
}

// LoadCookies 从文件加载 Cookies
func (m *Manager) LoadCookies() ([]Cookie, error) {
	data, err := os.ReadFile(m.cookieFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取 Cookie 文件失败: %w", err)
	}

	var cookies []Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("解析 Cookie 文件失败: %w", err)
	}

	return cookies, nil
}

// HasCookies 检查是否有保存的 Cookies
func (m *Manager) HasCookies() bool {
	_, err := os.Stat(m.cookieFile)
	return err == nil
}

// ClearCookies 清除保存的 Cookies
func (m *Manager) ClearCookies() error {
	if !m.HasCookies() {
		return nil
	}

	if err := os.Remove(m.cookieFile); err != nil {
		return fmt.Errorf("删除 Cookie 文件失败: %w", err)
	}

	// 同时清除浏览器配置文件目录
	browserProfileDir := filepath.Join(filepath.Dir(m.cookieFile), "browser_profile")
	if _, err := os.Stat(browserProfileDir); err == nil {
		if err := os.RemoveAll(browserProfileDir); err != nil {
			return fmt.Errorf("删除浏览器配置目录失败: %w", err)
		}
	}

	return nil
}

// GetCookieFile 获取 Cookie 文件路径
func (m *Manager) GetCookieFile() string {
	return m.cookieFile
}

// GetBrowserProfileDir 获取浏览器配置目录
func (m *Manager) GetBrowserProfileDir() string {
	return filepath.Join(filepath.Dir(m.cookieFile), "browser_profile")
}
