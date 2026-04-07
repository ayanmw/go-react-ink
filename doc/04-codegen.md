# 代码生成器设计

> 将 JSX AST 转换为 Go 源码

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  Code Generator 架构                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入:                                                     │
│   ├── 源文件信息 (SourceFile)                               │
│   ├── JSX AST ([]Node)                                      │
│   └── 代码段 ([]Segment)                                    │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │ Source File │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐    ┌─────────────┐                      │
│   │  Segment    │───▶│   Output    │                      │
│   │  Iterator   │    │   Buffer    │                      │
│   └─────────────┘    └─────────────┘                      │
│         │                    ▲                             │
│         │                    │                             │
│         ▼                    │                             │
│   ┌─────────────┐    ┌─────────────┐                      │
│   │  Go Code?   │───▶│   Write     │                      │
│   │  or JSX?    │    │   Direct    │                      │
│   └─────────────┘    └─────────────┘                      │
│         │                                                  │
│         │ JSX                                              │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Generator  │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ├──▶ generateElement()                              │
│         ├──▶ generateProps()                                │
│         ├──▶ generateChildren()                             │
│         └──▶ generateExpression()                           │
│                                                             │
│   输出: 生成的 Go 源码                                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心数据结构

```go
// 生成器主结构
type Generator struct {
    // 输出缓冲
    buf       *bytes.Buffer

    // 缩进管理
    indent   int
    indentStr string  // 默认 "    " (4空格)

    // 导入收集
    imports  *ImportCollector

    // 包信息
    pkgName  string

    // 配置
    config   Config
}

// 导入收集器
type ImportCollector struct {
    standard  map[string]struct{}  // 标准库
    thirdParty map[string]struct{} // 第三方库
    local     map[string]struct{}  // 本地包
    aliases   map[string]string    // 别名
}

// 生成配置
type Config struct {
    // 是否生成调试注释
    DebugComments bool

    // Props 结构体名称映射
    PropsNaming func(tagName string) string

    // 组件包路径
    InkPackage string  // 默认 "github.com/xxx/go-ink"
}
```

---

## 核心方法

### 主入口

```go
func (g *Generator) Generate(file *SourceFile) ([]byte, error) {
    // 1. 初始化
    g.buf = &bytes.Buffer{}
    g.imports = NewImportCollector()

    // 2. 写入包声明
    g.write("package %s\n\n", file.PackageName)

    // 3. 处理代码段
    for _, seg := range file.Segments {
        switch seg.Type {
        case SegmentGoCode:
            // Go 代码直接输出
            g.write("%s", seg.Content)

        case SegmentJSX:
            // JSX 需要转换
            ast, err := ParseJSX(seg.Content)
            if err != nil {
                return nil, err
            }
            g.generateElement(ast)
        }
    }

    // 4. 生成导入块 (插入到文件开头)
    output := g.buf.String()
    if len(g.imports.standard) > 0 ||
       len(g.imports.thirdParty) > 0 {
        importBlock := g.imports.Generate()
        output = insertImports(output, importBlock)
    }

    // 5. 格式化输出
    formatted, err := format.Source([]byte(output))
    if err != nil {
        return []byte(output), nil
    }

    return formatted, nil
}
```

### 元素生成

```go
func (g *Generator) generateElement(elem *JSXElement) {
    // 1. 确定组件函数名
    componentFunc := g.resolveComponent(elem.TagName)

    // 2. 生成 Props
    propsCode := g.generateProps(elem.Props)

    // 3. 生成 Children
    childrenCode := g.generateChildren(elem.Children)

    // 4. 组装输出
    if len(elem.Children) == 0 {
        // 无子节点
        g.write("%s(%s)", componentFunc, propsCode)
    } else {
        // 有子节点
        g.write("%s(%s,\n", componentFunc, propsCode)
        g.increaseIndent()
        g.writeIndent()
        g.write("%s,\n", childrenCode)
        g.decreaseIndent()
        g.writeIndent()
        g.write(")")
    }
}

// 解析组件名
func (g *Generator) resolveComponent(tagName string) string {
    // 内置组件
    builtins := map[string]string{
        "Box":    "ink.Box",
        "Text":   "ink.Text",
        "Spacer": "ink.Spacer",
        "Newline": "ink.Newline",
    }

    if name, ok := builtins[tagName]; ok {
        g.imports.Add(g.config.InkPackage)
        return name
    }

    // 自定义组件 (首字母大写，Go 可导出)
    return tagName
}
```

