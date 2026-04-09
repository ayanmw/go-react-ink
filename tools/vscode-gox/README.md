# GoX - JSX for Go

VSCode 扩展，为 `.gox` 文件提供语法高亮和语言支持。

## 功能

### 语法高亮
- Go 代码语法高亮
- JSX 标签和属性高亮
- 嵌入表达式 `{...}` 高亮

### 命令

| 命令 | 快捷键 | 说明 |
|------|--------|------|
| `GoX: Compile Current File` | `Ctrl+Shift+G` | 编译当前 .gox 文件 |
| `GoX: Compile All .gox Files` | - | 编译工作区所有 .gox 文件 |
| `GoX: Watch and Compile` | - | 监听文件变化自动编译 |
| `GoX: Run Lint` | - | 运行 gofmt, go vet, golint |

### 配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `gox.goxPath` | `gox` | gox 编译器路径 |
| `gox.compileOnSave` | `true` | 保存时自动编译 |
| `gox.componentPackage` | `ink` | 默认组件包名 |
| `gox.lsp.enabled` | `true` | 启用 LSP 功能 |
| `gox.lsp.path` | `gox-lsp` | LSP 服务器路径 |

## 安装

### 从 VSIX 安装

1. 下载 `.vsix` 文件
2. 打开 VSCode
3. 按 `Ctrl+Shift+P`，输入 `Install from VSIX`
4. 选择下载的文件

### 从源码构建

```bash
cd tools/vscode-gox
bun install
bun run compile
vsce package  # 需要: bun install -g @vscode/vsce
```

## 使用

1. 创建 `.gox` 文件：

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
        </Box>
    )
}
```

2. 保存文件，自动编译为 `.go`

3. 运行生成的 Go 代码：

```bash
go run app.go
```

## LSP 支持

扩展支持通过 LSP 服务器提供：

- 代码补全
- 悬停文档
- 错误诊断
- 跳转定义

确保 `gox-lsp` 在 PATH 中可用。

## 相关链接

- [Go-React-Ink](https://github.com/ayanmw/go-react-ink) - 主项目
- [GoLand Plugin](https://github.com/ayanmw/go-react-ink/tree/main/tools/goland-gox) - GoLand IDE 插件

## 许可证

MIT License