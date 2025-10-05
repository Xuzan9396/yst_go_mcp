# chromedp 技术说明

## 底层依赖

chromedp 使用 **Chrome DevTools Protocol (CDP)**，这是一个标准协议。

### 支持的浏览器

| 浏览器 | 支持 | 说明 |
|--------|------|------|
| Google Chrome | ✅ 最佳 | 推荐使用，兼容性最好 |
| Chromium | ✅ 推荐 | 开源版本，功能相同 |
| Microsoft Edge | ✅ 支持 | 基于 Chromium，可用 |
| Brave | ✅ 支持 | 基于 Chromium，可用 |
| Opera | ⚠️ 部分 | 基于 Chromium，部分功能可用 |
| Firefox | ❌ 不支持 | 使用不同的协议 |

## chromedp 与 Playwright 对比

| 特性 | chromedp | Playwright |
|------|----------|------------|
| **语言** | Go 原生 | Node.js/Python/Java/.NET |
| **依赖** | 系统 Chrome/Chromium | 自带浏览器驱动（需下载） |
| **安装大小** | ~15MB (仅库) | ~300MB (含浏览器) |
| **启动速度** | 极快 (<100ms) | 慢 (~2秒) |
| **浏览器查找** | 自动查找系统浏览器 | 需要 `playwright install` |
| **持久化会话** | ✅ UserDataDir | ✅ launch_persistent_context |
| **多浏览器** | ❌ 仅 Chromium 系 | ✅ Chrome/Firefox/Safari |
| **API 复杂度** | 中等 | 简单 |
| **性能** | 优秀 | 良好 |

## 安装要求

### macOS
```bash
# 检查是否已安装 Chrome
which google-chrome-stable || \
ls "/Applications/Google Chrome.app"

# 或安装 Chromium
brew install --cask chromium
```

### Linux
```bash
# Debian/Ubuntu
sudo apt install chromium-browser

# 或 Chrome
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome-stable_current_amd64.deb
```

### Windows
```powershell
# 检查 Chrome
Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe"

# 或下载安装
# https://www.google.com/chrome/
```

## chromedp 自动查找浏览器的位置

chromedp 会按顺序查找：

1. **环境变量** `CHROMEDP_BROWSER_PATH`
2. **系统 Chrome**：
   - macOS: `/Applications/Google Chrome.app`
   - Linux: `/usr/bin/google-chrome-stable`
   - Windows: `C:\Program Files\Google\Chrome\Application\chrome.exe`
3. **系统 Chromium**：
   - macOS: `/Applications/Chromium.app`
   - Linux: `/usr/bin/chromium-browser`
   - Windows: `C:\Program Files\Chromium\Application\chromium.exe`

## 使用 UserDataDir 的好处

```go
chromedp.UserDataDir("/path/to/profile")
```

### 优势
1. ✅ **登录状态持久化** - 登录一次，永久有效
2. ✅ **无需每次登录** - 浏览器记住所有登录信息
3. ✅ **Cookie 自动管理** - 不需要手动导出/导入
4. ✅ **完整浏览器状态** - 包括扩展、设置、缓存
5. ✅ **支持多会话** - 不同目录 = 不同身份

### 注意事项
⚠️ **安全性**：UserDataDir 包含敏感信息，需要妥善保管
⚠️ **磁盘空间**：单个 profile 可能占用几十到几百 MB
⚠️ **并发限制**：同一个 UserDataDir 不能被多个进程同时使用

## 当前项目配置

我们的代码已经正确配置了持久化会话：

```go
// internal/browser/login.go
chromedp.UserDataDir(l.cookieManager.GetBrowserProfileDir())
```

**存储位置**：
- 开发：`./data/browser_profile/`
- 生产：`~/.yst_go_mcp/data/browser_profile/`

**工作流程**：
1. 第一次运行 `browser_login` → 打开浏览器 → 手动登录
2. 浏览器自动保存所有状态到 UserDataDir
3. 第二次运行 → chromedp 加载 UserDataDir → **已登录状态**
4. 无需再次登录，直接采集数据

## 为什么窗口不可见？

可能原因：
1. **窗口在后台** - 检查 Dock/任务栏
2. **窗口在其他桌面** - macOS Mission Control (F3)
3. **SSH 连接无显示** - 需要本地图形环境
4. **headless 模式** - 检查代码是否设置了 headless=true

## 强制窗口显示的方法

```go
// macOS 激活 Chrome 窗口
exec.Command("osascript", "-e",
  `tell application "Google Chrome" to activate`).Run()

// 或使用 --new-window 标志
chromedp.Flag("new-window", true)
```
