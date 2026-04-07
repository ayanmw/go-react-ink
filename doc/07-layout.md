# Flexbox 布局引擎设计

> 实现 CSS Flexbox 核心算法，为终端 UI 提供布局能力

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  Flexbox Layout 架构                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: 元素树 (来自 Reconciler)                              │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │  Element    │                                          │
│   │  Tree       │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Create     │  创建布局节点                             │
│   │  Nodes      │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Calculate  │  计算布局                                 │
│   │  Layout     │  ├─ Measure Text                         │
│   │             │  ├─ Resolve Flex                         │
│   │             │  └─ Position Children                    │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Output     │  布局结果                                 │
│   │  (x, y,     │                                          │
│   │   w, h)     │                                          │
│   └─────────────┘                                          │
│                                                             │
│   输出: 每个元素的绝对位置和尺寸                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 布局节点结构

```go
// 布局节点
type LayoutNode struct {
    // 树结构
    Parent   *LayoutNode
    Children []*LayoutNode
    
    // Flex 属性
    Style FlexStyle
    
    // 计算结果
    Layout LayoutResult
    
    // 测量
    MeasureFunc MeasureFunc  // 自定义测量函数
    
    // 关联
    Element any  // 关联的 Fiber 或组件实例
}

// Flex 样式
type FlexStyle struct {
    // 方向
    Direction FlexDirection  // row, column
    
    // 主轴对齐
    JustifyContent JustifyContent  // flex-start, center, space-between...
    
    // 交叉轴对齐
    AlignItems     AlignItems      // stretch, center, flex-start...
    AlignSelf      AlignItems      // 覆盖父级 AlignItems
    
    // 换行
    FlexWrap FlexWrap  // nowrap, wrap
    
    // 尺寸
    Width   *float64  // 固定宽度
    Height  *float64  // 固定高度
    
    MinWidth  *float64
    MinHeight *float64
    MaxWidth  *float64
    MaxHeight *float64
    
    // Flex 属性
    FlexGrow   float64  // 默认 0
    FlexShrink float64  // 默认 1
    FlexBasis  *float64  // 默认 auto
    
    // 边距
    Margin  EdgeInsets
    Padding EdgeInsets
    
    // 边框
    Border EdgeInsets
}

// 布局结果
type LayoutResult struct {
    // 位置 (相对于父节点)
    X float64
    Y float64
    
    // 尺寸
    Width  float64
    Height float64
    
    // 是否计算完成
    Computed bool
}

// 边距
type EdgeInsets struct {
    Top    float64
    Right  float64
    Bottom float64
    Left   float64
}
```

---

## 枚举定义

```go
// Flex 方向
type FlexDirection int

const (
    DirectionRow FlexDirection = iota
    DirectionColumn
    DirectionRowReverse
    DirectionColumnReverse
)

// 主轴对齐
type JustifyContent int

const (
    JustifyFlexStart JustifyContent = iota
    JustifyFlexEnd
    JustifyCenter
    JustifySpaceBetween
    JustifySpaceAround
    JustifySpaceEvenly
)

// 交叉轴对齐
type AlignItems int

const (
    AlignStretch AlignItems = iota
    AlignFlexStart
    AlignFlexEnd
    AlignCenter
    AlignBaseline
)

// 换行
type FlexWrap int

const (
    WrapNoWrap FlexWrap = iota
    WrapWrap
    WrapWrapReverse
)
```

---

## 布局计算核心

