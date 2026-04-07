# 工具链集成设计

> gox-compiler 的工具链支持

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  工具链整体架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   开发者工作流:                                             │
│                                                             │
│   ┌─────────────┐                                          │
│   │ 编辑 .gox  │                                          │
│   │   文件     │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────────────────────────────────────────────┐  │
│   │              工具链处理                               │  │
│   ├─────────────────────────────────────────────────────┤  │
│   │                                                      │  │
│   │   ┌───────────┐    ┌───────────┐    ┌───────────┐  │  │
│   │   │  gox      │    │  gox      │    │   IDE     │  │  │
│   │   │  generate │ or │  watch    │ or │  插件     │  │  │
│   │   └───────────┘    └───────────┘    └───────────┘  │  │
│   │         │                 │                 │        │  │
│   │         └─────────────────┴─────────────────┘        │  │
│   │                           │                          │  │
│   │                           ▼                          │  │
│   │                   ┌───────────┐                      │  │
│   │                   │ Compiler  │                      │  │
│   │                   │  Pipeline │                      │  │
│   │                   └───────────┘                      │  │
│   │                           │                          │  │
│   │                           ▼                          │  │
│   │                   ┌───────────┐                      │  │
│   │                   │  生成     │                      │  │
│   │                   │  .go 文件 │                      │  │
│   │                   └───────────┘                      │  │
│   │                                                      │  │
│   └─────────────────────────────────────────────────────┘  │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  go build   │                                          │
│   │  go run     │                                          │
│   └─────────────┘                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 方式一：go generate 集成

### 使用方式

```
项目结构:
myapp/
├── main.go
├── ui.gox              # JSX 源文件
└── ui_gen.go           # 生成的 Go 文件 (可 gitignore)

main.go:
package main

//go:generate gox-compiler ui.gox

import "github.com/xxx/go-ink"

func main() {
    ink.Render(UI())
}

运行生成:
$ go generate
# 或
$ go generate ./...
```

### CLI 设计

```
命令格式:
gox-compiler [flags] <input> [output]

参数:
├── input:   输入文件或目录
├── output:  输出文件或目录 (可选)
└── flags:   选项标志

选项:
├── -o, --output:     指定输出路径
├── -p, --package:    指定包名
├── -w, --watch:      监听模式
├── -f, --force:      强制覆盖
├── -d, --debug:      调试模式 (生成注释)
├── -v, --verbose:    详细输出
├── --version:        显示版本
└── --help:           显示帮助

示例:
# 单文件
$ gox-compiler ui.gox
# 输出: ui.gox.go

# 指定输出
$ gox-compiler ui.gox -o ui_gen.go

# 整个目录
$ gox-compiler ./ui/
# 编译所有 .gox 文件

# 监听模式
$ gox-compiler -w ./ui/
# 文件变化时自动重新编译
```

### CLI 实现

```go
// cmd/gox-compiler/main.go

package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

var (
    flagOutput   = flag.String("o", "", "output path")
    flagPackage  = flag.String("p", "", "package name")
    flagWatch    = flag.Bool("w", false, "watch mode")
    flagForce    = flag.Bool("f", false, "force overwrite")
    flagDebug    = flag.Bool("d", false, "debug mode")
    flagVerbose  = flag.Bool("v", false, "verbose output")
    flagVersion  = flag.Bool("version", false, "show version")
)

func main() {
    flag.Usage = printUsage
    flag.Parse()

    if *flagVersion {
        fmt.Println("gox-compiler v1.0.0")
        return
    }

    args := flag.Args()
    if len(args) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    input := args[0]
    output := *flagOutput
    if len(args) > 1 && output == "" {
        output = args[1]
    }

    config := &compiler.Config{
        OutputPath: output,
        PackageName: *flagPackage,
        Force:      *flagForce,
        Debug:      *flagDebug,
        Verbose:    *flagVerbose,
    }

    if *flagWatch {
        if err := watch(input, config); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    } else {
        if err := compile(input, config); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    }
}
```

### 编译逻辑

