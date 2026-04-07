# JSX 解析器设计

> gox-compiler 的核心：将 JSX 语法解析为 AST

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  JSX 解析器架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: .gox 文件源码                                        │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐│
│   │ Scanner │───▶│ Lexer   │───▶│ Parser  │───▶│ AST     ││
│   │ 源码扫描 │    │ 词法分析 │    │ 语法分析 │    │ 抽象语法树││
│   └─────────┘    └─────────┘    └─────────┘    └─────────┘│
│                                                             │
│   输出: JSX AST                                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 第一阶段：Scanner (源码扫描)

### 职责

识别 Go 代码中的 JSX 区域

### 核心任务

1. 识别 JSX 起始标记: `return (` 后的 `<`
2. 识别 JSX 结束标记: 匹配的 `)`
3. 提取纯 Go 代码区域 (保持不变)
4. 提取 JSX 区域 (交给 Lexer 处理)

### 输入输出示例

```
输入:
func App() ink.Element {
    name := "World"
    return (
        <Text>Hello {name}</Text>
    )
}

输出: []Segment
[
  {Type: GoCode, Content: "func App... return ("},
  {Type: JSX, Content: "<Text>Hello {name}</Text>"},
  {Type: GoCode, Content: ") }"},
]
```

### Scanner 状态机

```
                    ┌─────────────┐
                    │   Start     │
                    │  (Go代码)   │
                    └─────────────┘
                          │
                          │ 检测到 "return (" 或 "("
                          │ 后跟 "<"
                          ▼
                    ┌─────────────┐
              ┌────▶│   JSX       │◀────┐
              │     │   区域      │     │
              │     └─────────────┘     │
              │           │             │
              │           │ 嵌套 "("   │ 匹配 ")"
              │           │ 或 "{"     │
              │           ▼             │
              │     ┌─────────────┐     │
              │     │  Nested     │─────┘
              │     │  (嵌套)     │
              │     └─────────────┘
              │           │
              │           │ 嵌套结束
              └───────────┘
                          │
                          │ JSX 结束 (匹配的 ")")
                          ▼
                    ┌─────────────┐
                    │   End       │
                    │  (Go代码)   │
                    └─────────────┘
```

---

## 第二阶段：Lexer (词法分析)

### 职责

将 JSX 区域转换为 Token 流

### Token 类型定义

```go
type TokenType int

const (
    // 标签相关
    TokenTagStart      TokenType = iota  // <
    TokenTagEnd                       // >
    TokenTagSelfClose                 // />
    TokenTagCloseStart                // </
    TokenTagName                     // Box, Text

    // 属性相关
    TokenAttrName                    // flexDirection
    TokenAttrValue                   // "column"
    TokenExprStart                   // {
    TokenExprEnd                     // }
    TokenExprContent                 // count * 2
    TokenSpread                      // ...props

    // 内容相关
    TokenText                        // Hello World

    // 特殊
    TokenEOF                         // 文件结束
    TokenError                       // 错误
)

type Token struct {
    Type     TokenType
    Value    string
    Position Position  // 行号、列号
}
```

### Lexer 示例

```
输入:
<Box flexDirection="column" padding={10}>
    <Text color="green">Hello {name}</Text>
</Box>

输出 Token 流:
 #  │ Type              │ Value
────┼───────────────────┼────────────────────────────────
 1  │ TokenTagStart     │ "<"
 2  │ TokenTagName      │ "Box"
 3  │ TokenAttrName     │ "flexDirection"
 4  │ TokenAttrValue    │ "column"
 5  │ TokenAttrName     │ "padding"
 6  │ TokenExprStart    │ "{"
 7  │ TokenExprContent  │ "10"
 8  │ TokenExprEnd      │ "}"
 9  │ TokenTagEnd       │ ">"
10  │ TokenTagStart     │ "<"
11  │ TokenTagName      │ "Text"
12  │ TokenAttrName     │ "color"
13  │ TokenAttrValue    │ "green"
14  │ TokenTagEnd       │ ">"
15  │ TokenText         │ "Hello "
16  │ TokenExprStart    │ "{"
17  │ TokenExprContent  │ "name"
18  │ TokenExprEnd      │ "}"
19  │ TokenTagCloseStart│ "</"
20  │ TokenTagName      │ "Text"
21  │ TokenTagEnd       │ ">"
22  │ TokenTagCloseStart│ "</"
23  │ TokenTagName      │ "Box"
24  │ TokenTagEnd       │ ">"
25  │ TokenEOF          │ ""
```

