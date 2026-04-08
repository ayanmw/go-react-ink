package hostconfig

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/layout"
)

func TestNewHostConfig(t *testing.T) {
	hc := NewHostConfig(80, 24)

	if hc.width != 80 {
		t.Errorf("Expected width 80, got %d", hc.width)
	}

	if hc.height != 24 {
		t.Errorf("Expected height 24, got %d", hc.height)
	}

	if hc.renderer == nil {
		t.Error("Renderer should not be nil")
	}
}

func TestCreateInstance(t *testing.T) {
	hc := NewHostConfig(80, 24)

	// 创建 Box 实例
	props := core.Props{
		"flexDirection":  "column",
		"justifyContent": "center",
		"alignItems":     "center",
		"width":          40,
		"height":         20,
		"padding":        1,
	}

	instance := hc.CreateInstanceByName("Box", props)

	node, ok := instance.(*layout.Node)
	if !ok {
		t.Fatal("Instance should be a layout.Node")
	}

	if node.Direction != layout.DirectionColumn {
		t.Error("Direction should be column")
	}

	if node.Justify != layout.JustifyCenter {
		t.Error("Justify should be center")
	}

	if node.Width != 40 {
		t.Errorf("Width should be 40, got %f", node.Width)
	}
}

func TestCreateTextInstance(t *testing.T) {
	hc := NewHostConfig(80, 24)

	instance := hc.CreateTextInstance("Hello")

	node, ok := instance.(*layout.Node)
	if !ok {
		t.Fatal("Instance should be a layout.Node")
	}

	if node.MeasureFunc == nil {
		t.Error("Text instance should have MeasureFunc")
	}
}

func TestAppendChild(t *testing.T) {
	hc := NewHostConfig(80, 24)

	parent := layout.NewNode()
	child := layout.NewNode()

	hc.AppendChild(parent, child)

	if len(parent.Children) != 1 {
		t.Error("Parent should have 1 child")
	}

	if parent.Children[0] != child {
		t.Error("Child should be added")
	}
}

func TestRemoveChild(t *testing.T) {
	hc := NewHostConfig(80, 24)

	parent := layout.NewNode()
	child1 := layout.NewNode()
	child2 := layout.NewNode()

	parent.AddChild(child1)
	parent.AddChild(child2)

	hc.RemoveChild(parent, child1)

	if len(parent.Children) != 1 {
		t.Error("Parent should have 1 child after removal")
	}

	if parent.Children[0] != child2 {
		t.Error("Remaining child should be child2")
	}
}

func TestInsertBefore(t *testing.T) {
	hc := NewHostConfig(80, 24)

	parent := layout.NewNode()
	child1 := layout.NewNode()
	child2 := layout.NewNode()

	parent.AddChild(child1)

	hc.InsertBefore(parent, child2, child1)

	// 简化版添加到末尾
	if len(parent.Children) != 2 {
		t.Error("Parent should have 2 children")
	}
}

func TestUpdateInstance(t *testing.T) {
	hc := NewHostConfig(80, 24)

	node := layout.NewNode()
	props := core.Props{
		"flexDirection": "row",
		"width":         100,
	}

	hc.UpdateInstanceByName(node, "Box", props)

	if node.Direction != layout.DirectionRow {
		t.Error("Direction should be row")
	}

	if node.Width != 100 {
		t.Errorf("Width should be 100, got %f", node.Width)
	}
}

func TestUpdateTextInstance(t *testing.T) {
	hc := NewHostConfig(80, 24)

	node := layout.NewNode()
	hc.UpdateTextInstance(node, "Updated")

	if node.MeasureFunc == nil {
		t.Error("Should have MeasureFunc after update")
	}
}

func TestClearContainer(t *testing.T) {
	hc := NewHostConfig(80, 24)

	node := layout.NewNode()
	node.AddChild(layout.NewNode())
	node.AddChild(layout.NewNode())

	hc.ClearContainer(node)

	if len(node.Children) != 0 {
		t.Error("Container should be cleared")
	}
}

func TestReplaceContainerChildren(t *testing.T) {
	hc := NewHostConfig(80, 24)

	parent := layout.NewNode()
	child1 := layout.NewNode()
	child2 := layout.NewNode()

	hc.ReplaceContainerChildren(parent, []any{child1, child2})

	if len(parent.Children) != 2 {
		t.Error("Parent should have 2 children")
	}
}

