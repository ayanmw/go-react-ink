// Package lexer 实现词法分析，将 JSX 文本转换为 Token 流
package lexer

import (
	"fmt"
	"unicode"
)

// TokenType Token 类型
type TokenType int

const (
	// TokenEOF marks end of input
	TokenEOF TokenType = iota
	// TokenError marks an error token
	TokenError
	// TokenTagStart marks tag start (<)
	TokenTagStart
	// TokenTagEnd marks tag end (>)
	TokenTagEnd
	// TokenTagSelfClose marks self-closing tag (/>)
	TokenTagSelfClose
	// TokenTagCloseStart marks closing tag start (</)
	TokenTagCloseStart
	// TokenTagName marks tag name (Box, Text, Fragment)
	TokenTagName

	// TokenAttrName marks attribute name
	TokenAttrName
	// TokenAttrValue marks attribute value
	TokenAttrValue
	// TokenExprStart marks expression start ({)
	TokenExprStart
	// TokenExprEnd marks expression end (})
	TokenExprEnd
	// TokenExprContent marks expression content
	TokenExprContent
	// TokenSpread marks spread operator (...props)
	TokenSpread

	// TokenText marks text content
	TokenText

	// TokenComment marks comment ({/* comment */})
	TokenComment
)

// Token 词法单元
type Token struct {
	Type     TokenType
	Value    string
	Line     int
	Column   int
	Position int // 字符位置
}

// Lexer 词法分析器
type Lexer struct {
	source []rune
	pos    int
	line   int
	column int
	tokens []Token
}

// New 创建新的词法分析器
func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
		pos:    0,
		line:   1,
		column: 1,
		tokens: make([]Token, 0),
	}
}

// Lex 执行词法分析
func (l *Lexer) Lex() []Token {
	for l.pos < len(l.source) {
		l.scanToken()
	}

	// 添加 EOF
	l.tokens = append(l.tokens, Token{
		Type:   TokenEOF,
		Line:   l.line,
		Column: l.column,
	})

	return l.tokens
}

// scanToken 扫描下一个 Token
func (l *Lexer) scanToken() {
	ch := l.peek()

	switch ch {
	case '<':
		l.scanTagStart()
	case '>':
		l.emit(TokenTagEnd, ">")
		l.advance()
	case '{':
		l.scanExpression()
	case '"', '\'':
		l.scanAttrValue(ch)
	default:
		if unicode.IsSpace(ch) {
			l.scanText()
		} else if unicode.IsLetter(ch) || ch == '_' {
			l.scanAttrOrText()
		} else {
			l.advance()
		}
	}
}

// scanTagStart 扫描标签开始
func (l *Lexer) scanTagStart() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	l.advance() // 跳过 <

	next := l.peek()

	switch {
	case next == '/':
		// 闭合标签 </TagName>
		l.advance() // 跳过 /
		l.emitAt(TokenTagCloseStart, "</", startLine, startCol, startPos)
		l.scanTagName()
		// 跳过空白
		l.skipWhitespace()
		// 期待 >
		if l.peek() == '>' {
			l.emit(TokenTagEnd, ">")
			l.advance()
		}

	case next == '>':
		// Fragment <>
		l.advance() // 跳过 >
		l.emitAt(TokenTagStart, "<>", startLine, startCol, startPos)

	case unicode.IsUpper(next):
		// 开始标签 <TagName>
		l.emitAt(TokenTagStart, "<", startLine, startCol, startPos)
		l.scanTagName()

		// 扫描属性
		l.scanAttrs()

	default:
		// 其他情况，当作文本
		l.emitAt(TokenText, "<", startLine, startCol, startPos)
	}
}

// scanTagName 扫描标签名
func (l *Lexer) scanTagName() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	for l.pos < len(l.source) {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	if l.pos > startPos {
		value := string(l.source[startPos:l.pos])
		l.emitAt(TokenTagName, value, startLine, startCol, startPos)
	}
}

