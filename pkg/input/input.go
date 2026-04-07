// Package input 实现输入处理 Hooks
package input

import (
	"sync"

	"github.com/anmingwei/go-ink/pkg/hooks"
)

// Key 输入键
type Key struct {
	Name      string
	Sequence  string
	Shift     bool
	Ctrl      bool
	Alt       bool
	Meta      bool
}

// InputHook 输入 Hook
type InputHook struct {
	context   *hooks.HookContext
	handlers  map[string]func(Key) // key name -> handler
	allHandler func(Key)
	mu        sync.RWMutex
}

// UseInput 输入 Hook
func UseInput(ctx *hooks.HookContext, handler func(Key)) *InputHook {
	inputHook := &InputHook{
		context:    ctx,
		handlers:   make(map[string]func(Key)),
		allHandler: handler,
	}
	return inputHook
}

// Handle 注册特定键处理
func (h *InputHook) Handle(keyName string, handler func(Key)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[keyName] = handler
}

// ProcessKey 处理输入键
func (h *InputHook) ProcessKey(key Key) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 先检查特定键处理
	if handler, ok := h.handlers[key.Name]; ok {
		handler(key)
		return
	}

	// 使用全局处理
	if h.allHandler != nil {
		h.allHandler(key)
	}
}

// AppHook 应用 Hook
type AppHook struct {
	Exit   func()
	Stdin  interface{}
	Stdout interface{}
}

// UseApp 应用 Hook
func UseApp(ctx *hooks.HookContext) *AppHook {
	return &AppHook{
		Exit: func() {
			// 退出应用
		},
	}
}

// FocusHook 焦点 Hook
type FocusHook struct {
	isFocused bool
	context   *hooks.HookContext
	onFocus   func()
	onBlur    func()
}

// UseFocus 焦点 Hook
func UseFocus(ctx *hooks.HookContext) (isFocused bool, focus func()) {
	focusHook := &FocusHook{
		context: ctx,
	}

	focus = func() {
		if !focusHook.isFocused {
			focusHook.isFocused = true
			if focusHook.onFocus != nil {
				focusHook.onFocus()
			}
		}
	}

	return focusHook.isFocused, focus
}

// UseFocusWithEvents 带事件的焦点 Hook
func UseFocusWithEvents(ctx *hooks.HookContext, onFocus func(), onBlur func()) (isFocused bool, focus func()) {
	focusHook := &FocusHook{
		context: ctx,
		onFocus: onFocus,
		onBlur:  onBlur,
	}

	focus = func() {
		if !focusHook.isFocused {
			focusHook.isFocused = true
			if focusHook.onFocus != nil {
				focusHook.onFocus()
			}
		}
	}

	return focusHook.isFocused, focus
}

// FocusManager 焦点管理器
type FocusManager struct {
	focusables map[string]bool
	activeId   string
	mu         sync.RWMutex
}

// UseFocusManager 焦点管理器 Hook
func UseFocusManager(ctx *hooks.HookContext) *FocusManager {
	return &FocusManager{
		focusables: make(map[string]bool),
	}
}

// Focus 聚焦指定 ID
func (m *FocusManager) Focus(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeId = id
}

// Blur 取消焦点
func (m *FocusManager) Blur() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeId = ""
}

// FocusNext 聚焦下一个
func (m *FocusManager) FocusNext() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 简化版：循环焦点
	ids := make([]string, 0, len(m.focusables))
	for id := range m.focusables {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return
	}

	if m.activeId == "" {
		m.activeId = ids[0]
		return
	}

	for i, id := range ids {
		if id == m.activeId {
			next := (i + 1) % len(ids)
			m.activeId = ids[next]
			return
		}
	}
}

// FocusPrev 聚焦上一个
func (m *FocusManager) FocusPrev() {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := make([]string, 0, len(m.focusables))
	for id := range m.focusables {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return
	}

	if m.activeId == "" {
		m.activeId = ids[len(ids)-1]
		return
	}

	for i, id := range ids {
		if id == m.activeId {
			prev := (i - 1 + len(ids)) % len(ids)
			m.activeId = ids[prev]
			return
		}
	}
}

// Register 注册可聚焦元素
func (m *FocusManager) Register(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.focusables[id] = true
}

// Unregister 注销可聚焦元素
func (m *FocusManager) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.focusables, id)
}

// IsFocused 检查是否聚焦
func (m *FocusManager) IsFocused(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeId == id
}

// CursorHook 光标 Hook
type CursorHook struct {
	visible bool
	x, y    int
	shape   string // block, underline, bar
}

// UseCursor 光标 Hook
func UseCursor(ctx *hooks.HookContext) *CursorHook {
	return &CursorHook{
		visible: true,
		shape:   "block",
	}
}

// Show 显示光标
func (h *CursorHook) Show() {
	h.visible = true
}

// Hide 隐藏光标
func (h *CursorHook) Hide() {
	h.visible = false
}

// SetPosition 设置位置
func (h *CursorHook) SetPosition(x, y int) {
	h.x = x
	h.y = y
}

// SetShape 设置形状
func (h *CursorHook) SetShape(shape string) {
	h.shape = shape
}

// AnimationHook 动画 Hook
type AnimationHook struct {
	frame    int
	playing  bool
	fps      int
	context  *hooks.HookContext
}

// UseAnimation 动画 Hook
func UseAnimation(ctx *hooks.HookContext, fps int) *AnimationHook {
	return &AnimationHook{
		fps:     fps,
		context: ctx,
	}
}

// Play 播放
func (h *AnimationHook) Play() {
	h.playing = true
}

// Pause 暂停
func (h *AnimationHook) Pause() {
	h.playing = false
}

// Stop 停止
func (h *AnimationHook) Stop() {
	h.playing = false
	h.frame = 0
}

// NextFrame 下一帧
func (h *AnimationHook) NextFrame() {
	h.frame++
}

// Reset 重置
func (h *AnimationHook) Reset() {
	h.frame = 0
}

// Frame 获取当前帧
func (h *AnimationHook) Frame() int {
	return h.frame
}

// IsPlaying 是否播放中
func (h *AnimationHook) IsPlaying() bool {
	return h.playing
}

// StdoutHook 标准输出 Hook
type StdoutHook struct {
	width  int
	height int
}

// UseStdout 标准输出 Hook
func UseStdout(ctx *hooks.HookContext) *StdoutHook {
	return &StdoutHook{
		width:  80,
		height: 24,
	}
}

// Width 获取宽度
func (h *StdoutHook) Width() int {
	return h.width
}

// Height 获取高度
func (h *StdoutHook) Height() int {
	return h.height
}