```go
// 布局计算器
type LayoutCalculator struct {
    // 配置
    config Config
}

// 计算布局
func (c *LayoutCalculator) CalculateLayout(root *LayoutNode, 
    availableWidth, availableHeight float64) {
    
    // 1. 测量阶段
    c.measureNode(root, availableWidth, availableHeight)
    
    // 2. 布局阶段
    c.layoutNode(root, 0, 0, availableWidth, availableHeight)
}

// 测量节点 (确定尺寸)
func (c *LayoutCalculator) measureNode(node *LayoutNode,
    availableWidth, availableHeight float64) {
    
    // 1. 确定约束
    widthConstraint := c.resolveWidthConstraint(node, availableWidth)
    heightConstraint := c.resolveHeightConstraint(node, availableHeight)
    
    // 2. 根据节点类型测量
    if node.MeasureFunc != nil {
        // 自定义测量 (文本节点)
        size := node.MeasureFunc(widthConstraint, heightConstraint)
        node.Layout.Width = size.Width
        node.Layout.Height = size.Height
    } else if len(node.Children) == 0 {
        // 叶子节点
        node.Layout.Width = c.resolveDimension(node.Style.Width, widthConstraint)
        node.Layout.Height = c.resolveDimension(node.Style.Height, heightConstraint)
    } else {
        // 容器节点: 测量子节点后确定
        c.measureContainer(node, widthConstraint, heightConstraint)
    }
    
    // 3. 应用 min/max 约束
    node.Layout.Width = c.clamp(node.Layout.Width, node.Style.MinWidth, node.Style.MaxWidth)
    node.Layout.Height = c.clamp(node.Layout.Height, node.Style.MinHeight, node.Style.MaxHeight)
    
    node.Layout.Computed = true
}

// 测量容器
func (c *LayoutCalculator) measureContainer(node *LayoutNode,
    widthConstraint, heightConstraint float64) {
    
    direction := node.Style.Direction
    isRow := direction == DirectionRow || direction == DirectionRowReverse
    
    // 1. 计算内容可用空间
    contentWidth := widthConstraint - node.Style.Padding.Horizontal() - node.Style.Border.Horizontal()
    contentHeight := heightConstraint - node.Style.Padding.Vertical() - node.Style.Border.Vertical()
    
    // 2. 测量所有子节点
    for _, child := range node.Children {
        if isRow {
            c.measureNode(child, contentWidth, contentHeight)
        } else {
            c.measureNode(child, contentWidth, contentHeight)
        }
    }
    
    // 3. 解析 Flex 属性
    c.resolveFlex(node, isRow, contentWidth, contentHeight)
    
    // 4. 计算容器尺寸
    if node.Style.Width == nil {
        node.Layout.Width = c.calculateContainerWidth(node, isRow)
    } else {
        node.Layout.Width = *node.Style.Width
    }
    
    if node.Style.Height == nil {
        node.Layout.Height = c.calculateContainerHeight(node, isRow)
    } else {
        node.Layout.Height = *node.Style.Height
    }
}
```

---

## Flex 解析

```go
// 解析 Flex 属性
func (c *LayoutCalculator) resolveFlex(node *LayoutNode, isRow bool,
    availableWidth, availableHeight float64) {
    
    // 主轴和交叉轴可用空间
    mainAxisAvailable := availableWidth
    crossAxisAvailable := availableHeight
    if !isRow {
        mainAxisAvailable, crossAxisAvailable = crossAxisAvailable, mainAxisAvailable
    }
    
    // 1. 计算固定尺寸子节点占用空间
    var totalFixedMain float64
    var totalFlexGrow float64
    var totalFlexShrink float64
    
    for _, child := range node.Children {
        childMain := child.Layout.Width
        if !isRow {
            childMain = child.Layout.Height
        }
        
        if child.Style.FlexGrow > 0 {
            totalFlexGrow += child.Style.FlexGrow
        } else if childMain > 0 {
            totalFixedMain += childMain + child.Style.Margin.Horizontal()
            if !isRow {
                totalFixedMain = childMain + child.Style.Margin.Vertical()
            }
        }
        
        totalFlexShrink += child.Style.FlexShrink
    }
    
    // 2. 计算剩余空间
    remainingMain := mainAxisAvailable - totalFixedMain
    
    // 3. 分配 Flex 空间
    if totalFlexGrow > 0 && remainingMain > 0 {
        // 有剩余空间，按 flexGrow 分配
        for _, child := range node.Children {
            if child.Style.FlexGrow > 0 {
                flexMain := remainingMain * (child.Style.FlexGrow / totalFlexGrow)
                if isRow {
                    child.Layout.Width = flexMain
                } else {
                    child.Layout.Height = flexMain
                }
            }
        }
    } else if totalFlexShrink > 0 && remainingMain < 0 {
        // 空间不足，按 flexShrink 收缩
        for _, child := range node.Children {
            if child.Style.FlexShrink > 0 {
                childMain := child.Layout.Width
                if !isRow {
                    childMain = child.Layout.Height
                }
                shrinkFactor := childMain * child.Style.FlexShrink / totalFlexShrink
                shrinkAmount := -remainingMain * shrinkFactor
                if isRow {
                    child.Layout.Width = childMain - shrinkAmount
                } else {
                    child.Layout.Height = childMain - shrinkAmount
                }
            }
        }
    }
}
```

---

## 子节点定位

