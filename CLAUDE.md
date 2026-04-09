# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Go-React-Ink 是使用 Golang 完整复刻 [React Ink](https://github.com/vadimdemedes/ink) 框架的终端 UI 开发库，核心特性：

- **JSX 语法编译**: .gox 文件经过编译器 (Scanner → Lexer → Parser → CodeGen) 转换为标准 Go 代码
- **Fiber 协调器**: 采用与 React 一致的 Fiber 架构，支持可中断渲染
- **Flexbox 布局**: 纯 Go 实现的终端场景简化版 Flexbox 引擎
- **IDE 支持**: 包含 VSCode 和 GoLand 插件，提供语法高亮和 LSP 支持

## 常用命令

### 构建
- `make build` - 构建所有 Go 包
- `make build-gox` - 构建 gox 编译器到 `bin/gox.exe`
- `make build-lsp` - 构建 LSP 服务器到 `bin/gox-lsp.exe`
- `make build-all` - 构建所有组件（Go 包 + gox + LSP）

### 测试
- `make test` - 运行所有测试
- `make test-coverage` - 生成覆盖率报告 (coverage.html)
- `make test-short` - 快速测试（跳过长时间测试）
- `make test e2e` - 运行端到端测试
- `make test-examples` - 运行 examples/ 目录下的独立测试用例

### 插件构建 (优先使用 bun)
- `make plugin` - 构建所有插件
- `make plugin Package` - 打包插件为可发布文件
- `make plugin-goland` - 构建 GoLand 插件 (需要 JAVA_HOME)
- `make plugin-vscode` - 构建 VSCode 插件 (使用 bun)

**注意**: 项目优先使用 bun，除非 npm 可用但 bun 不可用。GoLand 插件构建依赖 JAVA_HOME 环境变量，默认路径为 `C:/Users/anmingwei/AppData/Local/Programs/PyCharm/jbr`。

### 示例运行
- `make run-basic` - 运行基础示例
- `make run-counter` - 运行计数器示例
- `make run-layout` - 运行布局示例
- 在 examples/ 目录下运行: `make run-all` / `make build-all`

### 代码质量
- `make fmt` - 格式化代码
- `make vet` - 运行 go vet
- `make lint` - 运行 golangci-lint (可选)
- `make check` - 完整检查（格式化 + lint + 测试）

## 架构

### 编译时 (.gox → .go)
```
Scanner → Lexer → Parser (AST) → CodeGen → .go
```

### 运行时
- **Core**: 核心类型 (Element, Node 等)
- **Reconciler**: Fiber 协调器，调度器和工作循环
- **Layout**: Flexbox 布局引擎
- **Renderer**: 终端渲染器 (Buffer + Diff + ANSI)
- **Components**: 内置组件
- **Hooks**: 状态管理 Hooks

### 模块依赖
- reconciler 依赖 core, hooks
- layout 依赖 core
- renderer 依赖 core, layout
- components 依赖 core
- hooks 依赖 core, reconciler

## 测试用例

examples/ 目录下的示例程序包含独立的 go.mod 文件，使用 replace 指定本地 go-react-ink 目录：

```go
replace github.com/ayanmw/go-react-ink => ..
```

这允许使用最新代码进行完整独立测试。在 examples/ 下运行 `make test-examples` 会测试所有带 go.mod 的示例。

## 特定行为

- **目录构建**: build-lsp 需要进入 `tools/lsp-server` 目录编译 (工作区目录对构建有影响)
- **VSCode 插件保存**: VSCode-go 插件安装后，保存 .gox 文件应输出为 .go (不是 .goo)
- **GoLand 构建路径**: GoLand 插件构建需要指定 JAVA_HOME: `C:/Users/anmingwei/AppData/Local/Programs/PyCharm/jbr`

## 文档

- `doc/` - 设计文档 (01-overview.md 到 10-hooks.md)
- `openspec/specs/go-ink-spec.md` - Go-Ink 规范
- `reference/react-ink/` - React Ink 原版代码（参考）

## 重要约束

1. **bun 优先**: 所有 Node.js 相关工具链（VSCode 插件等）优先使用 bun，除非 bun 不可用
2. **目录感知**: 部分构建命令需要在特定目录下执行 (如 LSP 构建需要在 tools/lsp-server)
3. **GoLand 路径**: GoLand 构建 JAVA_HOME 硬编码为特定路径，修改需同步更新 Makefile
