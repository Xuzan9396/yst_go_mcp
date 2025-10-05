# YST Go MCP - KPI 日报自动采集 MCP 服务器

基于 Go + chromedp 开发的日报数据自动采集工具，支持从 KPI 系统批量采集指定月份范围的日报数据。

## 功能特性

- ✅ **自动登录**：使用 chromedp 自动打开浏览器，完成 Google OAuth 登录
- ✅ **持久化会话**：登录一次长期有效，会话数据自动保存到 `~/.yst_go_mcp/`
- ✅ **批量采集**：支持一次性采集多个月份的日报数据
- ✅ **格式化输出**：自动生成结构化 Markdown 报告
- ✅ **跨平台支持**：macOS / Linux / Windows 全平台编译

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
| `browser_login` | 启动浏览器进行登录 | `timeout` (可选，默认 300 秒) |
| `collect_reports` | 采集日报数据 | `start_month` (必需)、`end_month` (必需)、`output_file` (可选) |
| `clear_saved_cookies` | 清除登录信息 | 无 |

## 使用示例

### 在 Claude Desktop 中使用

```
使用 yst-go-mcp 采集 2025-01 到 2025-03 的日报
```

首次使用会自动打开浏览器，完成 Google 登录后，系统会自动采集数据。

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

- **版本**: 0.1.0
- **协议版本**: MCP 2024-11-05
- **仓库**: https://github.com/Xuzan9396/yst_go_mcp

## 从 Python 版本迁移

本项目从 Python + playwright 迁移而来，主要改进：
- ✅ 更快的启动速度（Go 编译二进制 vs Python 解释器）
- ✅ 更小的依赖（chromedp vs playwright）
- ✅ 单一二进制文件（无需 Python 环境）
- ✅ 更好的性能（原生 Go vs Python）

## 常见问题

### Q: 浏览器启动失败？
A: 确保系统已安装 Chrome 或 Chromium 浏览器。

### Q: Cookie 过期？
A: 运行 `clear_saved_cookies` 清除后重新登录。

### Q: 如何查看日志？
A: 服务器日志会输出到标准错误输出（stderr）。

## License

MIT
