// Package fiber 实现 Fiber 协调器架构
package fiber

import (
	"time"
)

// Priority 优先级 (简化版 Lane 系统)
type Priority int

const (
	// PriorityNormal is default priority
	PriorityNormal Priority = iota
	// PriorityHigh is for urgent updates
	PriorityHigh
	// PriorityLow is for deferred updates
	PriorityLow
)

// WorkTag Fiber 工作类型
type WorkTag int

const (
	// TagFunctionComponent marks function component
	TagFunctionComponent WorkTag = iota
	// TagHostComponent marks host component
	TagHostComponent
	// TagHostText marks host text
	TagHostText
	// TagFragment marks fragment
	TagFragment
	// TagRoot marks root node
	TagRoot
)

// EffectTag 副作用标记
type EffectTag int

const (
	// EffectNoEffect means no effect
	EffectNoEffect EffectTag = 0
	// EffectPlacement marks insertion
	EffectPlacement EffectTag = 1 << iota
	// EffectUpdate marks update
	EffectUpdate
	// EffectDeletion marks deletion
	EffectDeletion
	// EffectRef marks ref update
	EffectRef
)

// Fiber Fiber 节点
type Fiber struct {
	// 树结构
	Tag           WorkTag
	Key           string
	ElementID     string // React 内部 ID
	Type          any    // 组件类型
	PendingProps  any    // 新属性
	MemoizedProps any    // 旧属性

	// Fiber 树
	Return  *Fiber // 父节点
	Child   *Fiber // 第一个子节点
	Sibling *Fiber // 下一个兄弟节点
	Index   int    // 子节点索引

	// 状态
	StateNode     any // 实例节点 (DOM, 组件实例)
	MemoizedState any // Hooks 状态链表

	// 副作用
	Flags        EffectTag // 当前副作用
	SubtreeFlags EffectTag // 子树副作用
	Deletions    []*Fiber  // 待删除节点

	// 替代树 (双缓冲)
	Alternate *Fiber

	// 调度
	Lanes      Priority // 剩余工作优先级
	ChildLanes Priority // 子节点工作优先级

	// 性能追踪
	ActualStartTime time.Time
	ActualDuration  time.Duration
}

// RootFiber 根节点容器
type RootFiber struct {
	ContainerInfo any    // 容器信息 (终端实例)
	Current       *Fiber // 当前 Fiber 树
	FinishedWork  *Fiber // 完成的工作

	// 调度
	PendingLanes Priority
}

// NewFiber 创建新的 Fiber 节点
func NewFiber(tag WorkTag, pendingProps any, key string) *Fiber {
	return &Fiber{
		Tag:          tag,
		Key:          key,
		PendingProps: pendingProps,
		Flags:        EffectNoEffect,
	}
}

// CreateFiberFromElement 从元素创建 Fiber
func CreateFiberFromElement(elementType any, props any, key string) *Fiber {
	var tag WorkTag
	switch v := elementType.(type) {
	case string:
		tag = TagHostComponent
	case func(any) any:
		tag = TagFunctionComponent
	default:
		_ = v
		tag = TagFunctionComponent
	}

	return NewFiber(tag, props, key)
}

// CreateHostTextFiber 创建文本 Fiber
func CreateHostTextFiber(text string) *Fiber {
	fiber := NewFiber(TagHostText, text, "")
	fiber.StateNode = text
	return fiber
}

// CreateWorkInProgress 创建工作副本 (双缓冲)
func CreateWorkInProgress(current *Fiber, pendingProps any) *Fiber {
	var workInProgress *Fiber

	if current.Alternate != nil {
		// 复用已有副本
		workInProgress = current.Alternate
		workInProgress.PendingProps = pendingProps
		workInProgress.Flags = EffectNoEffect
		workInProgress.SubtreeFlags = EffectNoEffect
		workInProgress.Deletions = nil
	} else {
		// 创建新副本
		workInProgress = NewFiber(current.Tag, pendingProps, current.Key)
		workInProgress.Type = current.Type
		workInProgress.StateNode = current.StateNode
		workInProgress.Alternate = current
		current.Alternate = workInProgress
	}

	// 复制树结构
	workInProgress.Child = nil
	workInProgress.Sibling = nil
	workInProgress.Index = 0
	workInProgress.Return = nil

	return workInProgress
}

// AppendChild 添加子节点
func (f *Fiber) AppendChild(child *Fiber) {
	child.Return = f

	if f.Child == nil {
		f.Child = child
		return
	}

	// 找到最后一个子节点
	last := f.Child
	for last.Sibling != nil {
		last = last.Sibling
	}
	last.Sibling = child
}

// RemoveChild 移除子节点
func (f *Fiber) RemoveChild(child *Fiber) {
	if f.Child == nil {
		return
	}

	if f.Child == child {
		f.Child = child.Sibling
		child.Return = nil
		child.Sibling = nil
		return
	}

	// 在兄弟链中查找
	prev := f.Child
	for prev.Sibling != nil {
		if prev.Sibling == child {
			prev.Sibling = child.Sibling
			child.Return = nil
			child.Sibling = nil
			return
		}
		prev = prev.Sibling
	}
}

// FindFirstDefibr 找到第一个子节点
func (f *Fiber) FindFirstDefibr() *Fiber {
	return f.Child
}

// FindNextDefibr 找到下一个工作节点
func (f *Fiber) FindNextDefibr() *Fiber {
	if f.Child != nil {
		return f.Child
	}

	// 向上查找兄弟节点
	node := f
	for node != nil {
		if node.Sibling != nil {
			return node.Sibling
		}
		node = node.Return
	}

	return nil
}

// HasEffect 检查是否有副作用
func (f *Fiber) HasEffect() bool {
	return f.Flags != EffectNoEffect || f.SubtreeFlags != EffectNoEffect
}

// MarkEffect 标记副作用
func (f *Fiber) MarkEffect(flag EffectTag) {
	f.Flags |= flag
}

// ClearEffect 清除副作用
func (f *Fiber) ClearEffect() {
	f.Flags = EffectNoEffect
	f.SubtreeFlags = EffectNoEffect
}

// String 返回 Fiber 的字符串表示
func (f *Fiber) String() string {
	var typeName string
	switch f.Tag {
	case TagFunctionComponent:
		typeName = "FunctionComponent"
	case TagHostComponent:
		typeName = "HostComponent"
	case TagHostText:
		typeName = "HostText"
	case TagFragment:
		typeName = "Fragment"
	case TagRoot:
		typeName = "Root"
	default:
		typeName = "Unknown"
	}

	return typeName + "(" + f.Key + ")"
}
