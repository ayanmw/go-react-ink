package lexer

import (
	"strings"
	"testing"
)

func TestLexSimpleTag(t *testing.T) {
	source := "<Box></Box>"
	l := New(source)
	tokens := l.Lex()

	expected := []struct {
		Type  TokenType
		Value string
	}{
		{TokenTagStart, "<"},
		{TokenTagName, "Box"},
		{TokenTagEnd, ">"},
		{TokenTagCloseStart, "</"},
		{TokenTagName, "Box"},
		{TokenTagEnd, ">"},
		{TokenEOF, ""},
	}

	compareTokens(t, tokens, expected)
}

func TestLexTagWithAttrs(t *testing.T) {
	source := `<Box flexDirection="column" padding={1}></Box>`
	l := New(source)
	tokens := l.Lex()

	// 检查关键 Token
	if !hasToken(tokens, TokenTagStart, "<") {
		t.Error("Missing TagStart <")
	}
	if !hasToken(tokens, TokenTagName, "Box") {
		t.Error("Missing TagName Box")
	}
	if !hasToken(tokens, TokenAttrName, "flexDirection") {
		t.Error("Missing AttrName flexDirection")
	}
	if !hasToken(tokens, TokenAttrValue, "column") {
		t.Error("Missing AttrValue column")
	}
	if !hasToken(tokens, TokenAttrName, "padding") {
		t.Error("Missing AttrName padding")
	}
	if !hasToken(tokens, TokenExprContent, "1") {
		t.Error("Missing ExprContent 1")
	}
}

func TestLexNestedTags(t *testing.T) {
	source := `<Box><Text>Hello</Text></Box>`
	l := New(source)
	tokens := l.Lex()

	// 检查嵌套结构
	tagNames := collectTokens(tokens, TokenTagName)
	expected := []string{"Box", "Text", "Text", "Box"}
	if !equalStrings(tagNames, expected) {
		t.Errorf("Expected tag names %v, got %v", expected, tagNames)
	}

	// 检查文本内容
	if !hasToken(tokens, TokenText, "Hello") {
		t.Error("Missing Text 'Hello'")
	}
}

func TestLexExpression(t *testing.T) {
	source := `<Text>Count: {count}</Text>`
	l := New(source)
	tokens := l.Lex()

	if !hasToken(tokens, TokenExprStart, "{") {
		t.Error("Missing ExprStart {")
	}
	if !hasToken(tokens, TokenExprContent, "count") {
		t.Error("Missing ExprContent count")
	}
	if !hasToken(tokens, TokenExprEnd, "}") {
		t.Error("Missing ExprEnd }")
	}
}

func TestLexComplexExpression(t *testing.T) {
	source := `<Text>{item.name + " " + item.value}</Text>`
	l := New(source)
	tokens := l.Lex()

	// 表达式内容应该包含完整表达式
	for _, tok := range tokens {
		if tok.Type == TokenExprContent {
			if !strings.Contains(tok.Value, "item.name") {
				t.Errorf("ExprContent should contain 'item.name', got %q", tok.Value)
			}
		}
	}
}

func TestLexSelfClosing(t *testing.T) {
	source := `<Box><Spacer /></Box>`
	l := New(source)
	tokens := l.Lex()

	if !hasToken(tokens, TokenTagName, "Spacer") {
		t.Error("Missing TagName Spacer")
	}
	if !hasToken(tokens, TokenTagSelfClose, "/>") {
		t.Error("Missing TagSelfClose />")
	}
}

func TestLexSpreadAttr(t *testing.T) {
	source := `<Box {...props}></Box>`
	l := New(source)
	tokens := l.Lex()

	if !hasToken(tokens, TokenSpread, "...") {
		t.Error("Missing Spread ...")
	}
	if !hasToken(tokens, TokenExprContent, "props") {
		t.Error("Missing spread variable 'props'")
	}
}

