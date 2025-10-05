# @xuzan/yst-go-mcp

YST Go MCP - KPI 日报自动采集 MCP 服务器

## 快速开始

### 使用 npx（推荐）

```bash
npx -y @xuzan/yst-go-mcp
```

### 安装到项目

```bash
npm install @xuzan/yst-go-mcp
```

## 配置到 Claude Desktop

编辑配置文件：
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

添加配置：

```json
{
  "mcpServers": {
    "yst-go-mcp": {
      "command": "npx",
      "args": ["-y", "@xuzan/yst-go-mcp"]
    }
  }
}
```

重启 Claude Desktop 即可。

## 功能特性

- ✅ **自动登录** - 使用 chromedp 自动化浏览器登录
- ✅ **持久化会话** - 登录一次长期有效
- ✅ **批量采集** - 支持多个月份的日报数据采集
- ✅ **格式化输出** - 自动生成 Markdown 报告
- ✅ **跨平台** - macOS / Linux / Windows

## 使用方法

在 Claude Desktop 中：

```
使用 yst-go-mcp 采集 2025-01 到 2025-03 的日报
```

首次使用会自动打开浏览器进行 Google 登录。

## 支持平台

- macOS (Intel / Apple Silicon)
- Linux (x64)
- Windows (x64 / ARM64)

## 依赖

- Chrome 或 Chromium 浏览器
- 网络可访问 `https://kpi.drojian.dev`

## 更多信息

- 📖 [完整文档](https://github.com/Xuzan9396/yst_go_mcp)
- 🐛 [问题反馈](https://github.com/Xuzan9396/yst_go_mcp/issues)
- 📝 [更新日志](https://github.com/Xuzan9396/yst_go_mcp/releases)

## License

MIT
