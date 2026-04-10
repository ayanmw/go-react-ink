# components Specification

## Purpose
定义 Go-Ink 内置组件的设计和实现，包括 Box、Text、Static、Newline、Spacer、Transform 等。

## Requirements

### Requirement: Box 组件实现 Flexbox 容器
Box SHALL 作为 Flexbox 容器，支持布局属性。

#### Scenario: flexDirection
- **WHEN** 用户设置 `<Box flexDirection="column">`
- **THEN** 子元素垂直排列

#### Scenario: justifyContent
- **WHEN** 用户设置 `<Box justifyContent="center">`
- **THEN** 子元素居中对齐

#### Scenario: alignItems
- **WHEN** 用户设置 `<Box alignItems="stretch">`
- **THEN** 子元素拉伸填充交叉轴

### Requirement: Text 组件显示文本
Text SHALL 显示文本内容，支持样式属性。

#### Scenario: color
- **WHEN** 用户设置 `<Text color="green">`
- **THEN** 文本显示为绿色

#### Scenario: bold
- **WHEN** 用户设置 `<Text bold>`
- **THEN** 文本加粗显示

#### Scenario: wrap
- **WHEN** 用户设置 `<Text wrap="truncate">`
- **THEN** 文本超出时截断

### Requirement: Static 组件保持输出
Static SHALL 保持内容不被后续渲染清除。

#### Scenario: 静态输出
- **WHEN** 用户使用 `<Static items={items}>`
- **THEN** 内容保持固定，不参与增量渲染

### Requirement: Newline 组件换行
Newline SHALL 输出指定数量的换行。

#### Scenario: 换行
- **WHEN** 用户使用 `<Newline count={2}>`
- **THEN** 输出 2 个换行符

### Requirement: Spacer 组件填充空间
Spacer SHALL 填充剩余空间。

#### Scenario: 填充
- **WHEN** 用户使用 `<Spacer>`
- **THEN** 组件填充父容器剩余空间

### Requirement: Transform 组件变换输出
Transform SHALL 对子元素输出应用变换函数。

#### Scenario: 变换
- **WHEN** 用户使用 `<Transform transform={fn}>`
- **THEN** 子元素输出经过 fn 变换

## Box Component

```go
// Box 组件
func Box(props Props, children ...Element) Element {
    return &BoxElement{
        Props:    props,
        Children: children,
    }
}

// BoxElement
type BoxElement struct {
    Props    Props
    Children []Element
}

// 布局属性
type BoxProps struct {
    // Flexbox
    FlexDirection  string  // row, column
    JustifyContent string  // flex-start, center, flex-end, space-between, space-around
    AlignItems     string  // stretch, flex-start, center, flex-end, baseline
    FlexWrap       string  // nowrap, wrap
    Gap            int
    
    // 尺寸
    Width      any  // int 或 "50%"
    Height     any
    MinWidth   int
    MaxWidth   int
    MinHeight  int
    MaxHeight  int
    
    // Flex
    FlexGrow   float64
    FlexShrink float64
    FlexBasis  any
    
    // 间距
    Padding    any  // int 或 EdgeInsets
    Margin     any
    
    // 边框
    BorderStyle string  // single, double, round, bold
    BorderColor string
}
```

## Text Component

```go
// Text 组件
func Text(props Props, children ...Element) Element {
    return &TextElement{
        Props:    props,
        Children: children,
    }
}

// TextElement
type TextElement struct {
    Props    Props
    Children []Element
}

// 文本属性
type TextProps struct {
    // 颜色
    Color           string  // green, red, blue, yellow, magenta, cyan, white, black, gray
    BackgroundColor string
    
    // 样式
    Bold         bool
    Italic       bool
    Underline    bool
    Strikethrough bool
    Dim          bool
    Inverse      bool
    
    // 换行
    Wrap string  // wrap, truncate, truncate-start, truncate-middle, truncate-end
}
```

## Static Component

```go
// Static 组件
func Static(props Props, children ...Element) Element {
    return &StaticElement{
        Props:    props,
        Children: children,
    }
}

// StaticElement
type StaticElement struct {
    Props    Props
    Children []Element
}

// Static 属性
type StaticProps struct {
    Items []any
}
```

## Newline Component

```go
// Newline 组件
func Newline(props Props) Element {
    count := props.GetInt("count", 1)
    return &NewlineElement{Count: count}
}

// NewlineElement
type NewlineElement struct {
    Count int
}
```

## Spacer Component

```go
// Spacer 组件
func Spacer(props Props) Element {
    return &SpacerElement{Props: props}
}

// SpacerElement
type SpacerElement struct {
    Props Props
}
```

## Transform Component

```go
// Transform 组件
func Transform(props Props, children ...Element) Element {
    return &TransformElement{
        Props:    props,
        Children: children,
    }
}

// TransformElement
type TransformElement struct {
    Props    Props
    Children []Element
}

// Transform 属性
type TransformProps struct {
    Transform func(string) string
}
```

## Module Structure

```
components/
├── box.go          # Box 组件
├── text.go         # Text 组件
├── static.go       # Static 组件
├── newline.go      # Newline 组件
├── spacer.go       # Spacer 组件
├── transform.go    # Transform 组件
└── props.go        # Props 类型
```

## React Ink 对齐

| 组件 | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| Box | ✅ | ✅ | Flexbox 容器 |
| Text | ✅ | ✅ | 文本显示 |
| Static | ✅ | ✅ | 静态输出 |
| Newline | ✅ | ✅ | 换行 |
| Spacer | ✅ | ✅ | 填充空间 |
| Transform | ✅ | ✅ | 输出变换 |

### Box 属性对齐

| 属性 | React Ink | Go-Ink |
|-----|-----------|--------|
| flexDirection | ✅ | ✅ |
| justifyContent | ✅ | ✅ |
| alignItems | ✅ | ✅ |
| flexWrap | ✅ | ✅ |
| gap | ✅ | ✅ |
| width/height | ✅ | ✅ |
| minWidth/maxWidth | ✅ | ✅ |
| flexGrow/flexShrink | ✅ | ✅ |
| padding/margin | ✅ | ✅ |
| borderStyle | ✅ | ✅ |

### Text 属性对齐

| 属性 | React Ink | Go-Ink |
|-----|-----------|--------|
| color | ✅ | ✅ |
| backgroundColor | ✅ | ✅ |
| bold | ✅ | ✅ |
| italic | ✅ | ✅ |
| underline | ✅ | ✅ |
| strikethrough | ✅ | ✅ |
| dimColor | ✅ | ✅ |
| inverse | ✅ | ✅ |
| wrap | ✅ | ✅ |