```go
// 布局节点 (确定位置)
func (c *LayoutCalculator) layoutNode(node *LayoutNode,
    x, y, availableWidth, availableHeight float64) {
    
    // 设置位置
    node.Layout.X = x + node.Style.Margin.Left + node.Style.Border.Left
    node.Layout.Y = y + node.Style.Margin.Top + node.Style.Border.Top
    
    // 布局子节点
    if len(node.Children) > 0 {
        c.layoutChildren(node)
    }
}

// 布局子节点
func (c *LayoutCalculator) layoutChildren(node *LayoutNode) {
    direction := node.Style.Direction
    isRow := direction == DirectionRow || direction == DirectionRowReverse
    isReverse := direction == DirectionRowReverse || direction == DirectionColumnReverse
    
    // 内容区域起点
    contentX := node.Style.Padding.Left + node.Style.Border.Left
    contentY := node.Style.Padding.Top + node.Style.Border.Top
    contentWidth := node.Layout.Width - node.Style.Padding.Horizontal() - node.Style.Border.Horizontal()
    contentHeight := node.Layout.Height - node.Style.Padding.Vertical() - node.Style.Border.Vertical()
    
    // 计算主轴总尺寸
    mainSize := c.calculateMainSize(node, isRow)
    
    // 计算起始位置 (根据 justifyContent)
    mainOffset := c.calculateJustifyOffset(node.Style.JustifyContent, 
        isRow ? contentWidth : contentHeight, mainSize, len(node.Children))
    
    // 定位每个子节点
    currentMain := mainOffset
    maxCross := 0.0
    
    children := node.Children
    if isReverse {
        children = reverseSlice(children)
    }
    
    for _, child := range children {
        // 交叉轴对齐
        crossOffset := c.calculateAlignOffset(node.Style.AlignItems, child.Style.AlignSelf,
            isRow ? contentHeight : contentWidth,
            isRow ? child.Layout.Height : child.Layout.Width)
        
        // 设置位置
        childX := contentX
        childY := contentY
        
        if isRow {
            childX += currentMain
            childY += crossOffset
        } else {
            childX += crossOffset
            childY += currentMain
        }
        
        // 递归布局
        c.layoutNode(child, childX, childY, child.Layout.Width, child.Layout.Height)
        
        // 更新主轴位置
        childMain := child.Layout.Width
        if !isRow {
            childMain = child.Layout.Height
        }
        currentMain += childMain + child.Style.Margin.Horizontal()
        if !isRow {
            currentMain += childMain + child.Style.Margin.Vertical()
        }
        
        // 更新最大交叉轴尺寸
        childCross := child.Layout.Height
        if isRow {
            childCross = child.Layout.Height
        } else {
            childCross = child.Layout.Width
        }
        if childCross > maxCross {
            maxCross = childCross
        }
    }
}
```

---

## 对齐计算

```go
// 计算主轴对齐偏移
func (c *LayoutCalculator) calculateJustifyOffset(justify JustifyContent,
    containerMain, childrenMain float64, childCount int) float64 {
    
    freeSpace := containerMain - childrenMain
    
    switch justify {
    case JustifyFlexStart:
        return 0
    case JustifyFlexEnd:
        return freeSpace
    case JustifyCenter:
        return freeSpace / 2
    case JustifySpaceBetween:
        if childCount <= 1 {
            return 0
        }
        return 0  // 第一个节点在起点
    case JustifySpaceAround:
        if childCount == 0 {
            return 0
        }
        return freeSpace / (2 * float64(childCount))
    case JustifySpaceEvenly:
        if childCount == 0 {
            return 0
        }
        return freeSpace / float64(childCount+1)
    }
    return 0
}

// 计算交叉轴对齐偏移
func (c *LayoutCalculator) calculateAlignOffset(alignItems, alignSelf AlignItems,
    containerCross, childCross float64) float64 {
    
    align := alignItems
    if alignSelf != AlignStretch {
        align = alignSelf
    }
    
    switch align {
    case AlignFlexStart:
        return 0
    case AlignFlexEnd:
        return containerCross - childCross
    case AlignCenter:
        return (containerCross - childCross) / 2
    case AlignStretch:
        // Stretch 在测量阶段处理
        return 0
    case AlignBaseline:
        // Baseline 需要额外计算
        return 0
    }
    return 0
}
```

---

## 文本测量

```go
// 文本测量函数
type MeasureFunc func(availableWidth, availableHeight float64) Size

// 创建文本测量函数
func CreateTextMeasureFunc(text string, style TextStyle) MeasureFunc {
    return func(availableWidth, availableHeight float64) Size {
        // 1. 计算文本宽度
        width := measureTextWidth(text, style)
        
        // 2. 处理换行
        if style.Wrap && width > availableWidth {
            lines := wrapText(text, style, availableWidth)
            width = availableWidth
            height := float64(len(lines))
            return Size{Width: width, Height: height}
        }
        
        // 3. 单行文本
        return Size{Width: width, Height: 1}
    }
}

// 测量文本宽度 (考虑 ANSI 和宽字符)
func measureTextWidth(text string, style TextStyle) float64 {
    width := 0.0
    
    // 去除 ANSI 转义序列
    cleanText := stripAnsi(text)
    
    for _, r := range cleanText {
        if r < 128 {
            // ASCII 字符
            width += 1
        } else {
            // 宽字符 (中文等)
            width += 2
        }
    }
    
    return width
}
```