```go
func compileFile(inputPath string, config *Config) error {
    // 1. 检查扩展名
    if filepath.Ext(inputPath) != ".gox" {
        return fmt.Errorf("input must be .gox file")
    }

    // 2. 确定输出路径
    outputPath := config.OutputPath
    if outputPath == "" {
        outputPath = inputPath + ".go"
    }

    // 3. 检查是否需要编译
    if !config.Force && !needsCompile(inputPath, outputPath) {
        if config.Verbose {
            fmt.Printf("Skipping %s (up to date)\n", inputPath)
        }
        return nil
    }

    // 4. 读取源文件
    source, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input: %w", err)
    }

    // 5. 解析和编译
    result, err := compiler.Compile(source, config)
    if err != nil {
        return formatError(err, inputPath, source)
    }

    // 6. 写入输出
    return os.WriteFile(outputPath, result, 0644)
}

// 检查是否需要重新编译 (基于修改时间)
func needsCompile(input, output string) bool {
    inInfo, err := os.Stat(input)
    if err != nil {
        return true
    }

    outInfo, err := os.Stat(output)
    if err != nil {
        return true  // 输出不存在
    }

    return inInfo.ModTime().After(outInfo.ModTime())
}
```

---

## 方式二：文件监听模式

### 监听实现

```go
import (
    "github.com/fsnotify/fsnotify"
)

func watch(input string, config *Config) error {
    // 1. 初始编译
    if err := compile(input, config); err != nil {
        fmt.Fprintf(os.Stderr, "Initial compile error: %v\n", err)
    }

    // 2. 创建监听器
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return fmt.Errorf("create watcher: %w", err)
    }
    defer watcher.Close()

    // 3. 确定监听路径
    watchPath := input
    info, _ := os.Stat(input)
    if info != nil && !info.IsDir() {
        watchPath = filepath.Dir(input)
    }

    // 4. 添加监听
    if err := watcher.Add(watchPath); err != nil {
        return fmt.Errorf("add watch: %w", err)
    }

    fmt.Printf("Watching %s for changes...\n", watchPath)
    fmt.Println("Press Ctrl+C to stop")

    // 5. 事件循环
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return nil
            }

            if event.Op&fsnotify.Write == fsnotify.Write {
                if filepath.Ext(event.Name) == ".gox" {
                    fmt.Printf("\n%s changed, recompiling...\n", event.Name)
                    if err := compileFile(event.Name, config); err != nil {
                        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                    } else {
                        fmt.Println("Done")
                    }
                }
            }

        case err, ok := <-watcher.Errors:
            if !ok {
                return nil
            }
            fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
        }
    }
}
```

### 增量编译缓存

```go
type CompileCache struct {
    mu       sync.RWMutex
    entries  map[string]CacheEntry
}

type CacheEntry struct {
    InputHash  []byte    // 输入文件 hash
    OutputPath string
    ModTime    time.Time
}

func compileWithCache(inputPath string, config *Config, cache *CompileCache) error {
    // 1. 计算输入 hash
    source, _ := os.ReadFile(inputPath)
    hash := sha256.Sum256(source)

    // 2. 检查缓存
    if entry, ok := cache.Get(inputPath); ok {
        if bytes.Equal(entry.InputHash, hash[:]) {
            return nil  // 输入未变化，跳过编译
        }
    }

    // 3. 执行编译
    if err := compileFile(inputPath, config); err != nil {
        return err
    }

    // 4. 更新缓存
    cache.Set(inputPath, CacheEntry{
        InputHash:  hash[:],
        OutputPath: config.OutputPath,
        ModTime:    time.Now(),
    })

    return nil
}
```

---

## 方式三：IDE 插件

### VSCode 插件

#### 插件结构

```
gox-vscode/
├── package.json           # 插件配置
├── src/
│   ├── extension.ts       # 插件入口
│   ├── provider.ts        # 功能提供者
│   ├── commands.ts        # 命令
│   └── utils.ts           # 工具函数
├── syntaxes/
│   └── gox.tmLanguage.json # 语法高亮
└── README.md
```

