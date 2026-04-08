package main

import (
	"testing"

	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/input"
)

// TestKeyEvent 测试键盘事件
func TestKeyEvent(t *testing.T) {
	ctx := hooks.NewHookContext()

	pressedKeys := []string{}
	inputHook := input.UseInput(ctx, func(k input.Key) {
		pressedKeys = append(pressedKeys, k.Name)
	})

	// 模拟按键
	inputHook.ProcessKey(input.Key{Name: "a"})
	inputHook.ProcessKey(input.Key{Name: "b"})
	inputHook.ProcessKey(input.Key{Name: "c"})

	if len(pressedKeys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(pressedKeys))
	}
}

// TestKeyModifiers 测试键修饰符事件
func TestKeyModifiers(t *testing.T) {
	ctx := hooks.NewHookContext()

	ctrlPressed := false
	inputHook := input.UseInput(ctx, func(k input.Key) {
		if k.Ctrl && k.Name == "c" {
			ctrlPressed = true
		}
	})

	inputHook.ProcessKey(input.Key{Name: "c", Ctrl: true})
	if !ctrlPressed {
		t.Error("Ctrl+C should be detected")
	}
}

// TestSpecialKeys 测试特殊键
func TestSpecialKeys(t *testing.T) {
	specialKeys := []string{
		"return", "escape", "tab", "backspace", "delete",
		"up", "down", "left", "right",
		"home", "end", "pageup", "pagedown",
	}

	ctx := hooks.NewHookContext()
	receivedKeys := []string{}

	inputHook := input.UseInput(ctx, func(k input.Key) {
		receivedKeys = append(receivedKeys, k.Name)
	})

	for _, key := range specialKeys {
		inputHook.ProcessKey(input.Key{Name: key})
	}

	if len(receivedKeys) != len(specialKeys) {
		t.Errorf("Expected %d special keys, got %d", len(specialKeys), len(receivedKeys))
	}
}

// TestFocusEvents 测试焦点事件
func TestFocusEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	focusCount := 0
	blurCount := 0

	_, _ = input.UseFocusWithEvents(ctx,
		func() { focusCount++ },
		func() { blurCount++ },
	)

	// 焦点事件通过 focus 函数触发
	_ = focusCount
	_ = blurCount
}

// TestFocusManagerEvents 测试焦点管理器事件
func TestFocusManagerEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	manager := input.UseFocusManager(ctx)
	manager.Register("input1")
	manager.Register("input2")

	// 聚焦第一个
	manager.Focus("input1")
	if !manager.IsFocused("input1") {
		t.Error("input1 should be focused")
	}

	// Tab 到下一个
	manager.FocusNext()
	if !manager.IsFocused("input2") {
		t.Error("input2 should be focused after FocusNext")
	}

	// Shift+Tab 到上一个
	manager.FocusPrev()
	if !manager.IsFocused("input1") {
		t.Error("input1 should be focused after FocusPrev")
	}
}

// TestMultipleInputHandlers 测试多个输入处理器
func TestMultipleInputHandlers(t *testing.T) {
	ctx := hooks.NewHookContext()

	results := []string{}

	// 全局处理器
	inputHook := input.UseInput(ctx, func(k input.Key) {
		results = append(results, "global:"+k.Name)
	})

	// 特定键处理器
	inputHook.Handle("escape", func(k input.Key) {
		results = append(results, "escape-handler")
	})

	inputHook.Handle("return", func(k input.Key) {
		results = append(results, "return-handler")
	})

	inputHook.ProcessKey(input.Key{Name: "escape"})
	inputHook.ProcessKey(input.Key{Name: "a"})
	inputHook.ProcessKey(input.Key{Name: "return"})

	expected := []string{"escape-handler", "global:a", "return-handler"}
	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Expected '%s', got '%s'", exp, results[i])
		}
	}
}

// TestInputSequence 测试输入序列
func TestInputSequence(t *testing.T) {
	ctx := hooks.NewHookContext()

	sequence := ""
	inputHook := input.UseInput(ctx, func(k input.Key) {
		sequence += k.Name
	})

	// 模拟输入序列
	keys := []input.Key{
		{Name: "h"}, {Name: "e"}, {Name: "l"}, {Name: "l"}, {Name: "o"},
	}

	for _, k := range keys {
		inputHook.ProcessKey(k)
	}

	if sequence != "hello" {
		t.Errorf("Expected 'hello', got '%s'", sequence)
	}
}

