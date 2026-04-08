package scanner

import (
	"strings"
	"testing"
)

func TestScanGoCodeOnly(t *testing.T) {
	source := `package main

func main() {
    fmt.Println("Hello")
}
`
	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(segments) != 1 {
		t.Fatalf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Type != SegmentGoCode {
		t.Errorf("Expected GoCode segment, got %v", segments[0].Type)
	}
}

func TestScanJSXInReturn(t *testing.T) {
	source := `return (
    <Box>
        <Text>Hello</Text>
    </Box>
)`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 应该有 3 段: GoCode + JSX + GoCode
	var jsxFound bool
	for _, seg := range segments {
		if seg.Type == SegmentJSX {
			jsxFound = true
			// 验证内容
			if !strings.Contains(seg.Content, "<Box>") {
				t.Error("JSX should contain <Box>")
			}
			if !strings.Contains(seg.Content, "<Text>") {
				t.Error("JSX should contain <Text>")
			}
		}
	}

	if !jsxFound {
		t.Error("No JSX segment found")
		for i, seg := range segments {
			t.Logf("Segment %d: %v", i, seg)
		}
	}
}

func TestScanMixed(t *testing.T) {
	source := `package main

func App() Element {
    name := "World"
    return (
        <Box>
            <Text>Hello {name}</Text>
        </Box>
    )
}
`
	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 应该有 3 段: GoCode + JSX + GoCode
	if len(segments) != 3 {
		t.Fatalf("Expected 3 segments, got %d", len(segments))
		for i, seg := range segments {
			t.Logf("Segment %d: %v", i, seg)
		}
	}

	// 第一段应该是 GoCode
	if segments[0].Type != SegmentGoCode {
		t.Errorf("Expected first segment to be GoCode, got %v", segments[0].Type)
	}

	// 第二段应该是 JSX
	if segments[1].Type != SegmentJSX {
		t.Errorf("Expected second segment to be JSX, got %v", segments[1].Type)
	}

	// 第三段应该是 GoCode
	if segments[2].Type != SegmentGoCode {
		t.Errorf("Expected third segment to be GoCode, got %v", segments[2].Type)
	}
}

func TestScanNestedJSX(t *testing.T) {
	source := `return (
    <Box flexDirection="column">
        <Box>
            <Text>Nested</Text>
        </Box>
    </Box>
)`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 找到 JSX 段
	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证 JSX 内容包含嵌套标签
	expected := []string{"<Box", "<Text>", "Nested", "</Text>", "</Box>"}
	for _, exp := range expected {
		if !strings.Contains(jsxSegment.Content, exp) {
			t.Errorf("JSX content missing expected string: %s", exp)
		}
	}
}

func TestScanExpressionInJSX(t *testing.T) {
	source := `return <Text>Count: {count}</Text>`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证表达式被正确捕获
	if !strings.Contains(jsxSegment.Content, "{count}") {
		t.Error("JSX segment should contain expression {count}")
	}
}

func TestScanStringWithQuotes(t *testing.T) {
	source := `return <Text color="green">Hello "World"</Text>`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证字符串中的引号被正确处理
	if !strings.Contains(jsxSegment.Content, `"World"`) {
		t.Error("JSX should preserve quoted string content")
	}
}

func TestScanGoCodeWithStrings(t *testing.T) {
	source := `package main

func main() {
    str := "not <JSX>"
    return (
        <Text>Real JSX</Text>
    )
}`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 验证字符串中的 < 不产生额外的 JSX 段
	jsxCount := 0
	for _, seg := range segments {
		if seg.Type == SegmentJSX {
			jsxCount++
		}
	}

	// 应该只有 1 个 JSX 段 (Real JSX)
	if jsxCount != 1 {
		t.Errorf("Expected 1 JSX segment, got %d", jsxCount)
	}

	// 验证 JSX 段包含正确的内容
	for _, seg := range segments {
		if seg.Type == SegmentJSX {
			if !strings.Contains(seg.Content, "Real JSX") {
				t.Error("JSX segment should contain 'Real JSX'")
			}
		}
	}
}

func TestScanCommentHandling(t *testing.T) {
	source := `package main

func main() {
    // This is not <JSX> in comment
    /* Neither is <JSX> here */
    return (
        <Text>Real JSX</Text>
    )
}`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 应该只有 1 个 JSX 段
	jsxCount := 0
	for _, seg := range segments {
		if seg.Type == SegmentJSX {
			jsxCount++
		}
	}

	if jsxCount != 1 {
		t.Errorf("Expected 1 JSX segment, got %d", jsxCount)
	}
}

func TestScanSelfClosing(t *testing.T) {
	source := `return (
    <Box>
        <Spacer />
        <Text>Content</Text>
    </Box>
)`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证自闭合标签被正确捕获
	if !strings.Contains(jsxSegment.Content, "<Spacer />") {
		t.Error("JSX should contain self-closing tag")
	}
}

func TestScanLineColumnTracking(t *testing.T) {
	source := `package main

func App() Element {
    return (
        <Box>
            <Text>Hello</Text>
        </Box>
    )
}`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 找到 JSX 段
	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证行号 (JSX 应该在第 5 行开始)
	if jsxSegment.Line < 4 || jsxSegment.Line > 6 {
		t.Errorf("Expected JSX to start around line 5, got line %d", jsxSegment.Line)
	}
}

func TestScanMultipleJSX(t *testing.T) {
	source := `func App() Element {
    a := <Text>A</Text>
    b := <Text>B</Text>
    return <Box>{a}{b}</Box>
}`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// 应该有多个 JSX 段
	jsxCount := 0
	for _, seg := range segments {
		if seg.Type == SegmentJSX {
			jsxCount++
		}
	}

	if jsxCount < 3 {
		t.Errorf("Expected at least 3 JSX segments, got %d", jsxCount)
		for i, seg := range segments {
			t.Logf("Segment %d: %v", i, seg)
		}
	}
}

func TestScanFragment(t *testing.T) {
	source := `return (
    <>
        <Text>A</Text>
        <Text>B</Text>
    </>
)`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证 Fragment 被正确捕获
	if !strings.Contains(jsxSegment.Content, "<>") {
		t.Error("JSX should contain Fragment <>")
	}
	if !strings.Contains(jsxSegment.Content, "</>") {
		t.Error("JSX should contain Fragment close </>")
	}
}

func TestScanConditionalRender(t *testing.T) {
	source := `return (
    <Box>
        {show && <Text>Visible</Text>}
    </Box>
)`

	s := New(source)
	segments, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var jsxSegment *Segment
	for i := range segments {
		if segments[i].Type == SegmentJSX {
			jsxSegment = &segments[i]
			break
		}
	}

	if jsxSegment == nil {
		t.Fatal("No JSX segment found")
	}

	// 验证条件渲染表达式被捕获
	if !strings.Contains(jsxSegment.Content, "show") {
		t.Error("JSX should contain condition variable")
	}
}