### 属性生成

```go
func (g *Generator) generateProps(props []Prop) string {
    if len(props) == 0 {
        return "ink.Props{}"
    }

    var buf bytes.Buffer
    buf.WriteString("ink.Props{\n")

    for _, prop := range props {
        // 处理展开属性
        if spread, ok := prop.Value.(SpreadExpr); ok {
            continue
        }

        // 生成属性名 (Go 风格: 首字母大写)
        goName := g.propToGoName(prop.Name)

        // 生成属性值
        goValue := g.generatePropValue(prop.Value, prop.Name)

        fmt.Fprintf(&buf, "    %s: %s,\n", goName, goValue)
    }

    buf.WriteString("    }")
    return buf.String()
}
```

### 属性名转换

```go
func (g *Generator) propToGoName(name string) string {
    // 特殊映射
    special := map[string]string{
        "flexDirection": "FlexDirection",
        "justifyContent": "JustifyContent",
        "alignItems":     "AlignItems",
        "flexWrap":       "FlexWrap",
        "flexGrow":       "FlexGrow",
        "flexShrink":     "FlexShrink",
        "flexBasis":      "FlexBasis",
        "paddingTop":     "PaddingTop",
        "paddingBottom":  "PaddingBottom",
        "paddingLeft":    "PaddingLeft",
        "paddingRight":   "PaddingRight",
        "marginTop":      "MarginTop",
        "marginBottom":   "MarginBottom",
        "marginLeft":     "MarginLeft",
        "marginRight":    "MarginRight",
        "borderStyle":    "BorderStyle",
        "borderColor":    "BorderColor",
    }

    if goName, ok := special[name]; ok {
        return goName
    }

    // 默认: 首字母大写
    return strings.Title(name)
}
```

### 枚举值转换

```go
func (g *Generator) generateStringLiteral(value, propName string) string {
    // 枚举类型属性
    enumProps := map[string]map[string]string{
        "flexDirection": {
            "row":    "ink.Row",
            "column": "ink.Column",
            "row-reverse":    "ink.RowReverse",
            "column-reverse": "ink.ColumnReverse",
        },
        "justifyContent": {
            "flex-start":    "ink.FlexStart",
            "flex-end":      "ink.FlexEnd",
            "center":        "ink.Center",
            "space-between": "ink.SpaceBetween",
            "space-around":  "ink.SpaceAround",
            "space-evenly":  "ink.SpaceEvenly",
        },
        "alignItems": {
            "stretch":     "ink.Stretch",
            "flex-start":  "ink.FlexStart",
            "flex-end":    "ink.FlexEnd",
            "center":      "ink.Center",
            "baseline":    "ink.Baseline",
        },
        "color": {
            "red":     "ink.Red",
            "green":   "ink.Green",
            "blue":    "ink.Blue",
            "yellow":  "ink.Yellow",
            "magenta": "ink.Magenta",
            "cyan":    "ink.Cyan",
            "white":   "ink.White",
            "black":   "ink.Black",
            "gray":    "ink.Gray",
        },
    }

    if enums, ok := enumProps[propName]; ok {
        if goValue, ok := enums[value]; ok {
            return goValue
        }
    }

    // 普通字符串
    return fmt.Sprintf("%q", value)
}
```

---

## 子节点生成

### 核心逻辑