#### package.json 配置

```json
{
  "name": "gox-vscode",
  "displayName": "GOX (Go JSX) Support",
  "description": "JSX syntax support for Go with go-ink",
  "version": "1.0.0",
  "publisher": "go-ink",
  "engines": { "vscode": "^1.80.0" },
  "categories": ["Programming Languages", "Formatters"],
  "activationEvents": [
    "onLanguage:gox",
    "onCommand:gox.compile"
  ],
  "main": "./out/extension.js",
  "contributes": {
    "languages": [{
      "id": "gox",
      "aliases": ["GOX", "Go JSX"],
      "extensions": [".gox"]
    }],
    "grammars": [{
      "language": "gox",
      "scopeName": "source.gox",
      "path": "./syntaxes/gox.tmLanguage.json"
    }],
    "commands": [
      { "command": "gox.compile", "title": "GOX: Compile Current File" },
      { "command": "gox.watch", "title": "GOX: Start Watch Mode" }
    ],
    "configuration": {
      "title": "GOX",
      "properties": {
        "gox.autoCompile": {
          "type": "boolean",
          "default": true,
          "description": "Automatically compile on save"
        },
        "gox.compilerPath": {
          "type": "string",
          "default": "gox-compiler",
          "description": "Path to gox-compiler binary"
        }
      }
    }
  }
}
```

#### 语法高亮配置

```json
// syntaxes/gox.tmLanguage.json
{
  "$schema": "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
  "name": "GOX",
  "patterns": [
    { "include": "#go-code" },
    { "include": "#jsx" }
  ],
  "repository": {
    "go-code": {
      "patterns": [
        {
          "match": "\\b(package|import|func|var|const|type|struct|if|else|for|range|return)\\b",
          "name": "keyword.control.go"
        },
        {
          "match": "\\b(string|int|float64|bool|error|nil)\\b",
          "name": "support.type.go"
        },
        {
          "match": "\"([^\"\\\\]|\\\\.)*\"",
          "name": "string.quoted.double.go"
        }
      ]
    },
    "jsx": {
      "patterns": [
        { "include": "#jsx-tag" },
        { "include": "#jsx-expression" }
      ]
    },
    "jsx-tag": {
      "begin": "(<)([A-Z][a-zA-Z0-9]*)",
      "beginCaptures": {
        "1": { "name": "punctuation.definition.tag.begin.jsx" },
        "2": { "name": "entity.name.tag.jsx" }
      },
      "end": "(</)($2)(>)|(/>)",
      "patterns": [
        { "include": "#jsx-attribute" },
        { "include": "#jsx" }
      ]
    },
    "jsx-attribute": {
      "match": "([a-zA-Z][a-zA-Z0-9]*)\\s*(=)\\s*",
      "captures": {
        "1": { "name": "entity.other.attribute-name.jsx" },
        "2": { "name": "keyword.operator.assignment.jsx" }
      }
    },
    "jsx-expression": {
      "begin": "\\{",
      "end": "\\}",
      "patterns": [
        { "include": "#go-code" },
        { "include": "#jsx" }
      ]
    }
  },
  "scopeName": "source.gox"
}
```

---

### LSP (Language Server Protocol)

#### LSP Server 结构

```
gox-lsp/
├── cmd/
│   └── gox-lsp/
│       └── main.go        # LSP server 入口
├── internal/
│   ├── server/
│   │   ├── server.go      # LSP server 实现
│   │   ├── handler.go     # 请求处理器
│   │   └── document.go    # 文档管理
│   ├── analysis/
│   │   ├── parser.go      # 解析
│   │   ├── diagnostics.go # 诊断
│   │   ├── completion.go  # 补全
│   │   └── hover.go       # 悬停信息
│   └── cache/
│       └── cache.go       # 分析缓存
└── go.mod
```

#### LSP Server 核心

