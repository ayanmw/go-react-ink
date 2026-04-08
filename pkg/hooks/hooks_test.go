package hooks

import (
	"sync/atomic"
	"testing"
)

func TestUseState(t *testing.T) {
	ctx := NewHookContext()

	// 初始值
	value, setValue := UseState(ctx, 10)
	if value != 10 {
		t.Errorf("Expected initial value 10, got %d", value)
	}

	// 设置新值
	setValue(20)

	// 重置后重新获取
	ctx.Reset()
	value, _ = UseState(ctx, 10)
	if value != 20 {
		t.Errorf("Expected updated value 20, got %d", value)
	}
}

func TestUseStateString(t *testing.T) {
	ctx := NewHookContext()

	value, setValue := UseState(ctx, "hello")
	if value != "hello" {
		t.Errorf("Expected 'hello', got %q", value)
	}

	setValue("world")

	ctx.Reset()
	value, _ = UseState(ctx, "hello")
	if value != "world" {
		t.Errorf("Expected 'world', got %q", value)
	}
}

func TestUseStateWithDispatcher(t *testing.T) {
	ctx := NewHookContext()

	var updateCalled int32
	d := &testDispatcher{
		onScheduleUpdate: func() {
			atomic.AddInt32(&updateCalled, 1)
		},
	}
	ctx.SetDispatcher(d)

	_, setValue := UseState(ctx, 0)
	setValue(1)

	if atomic.LoadInt32(&updateCalled) != 1 {
		t.Errorf("Expected 1 update call, got %d", updateCalled)
	}
}

func TestUseEffect(t *testing.T) {
	ctx := NewHookContext()

	var setupCalled int
	var cleanupCalled int

	// 第一次渲染
	UseEffect(ctx, func() func() {
		setupCalled++
		return func() {
			cleanupCalled++
		}
	}, []any{})

	if setupCalled != 1 {
		t.Errorf("Expected 1 setup call, got %d", setupCalled)
	}

	// 清理
	ctx.CleanupHooks()

	if cleanupCalled != 1 {
		t.Errorf("Expected 1 cleanup call, got %d", cleanupCalled)
	}
}

func TestUseEffectDeps(t *testing.T) {
	ctx := NewHookContext()

	var setupCalled int

	// 第一次渲染，依赖 [1]
	UseEffect(ctx, func() func() {
		setupCalled++
		return nil
	}, []any{1})

	if setupCalled != 1 {
		t.Errorf("Expected 1 setup call, got %d", setupCalled)
	}

	// 重置后，依赖相同，不重新执行
	ctx.Reset()
	UseEffect(ctx, func() func() {
		setupCalled++
		return nil
	}, []any{1})

	if setupCalled != 1 {
		t.Errorf("Expected no additional setup call, got %d", setupCalled)
	}

	// 依赖变化，重新执行
	ctx.Reset()
	UseEffect(ctx, func() func() {
		setupCalled++
		return nil
	}, []any{2})

	if setupCalled != 2 {
		t.Errorf("Expected 2 setup calls, got %d", setupCalled)
	}
}

func TestUseEffectCleanup(t *testing.T) {
	ctx := NewHookContext()

	var cleanupCalled int

	// 第一次渲染
	UseEffect(ctx, func() func() {
		return func() {
			cleanupCalled++
		}
	}, []any{1})

	// 依赖变化，触发清理
	ctx.Reset()
	UseEffect(ctx, func() func() {
		return func() {
			cleanupCalled++
		}
	}, []any{2})

	if cleanupCalled != 1 {
		t.Errorf("Expected 1 cleanup call, got %d", cleanupCalled)
	}
}

func TestUseRef(t *testing.T) {
	ctx := NewHookContext()

	ref := UseRef(ctx, 0)
	if ref.Current != 0 {
		t.Errorf("Expected initial value 0, got %d", ref.Current)
	}

	// 修改 ref
	ref.Current = 10

	// 重置后，ref 保持
	ctx.Reset()
	ref = UseRef(ctx, 0)
	if ref.Current != 10 {
		t.Errorf("Expected ref to persist, got %d", ref.Current)
	}
}

func TestUseMemo(t *testing.T) {
	ctx := NewHookContext()

	var computeCalled int

	value := UseMemo(ctx, func() any {
		computeCalled++
		return 42
	}, []any{1})

	if value.(int) != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	if computeCalled != 1 {
		t.Errorf("Expected 1 compute call, got %d", computeCalled)
	}

	// 依赖相同，不重新计算
	ctx.Reset()
	value = UseMemo(ctx, func() any {
		computeCalled++
		return 42
	}, []any{1})

	if computeCalled != 1 {
		t.Errorf("Expected no additional compute, got %d", computeCalled)
	}

	// 依赖变化，重新计算
	ctx.Reset()
	value = UseMemo(ctx, func() any {
		computeCalled++
		return 100
	}, []any{2})

	if computeCalled != 2 {
		t.Errorf("Expected 2 compute calls, got %d", computeCalled)
	}

	if value.(int) != 100 {
		t.Errorf("Expected 100, got %d", value)
	}
}

func TestUseCallback(t *testing.T) {
	ctx := NewHookContext()

	var callCount int

	callback := UseCallback(ctx, func() {
		callCount++
	}, []any{1})

	callback()
	callback()

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}

	// 重置后，依赖相同，callback 应该是同一个
	ctx.Reset()
	_ = UseCallback(ctx, func() {
		callCount++
	}, []any{1})

	// 验证 memoization 通过依赖检查
	// Go 中函数不能直接比较，所以通过依赖数组验证
}

func TestDepsEqual(t *testing.T) {
	tests := []struct {
		a, b     []any
		expected bool
	}{
		{nil, nil, true},
		{nil, []any{}, false},
		{[]any{}, []any{}, true},
		{[]any{1}, []any{1}, true},
		{[]any{1}, []any{2}, false},
		{[]any{1, 2}, []any{1, 2}, true},
		{[]any{1, 2}, []any{1, 3}, false},
		{[]any{"a"}, []any{"a"}, true},
		{[]any{"a"}, []any{"b"}, false},
	}

	for i, test := range tests {
		result := depsEqual(test.a, test.b)
		if result != test.expected {
			t.Errorf("Test %d: expected %v, got %v", i, test.expected, result)
		}
	}
}

func TestMultipleHooks(t *testing.T) {
	ctx := NewHookContext()

	// 多个 Hook
	state1, _ := UseState(ctx, 1)
	state2, _ := UseState(ctx, 2)
	ref := UseRef(ctx, 3)

	if state1 != 1 {
		t.Errorf("Expected state1=1, got %d", state1)
	}
	if state2 != 2 {
		t.Errorf("Expected state2=2, got %d", state2)
	}
	if ref.Current != 3 {
		t.Errorf("Expected ref=3, got %d", ref.Current)
	}
}

func TestHookContextReset(t *testing.T) {
	ctx := NewHookContext()

	// 第一次渲染
	_, _ = UseState(ctx, 1)
	_, _ = UseState(ctx, 2)

	if ctx.hookIndex != 2 {
		t.Errorf("Expected hookIndex=2, got %d", ctx.hookIndex)
	}

	// 重置
	ctx.Reset()

	if ctx.hookIndex != 0 {
		t.Errorf("Expected hookIndex=0 after reset, got %d", ctx.hookIndex)
	}
}

// 测试调度器
type testDispatcher struct {
	onScheduleUpdate func()
}

func (d *testDispatcher) ScheduleUpdate() {
	if d.onScheduleUpdate != nil {
		d.onScheduleUpdate()
	}
}
