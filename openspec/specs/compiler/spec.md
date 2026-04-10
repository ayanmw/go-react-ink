# compiler Specification

## Purpose
定义 Go-Ink 编译器的设计和实现，将 .gox 文件 (JSX) 编译为标准 Go 代码。

## Requirements

### Requirement: 编译器识别 JSX 语法
编译器 SHALL 识别 .gox 文件中的 JSX 语法并转换为 Go 函数调用。

#### Scenario: 识别 JSX 元素
- **WHEN** 编译器遇到 `<Box>...</Box>` 语法
- **THEN** 创建 JSXElement AST 节点

#### Scenario: 识别 JSX 属性
- **WHEN** 编译器遇到 `<Box flexDirection="column">`
- **THEN** 解析属性名为 `flexDirection`，属性值为 `"column"`

#### Scenario: 识别 JSX 表达式
- **WHEN** 编译器遇到 `<Text>{count}</Text>`
- **THEN** 解析 `{count}` 为 JSXExpressionContainer

### Requirement: 编译器支持嵌套 JSX
编译器 SHALL 支持任意深度的 JSX 嵌套。

#### Scenario: 嵌套元素
- **WHEN** 编译器遇到 `<Box><Text>Hello</Text></Box>`
- **THEN** 创建嵌套的 AST 结构

### Requirement: 编译器混合 Go 代码和 JSX
编译器 SHALL 在同一个文件中处理 Go 代码和 JSX。

#### Scenario: 函数返回 JSX
- **WHEN** 编译器遇到 `func App() Element { return <Box/> }`
- **THEN** 正确解析函数签名和 JSX 返回值

### Requirement: JSX 元素转换为函数调用
编译器 SHALL 将 JSX 元素转换为对应的 Go 函数调用。

#### Scenario: 简单元素
- **GIVEN** `<Text>Hello</Text>`
- **WHEN** 编译器处理该 JSX
- **THEN** 生成 `Text(Props{}, TextChildren("Hello"))`

#### Scenario: 带属性元素
- **GIVEN** `<Box flexDirection="column" width={100}>`
- **WHEN** 编译器处理该 JSX
- **THEN** 生成 `Box(Props{"flexDirection": "column", "width": 100})`

### Requirement: JSX 表达式转换为值
编译器 SHALL 将 JSX 中的 `{expression}` 转换为 Go 表达式。

#### Scenario: 变量表达式
- **GIVEN** `<Text>{count}</Text>`
- **WHEN** 编译器处理该 JSX
- **THEN** 生成 `Text(Props{}, count)`

#### Scenario: 函数调用表达式
- **GIVEN** `<Text>{getName()}</Text>`
- **WHEN** 编译器处理该 JSX
- **THEN** 生成 `Text(Props{}, getName())`

### Requirement: JSX 子元素转换为参数
编译器 SHALL 将 JSX 子元素转换为函数参数。

#### Scenario: 嵌套子元素
- **GIVEN** `<Box><Text>A</Text><Text>B</Text></Box>`
- **WHEN** 编译器处理该 JSX
- **THEN** 生成 `Box(Props{}, Text(Props{}, "A"), Text(Props{}, "B"))`

### Requirement: Props 转换
编译器 SHALL 将 JSX 属性转换为 Props map。

#### Scenario: 字符串属性
- **GIVEN** `color="green"`
- **WHEN** 编译器处理该属性
- **THEN** 生成 `"color": "green"`

#### Scenario: 表达式属性
- **GIVEN** `width={100}`
- **WHEN** 编译器处理该属性
- **THEN** 生成 `"width": 100`

#### Scenario: 布尔属性
- **GIVEN** `<Text bold>`
- **WHEN** 编译器处理该属性
- **THEN** 生成 `"bold": true`

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Compiler 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: .gox 源文件                                          │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │  Scanner    │  字符扫描                                 │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │   Lexer     │  词法分析 → Token 流                      │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │   Parser    │  语法分析 → AST                           │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  CodeGen    │  代码生成 → .go 文件                      │
│   └─────────────┘                                          │
│                                                             │
│   输出: 标准 Go 代码                                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Token Types

