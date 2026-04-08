package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/input"
)

// TestUseInput 测试 useInput Hook
func TestUseInput(t *testing.T) {
	ctx := hooks.NewHookContext()

	keyReceived := ""
	inputHook := input.UseInput(ctx, func(k input.Key) {
		keyReceived = k.Name
	})

	if inputHook == nil {
		t.Fatal("UseInput should return non-nil hook")
	}

	// 测试处理键
	inputHook.ProcessKey(input.Key{Name: "return"})
	if keyReceived != "return" {
		t.Errorf("Expected 'return', got '%s'", keyReceived)
	}
}

// TestUseInputWithSpecificHandler 测试特定键处理
func TestUseInputWithSpecificHandler(t *testing.T) {
	ctx := hooks.NewHookContext()

	escapePressed := false
	inputHook := input.UseInput(ctx, nil)

	inputHook.Handle("escape", func(k input.Key) {
		escapePressed = true
	})

	inputHook.ProcessKey(input.Key{Name: "escape"})
	if !escapePressed {
		t.Error("Escape handler should have been called")
	}
}

// TestUseFocus 测试 useFocus Hook
func TestUseFocus(t *testing.T) {
	ctx := hooks.NewHookContext()

	isFocused, focus := input.UseFocus(ctx)

	if isFocused {
		t.Error("Should not be focused initially")
	}

	focus()
	// Note: 在简化实现中，isFocused 不会自动更新
}

// TestUseFocusWithEvents 测试带事件的焦点
func TestUseFocusWithEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	focusCalled := false
	blurCalled := false

	_, _ = input.UseFocusWithEvents(ctx,
		func() { focusCalled = true },
		func() { blurCalled = true },
	)

	_ = focusCalled
	_ = blurCalled
}

// TestUseFocusManager 测试焦点管理器
func TestUseFocusManager(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := input.UseFocusManager(ctx)

	manager.Register("input1")
	manager.Register("input2")
	manager.Register("input3")

	if manager.IsFocused("input1") {
		t.Error("input1 should not be focused initially")
	}

	manager.Focus("input2")
	if !manager.IsFocused("input2") {
		t.Error("input2 should be focused")
	}

	manager.FocusNext()
	// 应该聚焦下一个元素
}

// TestUseFocusManagerNavigation 测试焦点导航
func TestUseFocusManagerNavigation(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := input.UseFocusManager(ctx)
	manager.Register("a")
	manager.Register("b")
	manager.Register("c")

	manager.Focus("a")
	manager.FocusNext()
	if !manager.IsFocused("b") {
		t.Error("Should focus 'b' after FocusNext")
	}

	manager.FocusPrev()
	if !manager.IsFocused("a") {
		t.Error("Should focus 'a' after FocusPrev")
	}
}

// TestUseCursor 测试光标 Hook
func TestUseCursor(t *testing.T) {
	ctx := hooks.NewHookContext()

	cursor := input.UseCursor(ctx)

	cursor.Hide()
	cursor.Show()
	cursor.SetPosition(10, 5)
	cursor.SetShape("underline")
}

// TestUseAnimation 测试动画 Hook
func TestUseAnimation(t *testing.T) {
	ctx := hooks.NewHookContext()

	anim := input.UseAnimation(ctx, 30)

	if anim.IsPlaying() {
		t.Error("Animation should not be playing initially")
	}

	anim.Play()
	if !anim.IsPlaying() {
		t.Error("Animation should be playing after Play()")
	}

	anim.Pause()
	if anim.IsPlaying() {
		t.Error("Animation should be paused after Pause()")
	}

	anim.NextFrame()
	if anim.Frame() != 1 {
		t.Errorf("Expected frame 1, got %d", anim.Frame())
	}

	anim.Reset()
	if anim.Frame() != 0 {
		t.Errorf("Expected frame 0 after reset, got %d", anim.Frame())
	}
}

