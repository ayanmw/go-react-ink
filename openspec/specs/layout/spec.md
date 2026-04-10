# layout Specification

## Purpose
定义 Go-Ink Flexbox 布局引擎的设计和实现，计算组件的位置和尺寸。

## Requirements

### Requirement: 布局引擎实现 Flexbox
布局引擎 SHALL 支持 CSS Flexbox 布局算法的核心子集。

#### Scenario: flexDirection
- **WHEN** 用户设置 flexDirection="row"
- **THEN** 子元素水平排列

#### Scenario: justifyContent
- **WHEN** 用户设置 justifyContent="center"
- **THEN** 子元素居中对齐

#### Scenario: alignItems
- **WHEN** 用户设置 alignItems="stretch"
- **THEN** 子元素拉伸填充交叉轴

#### Scenario: flex
- **WHEN** 用户设置 flexGrow=1
- **THEN** 子元素填充剩余空间

### Requirement: 布局引擎计算尺寸
布局引擎 SHALL 根据属性计算元素的最终尺寸。

#### Scenario: 固定尺寸
- **WHEN** 用户设置 width=100
- **THEN** 元素宽度为 100 列

#### Scenario: 百分比尺寸
- **WHEN** 用户设置 width="50%"
- **THEN** 元素宽度为父元素的 50%

#### Scenario: 自动尺寸
- **WHEN** 用户未设置尺寸
- **THEN** 元素根据内容自适应

### Requirement: 布局引擎计算位置
布局引擎 SHALL 根据属性计算元素的最终位置。

#### Scenario: padding
- **WHEN** 用户设置 padding=2
- **THEN** 元素内容与边框间距 2 列

#### Scenario: margin
- **WHEN** 用户设置 margin=1
- **THEN** 元素与其他元素间距 1 列

### Requirement: 布局引擎支持 gap
布局引擎 SHALL 支持 gap 属性控制子元素间距。

#### Scenario: gap
- **WHEN** 用户设置 gap=1
- **THEN** 子元素之间间距 1 列

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Layout 架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: 元素树 (来自 Reconciler)                              │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │  Element    │  元素树                                   │
│   │  Tree       │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Create     │  创建布局节点                              │
│   │  Nodes      │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Measure    │  测量文本尺寸                              │
│   │  Text       │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Calculate  │  计算 Flexbox 布局                         │
│   │  Layout     │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Output     │  布局结果                                  │
│   │  Nodes      │                                          │
│   └─────────────┘                                          │
│                                                             │
│   输出: 布局树 (每个节点有 x, y, width, height)               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Layout Node

```go
// 布局节点
type LayoutNode struct {
    // 属性
    Style Style
    
    // 计算结果
    X      float64
    Y      float64
    Width  float64
    Height float64
    
    // 子节点
    Children []*LayoutNode
    
    // 文本内容
    Text string
}

// 样式属性
type Style struct {
    // Flexbox
    FlexDirection  FlexDirection
    JustifyContent JustifyContent
    AlignItems     AlignItems
    FlexWrap       FlexWrap
    Gap            float64
    
    // 尺寸
    Width      Dimension
    Height     Dimension
    MinWidth   Dimension
    MaxWidth   Dimension
    MinHeight  Dimension
    MaxHeight  Dimension
    
    // Flex
    FlexGrow   float64
    FlexShrink float64
    FlexBasis  Dimension
    
    // 间距
    Padding   EdgeInsets
    Margin    EdgeInsets
}
```

## Flexbox Algorithm

```go
// 计算布局
func (n *LayoutNode) CalculateLayout() {
    // 1. 确定可用空间
    availableWidth := n.Width
    availableHeight := n.Height
    
    // 2. 测量子节点
    n.measureChildren()
    
    // 3. 计算主轴尺寸
    n.calculateMainAxis()
    
    // 4. 计算交叉轴尺寸
    n.calculateCrossAxis()
    
    // 5. 定位子节点
    n.positionChildren()
}

// 计算主轴
func (n *LayoutNode) calculateMainAxis() {
    flexDirection := n.Style.FlexDirection
    
    // 收集 flex 子节点
    var flexChildren []*LayoutNode
    var totalFlexGrow float64
    var fixedSize float64
    
    for _, child := range n.Children {
        if child.Style.FlexGrow > 0 {
            flexChildren = append(flexChildren, child)
            totalFlexGrow += child.Style.FlexGrow
        } else {
            fixedSize += child.getSize(flexDirection)
        }
    }
    
    // 分配剩余空间给 flex 子节点
    remainingSpace := n.getInnerSize(flexDirection) - fixedSize
    
    for _, child := range flexChildren {
        child.setSize(flexDirection, 
            remainingSpace * child.Style.FlexGrow / totalFlexGrow)
    }
}

// 定位子节点
func (n *LayoutNode) positionChildren() {
    x := n.Style.Padding.Left
    y := n.Style.Padding.Top
    
    for _, child := range n.Children {
        // 应用 justifyContent
        switch n.Style.JustifyContent {
        case JustifyCenter:
            child.X = x + (n.getInnerWidth() - child.Width) / 2
        case JustifyFlexEnd:
            child.X = x + n.getInnerWidth() - child.Width
        default:
            child.X = x
        }
        
        // 应用 alignItems
        switch n.Style.AlignItems {
        case AlignCenter:
            child.Y = y + (n.getInnerHeight() - child.Height) / 2
        case AlignFlexEnd:
            child.Y = y + n.getInnerHeight() - child.Height
        default:
            child.Y = y
        }
        
        // 更新下一个子节点位置
        if n.Style.FlexDirection == FlexRow {
            x += child.Width + n.Style.Gap
        } else {
            y += child.Height + n.Style.Gap
        }
    }
}
```

## Supported Properties

### Flex Direction
- `row` - 水平排列
- `column` - 垂直排列
- `row-reverse` - 反向水平排列
- `column-reverse` - 反向垂直排列

### Justify Content
- `flex-start` - 起点
- `center` - 居中
- `flex-end` - 终点
- `space-between` - 两端对齐
- `space-around` - 环绕

### Align Items
- `stretch` - 拉伸
- `flex-start` - 起点
- `center` - 居中
- `flex-end` - 终点
- `baseline` - 基线

### Flex
- `flexGrow` - 增长因子
- `flexShrink` - 收缩因子
- `flexBasis` - 初始尺寸

### Dimensions
- `width` / `height` - 固定值或百分比
- `minWidth` / `minHeight` - 最小尺寸
- `maxWidth` / `maxHeight` - 最大尺寸

### Spacing
- `padding` - 内边距
- `margin` - 外边距
- `gap` - 子元素间距

## Module Structure

```
layout/
├── node.go             # 布局节点
├── style.go            # 样式定义
├── calculate.go        # 布局计算
├── measure.go          # 文本测量
├── flex.go             # Flexbox 算法
└── types.go            # 类型定义
```

## React Ink 对齐

| 特性 | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| Flexbox | ✅ Yoga | ✅ | 简化实现 |
| flexDirection | ✅ | ✅ | row/column |
| justifyContent | ✅ | ✅ | 5 种对齐 |
| alignItems | ✅ | ✅ | 5 种对齐 |
| flexGrow | ✅ | ✅ | 增长因子 |
| flexShrink | ✅ | ✅ | 收缩因子 |
| gap | ✅ | ✅ | 子元素间距 |
| padding/margin | ✅ | ✅ | EdgeInsets |
| 百分比尺寸 | ✅ | ✅ | "50%" |
| 文本测量 | ✅ | ✅ | 自动换行 |