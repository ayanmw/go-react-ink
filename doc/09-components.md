# 内置组件设计

> 提供终端 UI 的基础构建块

## 组件概览

```
┌─────────────────────────────────────────────────────────────┐
│                    内置组件体系                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   基础组件:                                                  │
│   ├── Box        容器组件 (Flexbox 布局)                    │
│   ├── Text       文本组件 (样式化文本)                       │
│   ├── Spacer     占位组件 (弹性空间)                        │
│   └── Newline    换行组件                                   │
│                                                             │
│   高级组件:                                                  │
│   ├── Static     静态内容 (不参与重渲染)                    │
│   ├── Transform  变换组件 (位置/样式变换)                   │
│   └── Focus      焦点组件 (焦点管理)                        │
│                                                             │
│   输入组件:                                                  │
│   └── TextInput  文本输入 (待实现)                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心接口

```go
// 组件接口
type Component interface {
    // 渲染
    Render() Element
    
    // 生命周期
    Mount()
    Unmount()
    Update(props Props)
}

// 元素接口
type Element interface {
    // 类型
    Type() string
    
    // 属性
    Props() Props
    
    // 子节点
    Children() []Element
    
    // Key (用于 Diff)
    Key() any
}

// 属性
type Props map[string]any

// 获取属性值 (带默认值)
func (p Props) Get(key string, defaultValue any) any {
    if v, ok := p[key]; ok {
        return v
    }
    return defaultValue
}
```

---

## Box 组件

```go
// Box 组件 - Flexbox 容器
func Box(props Props, children ...Element) Element {
    return &BoxElement{
        props:    props,
        children: children,
    }
}

type BoxElement struct {
    props    Props
    children []Element
}

func (e *BoxElement) Type() string {
    return "Box"
}

func (e *BoxElement) Props() Props {
    return e.props
}

func (e *BoxElement) Children() []Element {
    return e.children
}

// Box 支持的属性
type BoxProps struct {
    // Flexbox 属性
    Direction      FlexDirection   // flexDirection
    JustifyContent JustifyContent  // justifyContent
    AlignItems     AlignItems      // alignItems
    Wrap           FlexWrap        // flexWrap
    
    // Flex
    Grow   float64  // flexGrow
    Shrink float64  // flexShrink
    Basis  *float64 // flexBasis
    
    // 尺寸
    Width   *float64
    Height  *float64
    MinWidth  *float64
    MinHeight *float64
    MaxWidth  *float64
    MaxHeight *float64
    
    // 边距
    Margin  EdgeInsets
    Padding EdgeInsets
    
    // 边框
    Border       BorderStyle
    BorderColor  Color
    BorderStyle  BorderType
    
    // 背景
    Background Color
    
    // 显示
    Display DisplayMode  // flex, none
}
```

### Box 使用示例

```go
// 水平布局
<Box flexDirection="row" gap={1}>
    <Text>A</Text>
    <Text>B</Text>
</Box>

// 垂直布局
<Box flexDirection="column" padding={1}>
    <Text>Header</Text>
    <Box flexGrow={1}>
        <Text>Content</Text>
    </Box>
    <Text>Footer</Text>
</Box>

// 居中
<Box flexDirection="row" justifyContent="center" alignItems="center">
    <Text>Centered</Text>
</Box>
```

---

## Text 组件

```go
// Text 组件 - 文本显示
func Text(props Props, children ...Element) Element {
    return &TextElement{
        props:    props,
        children: children,
    }
}

// Text 支持的属性
type TextProps struct {
    // 内容
    Children string  // 或通过子节点传递
    
    // 颜色
    Color        Color  // 前景色
    Background   Color  // 背景色
    
    // 样式
    Bold          bool
    Italic        bool
    Dim           bool
    Underline     bool
    Strikethrough bool
    Inverse       bool  // 反色
    
    // 换行
    Wrap      bool     // 自动换行
    WrapMode  WrapMode // 换行模式
    
    // 对齐 (仅多行时有效)
    TextAlign TextAlign
    
    // 溢出处理
    Overflow OverflowMode  // truncate, wrap
    
    // Flex
    Grow   float64
    Shrink float64
}
```

### Text 使用示例

```go
// 基础文本
<Text>Hello World</Text>

