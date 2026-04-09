# Makefile for React Ink testing and development

# 默认目标
.DEFAULT_GOAL := help

# Node 和 npm 命令
NODE := node
NPM := npm
NPX := npx
TSX := $(NPX) tsx

# 项目目录
EXAMPLES_DIR := examples
TEST_DIR := test
SRC_DIR := src
BUILD_DIR := build

# 示例列表
EXAMPLES := \
	borders \
	box-backgrounds \
	border-backgrounds \
	chat \
	concurrent-suspense \
	counter \
	cursor-ime \
	incremental-rendering \
	justify-content \
	render-throttle \
	router \
	select-input \
	static \
	subprocess-output \
	suspense \
	table \
	terminal-resize \
	use-animation \
	use-focus \
	use-focus-with-id \
	use-input \
	use-stderr \
	use-stdout \
	use-transition \
	alternate-screen \
	aria \
	jest

# 彩色输出
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_BLUE := \033[34m
COLOR_YELLOW := \033[33m
COLOR_CYAN := \033[36m

.PHONY: help
help: ## 显示帮助信息
	@echo "$(COLOR_BOLD)React Ink Makefile$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)可用命令:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_CYAN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(COLOR_YELLOW)说明:$(COLOR_RESET) 带 * 的命令会启动交互式程序，需要 Ctrl+C 退出"

.PHONY: build
build: ## 编译 TypeScript 源代码
	@echo "$(COLOR_BLUE)编译项目...$(COLOR_RESET)"
	@$(NPM) run build
	@echo "$(COLOR_GREEN)✓ 编译完成$(COLOR_RESET)"

.PHONY: clean
clean: ## 清理编译输出
	@echo "$(COLOR_BLUE)清理构建输出...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(COLOR_GREEN)✓ 清理完成$(COLOR_RESET)"

.PHONY: typecheck
typecheck: ## 运行 TypeScript 类型检查
	@echo "$(COLOR_BLUE)运行类型检查...$(COLOR_RESET)"
	@$(NPM) run typecheck
	@echo "$(COLOR_GREEN)✓ 类型检查通过$(COLOR_RESET)"

.PHONY: lint
lint: ## 运行 ESLint 检查
	@echo "$(COLOR_BLUE)运行代码检查...$(COLOR_RESET)"
	@$(NPM) run lint
	@echo "$(COLOR_GREEN)✓ Lint 检查完成$(COLOR_RESET)"

.PHONY: format
format: ## 使用 Prettier 格式化代码
	@echo "$(COLOR_BLUE)格式化代码...$(COLOR_RESET)"
	@$(NPM) run format
	@echo "$(COLOR_GREEN)✓ 格式化完成$(COLOR_RESET)"

.PHONY: test
test: ## 运行完整测试套件
	@echo "$(COLOR_BLUE)运行测试套件...$(COLOR_RESET)"
	@$(NPM) test
	@echo "$(COLOR_GREEN)✓ 测试完成$(COLOR_RESET)"

.PHONY: test-ava
test-ava: ## 运行 Ava 测试（不清理）
	@echo "$(COLOR_BLUE)运行 Ava 测试...$(COLOR_RESET)"
	@$(NPX) ava --serial
	@echo "$(COLOR_GREEN)✓ Ava 测试完成$(COLOR_RESET)"

.PHONY: check
check: typecheck lint test-ava ## 运行所有检查（类型检查 + Lint + 测试）
	@echo "$(COLOR_GREEN)✓ 所有检查通过$(COLOR_RESET)"

# ==============================================================================
# 示例运行命令
# ==============================================================================

.PHONY: demo
demo: ## * 运行计数器示例（最简单）
	@echo "$(COLOR_YELLOW)运行 counter 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/counter/counter.tsx

.PHONY: example-borders
example-borders: ## * 运行边框示例
	@echo "$(COLOR_YELLOW)运行 borders 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/borders/borders.tsx

.PHONY: example-box-backgrounds
example-box-backgrounds: ## * 运行背景色示例
	@echo "$(COLOR_YELLOW)运行 box-backgrounds 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/box-backgrounds/box-backgrounds.tsx

.PHONY: example-static
example-static: ## * 运行 Static 组件示例
	@echo "$(COLOR_YELLOW)运行 static 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/static/static.tsx

.PHONY: example-use-input
example-use-input: ## * 运行 useInput 示例
	@echo "$(COLOR_YELLOW)运行 use-input 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/use-input/use-input.tsx

.PHONY: example-use-focus
example-use-focus: ## * 运行 useFocus 示例
	@echo "$(COLOR_YELLOW)运行 use-focus 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYN)按 Tab 切换焦点$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/use-focus/use-focus.tsx

.PHONY: example-use-focus-id
example-use-focus-id: ## * 运行 use-focus-with-id 示例
	@echo "$(COLOR_YELLOW)运行 use-focus-with-id 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/use-focus-with-id/use-focus-with-id.tsx

.PHONY: example-use-animation
example-use-animation: ## * 运行 useAnimation 动画示例
	@echo "$(COLOR_YELLOW)运行 use-animation 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)按空格暂停/恢复$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/use-animation/use-animation.tsx

.PHONY: example-cursor-ime
example-cursor-ime: ## * 运行光标 IME 示例
	@echo "$(COLOR_YELLOW)运行 cursor-ime 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)可以输入中文/韩文等$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/cursor-ime/cursor-ime.tsx

.PHONY: example-aria
example-aria: ## * 运行 ARIA 屏幕阅读器示例
	@echo "$(COLOR_YELLOW)运行 aria 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)按空格切换$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/aria/aria.tsx

.PHONY: example-alternate-screen
example-alternate-screen: ## * 运行备用屏幕示例（独角兽蛇游戏）
	@echo "$(COLOR_YELLOW)运行 alternate-screen 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)方向键移动，吃 🌈 生长$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/alternate-screen/alternate-screen.tsx

.PHONY: example-suspense
example-suspense: ## 运行 Suspense 示例（自动退出）
	@echo "$(COLOR_YELLOW)运行 suspense 示例...$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/suspense/suspense.tsx

.PHONY: example-concurrent-suspense
example-concurrent-suspense: ## 运行 Concurrent Suspense 示例
	@echo "$(COLOR_YELLOW)运行 concurrent-suspense 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/concurrent-suspense/concurrent-suspense.tsx

.PHONY: example-use-transition
example-use-transition: ## 运行 useTransition 示例
	@echo "$(COLOR_YELLOW)运行 use-transition 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/use-transition/use-transition.tsx

.PHONY: example-select-input
example-select-input: ## * 运行下拉选择示例
	@echo "$(COLOR_YELLOW)运行 select-input 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)上下箭头选择$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/select-input/select-input.tsx

.PHONY: example-table
example-table: ## 运行表格示例（自动退出）
	@echo "$(COLOR_YELLOW)运行 table 示例...$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/table/table.tsx

.PHONY: example-chat
example-chat: ## * 运行聊天界面示例
	@echo "$(COLOR_YELLOW)运行 chat 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)输入消息，回车发送$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/chat/chat.tsx

.PHONY: example-jest
example-jest: ## 运行 Jest 风格测试示例
	@echo "$(COLOR_YELLOW)运行 jest 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/jest/jest.tsx

.PHONY: example-router
example-router: ## * 运行路由示例
	@echo "$(COLOR_YELLOW)运行 router 示例...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)回车切换页面$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/router/router.tsx

.PHONY: example-justify-content
example-justify-content: ## 运行 Flexbox 对齐示例
	@echo "$(COLOR_YELLOW)运行 justify-content 示例...$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/justify-content/justify-content.tsx

# ==============================================================================
# 批量运行示例
# ==============================================================================

.PHONY: run-all-examples
run-all-examples: ## 依次运行所有示例（需要按 Ctrl+C 退出每个）
	@echo "$(COLOR_YELLOW)注意: 需要为每个示例按 Ctrl+C 退出$(COLOR_RESET)"
	@$(foreach example,$(EXAMPLES),\
		echo "$(COLOR_BOLD)=== 运行: $(example) ===$(COLOR_RESET)"; \
		$(TSX) $(EXAMPLES_DIR)/$(example)/$(example).tsx || true; \
		sleep 1; \
	)

# 快速运行测试类示例（自动完成的）
.PHONY: test-examples-auto
test-examples-auto: table example-suspense justify-content ## 运行自动退出的示例
	@echo "$(COLOR_GREEN)✓ 自动类示例运行完成$(COLOR_RESET)"

# 快速运行交互式示例（需要 Ctrl+C）
.PHONY: test-examples-interactive
test-examples-interactive: demo example-borders example-static example-use-input example-use-focus example-cursor-ime example-aria ## 运行核心交互式示例
	@echo "$(COLOR_GREEN)✓ 交互式示例运行完成$(COLOR_RESET)"

# 运行所有 TUI 动画相关示例
.PHONY: test-tui-animations
test-tui-animations: example-use-animation example-alternate-screen ## 运行 TUI 动画示例
	@echo "$(COLOR_GREEN)✓ TUI 动画示例运行完成$(COLOR_RESET)"

# 运行所有 React 高级功能示例
.PHONY: test-advanced
test-advanced: example-concurrent-suspense example-use-transition example-select-input ## 运行高级功能示例
	@echo "$(COLOR_GREEN)✓ 高级功能示例运行完成$(COLOR_RESET)"

# ==============================================================================
# 开发相关
# ==============================================================================

.PHONY: watch
watch: ## 监听模式编译
	@echo "$(COLOR_BLUE)启动监听模式...$(COLOR_RESET)"
	@$(NPM) run dev

.PHONY: Dev
Dev: ## 进入开发工具模式（每5秒循环）
	@echo "$(COLOR_BOLD)进入开发工具模式（每5秒执行一次:typecheck lint format）$(COLOR_RESET)"
	@while true; do \
		$(MAKE) typecheck lint format; \
		echo "$(COLOR_GREEN)完成！等待5秒...$(COLOR_RESET)"; \
		sleep 5; \
	done

# ==============================================================================
# 信息相关
# ==============================================================================

.PHONY: info
info: ## 显示项目信息
	@echo "$(COLOR_BOLD)项目信息:$(COLOR_RESET)"
	@echo "  Node 版本:     $$($(NODE) --version)"
	@echo "  npm 版本:      $$($(NPM) --version)"
	@echo "  tsx 版本:      $$($(TSX) --version)"
	@echo ""
	@echo "$(COLOR_BOLD)统计数据:$(COLOR_RESET)"
	@echo "  源文件数:     $$($(MAKE) -s info-src-count)"
	@echo "  示例数:       $$($(MAKE) -s info-example-count)"
	@echo "  测试文件数:   $$($(MAKE) -s info-test-count)"

.PHONY: info-src-count
info-src-count:
	@find $(SRC_DIR) -name "*.ts" -o -name "*.tsx" | wc -l

.PHONY: info-example-count
info-example-count:
	@find $(EXAMPLES_DIR) -name "*.tsx" | wc -l

.PHONY: info-test-count
info-test-count:
	@find $(TEST_DIR) -name "*.ts" -o -name "*.tsx" | wc -l

.PHONY: list-examples
list-examples: ## 列出所有示例
	@echo "$(COLOR_BOLD)可用示例:$(COLOR_RESET)"
	@for example in $(EXAMPLES); do \
		echo "  - $$example"; \
	done

# ==============================================================================
# 专用测试命令
# ==============================================================================

.PHONY: test-components
test-components: ## 测试渲染组件
	@echo "$(COLOR_BLUE)测试组件渲染...$(COLOR_RESET)"
	@$(NPX) ava test/components.tsx

.PHONY: test-flex
test-flex: ## 测试 Flexbox 布局
	@echo "$(COLOR_BLUE)测试 Flexbox 布局...$(COLOR_RESET)"
	@$(NPX) ava test/flex.tsx
	@$(NPX) ava test/flex-*.tsx

.PHONY: test-hooks
test-hooks: ## 测试 Hooks
	@echo "$(COLOR_BLUE)测试 Hooks...$(COLOR_RESET)"
	@$(NPX) ava test/hooks.tsx

.PHONY: test-display
test-display: ## 测试显示功能
	@echo "$(COLOR_BLUE)测试显示功能...$(COLOR_RESET)"
	@$(NPX) ava test/display.tsx

.PHONY: test-focus
test-focus: ## 测试焦点管理
	@echo "$(COLOR_BLUE)测试焦点管理...$(COLOR_RESET)"
	@$(NPX) ava test/focus.tsx

.PHONY: test-exit
test-exit: ## 测试退出处理
	@echo "$(COLOR_BLUE)测试退出处理...$(COLOR_RESET)"
	@$(NPX) ava test/exit.tsx

.PHONY: test-borders
test-borders: ## 测试边框渲染
	@echo "$(COLOR_BLUE)测试边框渲染...$(COLOR_RESET)"
	@$(NPX) ava test/borders.tsx

.PHONY: run-example-by-name
run-example-by-name: ## 使用 tsx 运行指定示例: make run-example-by-name NAME=borders
	@$(MAKE) run-example-by-name-impl NAME=$(NAME)

.PHONY: run-example-by-name-impl
run-example-by-name-impl:
	@if [ -z "$(NAME)" ]; then \
		echo "$(COLOR_RED)错误: 请指定示例名称，例如: make run-example-by-name NAME=borders$(COLOR_RESET)"; \
		exit 1; \
	fi
	@if [ ! -d "$(EXAMPLES_DIR)/$(NAME)" ]; then \
		echo "$(COLOR_RED)错误: 示例 '$(NAME)' 不存在$(COLOR_RESET)"; \
		$(MAKE) list-examples; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)运行示例: $(NAME)$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)按 Ctrl+C 退出$(COLOR_RESET)"
	@$(TSX) $(EXAMPLES_DIR)/$(NAME)/$(NAME).tsx
