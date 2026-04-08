package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

// TestCreateElementJSX 测试 JSX 风格的 CreateElement
func TestCreateElementJSX(t *testing.T) {
	// 类似 JSX: <Box><Text>Hello</Text></Box>
	box := core.CreateElement(components.Box, core.Props{
		"flexDirection": "column",
	}, []core.Element{
		core.CreateElement(components.Text, core.Props{
			"children": "Hello",
		}, nil),
	})

	if box == nil {
		t.Fatal("CreateElement should return non-nil element")
	}

	result := box.Render()
	if result == "" {
		t.Error("Box should render content")
	}
}

// TestNestedJSXStyle 测试嵌套 JSX 风格
func TestNestedJSXStyle(t *testing.T) {
	// <Box flexDirection="column">
	//   <Box flexDirection="row">
	//     <Text>A</Text>
	//     <Text>B</Text>
	//   </Box>
	//   <Text>C</Text>
	// </Box>
	root := core.CreateElement(components.Box, core.Props{
		"flexDirection": "column",
	}, []core.Element{
		core.CreateElement(components.Box, core.Props{
			"flexDirection": "row",
		}, []core.Element{
			core.CreateElement(components.Text, core.Props{"children": "A"}, nil),
			core.CreateElement(components.Text, core.Props{"children": "B"}, nil),
		}),
		core.CreateElement(components.Text, core.Props{"children": "C"}, nil),
	})

	if root == nil {
		t.Fatal("Nested JSX should work")
	}

	result := root.Render()
	if result == "" {
		t.Error("Should render nested content")
	}
}

// TestSpreadAttributes 测试 spread 属性
func TestSpreadAttributes(t *testing.T) {
	baseProps := core.Props{
		"color": "green",
		"bold":  true,
	}

	text := core.CreateElement(components.Text, core.Props{
		"__spread__": baseProps,
		"children":   "Styled Text",
	}, nil)

	if text == nil {
		t.Fatal("Spread attributes should work")
	}
}

// TestFragmentWithChildren 测试 Fragment 子元素
func TestFragmentWithChildren(t *testing.T) {
	fragment := core.CreateElement(components.Fragment, core.Props{}, []core.Element{
		core.CreateElement(components.Text, core.Props{"children": "A"}, nil),
		core.CreateElement(components.Text, core.Props{"children": "B"}, nil),
		core.CreateElement(components.Text, core.Props{"children": "C"}, nil),
	})

	result := fragment.Render()
	if result != "ABC" {
		t.Errorf("Expected 'ABC', got '%s'", result)
	}
}

// TestConditionalRendering 测试条件渲染
func TestConditionalRendering(t *testing.T) {
	showExtra := true

	children := []core.Element{
		core.CreateElement(components.Text, core.Props{"children": "Always"}, nil),
	}

	if showExtra {
		children = append(children, core.CreateElement(components.Text, core.Props{"children": "Extra"}, nil))
	}

	box := core.CreateElement(components.Box, core.Props{}, children)
	if box == nil {
		t.Fatal("Conditional rendering should work")
	}
}

// TestListRendering 测试列表渲染
func TestListRendering(t *testing.T) {
	items := []string{"Item 1", "Item 2", "Item 3"}

	children := make([]core.Element, len(items))
	for i, item := range items {
		children[i] = core.CreateElement(components.Text, core.Props{
			"children": item,
			"key":      i,
		}, nil)
	}

	box := core.CreateElement(components.Box, core.Props{
		"flexDirection": "column",
	}, children)

	if box == nil {
		t.Fatal("List rendering should work")
	}
}

// TestMapRendering 测试 map 渲染
func TestMapRendering(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	children := make([]core.Element, 0)
	for _, n := range numbers {
		children = append(children, core.CreateElement(components.Text, core.Props{
			"children": n,
		}, nil))
	}

	box := core.CreateElement(components.Box, core.Props{}, children)
	if box == nil {
		t.Fatal("Map rendering should work")
	}
}

