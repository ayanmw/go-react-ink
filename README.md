# Go-React-Ink

[![Go Reference](https://pkg.go.dev/badge/github.com/ayanmw/go-react-ink.svg)](https://pkg.go.dev/github.com/ayanmw/go-react-ink)
[![Go Report Card](https://goreportcard.com/badge/github.com/ayanmw/go-react-ink)](https://goreportcard.com/report/github.com/ayanmw/go-react-ink)
[![CI](https://github.com/ayanmw/go-react-ink/actions/workflows/ci.yml/badge.svg)](https://github.com/ayanmw/go-react-ink/actions/workflows/ci.yml)

使用 Golang 完整复刻 [React Ink](https://github.com/vadimdemedes/ink) 框架，实现声明式终端 UI 开发体验。

## 对标版本

本项目基于 React Ink 进行复刻实现，当前对标版本：

| 项目 | 版本 | Commit | 日期 |
|------|------|--------|------|
| React Ink | v6.8.0 | `be1b1bb6ec65056e2ed60ef3c5ae642704b82d31` | 2025-01 |
| Go-React-Ink | v0.1.0 | - | 2026-04 |

> **注意**: 当 React Ink 发布新版本时，我们会持续同步更新。请关注 [CHANGELOG.md](./CHANGELOG.md) 获取最新同步状态。

## 特性

- 🎨 **JSX 语法** - 使用类似 React 的 JSX 语法编写终端 UI
- ⚡ **Fiber 架构** - 采用 React Fiber 协调器架构，支持可中断渲染
- 📐 **Flexbox 布局** - 纯 Go 实现的 Flexbox 布局引擎
- 🎯 **Hooks API** - 完整的 React-like Hooks 支持
- 🖥️ **跨平台** - 支持 Windows、macOS、Linux

## 安装

```bash
go get github.com/ayanmw/go-react-ink
```

安装 CLI 编译器:

```bash
go install github.com/ayanmw/go-react-ink/cmd/gox@latest
```

## 快速开始

### 1. 创建 .gox 文件

```go
// app.gox
package main

import (
    "github.com/ayanmw/go-react-ink/pkg/core"
    "github.com/ayanmw/go-react-ink/pkg/components"
)

func App() core.Element {
    return (
        <Box flexDirection="column">
            <Text color="green">Hello, World!</Text>
            <Text>Count: 42</Text>
        </Box>
    )
}

func main() {
    app := App()
    println(app.Render())
}
```

### 2. 编译

```bash
gox app.gox -o app.go
```

### 3. 运行

```bash
go run app.go
```

## 组件

### Box

容器组件，支持 Flexbox 布局:

```go
<Box
    flexDirection="column"    // row, column, row-reverse, column-reverse
    justifyContent="center"   // flex-start, flex-end, center, space-between
    alignItems="center"       // flex-start, flex-end, center, stretch
    padding={1}
    margin={1}
>
    {/* children */}
</Box>
```

### Text

文本组件:

```go
<Text
    color="green"      // 文本颜色
    bold={true}        // 粗体
    italic={true}      // 斜体
    underline={true}   // 下划线
>
    Hello, World!
</Text>
```

### Spacer

空白填充:

```go
<Box flexDirection="row">
    <Text>Left</Text>
    <Spacer />
    <Text>Right</Text>
</Box>
```

### Newline

换行:

```go
<Text>Hello</Text>
<Newline />
<Text>World</Text>
```

## Hooks

### useState

状态管理:

```go
func Counter() core.Element {
    count, setCount := hooks.UseState(ctx, 0)
    
    return (
        <Box>
            <Text>Count: {count}</Text>
        </Box>
    )
}
```

### useEffect

副作用处理:

```go
hooks.UseEffect(ctx, func() func() {
    // 初始化
    return func() {
        // 清理
    }
}, []any{deps})
```

### useInput

键盘输入处理:

```go
input := input.UseInput(ctx, func(key input.Key) {
    if key.Name == "enter" {
        // 处理回车
    }
})
```

### useFocus

焦点管理:

```go
isFocused, focus := input.UseFocus(ctx)
```

## 布局

Go-Ink 使用纯 Go 实现的 Flexbox 布局引擎:

```go
root := layout.NewNode()
root.Direction = layout.DirectionColumn
root.Width = 80
root.Height = 24

header := layout.NewNode()
header.Height = 1

content := layout.NewNode()
content.FlexGrow = 1

footer := layout.NewNode()
footer.Height = 1

root.AddChild(header)
root.AddChild(content)
root.AddChild(footer)

root.CalculateLayout(80, 24)
```

## 编译器

### 命令行工具

```bash
# 编译单个文件
gox app.gox -o app.go

# 编译目录
gox ./src

# 监听模式
gox -watch ./src

# 显示帮助
gox -help

# 显示版本
gox -version
```

### 编译选项

| 选项 | 说明 |
|------|------|
| `-o` | 输出文件路径 |
| `-pkg` | 组件包名 (默认: ink) |
| `-watch` | 监听文件变化 |
| `-version` | 显示版本 |
| `-help` | 显示帮助 |

## 架构

```
┌─────────────────────────────────────────────────────┐
│                    Go-Ink 架构                        │
├─────────────────────────────────────────────────────┤
│                                                      │
│  编译时 (.gox → .go)                                 │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐           │
│  │ Scanner │──▶│  Lexer  │──▶│ Parser  │           │
│  └─────────┘   └─────────┘   └─────────┘           │
│                     │                               │
│                     ▼                               │
│  ┌─────────┐   ┌─────────┐                         │
│  │ CodeGen │◀──│Compiler │                         │
│  └─────────┘   └─────────┘                         │
│                                                      │
│  运行时                                              │
│  ┌─────────────────────────────────────────────┐   │
│  │              Fiber Reconciler                 │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐   │   │
│  │  │ Scheduler│  │  WorkLoop│  │  Commit  │   │   │
│  │  └──────────┘  └──────────┘  └──────────┘   │   │
│  └─────────────────────────────────────────────┘   │
│                                                      │
│  ┌─────────────────────────────────────────────┐   │
│  │              Layout Engine                    │   │
│  │  Flexbox • Measure • Position • Align       │   │
│  └─────────────────────────────────────────────┘   │
│                                                      │
│  ┌─────────────────────────────────────────────┐   │
│  │              Renderer                         │   │
│  │  Buffer • Diff • ANSI • Terminal            │   │
│  └─────────────────────────────────────────────┘   │
│                                                      │
└─────────────────────────────────────────────────────┘
```

## API 兼容性

Go-Ink 致力于与 React Ink API 保持一致:

| 功能 | React Ink | Go-Ink | 状态 |
|------|-----------|--------|------|
| Box | ✅ | ✅ | 完成 |
| Text | ✅ | ✅ | 完成 |
| Spacer | ✅ | ✅ | 完成 |
| Newline | ✅ | ✅ | 完成 |
| Static | ✅ | ✅ | 完成 |
| Transform | ✅ | ✅ | 完成 |
| useState | ✅ | ✅ | 完成 |
| useEffect | ✅ | ✅ | 完成 |
| useRef | ✅ | ✅ | 完成 |
| useMemo | ✅ | ✅ | 完成 |
| useInput | ✅ | ✅ | 完成 |
| useApp | ✅ | ✅ | 完成 |
| useFocus | ✅ | ✅ | 完成 |
| useCursor | ✅ | ✅ | 完成 |
| useAnimation | ✅ | ✅ | 完成 |
| Flexbox | ✅ | ✅ | 完成 |

## 示例

查看 [examples](./examples) 目录获取更多示例。

### 计数器

```go
func Counter() core.Element {
    count, setCount := hooks.UseState(ctx, 0)
    
    input.UseInput(ctx, func(key input.Key) {
        if key.Name == "up" {
            setCount(count.(int) + 1)
        } else if key.Name == "down" {
            setCount(count.(int) - 1)
        }
    })
    
    return (
        <Box flexDirection="column">
            <Text>Count: {count}</Text>
            <Text dim>Use Up/Down arrows to change</Text>
        </Box>
    )
}
```

## 开发

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./pkg/layout/... -v

# 运行覆盖率
go test -cover ./...
```

### 构建

```bash
# 构建所有包
go build ./...

# 构建 CLI
go build -o gox ./cmd/gox
```

## 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详情。

## 许可证

MIT License - 详见 [LICENSE](./LICENSE) 文件

## 致谢

- [React Ink](https://github.com/vadimdemedes/ink) - 原始框架灵感
- [React](https://github.com/facebook/react) - Fiber 架构参考
- [Yoga](https://github.com/facebook/yoga) - Flexbox 布局参考