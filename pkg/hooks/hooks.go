// Package hooks 实现 React-like Hooks
package hooks

import (
	"sync"
)

// HookContext Hook 上下文
type HookContext struct {
	mu         sync.Mutex
	hookIndex  int
	hooks      []any
	dispatcher Dispatcher
}

// Hook Hook 接口
type Hook interface {
	Cleanup()
}

// Dispatcher 调度器接口
type Dispatcher interface {
	ScheduleUpdate()
}

// NewHookContext 创建 Hook 上下文
func NewHookContext() *HookContext {
	return &HookContext{
		hooks: make([]any, 0),
	}
}

// Reset 重置 Hook 索引 (每次渲染开始时调用)
func (c *HookContext) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hookIndex = 0
}

// SetDispatcher 设置调度器
func (c *HookContext) SetDispatcher(d Dispatcher) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dispatcher = d
}

// nextHook 获取下一个 Hook 索引
func (c *HookContext) nextHook() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	idx := c.hookIndex
	c.hookIndex++
	return idx
}

// scheduleUpdate 调度更新
func (c *HookContext) scheduleUpdate() {
	c.mu.Lock()
	d := c.dispatcher
	c.mu.Unlock()

	if d != nil {
		d.ScheduleUpdate()
	}
}

// StateHook State Hook
type StateHook struct {
	Value   any
	Setter  func(any)
	context *HookContext
}

// Cleanup 清理 (State 不需要清理)
func (h *StateHook) Cleanup() {}

// UseState 状态 Hook
func UseState(ctx *HookContext, initial any) (any, func(any)) {
	idx := ctx.nextHook()

	ctx.mu.Lock()
	var hook *StateHook
	if idx >= len(ctx.hooks) {
		// 创建新 Hook
		hook = &StateHook{
			Value:   initial,
			context: ctx,
		}
		hook.Setter = func(newValue any) {
			ctx.mu.Lock()
			hook.Value = newValue
			ctx.mu.Unlock()
			ctx.scheduleUpdate()
		}
		ctx.hooks = append(ctx.hooks, hook)
	} else {
		// 复用已有 Hook
		hook = ctx.hooks[idx].(*StateHook)
	}
	ctx.mu.Unlock()

	return hook.Value, hook.Setter
}

// EffectHook Effect Hook
type EffectHook struct {
	Setup     func() func() // 返回清理函数
	CleanupFn func()
	Deps      []any
	context   *HookContext
}

// Cleanup 执行清理
func (h *EffectHook) Cleanup() {
	if h.CleanupFn != nil {
		h.CleanupFn()
		h.CleanupFn = nil
	}
}

// UseEffect 副作用 Hook
func UseEffect(ctx *HookContext, setup func() func(), deps []any) {
	idx := ctx.nextHook()

	ctx.mu.Lock()
	var hook *EffectHook
	if idx >= len(ctx.hooks) {
		// 创建新 Hook
		hook = &EffectHook{
			Setup:   setup,
			Deps:    deps,
			context: ctx,
		}
		ctx.hooks = append(ctx.hooks, hook)
	} else {
		// 复用已有 Hook
		hook = ctx.hooks[idx].(*EffectHook)

		// 检查依赖是否变化
		if depsEqual(hook.Deps, deps) {
			// 依赖未变化，跳过
			ctx.mu.Unlock()
			return
		}

		// 执行旧清理
		hook.Cleanup()

		// 更新 Hook
		hook.Setup = setup
		hook.Deps = deps
	}
	ctx.mu.Unlock()

	// 执行 setup
	if setup != nil {
		cleanup := setup()
		ctx.mu.Lock()
		hook.CleanupFn = cleanup
		ctx.mu.Unlock()
	}
}

// UseLayoutEffect 布局副作用 Hook (同步执行)
func UseLayoutEffect(ctx *HookContext, setup func() func(), deps []any) {
	// 简化版：与 UseEffect 相同
	UseEffect(ctx, setup, deps)
}

// RefHook Ref Hook
type RefHook struct {
	Current any
}

// Cleanup 清理 (Ref 不需要清理)
func (h *RefHook) Cleanup() {}

// UseRef 引用 Hook
func UseRef(ctx *HookContext, initial any) *RefHook {
	idx := ctx.nextHook()

	ctx.mu.Lock()
	var hook *RefHook
	if idx >= len(ctx.hooks) {
		// 创建新 Hook
		hook = &RefHook{Current: initial}
		ctx.hooks = append(ctx.hooks, hook)
	} else {
		// 复用已有 Hook
		hook = ctx.hooks[idx].(*RefHook)
	}
	ctx.mu.Unlock()

	return hook
}

// MemoHook Memo Hook
type MemoHook struct {
	Value any
	Deps  []any
}

// Cleanup 清理 (Memo 不需要清理)
func (h *MemoHook) Cleanup() {}

// UseMemo 记忆化 Hook
func UseMemo(ctx *HookContext, factory func() any, deps []any) any {
	idx := ctx.nextHook()

	ctx.mu.Lock()
	var hook *MemoHook
	if idx >= len(ctx.hooks) {
		// 创建新 Hook
		hook = &MemoHook{
			Value: factory(),
			Deps:  deps,
		}
		ctx.hooks = append(ctx.hooks, hook)
	} else {
		// 复用已有 Hook
		hook = ctx.hooks[idx].(*MemoHook)

		// 检查依赖是否变化
		if !depsEqual(hook.Deps, deps) {
			// 依赖变化，重新计算
			hook.Value = factory()
			hook.Deps = deps
		}
	}
	ctx.mu.Unlock()

	return hook.Value
}

// UseCallback 回调记忆化 Hook
func UseCallback(ctx *HookContext, callback func(), deps []any) func() {
	return UseMemo(ctx, func() any { return callback }, deps).(func())
}

// depsEqual 比较依赖数组
func depsEqual(a, b []any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
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

// CleanupHooks 清理所有 Hook
func (c *HookContext) CleanupHooks() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, hook := range c.hooks {
		if h, ok := hook.(Hook); ok {
			h.Cleanup()
		}
	}
}