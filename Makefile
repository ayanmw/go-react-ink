# Go-React-Ink Makefile
# 统一管理所有构建、测试、发布脚本

SHELL := /bin/bash
GO ?= go
GOPROXY ?= https://goproxy.io,direct

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 目录
PROJECT_DIR := $(shell pwd)
BIN_DIR := bin
TOOLS_DIR := tools

# 颜色输出
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.PHONY: help
help: ## 显示帮助信息
	@echo "Go-React-Ink 构建命令"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "示例: make build    # 构建项目"
	@echo "      make test     # 运行测试"
	@echo "      make plugin   # 构建 GoLand 插件"

# ============================================
# 构建命令
# ============================================

.PHONY: build
build: ## 构建所有包
	@echo "$(GREEN)构建项目...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) build ./...
	@echo "$(GREEN)✓ 构建完成$(RESET)"

.PHONY: build-gox
build-gox: ## 构建 gox 编译器
	@echo "$(GREEN)构建 gox 编译器...$(RESET)"
	mkdir -p $(BIN_DIR)
	GOPROXY=$(GOPROXY) $(GO) build $(LDFLAGS) -o $(BIN_DIR)/gox.exe ./cmd/gox
	@echo "$(GREEN)✓ gox 编译器已构建: $(BIN_DIR)/gox.exe$(RESET)"

.PHONY: build-lsp
build-lsp: ## 构建 LSP Server
	@echo "$(GREEN)构建 LSP Server...$(RESET)"
	mkdir -p $(BIN_DIR)
	cd $(TOOLS_DIR)/lsp-server && GOPROXY=$(GOPROXY) $(GO) build $(LDFLAGS) -o ../../$(BIN_DIR)/gox-lsp.exe .
	@echo "$(GREEN)✓ LSP Server 已构建: $(BIN_DIR)/gox-lsp.exe$(RESET)"

.PHONY: build-all
build-all: build build-gox build-lsp ## 构建所有组件

# ============================================
# 测试命令
# ============================================

.PHONY: test
test: ## 运行所有测试
	@echo "$(GREEN)运行测试...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) test ./... -v
	@echo "$(GREEN)✓ 测试完成$(RESET)"

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "$(GREEN)生成覆盖率报告...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ 覆盖率报告: coverage.html$(RESET)"

.PHONY: test-short
test-short: ## 快速测试（跳过长时间测试）
	@echo "$(GREEN)快速测试...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) test ./... -short

.PHONY: test-e2e
test-e2e: ## 运行端到端测试
	@echo "$(GREEN)运行 E2E 测试...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) test ./test/e2e/... -v

.PHONY: test-examples
test-examples: ## 运行示例测试
	@echo "$(GREEN)运行示例测试...$(RESET)"
	@for dir in test/examples/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "测试 $$dir"; \
			cd $$dir && $(GO) test ./... && cd -; \
		fi; \
	done

# ============================================
# 代码检查
# ============================================

.PHONY: fmt
fmt: ## 格式化代码
	@echo "$(GREEN)格式化代码...$(RESET)"
	$(GO) fmt ./...
	gofmt -s -w .
	@echo "$(GREEN)✓ 格式化完成$(RESET)"

.PHONY: vet
vet: ## 运行 go vet
	@echo "$(GREEN)运行 go vet...$(RESET)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ go vet 通过$(RESET)"

.PHONY: lint
lint: vet ## 运行代码检查
	@echo "$(GREEN)运行 lint...$(RESET)"
	@which golangci-lint > /dev/null || { echo "$(YELLOW)golangci-lint 未安装，跳过$(RESET)"; exit 0; }
	golangci-lint run --timeout 5m
	@echo "$(GREEN)✓ Lint 检查完成$(RESET)"

.PHONY: check
check: fmt lint test-short ## 完整检查（格式化 + lint + 测试）

# ============================================
# 插件构建
# ============================================

