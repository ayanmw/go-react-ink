package compiler

import (
	"strings"
	"testing"

	"github.com/anmingwei/go-ink/internal/codegen"
	"github.com/anmingwei/go-ink/internal/parser"
)

func TestCompileSimpleElement(t *testing.T) {
	c := New("ink")
	output, err := c.Compile("test.gox", `<Box></Box>`)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证输出包含 CreateElement
	if !strings.Contains(string(output), "CreateElement") {
		t.Error("Output should contain CreateElement")
	}
}

func TestCompileElementWithChildren(t *testing.T) {
	c := New("ink")
	source := `<Box><Text>Hello</Text></Box>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证输出包含嵌套结构
	content := string(output)
	if !strings.Contains(content, "ink.Box") {
		t.Error("Output should contain 'ink.Box'")
	}
	if !strings.Contains(content, "ink.Text") {
		t.Error("Output should contain 'ink.Text'")
	}
	if !strings.Contains(content, "Hello") {
		t.Error("Output should contain 'Hello'")
	}
}

func TestCompileWithGoCode(t *testing.T) {
	c := New("ink")
	source := `package main

func App() Element {
	return <Box><Text>Hello</Text></Box>
}
`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证 Go 代码保留
	content := string(output)
	if !strings.Contains(content, "package main") {
		t.Error("Output should contain 'package main'")
	}
	if !strings.Contains(content, "func App()") {
		t.Error("Output should contain 'func App()'")
	}
}

func TestCompileExpression(t *testing.T) {
	c := New("ink")
	source := `<Text>Count: {count}</Text>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证表达式生成
	content := string(output)
	if !strings.Contains(content, "count") {
		t.Error("Output should contain expression 'count'")
	}
}

func TestCompileAttributes(t *testing.T) {
	c := New("ink")
	source := `<Box flexDirection="column" padding={1}></Box>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证属性生成
	content := string(output)
	if !strings.Contains(content, `"flexDirection"`) {
		t.Error("Output should contain 'flexDirection' attribute")
	}
	if !strings.Contains(content, `"column"`) {
		t.Error("Output should contain 'column' value")
	}
}

func TestCompileSelfClosing(t *testing.T) {
	c := New("ink")
	source := `<Box><Spacer /></Box>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证自闭合标签
	content := string(output)
	if !strings.Contains(content, "ink.Spacer") {
		t.Error("Output should contain 'ink.Spacer'")
	}
}

func TestCompileSpread(t *testing.T) {
	c := New("ink")
	source := `<Box {...props}></Box>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证 spread 生成
	content := string(output)
	if !strings.Contains(content, "Spread(props)") {
		t.Error("Output should contain 'Spread(props)'")
	}
}

func TestCompileFragment(t *testing.T) {
	c := New("ink")
	source := `<><Text>A</Text><Text>B</Text></>`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证 fragment 处理
	content := string(output)
	if !strings.Contains(content, "A") || !strings.Contains(content, "B") {
		t.Error("Output should contain children content")
	}
}

func TestCompileMultipleJSX(t *testing.T) {
	c := New("ink")
	source := `func App() Element {
	a := <Text>A</Text>
	b := <Text>B</Text>
	return <Box>{a}{b}</Box>
}
`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证多个 JSX 编译
	content := string(output)
	if !strings.Contains(content, "a :=") {
		t.Error("Output should contain variable assignment 'a :='")
	}
	if !strings.Contains(content, "b :=") {
		t.Error("Output should contain variable assignment 'b :='")
	}
}

func TestCompileImports(t *testing.T) {
	c := New("ink")
	source := `package main

func App() Element {
	return <Box></Box>
}
`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// 验证导入生成
	content := string(output)
	if !strings.Contains(content, "import") {
		t.Error("Output should contain import statement")
	}
	if !strings.Contains(content, "github.com/anmingwei/go-ink/pkg/core") {
		t.Error("Output should contain core package import")
	}
}

func TestConvertAST(t *testing.T) {
	// 测试 AST 转换
	parserNode := parser.Node{
		Type:     parser.NodeElement,
		TagName:  "Box",
		Children: []parser.Node{
			{Type: parser.NodeText, Value: "Hello"},
		},
	}

	c := New("ink")
	cgNode := c.convertAST(parserNode)

	// 验证转换
	if cgNode.Type != codegen.NodeElement {
		t.Error("Type should be NodeElement")
	}
	if cgNode.TagName != "Box" {
		t.Error("TagName should be 'Box'")
	}
	if len(cgNode.Children) != 1 {
		t.Error("Should have 1 child")
	}
	if cgNode.Children[0].Value != "Hello" {
		t.Error("Child value should be 'Hello'")
	}
}

func TestFixImportsNoPackage(t *testing.T) {
	c := New("ink")
	// 需要包含 core. 或 ink. 才会触发导入添加
	source := []byte("func main() { core.CreateElement(nil, nil, nil) }")
	output := c.fixImports(source)

	// 没有 package 声明，应该添加
	if !strings.Contains(string(output), "package main") {
		t.Error("Should add package declaration")
	}
}

func TestFixImportsWithExisting(t *testing.T) {
	c := New("ink")
	source := []byte(`package main

import "fmt"

func main() {}
`)
	output := c.fixImports(source)

	// 有现有 import，应该保持
	content := string(output)
	if !strings.Contains(content, "import") {
		t.Error("Should contain import")
	}
	if !strings.Contains(content, "fmt") {
		t.Error("Should preserve existing import 'fmt'")
	}
}