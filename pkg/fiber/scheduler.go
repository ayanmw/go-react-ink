// Package fiber 实现调度器和工作循环
package fiber

import (
	"fmt"
	"time"
)

// Scheduler 调度器
type Scheduler struct {
	root      *RootFiber
	workRoot  *Fiber  // 当前工作树
	nextUnitOfWork *Fiber // 下一个工作单元

	// 调度状态
	isWorking bool
	isCommitting bool
	pendingLanes Priority

	// 时间切片
	timeSlice     time.Duration // 16ms 时间切片
	startTime     time.Time
	didYield      bool

	// 副作用
	pendingEffects []*Fiber
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		timeSlice:      16 * time.Millisecond,
		pendingEffects: make([]*Fiber, 0),
	}
}

// ScheduleUpdate 调度更新
func (s *Scheduler) ScheduleUpdate(root *RootFiber) {
	s.root = root

	// 创建工作树
	current := root.Current
	s.workRoot = CreateWorkInProgress(current, current.PendingProps)

	// 开始工作循环
	s.nextUnitOfWork = s.workRoot
	s.isWorking = true
	s.isCommitting = false
}

// RequestWork 请求工作单元
func (s *Scheduler) RequestWork() *Fiber {
	return s.nextUnitOfWork
}

// PerformUnitOfWork 执行工作单元
func (s *Scheduler) PerformUnitOfWork(unit *Fiber) *Fiber {
	// 检查时间切片
	s.checkTimeSlice()

	// 执行 beginWork
	next := s.beginWork(unit)

	if next == nil {
		// 完成当前节点
		s.completeWork(unit)
	}

	// 找下一个工作单元
	s.nextUnitOfWork = s.findNextUnitOfWork(unit)

	return s.nextUnitOfWork
}

// beginWork 开始工作
func (s *Scheduler) beginWork(unit *Fiber) *Fiber {
	// 设置开始时间
	unit.ActualStartTime = time.Now()

	switch unit.Tag {
	case TagFunctionComponent:
		return s.beginFunctionComponent(unit)

	case TagHostComponent:
		return s.beginHostComponent(unit)

	case TagHostText:
		return s.beginHostText(unit)

	case TagFragment:
		return s.beginFragment(unit)

	case TagRoot:
		return s.beginRoot(unit)

	default:
		return nil
	}
}

// beginFunctionComponent 函数组件开始
func (s *Scheduler) beginFunctionComponent(unit *Fiber) *Fiber {
	// 获取组件函数
	component := unit.Type
	if component == nil {
		return nil
	}

	// 执行组件函数
	props := unit.PendingProps

	// 简化版：直接调用组件函数
	if fn, ok := component.(func(any) any); ok {
		result := fn(props)
		s.reconcileChildren(unit, result)
	} else {
		// 组件为空，创建空子树
		unit.Child = nil
	}

	// 返回第一个子节点
	return unit.Child
}

// beginHostComponent Host 组件开始
func (s *Scheduler) beginHostComponent(unit *Fiber) *Fiber {
	// 简化版：直接处理 props
	props := unit.PendingProps

	if props == nil {
		props = unit.MemoizedProps
	}

	// 标记需要更新
	if unit.MemoizedProps != nil && props != nil {
		if !s.shallowEqual(unit.MemoizedProps, props) {
			unit.MarkEffect(EffectUpdate)
		}
	}

	// 处理子节点
	if children := s.getChildren(props); children != nil {
		s.reconcileChildren(unit, children)
	}

	return unit.Child
}

// beginHostText Host 文本开始
func (s *Scheduler) beginHostText(unit *Fiber) *Fiber {
	props := unit.PendingProps

	if props == nil {
		props = unit.MemoizedProps
	}

	// 文本节点不需要子节点
	if unit.MemoizedProps != nil && props != nil {
		oldText := s.getText(unit.MemoizedProps)
		newText := s.getText(props)
		if oldText != newText {
			unit.MarkEffect(EffectUpdate)
			unit.StateNode = newText
		}
	} else {
		unit.StateNode = s.getText(props)
	}

	return nil
}

// beginFragment Fragment 开始
func (s *Scheduler) beginFragment(unit *Fiber) *Fiber {
	// 处理子节点
	if children := s.getChildren(unit.PendingProps); children != nil {
		s.reconcileChildren(unit, children)
	}

	return unit.Child
}

// beginRoot Root 开始
func (s *Scheduler) beginRoot(unit *Fiber) *Fiber {
	// 处理子节点
	if children := s.getChildren(unit.PendingProps); children != nil {
		s.reconcileChildren(unit, children)
	}

	return unit.Child
}

// completeWork 完成工作
func (s *Scheduler) completeWork(unit *Fiber) {
	// 记录耗时
	unit.ActualDuration = time.Since(unit.ActualStartTime)

	// 收集子树副作用
	if unit.Child != nil {
		unit.SubtreeFlags |= unit.Child.SubtreeFlags | unit.Child.Flags
	}

	// 向上传递副作用
	if unit.Return != nil {
		unit.Return.SubtreeFlags |= unit.SubtreeFlags | unit.Flags
	}

	// 添加到待提交列表
	if unit.Flags != EffectNoEffect {
		s.pendingEffects = append(s.pendingEffects, unit)
	}
}

