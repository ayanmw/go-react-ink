package parser

import (
	"testing"

	"github.com/anmingwei/go-ink/internal/lexer"
)

func TestParseSimpleElement(t *testing.T) {
	source := "<Box></Box>"
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if node.Type != NodeElement {
		t.Errorf("Expected Element node, got %v", node.Type)
	}

	if node.TagName != "Box" {
		t.Errorf("Expected tag name 'Box', got %q", node.TagName)
	}

	if len(node.Children) != 0 {
		t.Errorf("Expected no children, got %d", len(node.Children))
	}
}

func TestParseElementWithAttrs(t *testing.T) {
	source := `<Box flexDirection="column" padding={1}></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if node.TagName != "Box" {
		t.Errorf("Expected tag name 'Box', got %q", node.TagName)
	}

	// 检查属性
	if len(node.Attributes) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(node.Attributes))
	}

	flexDir, ok := node.Attributes["flexDirection"]
	if !ok {
		t.Error("Missing attribute 'flexDirection'")
	} else {
		if flexDir.IsExpression {
			t.Error("flexDirection should not be expression")
		}
		if flexDir.Value != "column" {
			t.Errorf("Expected flexDirection 'column', got %q", flexDir.Value)
		}
	}

	padding, ok := node.Attributes["padding"]
	if !ok {
		t.Error("Missing attribute 'padding'")
	} else {
		if !padding.IsExpression {
			t.Error("padding should be expression")
		}
		if padding.Value != "1" {
			t.Errorf("Expected padding '1', got %q", padding.Value)
		}
	}
}

func TestParseNestedElements(t *testing.T) {
	source := `<Box><Text>Hello</Text></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(node.Children))
	}

	child := node.Children[0]
	if child.Type != NodeElement {
		t.Errorf("Expected Element child, got %v", child.Type)
	}
	if child.TagName != "Text" {
		t.Errorf("Expected child tag 'Text', got %q", child.TagName)
	}

	// 检查 Text 的子节点
	if len(child.Children) != 1 {
		t.Fatalf("Expected 1 grandchild, got %d", len(child.Children))
	}

	text := child.Children[0]
	if text.Type != NodeText {
		t.Errorf("Expected Text node, got %v", text.Type)
	}
	if text.Value != "Hello" {
		t.Errorf("Expected text 'Hello', got %q", text.Value)
	}
}

func TestParseTextContent(t *testing.T) {
	source := `<Text>Hello World</Text>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(node.Children))
	}

	text := node.Children[0]
	if text.Type != NodeText {
		t.Errorf("Expected Text node, got %v", text.Type)
	}
	if text.Value != "Hello World" {
		t.Errorf("Expected text 'Hello World', got %q", text.Value)
	}
}

func TestParseExpression(t *testing.T) {
	source := `<Text>Count: {count}</Text>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 应该有两个子节点: 文本 + 表达式
	if len(node.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(node.Children))
	}

	text := node.Children[0]
	if text.Type != NodeText {
		t.Errorf("Expected first child to be Text, got %v", text.Type)
	}

	expr := node.Children[1]
	if expr.Type != NodeExpression {
		t.Errorf("Expected second child to be Expression, got %v", expr.Type)
	}
	if expr.Value != "count" {
		t.Errorf("Expected expression 'count', got %q", expr.Value)
	}
}

func TestParseSelfClosing(t *testing.T) {
	source := `<Box><Spacer /></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(node.Children))
	}

	child := node.Children[0]
	if child.Type != NodeElement {
		t.Errorf("Expected Element child, got %v", child.Type)
	}
	if child.TagName != "Spacer" {
		t.Errorf("Expected tag 'Spacer', got %q", child.TagName)
	}
	if len(child.Children) != 0 {
		t.Errorf("Self-closing should have no children")
	}
}

func TestParseSpreadAttr(t *testing.T) {
	source := `<Box {...props}></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Spread) != 1 {
		t.Fatalf("Expected 1 spread attribute, got %d", len(node.Spread))
	}

	if node.Spread[0] != "props" {
		t.Errorf("Expected spread 'props', got %q", node.Spread[0])
	}
}

func TestParseFragment(t *testing.T) {
	source := `<>Content</>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if node.Type != NodeFragment {
		t.Errorf("Expected Fragment node, got %v", node.Type)
	}

	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(node.Children))
	}

	text := node.Children[0]
	if text.Type != NodeText {
		t.Errorf("Expected Text child, got %v", text.Type)
	}
}

func TestParseMultipleChildren(t *testing.T) {
	source := `<Box><Text>A</Text><Text>B</Text><Text>C</Text></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Children) != 3 {
		t.Fatalf("Expected 3 children, got %d", len(node.Children))
	}

	expected := []string{"A", "B", "C"}
	for i, child := range node.Children {
		if child.TagName != "Text" {
			t.Errorf("Child %d: expected Text, got %s", i, child.TagName)
		}
		if len(child.Children) != 1 {
			t.Errorf("Child %d: expected 1 grandchild", i)
			continue
		}
		text := child.Children[0]
		if text.Value != expected[i] {
			t.Errorf("Child %d: expected text %q, got %q", i, expected[i], text.Value)
		}
	}
}

func TestParseDeepNesting(t *testing.T) {
	source := `<A><B><C><D>Deep</D></C></B></A>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 遍历到最深层
	current := node
	expectedTags := []string{"A", "B", "C", "D"}
	for i, expected := range expectedTags {
		if current.TagName != expected {
			t.Errorf("Level %d: expected tag %s, got %s", i, expected, current.TagName)
		}
		if i < len(expectedTags)-1 {
			if len(current.Children) != 1 {
				t.Fatalf("Level %d: expected 1 child", i)
			}
			current = current.Children[0]
		}
	}

	// 检查最深层文本
	if len(current.Children) != 1 {
		t.Fatalf("Expected 1 text child at deepest level")
	}
	if current.Children[0].Type != NodeText {
		t.Errorf("Expected Text node at deepest level")
	}
	if current.Children[0].Value != "Deep" {
		t.Errorf("Expected text 'Deep', got %q", current.Children[0].Value)
	}
}

func TestParseComplexExpression(t *testing.T) {
	source := `<Text>{item.name + " " + item.value}</Text>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(node.Children))
	}

	expr := node.Children[0]
	if expr.Type != NodeExpression {
		t.Errorf("Expected Expression node, got %v", expr.Type)
	}
}

func TestParseMixedContent(t *testing.T) {
	source := `<Text>Before {expr} After</Text>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 应该有三个子节点
	if len(node.Children) != 3 {
		t.Fatalf("Expected 3 children, got %d", len(node.Children))
	}

	if node.Children[0].Type != NodeText {
		t.Errorf("First child should be Text")
	}
	if node.Children[1].Type != NodeExpression {
		t.Errorf("Second child should be Expression")
	}
	if node.Children[2].Type != NodeText {
		t.Errorf("Third child should be Text")
	}
}

func TestParseBooleanAttr(t *testing.T) {
	source := `<Box flexGrow></Box>`
	l := lexer.New(source)
	tokens := l.Lex()

	p := New(tokens)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	flexGrow, ok := node.Attributes["flexGrow"]
	if !ok {
		t.Error("Missing attribute 'flexGrow'")
	} else {
		if flexGrow.Value != "true" {
			t.Errorf("Boolean attribute should default to 'true', got %q", flexGrow.Value)
		}
	}
}
