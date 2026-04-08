// Package hostconfig 实现 React Reconciler Host Config
package hostconfig

import (
	"github.com/anmingwei/go-ink/pkg/core"
	"github.com/anmingwei/go-ink/pkg/fiber"
	"github.com/anmingwei/go-ink/pkg/layout"
	"github.com/anmingwei/go-ink/pkg/renderer"
)

// ComponentType 组件类型
type ComponentType int

const (
	ComponentTypeBox ComponentType = iota
	ComponentTypeText
	ComponentTypeSpacer
	ComponentTypeNewline
	ComponentTypeUnknown
)

// HostConfig Host 配置
type HostConfig struct {
	renderer    *renderer.Renderer
	rootNode    *layout.Node
	rootElement core.Element
	width       int
	height      int
}

// NewHostConfig 创建 Host 配置
func NewHostConfig(width, height int) *HostConfig {
	return &HostConfig{
		renderer: renderer.NewRenderer(width, height),
		width:    width,
		height:   height,
	}
}

// GetComponentType 获取组件类型
func GetComponentType(type_ any) ComponentType {
	// 通过名称或指针比较
	typeStr := getTypeString(type_)
	switch typeStr {
	case "Box":
		return ComponentTypeBox
	case "Text":
		return ComponentTypeText
	case "Spacer":
		return ComponentTypeSpacer
	case "Newline":
		return ComponentTypeNewline
	default:
		return ComponentTypeUnknown
	}
}

func getTypeString(type_ any) string {
	// 简化版：返回未知类型，实际使用时传入组件名称
	return "Unknown"
}

// CreateInstanceByName 通过名称创建实例
func (h *HostConfig) CreateInstanceByName(name string, props core.Props) any {
	node := layout.NewNode()
	switch name {
	case "Box":
		h.applyBoxProps(node, props)
	case "Text":
		h.applyTextProps(node, props)
	}
	return node
}

// UpdateInstanceByName 通过名称更新实例
func (h *HostConfig) UpdateInstanceByName(instance any, name string, props core.Props) {
	if node, ok := instance.(*layout.Node); ok {
		switch name {
		case "Box":
			h.applyBoxProps(node, props)
		case "Text":
			h.applyTextProps(node, props)
		}
	}
}

// CreateInstance 创建实例
func (h *HostConfig) CreateInstance(type_ any, props core.Props) any {
	// 创建布局节点
	node := layout.NewNode()

	// 根据类型设置属性
	compType := GetComponentType(type_)
	switch compType {
	case ComponentTypeBox:
		h.applyBoxProps(node, props)
	case ComponentTypeText:
		h.applyTextProps(node, props)
	}

	return node
}

// applyBoxProps 应用 Box 属性
func (h *HostConfig) applyBoxProps(node *layout.Node, props core.Props) {
	if v, ok := props["flexDirection"].(string); ok {
		switch v {
		case "row":
			node.Direction = layout.DirectionRow
		case "column":
			node.Direction = layout.DirectionColumn
		case "row-reverse":
			node.Direction = layout.DirectionRowReverse
		case "column-reverse":
			node.Direction = layout.DirectionColumnReverse
		}
	}

	if v, ok := props["justifyContent"].(string); ok {
		switch v {
		case "flex-start":
			node.Justify = layout.JustifyFlexStart
		case "flex-end":
			node.Justify = layout.JustifyFlexEnd
		case "center":
			node.Justify = layout.JustifyCenter
		case "space-between":
			node.Justify = layout.JustifySpaceBetween
		case "space-around":
			node.Justify = layout.JustifySpaceAround
		}
	}

	if v, ok := props["alignItems"].(string); ok {
		switch v {
		case "flex-start":
			node.AlignItems = layout.AlignFlexStart
		case "flex-end":
			node.AlignItems = layout.AlignFlexEnd
		case "center":
			node.AlignItems = layout.AlignCenter
		case "stretch":
			node.AlignItems = layout.AlignStretch
		}
	}

	if v, ok := props["flexGrow"].(int); ok {
		node.FlexGrow = float64(v)
	}

	if v, ok := props["flexShrink"].(int); ok {
		node.FlexShrink = float64(v)
	}

	if v, ok := props["width"].(int); ok {
		node.Width = float64(v)
	}

	if v, ok := props["height"].(int); ok {
		node.Height = float64(v)
	}

	if v, ok := props["padding"].(int); ok {
		node.SetPadding(float64(v), float64(v), float64(v), float64(v))
	}

	if v, ok := props["margin"].(int); ok {
		node.SetMargin(float64(v), float64(v), float64(v), float64(v))
	}
}

// applyTextProps 应用 Text 属性
func (h *HostConfig) applyTextProps(node *layout.Node, props core.Props) {
	// 文本节点使用测量函数
	node.MeasureFunc = func(width float64) (minWidth, maxWidth, height float64) {
		if children, ok := props["children"].(string); ok {
			// 简化版：假设每个字符宽度为 1
			textWidth := float64(len(children))
			return textWidth, textWidth, 1
		}
		return 0, 0, 1
	}
}

// AppendChild 添加子节点
func (h *HostConfig) AppendChild(parent, child any) {
	if parentNode, ok := parent.(*layout.Node); ok {
		if childNode, ok := child.(*layout.Node); ok {
			parentNode.AddChild(childNode)
		}
	}
}