// scanAttrs 扫描属性直到标签结束
func (l *Lexer) scanAttrs() {
	for l.pos < len(l.source) {
		l.skipWhitespace()

		ch := l.peek()

		switch {
		case ch == '>':
			l.emit(TokenTagEnd, ">")
			l.advance()
			return

		case ch == '/' && l.peekNext() == '>':
			l.emit(TokenTagSelfClose, "/>")
			l.advance()
			l.advance()
			return

		case ch == '{':
			l.scanExpression()

		case ch == '.' && l.peekNext() == '.' && l.pos+2 < len(l.source) && l.source[l.pos+2] == '.':
			// Spread ...props
			l.scanSpread()

		case unicode.IsLetter(ch) || ch == '_' || ch == ':':
			l.scanAttrName()

		default:
			l.advance()
		}
	}
}

// scanAttrName 扫描属性名
func (l *Lexer) scanAttrName() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	for l.pos < len(l.source) {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' || ch == ':' {
			l.advance()
		} else {
			break
		}
	}

	if l.pos > startPos {
		value := string(l.source[startPos:l.pos])
		l.emitAt(TokenAttrName, value, startLine, startCol, startPos)
	}

	// 跳过空白和 =
	l.skipWhitespace()
	if l.peek() == '=' {
		l.advance()
		l.skipWhitespace()
		l.scanAttrValueOrExpr()
	}
}

// scanAttrValueOrExpr 扫描属性值或表达式
func (l *Lexer) scanAttrValueOrExpr() {
	ch := l.peek()
	if ch == '{' {
		l.scanExpression()
	} else if ch == '"' || ch == '\'' {
		l.scanAttrValue(ch)
	}
}

// scanSpread 扫描展开属性 ...props
func (l *Lexer) scanSpread() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	l.advance() // .
	l.advance() // .
	l.advance() // .
	l.emitAt(TokenSpread, "...", startLine, startCol, startPos)

	// 扫描变量名
	l.skipWhitespace()
	varStart := l.pos
	for l.pos < len(l.source) {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.' {
			l.advance()
		} else {
			break
		}
	}
	if l.pos > varStart {
		value := string(l.source[varStart:l.pos])
		l.emit(TokenExprContent, value)
	}
}

// scanExpression 扫描 JSX 表达式 {...}
func (l *Lexer) scanExpression() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	l.advance() // 跳过 {

	// 检查是否是 spread ...props
	if l.peek() == '.' && l.peekNext() == '.' && l.pos+2 < len(l.source) && l.source[l.pos+2] == '.' {
		l.emitAt(TokenExprStart, "{", startLine, startCol, startPos)
		l.advance() // .
		l.advance() // .
		l.advance() // .
		l.emit(TokenSpread, "...")

		// 扫描变量名
		l.skipWhitespace()
		varStart := l.pos
		for l.pos < len(l.source) {
			ch := l.peek()
			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.' {
				l.advance()
			} else {
				break
			}
		}
		if l.pos > varStart {
			value := string(l.source[varStart:l.pos])
			l.emit(TokenExprContent, value)
		}

		// 期待 }
		l.skipWhitespace()
		if l.peek() == '}' {
			l.emit(TokenExprEnd, "}")
			l.advance()
		}
		return
	}

	l.emitAt(TokenExprStart, "{", startLine, startCol, startPos)

	// 扫描表达式内容
	depth := 1
	exprStart := l.pos

	for l.pos < len(l.source) && depth > 0 {
		ch := l.peek()
		switch ch {
		case '{':
			depth++
			l.advance()
		case '}':
			depth--
			if depth == 0 {
				// 表达式结束
				if l.pos > exprStart {
					value := string(l.source[exprStart:l.pos])
					l.emit(TokenExprContent, value)
				}
				l.emit(TokenExprEnd, "}")
				l.advance()
			} else {
				l.advance()
			}
		case '"', '\'', '`':
			l.skipString(ch)
		case '<':
			// 嵌套 JSX
			if l.peekNext() != '!' {
				l.scanTagStart()
			} else {
				l.advance()
			}
		default:
			l.advance()
		}
	}
}

