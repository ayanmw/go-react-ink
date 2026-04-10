# Go-Ink 复刻项目设计文档

> 使用 Golang 复刻 React Ink 完整框架实现

## 项目目标

- 在 Go 语言中复刻 React Ink 的声明式 + Flexbox 体验
- 零学习成本：React Ink 开发者可直接使用
- 跨平台兼容：Windows / macOS / Linux

---

## 核心架构概览

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go-Ink Architecture                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   用户代码层                                                     │
│   ┌─────────────────┐                                           │
│   │ .gox 文件 (JSX) │  ← JSX 语法，React 风格                   │
│   └─────────────────┘                                           │
│           │                                                     │
│           ▼  gox-compiler (预编译)                              │
│   ┌─────────────────┐                                           │
│   │   .go 文件      │  ← 生成的 Go 代码                         │
│   └─────────────────┘                                           │
│           │                                                     │
│           ▼                                                     │
│   ┌─────────────────┐                                           │
│   │ React Reconciler│  ← Fiber Tree, Diff算法, Scheduler        │
│   │   (协调器)       │                                           │
│   └─────────────────┘                                           │
│           │                                                     │
│           ▼                                                     │
│   ┌─────────────────┐                                           │
│   │  Ink Renderer   │  ← Host Config 实现                       │
│   │  (自定义宿主)    │                                           │
│   └─────────────────┘                                           │
│           │                                                     │
│     ┌─────┴─────┐                                               │
│     ▼           ▼                                               │
│   ┌─────┐   ┌───────────┐                                       │
│   │Yoga │   │ Terminal  │                                       │
│   │Layout│   │ Output   │                                       │
│   │(Flex)│   │ (ANSI)   │                                       │
│   └─────┘   └───────────┘                                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 核心模块

| 模块 | 功能 | 状态 |
|-----|------|------|
| gox-compiler | JSX → Go 预编译器 | ✅ 完成 |
| reconciler | React Fiber 协调器 | ✅ 完成 |
| layout | Flexbox 布局引擎 | ✅ 完成 |
| renderer | 终端渲染器 | ✅ 完成 |
| components | 内置组件 (Box, Text, etc.) | ✅ 完成 |
| hooks | React Hooks 复刻 | ✅ 完成 |

---

## 核心技术特性

### 1. 备用屏幕缓冲 (Alternate Screen Buffer)

使用 ANSI `\x1b[?1049h` 启用备用屏幕，`\x1b[?1049l` 恢复原屏幕。

**优势**:
- TUI 运行时不影响原终端内容
- 退出后自动恢复原终端状态
- 消除滚动闪烁

**配置**:
```go
opts := &ink.RenderOptions{
    AlternateScreen: true, // 启用备用屏幕
}
```

### 2. TTY 检测

自动检测 stdout 是否为终端，非终端时禁用交互模式。

**实现**:
```go
// 检测文件描述符是否为终端
func IsTerminal(file *os.File) bool {
    fi, _ := file.Stat()
    return (fi.Mode() & os.ModeCharDevice) != 0
}
```

**行为**:
- 终端: 启用交互模式、增量渲染
- 管道/重定向: 禁用交互模式，直接输出

### 3. 信号处理

内置 SIGINT/SIGTERM 处理，Ctrl+C 可靠退出。

**配置**:
```go
opts := &ink.RenderOptions{
    ExitOnCtrlC: true, // 默认启用
}
```

### 4. 增量渲染

只更新变化的行，最小化终端输出。

**实现**:
- 双缓冲对比 (Previous vs Current)
- 逐行差异计算
- ANSI 序列优化

**配置**:
```go
opts := &ink.RenderOptions{
    IncrementalRendering: true, // 默认启用
}
```

---

## 设计决策

### 1. JSX 支持

**决策**: 使用预处理器 方案

- `.gox` 文件包含 JSX 语法
- `gox-compiler` 预编译为标准 `.go` 文件
- 编译时类型检查
- 零运行时开销

**理由**:
- 完全兼容 React Ink 的 JSX 语法
- React 开发者零学习成本
- Go 编译时类型安全

### 2. Flexbox 布局

**决策**: Go 核心子集实现 (80% 常用场景)

支持的属性:
- flexDirection (row/column)
- justifyContent (flex-start/center/space-between/etc.)
- alignItems (stretch/center/flex-start/etc.)
- flex (grow/shrink/basis)
- gap (rowGap/columnGap)
- padding/margin
- width/height (固定值 + 百分比)

**理由**:
- 终端 UI 布局需求相对简单
- 可简化 30-40% 复杂度
- 渐进增强策略

### 3. 渲染策略

**决策**: 复刻 Ink 的增量渲染

流程:
1. React 状态变更
2. Reconciler 触发重新渲染
3. 构建/更新组件树
4. Yoga 计算布局
5. 渲染到 Output Buffer
6. Diff: Current vs Previous
7. 生成最小 ANSI 指令集
8. 写入 stdout

---

## 项目结构

```
go-ink/
├── cmd/
│   └── gox-compiler/          # JSX 编译器 CLI
│
├── pkg/
│   ├── compiler/              # 编译器核心
│   │   ├── scanner/           # 源码扫描
│   │   ├── lexer/             # 词法分析
│   │   ├── parser/            # 语法分析
│   │   ├── codegen/           # 代码生成
│   │   └── sourcemap/         # Source Map
│   │
│   ├── reconciler/            # React 协调器
│   │   ├── fiber.go           # Fiber 树
│   │   ├── scheduler.go       # 任务调度
│   │   ├── diff.go            # Diff 算法
│   │   └── work_loop.go       # 渲染循环
│   │
│   ├── layout/                # Flexbox 布局
│   │   ├── yoga.go            # 布局计算
│   │   ├── node.go            # 布局节点
│   │   ├── measure.go         # 尺寸测量
│   │   └── cache.go           # 布局缓存
│   │
│   ├── renderer/              # 终端渲染
│   │   ├── buffer.go          # Output Buffer
│   │   ├── diff.go            # 帧差异计算
│   │   ├── ansi.go            # ANSI 指令
│   │   └── output.go          # stdout 写入
│   │
│   ├── components/            # 内置组件
│   │   ├── box.go
│   │   ├── text.go
│   │   ├── static.go
│   │   └── ...
│   │
│   ├── hooks/                 # React Hooks
│   │   ├── state.go           # useState
│   │   ├── effect.go          # useEffect
│   │   ├── input.go           # useInput
│   │   └── ...
│   │
│   └── core/
│       ├── element.go         # React.Element
│       ├── component.go       # 组件接口
│       └── props.go           # 属性定义
│
├── tooling/
│   ├── vscode/                # VSCode 插件
│   ├── intellij/              # GoLand 插件
│   └── lsp/                   # Language Server
│
└── doc/
    ├── 01-overview.md         # 概览 (本文档)
    ├── 02-jsx-parser.md       # JSX 解析器设计
    ├── 03-expr-transform.md   # 表达式转换
    ├── 04-codegen.md          # 代码生成器
    ├── 05-toolchain.md        # 工具链集成
    └── ...
```

---

## 参考资料

- [React Ink](https://github.com/vadimdemedes/ink) - 原始框架
- [React Reconciler](https://github.com/facebook/react/tree/main/packages/react-reconciler) - React 协调器
- [Yoga Layout](https://github.com/facebook/yoga) - Flexbox 布局引擎
- [Templ](https://github.com/a-h/templ) - Go HTML 模板引擎 (参考架构)