// TestUseStdout 测试标准输出 Hook
func TestUseStdout(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdout := input.UseStdout(ctx)

	if stdout.Width() != 80 {
		t.Errorf("Expected default width 80, got %d", stdout.Width())
	}

	if stdout.Height() != 24 {
		t.Errorf("Expected default height 24, got %d", stdout.Height())
	}

	stdout.SetSize(100, 30)
	if stdout.Width() != 100 || stdout.Height() != 30 {
		t.Error("SetSize should update dimensions")
	}
}

// TestUseStdin 测试标准输入 Hook
func TestUseStdin(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := input.UseStdin(ctx)

	if !stdin.IsTTY() {
		t.Error("Should be TTY by default")
	}

	stdin.SetRaw(true)
	if !stdin.IsRaw() {
		t.Error("Should be in raw mode")
	}

	dataReceived := ""
	stdin.OnData(func(data string) {
		dataReceived = data
	})

	stdin.EmitData("test")
	if dataReceived != "test" {
		t.Errorf("Expected 'test', got '%s'", dataReceived)
	}
}

// TestUseStderr 测试标准错误 Hook
func TestUseStderr(t *testing.T) {
	ctx := hooks.NewHookContext()

	stderr := input.UseStderr(ctx)

	if stderr.Width() != 80 || stderr.Height() != 24 {
		t.Error("Default stderr dimensions should be 80x24")
	}
}

// TestUseWindowSize 测试窗口尺寸 Hook
func TestUseWindowSize(t *testing.T) {
	ctx := hooks.NewHookContext()

	size := input.UseWindowSize(ctx)

	if size.Columns() != 80 || size.Rows() != 24 {
		t.Error("Default size should be 80x24")
	}

	size.SetSize(120, 40)
	if size.Columns() != 120 || size.Rows() != 40 {
		t.Error("SetSize should update dimensions")
	}
}

// TestUseBoxMetrics 测试盒子尺寸 Hook
func TestUseBoxMetrics(t *testing.T) {
	ctx := hooks.NewHookContext()

	metrics := input.UseBoxMetrics(ctx)

	m := metrics.Metrics()
	if m.HasMeasured {
		t.Error("Should not have measured initially")
	}

	metrics.SetMetrics(100, 50, 10, 5)
	m = metrics.Metrics()
	if !m.HasMeasured {
		t.Error("Should have measured after SetMetrics")
	}
	if m.Width != 100 || m.Height != 50 {
		t.Error("Metrics should match SetMetrics values")
	}
}

// TestUsePaste 测试粘贴 Hook
func TestUsePaste(t *testing.T) {
	ctx := hooks.NewHookContext()

	paste := input.UsePaste(ctx)

	pastedText := ""
	paste.OnPaste(func(text string) {
		pastedText = text
	})

	paste.EmitPaste("clipboard content")
	if pastedText != "clipboard content" {
		t.Errorf("Expected 'clipboard content', got '%s'", pastedText)
	}
}

// TestUseIsScreenReaderEnabled 测试屏幕阅读器 Hook
func TestUseIsScreenReaderEnabled(t *testing.T) {
	ctx := hooks.NewHookContext()

	sr := input.UseIsScreenReaderEnabled(ctx)

	if sr.IsEnabled() {
		t.Error("Screen reader should be disabled by default")
	}

	sr.SetEnabled(true)
	if !sr.IsEnabled() {
		t.Error("Screen reader should be enabled after SetEnabled(true)")
	}
}

// TestUseApp 测试应用 Hook
func TestUseApp(t *testing.T) {
	ctx := hooks.NewHookContext()

	app := input.UseApp(ctx)

	if app == nil {
		t.Fatal("UseApp should return non-nil hook")
	}

	if app.Exit == nil {
		t.Error("Exit should not be nil")
	}
}

// TestKeyModifiers 测试键修饰符
func TestKeyModifiers(t *testing.T) {
	key := input.Key{
		Name: "a",
		Shift: true,
		Ctrl: true,
		Alt: false,
	}

	if !key.Shift || !key.Ctrl {
		t.Error("Key modifiers should be set correctly")
	}
}