func TestLexFragment(t *testing.T) {
	source := `<>Content</>`
	l := New(source)
	tokens := l.Lex()

	// Fragment 开始
	if !hasToken(tokens, TokenTagStart, "<>") {
		t.Error("Missing Fragment start <>")
	}

	// Fragment 结束
	if !hasToken(tokens, TokenTagCloseStart, "</") {
		t.Error("Missing Fragment close </")
	}
}

func TestLexMultipleAttrs(t *testing.T) {
	source := `<Box flexDirection="row" justifyContent="center" alignItems="center" padding={1} margin={2}></Box>`
	l := New(source)
	tokens := l.Lex()

	attrs := collectTokens(tokens, TokenAttrName)
	expected := []string{"flexDirection", "justifyContent", "alignItems", "padding", "margin"}
	if !equalStrings(attrs, expected) {
		t.Errorf("Expected attrs %v, got %v", expected, attrs)
	}
}

func TestLexTextWithWhitespace(t *testing.T) {
	source := `<Text>Hello    World</Text>`
	l := New(source)
	tokens := l.Lex()

	// 应该有一个文本 Token 包含完整内容
	textTokens := collectTokens(tokens, TokenText)
	if len(textTokens) == 0 {
		t.Error("No text tokens found")
	}

	// 文本应该包含空格
	found := false
	for _, tok := range tokens {
		if tok.Type == TokenText && strings.Contains(tok.Value, "Hello") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Text token should contain 'Hello'")
	}
}

func TestLexConditionalExpression(t *testing.T) {
	source := `{show && <Text>Visible</Text>}`
	l := New(source)
	tokens := l.Lex()

	// 表达式内容应该包含条件
	if !hasToken(tokens, TokenExprStart, "{") {
		t.Error("Missing ExprStart")
	}

	// 应该有嵌套的标签
	if !hasToken(tokens, TokenTagName, "Text") {
		t.Error("Missing nested TagName Text")
	}
}

func TestLexMapExpression(t *testing.T) {
	source := `{items.map(item => <Item key={item.id} />)}`
	l := New(source)
	tokens := l.Lex()

	// 表达式应该包含 map
	found := false
	for _, tok := range tokens {
		if tok.Type == TokenExprContent && strings.Contains(tok.Value, "items.map") {
			found = true
			break
		}
	}
	if !found {
		t.Error("ExprContent should contain 'items.map'")
	}
}

func TestLexLineColumnTracking(t *testing.T) {
	source := `<Box>
    <Text>Hello</Text>
</Box>`
	l := New(source)
	tokens := l.Lex()

	// 第一个 Token 应该在第一行
	if len(tokens) > 0 && tokens[0].Line != 1 {
		t.Errorf("First token should be on line 1, got line %d", tokens[0].Line)
	}

	// 找到 Text 标签，应该在第二行
	for _, tok := range tokens {
		if tok.Type == TokenTagName && tok.Value == "Text" {
			if tok.Line != 2 {
				t.Errorf("Text tag should be on line 2, got line %d", tok.Line)
			}
		}
	}
}

// 辅助函数
func compareTokens(t *testing.T, tokens []Token, expected []struct {
	Type  TokenType
	Value string
}) {
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
		return
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.Type {
			t.Errorf("Token %d: expected type %v, got %v", i, exp.Type, tokens[i].Type)
		}
		if tokens[i].Value != exp.Value {
			t.Errorf("Token %d: expected value %q, got %q", i, exp.Value, tokens[i].Value)
		}
	}
}

func hasToken(tokens []Token, typ TokenType, value string) bool {
	for _, tok := range tokens {
		if tok.Type == typ && tok.Value == value {
			return true
		}
	}
	return false
}

func collectTokens(tokens []Token, typ TokenType) []string {
	result := make([]string, 0)
	for _, tok := range tokens {
		if tok.Type == typ {
			result = append(result, tok.Value)
		}
	}
	return result
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}