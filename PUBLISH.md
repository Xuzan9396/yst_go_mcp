# 🚀 快速发布指南

## 一键发布流程（5 分钟）

### 第 1 步：准备 NPM Token

1. 登录 NPM：https://www.npmjs.com
2. 创建 Token：https://www.npmjs.com/settings/{你的用户名}/tokens
   - 点击 "Generate New Token"
   - 选择 "Automation"
   - 复制 token

3. 配置到 GitHub：
   - 访问：https://github.com/Xuzan9396/yst_go_mcp/settings/secrets/actions
   - 点击 "New repository secret"
   - Name: `NPM_TOKEN`
   - Value: 粘贴刚才的 token
   - 点击 "Add secret"

### 第 2 步：更新版本号

编辑 `npm/package.json`：

```json
{
  "version": "0.1.1"  ← 修改这里
}
```

### 第 3 步：提交并创建 Tag

```bash
# 1. 提交更改
git add .
git commit -m "chore: bump version to v0.1.1"
git push origin main

# 2. 创建并推送 tag
git tag v0.1.1
git push origin v0.1.1
```

### 第 4 步：等待自动构建

访问 GitHub Actions 查看进度：
```
https://github.com/Xuzan9396/yst_go_mcp/actions
```

**自动完成的任务**：
- ✅ 编译 macOS arm64
- ✅ 编译 macOS amd64
- ✅ 编译 Linux amd64
- ✅ 编译 Windows amd64
- ✅ 编译 Windows arm64
- ✅ 创建 GitHub Release
- ✅ 上传 5 个二进制文件
- ✅ 发布到 NPM

**预计耗时**: 5-10 分钟

### 第 5 步：验证发布

```bash
# 验证 GitHub Release
open https://github.com/Xuzan9396/yst_go_mcp/releases

# 验证 NPM
npm view @xuzan/yst-go-mcp

# 测试 npx
npx -y @xuzan/yst-go-mcp
```

---

## 📋 发布前检查清单

复制此清单，逐项检查：

```
[ ] 所有代码已提交
[ ] 所有测试通过（make test-full）
[ ] 更新了版本号（npm/package.json）
[ ] NPM_TOKEN 已配置到 GitHub Secrets
[ ] 当前在 main 分支
[ ] 已拉取最新代码（git pull）
```

---

## 🎯 快速命令

```bash
# 测试编译
make build

# 测试所有平台
make build-all

# 运行测试
cd test_mcp && ./run_all_tests.sh

# 发布流程
git add . && git commit -m "chore: bump version to v0.1.1"
git push origin main
git tag v0.1.1 && git push origin v0.1.1
```

---

## 📦 发布后使用

用户可通过以下方式使用：

### 方式 1：npx（推荐）

```bash
npx -y @xuzan/yst-go-mcp
```

### 方式 2：Claude Desktop

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

### 方式 3：手动下载

```bash
# macOS Apple Silicon
curl -L https://github.com/Xuzan9396/yst_go_mcp/releases/latest/download/yst-go-mcp-darwin-arm64 -o yst-go-mcp
chmod +x yst-go-mcp
./yst-go-mcp
```

---

## ⚠️ 常见问题

### Q: GitHub Actions 构建失败？
A: 检查 Actions 日志，通常是依赖问题或 Go 版本问题

### Q: NPM 发布失败？
A: 检查 NPM_TOKEN 是否正确配置，版本号是否重复

### Q: 如何回滚版本？
A: 删除 tag 并在 NPM 上废弃版本
```bash
git tag -d v0.1.1
git push origin :refs/tags/v0.1.1
npm deprecate @xuzan/yst-go-mcp@0.1.1 "请使用最新版本"
```

---

## 📞 需要帮助？

- 📖 详细发布指南：[RELEASE.md](RELEASE.md)
- 🐛 问题反馈：https://github.com/Xuzan9396/yst_go_mcp/issues