```go
type TokenType int

const (
    // Go Tokens
    TokenIdent TokenType = iota
    TokenNumber
    TokenString
    TokenOperator
    TokenKeyword
    TokenPunctuation
    
    // JSX Tokens
    TokenJSXOpenTag      // <
    TokenJSXCloseTag     // >
    TokenJSXSelfClose    // />
    TokenJSXEndTag       // </
    TokenJSXText         // 文本内容
    TokenJSXExprStart    // {
    TokenJSXExprEnd      // }
)
```

## AST Nodes

```go
// JSX 元素节点
type JSXElement struct {
    OpeningElement *JSXOpeningElement
    Children       []JSXChild
    ClosingElement *JSXClosingElement
}

// JSX 开始标签
type JSXOpeningElement struct {
    Name       *JSXIdentifier
    Attributes []JSXAttribute
    SelfClosing bool
}

// JSX 属性
type JSXAttribute struct {
    Name  *JSXIdentifier
    Value JSXAttributeValue
}

// JSX 表达式容器
type JSXExpressionContainer struct {
    Expression Expression
}
```

## Transform Rules

### Element Transform

```
JSX:
<Element prop="value">child</Element>

Go:
Element(Props{"prop": "value"}, "child")
```

### Nested Element Transform

```
JSX:
<Box>
  <Text>A</Text>
  <Text>B</Text>
</Box>

Go:
Box(Props{},
  Text(Props{}, "A"),
  Text(Props{}, "B"),
)
```

### Expression Transform

```
JSX:
<Text>{count + 1}</Text>

Go:
Text(Props{}, count + 1)
```

### Fragment Transform

```
JSX:
<>
  <Text>A</Text>
  <Text>B</Text>
</>

Go:
Fragment(Props{},
  Text(Props{}, "A"),
  Text(Props{}, "B"),
)
```

## Edge Cases

### Spread Attributes
```
JSX:
<Box {...props} />

Go:
Box(mergeProps(props))
```

### Conditional Rendering
```
JSX:
{show && <Text>Hello</Text>}

Go:
func() Element {
  if show {
    return Text(Props{}, "Hello")
  }
  return nil
}()
```

### Map/List Rendering
```
JSX:
{items.map(item => <Text>{item.name}</Text>)}

Go:
MapElements(items, func(item Item) Element {
  return Text(Props{}, item.name)
})
```

## Toolchain Integration

### gox-compiler CLI

```bash
# 编译单个文件
gox-compiler input.gox -o output.go

# 编译目录
gox-compiler ./src/...

# 监听模式
gox-compiler -watch ./src/
```

### go:generate 注释

```go
//go:generate gox-compiler $GOFILE
```

## Module Structure

```
compiler/
├── scanner/
│   └── scanner.go        # 字符扫描
├── lexer/
│   ├── lexer.go          # 词法分析
│   └── token.go          # Token 定义
├── parser/
│   ├── parser.go         # 语法分析
│   ├── jsx_parser.go     # JSX 解析
│   └── ast.go            # AST 定义
├── codegen/
│   ├── codegen.go        # 代码生成
│   └── template.go       # 代码模板
└── cli/
    └── main.go           # CLI 入口
```

## React Ink 对齐

| 特性 | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| JSX 语法 | ✅ Babel | ✅ gox-compiler | 编译时转换 |
| 嵌套元素 | ✅ | ✅ | AST 嵌套结构 |
| 表达式 `{}` | ✅ | ✅ | Go 表达式 |
| 属性展开 `{...props}` | ✅ | ✅ | mergeProps |
| 条件渲染 `{cond && <A/>}` | ✅ | ✅ | if 表达式 |
| 列表渲染 `{items.map()}` | ✅ | ✅ | MapElements |
| Fragment `<>` | ✅ | ✅ | Fragment 组件 |