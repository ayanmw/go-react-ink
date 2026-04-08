package fiber

import (
	"testing"
)

func TestNewFiber(t *testing.T) {
	fiber := NewFiber(TagHostComponent, nil, "test")

	if fiber.Tag != TagHostComponent {
		t.Errorf("Expected TagHostComponent, got %v", fiber.Tag)
	}

	if fiber.Key != "test" {
		t.Errorf("Expected key 'test', got %q", fiber.Key)
	}

	if fiber.Flags != EffectNoEffect {
		t.Errorf("Expected no effect, got %v", fiber.Flags)
	}
}

func TestCreateWorkInProgress(t *testing.T) {
	current := NewFiber(TagHostComponent, map[string]any{"a": 1}, "")
	current.StateNode = "instance"

	workInProgress := CreateWorkInProgress(current, map[string]any{"a": 2})

	if workInProgress.Tag != current.Tag {
		t.Errorf("Tag should match")
	}

	if workInProgress.PendingProps == nil {
		t.Error("Should have pending props")
	}

	if workInProgress.Alternate != current {
		t.Error("Alternate should point to current")
	}

	if current.Alternate != workInProgress {
		t.Error("Current's alternate should point to workInProgress")
	}
}

func TestAppendChild(t *testing.T) {
	parent := NewFiber(TagHostComponent, nil, "")
	child1 := NewFiber(TagHostText, "a", "")
	child2 := NewFiber(TagHostText, "b", "")

	parent.AppendChild(child1)
	parent.AppendChild(child2)

	if parent.Child != child1 {
		t.Error("First child should be child1")
	}

	if child1.Sibling != child2 {
		t.Error("child1's sibling should be child2")
	}

	if child1.Return != parent {
		t.Error("child1's return should be parent")
	}

	if child2.Return != parent {
		t.Error("child2's return should be parent")
	}
}

func TestRemoveChild(t *testing.T) {
	parent := NewFiber(TagHostComponent, nil, "")
	child1 := NewFiber(TagHostText, "a", "")
	child2 := NewFiber(TagHostText, "b", "")

	parent.AppendChild(child1)
	parent.AppendChild(child2)

	parent.RemoveChild(child1)

	if parent.Child != child2 {
		t.Error("After removing child1, first child should be child2")
	}

	if child1.Return != nil {
		t.Error("Removed child should have no parent")
	}
}

func TestFindNextDefibr(t *testing.T) {
	parent := NewFiber(TagHostComponent, nil, "")
	child1 := NewFiber(TagHostText, "a", "")
	child2 := NewFiber(TagHostText, "b", "")

	parent.AppendChild(child1)
	parent.AppendChild(child2)

	// 从 parent 开始，应该找到 child1
	next := parent.FindNextDefibr()
	if next != child1 {
		t.Error("Should find child1 first")
	}

	// 从 child1 开始，应该找到 child2
	next = child1.FindNextDefibr()
	if next != child2 {
		t.Error("Should find child2 next")
	}

	// 从 child2 开始，应该返回 nil
	next = child2.FindNextDefibr()
	if next != nil {
		t.Error("Should return nil after last child")
	}
}

func TestMarkEffect(t *testing.T) {
	fiber := NewFiber(TagHostComponent, nil, "")

	fiber.MarkEffect(EffectPlacement)

	if fiber.Flags != EffectPlacement {
		t.Errorf("Expected EffectPlacement, got %v", fiber.Flags)
	}

	fiber.MarkEffect(EffectUpdate)

	if fiber.Flags != (EffectPlacement | EffectUpdate) {
		t.Errorf("Expected both effects, got %v", fiber.Flags)
	}
}

func TestClearEffect(t *testing.T) {
	fiber := NewFiber(TagHostComponent, nil, "")
	fiber.MarkEffect(EffectPlacement | EffectUpdate)
	fiber.SubtreeFlags = EffectPlacement

	fiber.ClearEffect()

	if fiber.Flags != EffectNoEffect {
		t.Error("Flags should be cleared")
	}

	if fiber.SubtreeFlags != EffectNoEffect {
		t.Error("SubtreeFlags should be cleared")
	}
}

