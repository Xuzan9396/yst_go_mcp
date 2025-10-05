# YST Go MCP Makefile

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
LDFLAGS := -s -w
BUILD_DIR := build
BINARY_NAME := yst-go-mcp

# 平台列表
PLATFORMS := darwin-amd64 darwin-arm64 linux-amd64 windows-amd64 windows-arm64

.PHONY: all build clean install test build-all release help

# 默认目标
all: build

# 编译当前平台
build:
	@echo "🔨 编译 $(BINARY_NAME)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/yst-go-mcp
	@echo "✅ 编译完成: $(BINARY_NAME)"

# 编译所有平台
build-all: clean
	@echo "🔨 编译所有平台..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'-' -f1); \
		GOARCH=$$(echo $$platform | cut -d'-' -f2); \
		output=$(BUILD_DIR)/$(BINARY_NAME)-$$platform; \
		if [ "$$GOOS" = "windows" ]; then output=$$output.exe; fi; \
		echo "  → 编译 $$platform..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags="$(LDFLAGS)" -o $$output ./cmd/yst-go-mcp || exit 1; \
		ls -lh $$output; \
	done
	@echo "✅ 所有平台编译完成！"
	@ls -lh $(BUILD_DIR)/

# 清理构建产物
clean:
	@echo "🧹 清理构建产物..."
	@rm -rf $(BUILD_DIR) $(BINARY_NAME) yst-go-mcp-*
	@echo "✅ 清理完成"

# 安装到本地
install: build
	@echo "📦 安装到 /usr/local/bin..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "✅ 安装完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 运行 MCP 客户端测试
test-mcp:
	@echo "🧪 运行 MCP 测试..."
	@cd test_mcp && go run TestClient.go

# 运行完整测试
test-full:
	@echo "🧪 运行完整测试..."
	@cd test_mcp && ./run_all_tests.sh

# 创建 release
release: build-all
	@echo "📦 准备 release..."
	@VERSION_TAG=$$(echo $(VERSION) | sed 's/^v//'); \
	echo "版本: $$VERSION_TAG"

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 格式化完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	go vet ./...
	@echo "✅ 检查完成"

# 下载依赖
deps:
	@echo "📥 下载依赖..."
	go mod download
	go mod tidy
	@echo "✅ 依赖更新完成"

# 显示帮助
help:
	@echo "YST Go MCP Makefile"
	@echo ""
	@echo "使用方法:"
	@echo "  make build       - 编译当前平台"
	@echo "  make build-all   - 编译所有平台"
	@echo "  make clean       - 清理构建产物"
	@echo "  make install     - 安装到本地"
	@echo "  make test        - 运行测试"
	@echo "  make test-mcp    - 运行 MCP 测试"
	@echo "  make test-full   - 运行完整测试"
	@echo "  make release     - 创建 release"
	@echo "  make fmt         - 格式化代码"
	@echo "  make lint        - 代码检查"
	@echo "  make deps        - 下载依赖"
	@echo "  make help        - 显示帮助"
	@echo ""
	@echo "版本: $(VERSION)"
