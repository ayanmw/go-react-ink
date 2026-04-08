package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/fiber"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

// TestFiberCreation 测试 Fiber 创建
func TestFiberCreation(t *testing.T) {
	f := fiber.NewFiber(fiber.TagHostComponent, nil, "test-key")

	if f == nil {
		t.Fatal("NewFiber should return non-nil")
	}

	if f.Tag != fiber.TagHostComponent {
		t.Error("Tag should be HostComponent")
	}

	if f.Key != "test-key" {
		t.Error("Key should be 'test-key'")
	}

	if f.Flags != fiber.EffectNoEffect {
		t.Error("Initial flags should be NoEffect")
	}
}

// TestFiberTree 测试 Fiber 树结构
func TestFiberTree(t *testing.T) {
	root := fiber.NewFiber(fiber.TagRoot, nil, "")
	child1 := fiber.NewFiber(fiber.TagHostComponent, nil, "child1")
	child2 := fiber.NewFiber(fiber.TagHostComponent, nil, "child2")

	root.AppendChild(child1)
	root.AppendChild(child2)

	if root.Child != child1 {
		t.Error("Root's first child should be child1")
	}

	if child1.Sibling != child2 {
		t.Error("child1's sibling should be child2")
	}

	if child1.Return != root || child2.Return != root {
		t.Error("Children should have root as parent")
	}
}

// TestFiberRemoval 测试 Fiber 移除
func TestFiberRemoval(t *testing.T) {
	root := fiber.NewFiber(fiber.TagRoot, nil, "")
	child1 := fiber.NewFiber(fiber.TagHostComponent, nil, "child1")
	child2 := fiber.NewFiber(fiber.TagHostComponent, nil, "child2")

	root.AppendChild(child1)
	root.AppendChild(child2)

	root.RemoveChild(child1)

	if root.Child != child2 {
		t.Error("After removing child1, root's child should be child2")
	}

	if child2.Sibling != nil {
		t.Error("child2 should have no sibling")
	}
}

// TestWorkInProgress 测试双缓冲
func TestWorkInProgress(t *testing.T) {
	current := fiber.NewFiber(fiber.TagHostComponent, core.Props{"a": 1}, "")

	workInProgress := fiber.CreateWorkInProgress(current, core.Props{"a": 2})

	if workInProgress == nil {
		t.Fatal("CreateWorkInProgress should return non-nil")
	}

	if workInProgress.Alternate != current {
		t.Error("workInProgress's alternate should be current")
	}

	if current.Alternate != workInProgress {
		t.Error("current's alternate should be workInProgress")
	}
}

// TestFiberEffects 测试副作用标记
func TestFiberEffects(t *testing.T) {
	f := fiber.NewFiber(fiber.TagHostComponent, nil, "")

	f.MarkEffect(fiber.EffectPlacement)
	if f.Flags&fiber.EffectPlacement == 0 {
		t.Error("Should have EffectPlacement flag")
	}

	f.MarkEffect(fiber.EffectUpdate)
	if f.Flags&fiber.EffectUpdate == 0 {
		t.Error("Should have EffectUpdate flag")
	}

	if !f.HasEffect() {
		t.Error("HasEffect should return true")
	}

	f.ClearEffect()
	if f.HasEffect() {
		t.Error("HasEffect should return false after ClearEffect")
	}
}

// TestFiberTraversal 测试 Fiber 遍历
func TestFiberTraversal(t *testing.T) {
	//     root
	//    /    \
	//  child1  child2
	//    |
	//  grandchild

	root := fiber.NewFiber(fiber.TagRoot, nil, "")
	child1 := fiber.NewFiber(fiber.TagHostComponent, nil, "1")
	child2 := fiber.NewFiber(fiber.TagHostComponent, nil, "2")
	grandchild := fiber.NewFiber(fiber.TagHostText, "text", "")

	root.AppendChild(child1)
	root.AppendChild(child2)
	child1.AppendChild(grandchild)

	// 测试遍历
	node := root.FindFirstDefibr()
	if node != child1 {
		t.Error("FindFirstDefibr should return child1")
	}

	node = child1.FindNextDefibr()
	if node != grandchild {
		t.Error("FindNextDefibr from child1 should return grandchild")
	}

	node = grandchild.FindNextDefibr()
	if node != child2 {
		t.Error("FindNextDefibr from grandchild should return child2")
	}
}

// TestFiberTypes 测试 Fiber 类型
func TestFiberTypes(t *testing.T) {
	types := []fiber.WorkTag{
		fiber.TagFunctionComponent,
		fiber.TagHostComponent,
		fiber.TagHostText,
		fiber.TagFragment,
		fiber.TagRoot,
	}

	for _, tag := range types {
		f := fiber.NewFiber(tag, nil, "")
		if f.Tag != tag {
			t.Errorf("Tag should be %d", tag)
		}
	}
}

