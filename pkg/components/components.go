// Package components 提供内置组件
package components

import (
	"github.com/ayanmw/go-react-ink/pkg/core"
)

// BoxProps Box 组件属性
type BoxProps struct {
	// Flexbox 属性
	FlexDirection  string // row, column, row-reverse, column-reverse
	JustifyContent string // flex-start, flex-end, center, space-between, space-around
	AlignItems     string // flex-start, flex-end, center, stretch, baseline
	AlignSelf      string // flex-start, flex-end, center, stretch, baseline
	FlexWrap       string // nowrap, wrap, wrap-reverse
	FlexGrow       int
	FlexShrink     int
	FlexBasis      string
	Order          int

	// 尺寸
	Width     int
	Height    int
	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int

	// 边距
	Margin        int
	MarginTop     int
	MarginRight   int
	MarginBottom  int
	MarginLeft    int
	Padding       int
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int
	PaddingLeft   int

	// 边框
	BorderStyle  string
	BorderColor  string
	BorderTop    bool
	BorderRight  bool
	BorderBottom bool
	BorderLeft   bool

	// 其他
	Display  string // flex, none
	Overflow string // visible, hidden
}

// Box Box 容器组件
func Box(props core.Props, children []core.Element) core.Element {
	return &BoxElement{
		Props:    props,
		Children: children,
	}
}

// BoxElement Box 元素
type BoxElement struct {
	Props    core.Props
	Children []core.Element
}

// Render 渲染 Box
func (e *BoxElement) Render() string {
	// 简化版：渲染子元素
	var content string
	for _, child := range e.Children {
		content += child.Render()
	}
	return content
}

// TextProps Text 组件属性
type TextProps struct {
	// 内容
	Children string

	// 样式
	Color      string
	Background string
	Bold       bool
	Italic     bool
	Underline  bool
	Dim        bool
	Blink      bool
	Reverse    bool

	// 对齐
	Wrap     string // truncate, wrap, wrap-middle, wrap-end
	Overflow string // truncate, ellipsis
}

// Text 文本组件
func Text(props core.Props, children []core.Element) core.Element {
	// 获取文本内容
	var content string
	if props != nil {
		if c, ok := props["children"].(string); ok {
			content = c
		}
	}
	if content == "" && len(children) > 0 {
		for _, child := range children {
			content += child.Render()
		}
	}

	return &TextElement{
		Content: content,
		Props:   props,
	}
}

// TextElement 文本元素
type TextElement struct {
	Content string
	Props   core.Props
}

// Render 渲染文本
func (e *TextElement) Render() string {
	return e.Content
}

// Spacer 空白填充组件
func Spacer(props core.Props, children []core.Element) core.Element {
	return &SpacerElement{Props: props}
}

// SpacerElement 空白元素
type SpacerElement struct {
	Props core.Props
}

// Render 渲染空白
func (e *SpacerElement) Render() string {
	return " "
}

// Newline 换行组件
func Newline(props core.Props, children []core.Element) core.Element {
	return &NewlineElement{Props: props}
}

// NewlineElement 换行元素
type NewlineElement struct {
	Props core.Props
}

// Render 渲染换行
func (e *NewlineElement) Render() string {
	return "\n"
}

// StaticProps Static 组件属性
type StaticProps struct {
	Items []string
}

// Static 静态内容组件 (不参与重渲染)
func Static(props core.Props, children []core.Element) core.Element {
	return &StaticElement{
		Props:    props,
		Children: children,
	}
}

// StaticElement 静态元素
type StaticElement struct {
	Props    core.Props
	Children []core.Element
}

// Render 渲染静态内容
func (e *StaticElement) Render() string {
	var content string
	for _, child := range e.Children {
		content += child.Render()
	}
	return content
}

// TransformProps Transform 组件属性
type TransformProps struct {
	Transform func(string) string
}

// Transform 文本转换组件
func Transform(props core.Props, children []core.Element) core.Element {
	// 获取子元素内容
	var content string
	for _, child := range children {
		content += child.Render()
	}

	// 应用转换
	if props != nil {
		if fn, ok := props["transform"].(func(string) string); ok {
			content = fn(content)
		}
	}

	return &TextElement{Content: content}
}

// Fragment Fragment 组件
func Fragment(props core.Props, children []core.Element) core.Element {
	return &FragmentElement{Children: children}
}

// FragmentElement Fragment 元素
type FragmentElement struct {
	Children []core.Element
}

// Render 渲染 Fragment
func (e *FragmentElement) Render() string {
	var content string
	for _, child := range e.Children {
		content += child.Render()
	}
	return content
}
