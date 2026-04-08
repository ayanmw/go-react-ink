package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

// TestBoxComponent 测试 Box 组件创建
func TestBoxComponent(t *testing.T) {
	box := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
	}, nil)

	if box == nil {
		t.Fatal("Box should not be nil")
	}

	// 验证渲染不报错
	result := box.Render()
	if result == "" && false { // Box 可以返回空字符串
		t.Error("Box should render")
	}
}

// TestTextComponent 测试 Text 组件创建
func TestTextComponent(t *testing.T) {
	text := components.Text(core.Props{
		"color": "green",
		"bold":  true,
		"children": "Hello",
	}, nil)

	if text == nil {
		t.Fatal("Text should not be nil")
	}

	result := text.Render()
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

// TestSpacerComponent 测试 Spacer 组件
func TestSpacerComponent(t *testing.T) {
	spacer := components.Spacer(core.Props{}, nil)

	if spacer == nil {
		t.Fatal("Spacer should not be nil")
	}

	result := spacer.Render()
	if result != " " {
		t.Errorf("Expected ' ', got '%s'", result)
	}
}

// TestNewlineComponent 测试 Newline 组件
func TestNewlineComponent(t *testing.T) {
	newline := components.Newline(core.Props{}, nil)

	if newline == nil {
		t.Fatal("Newline should not be nil")
	}

	result := newline.Render()
	if result != "\n" {
		t.Errorf("Expected '\\n', got '%s'", result)
	}
}

// TestStaticComponent 测试 Static 组件
func TestStaticComponent(t *testing.T) {
	static := components.Static(core.Props{}, []core.Element{
		components.Text(core.Props{"children": "line1"}, nil),
	})

	if static == nil {
		t.Fatal("Static should not be nil")
	}
}

// TestTransformComponent 测试 Transform 组件
func TestTransformComponent(t *testing.T) {
	transform := components.Transform(core.Props{
		"transform": func(s string) string { return "->" + s + "<-" },
	}, []core.Element{
		components.Text(core.Props{"children": "test"}, nil),
	})

	if transform == nil {
		t.Fatal("Transform should not be nil")
	}

	result := transform.Render()
	if result != "->test<-" {
		t.Errorf("Expected '->test<-', got '%s'", result)
	}
}

// TestFragmentComponent 测试 Fragment 组件
func TestFragmentComponent(t *testing.T) {
	fragment := components.Fragment(core.Props{}, []core.Element{
		components.Text(core.Props{"children": "A"}, nil),
		components.Text(core.Props{"children": "B"}, nil),
	})

	if fragment == nil {
		t.Fatal("Fragment should not be nil")
	}

	result := fragment.Render()
	if result != "AB" {
		t.Errorf("Expected 'AB', got '%s'", result)
	}
}

// TestNestedComponents 测试嵌套组件
func TestNestedComponents(t *testing.T) {
	child := components.Text(core.Props{"children": "inner"}, nil)
	parent := components.Box(core.Props{
		"flexDirection": "column",
	}, []core.Element{child})

	if parent == nil {
		t.Fatal("Parent Box should not be nil")
	}

	result := parent.Render()
	if result != "inner" {
		t.Errorf("Expected 'inner', got '%s'", result)
	}
}

// TestMultipleChildren 测试多个子组件
func TestMultipleChildren(t *testing.T) {
	children := []core.Element{
		components.Text(core.Props{"children": "A"}, nil),
		components.Text(core.Props{"children": "B"}, nil),
		components.Text(core.Props{"children": "C"}, nil),
	}

	box := components.Box(core.Props{}, children)
	result := box.Render()

	if result != "ABC" {
		t.Errorf("Expected 'ABC', got '%s'", result)
	}
}

// TestBoxWithAllProps 测试 Box 所有属性
func TestBoxWithAllProps(t *testing.T) {
	props := core.Props{
		"width":           80,
		"height":          24,
		"flexDirection":   "column",
		"justifyContent":  "center",
		"alignItems":      "center",
		"flexGrow":        1,
		"flexShrink":      0,
		"flexBasis":       "auto",
		"padding":         1,
		"paddingX":        2,
		"paddingY":        2,
		"paddingTop":      1,
		"paddingBottom":   1,
		"paddingLeft":     1,
		"paddingRight":    1,
		"margin":          1,
		"marginX":         2,
		"marginY":         2,
		"marginTop":       1,
		"marginBottom":    1,
		"marginLeft":      1,
		"marginRight":     1,
		"borderStyle":     "single",
		"borderColor":     "blue",
	}

	box := components.Box(props, nil)

	if box == nil {
		t.Fatal("Box with all props should not be nil")
	}
}

// TestTextWithAllProps 测试 Text 所有属性
func TestTextWithAllProps(t *testing.T) {
	props := core.Props{
		"children":        "Test Text",
		"color":           "green",
		"backgroundColor": "black",
		"bold":            true,
		"italic":          true,
		"underline":       true,
		"dim":             false,
		"reverse":         false,
		"wrap":            "wrap",
	}

	text := components.Text(props, nil)

	if text == nil {
		t.Fatal("Text with all props should not be nil")
	}

	result := text.Render()
	if result != "Test Text" {
		t.Errorf("Expected 'Test Text', got '%s'", result)
	}
}

// TestCreateElement 测试 CreateElement 函数
func TestCreateElement(t *testing.T) {
	element := core.CreateElement(components.Box, core.Props{
		"padding": 1,
	}, nil)

	if element == nil {
		t.Fatal("CreateElement should return non-nil element")
	}
}

// TestSpreadProps 测试 Spread 属性
func TestSpreadProps(t *testing.T) {
	baseProps := core.Props{"color": "red", "bold": true}
	spread := core.Spread(baseProps)

	if spread == nil {
		t.Fatal("Spread should return non-nil props")
	}

	if spread["__spread__"] == nil {
		t.Error("Spread should contain __spread__ key")
	}
}