.PHONY: plugin
plugin: plugin-goland plugin-vscode ## 构建所有插件

.PHONY: plugin-goland
plugin-goland: ## 构建 GoLand 插件
	@echo "$(GREEN)构建 GoLand 插件...$(RESET)"
	@if [ ! -f "$(TOOLS_DIR)/goland-gox/gradlew" ]; then \
		echo "$(YELLOW)Gradle wrapper 不存在，检查 Gradle...$(RESET)"; \
		if ! command -v gradle &> /dev/null; then \
			echo "$(YELLOW)Gradle 未安装，尝试安装...$(RESET)"; \
			if command -v sdk &> /dev/null; then \
				sdk install gradle; \
			elif command -v brew &> /dev/null; then \
				brew install gradle; \
			elif command -v apt-get &> /dev/null; then \
				sudo apt-get install -y gradle; \
			else \
				echo "$(RED)✗ 无法自动安装 Gradle，请手动安装: https://gradle.org/install/$(RESET)"; \
				exit 1; \
			fi; \
		fi; \
		cd $(TOOLS_DIR)/goland-gox && gradle wrapper; \
	fi
	@cd $(TOOLS_DIR)/goland-gox && ./gradlew buildPlugin
	@echo "$(GREEN)✓ GoLand 插件: $(TOOLS_DIR)/goland-gox/build/distributions/$(RESET)"

.PHONY: plugin-vscode
plugin-vscode: ## 构建 VSCode 插件
	@echo "$(GREEN)构建 VSCode 插件...$(RESET)"
	@if ! command -v bun &> /dev/null; then \
		echo "$(YELLOW)Bun 未安装，尝试安装...$(RESET)"; \
		if command -v curl &> /dev/null; then \
			curl -fsSL https://bun.sh/install | bash; \
		elif command -v powershell &> /dev/null; then \
			powershell -c "irm bun.sh/install.ps1 | iex"; \
		else \
			echo "$(RED)✗ 无法自动安装 Bun，请手动安装: https://bun.sh$(RESET)"; \
			exit 1; \
		fi; \
	fi
	@cd $(TOOLS_DIR)/vscode-gox && bun install && bun run compile
	@echo "$(GREEN)✓ VSCode 插件构建完成$(RESET)"

