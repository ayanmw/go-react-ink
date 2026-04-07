// Package core 提供 Go-Ink 核心类型和函数
package core

import (
	"fmt"
)

// Element 元素接口
type Element interface {
	Render() string
}

// Props 属性映射
type Props map[string]any

// Component 组件函数类型
type Component func(props Props, children []Element) Element

// CreateElement 创建元素 (类似 React.createElement)
func CreateElement(component Component, props Props, children []Element) Element {
	// 处理 spread 属性
	finalProps := props
	for _, spread := range extractSpread(props) {
		for k, v := range spread {
			if finalProps == nil {
				finalProps = Props{}
			}
			finalProps[k] = v
		}
	}

	return component(finalProps, children)
}

// Spread 创建 spread 属性
func Spread(props Props) Props {
	return Props{"__spread__": props}
}

// extractSpread 提取 spread 属性
func extractSpread(props Props) []Props {
	var spreads []Props
	if props == nil {
		return spreads
	}

	if spread, ok := props["__spread__"]; ok {
		if sp, ok := spread.(Props); ok {
			spreads = append(spreads, sp)
		}
	}

	return spreads
}

// TextComponent 文本组件
var Text = func(props Props, children []Element) Element {
	text := ""
	if children != nil && len(children) > 0 {
		for _, child := range children {
			text += child.Render()
		}
	}
	if props != nil {
		if content, ok := props["children"].(string); ok {
			text = content
		}
	}
	return &TextElement{Content: text}
}

// TextElement 文本元素
type TextElement struct {
	Content string
}

// Render 渲染文本
func (t *TextElement) Render() string {
	return t.Content
}

// BoxComponent Box 容器组件
var Box = func(props Props, children []Element) Element {
	return &BoxElement{
		Props:    props,
		Children: children,
	}
}

// BoxElement Box 元素
type BoxElement struct {
	Props    Props
	Children []Element
}

// Render 渲染 Box
func (b *BoxElement) Render() string {
	var content string
	if b.Children != nil {
		for _, child := range b.Children {
			content += child.Render()
		}
	}
	return fmt.Sprintf("[Box:%s]", content)
}

// FragmentComponent Fragment 组件
var Fragment = func(props Props, children []Element) Element {
	return &FragmentElement{Children: children}
}

// FragmentElement Fragment 元素
type FragmentElement struct {
	Children []Element
}

// Render 渲染 Fragment
func (f *FragmentElement) Render() string {
	var content string
	if f.Children != nil {
		for _, child := range f.Children {
			content += child.Render()
		}
	}
	return content
}

// String 返回 Element 的字符串表示（使用 Render 方法）
func RenderString(e Element) string {
	return e.Render()
}