// 样式化
<Text color="green" bold={true}>Success</Text>

// 颜色变量
<Text color={Color(255, 0, 0)}>Red Text</Text>

// 动态内容
<Text>Count: {count}</Text>

// 条件样式
<Text color={isActive ? "green" : "gray"}>
    Status
</Text>
```

---

## Spacer 组件

```go
// Spacer 组件 - 弹性占位
func Spacer(props Props) Element {
    return &SpacerElement{
        props: props,
    }
}

// Spacer 属性
type SpacerProps struct {
    // 尺寸
    Width  *float64
    Height *float64
    
    // Flex
    Grow   float64  // 默认 1
    Shrink float64
}
```

### Spacer 使用示例

```go
// 两端对齐
<Box flexDirection="row">
    <Text>Left</Text>
    <Spacer />
    <Text>Right</Text>
</Box>

// 垂直间隔
<Box flexDirection="column">
    <Text>Top</Text>
    <Spacer height={2} />
    <Text>Bottom</Text>
</Box>
```

---

## Newline 组件

```go
// Newline 组件 - 换行
func Newline(props Props) Element {
    return &NewlineElement{
        props: props,
    }
}

// Newline 属性
type NewlineProps struct {
    Count int  // 换行数，默认 1
}
```

### Newline 使用示例

```go
<Text>
    Line 1
    <Newline />
    Line 2
    <Newline count={2} />
    Line 4
</Text>
```

---

## Static 组件

```go
// Static 组件 - 静态内容 (不参与重渲染)
func Static(props Props, children ...Element) Element {
    return &StaticElement{
        props:    props,
        children: children,
    }
}

// Static 属性
type StaticProps struct {
    // 标识 (用于更新)
    Key string
    
    // 子节点
    Children []Element
}
```

### Static 使用示例

```go
// 日志输出 (只追加，不重绘)
<Box flexDirection="column">
    <Static>
        {logs.map(log => (
            <Text key={log.id}>{log.message}</Text>
        ))}
    </Static>
    <Text>Input: {input}</Text>
</Box>
```

---

## Transform 组件

```go
// Transform 组件 - 变换
func Transform(props Props, children ...Element) Element {
    return &TransformElement{
        props:    props,
        children: children,
    }
}

// Transform 属性
type TransformProps struct {
    // 位置变换
    OffsetX int
    OffsetY int
    
    // 样式变换
    StyleFunc func(CellStyle) CellStyle
    
    // 内容变换
    TransformFunc func(string) string
}
```

### Transform 使用示例

```go
// 偏移
<Transform offsetX={2} offsetY={1}>
    <Text>Shifted</Text>
</Transform>

// 样式变换 (高亮)
<Transform styleFunc={s => { s.Inverse = true; return s }}>
    <Text>Highlighted</Text>
</Transform>

// 内容变换 (大写)
<Transform transformFunc={strings.ToUpper}>
    <Text>{text}</Text>
</Transform>
```

---

## 辅助函数

```go
// Children - 子节点列表
func Children(elements ...Element) Element {
    return &FragmentElement{
        children: elements,
    }
}

// Cond - 条件渲染
func Cond(condition bool, ifTrue, ifFalse Element) Element {
    if condition {
        return ifTrue
    }
    return ifFalse
}

// Map - 列表渲染
func Map[T any](items []T, render func(item T, index int) Element) Element {
    children := make([]Element, len(items))
    for i, item := range items {
        children[i] = render(item, i)
    }
    return &FragmentElement{
        children: children,
    }
}

// TextContent - 纯文本内容
func TextContent(text string) Element {
    return &TextElement{
        props: Props{},
        children: []Element{&TextContentNode{text: text}},
    }
}