```go
type GOXServer struct {
    client    protocol.Client
    documents map[string]*Document
    cache     *AnalysisCache
    mu        sync.RWMutex
}

// 文档打开
func (s *GOXServer) DidOpen(params protocol.DidOpenTextDocumentParams) error {
    doc := &Document{
        URI:     string(params.TextDocument.URI),
        Content: params.TextDocument.Text,
        Version: params.TextDocument.Version,
    }

    s.analyze(doc)

    s.mu.Lock()
    s.documents[doc.URI] = doc
    s.mu.Unlock()

    s.publishDiagnostics(doc)
    return nil
}

// 补全
func (s *GOXServer) Completion(params protocol.CompletionParams) (*protocol.CompletionList, error) {
    s.mu.RLock()
    doc := s.documents[string(params.TextDocument.URI)]
    s.mu.RUnlock()

    items := s.computeCompletions(doc, params.Position)
    return &protocol.CompletionList{Items: items}, nil
}
```

---

## 构建系统集成

### Makefile

```makefile
.PHONY: all build clean generate

all: generate build

generate:
    gox-compiler ./ui/

build: generate
    go build -o bin/app ./...

run: generate
    go run ./...

clean:
    rm -f ui/*.gox.go
```

### Taskfile

```yaml
version: '3'

tasks:
  generate:
    cmds:
      - gox-compiler ./ui/
    sources:
      - ui/**/*.gox
    generates:
      - ui/**/*.gox.go

  build:
    deps: [generate]
    cmds:
      - go build -o bin/app ./...

  dev:
    deps: [generate]
    cmds:
      - gox-compiler -w ./ui/
      - go run ./...
```

---

## 错误报告

### 错误格式

```
Error: mismatched closing tag

  5 |     <Box>
  6 |         <Text>Hello</Text>
  7 |     </Text>
    |     ^^^^^^^^
  8 | </Box>

ui.gox:7:5: expected </Box>, found </Text>
```

### 错误格式化实现

```go
func formatParseError(err *ParseError, path string, source []byte) error {
    lines := bytes.Split(source, []byte("\n"))

    var buf bytes.Buffer

    fmt.Fprintf(&buf, "Error: %s\n\n", err.Message)

    start := max(0, err.Line-2)
    end := min(len(lines), err.Line+2)

    for i := start; i < end; i++ {
        fmt.Fprintf(&buf, "  %d | ", i+1)
        buf.Write(lines[i])
        buf.WriteByte('\n')

        if i == err.Line-1 {
            buf.WriteString("    | ")
            buf.Write(bytes.Repeat([]byte(" "), err.Column-1))
            buf.Write(bytes.Repeat([]byte("^"), err.Length))
            buf.WriteByte('\n')
        }
    }

    fmt.Fprintf(&buf, "\n%s:%d:%d: %s\n",
        path, err.Line, err.Column, err.Message)

    return errors.New(buf.String())
}
```

---

## Source Map

```go
type SourceMap struct {
    File     string        // 生成的 .go 文件
    Source   string        // 源 .gox 文件
    Mappings []Mapping     // 位置映射
}

type Mapping struct {
    GeneratedLine   int    // 生成代码行号
    GeneratedColumn int    // 生成代码列号
    SourceLine      int    // 源代码行号
    SourceColumn    int    // 源代码列号
    SourceName      string // 标识符名 (可选)
}
```

---

## 总结

### 工具链集成方式

| 方式 | 用途 | 特点 |
|-----|------|------|
| go generate | CI/CD | 标准 Go 工作流 |
| 文件监听 | 开发时 | 自动编译，提升体验 |
| VSCode 插件 | 编辑器 | 语法高亮、补全、错误提示 |
| GoLand 插件 | IDE | 完整 IDE 支持 |
| LSP | 跨编辑器 | Vim/Emacs/VSCode 通用 |

### 核心功能

1. **编译**
   - 单文件/目录编译
   - 增量编译 (修改时间/hash)
   - 并行编译

2. **监听**
   - 文件变化检测
   - 自动重新编译
   - 错误实时反馈

3. **IDE 支持**
   - 语法高亮
   - 代码补全
   - 错误诊断
   - 悬停信息

4. **调试**
   - Source Map
   - 位置映射
   - 错误报告