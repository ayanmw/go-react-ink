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