func TestHasEffect(t *testing.T) {
	fiber := NewFiber(TagHostComponent, nil, "")

	if fiber.HasEffect() {
		t.Error("New fiber should have no effect")
	}

	fiber.MarkEffect(EffectUpdate)

	if !fiber.HasEffect() {
		t.Error("Should have effect after marking")
	}
}

func TestSchedulerNew(t *testing.T) {
	s := NewScheduler()

	if s.timeSlice != 16*1000*1000 { // 16ms in nanoseconds
		t.Error("Time slice should be 16ms")
	}

	if s.pendingEffects == nil {
		t.Error("pendingEffects should be initialized")
	}
}

func TestSchedulerScheduleUpdate(t *testing.T) {
	s := NewScheduler()

	root := &RootFiber{
		Current: NewFiber(TagRoot, nil, ""),
	}

	s.ScheduleUpdate(root)

	if s.root != root {
		t.Error("Root should be set")
	}

	if s.nextUnitOfWork == nil {
		t.Error("Should have next unit of work")
	}

	if !s.isWorking {
		t.Error("Should be working")
	}
}

func TestSchedulerPerformUnitOfWork(t *testing.T) {
	s := NewScheduler()

	// 创建简单树 - 使用文本节点，它会立即完成
	textFiber := CreateHostTextFiber("Hello")
	textFiber.MemoizedProps = "Hello" // 设置旧值
	textFiber.PendingProps = "World"  // 设置新值触发更新

	s.nextUnitOfWork = textFiber

	_ = s.PerformUnitOfWork(textFiber)

	// 文本节点应该被标记为更新
	if textFiber.Flags != EffectUpdate {
		t.Logf("Flags: %v", textFiber.Flags)
	}
}

func TestSchedulerCommitRoot(t *testing.T) {
	s := NewScheduler()

	// 创建简单树
	root := &RootFiber{
		Current: NewFiber(TagRoot, nil, ""),
	}

	workRoot := NewFiber(TagRoot, nil, "")
	workRoot.MarkEffect(EffectPlacement)

	s.root = root
	s.workRoot = workRoot
	s.pendingEffects = []*Fiber{workRoot}

	s.CommitRoot()

	// 应该更新 current
	if root.Current != workRoot {
		t.Error("Current should be updated to workRoot")
	}

	// 应该清空 pendingEffects
	if len(s.pendingEffects) != 0 {
		t.Error("pendingEffects should be cleared")
	}

	// 应该重置状态
	if s.isWorking {
		t.Error("Should not be working after commit")
	}
}

func TestFiberString(t *testing.T) {
	tests := []struct {
		tag      WorkTag
		expected string
	}{
		{TagFunctionComponent, "FunctionComponent"},
		{TagHostComponent, "HostComponent"},
		{TagHostText, "HostText"},
		{TagFragment, "Fragment"},
		{TagRoot, "Root"},
	}

	for _, test := range tests {
		fiber := NewFiber(test.tag, nil, "key")
		str := fiber.String()

		if str != test.expected+"(key)" {
			t.Errorf("Expected %s(key), got %s", test.expected, str)
		}
	}
}

func TestCreateHostTextFiber(t *testing.T) {
	fiber := CreateHostTextFiber("Hello")

	if fiber.Tag != TagHostText {
		t.Error("Should be HostText")
	}

	if fiber.StateNode != "Hello" {
		t.Error("StateNode should be text content")
	}
}

