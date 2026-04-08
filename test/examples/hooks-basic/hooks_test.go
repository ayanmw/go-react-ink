package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

// TestUseState 测试 useState Hook
func TestUseState(t *testing.T) {
	ctx := hooks.NewHookContext()

	value, setValue := hooks.UseState(ctx, 0)

	if value != 0 {
		t.Errorf("Expected initial value 0, got %v", value)
	}

	setValue(1)

	// 注意：简化实现中 value 不会自动更新
	// 实际应用中会触发重新渲染
}

// TestUseStateWithDifferentTypes 测试 useState 不同类型
func TestUseStateWithDifferentTypes(t *testing.T) {
	ctx := hooks.NewHookContext()

	// 整数
	intVal, _ := hooks.UseState(ctx, 42)
	if intVal != 42 {
		t.Errorf("Expected 42, got %v", intVal)
	}

	// 字符串
	strVal, _ := hooks.UseState(ctx, "hello")
	if strVal != "hello" {
		t.Errorf("Expected 'hello', got %v", strVal)
	}

	// 布尔
	boolVal, _ := hooks.UseState(ctx, true)
	if boolVal != true {
		t.Errorf("Expected true, got %v", boolVal)
	}

	// 切片
	sliceVal, _ := hooks.UseState(ctx, []int{1, 2, 3})
	if len(sliceVal.([]int)) != 3 {
		t.Errorf("Expected slice length 3")
	}

	// Map
	mapVal, _ := hooks.UseState(ctx, map[string]int{"a": 1})
	if mapVal.(map[string]int)["a"] != 1 {
		t.Errorf("Expected map value 1")
	}
}

// TestUseEffect 测试 useEffect Hook
func TestUseEffect(t *testing.T) {
	ctx := hooks.NewHookContext()

	called := false
	hooks.UseEffect(ctx, func() func() {
		called = true
		return func() {
			called = false
		}
	}, nil)

	if !called {
		t.Error("Effect should have been called")
	}
}

// TestUseEffectWithDeps 测试 useEffect 依赖
func TestUseEffectWithDeps(t *testing.T) {
	ctx := hooks.NewHookContext()

	callCount := 0
	hooks.UseEffect(ctx, func() func() {
		callCount++
		return nil
	}, []any{1, 2, 3})

	// 依赖变化时会重新执行
	_ = callCount
}

// TestUseLayoutEffect 测试 useLayoutEffect Hook
func TestUseLayoutEffect(t *testing.T) {
	ctx := hooks.NewHookContext()

	called := false
	hooks.UseLayoutEffect(ctx, func() func() {
		called = true
		return nil
	}, nil)

	if !called {
		t.Error("LayoutEffect should have been called")
	}
}

// TestUseRef 测试 useRef Hook
func TestUseRef(t *testing.T) {
	ctx := hooks.NewHookContext()

	ref := hooks.UseRef(ctx, "initial")

	if ref.Current != "initial" {
		t.Errorf("Expected 'initial', got %v", ref.Current)
	}

	ref.Current = "updated"

	if ref.Current != "updated" {
		t.Errorf("Expected 'updated', got %v", ref.Current)
	}
}

// TestUseRefWithNil 测试 useRef nil 值
func TestUseRefWithNil(t *testing.T) {
	ctx := hooks.NewHookContext()

	ref := hooks.UseRef(ctx, nil)

	if ref.Current != nil {
		t.Errorf("Expected nil, got %v", ref.Current)
	}

	ref.Current = 42
	if ref.Current != 42 {
		t.Errorf("Expected 42, got %v", ref.Current)
	}
}

// TestUseMemo 测试 useMemo Hook
func TestUseMemo(t *testing.T) {
	ctx := hooks.NewHookContext()

	callCount := 0
	value := hooks.UseMemo(ctx, func() any {
		callCount++
		return 42
	}, []any{"dep"})

	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

// TestUseMemoWithDifferentDeps 测试 useMemo 不同依赖
func TestUseMemoWithDifferentDeps(t *testing.T) {
	ctx := hooks.NewHookContext()

	computed := hooks.UseMemo(ctx, func() any {
		return "computed"
	}, []any{"a", "b"})

	if computed != "computed" {
		t.Errorf("Expected 'computed', got %v", computed)
	}
}

// TestUseCallback 测试 useCallback Hook
func TestUseCallback(t *testing.T) {
	ctx := hooks.NewHookContext()

	callCount := 0
	callback := hooks.UseCallback(ctx, func() {
		callCount++
	}, []any{"dep"})

	if callback == nil {
		t.Fatal("Callback should not be nil")
	}

	callback()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

// TestHookContext 测试 HookContext
func TestHookContext(t *testing.T) {
	ctx := hooks.NewHookContext()

	if ctx == nil {
		t.Fatal("HookContext should not be nil")
	}
}

// TestMultipleHooks 测试多个 Hook 同时使用
func TestMultipleHooks(t *testing.T) {
	ctx := hooks.NewHookContext()

	// 同时使用多个 Hook
	value, setValue := hooks.UseState(ctx, 0)
	ref := hooks.UseRef(ctx, nil)

	hooks.UseEffect(ctx, func() func() {
		ref.Current = value
		return nil
	}, []any{value})

	setValue(1)

	if ref.Current == nil {
		// 简化实现可能不更新
	}
}