// TestComponentComposition 测试组件组合
func TestComponentComposition(t *testing.T) {
	// 定义一个自定义组件
	headerComponent := func(props core.Props, children []core.Element) core.Element {
		return core.CreateElement(components.Box, core.Props{
			"flexDirection": "row",
			"padding":       1,
		}, []core.Element{
			core.CreateElement(components.Text, core.Props{
				"color":    "blue",
				"bold":     true,
				"children": props["title"],
			}, nil),
		})
	}

	header := core.CreateElement(headerComponent, core.Props{
		"title": "My App",
	}, nil)

	if header == nil {
		t.Fatal("Component composition should work")
	}
}

// TestPropsMerging 测试属性合并
func TestPropsMerging(t *testing.T) {
	defaultProps := core.Props{
		"color":     "white",
		"bold":      false,
		"underline": false,
	}

	userProps := core.Props{
		"color": "green",
		"bold":  true,
	}

	// 合并属性 (用户属性覆盖默认值)
	merged := core.Props{}
	for k, v := range defaultProps {
		merged[k] = v
	}
	for k, v := range userProps {
		merged[k] = v
	}

	if merged["color"] != "green" {
		t.Error("User props should override default props")
	}
	if !merged["bold"].(bool) {
		t.Error("User bold should override default")
	}
	if merged["underline"] != false {
		t.Error("Default underline should remain")
	}
}

// TestChildrenProp 测试 children 属性
func TestChildrenProp(t *testing.T) {
	children := []core.Element{
		core.CreateElement(components.Text, core.Props{"children": "Child 1"}, nil),
		core.CreateElement(components.Text, core.Props{"children": "Child 2"}, nil),
	}

	box := core.CreateElement(components.Box, core.Props{
		"children": children,
	}, nil)

	if box == nil {
		t.Fatal("Children prop should work")
	}
}

// TestTextInterpolation 测试文本插值
func TestTextInterpolation(t *testing.T) {
	name := "World"
	count := 42

	// 模拟模板字符串
	text := core.CreateElement(components.Text, core.Props{
		"children": "Hello " + name + ", count: " + string(rune(count+'0')),
	}, nil)

	if text == nil {
		t.Fatal("Text interpolation should work")
	}
}

// TestStyleObject 测试样式对象
func TestStyleObject(t *testing.T) {
	style := map[string]any{
		"color":           "red",
		"backgroundColor": "black",
		"padding":         1,
	}

	text := core.CreateElement(components.Text, core.Props{
		"children": "Styled",
		"style":    style,
	}, nil)

	if text == nil {
		t.Fatal("Style object should work")
	}
}

// TestEmptyChildren 测试空子元素
func TestEmptyChildren(t *testing.T) {
	box := core.CreateElement(components.Box, core.Props{}, nil)
	if box == nil {
		t.Fatal("Empty children should work")
	}

	box2 := core.CreateElement(components.Box, core.Props{}, []core.Element{})
	if box2 == nil {
		t.Fatal("Empty children slice should work")
	}
}

// TestKeyProp 测试 key 属性
func TestKeyProp(t *testing.T) {
	items := []struct {
		id    string
		value string
	}{
		{"a", "Item A"},
		{"b", "Item B"},
		{"c", "Item C"},
	}

	children := make([]core.Element, len(items))
	for i, item := range items {
		children[i] = core.CreateElement(components.Text, core.Props{
			"key":     item.id,
			"children": item.value,
		}, nil)
	}

	box := core.CreateElement(components.Box, core.Props{}, children)
	if box == nil {
		t.Fatal("Key prop should work")
	}
}

// TestRefProp 测试 ref 属性
func TestRefProp(t *testing.T) {
	ref := &struct {
		Current any
	}{}

	box := core.CreateElement(components.Box, core.Props{
		"ref": ref,
	}, nil)

	if box == nil {
		t.Fatal("Ref prop should work")
	}
}
