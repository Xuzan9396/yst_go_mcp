# YST Go MCP - KPI 日报自动采集 MCP 服务器

基于 Go + chromedp 开发的日报数据自动采集工具，支持从 KPI 系统批量采集指定月份范围的日报数据。

## 功能特性

- ✅ **智能自动化**：自动检测登录状态，未登录时自动触发浏览器登录，一键完成采集
- ✅ **自动登录**：使用 chromedp 自动打开浏览器，完成 Google OAuth 登录
- ✅ **持久化会话**：登录一次长期有效，会话数据自动保存到 `~/.yst_go_mcp/`
- ✅ **批量采集**：支持一次性采集多个月份的日报数据
- ✅ **格式化输出**：自动生成结构化 Markdown 报告
- ✅ **跨平台支持**：macOS / Linux / Windows 全平台编译
- ✅ **灵活超时**：首次登录最长支持 6 分钟超时（默认），适应复杂的认证流程

## 环境要求

- **Go**: >= 1.25.0
- **浏览器**: Chrome/Chromium（chromedp 会自动查找）
- **网络**: 能访问 `https://kpi.drojian.dev`

## 快速开始

### 1. 编译

```bash
cd /Users/admin/go/empty/python/yst_go_mcp

# 编译
go build -o yst-go-mcp ./cmd/yst-go-mcp

# 测试运行
./yst-go-mcp
```

### 2. 配置到 Claude Desktop

编辑 Claude Desktop 的 MCP 配置文件：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

添加以下配置：

```json
{
  "mcpServers": {
    "yst-go-mcp": {
      "command": "/Users/admin/go/empty/python/yst_go_mcp/yst-go-mcp"
    }
  }
}
```

保存后，重启 Claude Desktop 即可生效。

## MCP 工具列表

| 工具名称 | 功能说明 | 参数 |
|---------|---------|------|
| `auto_collect_reports` | **🚀 自动采集日报（推荐）** - 自动检测登录状态，未登录时自动启动浏览器登录，登录成功后自动采集数据 | `start_month` (必需)、`end_month` (必需)、`output_file` (可选)、`login_timeout` (可选，默认 360 秒) |
| `collect_reports` | 采集日报数据（需要已登录） | `start_month` (必需)、`end_month` (必需)、`output_file` (可选) |
| `browser_login` | 启动浏览器进行登录 | `timeout` (可选，默认 360 秒) |
| `clear_saved_cookies` | 清除登录信息 | 无 |

## 使用示例

### 方式一：自动化采集（推荐）

**一键完成登录+采集**，推荐首次使用或 Cookie 过期时使用：

```
使用 auto_collect_reports 采集 2025-01 到 2025-03 的日报
```

系统会自动：
1. 检测是否已登录
2. 如果未登录，自动打开浏览器进行 Google 登录
3. 登录成功后自动采集数据

### 方式二：手动分步操作

**首次登录**：

```
使用 browser_login 登录
```

**采集数据**：

```
使用 collect_reports 采集 2025-01 到 2025-03 的日报
```

### 使用 Go 客户端测试

```bash
# 运行简单测试
go run test_client.go

# 运行完整测试
go run test_full.go

# 运行 STDIO 协议测试
go run test_stdio.go
```

## 项目结构

```
python/yst_go_mcp/
├── cmd/yst-go-mcp/
│   └── main.go                 # MCP Server 主程序
├── internal/
│   ├── cookie/
│   │   └── manager.go          # Cookie 管理
│   ├── browser/
│   │   └── login.go            # 浏览器登录 (chromedp)
│   └── collector/
│       └── collector.go        # 日报采集
├── test_client.go              # 基础测试客户端
├── test_full.go                # 完整测试客户端
├── test_stdio.go               # STDIO 协议测试
├── go.mod                      # Go 依赖
└── yst-go-mcp                  # 编译后的二进制 (15MB)
```

## 数据目录

### 开发环境
- Cookie: `./data/cookies.json`
- 浏览器配置: `./data/browser_profile/`
- 输出文件: `./output/new.md`

### 打包后
- Cookie: `~/.yst_go_mcp/data/cookies.json`
- 浏览器配置: `~/.yst_go_mcp/data/browser_profile/`
- 输出文件: `~/.yst_go_mcp/output/new.md`

## 技术栈

- **语言**: Go 1.25.0
- **MCP 库**: github.com/mark3labs/mcp-go v0.6.0
- **浏览器自动化**: github.com/chromedp/chromedp v0.11.2
- **HTML 解析**: github.com/PuerkitoBio/goquery v1.10.1

## 测试结果

```bash
$ go run test_stdio.go

=== 简单 STDIO 测试 ===

发送初始化请求...
收到: {"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05",...}}

发送列出工具请求...
收到: {"jsonrpc":"2.0","id":2,"result":{"tools":[...]}}

发送调用工具请求...
收到: {"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"✓ Cookie 和浏览器数据已清除"}]}}
```

✅ 所有核心功能测试通过！

## 开发者

- **版本**: 0.0.3
- **协议版本**: MCP 2024-11-05
- **仓库**: https://github.com/Xuzan9396/yst_go_mcp

## 从 Python 版本迁移

本项目从 Python + playwright 迁移而来，主要改进：
- ✅ 更快的启动速度（Go 编译二进制 vs Python 解释器）
- ✅ 更小的依赖（chromedp vs playwright）
- ✅ 单一二进制文件（无需 Python 环境）
- ✅ 更好的性能（原生 Go vs Python）

## 常见问题

### Q: 首次使用应该用哪个工具？
A: **推荐使用 `auto_collect_reports`**，它会自动检测登录状态，未登录时自动打开浏览器登录，然后自动采集数据，一步到位。

### Q: 登录超时时间太短怎么办？
A: 默认超时时间为 360 秒（6 分钟）。如需更长时间，可在调用时指定 `login_timeout` 参数，例如 `login_timeout: 600`（10 分钟）。

### Q: 浏览器启动失败？
A: 确保系统已安装 Chrome 或 Chromium 浏览器。

### Q: Cookie 过期？
A: 使用 `auto_collect_reports` 会自动检测并重新登录。也可以手动运行 `clear_saved_cookies` 清除后重新登录。

### Q: 如何查看日志？
A: 服务器日志会输出到标准错误输出（stderr）。

### Q: `auto_collect_reports` 和 `collect_reports` 有什么区别？
A:
- `auto_collect_reports`：智能自动化工具，会自动检测登录状态，未登录时自动触发登录，**推荐使用**
- `collect_reports`：简单采集工具，需要确保已登录，否则会报错

## License

MIT