.PHONY: plugin-package
plugin-package: plugin-goland plugin-vscode ## 打包插件为可发布文件
	@echo "$(GREEN)打包插件...$(RESET)"
	@mkdir -p $(BIN_DIR)/plugins
	@# GoLand 插件
	@if [ -d "$(TOOLS_DIR)/goland-gox/build/distributions" ]; then \
		cp $(TOOLS_DIR)/goland-gox/build/distributions/*.zip $(BIN_DIR)/plugins/ 2>/dev/null || true; \
	fi
	@# VSCode 插件
	@cd $(TOOLS_DIR)/vscode-gox && bun x vsce package --allow-missing-repository --out ../../$(BIN_DIR)/plugins/
	@echo "$(GREEN)✓ 插件已打包到: $(BIN_DIR)/plugins/$(RESET)"
	@ls -la $(BIN_DIR)/plugins/ 2>/dev/null || echo "$(YELLOW)没有生成插件包$(RESET)"

.PHONY: plugin-run-ide
plugin-run-ide: ## 运行 GoLand 插件开发模式
	@echo "$(GREEN)启动 GoLand 插件开发模式...$(RESET)"
	cd $(TOOLS_DIR)/goland-gox && ./gradlew runIde

# ============================================
# 依赖管理
# ============================================

.PHONY: deps
deps: ## 下载依赖
	@echo "$(GREEN)下载依赖...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) mod download
	@echo "$(GREEN)✓ 依赖下载完成$(RESET)"

.PHONY: deps-update
deps-update: ## 更新依赖
	@echo "$(GREEN)更新依赖...$(RESET)"
	GOPROXY=$(GOPROXY) $(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ 依赖更新完成$(RESET)"

.PHONY: tidy
tidy: ## 整理依赖
	$(GO) mod tidy

# ============================================
# 示例运行
# ============================================

.PHONY: run-basic
run-basic: build-gox ## 运行基础示例
	@echo "$(GREEN)运行基础示例...$(RESET)"
	cd examples/basic && $(GO) run main.go

.PHONY: run-counter
run-counter: build-gox ## 运行计数器示例
	@echo "$(GREEN)运行计数器示例...$(RESET)"
	cd examples/counter && $(GO) run main.go

.PHONY: run-layout
run-layout: build-gox ## 运行布局示例
	@echo "$(GREEN)运行布局示例...$(RESET)"
	cd examples/layout && $(GO) run main.go

# ============================================
# 清理
# ============================================

.PHONY: clean
clean: ## 清理构建产物
	@echo "$(YELLOW)清理构建产物...$(RESET)"
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	find . -name "*.exe" -delete
	find . -name "*.test" -delete
	@echo "$(GREEN)✓ 清理完成$(RESET)"

.PHONY: clean-all
clean-all: clean ## 深度清理（包括插件构建产物）
	@echo "$(YELLOW)深度清理...$(RESET)"
	rm -rf $(TOOLS_DIR)/goland-gox/build
	rm -rf $(TOOLS_DIR)/goland-gox/.gradle
	rm -rf $(TOOLS_DIR)/vscode-gox/node_modules
	rm -rf $(TOOLS_DIR)/vscode-gox/out
	@echo "$(GREEN)✓ 深度清理完成$(RESET)"

# ============================================
# 发布
# ============================================

.PHONY: release
release: build-all test ## 准备发布
	@echo "$(GREEN)准备发布 v$(VERSION)...$(RESET)"
	@echo "构建产物:"
	@ls -la $(BIN_DIR)/

.PHONY: install
install: build-gox ## 安装 gox 到 $GOPATH/bin
	@echo "$(GREEN)安装 gox...$(RESET)"
	cp $(BIN_DIR)/gox.exe $(shell go env GOPATH)/bin/
	@echo "$(GREEN)✓ gox 已安装到 $(shell go env GOPATH)/bin/$(RESET)"

# ============================================
# 文档
# ============================================

.PHONY: docs
docs: ## 生成文档
	@echo "$(GREEN)生成 API 文档...$(RESET)"
	$(GO) doc -all ./pkg/core > docs/api/core.md 2>/dev/null || true
	$(GO) doc -all ./pkg/hooks > docs/api/hooks.md 2>/dev/null || true
	$(GO) doc -all ./pkg/components > docs/api/components.md 2>/dev/null || true
	@echo "$(GREEN)✓ 文档生成完成$(RESET)"

# ============================================
# 开发
# ============================================

.PHONY: watch
watch: ## 监听文件变化并自动构建
	@echo "$(GREEN)监听模式...$(RESET)"
	@which reflex > /dev/null || { echo "请安装 reflex: go install github.com/cespare/reflex@latest"; exit 1; }
	reflex -r '\.go$$' -s -- make build

.PHONY: dev
dev: deps build-gox build-lsp ## 开发环境准备
	@echo "$(GREEN)开发环境准备完成$(RESET)"
	@echo "可用命令:"
	@echo "  make test       - 运行测试"
	@echo "  make watch      - 监听模式"
	@echo "  make plugin     - 构建插件"

# ============================================
# 信息
# ============================================

.PHONY: version
version: ## 显示版本信息
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Go 版本: $$(go version)"

.PHONY: info
info: ## 显示项目信息
	@echo "项目: go-react-ink"
	@echo "版本: $(VERSION)"
	@echo "Go: $$(go version)"
	@echo "GOPATH: $$(go env GOPATH)"
	@echo "GOPROXY: $(GOPROXY)"

.DEFAULT_GOAL := help