// scanAttrValue 扫描属性值 "..."
func (l *Lexer) scanAttrValue(quote rune) {
	l.advance() // 跳过开始引号

	startLine := l.line
	startCol := l.column
	startPos := l.pos

	valueStart := l.pos
	for l.pos < len(l.source) {
		ch := l.peek()
		if ch == '\\' {
			l.advance()
			if l.pos < len(l.source) {
				l.advance()
			}
			continue
		}
		if ch == quote {
			break
		}
		l.advance()
	}

	value := string(l.source[valueStart:l.pos])
	l.emitAt(TokenAttrValue, value, startLine, startCol, startPos)

	if l.peek() == quote {
		l.advance() // 跳过结束引号
	}
}

// scanAttrOrText 扫描属性或文本
func (l *Lexer) scanAttrOrText() {
	// 检查是否是属性 (后面跟着 =)
	savedPos := l.pos

	// 扫描名称
	for l.pos < len(l.source) {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' || ch == ':' {
			l.advance()
		} else {
			break
		}
	}

	// 检查后面是否是 =
	l.skipWhitespace()
	if l.peek() == '=' {
		// 是属性，发射属性名
		l.pos = savedPos
		l.scanAttrName()
		return
	}

	// 不是属性，当作文本
	l.pos = savedPos
	l.scanText()
}

// scanText 扫描文本内容
func (l *Lexer) scanText() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	for l.pos < len(l.source) {
		ch := l.peek()

		// 遇到特殊字符停止
		if ch == '<' || ch == '{' || ch == '>' {
			break
		}

		l.advance()
	}

	if l.pos > startPos {
		value := string(l.source[startPos:l.pos])
		// 只有非纯空白才发射
		if !isAllWhitespace(value) {
			l.emitAt(TokenText, value, startLine, startCol, startPos)
		}
	}
}

// skipWhitespace 跳过空白
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.source) && unicode.IsSpace(l.peek()) {
		l.advance()
	}
}

// skipString 跳过字符串
func (l *Lexer) skipString(quote rune) {
	l.advance() // 开始引号
	for l.pos < len(l.source) {
		ch := l.peek()
		if ch == '\\' {
			l.advance()
			if l.pos < len(l.source) {
				l.advance()
			}
			continue
		}
		if ch == quote {
			l.advance()
			return
		}
		l.advance()
	}
}

// emit 发射 Token
func (l *Lexer) emit(typ TokenType, value string) {
	l.tokens = append(l.tokens, Token{
		Type:     typ,
		Value:    value,
		Line:     l.line,
		Column:   l.column,
		Position: l.pos,
	})
}

// emitAt 在指定位置发射 Token
func (l *Lexer) emitAt(typ TokenType, value string, line, col, pos int) {
	l.tokens = append(l.tokens, Token{
		Type:     typ,
		Value:    value,
		Line:     line,
		Column:   col,
		Position: pos,
	})
}

// advance 前进一步
func (l *Lexer) advance() {
	if l.pos >= len(l.source) {
		return
	}
	if l.source[l.pos] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.pos++
}

// peek 查看当前字符
func (l *Lexer) peek() rune {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

// peekNext 查看下一个字符
func (l *Lexer) peekNext() rune {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.pos+1]
}

// isAllWhitespace 检查是否全是空白
func isAllWhitespace(s string) bool {
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			return false
		}
	}
	return true
}

// String 返回 Token 的字符串表示
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) @ %d:%d", t.Type.String(), t.Value, t.Line, t.Column)
}

// String 返回 TokenType 的字符串表示
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "Error"
	case TokenTagStart:
		return "TagStart"
	case TokenTagEnd:
		return "TagEnd"
	case TokenTagSelfClose:
		return "TagSelfClose"
	case TokenTagCloseStart:
		return "TagCloseStart"
	case TokenTagName:
		return "TagName"
	case TokenAttrName:
		return "AttrName"
	case TokenAttrValue:
		return "AttrValue"
	case TokenExprStart:
		return "ExprStart"
	case TokenExprEnd:
		return "ExprEnd"
	case TokenExprContent:
		return "ExprContent"
	case TokenSpread:
		return "Spread"
	case TokenText:
		return "Text"
	case TokenComment:
		return "Comment"
	default:
		return fmt.Sprintf("Token(%d)", t)
	}
}