```go
func (g *Generator) generateChildren(children []Node) string {
    if len(children) == 0 {
        return ""
    }

    // 分类处理
    var elements []string
    var hasComplexExpr bool

    for _, child := range children {
        switch c := child.(type) {
        case *JSXElement:
            // 子元素
            elements = append(elements, g.generateElementCode(c))

        case *JSXText:
            // 纯文本
            elements = append(elements,
                fmt.Sprintf("ink.TextContent(%q)", c.Value))

        case *JSXExpressionContainer:
            // 表达式
            code := g.generateChildExpr(c.Expr)
            if code != "" {
                elements = append(elements, code)
            }
            if g.isComplexExpr(c.Expr) {
                hasComplexExpr = true
            }

        case *JSXFragment:
            // Fragment
            elements = append(elements, g.generateFragment(c))
        }
    }

    // 组装输出
    if len(elements) == 1 && !hasComplexExpr {
        return elements[0]
    }

    // 多子节点
    var buf bytes.Buffer
    buf.WriteString("ink.Children{\n")
    for _, elem := range elements {
        buf.WriteString("    ")
        buf.WriteString(elem)
        buf.WriteString(",\n")
    }
    buf.WriteString("    }")
    return buf.String()
}
```

### 条件渲染生成

```go
// {cond && <Element/>}
func (g *Generator) generateConditionalAnd(parsed ParsedExpr) string {
    cond := g.convertExpr(parsed.Condition)
    element := g.generateElementCode(parsed.Element)

    return fmt.Sprintf("ink.Cond(%s,\n        %s,\n        nil,\n    )",
        cond, element)
}

// {cond ? <A/> : <B/>}
func (g *Generator) generateConditionalTernary(parsed ParsedExpr) string {
    cond := g.convertExpr(parsed.Condition)
    ifTrue := g.generateElementCode(parsed.IfTrue)
    ifFalse := g.generateElementCode(parsed.IfFalse)

    return fmt.Sprintf("ink.Cond(%s,\n        %s,\n        %s,\n    )",
        cond, ifTrue, ifFalse)
}
```

### 列表渲染生成

```go
// {items.map((item, index) => <Item key={index} .../>)}
func (g *Generator) generateMapExpr(parsed ParsedExpr) string {
    items := parsed.Items
    callback := parsed.Callback

    // 解析回调参数
    itemVar := callback.Params[0]
    indexVar := ""
    if len(callback.Params) > 1 {
        indexVar = callback.Params[1]
    }

    // 生成循环体
    body := g.generateMapBody(callback.Body, itemVar, indexVar)

    // 确定 items 类型
    itemType := g.inferItemType(items)

    // 生成 ink.Map 调用
    if indexVar != "" {
        return fmt.Sprintf("ink.Map(%s, func(%s %s, %s int) ink.Element {\n%s\n    })",
            items, itemVar, itemType, indexVar, body)
    }
    return fmt.Sprintf("ink.Map(%s, func(%s %s) ink.Element {\n%s\n    })",
        items, itemVar, itemType, body)
}
```

---

## 表达式转换

```go
func (g *Generator) convertExpr(expr string) string {
    // 1. 运算符转换
    expr = strings.ReplaceAll(expr, "===", "==")
    expr = strings.ReplaceAll(expr, "!==", "!=")

    // 2. 属性访问转换
    expr = g.convertPropertyAccess(expr)

    // 3. 方法调用转换
    expr = g.convertMethodCalls(expr)

    // 4. 内置函数转换
    expr = g.convertBuiltins(expr)

    return expr
}

// 属性访问转换
func (g *Generator) convertPropertyAccess(expr string) string {
    // .length → len()
    re := regexp.MustCompile(`(\w+)\.length`)
    expr = re.ReplaceAllString(expr, "len($1)")

    // .property → .Property (首字母大写)
    re = regexp.MustCompile(`\.([a-z])([a-zA-Z]*)`)
    expr = re.ReplaceAllStringFunc(expr, func(m string) string {
        prop := m[1:]
        if prop == "length" {
            return m
        }
        return "." + strings.Title(prop)
    })

    return expr
}
```

