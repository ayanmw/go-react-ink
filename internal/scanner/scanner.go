// Package scanner 实现源码扫描，识别 .gox 文件中的 JSX 区域
package scanner

import (
	"fmt"
	"strings"
	"unicode"
)

// SegmentType 代码段类型
type SegmentType int

const (
	// SegmentGoCode 纯 Go 代码
	SegmentGoCode SegmentType = iota
	// SegmentJSX JSX 区域
	SegmentJSX
)

// Segment 代码段
type Segment struct {
	Type    SegmentType
	Content string
	Line    int // 起始行号
	Column  int // 起始列号
}

// Scanner 源码扫描器
type Scanner struct {
	source   []rune
	pos      int
	line     int
	column   int
	segments []Segment
}

// New 创建新的扫描器
func New(source string) *Scanner {
	return &Scanner{
		source:   []rune(source),
		pos:      0,
		line:     1,
		column:   1,
		segments: make([]Segment, 0),
	}
}

// Scan 执行扫描，返回代码段列表
func (s *Scanner) Scan() ([]Segment, error) {
	for s.pos < len(s.source) {
		// 检测 JSX 起始标记
		if s.isJSXStart() {
			if err := s.scanJSX(); err != nil {
				return nil, err
			}
		} else {
			s.scanGoCode()
		}
	}

	return s.segments, nil
}

// isJSXStart 检测是否为 JSX 起始位置
func (s *Scanner) isJSXStart() bool {
	if s.peek() != '<' {
		return false
	}

	// 检查 < 后是否为大写字母 (组件名如 Box, Text) 或 Fragment <>
	nextPos := s.pos + 1
	if nextPos >= len(s.source) {
		return false
	}

	next := s.source[nextPos]
	// 支持组件名 (大写字母) 或 Fragment (>)
	if !unicode.IsUpper(next) && next != '>' {
		return false
	}

	// 检查前面是否是有效的 JSX 上下文
	return s.isValidJSXContext()
}

// isValidJSXContext 检查当前上下文是否允许 JSX
func (s *Scanner) isValidJSXContext() bool {
	// 向前查找有效的 JSX 上下文
	pos := s.pos - 1

	// 跳过空白
	for pos >= 0 && isWhitespace(s.source[pos]) {
		pos--
	}

	if pos < 0 {
		// 文件开头，允许顶层 JSX
		return true
	}

	ch := s.source[pos]

	// 检查各种有效的 JSX 上下文
	switch ch {
	case '(', '{', ',', ':':
		return true
	case '=':
		return true
	case 'n':
		// 检查是否是 return
		if pos >= 5 {
			word := string(s.source[pos-5 : pos+1])
			if word == "return" {
				return true
			}
		}
	}

	return false
}

// scanJSX 扫描 JSX 区域
func (s *Scanner) scanJSX() error {
	startLine := s.line
	startColumn := s.column
	start := s.pos

	// 使用栈来跟踪标签
	tagStack := make([]string, 0)

	for s.pos < len(s.source) {
		ch := s.peek()

		switch ch {
		case '<':
			// 检查是什么类型的标签
			if s.peekNext() == '/' {
				// 闭合标签 </TagName>
				s.advance() // <
				s.advance() // /
				tagName := s.scanTagName()
				if len(tagStack) > 0 && tagStack[len(tagStack)-1] == tagName {
					tagStack = tagStack[:len(tagStack)-1]
				}
				// 跳过到 >
				for s.pos < len(s.source) && s.peek() != '>' {
					s.advance()
				}
				if s.peek() == '>' {
					s.advance()
				}
			} else if s.peekNext() == '>' {
				// Fragment <>
				s.advance() // <
				s.advance() // >
				// Fragment 不加入标签栈，用特殊标记
				tagStack = append(tagStack, "<>")
			} else if s.peekNext() == '!' && s.pos+2 < len(s.source) && s.source[s.pos+2] == '-' {
				// HTML 注释 <!--
				s.skipHTMLComment()
			} else if unicode.IsUpper(s.peekNext()) {
				// 开始标签 <TagName>
				s.advance() // <
				tagName := s.scanTagName()
				tagStack = append(tagStack, tagName)
			} else {
				s.advance()
			}

		case '/':
			// 自闭合标签 />
			if s.peekNext() == '>' {
				s.advance() // /
				s.advance() // >
				if len(tagStack) > 0 {
					tagStack = tagStack[:len(tagStack)-1]
				}
			} else {
				s.advance()
			}

		case '>':
			// 普通结束
			s.advance()

		case '{':
			// JSX 表达式
			s.advance()
			s.skipExpression()

		case '"', '\'', '`':
			// 字符串
			s.skipString(ch)

		default:
			s.advance()
		}

		// 当标签栈为空时，JSX 区域结束
		if len(tagStack) == 0 && s.pos > start {
			break
		}
	}

	content := string(s.source[start:s.pos])
	s.segments = append(s.segments, Segment{
		Type:    SegmentJSX,
		Content: content,
		Line:    startLine,
		Column:  startColumn,
	})

	return nil
}

