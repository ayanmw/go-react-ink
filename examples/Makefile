# Go-Ink Examples Makefile
# 运行和构建所有示例程序

.PHONY: all run-all build-all clean help
.PHONY: run-hello run-counter run-todo run-flexbox run-hooks run-styling run-input run-animation
.PHONY: build-hello build-counter build-todo build-flexbox build-hooks build-styling build-input build-animation

# 所有示例
EXAMPLES := hello-world counter-app todo-list flexbox-layout hooks-demo text-styling interactive-input animation-demo

# 默认目标
all: help

# 帮助信息
help:
	@echo "Go-Ink Examples"
	@echo ""
	@echo "Usage:"
	@echo "  make run-all        - 运行所有示例"
	@echo "  make build-all     - 构建所有示例"
	@echo "  make clean         - 清理构建产物"
	@echo ""
	@echo "Run specific example:"
	@echo "  make run-hello     - 运行 hello-world"
	@echo "  make run-counter   - 运行 counter-app"
	@echo "  make run-todo      - 运行 todo-list"
	@echo "  make run-flexbox   - 运行 flexbox-layout"
	@echo "  make run-hooks     - 运行 hooks-demo"
	@echo "  make run-styling   - 运行 text-styling"
	@echo "  make run-input     - 运行 interactive-input"
	@echo "  make run-animation - 运行 animation-demo"
	@echo ""
	@echo "Build specific example:"
	@echo "  make build-hello   - 构建 hello-world"
	@echo "  make build-counter - 构建 counter-app"
	@echo "  ..."

# 运行所有示例
run-all:
	@echo "=== Running all examples ==="
	@for dir in $(EXAMPLES); do \
		echo ""; \
		echo ">>> Running $$dir..."; \
		cd $$dir && go run main.go; \
		cd ..; \
	done

# 构建所有示例
build-all:
	@echo "=== Building all examples ==="
	@for dir in $(EXAMPLES); do \
		echo "Building $$dir..."; \
		cd $$dir && go build -o ../bin/$$dir main.go; \
		cd ..; \
	done

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@for dir in $(EXAMPLES); do \
		cd $$dir && rm -f $$dir; \
		cd ..; \
	done

# 运行单个示例
run-hello:
	@cd hello-world && go run main.go

run-counter:
	@cd counter-app && go run main.go

run-todo:
	@cd todo-list && go run main.go

run-flexbox:
	@cd flexbox-layout && go run main.go

run-hooks:
	@cd hooks-demo && go run main.go

run-styling:
	@cd text-styling && go run main.go

run-input:
	@cd interactive-input && go run main.go

run-animation:
	@cd animation-demo && go run main.go

# 构建单个示例
build-hello:
	@cd hello-world && go build -o hello-world main.go

build-counter:
	@cd counter-app && go build -o counter-app main.go

build-todo:
	@cd todo-list && go build -o todo-list main.go

build-flexbox:
	@cd flexbox-layout && go build -o flexbox-layout main.go

build-hooks:
	@cd hooks-demo && go build -o hooks-demo main.go

build-styling:
	@cd text-styling && go build -o text-styling main.go

build-input:
	@cd interactive-input && go build -o interactive-input main.go

build-animation:
	@cd animation-demo && go build -o animation-demo main.go
