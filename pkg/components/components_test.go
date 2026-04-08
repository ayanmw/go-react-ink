package components

import (
	"strings"
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/core"
)

func TestBox(t *testing.T) {
	props := core.Props{"flexDirection": "column"}
	children := []core.Element{&TextElement{Content: "Hello"}}

	element := Box(props, children)

	if element == nil {
		t.Fatal("Box should return element")
	}

	output := element.Render()
	if !strings.Contains(output, "Hello") {
		t.Error("Box should render children")
	}
}

func TestText(t *testing.T) {
	props := core.Props{"children": "Hello World"}

	element := Text(props, nil)

	if element == nil {
		t.Fatal("Text should return element")
	}

	output := element.Render()
	if output != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", output)
	}
}

func TestTextWithChildren(t *testing.T) {
	children := []core.Element{&TextElement{Content: "Child"}}

	element := Text(nil, children)

	output := element.Render()
	if output != "Child" {
		t.Errorf("Expected 'Child', got %q", output)
	}
}

func TestSpacer(t *testing.T) {
	element := Spacer(nil, nil)

	output := element.Render()
	if output != " " {
		t.Errorf("Expected ' ', got %q", output)
	}
}

func TestNewline(t *testing.T) {
	element := Newline(nil, nil)

	output := element.Render()
	if output != "\n" {
		t.Errorf("Expected newline, got %q", output)
	}
}

func TestStatic(t *testing.T) {
	children := []core.Element{
		&TextElement{Content: "Line 1\n"},
		&TextElement{Content: "Line 2\n"},
	}

	element := Static(nil, children)

	output := element.Render()
	if !strings.Contains(output, "Line 1") || !strings.Contains(output, "Line 2") {
		t.Error("Static should render all children")
	}
}

func TestTransform(t *testing.T) {
	children := []core.Element{&TextElement{Content: "hello"}}
	props := core.Props{
		"transform": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	element := Transform(props, children)

	output := element.Render()
	if output != "HELLO" {
		t.Errorf("Expected 'HELLO', got %q", output)
	}
}

func TestFragment(t *testing.T) {
	children := []core.Element{
		&TextElement{Content: "A"},
		&TextElement{Content: "B"},
	}

	element := Fragment(nil, children)

	output := element.Render()
	if output != "AB" {
		t.Errorf("Expected 'AB', got %q", output)
	}
}

func TestBoxElementRender(t *testing.T) {
	box := &BoxElement{
		Children: []core.Element{
			&TextElement{Content: "Hello"},
			&TextElement{Content: " "},
			&TextElement{Content: "World"},
		},
	}

	output := box.Render()
	if output != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", output)
	}
}

func TestTextElementRender(t *testing.T) {
	text := &TextElement{Content: "Test"}

	output := text.Render()
	if output != "Test" {
		t.Errorf("Expected 'Test', got %q", output)
	}
}

func TestSpacerElementRender(t *testing.T) {
	spacer := &SpacerElement{}

	output := spacer.Render()
	if output != " " {
		t.Errorf("Expected ' ', got %q", output)
	}
}

func TestNewlineElementRender(t *testing.T) {
	newline := &NewlineElement{}

	output := newline.Render()
	if output != "\n" {
		t.Errorf("Expected newline, got %q", output)
	}
}

func TestStaticElementRender(t *testing.T) {
	static := &StaticElement{
		Children: []core.Element{
			&TextElement{Content: "Static"},
		},
	}

	output := static.Render()
	if output != "Static" {
		t.Errorf("Expected 'Static', got %q", output)
	}
}

func TestFragmentElementRender(t *testing.T) {
	fragment := &FragmentElement{
		Children: []core.Element{
			&TextElement{Content: "A"},
			&TextElement{Content: "B"},
			&TextElement{Content: "C"},
		},
	}

	output := fragment.Render()
	if output != "ABC" {
		t.Errorf("Expected 'ABC', got %q", output)
	}
}

func TestBoxWithProps(t *testing.T) {
	props := core.Props{
		"flexDirection":  "column",
		"justifyContent": "center",
		"alignItems":     "center",
		"padding":        1,
	}

	element := Box(props, nil)

	if element == nil {
		t.Fatal("Box should return element")
	}

	box := element.(*BoxElement)
	if box.Props["flexDirection"] != "column" {
		t.Error("Props should be preserved")
	}
}

func TestTextWithStyle(t *testing.T) {
	props := core.Props{
		"children":  "Styled",
		"color":     "green",
		"bold":      true,
		"underline": true,
	}

	element := Text(props, nil)

	text := element.(*TextElement)
	if text.Content != "Styled" {
		t.Error("Content should be 'Styled'")
	}
	if text.Props["color"] != "green" {
		t.Error("Color should be green")
	}
}