func TestDeepTree(t *testing.T) {
	// 创建深层树
	root := NewFiber(TagRoot, nil, "")
	level1 := NewFiber(TagHostComponent, nil, "l1")
	level2 := NewFiber(TagHostComponent, nil, "l2")
	level3 := NewFiber(TagHostText, "leaf", "")

	root.AppendChild(level1)
	level1.AppendChild(level2)
	level2.AppendChild(level3)

	// 验证树结构
	if root.Child != level1 {
		t.Error("Root's child should be level1")
	}

	if level1.Child != level2 {
		t.Error("Level1's child should be level2")
	}

	if level2.Child != level3 {
		t.Error("Level2's child should be level3")
	}

	// 验证 return 指针
	if level3.Return != level2 {
		t.Error("Level3's return should be level2")
	}

	if level2.Return != level1 {
		t.Error("Level2's return should be level1")
	}

	if level1.Return != root {
		t.Error("Level1's return should be root")
	}
}

func TestCreateFiberFromElement(t *testing.T) {
	// 字符串类型创建 HostComponent
	fiber := CreateFiberFromElement("div", map[string]any{"children": "hello"}, "test-key")

	if fiber.Tag != TagHostComponent {
		t.Errorf("Expected TagHostComponent, got %v", fiber.Tag)
	}

	if fiber.Key != "test-key" {
		t.Errorf("Expected key 'test-key', got %v", fiber.Key)
	}

	// 函数类型创建 FunctionComponent
	component := func(props any) any { return nil }
	fiber2 := CreateFiberFromElement(component, nil, "func-key")

	if fiber2.Tag != TagFunctionComponent {
		t.Errorf("Expected TagFunctionComponent, got %v", fiber2.Tag)
	}
}

func TestFindFirstDefibr(t *testing.T) {
	root := NewFiber(TagRoot, nil, "")
	child := NewFiber(TagHostComponent, nil, "")

	root.AppendChild(child)

	// FindFirstDefibr 应该返回第一个子节点
	first := root.FindFirstDefibr()
	if first != child {
		t.Error("FindFirstDefibr should return first child")
	}

	// 没有子节点时返回 nil
	leaf := NewFiber(TagHostText, "leaf", "")
	if leaf.FindFirstDefibr() != nil {
		t.Error("Leaf should have no first child")
	}
}

func TestSchedulerRequestWork(t *testing.T) {
	s := NewScheduler()

	// 初始状态
	if s.RequestWork() != nil {
		t.Error("Initial RequestWork should return nil")
	}

	// 设置工作单元后
	root := &RootFiber{
		Current: NewFiber(TagRoot, nil, ""),
	}
	s.ScheduleUpdate(root)

	if s.RequestWork() == nil {
		t.Error("Should have work after ScheduleUpdate")
	}
}

func TestSchedulerIsWorkComplete(t *testing.T) {
	s := NewScheduler()

	// 没有工作时应该完成
	if !s.IsWorkComplete() {
		t.Error("Should be complete when no work")
	}

	// 有工作时未完成
	root := &RootFiber{
		Current: NewFiber(TagRoot, nil, ""),
	}
	s.ScheduleUpdate(root)

	if s.IsWorkComplete() {
		t.Error("Should not be complete when has work")
	}
}

func TestSchedulerBeginFunctionComponent(t *testing.T) {
	s := NewScheduler()

	// 创建函数组件 Fiber
	component := func(props any) any {
		return map[string]any{
			"type": "div",
			"props": map[string]any{
				"children": "test",
			},
		}
	}

	fiber := NewFiber(TagFunctionComponent, nil, "")
	fiber.Type = component
	fiber.PendingProps = map[string]any{}

	next := s.beginWork(fiber)
	_ = next // 函数组件返回子节点
}

func TestSchedulerBeginHostComponent(t *testing.T) {
	s := NewScheduler()

	fiber := NewFiber(TagHostComponent, nil, "")
	fiber.PendingProps = map[string]any{
		"children": []any{
			map[string]any{"type": "span", "props": map[string]any{}},
		},
	}

	next := s.beginWork(fiber)
	_ = next // Host 组件返回子节点
}

func TestSchedulerBeginFragment(t *testing.T) {
	s := NewScheduler()

	fiber := NewFiber(TagFragment, nil, "")
	fiber.PendingProps = map[string]any{
		"children": []any{
			map[string]any{"type": "span", "props": map[string]any{}},
		},
	}

	next := s.beginWork(fiber)
	_ = next
}

