// Package parser 实现语法分析，将 Token 流转换为 AST
package parser

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/internal/lexer"
)

// NodeType AST 节点类型
type NodeType int

const (
	// NodeElement represents an element node
	NodeElement NodeType = iota
	// NodeText represents a text node
	NodeText
	// NodeExpression represents an expression node
	NodeExpression
	// NodeFragment represents a fragment node
	NodeFragment
)

// Node AST 节点
type Node struct {
	Type       NodeType
	TagName    string               // Element
	Attributes map[string]AttrValue // Element
	Children   []Node               // Element, Fragment
	Value      string               // Text, Expression
	Spread     []string             // Element: spread attributes
	Line       int
	Column     int
}

// AttrValue 属性值
type AttrValue struct {
	IsExpression bool
	Value        string // 字符串值或表达式内容
}

// Parser 语法分析器
type Parser struct {
	tokens []lexer.Token
	pos    int
	errors []Error
}

// Error 解析错误
type Error struct {
	Message string
	Line    int
	Column  int
}

// New 创建新的解析器
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: make([]Error, 0),
	}
}

// Parse 执行语法分析
func (p *Parser) Parse() (Node, error) {
	node := p.parseRoot()

	if len(p.errors) > 0 {
		return node, fmt.Errorf("parse errors: %v", p.errors)
	}

	return node, nil
}

// parseRoot 解析根节点
func (p *Parser) parseRoot() Node {
	// 跳过前导空白/非标签 Token
	for p.pos < len(p.tokens) {
		tok := p.current()
		if tok.Type == lexer.TokenTagStart || tok.Type == lexer.TokenEOF {
			break
		}
		p.advance()
	}

	// 解析根元素
	if p.current().Type == lexer.TokenTagStart {
		return p.parseElement()
	}

	// 空文档
	return Node{Type: NodeFragment}
}

// parseElement 解析元素
func (p *Parser) parseElement() Node {
	startTok := p.current()

	// 检查是否是 Fragment
	if startTok.Value == "<>" {
		return p.parseFragment()
	}

	node := Node{
		Type:       NodeElement,
		Attributes: make(map[string]AttrValue),
		Spread:     make([]string, 0),
		Children:   make([]Node, 0),
		Line:       startTok.Line,
		Column:     startTok.Column,
	}

	p.advance() // 跳过 <

	// 解析标签名
	if p.current().Type == lexer.TokenTagName {
		node.TagName = p.current().Value
		p.advance()
	}

	// 解析属性
	p.parseAttributes(&node)

	// 检查是否自闭合
	if p.current().Type == lexer.TokenTagSelfClose {
		p.advance()
		return node
	}

	// 跳过 >
	if p.current().Type == lexer.TokenTagEnd {
		p.advance()
	}

	// 解析子节点
	p.parseChildren(&node)

	return node
}

// parseFragment 解析 Fragment
func (p *Parser) parseFragment() Node {
	node := Node{
		Type:     NodeFragment,
		Children: make([]Node, 0),
		Line:     p.current().Line,
		Column:   p.current().Column,
	}

	p.advance() // 跳过 <>

	// 解析子节点
	p.parseChildren(&node)

	return node
}