// reconcileChildren 协调子节点
func (s *Scheduler) reconcileChildren(parent *Fiber, children any) {
	if children == nil {
		parent.Child = nil
		return
	}

	// 处理单个元素或数组
	switch v := children.(type) {
	case *Fiber:
		// 单个 Fiber
		s.reconcileSingleChild(parent, v)

	case []*Fiber:
		// Fiber 数组
		s.reconcileChildrenArray(parent, v)

	case []any:
		// 元素数组
		s.reconcileElementArray(parent, v)

	case string:
		// 文本
		textFiber := CreateHostTextFiber(v)
		parent.AppendChild(textFiber)

	default:
		_ = v
		parent.Child = nil
	}
}

// reconcileSingleChild 协调单个子节点
func (s *Scheduler) reconcileSingleChild(parent *Fiber, child *Fiber) {
	// 简化版：直接替换
	parent.Child = nil
	parent.AppendChild(child)

	// 标记插入
	child.MarkEffect(EffectPlacement)
}

// reconcileChildrenArray 协调子节点数组
func (s *Scheduler) reconcileChildrenArray(parent *Fiber, children []*Fiber) {
	parent.Child = nil

	for _, child := range children {
		parent.AppendChild(child)
		child.MarkEffect(EffectPlacement)
	}
}

// reconcileElementArray 协调元素数组
func (s *Scheduler) reconcileElementArray(parent *Fiber, elements []any) {
	parent.Child = nil

	for i, element := range elements {
		// 创建 Fiber
		key := fmt.Sprintf("%d", i)
		fiber := s.createFiberFromElement(element, key)
		parent.AppendChild(fiber)
		fiber.MarkEffect(EffectPlacement)
	}
}

// findNextUnitOfWork 找到下一个工作单元
func (s *Scheduler) findNextUnitOfWork(completedUnit *Fiber) *Fiber {
	if s.didYield {
		// 时间切片耗尽，暂停工作
		return nil
	}

	// 检查兄弟节点
	if completedUnit.Sibling != nil {
		return completedUnit.Sibling
	}

	// 向上查找
	unit := completedUnit.Return
	for unit != nil {
		// 完成父节点
		s.completeWork(unit)

		if unit.Sibling != nil {
			return unit.Sibling
		}
		unit = unit.Return
	}

	// 工作完成
	return nil
}

// checkTimeSlice 检查时间切片
func (s *Scheduler) checkTimeSlice() {
	if s.startTime.IsZero() {
		s.startTime = time.Now()
		return
	}

	if time.Since(s.startTime) > s.timeSlice {
		s.didYield = true
	}
}

// IsWorkComplete 检查工作是否完成
func (s *Scheduler) IsWorkComplete() bool {
	return s.nextUnitOfWork == nil
}

// CommitRoot 提交根节点
func (s *Scheduler) CommitRoot() {
	s.isCommitting = true

	// 执行副作用
	for _, unit := range s.pendingEffects {
		s.commitWork(unit)
	}

	// 清空副作用列表
	s.pendingEffects = make([]*Fiber, 0)

	// 更新根节点
	if s.root != nil {
		s.root.Current = s.workRoot
	}

	s.isCommitting = false
	s.isWorking = false
}

// commitWork 执行副作用
func (s *Scheduler) commitWork(unit *Fiber) {
	flags := unit.Flags

	if flags&EffectPlacement != 0 {
		s.commitPlacement(unit)
	}

	if flags&EffectUpdate != 0 {
		s.commitUpdate(unit)
	}

	if flags&EffectDeletion != 0 {
		s.commitDeletion(unit)
	}

	// 清除副作用
	unit.ClearEffect()
}

// commitPlacement 提交插入
func (s *Scheduler) commitPlacement(unit *Fiber) {
	// 找到父节点容器
	parent := unit.Return
	for parent != nil && parent.Tag != TagHostComponent && parent.Tag != TagRoot {
		parent = parent.Return
	}

	if parent == nil {
		return
	}

	// 在实际实现中，这里会创建实例并插入到容器
	// 简化版：仅记录
	unit.StateNode = unit.Type
}

// commitUpdate 提交更新
func (s *Scheduler) commitUpdate(unit *Fiber) {
	// 更新属性
	unit.MemoizedProps = unit.PendingProps

	// 在实际实现中，这里会更新实例属性
}

// commitDeletion 提交删除
func (s *Scheduler) commitDeletion(unit *Fiber) {
	// 递归删除子节点
	child := unit.Child
	for child != nil {
		s.commitDeletion(child)
		child = child.Sibling
	}

	// 清除节点
	unit.StateNode = nil
	unit.MemoizedProps = nil
	unit.MemoizedState = nil
}

// 辅助方法

func (s *Scheduler) shallowEqual(a, b any) bool {
	// 简化版：直接比较
	return a == b
}

func (s *Scheduler) getChildren(props any) any {
	// 简化版：从 props 获取 children
	if props == nil {
		return nil
	}

	// 如果 props 是 map
	if m, ok := props.(map[string]any); ok {
		return m["children"]
	}

	return nil
}

func (s *Scheduler) getText(props any) string {
	if props == nil {
		return ""
	}

	if str, ok := props.(string); ok {
		return str
	}

	if m, ok := props.(map[string]any); ok {
		if children, ok := m["children"]; ok {
			if str, ok := children.(string); ok {
				return str
			}
		}
	}

	return ""
}

func (s *Scheduler) createFiberFromElement(element any, key string) *Fiber {
	// 简化版：根据元素类型创建 Fiber
	if str, ok := element.(string); ok {
		return CreateHostTextFiber(str)
	}

	if fiber, ok := element.(*Fiber); ok {
		fiber.Key = key
		return fiber
	}

	// 其他情况：创建 Function Fiber
	return NewFiber(TagFunctionComponent, element, key)
}