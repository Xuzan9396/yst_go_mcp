# YST Go MCP 发布指南

## 发布流程

### 1. 准备发布

```bash
# 确保所有更改已提交
git status

# 确保在 main 分支
git checkout main

# 拉取最新代码
git pull origin main

# 运行测试
make test-full
```

### 2. 更新版本号

编辑 `npm/package.json`，更新版本号：

```json
{
  "version": "0.1.1"  // 更新这里
}
```

### 3. 创建 Git Tag

```bash
# 创建 tag（版本号与 package.json 一致）
git tag v0.1.1

# 推送 tag 到远程
git push origin v0.1.1
```

### 4. 自动构建和发布

GitHub Actions 会自动：
1. ✅ 编译 5 个平台的二进制文件
2. ✅ 创建 GitHub Release
3. ✅ 上传二进制到 Release
4. ✅ 发布到 NPM

**查看进度**：
https://github.com/Xuzan9396/yst_go_mcp/actions

### 5. 验证发布

#### 验证 GitHub Release

```bash
# 访问 Release 页面
open https://github.com/Xuzan9396/yst_go_mcp/releases
```

检查：
- ✅ 5 个平台的二进制文件都已上传
- ✅ Release 说明正确显示
- ✅ 文件大小正常（约 15MB）

#### 验证 NPM 发布

```bash
# 查看 NPM 包信息
npm view @xuzan/yst-go-mcp

# 测试 npx 安装
npx -y @xuzan/yst-go-mcp@latest
```

### 6. 测试各平台

#### macOS (本地测试)

```bash
# 下载并测试
curl -L https://github.com/Xuzan9396/yst_go_mcp/releases/latest/download/yst-go-mcp-darwin-arm64 -o yst-go-mcp
chmod +x yst-go-mcp
./yst-go-mcp
```

#### NPM 包测试

```bash
# 测试 npx
npx -y @xuzan/yst-go-mcp

# 测试 npm install
mkdir test-install && cd test-install
npm init -y
npm install @xuzan/yst-go-mcp
npx yst-go-mcp
```

## GitHub Secrets 配置

需要在 GitHub 仓库设置中配置以下 Secrets：

1. **NPM_TOKEN** - NPM 发布令牌
   - 访问 https://www.npmjs.com/settings/{username}/tokens
   - 创建 "Automation" 类型的 token
   - 复制 token 到 GitHub Secrets

2. **GITHUB_TOKEN** - 自动提供，无需配置

### 配置步骤

1. 访问仓库设置
   ```
   https://github.com/Xuzan9396/yst_go_mcp/settings/secrets/actions
   ```

2. 点击 "New repository secret"

3. 添加 NPM_TOKEN：
   - Name: `NPM_TOKEN`
   - Value: 你的 NPM token

## 本地构建测试

### 编译所有平台

```bash
make build-all
```

输出在 `build/` 目录：
```
build/
├── yst-go-mcp-darwin-amd64
├── yst-go-mcp-darwin-arm64
├── yst-go-mcp-linux-amd64
├── yst-go-mcp-windows-amd64.exe
└── yst-go-mcp-windows-arm64.exe
```

### 测试 NPM 安装脚本

```bash
cd npm
node install.js
```

## 版本管理规范

使用语义化版本：`MAJOR.MINOR.PATCH`

- **MAJOR**: 不兼容的 API 变更
- **MINOR**: 新增功能，向后兼容
- **PATCH**: 问题修复，向后兼容

例如：
- `v0.1.0` → 初始版本
- `v0.1.1` → 修复 bug
- `v0.2.0` → 新增功能
- `v1.0.0` → 正式版本

## 发布检查清单

发布前检查：

- [ ] 所有测试通过
- [ ] 更新了 CHANGELOG
- [ ] 更新了版本号（package.json）
- [ ] 更新了 README（如有必要）
- [ ] 提交了所有更改
- [ ] 创建了正确的 git tag

发布后检查：

- [ ] GitHub Actions 构建成功
- [ ] GitHub Release 创建成功
- [ ] 5 个平台二进制文件都已上传
- [ ] NPM 包发布成功
- [ ] npx 命令可正常工作
- [ ] 更新了文档

## 回滚发布

### 删除 GitHub Release

```bash
# 删除 tag
git tag -d v0.1.1
git push origin :refs/tags/v0.1.1

# 在 GitHub 网页上删除 Release
```

### 撤销 NPM 发布

```bash
# 72 小时内可以撤销
npm unpublish @xuzan/yst-go-mcp@0.1.1

# 或废弃版本
npm deprecate @xuzan/yst-go-mcp@0.1.1 "版本有问题，请使用 @latest"
```

## 常见问题

### Q: GitHub Actions 构建失败？

A: 检查：
1. Go 版本是否正确（1.25.0）
2. 依赖是否能正常下载
3. 查看 Actions 日志详细错误

### Q: NPM 发布失败？

A: 检查：
1. NPM_TOKEN 是否正确配置
2. 版本号是否已存在
3. 包名是否可用

### Q: 二进制文件无法运行？

A: 检查：
1. 平台/架构是否匹配
2. 是否有执行权限（chmod +x）
3. 是否缺少依赖（如 Chrome）

## 下一版本计划

- [ ] 添加更多工具
- [ ] 优化错误处理
- [ ] 支持配置文件
- [ ] 添加更多测试