// TestStdinDataEvents 测试标准输入数据事件
func TestStdinEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := input.UseStdin(ctx)

	receivedData := ""
	stdin.OnData(func(data string) {
		receivedData += data
	})

	stdin.EmitData("line1\n")
	stdin.EmitData("line2\n")

	if receivedData != "line1\nline2\n" {
		t.Errorf("Expected 'line1\\nline2\\n', got '%s'", receivedData)
	}
}

// TestPasteEvents 测试粘贴事件
func TestPasteEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	paste := input.UsePaste(ctx)

	pastedContent := ""
	paste.OnPaste(func(text string) {
		pastedContent = text
	})

	paste.EmitPaste("multi\nline\ntext")

	if pastedContent != "multi\nline\ntext" {
		t.Errorf("Expected 'multi\\nline\\ntext', got '%s'", pastedContent)
	}
}

// TestWindowSizeEvents 测试窗口尺寸事件
func TestWindowSizeEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	size := input.UseWindowSize(ctx)

	initialCols := size.Columns()
	initialRows := size.Rows()

	size.SetSize(120, 40)

	if size.Columns() == initialCols || size.Rows() == initialRows {
		t.Error("Size should change after SetSize")
	}
}

// TestAnimationEvents 测试动画事件
func TestAnimationEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	anim := input.UseAnimation(ctx, 60)

	anim.Play()
	if !anim.IsPlaying() {
		t.Error("Animation should be playing")
	}

	// 模拟帧更新
	anim.NextFrame()
	anim.NextFrame()
	anim.NextFrame()

	if anim.Frame() != 3 {
		t.Errorf("Expected frame 3, got %d", anim.Frame())
	}

	anim.Stop()
	if anim.IsPlaying() {
		t.Error("Animation should stop after Stop()")
	}
	if anim.Frame() != 0 {
		t.Error("Frame should reset to 0 after Stop()")
	}
}

// TestCursorEvents 测试光标事件
func TestCursorEvents(t *testing.T) {
	ctx := hooks.NewHookContext()

	cursor := input.UseCursor(ctx)

	// 显示/隐藏
	cursor.Hide()
	cursor.Show()

	// 移动
	cursor.SetPosition(10, 5)

	// 改变形状
	cursor.SetShape("block")
	cursor.SetShape("underline")
	cursor.SetShape("bar")
}

// TestMultipleCallbacks 测试多个回调
func TestMultipleCallbacks(t *testing.T) {
	ctx := hooks.NewHookContext()

	stdin := input.UseStdin(ctx)

	callCount := 0
	stdin.OnData(func(data string) {
		callCount++
	})
	stdin.OnData(func(data string) {
		callCount++
	})
	stdin.OnData(func(data string) {
		callCount++
	})

	stdin.EmitData("test")

	if callCount != 3 {
		t.Errorf("Expected 3 callbacks, got %d", callCount)
	}
}

// TestKeySequence 测试键序列 (快捷键)
func TestKeySequence(t *testing.T) {
	ctx := hooks.NewHookContext()

	quitDetected := false
	inputHook := input.UseInput(ctx, nil)

	// Ctrl+Q 退出
	inputHook.Handle("q", func(k input.Key) {
		if k.Ctrl {
			quitDetected = true
		}
	})

	inputHook.ProcessKey(input.Key{Name: "q", Ctrl: true})
	if !quitDetected {
		t.Error("Ctrl+Q should trigger quit")
	}
}

// TestEventPropagation 测试事件传播
func TestEventPropagation(t *testing.T) {
	ctx := hooks.NewHookContext()

	handlerOrder := []string{}

	// 先注册全局处理器
	inputHook := input.UseInput(ctx, func(k input.Key) {
		handlerOrder = append(handlerOrder, "global")
	})

	// 再注册特定键处理器
	inputHook.Handle("a", func(k input.Key) {
		handlerOrder = append(handlerOrder, "specific")
	})

	// 特定键处理器优先
	inputHook.ProcessKey(input.Key{Name: "a"})
	if handlerOrder[0] != "specific" {
		t.Error("Specific handler should be called first")
	}

	// 其他键使用全局处理器
	handlerOrder = []string{}
	inputHook.ProcessKey(input.Key{Name: "b"})
	if handlerOrder[0] != "global" {
		t.Error("Global handler should handle non-specific keys")
	}
}

// TestAppExit 测试应用退出事件
func TestAppExit(t *testing.T) {
	ctx := hooks.NewHookContext()

	app := input.UseApp(ctx)

	// Exit 函数应该可用
	if app.Exit == nil {
		t.Error("Exit function should not be nil")
	}

	// 调用退出 (不会真的退出)
	app.Exit()
}