// RemoveChild 移除子节点
func (h *HostConfig) RemoveChild(parent, child any) {
	if parentNode, ok := parent.(*layout.Node); ok {
		if childNode, ok := child.(*layout.Node); ok {
			parentNode.RemoveChild(childNode)
		}
	}
}

// InsertBefore 在指定节点前插入
func (h *HostConfig) InsertBefore(parent, child, beforeChild any) {
	// 简化版：添加到末尾
	h.AppendChild(parent, child)
}

// CreateTextInstance 创建文本实例
func (h *HostConfig) CreateTextInstance(text string) any {
	node := layout.NewNode()
	node.MeasureFunc = func(width float64) (minWidth, maxWidth, height float64) {
		textWidth := float64(len(text))
		return textWidth, textWidth, 1
	}
	return node
}

// UpdateInstance 更新实例
func (h *HostConfig) UpdateInstance(instance any, type_ any, props core.Props) {
	if node, ok := instance.(*layout.Node); ok {
		compType := GetComponentType(type_)
		switch compType {
		case ComponentTypeBox:
			h.applyBoxProps(node, props)
		case ComponentTypeText:
			h.applyTextProps(node, props)
		}
	}
}

// UpdateTextInstance 更新文本实例
func (h *HostConfig) UpdateTextInstance(instance any, text string) {
	if node, ok := instance.(*layout.Node); ok {
		node.MeasureFunc = func(width float64) (minWidth, maxWidth, height float64) {
			textWidth := float64(len(text))
			return textWidth, textWidth, 1
		}
	}
}

// FinalizeChildren 完成子节点处理
func (h *HostConfig) FinalizeChildren(instance any, children []any) {
	// 子节点已在 AppendChild 中处理
}

// ReplaceContainerChildren 替换容器子节点
func (h *HostConfig) ReplaceContainerChildren(container any, children []any) {
	if node, ok := container.(*layout.Node); ok {
		node.Children = nil
		for _, child := range children {
			if childNode, ok := child.(*layout.Node); ok {
				node.AddChild(childNode)
			}
		}
	}
}

// PrepareUpdate 准备更新
func (h *HostConfig) PrepareUpdate(instance any, type_ any, oldProps, newProps core.Props) any {
	// 返回需要更新的属性
	return newProps
}

// ShouldSetTextContent 是否应该设置文本内容
func (h *HostConfig) ShouldSetTextContent(type_ any, props core.Props) bool {
	// Text 组件应该设置文本内容
	if str, ok := type_.(string); ok {
		return str == "Text"
	}
	return false
}

// CanHydrate 是否可以水合
func (h *HostConfig) CanHydrate(instance any) bool {
	return false
}

// ClearContainer 清空容器
func (h *HostConfig) ClearContainer(container any) {
	if node, ok := container.(*layout.Node); ok {
		node.Children = nil
	}
}

// DetachInstance 分离实例
func (h *HostConfig) DetachInstance(instance any) {
	// 清理实例
}

// DetachTextInstance 分离文本实例
func (h *HostConfig) DetachTextInstance(instance any) {
	// 清理文本实例
}

// GetPublicInstance 获取公共实例
func (h *HostConfig) GetPublicInstance(instance any) any {
	return instance
}

// PrepareForCommit 准备提交
func (h *HostConfig) PrepareForCommit(container any) any {
	return nil
}

// ResetAfterCommit 提交后重置
func (h *HostConfig) ResetAfterCommit(container any) {
	// 计算布局
	if node, ok := container.(*layout.Node); ok {
		node.CalculateLayout(float64(h.width), float64(h.height))
	}
}

// ScheduleTimeout 调度超时
func (h *HostConfig) ScheduleTimeout(callback func(), delay int) any {
	// 简化版：立即执行
	go callback()
	return nil
}

// CancelTimeout 取消超时
func (h *HostConfig) CancelTimeout(timeoutID any) {
	// 简化版：无操作
}

// NoTimeout 无超时
func (h *HostConfig) NoTimeout() any {
	return nil
}

// IsPrimaryRenderer 是否主渲染器
func (h *HostConfig) IsPrimaryRenderer() bool {
	return true
}

// GetRootHostContext 获取根 Host 上下文
func (h *HostConfig) GetRootHostContext(rootContainerInstance any) any {
	return h
}

// GetChildHostContext 获取子 Host 上下文
func (h *HostConfig) GetChildHostContext(parentHostContext, type_ any, rootContainerInstance any) any {
	return parentHostContext
}

// PushTreeContext 推入树上下文
func (h *HostConfig) PushTreeContext(fiber *fiber.Fiber) {
	// 简化版：无操作
}

// PopTreeContext 弹出树上下文
func (h *HostConfig) PopTreeContext() {
	// 简化版：无操作
}

// SetRoot 设置根节点
func (h *HostConfig) SetRoot(root any) {
	if node, ok := root.(*layout.Node); ok {
		h.rootNode = node
	}
}

// Render 渲染
func (h *HostConfig) Render() string {
	if h.rootNode != nil {
		h.rootNode.CalculateLayout(float64(h.width), float64(h.height))
	}
	return h.renderer.Render()
}

// Resize 调整大小
func (h *HostConfig) Resize(width, height int) {
	h.width = width
	h.height = height
	h.renderer.Resize(width, height)
}

// GetRenderer 获取渲染器
func (h *HostConfig) GetRenderer() *renderer.Renderer {
	return h.renderer
}

// GetLayoutRoot 获取布局根节点
func (h *HostConfig) GetLayoutRoot() *layout.Node {
	return h.rootNode
}