// Package input 实现输入处理 Hooks
package input

import (
	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"sync"
)

// Key 输入键
type Key struct {
	Name     string
	Sequence string
	Shift    bool
	Ctrl     bool
	Alt      bool
	Meta     bool
}

// Hook 输入 Hook
type Hook struct {
	context    *hooks.HookContext
	handlers   map[string]func(Key) // key name -> handler
	allHandler func(Key)
	mu         sync.RWMutex
}

// UseInput 输入 Hook
func UseInput(ctx *hooks.HookContext, handler func(Key)) *Hook {
	inputHook := &Hook{
		context:    ctx,
		handlers:   make(map[string]func(Key)),
		allHandler: handler,
	}
	return inputHook
}

// Handle 注册特定键处理
func (h *Hook) Handle(keyName string, handler func(Key)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[keyName] = handler
}

// ProcessKey 处理输入键
func (h *Hook) ProcessKey(key Key) {
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

// AppHook 应用 Hook (deprecated: use AppController from app.go)
type AppHook struct {
	Exit   func()
	Stdin  interface{}
	Stdout interface{}
}

// UseAppLegacy 应用 Hook (deprecated: use UseApp from app.go)
func UseAppLegacy(ctx *hooks.HookContext) *AppHook {
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
	activeID   string
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
	m.activeID = id
}

// Blur 取消焦点
func (m *FocusManager) Blur() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeID = ""
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

	if m.activeID == "" {
		m.activeID = ids[0]
		return
	}

	for i, id := range ids {
		if id == m.activeID {
			next := (i + 1) % len(ids)
			m.activeID = ids[next]
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

	if m.activeID == "" {
		m.activeID = ids[len(ids)-1]
		return
	}

	for i, id := range ids {
		if id == m.activeID {
			prev := (i - 1 + len(ids)) % len(ids)
			m.activeID = ids[prev]
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
	return m.activeID == id
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
	frame   int
	playing bool
	fps     int
	context *hooks.HookContext
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

// Width 返回宽度
func (h *StdoutHook) Width() int {
	return h.width
}

// Height 返回高度
func (h *StdoutHook) Height() int {
	return h.height
}

// SetSize 设置尺寸
func (h *StdoutHook) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// StdinHook 标准输入 Hook
type StdinHook struct {
	context   *hooks.HookContext
	isTTY     bool
	isRaw     bool
	callbacks []func(string)
	mu        sync.RWMutex
}

// UseStdin 标准输入 Hook
func UseStdin(ctx *hooks.HookContext) *StdinHook {
	return &StdinHook{
		context: ctx,
		isTTY:   true,
		isRaw:   false,
	}
}

// IsTTY 返回是否为 TTY
func (h *StdinHook) IsTTY() bool {
	return h.isTTY
}

// SetRaw 设置原始模式
func (h *StdinHook) SetRaw(raw bool) {
	h.isRaw = raw
}

// IsRaw 返回是否为原始模式
func (h *StdinHook) IsRaw() bool {
	return h.isRaw
}

// OnData 注册数据回调
func (h *StdinHook) OnData(callback func(string)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.callbacks = append(h.callbacks, callback)
}

// EmitData 触发数据
func (h *StdinHook) EmitData(data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, cb := range h.callbacks {
		cb(data)
	}
}

// StderrHook 标准错误 Hook
type StderrHook struct {
	width  int
	height int
}

// UseStderr 标准错误 Hook
func UseStderr(ctx *hooks.HookContext) *StderrHook {
	return &StderrHook{
		width:  80,
		height: 24,
	}
}

// Width 返回宽度
func (h *StderrHook) Width() int {
	return h.width
}

// Height 返回高度
func (h *StderrHook) Height() int {
	return h.height
}

// WindowSize 窗口尺寸
type WindowSize struct {
	Columns int
	Rows    int
}

// WindowSizeHook 窗口尺寸 Hook
type WindowSizeHook struct {
	size    WindowSize
	context *hooks.HookContext
	mu      sync.RWMutex
}

// UseWindowSize 窗口尺寸 Hook
func UseWindowSize(ctx *hooks.HookContext) *WindowSizeHook {
	return &WindowSizeHook{
		size: WindowSize{
			Columns: 80,
			Rows:    24,
		},
		context: ctx,
	}
}

// Size 返回尺寸
func (h *WindowSizeHook) Size() WindowSize {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.size
}

// SetSize 设置尺寸
func (h *WindowSizeHook) SetSize(columns, rows int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.size.Columns = columns
	h.size.Rows = rows
}

// Columns 返回列数
func (h *WindowSizeHook) Columns() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.size.Columns
}

// Rows 返回行数
func (h *WindowSizeHook) Rows() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.size.Rows
}

// BoxMetrics 盒子尺寸信息
type BoxMetrics struct {
	Width       float64
	Height      float64
	Left        float64
	Top         float64
	HasMeasured bool
}

// BoxMetricsHook 盒子尺寸 Hook
type BoxMetricsHook struct {
	metrics BoxMetrics
	context *hooks.HookContext
	mu      sync.RWMutex
}

// UseBoxMetrics 盒子尺寸 Hook
// 用于获取元素的布局信息 (类似 React Ink 的 useBoxMetrics)
func UseBoxMetrics(ctx *hooks.HookContext) *BoxMetricsHook {
	return &BoxMetricsHook{
		metrics: BoxMetrics{
			Width:       0,
			Height:      0,
			Left:        0,
			Top:         0,
			HasMeasured: false,
		},
		context: ctx,
	}
}

// Metrics 返回尺寸信息
func (h *BoxMetricsHook) Metrics() BoxMetrics {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.metrics
}

// SetMetrics 设置尺寸信息
func (h *BoxMetricsHook) SetMetrics(width, height, left, top float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.metrics.Width = width
	h.metrics.Height = height
	h.metrics.Left = left
	h.metrics.Top = top
	h.metrics.HasMeasured = true
}

// SetHasMeasured 设置是否已测量
func (h *BoxMetricsHook) SetHasMeasured(measured bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.metrics.HasMeasured = measured
}

// PasteHook 粘贴 Hook
type PasteHook struct {
	context   *hooks.HookContext
	callbacks []func(string)
	mu        sync.RWMutex
}

// UsePaste 粘贴 Hook
func UsePaste(ctx *hooks.HookContext) *PasteHook {
	return &PasteHook{
		context: ctx,
	}
}

// OnPaste 注册粘贴回调
func (h *PasteHook) OnPaste(callback func(string)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.callbacks = append(h.callbacks, callback)
}

// EmitPaste 触发粘贴
func (h *PasteHook) EmitPaste(text string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, cb := range h.callbacks {
		cb(text)
	}
}

// ScreenReaderHook 屏幕阅读器 Hook
type ScreenReaderHook struct {
	enabled bool
	context *hooks.HookContext
}

// UseIsScreenReaderEnabled 屏幕阅读器 Hook
func UseIsScreenReaderEnabled(ctx *hooks.HookContext) *ScreenReaderHook {
	return &ScreenReaderHook{
		enabled: false,
		context: ctx,
	}
}

// IsEnabled 返回是否启用
func (h *ScreenReaderHook) IsEnabled() bool {
	return h.enabled
}

// SetEnabled 设置是否启用
func (h *ScreenReaderHook) SetEnabled(enabled bool) {
	h.enabled = enabled
}