---

## 支持的属性

### 完整支持

| 属性 | 说明 |
|-----|------|
| flexDirection | row, column |
| justifyContent | flex-start, center, flex-end, space-between |
| alignItems | stretch, flex-start, center, flex-end |
| flexWrap | nowrap, wrap |
| flexGrow | 弹性增长 |
| flexShrink | 弹性收缩 |
| width | 固定宽度 |
| height | 固定高度 |
| padding | 内边距 |
| margin | 外边距 |

### 部分支持

| 属性 | 说明 |
|-----|------|
| flexBasis | 仅支持 auto |
| minWidth/maxWidth | 最小/最大宽度约束 |
| minHeight/maxHeight | 最小/最大高度约束 |
| alignSelf | 覆盖父级 alignItems |

### 不支持 (终端场景不需要)

| 属性 | 原因 |
|-----|------|
| gap | 可用 margin 替代 |
| order | 终端场景极少需要 |
| alignContent | 单行场景不需要 |
| aspectRatio | 终端字符比例固定 |

---

## 模块结构

```
layout/
├── layout.go          # 布局计算入口
├── node.go            # 布局节点定义
├── style.go           # Flex 样式定义
├── enums.go           # 枚举定义
├── measure.go         # 测量阶段
├── resolve.go         # Flex 解析
├── position.go        # 子节点定位
├── align.go           # 对齐计算
├── text.go            # 文本测量
├── constraints.go     # 尺寸约束
└── cache.go           # 布局缓存
```

---

## 实现决策

### 实现路线 (倒序渐进式验证)

```
┌─────────────────────────────────────────────────────────────┐
│              Flexbox 实现路线 (倒序验证)                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   Phase 1: CGO + yoga-layout (最快验证)                     │
│   ├── 目标: 快速跑通整个流程                                 │
│   ├── 优点: 1-2 天即可集成                                   │
│   ├── 缺点: CGO 依赖，跨平台编译复杂                         │
│   └── 验证: 所有测试用例通过                                 │
│                                                             │
│   Phase 2: Go 核心子集实现 (主力实现)                        │
│   ├── 目标: 替换 CGO，纯 Go 实现                             │
│   ├── 范围: 终端场景 80% 常用属性                            │
│   │   ├── flexDirection (row/column)                        │
│   │   ├── justifyContent (5种对齐)                          │
│   │   ├── alignItems (5种对齐)                              │
│   │   ├── flexWrap                                          │
│   │   ├── flexGrow / flexShrink                             │
│   │   ├── width / height                                    │
│   │   └── padding / margin                                  │
│   ├── 优点: 纯 Go，无 CGO 依赖                               │
│   └── 验证: 测试用例通过，性能可接受                         │
│                                                             │
│   Phase 3: 完整实现 (按需扩展)                               │
│   ├── 目标: 100% 属性支持                                    │
│   ├── 范围: 剩余 20% 属性                                    │
│   │   ├── alignContent                                      │
│   │   ├── gap                                               │
│   │   ├── order                                             │
│   │   └── baseline 对齐                                     │
│   └── 条件: Phase 2 测试全部通过后                           │
│                                                             │
│   回退机制:                                                  │
│   └── 如果 Phase 2 某属性实现困难，可回退到 CGO              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 为什么倒序？

1. **降低风险**: 先验证整体架构可行，再优化实现
2. **快速迭代**: CGO 方案 1-2 天即可跑通
3. **明确目标**: 测试用例通过后，知道 Go 实现需要达到的效果
4. **渐进优化**: 从依赖 → 自主实现，逐步降低复杂度

---

## 总结

### 核心要点

1. **两阶段算法**
   - 测量阶段: 确定节点尺寸
   - 布局阶段: 确定节点位置

2. **Flex 核心**
   - flexGrow: 剩余空间分配
   - flexShrink: 空间不足收缩
   - 主轴/交叉轴分离计算

3. **对齐系统**
   - justifyContent: 主轴对齐
   - alignItems: 交叉轴对齐
   - alignSelf: 单节点覆盖

4. **文本处理**
   - 宽字符支持 (中文占 2 格)
   - 自动换行
   - ANSI 序列处理

5. **终端适配**
   - 字符为单位的布局
   - 简化属性集 (30-40% 简化)
   - 性能优化 (缓存、增量)