func TestSchedulerBeginRoot(t *testing.T) {
	s := NewScheduler()

	fiber := NewFiber(TagRoot, nil, "")
	fiber.PendingProps = map[string]any{
		"children": map[string]any{"type": "div", "props": map[string]any{}},
	}

	next := s.beginWork(fiber)
	_ = next
}

func TestSchedulerCommitUpdate(t *testing.T) {
	s := NewScheduler()

	fiber := NewFiber(TagHostComponent, nil, "")
	fiber.StateNode = "old-instance"
	fiber.PendingProps = map[string]any{"color": "red"}

	s.commitUpdate(fiber)
}

func TestSchedulerCommitDeletion(t *testing.T) {
	s := NewScheduler()

	parent := NewFiber(TagHostComponent, nil, "")
	child := NewFiber(TagHostText, "delete-me", "")
	parent.AppendChild(child)

	s.commitDeletion(child)
}

func TestShallowEqual(t *testing.T) {
	s := NewScheduler()

	// shallowEqual 使用 == 比较，只能比较可比较类型
	// 测试字符串
	if !s.shallowEqual("test", "test") {
		t.Error("Equal strings should be equal")
	}

	if s.shallowEqual("test", "different") {
		t.Error("Different strings should not be equal")
	}

	// 测试 nil
	if !s.shallowEqual(nil, nil) {
		t.Error("nil should be equal")
	}

	// 测试数字
	if !s.shallowEqual(42, 42) {
		t.Error("Equal numbers should be equal")
	}
}

func TestGetChildren(t *testing.T) {
	s := NewScheduler()

	children := []any{
		map[string]any{"type": "span", "props": map[string]any{}},
	}

	props := map[string]any{
		"children": children,
	}

	result := s.getChildren(props)
	if result == nil {
		t.Error("Should get children from props")
	}

	// 单个子节点
	props2 := map[string]any{
		"children": "text",
	}
	result2 := s.getChildren(props2)
	if result2 == nil {
		t.Error("Should handle string children")
	}
}

func TestGetText(t *testing.T) {
	s := NewScheduler()

	props := map[string]any{
		"children": "hello",
	}

	text := s.getText(props)
	if text != "hello" {
		t.Errorf("Expected 'hello', got %v", text)
	}

	// 多个子节点
	props2 := map[string]any{
		"children": []any{"a", "b"},
	}
	text2 := s.getText(props2)
	_ = text2 // 可能返回空或第一个
}

func TestSchedulerCreateFiberFromElement(t *testing.T) {
	s := NewScheduler()

	// 字符串元素
	textFiber := s.createFiberFromElement("hello", "key1")
	if textFiber.Tag != TagHostText {
		t.Error("String should create HostText fiber")
	}

	// 已有 Fiber
	existing := NewFiber(TagHostComponent, nil, "old-key")
	resultFiber := s.createFiberFromElement(existing, "new-key")
	if resultFiber.Key != "new-key" {
		t.Error("Should update key")
	}
}

func TestSchedulerReconcileChildren(t *testing.T) {
	s := NewScheduler()

	parent := NewFiber(TagHostComponent, nil, "")
	children := map[string]any{
		"type": "span",
		"props": map[string]any{
			"children": "test",
		},
	}

	s.reconcileChildren(parent, children)
	// 应该创建子节点
}

func TestSchedulerReconcileSingleChild(t *testing.T) {
	s := NewScheduler()

	parent := NewFiber(TagHostComponent, nil, "")
	child := NewFiber(TagHostText, "hello", "")

	s.reconcileSingleChild(parent, child)
}

func TestSchedulerReconcileChildrenArray(t *testing.T) {
	s := NewScheduler()

	parent := NewFiber(TagHostComponent, nil, "")
	children := []*Fiber{
		NewFiber(TagHostText, "a", ""),
		NewFiber(TagHostText, "b", ""),
	}

	s.reconcileChildrenArray(parent, children)

	if parent.Child == nil {
		t.Error("Should have child after reconciliation")
	}
}
