// Package layout 实现 Flexbox 布局引擎 (纯 Go 实现)
package layout

import (
	"math"
)

// Direction Flex 方向
type Direction int

const (
	DirectionRow Direction = iota
	DirectionRowReverse
	DirectionColumn
	DirectionColumnReverse
)

// Justify 主轴对齐
type Justify int

const (
	JustifyFlexStart Justify = iota
	JustifyFlexEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Align 交叉轴对齐
type Align int

const (
	AlignAuto Align = iota
	AlignFlexStart
	AlignFlexEnd
	AlignCenter
	AlignStretch
	AlignBaseline
)

// Wrap 换行模式
type Wrap int

const (
	WrapNoWrap Wrap = iota
	WrapWrap
	WrapWrapReverse
)

// Position 定位类型
type Position int

const (
	PositionRelative Position = iota
	PositionAbsolute
)

// Overflow 溢出处理
type Overflow int

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

// Display 显示类型
type Display int

const (
	DisplayFlex Display = iota
	DisplayNone
)

// Node 布局节点
type Node struct {
	// 样式属性
	Direction     Direction
	Justify       Justify
	AlignItems    Align
	AlignSelf     Align
	AlignContent  Align
	Wrap          Wrap
	Display       Display
	Position      Position
	Overflow      Overflow

	// 尺寸
	Width    float64
	Height   float64
	MinWidth  float64
	MinHeight float64
	MaxWidth  float64
	MaxHeight float64

	// 边距
	Margin       [4]float64 // top, right, bottom, left
	Padding      [4]float64
	Border       [4]float64

	// Flex 属性
	FlexGrow   float64
	FlexShrink float64
	FlexBasis  float64

	// 比例
	AspectRatio float64

	// 测量函数
	MeasureFunc func(width float64) (minWidth, maxWidth, height float64)

	// 子节点
	Children []*Node

	// 计算结果
	Layout Layout
}

// Layout 计算后的布局
type Layout struct {
	X           float64
	Y           float64
	Width       float64
	Height      float64
	ComputedBasis float64
}

// NewNode 创建新节点
func NewNode() *Node {
	return &Node{
		Direction:    DirectionRow,
		Justify:      JustifyFlexStart,
		AlignItems:   AlignStretch,
		AlignSelf:    AlignAuto,
		AlignContent: AlignStretch,
		Wrap:         WrapNoWrap,
		Display:      DisplayFlex,
		Position:     PositionRelative,
		Overflow:     OverflowVisible,
		FlexGrow:     0,
		FlexShrink:   1,
		FlexBasis:    -1, // auto
	}
}

// AddChild 添加子节点
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// RemoveChild 移除子节点
func (n *Node) RemoveChild(child *Node) {
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}

// CalculateLayout 计算布局
func (n *Node) CalculateLayout(width, height float64) {
	// 设置容器尺寸
	n.Layout.Width = width
	n.Layout.Height = height

	// 开始布局计算
	n.calculateNodeLayout(width, height)
}

// calculateNodeLayout 计算节点布局
func (n *Node) calculateNodeLayout(width, height float64) {
	if n.Display == DisplayNone {
		return
	}

	// 应用边距
	contentWidth := width - n.Padding[1] - n.Padding[3] - n.Border[1] - n.Border[3]
	contentHeight := height - n.Padding[0] - n.Padding[2] - n.Border[0] - n.Border[2]

	if contentWidth < 0 {
		contentWidth = 0
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	// 处理子节点
	if len(n.Children) == 0 {
		// 叶子节点，使用测量函数或固定尺寸
		if n.MeasureFunc != nil {
			minW, maxW, h := n.MeasureFunc(contentWidth)
			_ = minW
			_ = maxW
			n.Layout.Height = h + n.Padding[0] + n.Padding[2] + n.Border[0] + n.Border[2]
		} else {
			if n.Width > 0 {
				n.Layout.Width = n.Width
			}
			if n.Height > 0 {
				n.Layout.Height = n.Height
			}
		}
		return
	}

	// Flexbox 布局
	if n.Direction == DirectionRow || n.Direction == DirectionRowReverse {
		n.layoutRow(contentWidth, contentHeight)
	} else {
		n.layoutColumn(contentWidth, contentHeight)
	}

	// 递归计算子节点
	for _, child := range n.Children {
		child.calculateNodeLayout(child.Layout.Width, child.Layout.Height)
	}
}

// layoutRow 水平布局
func (n *Node) layoutRow(width, height float64) {
	children := n.getChildren()

	// 第一遍：计算固定尺寸和 flex basis
	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	usedWidth := 0.0

	for _, child := range children {
		if child.Display == DisplayNone {
			continue
		}

		// 计算子节点尺寸
		childWidth := child.getWidth(width)
		childHeight := child.getHeight(height)

		// 应用 flex basis
		if child.FlexBasis >= 0 {
			childWidth = child.FlexBasis
		}

		child.Layout.ComputedBasis = childWidth
		child.Layout.Width = childWidth
		child.Layout.Height = childHeight

		if child.FlexGrow > 0 {
			totalFlexGrow += child.FlexGrow
		} else {
			usedWidth += childWidth + child.Margin[1] + child.Margin[3]
		}

		if child.FlexShrink > 0 {
			totalFlexShrink += child.FlexShrink
		}
	}

	// 第二遍：分配剩余空间
	remainingWidth := width - usedWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	if totalFlexGrow > 0 && remainingWidth > 0 {
		// 分配 flex grow 空间
		for _, child := range children {
			if child.Display == DisplayNone || child.FlexGrow <= 0 {
				continue
			}
			flexWidth := (remainingWidth / totalFlexGrow) * child.FlexGrow
			child.Layout.Width += flexWidth
		}
	}

	// 第三遍：定位子节点
	x := n.Padding[3] + n.Border[3]
	y := n.Padding[0] + n.Border[0]

	// 应用 justify
	if n.Justify == JustifyCenter {
		totalWidth := 0.0
		for _, child := range children {
			if child.Display != DisplayNone {
				totalWidth += child.Layout.Width + child.Margin[1] + child.Margin[3]
			}
		}
		x += (width - totalWidth) / 2
	} else if n.Justify == JustifyFlexEnd {
		totalWidth := 0.0
		for _, child := range children {
			if child.Display != DisplayNone {
				totalWidth += child.Layout.Width + child.Margin[1] + child.Margin[3]
			}
		}
		x += width - totalWidth
	} else if n.Justify == JustifySpaceBetween && len(children) > 1 {
		totalWidth := 0.0
		for _, child := range children {
			if child.Display != DisplayNone {
				totalWidth += child.Layout.Width + child.Margin[1] + child.Margin[3]
			}
		}
		spacing := (width - totalWidth) / float64(len(children)-1)
		// 将在定位时应用
		_ = spacing
	}

	// 定位子节点
	for i, child := range children {
		if child.Display == DisplayNone {
			continue
		}

		child.Layout.X = x + child.Margin[3]
		child.Layout.Y = y + child.Margin[0]

		// 应用 align
		childHeight := child.Layout.Height
		switch child.AlignSelf {
		case AlignCenter:
			child.Layout.Y = y + (height-childHeight)/2
		case AlignFlexEnd:
			child.Layout.Y = y + height - childHeight - child.Margin[2]
		case AlignStretch:
			child.Layout.Height = height - child.Margin[0] - child.Margin[2]
		}

		x += child.Layout.Width + child.Margin[1] + child.Margin[3]

		// Space between
		if n.Justify == JustifySpaceBetween && i < len(children)-1 {
			totalWidth := 0.0
			for _, c := range children {
				if c.Display != DisplayNone {
					totalWidth += c.Layout.Width + c.Margin[1] + c.Margin[3]
				}
			}
			spacing := (width - totalWidth) / float64(len(children)-1)
			x += spacing
		}
	}
}

// layoutColumn 垂直布局
func (n *Node) layoutColumn(width, height float64) {
	children := n.getChildren()

	// 第一遍：计算固定尺寸和 flex basis
	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	usedHeight := 0.0

	for _, child := range children {
		if child.Display == DisplayNone {
			continue
		}

		childWidth := child.getWidth(width)
		childHeight := child.getHeight(height)

		if child.FlexBasis >= 0 {
			childHeight = child.FlexBasis
		}

		child.Layout.ComputedBasis = childHeight
		child.Layout.Width = childWidth
		child.Layout.Height = childHeight

		if child.FlexGrow > 0 {
			totalFlexGrow += child.FlexGrow
		} else {
			usedHeight += childHeight + child.Margin[0] + child.Margin[2]
		}

		if child.FlexShrink > 0 {
			totalFlexShrink += child.FlexShrink
		}
	}

	// 第二遍：分配剩余空间
	remainingHeight := height - usedHeight
	if remainingHeight < 0 {
		remainingHeight = 0
	}

	if totalFlexGrow > 0 && remainingHeight > 0 {
		for _, child := range children {
			if child.Display == DisplayNone || child.FlexGrow <= 0 {
				continue
			}
			flexHeight := (remainingHeight / totalFlexGrow) * child.FlexGrow
			child.Layout.Height += flexHeight
		}
	}

	// 第三遍：定位子节点
	x := n.Padding[3] + n.Border[3]
	y := n.Padding[0] + n.Border[0]

	// 应用 justify
	if n.Justify == JustifyCenter {
		totalHeight := 0.0
		for _, child := range children {
			if child.Display != DisplayNone {
				totalHeight += child.Layout.Height + child.Margin[0] + child.Margin[2]
			}
		}
		y += (height - totalHeight) / 2
	} else if n.Justify == JustifyFlexEnd {
		totalHeight := 0.0
		for _, child := range children {
			if child.Display != DisplayNone {
				totalHeight += child.Layout.Height + child.Margin[0] + child.Margin[2]
			}
		}
		y += height - totalHeight
	}

	// 定位子节点
	for i, child := range children {
		if child.Display == DisplayNone {
			continue
		}

		child.Layout.X = x + child.Margin[3]
		child.Layout.Y = y + child.Margin[0]

		// 应用 align
		childWidth := child.Layout.Width
		switch child.AlignSelf {
		case AlignCenter:
			child.Layout.X = x + (width-childWidth)/2
		case AlignFlexEnd:
			child.Layout.X = x + width - childWidth - child.Margin[1]
		case AlignStretch:
			child.Layout.Width = width - child.Margin[1] - child.Margin[3]
		}

		y += child.Layout.Height + child.Margin[0] + child.Margin[2]

		// Space between
		if n.Justify == JustifySpaceBetween && i < len(children)-1 {
			totalHeight := 0.0
			for _, c := range children {
				if c.Display != DisplayNone {
					totalHeight += c.Layout.Height + c.Margin[0] + c.Margin[2]
				}
			}
			spacing := (height - totalHeight) / float64(len(children)-1)
			y += spacing
		}
	}
}

// getChildren 获取子节点 (考虑 direction reverse)
func (n *Node) getChildren() []*Node {
	if n.Direction == DirectionRowReverse || n.Direction == DirectionColumnReverse {
		// 反转子节点
		children := make([]*Node, len(n.Children))
		copy(children, n.Children)
		for i, j := 0, len(children)-1; i < j; i, j = i+1, j-1 {
			children[i], children[j] = children[j], children[i]
		}
		return children
	}
	return n.Children
}

// getWidth 获取宽度
func (n *Node) getWidth(parentWidth float64) float64 {
	if n.Width > 0 {
		return n.Width
	}
	if n.FlexGrow > 0 {
		return 0 // 将由 flex 分配
	}
	return parentWidth - n.Margin[1] - n.Margin[3]
}

// getHeight 获取高度
func (n *Node) getHeight(parentHeight float64) float64 {
	if n.Height > 0 {
		return n.Height
	}
	if n.FlexGrow > 0 {
		return 0
	}
	return parentHeight - n.Margin[0] - n.Margin[2]
}

// SetMargin 设置边距
func (n *Node) SetMargin(top, right, bottom, left float64) {
	n.Margin = [4]float64{top, right, bottom, left}
}

// SetPadding 设置内边距
func (n *Node) SetPadding(top, right, bottom, left float64) {
	n.Padding = [4]float64{top, right, bottom, left}
}

// SetBorder 设置边框
func (n *Node) SetBorder(top, right, bottom, left float64) {
	n.Border = [4]float64{top, right, bottom, left}
}

// GetLayout 获取布局结果
func (n *Node) GetLayout() Layout {
	return n.Layout
}

// IsNaN 检查是否为 NaN
func isNaN(f float64) bool {
	return f != f
}

// Max 取最大值
func max(a, b float64) float64 {
	if a > b || isNaN(b) {
		return a
	}
	return b
}

// Min 取最小值
func min(a, b float64) float64 {
	if a < b || isNaN(b) {
		return a
	}
	return b
}

// Abs 取绝对值
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// Round 四舍五入
func round(f float64) float64 {
	return math.Floor(f + 0.5)
}