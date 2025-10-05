# ⚡ 快速开始

## 1️⃣ 使用 npx（最简单）

```bash
npx -y @xuzan/yst-go-mcp
```

## 2️⃣ 配置到 Claude Desktop

**macOS**:
```bash
# 编辑配置文件
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**添加配置**:
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

**重启 Claude Desktop**

## 3️⃣ 开始使用

在 Claude Desktop 中：

```
使用 yst-go-mcp 采集 2025-01 到 2025-03 的日报
```

首次使用会打开浏览器，完成 Google 登录即可。

---

## 🔧 开发者

### 从源码运行

```bash
git clone https://github.com/Xuzan9396/yst_go_mcp.git
cd yst_go_mcp
make build
./yst-go-mcp
```

### 运行测试

```bash
cd test_mcp
./run_all_tests.sh
```

### 本地开发

```bash
# 编译
make build

# 编译所有平台
make build-all

# 运行测试
make test-full

# 格式化代码
make fmt
```

---

## 📖 更多文档

- [完整 README](README.md)
- [发布指南](PUBLISH.md)
- [chromedp 技术说明](docs/chromedp_info.md)