---

## 第三阶段：Parser (语法分析)

### 职责

将 Token 流转换为 AST (抽象语法树)

### AST 节点定义

```go
// AST 节点接口
type Node interface {
    nodeType()
}

// JSX 元素节点
type JSXElement struct {
    TagName   string
    Props     []Prop
    Children  []Node
    Position  Position
}

// 属性
type Prop struct {
    Name      string
    Value     PropValue
    Position  Position
}

// 属性值 (字符串字面量 或 表达式)
type PropValue interface { propValue() }

type StringLiteral struct {
    Value string
}

type Expression struct {
    Expr string  // 原始表达式文本
}

type SpreadExpr struct {
    Expr string  // ...props 中的 props
}

// JSX 文本节点
type JSXText struct {
    Value    string
    Position Position
}

// JSX 表达式容器 {expression}
type JSXExpressionContainer struct {
    Expr     string
    Position Position
}

// JSX Fragment <>...</>
type JSXFragment struct {
    Children []Node
}
```

### Parser 递归下降算法

```go
// 主入口
func (p *Parser) Parse() *JSXElement {
    return p.parseElement()
}

// 解析元素: <TagName props>children</TagName>
func (p *Parser) parseElement() *JSXElement {
    elem := &JSXElement{}

    // 1. 期望 <
    p.expect(TokenTagStart)

    // 2. 解析标签名
    elem.TagName = p.expect(TokenTagName).Value

    // 3. 解析属性
    elem.Props = p.parseProps()

    // 4. 检查自闭合或解析子节点
    if p.peek().Type == TokenTagSelfClose {
        p.consume()  // 消耗 />
        return elem
    }

    p.expect(TokenTagEnd)  // 消耗 >

    // 5. 解析子节点
    elem.Children = p.parseChildren(elem.TagName)

    return elem
}

// 解析属性
func (p *Parser) parseProps() []Prop {
    props := []Prop{}

    for p.peek().Type != TokenTagEnd &&
        p.peek().Type != TokenTagSelfClose {

        // 处理展开属性 {...props}
        if p.peek().Type == TokenSpread {
            props = append(props, p.parseSpreadProp())
            continue
        }

        // 解析普通属性
        prop := Prop{}
        prop.Name = p.expect(TokenAttrName).Value

        // 属性值
        if p.peek().Type == TokenExprStart {
            // 表达式值 {expr}
            p.consume()  // {
            prop.Value = Expression{
                Expr: p.expect(TokenExprContent).Value,
            }
            p.expect(TokenExprEnd)  // }
        } else {
            // 字符串值 "value"
            prop.Value = StringLiteral{
                Value: p.expect(TokenAttrValue).Value,
            }
        }

        props = append(props, prop)
    }

    return props
}

// 解析子节点
func (p *Parser) parseChildren(parentTag string) []Node {
    children := []Node{}

    for {
        tok := p.peek()

        // 检查结束标签 </parentTag>
        if tok.Type == TokenTagCloseStart {
            p.consume()  // </
            tagName := p.expect(TokenTagName).Value
            p.expect(TokenTagEnd)  // >

            if tagName != parentTag {
                p.error("mismatched closing tag")
            }
            break
        }

        // 解析子节点
        switch tok.Type {
        case TokenTagStart:
            // 嵌套元素
            children = append(children, p.parseElement())

        case TokenText:
            // 文本节点
            children = append(children, &JSXText{
                Value: p.consume().Value,
            })

        case TokenExprStart:
            // 表达式容器 {expr}
            children = append(children, p.parseExpr())
        }
    }

    return children
}
```