// ExprContent - 表达式内容
func ExprContent(format string, args ...any) Element {
    text := fmt.Sprintf(format, args...)
    return TextContent(text)
}
```

---

## 颜色系统

```go
// 预定义颜色
var (
    // 基础色
    Black   = Color{Type: ColorBasic, Basic: 0}
    Red     = Color{Type: ColorBasic, Basic: 1}
    Green   = Color{Type: ColorBasic, Basic: 2}
    Yellow  = Color{Type: ColorBasic, Basic: 3}
    Blue    = Color{Type: ColorBasic, Basic: 4}
    Magenta = Color{Type: ColorBasic, Basic: 5}
    Cyan    = Color{Type: ColorBasic, Basic: 6}
    White   = Color{Type: ColorBasic, Basic: 7}
    
    // 亮色
    BrightBlack   = Color{Type: ColorBasic, Basic: 8}   // Gray
    BrightRed     = Color{Type: ColorBasic, Basic: 9}
    BrightGreen   = Color{Type: ColorBasic, Basic: 10}
    BrightYellow  = Color{Type: ColorBasic, Basic: 11}
    BrightBlue    = Color{Type: ColorBasic, Basic: 12}
    BrightMagenta = Color{Type: ColorBasic, Basic: 13}
    BrightCyan    = Color{Type: ColorBasic, Basic: 14}
    BrightWhite   = Color{Type: ColorBasic, Basic: 15}
    
    // 别名
    Gray = BrightBlack
)

// 256 色
func Color256(n int) Color {
    return Color{Type: ColorExtended, Extended: n}
}

// RGB 真彩色
func RGB(r, g, b uint8) Color {
    return Color{Type: ColorTrueColor, R: r, G: g, B: b}
}

// Hex 颜色
func Hex(hex string) Color {
    // 解析 #RRGGBB
    r, _ := strconv.ParseUint(hex[1:3], 16, 8)
    g, _ := strconv.ParseUint(hex[3:5], 16, 8)
    b, _ := strconv.ParseUint(hex[5:7], 16, 8)
    return RGB(uint8(r), uint8(g), uint8(b))
}
```

---

## 边框样式

```go
// 边框类型
type BorderType int

const (
    BorderNone BorderType = iota
    BorderSingle
    BorderDouble
    BorderRound
    BorderBold
)

// 边框字符
var borderChars = map[BorderType]struct {
    TopLeft, TopRight, BottomLeft, BottomRight rune
    Horizontal, Vertical                        rune
}{
    BorderSingle: {'┌', '┐', '└', '┘', '─', '│'},
    BorderDouble: {'╔', '╗', '╚', '╝', '═', '║'},
    BorderRound:  {'╭', '╮', '╰', '╯', '─', '│'},
    BorderBold:   {'┏', '┓', '┗', '┛', '━', '┃'},
}
```

---

## 模块结构

```
components/
├── component.go      # 组件接口
├── element.go        # 元素接口
├── props.go          # 属性定义
│
├── box/
│   ├── box.go        # Box 组件
│   └── props.go      # Box 属性
│
├── text/
│   ├── text.go       # Text 组件
│   └── props.go      # Text 属性
│
├── spacer/
│   └── spacer.go     # Spacer 组件
│
├── newline/
│   └── newline.go    # Newline 组件
│
├── static/
│   └── static.go     # Static 组件
│
├── transform/
│   └── transform.go  # Transform 组件
│
├── helpers/
│   ├── children.go   # Children 辅助
│   ├── cond.go       # Cond 辅助
│   ├── map.go        # Map 辅助
│   └── content.go    # TextContent/ExprContent
│
└── styles/
    ├── colors.go     # 颜色定义
    ├── border.go     # 边框样式
    └── edge.go       # 边距定义
```

---

## 总结

### 核心要点

1. **基础组件**
   - Box: Flexbox 容器
   - Text: 样式化文本
   - Spacer: 弹性占位
   - Newline: 换行控制

2. **高级组件**
   - Static: 静态内容优化
   - Transform: 灵活变换

3. **辅助函数**
   - Children: 子节点列表
   - Cond: 条件渲染
   - Map: 列表渲染
   - TextContent/ExprContent: 内容节点

4. **样式系统**
   - 颜色: 基础色/256色/真彩色
   - 边框: 多种样式
   - 文本: 粗体/斜体/下划线等

5. **设计原则**
   - 与 Ink API 一致
   - 类型安全
   - 可扩展