---

## 导入管理

```go
func (c *ImportCollector) Add(path string) {
    if strings.HasPrefix(path, "github.com/") {
        c.thirdParty[path] = struct{}{}
    } else if !strings.Contains(path, ".") {
        c.standard[path] = struct{}{}
    } else {
        c.local[path] = struct{}{}
    }
}

// 生成导入块
func (c *ImportCollector) Generate() string {
    var buf bytes.Buffer
    buf.WriteString("import (\n")

    // 标准库 (按字母排序)
    stdLibs := sortedKeys(c.standard)
    for _, lib := range stdLibs {
        buf.WriteString(fmt.Sprintf("    %q\n", lib))
    }

    // 第三方库
    thirdParty := sortedKeys(c.thirdParty)
    if len(thirdParty) > 0 && len(stdLibs) > 0 {
        buf.WriteString("\n")
    }
    for _, lib := range thirdParty {
        buf.WriteString(fmt.Sprintf("    %q\n", lib))
    }

    buf.WriteString(")")
    return buf.String()
}
```

---

## 完整示例

### 输入

```go
package main

import "github.com/xxx/go-ink"

func Counter() ink.Element {
    count := 0
    show := true

    return (
        <Box flexDirection="column" padding={1}>
            <Text color="green">
                Count: {count}
            </Text>
            {show && (
                <Text color="yellow">Visible</Text>
            )}
        </Box>
    )
}
```

### 输出

```go
package main

import (
    "fmt"

    "github.com/xxx/go-ink"
)

func Counter() ink.Element {
    count := 0
    show := true

    return ink.Box(ink.Props{
        FlexDirection: ink.Column,
        Padding:      1,
    },
        ink.Children{
            ink.Text(ink.Props{
                Color: ink.Green,
            },
                ink.TextContent("Count: "),
                ink.ExprContent(fmt.Sprintf("%d", count)),
            ),
            ink.Cond(show,
                ink.Text(ink.Props{
                    Color: ink.Yellow,
                },
                    ink.TextContent("Visible"),
                ),
                nil,
            ),
        },
    )
}
```

---

## 优化: Textf 简化

```
原始输出:
ink.Text(ink.Props{Color: ink.Green},
    ink.TextContent("Count: "),
    ink.ExprContent(fmt.Sprintf("%d", count)),
)

优化后:
ink.Textf(ink.Props{Color: ink.Green},
    "Count: %d", count,
)

优化条件:
- 子节点只包含 TextContent 和 ExprContent
- 无嵌套元素
- ExprContent 是简单变量引用
```

---

## 模块结构

```
codegen/
├── generator.go       # 主生成器
├── element.go         # 元素生成
├── props.go           # 属性生成
├── children.go        # 子节点生成
├── expr.go            # 表达式转换
├── imports.go         # 导入管理
├── output.go          # 输出缓冲
├── optimize.go        # 优化处理
└── config.go          # 配置
```

---

## 总结

### 核心要点

1. **分层生成**
   - generateElement: 元素 → 组件函数调用
   - generateProps: 属性 → Props 结构体
   - generateChildren: 子节点 → Children 或单个元素
   - generateExpr: 表达式 → Go 表达式

2. **特殊处理**
   - 枚举属性: "column" → ink.Column
   - 条件渲染: && → ink.Cond
   - 列表渲染: map → ink.Map
   - 文本插值: {name} → Textf 或 ExprContent

3. **导入管理**
   - 自动收集需要的导入
   - 按规范分组 (标准库/第三方/本地)
   - 插入到正确位置

4. **优化**
   - Textf 简化
   - 格式化输出 (gofmt)
   - 去除冗余代码

5. **错误处理**
   - 位置信息保留
   - 友好的错误报告
   - 部分成功处理