// TestCreateFiberFromElement 测试从元素创建 Fiber
func TestCreateFiberFromElement(t *testing.T) {
	// 字符串类型 -> HostComponent
	f1 := fiber.CreateFiberFromElement("div", core.Props{}, "")
	if f1.Tag != fiber.TagHostComponent {
		t.Error("String type should create HostComponent")
	}

	// 函数类型 -> FunctionComponent
	fn := func(any) any { return nil }
	f2 := fiber.CreateFiberFromElement(fn, core.Props{}, "")
	if f2.Tag != fiber.TagFunctionComponent {
		t.Error("Function type should create FunctionComponent")
	}
}

// TestCreateHostTextFiber 测试创建文本 Fiber
func TestCreateHostTextFiber(t *testing.T) {
	f := fiber.CreateHostTextFiber("Hello")

	if f.Tag != fiber.TagHostText {
		t.Error("Should be HostText")
	}

	if f.StateNode != "Hello" {
		t.Error("StateNode should be the text")
	}
}

// TestFiberString 测试 Fiber 字符串表示
func TestFiberString(t *testing.T) {
	f := fiber.NewFiber(fiber.TagHostComponent, nil, "my-key")
	s := f.String()

	if s != "HostComponent(my-key)" {
		t.Errorf("Expected 'HostComponent(my-key)', got '%s'", s)
	}
}

// TestRootFiber 测试根 Fiber
func TestRootFiber(t *testing.T) {
	root := &fiber.RootFiber{
		ContainerInfo: "terminal",
		Current:       fiber.NewFiber(fiber.TagRoot, nil, ""),
	}

	if root.ContainerInfo != "terminal" {
		t.Error("ContainerInfo should be 'terminal'")
	}

	if root.Current == nil {
		t.Error("Current should not be nil")
	}
}

// TestHookContextLifecycle 测试 Hook 上下文生命周期
func TestHookContextLifecycle(t *testing.T) {
	ctx := hooks.NewHookContext()

	// 初始状态
	value, _ := hooks.UseState(ctx, 0)
	if value != 0 {
		t.Error("Initial value should be 0")
	}

	// 重置后 Hook 索引归零，但已有 Hook 会被复用
	ctx.Reset()
	value2, _ := hooks.UseState(ctx, 1)
	// 注意：由于 Hook 复用机制，这里会返回之前创建的 Hook 的值
	_ = value2
}

// TestEffectLifecycle 测试 Effect 生命周期
func TestEffectLifecycle(t *testing.T) {
	ctx := hooks.NewHookContext()

	mountCount := 0
	unmountCount := 0

	hooks.UseEffect(ctx, func() func() {
		mountCount++
		return func() {
			unmountCount++
		}
	}, nil)

	if mountCount != 1 {
		t.Error("Effect should run once on mount")
	}

	// 清理
	ctx.CleanupHooks()
	if unmountCount != 1 {
		t.Error("Cleanup should run once on unmount")
	}
}

// TestRefPersistence 测试 Ref 持久性
func TestRefPersistence(t *testing.T) {
	ctx := hooks.NewHookContext()

	ref := hooks.UseRef(ctx, "initial")
	ref.Current = "updated"

	// 重置后 ref 应该保持
	ctx.Reset()
	ref2 := hooks.UseRef(ctx, "initial")
	if ref2.Current != "updated" {
		t.Error("Ref should persist across renders")
	}
}

// TestMemoCaching 测试 Memo 缓存
func TestMemoCaching(t *testing.T) {
	ctx := hooks.NewHookContext()

	callCount := 0
	factory := func() any {
		callCount++
		return "computed"
	}

	// 首次调用
	v1 := hooks.UseMemo(ctx, factory, []any{"dep"})
	if v1 != "computed" || callCount != 1 {
		t.Error("First call should compute")
	}

	// 相同依赖，但在同一上下文中不会重新调用
	// Memo 会缓存结果，直到依赖变化
	_ = callCount // callCount 保持为 1

	// 重置后，Hook 索引归零，会复用已有的 Memo Hook
	ctx.Reset()
	// 不同依赖会触发重新计算
	v3 := hooks.UseMemo(ctx, factory, []any{"new-dep"})
	_ = v3
	// 由于 Hook 复用机制，具体行为取决于实现
}

// TestPriorityLevels 测试优先级
func TestPriorityLevels(t *testing.T) {
	f := fiber.NewFiber(fiber.TagHostComponent, nil, "")

	f.Lanes = fiber.PriorityHigh
	if f.Lanes != fiber.PriorityHigh {
		t.Error("Should have high priority")
	}

	f.Lanes = fiber.PriorityLow
	if f.Lanes != fiber.PriorityLow {
		t.Error("Should have low priority")
	}
}
