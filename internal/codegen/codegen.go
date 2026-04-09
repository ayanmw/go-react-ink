// Package codegen 实现代码生成，将 AST 转换为 Go 代码
package codegen

import (
	"fmt"
	"strings"
)

// Generator 代码生成器
type Generator struct {
	builder      strings.Builder
	imports      map[string]bool
	indent       int
	componentPkg string
}

// New 创建新的代码生成器
func New(componentPkg string) *Generator {
	return &Generator{
		imports:      make(map[string]bool),
		indent:       0,
		componentPkg: componentPkg,
	}
}

// Generate 生成代码
func (g *Generator) Generate(node Node) string {
	g.builder.Reset()

	// 重置导入
	g.imports = make(map[string]bool)
	g.imports["github.com/ayanmw/go-react-ink/pkg/core"] = true

	// 生成元素代码
	code := g.generateNode(node)

	// 添加导入
	var importLines []string
	for imp := range g.imports {
		importLines = append(importLines, fmt.Sprintf("\t%q", imp))
	}

	// 构建完整代码
	var result strings.Builder
	if len(importLines) > 0 {
		result.WriteString("import (\n")
		for _, line := range importLines {
			result.WriteString(line)
			result.WriteString("\n")
		}
		result.WriteString(")\n\n")
	}
	result.WriteString(code)

	return result.String()
}

// generateNode 生成节点代码
func (g *Generator) generateNode(node Node) string {
	switch node.Type {
	case NodeElement:
		return g.generateElement(node)
	case NodeText:
		return g.generateText(node)
	case NodeExpression:
		return g.generateExpression(node)
	case NodeFragment:
		return g.generateFragment(node)
	default:
		return ""
	}
}

// generateElement 生成元素代码
func (g *Generator) generateElement(node Node) string {
	var code strings.Builder

	// 开始调用
	code.WriteString("core.CreateElement(\n")
	g.indent++
	code.WriteString(g.indentStr())

	// 组件名
	if g.componentPkg != "" {
		code.WriteString(fmt.Sprintf("%s.%s", g.componentPkg, node.TagName))
	} else {
		code.WriteString(node.TagName)
	}
	code.WriteString(",\n")

	// 属性
	code.WriteString(g.indentStr())
	code.WriteString("core.Props{\n")
	g.indent++
	for name, value := range node.Attributes {
		code.WriteString(g.indentStr())
		code.WriteString(fmt.Sprintf("%q: ", name))
		if value.IsExpression {
			code.WriteString(value.Value)
		} else {
			code.WriteString(fmt.Sprintf("%q", value.Value))
		}
		code.WriteString(",\n")
	}
	g.indent--
	code.WriteString(g.indentStr())
	code.WriteString("},\n")

	// Spread 属性
	for _, spread := range node.Spread {
		code.WriteString(g.indentStr())
		code.WriteString(fmt.Sprintf("core.Spread(%s),\n", spread))
	}

	// 子节点
	if len(node.Children) > 0 {
		code.WriteString(g.indentStr())
		code.WriteString("[]core.Element{\n")
		g.indent++
		for _, child := range node.Children {
			code.WriteString(g.indentStr())
			childCode := g.generateNode(child)
			code.WriteString(childCode)
			code.WriteString(",\n")
		}
		g.indent--
		code.WriteString(g.indentStr())
		code.WriteString("},\n")
	} else {
		code.WriteString(g.indentStr())
		code.WriteString("nil,\n")
	}

	g.indent--
	code.WriteString(g.indentStr())
	code.WriteString(")")

	return code.String()
}

// generateText 生成文本代码
func (g *Generator) generateText(node Node) string {
	return fmt.Sprintf("core.CreateElement(core.Text, core.Props{\"children\": %q}, nil)", node.Value)
}

// generateExpression 生成表达式代码
func (g *Generator) generateExpression(node Node) string {
	// 表达式直接返回其值
	return node.Value
}

// generateFragment 生成 Fragment 代码
func (g *Generator) generateFragment(node Node) string {
	if len(node.Children) == 0 {
		return "nil"
	}

	if len(node.Children) == 1 {
		return g.generateNode(node.Children[0])
	}

	// 多个子节点，返回数组
	var code strings.Builder
	code.WriteString("[]core.Element{\n")
	g.indent++
	for _, child := range node.Children {
		code.WriteString(g.indentStr())
		code.WriteString(g.generateNode(child))
		code.WriteString(",\n")
	}
	g.indent--
	code.WriteString(g.indentStr())
	code.WriteString("}")

	return code.String()
}

// indentStr 生成缩进字符串
func (g *Generator) indentStr() string {
	return strings.Repeat("\t", g.indent)
}

// GenerateFunction 生成组件函数
func GenerateFunction(name string, node Node, componentPkg string) string {
	gen := New(componentPkg)
	elementCode := gen.generateNode(node)

	var code strings.Builder
	code.WriteString(fmt.Sprintf("func %s() core.Element {\n", name))
	code.WriteString("\treturn ")
	code.WriteString(elementCode)
	code.WriteString("\n")
	code.WriteString("}")

	return code.String()
}

// GenerateReturn 生成 return 语句
func GenerateReturn(node Node, componentPkg string) string {
	gen := New(componentPkg)
	return "return " + gen.generateNode(node)
}

// GenerateAssignment 生成赋值语句
func GenerateAssignment(varName string, node Node, componentPkg string) string {
	gen := New(componentPkg)
	return varName + " := " + gen.generateNode(node)
}

// GenerateNode 生成节点代码（公开方法）
func (g *Generator) GenerateNode(node Node) string {
	return g.generateNode(node)
}
