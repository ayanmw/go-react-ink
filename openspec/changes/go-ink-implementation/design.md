# Design: Go-Ink Architecture

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Go-Ink 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   编译时:                                                    │
│                                                             │
│   .gox 文件 ──▶ Scanner ──▶ Lexer ──▶ Parser ──▶ AST       │
│                                               │              │
│                                               ▼              │
│                                          CodeGen ──▶ .go    │
│                                                             │
│   运行时:                                                    │
│                                                             │
│   用户组件 ──▶ Reconciler ──▶ Host Config                    │
│                    │                                        │
│         ┌──────────┴──────────┐                            │
│         ▼                     ▼                            │
│   Flexbox Layout      Terminal Renderer                     │
│         │                     │                            │
│         └──────────┬──────────┘                            │
│                    ▼                                        │
│              ANSI Output                                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Module Structure

```
go-ink/
├── cmd/
│   └── gox-compiler/          # CLI 工具
│       └── main.go
│
├── internal/
│   ├── scanner/               # 源码扫描
│   ├── lexer/                 # 词法分析
│   ├── parser/                # 语法分析
│   ├── codegen/               # 代码生成
│   └── compiler/              # 编译器集成
│
├── pkg/
│   ├── reconciler/            # Fiber 协调器
│   │   ├── fiber.go
│   │   ├── scheduler.go
│   │   ├── workloop.go
│   │   └── hostconfig.go
│   │
│   ├── layout/                # Flexbox 布局
│   │   ├── node.go
│   │   ├── calculator.go
│   │   └── yoga.go            # CGO 绑定 (Phase 1)
│   │
│   ├── renderer/              # 终端渲染
│   │   ├── buffer.go
│   │   ├── diff.go
│   │   └── ansi.go
│   │
│   ├── components/            # 内置组件
│   │   ├── box.go
│   │   ├── text.go
│   │   └── ...
│   │
│   ├── hooks/                 # Hooks
│   │   ├── state.go
│   │   ├── effect.go
│   │   └── input.go
│   │
│   └── core/                  # 核心类型
│       ├── element.go
│       ├── props.go
│       └── styles.go
│
├── tools/
│   ├── vscode/                # VSCode 插件
│   └── lsp/                   # LSP Server
│
└── reference/
    └── react-ink/             # 原版参考
```

## Data Flow

### 编译流程

```
.gox 文件
    │
    ▼
┌─────────────┐
│   Scanner   │  识别 JSX 区域
└─────────────┘
    │
    ▼
┌─────────────┐
│   Lexer     │  Token 流
└─────────────┘
    │
    ▼
┌─────────────┐
│   Parser    │  AST
└─────────────┘
    │
    ▼
┌─────────────┐
│  CodeGen    │  Go 源码
└─────────────┘
    │
    ▼
.go 文件
```

### 渲染流程

```
用户组件函数
    │
    ▼
┌─────────────┐
│ Reconciler  │  Fiber 树构建
│  (Render)   │  Diff 计算
└─────────────┘
    │
    ▼
┌─────────────┐
│   Commit    │  副作用提交
└─────────────┘
    │
    ▼
┌─────────────┐
│   Layout    │  Flexbox 计算
└─────────────┘
    │
    ▼
┌─────────────┐
│  Renderer   │  输出缓冲构建
│             │  Diff 计算
└─────────────┘
    │
    ▼
┌─────────────┐
│   ANSI      │  最小化输出
└─────────────┘
    │
    ▼
终端显示
```

## Key Interfaces

### Element

```go
type Element interface {
    Type() string
    Props() Props
    Children() []Element
    Key() any
}
```

### HostConfig

```go
type HostConfig interface {
    CreateInstance(typeName string, props Props) Instance
    CreateTextInstance(text string) Instance
    UpdateInstance(instance Instance, props Props)
    AppendChild(parent, child Instance)
    RemoveChild(parent, child Instance)
    ScheduleLayout(instance Instance)
    ScheduleRender()
}
```

### Component

```go
type Component func(props Props) Element

// 示例
func Counter(props Props) Element {
    count, setCount := useState(0)
    
    return Box(Props{},
        Text(Props{}, fmt.Sprintf("Count: %d", count)),
    )
}
```

## Error Handling

### 编译时错误

```
Error: mismatched closing tag

  5 |     <Box>
  6 |         <Text>Hello</Text>
  7 |     </Text>
    |     ^^^^^^^^
  8 | </Box>

ui.gox:7:5: expected </Box>, found </Text>
```

### 运行时错误

- 组件渲染错误 → 错误边界捕获
- 布局计算错误 → 回退到默认布局
- 终端输出错误 → 优雅降级

## Performance Considerations

1. **增量渲染**: 只更新变化的部分
2. **布局缓存**: 相同 Props 不重新计算
3. **ANSI 优化**: 合并连续单元格，减少序列
4. **时间切片**: 避免阻塞主线程

## Testing Strategy

1. **单元测试**: 每个模块独立测试
2. **集成测试**: 编译器 + 运行时
3. **对比测试**: 与 React Ink 结果对比
4. **性能测试**: Benchmark 对比

## Security Considerations

- 无网络请求
- 无文件系统写入 (除编译输出)
- 无用户代码执行
- ANSI 序列安全处理
