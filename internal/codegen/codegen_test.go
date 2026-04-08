package codegen

import (
	"strings"
	"testing"
)

func TestGenerateSimpleElement(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证生成代码包含 CreateElement
	if !strings.Contains(code, "CreateElement") {
		t.Error("Generated code should contain CreateElement")
	}

	// 验证包含 Box
	if !strings.Contains(code, "ink.Box") {
		t.Error("Generated code should contain 'ink.Box'")
	}
}

func TestGenerateElementWithAttrs(t *testing.T) {
	node := Node{
		Type:       NodeElement,
		TagName:    "Box",
		Attributes: map[string]AttrValue{
			"flexDirection": {IsExpression: false, Value: "column"},
			"padding":       {IsExpression: true, Value: "1"},
		},
		Children: []Node{},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证属性生成
	if !strings.Contains(code, `"flexDirection"`) {
		t.Error("Generated code should contain 'flexDirection' attribute")
	}
	if !strings.Contains(code, `"column"`) {
		t.Error("Generated code should contain 'column' value")
	}
	if !strings.Contains(code, "padding") {
		t.Error("Generated code should contain 'padding' attribute")
	}
	// 表达式值应该直接使用，不加引号
	if strings.Contains(code, `"1"`) && strings.Contains(code, "padding") {
		// 检查是否是表达式形式
		lines := strings.Split(code, "\n")
		for _, line := range lines {
			if strings.Contains(line, "padding") && strings.Contains(line, `"1"`) {
				t.Error("Expression value should not be quoted")
			}
		}
	}
}

func TestGenerateElementWithChildren(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{
			{Type: NodeElement, TagName: "Text", Children: []Node{}},
		},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证子节点生成
	if !strings.Contains(code, "ink.Text") {
		t.Error("Generated code should contain child 'ink.Text'")
	}
	if !strings.Contains(code, "[]core.Element") {
		t.Error("Generated code should contain children array")
	}
}

func TestGenerateText(t *testing.T) {
	node := Node{
		Type:  NodeText,
		Value: "Hello World",
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证文本生成
	if !strings.Contains(code, "Text") {
		t.Error("Generated code should contain 'Text'")
	}
	if !strings.Contains(code, "Hello World") {
		t.Error("Generated code should contain text content")
	}
}

func TestGenerateExpression(t *testing.T) {
	node := Node{
		Type:  NodeExpression,
		Value: "count",
	}

	gen := New("ink")
	code := gen.generateNode(node)

	// 表达式应该直接返回其值
	if code != "count" {
		t.Errorf("Expected 'count', got %q", code)
	}
}

func TestGenerateFragment(t *testing.T) {
	node := Node{
		Type:     NodeFragment,
		Children: []Node{
			{Type: NodeText, Value: "A"},
			{Type: NodeText, Value: "B"},
		},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// Fragment 多子节点应该生成数组
	if !strings.Contains(code, "[]core.Element") {
		t.Error("Fragment with multiple children should generate array")
	}
}

func TestGenerateFragmentSingleChild(t *testing.T) {
	node := Node{
		Type:     NodeFragment,
		Children: []Node{
			{Type: NodeText, Value: "Single"},
		},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// Fragment 单子节点应该直接返回
	if strings.Contains(code, "[]core.Element") {
		t.Error("Fragment with single child should not generate array")
	}
	if !strings.Contains(code, "Single") {
		t.Error("Generated code should contain child content")
	}
}

func TestGenerateSpread(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Spread:   []string{"props"},
		Children: []Node{},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证 spread 生成
	if !strings.Contains(code, "Spread(props)") {
		t.Error("Generated code should contain Spread(props)")
	}
}

func TestGenerateNestedElements(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{
			{
				Type:     NodeElement,
				TagName:  "Text",
				Children: []Node{
					{Type: NodeText, Value: "Hello"},
				},
			},
		},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证嵌套结构
	if !strings.Contains(code, "ink.Box") {
		t.Error("Generated code should contain 'ink.Box'")
	}
	if !strings.Contains(code, "ink.Text") {
		t.Error("Generated code should contain 'ink.Text'")
	}
	if !strings.Contains(code, "Hello") {
		t.Error("Generated code should contain 'Hello'")
	}
}

func TestGenerateFunction(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{},
	}

	code := GenerateFunction("App", node, "ink")

	// 验证函数生成
	if !strings.Contains(code, "func App()") {
		t.Error("Generated code should contain function declaration")
	}
	if !strings.Contains(code, "return") {
		t.Error("Generated code should contain return statement")
	}
}

func TestGenerateReturn(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{},
	}

	code := GenerateReturn(node, "ink")

	// 验证 return 语句
	if !strings.HasPrefix(code, "return ") {
		t.Error("Generated code should start with 'return'")
	}
}

func TestGenerateAssignment(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{},
	}

	code := GenerateAssignment("element", node, "ink")

	// 验证赋值语句
	if !strings.HasPrefix(code, "element :=") {
		t.Error("Generated code should start with 'element :='")
	}
}

func TestGenerateImports(t *testing.T) {
	node := Node{
		Type:     NodeElement,
		TagName:  "Box",
		Children: []Node{},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 验证导入
	if !strings.Contains(code, "import") {
		t.Error("Generated code should contain import statement")
	}
	if !strings.Contains(code, "github.com/ayanmw/go-react-ink/pkg/core") {
		t.Error("Generated code should contain core package import")
	}
}

func TestGenerateBooleanAttr(t *testing.T) {
	node := Node{
		Type:       NodeElement,
		TagName:    "Box",
		Attributes: map[string]AttrValue{
			"flexGrow": {IsExpression: false, Value: "true"},
		},
		Children: []Node{},
	}

	gen := New("ink")
	code := gen.Generate(node)

	// 布尔属性应该生成字符串
	if !strings.Contains(code, `"flexGrow"`) {
		t.Error("Generated code should contain 'flexGrow' attribute")
	}
	if !strings.Contains(code, `"true"`) {
		t.Error("Boolean attribute should be string 'true'")
	}
}