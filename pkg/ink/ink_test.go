// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/core"
)

// MockElement 用于测试的 Mock 元素
type MockElement struct {
	content string
}

func (m *MockElement) Render() string {
	return m.content
}

func TestNewInk(t *testing.T) {
	opts := DefaultRenderOptions()
	ink := NewInk(opts)

	if ink == nil {
		t.Fatal("NewInk returned nil")
	}

	if ink.options == nil {
		t.Error("options should not be nil")
	}

	if ink.scheduler == nil {
		t.Error("scheduler should not be nil")
	}

	if ink.log == nil {
		t.Error("log should not be nil")
	}

	if ink.animationManager == nil {
		t.Error("animationManager should not be nil")
	}
}

func TestRenderOptionsDefaults(t *testing.T) {
	// Test nil options
	opts := applyDefaults(nil)
	if opts.MaxFps != 30 {
		t.Errorf("Expected MaxFps 30, got %d", opts.MaxFps)
	}

	if opts.ExitOnCtrlC != true {
		t.Error("Expected ExitOnCtrlC to be true")
	}

	if opts.Interactive != true {
		t.Error("Expected Interactive to be true")
	}

	// Test custom options with some defaults
	custom := &RenderOptions{
		MaxFps: 60,
	}
	opts = applyDefaults(custom)
	if opts.MaxFps != 60 {
		t.Errorf("Expected MaxFps 60, got %d", opts.MaxFps)
	}
}

func TestFrameDuration(t *testing.T) {
	opts := &RenderOptions{MaxFps: 30}
	duration := opts.FrameDuration()
	expected := time.Second / 30

	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}
}

func TestLogUpdate(t *testing.T) {
	var buf bytes.Buffer
	log := NewLogUpdate(&buf)

	// Test first write
	err := log.Write("Hello\n")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if log.lastHeight != 1 {
		t.Errorf("Expected lastHeight 1, got %d", log.lastHeight)
	}

	// Test second write (should clear previous)
	err = log.Write("World\n")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check that output contains World
	output := buf.String()
	if !strings.Contains(output, "World") {
		t.Error("Expected output to contain 'World'")
	}
}

func TestLogUpdateClear(t *testing.T) {
	var buf bytes.Buffer
	log := NewLogUpdate(&buf)

	// Write something
	log.Write("Test\n")

	// Clear
	err := log.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if log.lastHeight != 0 {
		t.Errorf("Expected lastHeight 0 after clear, got %d", log.lastHeight)
	}
}

func TestAnimationManager(t *testing.T) {
	manager := NewAnimationManager()

	if manager == nil {
		t.Fatal("NewAnimationManager returned nil")
	}

	// Test HasActiveAnimations with no subscribers
	if manager.HasActiveAnimations() {
		t.Error("Expected no active animations")
	}

	// Add a subscriber
	sub := &AnimationSubscriber{
		Fps:       10,
		IsPlaying: true,
		LastTick:  time.Now(),
	}
	manager.Subscribe(sub)

	// Test HasActiveAnimations with active subscriber
	if !manager.HasActiveAnimations() {
		t.Error("Expected active animations")
	}

	// Test Tick
	manager.Tick()

	// Unsubscribe
	manager.Unsubscribe(sub)

	if manager.HasActiveAnimations() {
		t.Error("Expected no active animations after unsubscribe")
	}
}

func TestAnimationSubscriberTick(t *testing.T) {
	manager := NewAnimationManager()

	callbackCalled := false
	sub := &AnimationSubscriber{
		Fps:       1000, // Very high FPS for testing
		IsPlaying: true,
		LastTick:  time.Now().Add(-time.Second), // Set to past so tick will fire
		Callback: func(frame int) {
			callbackCalled = true
		},
	}

	manager.Subscribe(sub)

	// Wait a bit and tick
	time.Sleep(time.Millisecond)
	manager.Tick()

	if !callbackCalled {
		t.Error("Expected callback to be called")
	}

	if sub.Frame != 1 {
		t.Errorf("Expected Frame 1, got %d", sub.Frame)
	}
}

func TestInkRender(t *testing.T) {
	var buf bytes.Buffer
	opts := &RenderOptions{
		Stdout:      &buf,
		MaxFps:      30,
		Interactive: false, // Disable for testing
	}

	ink := NewInk(opts)

	// Create mock element
	element := &MockElement{content: "Test Output"}

	// Render
	ink.Render(element)

	// Check that element was stored
	if ink.rootNode == nil {
		t.Error("Expected rootNode to be set")
	}
}

func TestInkUnmount(t *testing.T) {
	var buf bytes.Buffer
	opts := &RenderOptions{
		Stdout:      &buf,
		MaxFps:      30,
		Interactive: false,
	}

	ink := NewInk(opts)

	// Unmount should not panic
	ink.Unmount()

	if !ink.isUnmounted {
		t.Error("Expected isUnmounted to be true")
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"hello", 1},
		{"hello\n", 1},
		{"hello\nworld", 2},
		{"hello\nworld\n", 2},
		{"a\nb\nc", 3},
	}

	for _, test := range tests {
		result := countLines(test.input)
		if result != test.expected {
			t.Errorf("countLines(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestANSISequences(t *testing.T) {
	// Test that ANSI constants are defined
	if ANSIHideCursor == "" {
		t.Error("ANSIHideCursor should not be empty")
	}

	if ANSIShowCursor == "" {
		t.Error("ANSIShowCursor should not be empty")
	}

	if ANSIClearScreen == "" {
		t.Error("ANSIClearScreen should not be empty")
	}

	// Test MoveCursorUp
	result := MoveCursorUp(1)
	if result == "" {
		t.Error("MoveCursorUp(1) should not be empty")
	}

	// Test MoveCursorDown
	result = MoveCursorDown(1)
	if result == "" {
		t.Error("MoveCursorDown(1) should not be empty")
	}

	// Test MoveCursorTo
	result = MoveCursorTo(1, 1)
	if result == "" {
		t.Error("MoveCursorTo(1, 1) should not be empty")
	}
}

// Test that core.Element interface is satisfied
var _ core.Element = &MockElement{}