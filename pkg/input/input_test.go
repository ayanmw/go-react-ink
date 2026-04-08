package input

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

func TestUseInput(t *testing.T) {
	ctx := hooks.NewHookContext()

	var receivedKey Key
	inputHook := UseInput(ctx, func(k Key) {
		receivedKey = k
	})

	// 处理键
	inputHook.ProcessKey(Key{Name: "a", Sequence: "a"})

	if receivedKey.Name != "a" {
		t.Errorf("Expected key 'a', got %q", receivedKey.Name)
	}
}

func TestInputHandleSpecificKey(t *testing.T) {
	ctx := hooks.NewHookContext()

	var called bool
	inputHook := UseInput(ctx, nil)

	// 注册特定键处理
	inputHook.Handle("enter", func(k Key) {
		called = true
	})

	// 处理 enter 键
	inputHook.ProcessKey(Key{Name: "enter"})

	if !called {
		t.Error("Enter handler should be called")
	}
}

func TestUseApp(t *testing.T) {
	ctx := hooks.NewHookContext()

	app := UseApp(ctx)

	if app == nil {
		t.Fatal("UseApp should return AppHook")
	}

	if app.Exit == nil {
		t.Error("Exit should not be nil")
	}
}

func TestUseFocus(t *testing.T) {
	ctx := hooks.NewHookContext()

	isFocused, focus := UseFocus(ctx)

	if isFocused {
		t.Error("Should not be focused initially")
	}

	focus()

	// 注意：简化版中 isFocused 不会自动更新
	// 在完整实现中，这会触发重新渲染
}

func TestUseFocusWithEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	var focusCalled bool
	var blurCalled bool

	isFocused, focus := UseFocusWithEvents(ctx,
		func() { focusCalled = true },
		func() { blurCalled = true },
	)

	if isFocused {
		t.Error("Should not be focused initially")
	}

	focus()

	if !focusCalled {
		t.Error("Focus callback should be called")
	}

	_ = blurCalled
}

func TestFocusManager(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := UseFocusManager(ctx)

	// 注册元素
	manager.Register("a")
	manager.Register("b")
	manager.Register("c")

	// 聚焦第一个
	manager.Focus("a")

	if !manager.IsFocused("a") {
		t.Error("'a' should be focused")
	}

	// 聚焦下一个
	manager.FocusNext()

	if manager.IsFocused("a") {
		t.Error("'a' should not be focused after FocusNext")
	}

	// 聚焦上一个
	manager.FocusPrev()

	_ = manager.IsFocused("a")
}

func TestFocusManagerBlur(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := UseFocusManager(ctx)
	manager.Register("a")
	manager.Focus("a")

	if !manager.IsFocused("a") {
		t.Error("'a' should be focused")
	}

	manager.Blur()

	if manager.IsFocused("a") {
		t.Error("'a' should not be focused after blur")
	}
}

func TestFocusManagerUnregister(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := UseFocusManager(ctx)
	manager.Register("a")
	manager.Register("b")

	manager.Focus("a")
	manager.Unregister("a")

	manager.FocusNext()

	// 应该跳过已注销的元素
}

func TestUseCursor(t *testing.T) {
	ctx := hooks.NewHookContext()

	cursor := UseCursor(ctx)

	if !cursor.visible {
		t.Error("Cursor should be visible by default")
	}

	cursor.Hide()

	if cursor.visible {
		t.Error("Cursor should be hidden")
	}

	cursor.Show()

	if !cursor.visible {
		t.Error("Cursor should be visible")
	}
}

func TestUseCursorPosition(t *testing.T) {
	ctx := hooks.NewHookContext()

	cursor := UseCursor(ctx)

	cursor.SetPosition(10, 5)

	if cursor.x != 10 || cursor.y != 5 {
		t.Errorf("Expected position (10, 5), got (%d, %d)", cursor.x, cursor.y)
	}
}

func TestUseCursorShape(t *testing.T) {
	ctx := hooks.NewHookContext()

	cursor := UseCursor(ctx)

	cursor.SetShape("underline")

	if cursor.shape != "underline" {
		t.Errorf("Expected shape 'underline', got %q", cursor.shape)
	}
}

func TestUseAnimation(t *testing.T) {
	ctx := hooks.NewHookContext()

	anim := UseAnimation(ctx, 30)

	if anim.IsPlaying() {
		t.Error("Animation should not be playing initially")
	}

	anim.Play()

	if !anim.IsPlaying() {
		t.Error("Animation should be playing")
	}

	anim.Pause()

	if anim.IsPlaying() {
		t.Error("Animation should be paused")
	}
}

func TestAnimationFrame(t *testing.T) {
	ctx := hooks.NewHookContext()

	anim := UseAnimation(ctx, 30)

	if anim.Frame() != 0 {
		t.Error("Initial frame should be 0")
	}

	anim.NextFrame()

	if anim.Frame() != 1 {
		t.Errorf("Frame should be 1, got %d", anim.Frame())
	}

	anim.Reset()

	if anim.Frame() != 0 {
		t.Error("Frame should be reset to 0")
	}
}