---

## 复杂语法处理

### 1. 条件渲染

```
语法: {condition && <Element/>}
      {condition ? <A/> : <B/>}

JSX 输入:
{show && <Text>Visible</Text>}

Token 流:
{, show && <Text>Visible</Text>, }

AST:
JSXExpressionContainer{
    Expr: "show && <Text>Visible</Text>"
}

关键: 表达式内部可能包含 JSX，需要递归解析
```

### 2. 列表渲染

```
语法: {items.map(item => <Item key={item.id} {...item}/>)}

JSX 输入:
{items.map((item, i) => (
    <Item key={i} {...item}/>
))}

AST:
JSXExpressionContainer{
    Expr: "items.map((item, i) => (
        <Item key={i} {...item}/>
    ))"
}

关键点:
1. 识别箭头函数内的 JSX
2. 处理展开属性 {...item}
3. 保留 key 属性 (React 特殊属性)
```

### 3. Fragment

```
语法: <>...</> 或 <Fragment>...</Fragment>

JSX 输入:
<>
    <Text>A</Text>
    <Text>B</Text>
</>

AST:
JSXFragment{
    Children: [
        JSXElement{TagName: "Text", ...},
        JSXElement{TagName: "Text", ...},
    ],
}

Lexer 特殊处理:
- "<>" 识别为 TokenFragmentStart
- "</>" 识别为 TokenFragmentEnd
```

### 4. 属性展开

```
语法: <Element {...props}/>

JSX 输入:
<Box {...baseProps} padding={10}/>

AST:
JSXElement{
    Props: [
        SpreadExpr{Expr: "baseProps"},
        {Name: "padding", Value: Expression{Expr: "10"}},
    ],
}
```

---

## 错误处理

### 错误类型

1. **词法错误**
   - 未闭合的字符串: `"hello`
   - 未闭合的表达式: `{count`
   - 无效字符: `<@`

2. **语法错误**
   - 未闭合的标签: `<Box>...</Text>`
   - 标签不匹配: `<Box></Text>`
   - 缺少属性值: `<Box color>`
   - 无效的子节点位置

### 错误报告格式

```
Error: mismatched closing tag

  3 |     <Box>
  4 |         <Text>Hello</Text>
  5 |     </Text>
    |     ^^^^^^^^
  6 |

expected </Box>, found </Text>
```

### 错误恢复策略

- 同步点: 标签结束 `>` 或 `/>`
- 跳过错误 token 直到同步点
- 继续解析以报告更多错误

---

## 模块结构

```
gox-compiler/
├── cmd/
│   └── gox/
│       └── main.go        # CLI 入口
│
├── scanner/
│   ├── scanner.go         # 源码扫描
│   └── segment.go         # 代码段定义
│
├── lexer/
│   ├── lexer.go           # 词法分析器
│   ├── token.go           # Token 定义
│   └── state.go           # 状态机
│
├── parser/
│   ├── parser.go          # 语法分析器
│   ├── ast.go             # AST 节点定义
│   └── error.go           # 错误处理
│
├── codegen/
│   ├── generator.go       # Go 代码生成
│   └── template.go        # 代码模板
│
└── internal/
    └── position.go        # 位置信息
```

---

## 总结

### 核心要点

1. **三阶段处理**: Scanner → Lexer → Parser
   - Scanner: 识别 Go 代码中的 JSX 区域
   - Lexer: Token 流生成
   - Parser: AST 构建

2. **递归下降解析器**
   - 自然处理嵌套结构
   - 易于理解和扩展
   - 良好的错误报告

3. **表达式处理策略**
   - 保留原始表达式文本
   - 检测表达式内的嵌套 JSX
   - 交给代码生成阶段处理语义

4. **错误处理**
   - 位置信息跟踪
   - 友好的错误报告
   - 错误恢复以继续解析
