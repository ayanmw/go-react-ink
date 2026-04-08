package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/layout"
)

// TestFlexDirection 测试 Flex 方向
func TestFlexDirection(t *testing.T) {
	directions := []string{"row", "column", "row-reverse", "column-reverse"}

	for _, dir := range directions {
		box := components.Box(core.Props{
			"flexDirection": dir,
		}, nil)

		if box == nil {
			t.Errorf("Box with flexDirection '%s' should not be nil", dir)
		}
	}
}

// TestJustifyContent 测试 JustifyContent
func TestJustifyContent(t *testing.T) {
	values := []string{"flex-start", "flex-end", "center", "space-between", "space-around", "space-evenly"}

	for _, val := range values {
		box := components.Box(core.Props{
			"justifyContent": val,
		}, nil)

		if box == nil {
			t.Errorf("Box with justifyContent '%s' should not be nil", val)
		}
	}
}

// TestAlignItems 测试 AlignItems
func TestAlignItems(t *testing.T) {
	values := []string{"flex-start", "flex-end", "center", "stretch", "baseline"}

	for _, val := range values {
		box := components.Box(core.Props{
			"alignItems": val,
		}, nil)

		if box == nil {
			t.Errorf("Box with alignItems '%s' should not be nil", val)
		}
	}
}

// TestFlexWrap 测试 FlexWrap
func TestFlexWrap(t *testing.T) {
	values := []string{"nowrap", "wrap", "wrap-reverse"}

	for _, val := range values {
		box := components.Box(core.Props{
			"flexWrap": val,
		}, nil)

		if box == nil {
			t.Errorf("Box with flexWrap '%s' should not be nil", val)
		}
	}
}

// TestFlexGrowShrink 测试 FlexGrow 和 FlexShrink
func TestFlexGrowShrink(t *testing.T) {
	box := components.Box(core.Props{
		"flexGrow":   1,
		"flexShrink": 0,
		"flexBasis":  "auto",
	}, nil)

	if box == nil {
		t.Fatal("Box with flex props should not be nil")
	}
}

// TestWidthHeight 测试宽高
func TestWidthHeight(t *testing.T) {
	box := components.Box(core.Props{
		"width":  80,
		"height": 24,
	}, nil)

	if box == nil {
		t.Fatal("Box with width/height should not be nil")
	}
}

// TestMinMaxDimensions 测试最小最大尺寸
func TestMinMaxDimensions(t *testing.T) {
	box := components.Box(core.Props{
		"minWidth":  20,
		"maxWidth":  100,
		"minHeight": 10,
		"maxHeight": 50,
	}, nil)

	if box == nil {
		t.Fatal("Box with min/max dimensions should not be nil")
	}
}

// TestPadding 测试内边距
func TestPadding(t *testing.T) {
	// 简写
	box1 := components.Box(core.Props{"padding": 2}, nil)
	if box1 == nil {
		t.Error("Box with padding should not be nil")
	}

	// 分方向
	box2 := components.Box(core.Props{
		"paddingTop":    1,
		"paddingRight":  2,
		"paddingBottom": 3,
		"paddingLeft":   4,
	}, nil)
	if box2 == nil {
		t.Error("Box with directional padding should not be nil")
	}

	// 轴向
	box3 := components.Box(core.Props{
		"paddingX": 2,
		"paddingY": 1,
	}, nil)
	if box3 == nil {
		t.Error("Box with axis padding should not be nil")
	}
}

// TestMargin 测试外边距
func TestMargin(t *testing.T) {
	// 简写
	box1 := components.Box(core.Props{"margin": 2}, nil)
	if box1 == nil {
		t.Error("Box with margin should not be nil")
	}

	// 分方向
	box2 := components.Box(core.Props{
		"marginTop":    1,
		"marginRight":  2,
		"marginBottom": 3,
		"marginLeft":   4,
	}, nil)
	if box2 == nil {
		t.Error("Box with directional margin should not be nil")
	}

	// 轴向
	box3 := components.Box(core.Props{
		"marginX": 2,
		"marginY": 1,
	}, nil)
	if box3 == nil {
		t.Error("Box with axis margin should not be nil")
	}
}

// TestBorder 测试边框
func TestBorder(t *testing.T) {
	styles := []string{"single", "double", "round", "bold", "none"}

	for _, style := range styles {
		box := components.Box(core.Props{
			"borderStyle": style,
			"borderColor": "blue",
		}, nil)

		if box == nil {
			t.Errorf("Box with borderStyle '%s' should not be nil", style)
		}
	}
}

// TestLayoutNode 测试布局节点
func TestLayoutNode(t *testing.T) {
	node := layout.NewNode()
	node.Width = 100
	node.Height = 50
	node.Direction = layout.DirectionColumn

	if node == nil {
		t.Fatal("Layout node should not be nil")
	}

	if node.Width != 100 {
		t.Errorf("Expected width 100, got %f", node.Width)
	}
}

// TestLayoutCalculation 测试布局计算
func TestLayoutCalculation(t *testing.T) {
	root := layout.NewNode()
	root.Width = 80
	root.Height = 24
	root.Direction = layout.DirectionColumn

	child1 := layout.NewNode()
	child1.Height = 1
	root.AddChild(child1)

	child2 := layout.NewNode()
	child2.FlexGrow = 1
	root.AddChild(child2)

	child3 := layout.NewNode()
	child3.Height = 1
	root.AddChild(child3)

	root.CalculateLayout(80, 24)

	// 验证布局计算
	if root.Layout.Width != 80 {
		t.Errorf("Expected root width 80, got %f", root.Layout.Width)
	}
}

// TestNestedLayout 测试嵌套布局
func TestNestedLayout(t *testing.T) {
	root := layout.NewNode()
	root.Direction = layout.DirectionColumn
	root.Width = 80
	root.Height = 24

	row := layout.NewNode()
	row.Direction = layout.DirectionRow
	row.Height = 1
	root.AddChild(row)

	left := layout.NewNode()
	left.Width = 40
	row.AddChild(left)

	right := layout.NewNode()
	right.FlexGrow = 1
	row.AddChild(right)

	root.CalculateLayout(80, 24)

	if root.Layout.Width != 80 {
		t.Errorf("Expected root width 80, got %f", root.Layout.Width)
	}
}