// parseAttributes 解析属性
func (p *Parser) parseAttributes(node *Node) {
	for p.pos < len(p.tokens) {
		tok := p.current()

		switch tok.Type {
		case lexer.TokenTagEnd, lexer.TokenTagSelfClose:
			return

		case lexer.TokenAttrName:
			attrName := tok.Value
			p.advance()

			// 查找属性值
			value := AttrValue{Value: "true"} // 默认值

			// 跳过可能的 =
			if p.current().Type == lexer.TokenError {
				p.advance()
			}

			// 检查下一个 Token
			switch p.current().Type {
			case lexer.TokenAttrValue:
				value = AttrValue{
					IsExpression: false,
					Value:        p.current().Value,
				}
				p.advance()

			case lexer.TokenExprStart:
				// 表达式值
				p.advance() // 跳过 {
				if p.current().Type == lexer.TokenExprContent {
					value = AttrValue{
						IsExpression: true,
						Value:        p.current().Value,
					}
					p.advance()
				}
				if p.current().Type == lexer.TokenExprEnd {
					p.advance() // 跳过 }
				}
			}

			node.Attributes[attrName] = value

		case lexer.TokenSpread:
			p.advance() // 跳过 ...
			if p.current().Type == lexer.TokenExprContent {
				node.Spread = append(node.Spread, p.current().Value)
				p.advance()
			}

		case lexer.TokenExprStart:
			// 表达式属性 {...props}
			p.advance() // 跳过 {
			if p.current().Type == lexer.TokenSpread {
				p.advance() // 跳过 ...
				if p.current().Type == lexer.TokenExprContent {
					node.Spread = append(node.Spread, p.current().Value)
					p.advance()
				}
			}
			if p.current().Type == lexer.TokenExprEnd {
				p.advance() // 跳过 }
			}

		default:
			p.advance()
		}
	}
}

// parseChildren 解析子节点
func (p *Parser) parseChildren(node *Node) {
	for p.pos < len(p.tokens) {
		tok := p.current()

		switch tok.Type {
		case lexer.TokenTagCloseStart:
			// 闭合标签
			p.advance() // 跳过 </
			p.advance() // 跳过标签名
			if p.current().Type == lexer.TokenTagEnd {
				p.advance() // 跳过 >
			}
			return

		case lexer.TokenTagStart:
			// 子元素
			child := p.parseElement()
			node.Children = append(node.Children, child)

		case lexer.TokenText:
			// 文本节点
			child := Node{
				Type:   NodeText,
				Value:  tok.Value,
				Line:   tok.Line,
				Column: tok.Column,
			}
			node.Children = append(node.Children, child)
			p.advance()

		case lexer.TokenExprStart:
			// 表达式节点
			child := p.parseExpression()
			node.Children = append(node.Children, child)

		case lexer.TokenEOF:
			return

		default:
			p.advance()
		}
	}
}

// parseExpression 解析表达式
func (p *Parser) parseExpression() Node {
	node := Node{
		Type:   NodeExpression,
		Line:   p.current().Line,
		Column: p.current().Column,
	}

	p.advance() // 跳过 {

	// 收集表达式内容
	var content string
	for p.pos < len(p.tokens) {
		tok := p.current()

		switch tok.Type {
		case lexer.TokenExprEnd:
			p.advance()
			node.Value = content
			return node

		case lexer.TokenExprContent:
			content += tok.Value
			p.advance()

		case lexer.TokenTagStart:
			// 嵌套 JSX
			child := p.parseElement()
			content += fmt.Sprintf("<%s>", child.TagName)
			node.Value = content
			// 将嵌套元素作为子节点
			node.Children = append(node.Children, child)

		default:
			p.advance()
		}
	}

	node.Value = content
	return node
}

// current 获取当前 Token
func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

// advance 前进一步
func (p *Parser) advance() {
	p.pos++
}

// addError 添加错误
func (p *Parser) addError(message string, line, column int) {
	p.errors = append(p.errors, Error{
		Message: message,
		Line:    line,
		Column:  column,
	})
}

// String 返回节点的字符串表示
func (n Node) String() string {
	switch n.Type {
	case NodeElement:
		return fmt.Sprintf("Element(%s, attrs=%d, children=%d)", n.TagName, len(n.Attributes), len(n.Children))
	case NodeText:
		return fmt.Sprintf("Text(%q)", n.Value)
	case NodeExpression:
		return fmt.Sprintf("Expr(%q)", n.Value)
	case NodeFragment:
		return fmt.Sprintf("Fragment(children=%d)", len(n.Children))
	default:
		return fmt.Sprintf("Node(%d)", n.Type)
	}
}