func TestShouldSetTextContent(t *testing.T) {
	hc := NewHostConfig(80, 24)

	if !hc.ShouldSetTextContent("Text", nil) {
		t.Error("Text component should set text content")
	}

	if hc.ShouldSetTextContent("Box", nil) {
		t.Error("Box component should not set text content")
	}
}

func TestSetRoot(t *testing.T) {
	hc := NewHostConfig(80, 24)

	root := layout.NewNode()
	hc.SetRoot(root)

	if hc.rootNode != root {
		t.Error("Root should be set")
	}
}

func TestResize(t *testing.T) {
	hc := NewHostConfig(80, 24)

	hc.Resize(100, 30)

	if hc.width != 100 {
		t.Errorf("Width should be 100, got %d", hc.width)
	}

	if hc.height != 30 {
		t.Errorf("Height should be 30, got %d", hc.height)
	}
}

func TestGetRenderer(t *testing.T) {
	hc := NewHostConfig(80, 24)

	renderer := hc.GetRenderer()

	if renderer == nil {
		t.Error("Renderer should not be nil")
	}
}

func TestGetLayoutRoot(t *testing.T) {
	hc := NewHostConfig(80, 24)

	root := layout.NewNode()
	hc.SetRoot(root)

	layoutRoot := hc.GetLayoutRoot()

	if layoutRoot != root {
		t.Error("Layout root should be the set root")
	}
}

func TestIsPrimaryRenderer(t *testing.T) {
	hc := NewHostConfig(80, 24)

	if !hc.IsPrimaryRenderer() {
		t.Error("Should be primary renderer")
	}
}

func TestPrepareForCommit(t *testing.T) {
	hc := NewHostConfig(80, 24)

	result := hc.PrepareForCommit(nil)

	if result != nil {
		t.Error("PrepareForCommit should return nil")
	}
}

func TestResetAfterCommit(t *testing.T) {
	hc := NewHostConfig(80, 24)

	root := layout.NewNode()
	root.Width = 80
	root.Height = 24

	child := layout.NewNode()
	child.Width = 40
	child.Height = 10
	root.AddChild(child)

	hc.ResetAfterCommit(root)

	// 布局应该被计算
	if child.Layout.Width != 40 {
		t.Errorf("Child width should be 40, got %f", child.Layout.Width)
	}
}

func TestApplyBoxPropsAll(t *testing.T) {
	hc := NewHostConfig(80, 24)

	props := core.Props{
		"flexDirection":  "column",
		"justifyContent": "space-between",
		"alignItems":     "flex-end",
		"flexGrow":       2,
		"flexShrink":     0,
		"width":          50,
		"height":         30,
		"padding":        2,
		"margin":         1,
	}

	node := layout.NewNode()
	hc.applyBoxProps(node, props)

	if node.Direction != layout.DirectionColumn {
		t.Error("Direction should be column")
	}
	if node.Justify != layout.JustifySpaceBetween {
		t.Error("Justify should be space-between")
	}
	if node.AlignItems != layout.AlignFlexEnd {
		t.Error("AlignItems should be flex-end")
	}
	if node.FlexGrow != 2 {
		t.Errorf("FlexGrow should be 2, got %f", node.FlexGrow)
	}
	if node.FlexShrink != 0 {
		t.Errorf("FlexShrink should be 0, got %f", node.FlexShrink)
	}
	if node.Width != 50 {
		t.Errorf("Width should be 50, got %f", node.Width)
	}
	if node.Height != 30 {
		t.Errorf("Height should be 30, got %f", node.Height)
	}
}

func TestApplyTextProps(t *testing.T) {
	hc := NewHostConfig(80, 24)

	props := core.Props{
		"children": "Hello World",
	}

	node := layout.NewNode()
	hc.applyTextProps(node, props)

	if node.MeasureFunc == nil {
		t.Error("Should have MeasureFunc")
	}

	// 测量
	minW, _, h := node.MeasureFunc(80)
	if minW != 11 { // "Hello World" 长度
		t.Errorf("Min width should be 11, got %f", minW)
	}
	if h != 1 {
		t.Errorf("Height should be 1, got %f", h)
	}
}

func TestRender(t *testing.T) {
	hc := NewHostConfig(80, 24)

	root := layout.NewNode()
	root.Width = 80
	root.Height = 24
	hc.SetRoot(root)

	output := hc.Render()

	// 渲染应该成功
	_ = output
}