// scanTagName 扫描标签名
func (s *Scanner) scanTagName() string {
	start := s.pos
	for s.pos < len(s.source) {
		ch := s.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' {
			s.advance()
		} else {
			break
		}
	}
	return string(s.source[start:s.pos])
}

// skipHTMLComment 跳过 HTML 注释
func (s *Scanner) skipHTMLComment() {
	s.advance() // <
	s.advance() // !
	s.advance() // -
	s.advance() // -

	for s.pos < len(s.source)-2 {
		if s.peek() == '-' && s.peekNext() == '-' && s.source[s.pos+2] == '>' {
			s.advance()
			s.advance()
			s.advance()
			return
		}
		s.advance()
	}
}

// scanGoCode 扫描 Go 代码直到遇到 JSX
func (s *Scanner) scanGoCode() {
	startLine := s.line
	startColumn := s.column
	start := s.pos

	for s.pos < len(s.source) {
		// 检查是否遇到 JSX
		if s.isJSXStart() {
			break
		}

		// 处理字符串
		ch := s.peek()
		if ch == '"' || ch == '\'' || ch == '`' {
			s.skipString(ch)
			continue
		}

		// 处理注释
		if ch == '/' {
			if s.peekNext() == '/' {
				s.skipLineComment()
				continue
			}
			if s.peekNext() == '*' {
				s.skipBlockComment()
				continue
			}
		}

		s.advance()
	}

	if s.pos > start {
		content := string(s.source[start:s.pos])
		s.segments = append(s.segments, Segment{
			Type:    SegmentGoCode,
			Content: content,
			Line:    startLine,
			Column:  startColumn,
		})
	}
}

// skipExpression 跳过 JSX 表达式 {...}
func (s *Scanner) skipExpression() {
	depth := 1
	for s.pos < len(s.source) && depth > 0 {
		ch := s.peek()
		switch ch {
		case '{':
			depth++
			s.advance()
		case '}':
			depth--
			s.advance()
		case '"', '\'', '`':
			s.skipString(ch)
		default:
			s.advance()
		}
	}
}

// skipString 跳过字符串字面量
func (s *Scanner) skipString(quote rune) {
	s.advance() // 开始引号

	for s.pos < len(s.source) {
		ch := s.peek()
		if ch == '\\' {
			s.advance() // 跳过转义
			if s.pos < len(s.source) {
				s.advance()
			}
			continue
		}
		if ch == quote {
			s.advance() // 结束引号
			return
		}
		s.advance()
	}
}

// skipLineComment 跳过行注释
func (s *Scanner) skipLineComment() {
	s.advance() // /
	s.advance() // /

	for s.pos < len(s.source) && s.peek() != '\n' {
		s.advance()
	}
}

// skipBlockComment 跳过块注释
func (s *Scanner) skipBlockComment() {
	s.advance() // /
	s.advance() // *

	for s.pos < len(s.source) {
		if s.peek() == '*' && s.peekNext() == '/' {
			s.advance() // *
			s.advance() // /
			return
		}
		s.advance()
	}
}

// advance 前进一步
func (s *Scanner) advance() {
	if s.pos >= len(s.source) {
		return
	}

	if s.source[s.pos] == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}
	s.pos++
}

// peek 查看当前字符
func (s *Scanner) peek() rune {
	if s.pos >= len(s.source) {
		return 0
	}
	return s.source[s.pos]
}

// peekNext 查看下一个字符
func (s *Scanner) peekNext() rune {
	if s.pos+1 >= len(s.source) {
		return 0
	}
	return s.source[s.pos+1]
}

// isWhitespace 判断是否为空白字符
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// String 返回 Segment 的字符串表示
func (s Segment) String() string {
	var typ string
	switch s.Type {
	case SegmentGoCode:
		typ = "GoCode"
	case SegmentJSX:
		typ = "JSX"
	}
	return fmt.Sprintf("[%s @ %d:%d] %q", typ, s.Line, s.Column, truncate(s.Content, 50))
}

func truncate(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max] + "..."
}

// 去除导入冲突
var _ = strings.TrimSpace("")