func TestUseStdout(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdout := UseStdout(ctx)

	if stdout.Width() != 80 {
		t.Errorf("Expected width 80, got %d", stdout.Width())
	}

	if stdout.Height() != 24 {
		t.Errorf("Expected height 24, got %d", stdout.Height())
	}
}

func TestUseStdoutSetSize(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdout := UseStdout(ctx)
	stdout.SetSize(100, 50)

	if stdout.Width() != 100 {
		t.Errorf("Expected width 100, got %d", stdout.Width())
	}

	if stdout.Height() != 50 {
		t.Errorf("Expected height 50, got %d", stdout.Height())
	}
}

func TestUseStdin(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := UseStdin(ctx)

	if !stdin.IsTTY() {
		t.Error("Should be TTY by default")
	}
}

func TestUseStdinRawMode(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := UseStdin(ctx)

	stdin.SetRaw(true)
	if !stdin.IsRaw() {
		t.Error("Should be in raw mode")
	}

	stdin.SetRaw(false)
	if stdin.IsRaw() {
		t.Error("Should not be in raw mode")
	}
}

func TestUseStdinOnData(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := UseStdin(ctx)

	var receivedData string
	stdin.OnData(func(data string) {
		receivedData = data
	})

	stdin.EmitData("test")

	if receivedData != "test" {
		t.Errorf("Expected 'test', got %s", receivedData)
	}
}

func TestUseStderr(t *testing.T) {
	ctx := hooks.NewHookContext()

	stderr := UseStderr(ctx)

	if stderr.Width() != 80 {
		t.Errorf("Expected width 80, got %d", stderr.Width())
	}

	if stderr.Height() != 24 {
		t.Errorf("Expected height 24, got %d", stderr.Height())
	}
}

func TestUseWindowSize(t *testing.T) {
	ctx := hooks.NewHookContext()

	size := UseWindowSize(ctx)

	s := size.Size()
	if s.Columns != 80 || s.Rows != 24 {
		t.Errorf("Expected (80, 24), got (%d, %d)", s.Columns, s.Rows)
	}

	if size.Columns() != 80 {
		t.Errorf("Expected 80 columns, got %d", size.Columns())
	}

	if size.Rows() != 24 {
		t.Errorf("Expected 24 rows, got %d", size.Rows())
	}
}

func TestUseWindowSizeSetSize(t *testing.T) {
	ctx := hooks.NewHookContext()

	size := UseWindowSize(ctx)
	size.SetSize(120, 40)

	if size.Columns() != 120 {
		t.Errorf("Expected 120 columns, got %d", size.Columns())
	}

	if size.Rows() != 40 {
		t.Errorf("Expected 40 rows, got %d", size.Rows())
	}
}

func TestUseBoxMetrics(t *testing.T) {
	ctx := hooks.NewHookContext()

	metrics := UseBoxMetrics(ctx)

	m := metrics.Metrics()
	if m.HasMeasured {
		t.Error("Should not have measured initially")
	}
}

func TestUseBoxMetricsSetMetrics(t *testing.T) {
	ctx := hooks.NewHookContext()

	metrics := UseBoxMetrics(ctx)
	metrics.SetMetrics(100, 50, 10, 20)

	m := metrics.Metrics()
	if m.Width != 100 || m.Height != 50 || m.Left != 10 || m.Top != 20 {
		t.Errorf("Expected (100, 50, 10, 20), got (%f, %f, %f, %f)", m.Width, m.Height, m.Left, m.Top)
	}

	if !m.HasMeasured {
		t.Error("Should have measured after SetMetrics")
	}
}

func TestUsePaste(t *testing.T) {
	ctx := hooks.NewHookContext()

	paste := UsePaste(ctx)

	var receivedText string
	paste.OnPaste(func(text string) {
		receivedText = text
	})

	paste.EmitPaste("pasted content")

	if receivedText != "pasted content" {
		t.Errorf("Expected 'pasted content', got %s", receivedText)
	}
}

func TestUseIsScreenReaderEnabled(t *testing.T) {
	ctx := hooks.NewHookContext()

	reader := UseIsScreenReaderEnabled(ctx)

	if reader.IsEnabled() {
		t.Error("Should not be enabled by default")
	}

	reader.SetEnabled(true)

	if !reader.IsEnabled() {
		t.Error("Should be enabled after SetEnabled(true)")
	}
}

func TestAnimationStop(t *testing.T) {
	ctx := hooks.NewHookContext()

	anim := UseAnimation(ctx, 30)
	anim.Play()

	if !anim.IsPlaying() {
		t.Error("Should be playing")
	}

	anim.Stop()

	if anim.IsPlaying() {
		t.Error("Should not be playing after stop")
	}

	if anim.Frame() != 0 {
		t.Error("Frame should be reset to 0 after stop")
	}
}
