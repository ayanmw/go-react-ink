// Package tcell 提供终端抽象层接口
// 实际使用时需要安装 github.com/gdamore/tcell/v2
package tcell

import (
	"time"
)

// Screen 终端屏幕接口
type Screen interface {
	Init() error
	Fini()
	Size() (int, int)
	Clear()
	SetContent(x, y int, ch rune, comb []rune, style Style)
	Show()
	Sync()
	EnableMouse()
	DisableMouse()
	PollEvent() Event
	PostEvent(ev Event)
}

// Style 样式
type Style struct {
	fg       Color
	bg       Color
	bold     bool
	italic   bool
	underline bool
	dim      bool
	blink    bool
	reverse  bool
}

// NewStyle 创建新样式
func NewStyle() Style {
	return Style{}
}

// Foreground 设置前景色
func (s Style) Foreground(c Color) Style {
	s.fg = c
	return s
}

// Background 设置背景色
func (s Style) Background(c Color) Style {
	s.bg = c
	return s
}

// Bold 设置粗体
func (s Style) Bold(on bool) Style {
	s.bold = on
	return s
}

// Italic 设置斜体
func (s Style) Italic(on bool) Style {
	s.italic = on
	return s
}

// Underline 设置下划线
func (s Style) Underline(on bool) Style {
	s.underline = on
	return s
}

// Dim 设置变暗
func (s Style) Dim(on bool) Style {
	s.dim = on
	return s
}

// Blink 设置闪烁
func (s Style) Blink(on bool) Style {
	s.blink = on
	return s
}

// Reverse 设置反转
func (s Style) Reverse(on bool) Style {
	s.reverse = on
	return s
}

// Color 颜色
type Color int32

const (
	ColorDefault Color = -1
	ColorBlack   Color = 0
	ColorRed     Color = 1
	ColorGreen   Color = 2
	ColorYellow  Color = 3
	ColorBlue    Color = 4
	ColorMagenta Color = 5
	ColorCyan    Color = 6
	ColorWhite   Color = 7
)

// NewRGBColor 创建 RGB 颜色
func NewRGBColor(r, g, b int32) Color {
	return Color((1 << 24) | (r << 16) | (g << 8) | b)
}

// Event 事件接口
type Event interface {
	When() time.Time
}

// EventKey 键盘事件
type EventKey struct {
	when time.Time
	key  Key
	ch   rune
	mod  ModMask
}

// When 返回事件时间
func (e *EventKey) When() time.Time {
	return e.when
}

// Key 返回键码
func (e *EventKey) Key() Key {
	return e.key
}

// Rune 返回字符
func (e *EventKey) Rune() rune {
	return e.ch
}

// Modifiers 返回修饰键
func (e *EventKey) Modifiers() ModMask {
	return e.mod
}

// Name 返回键名
func (e *EventKey) Name() string {
	if e.key == KeyRune {
		return string(e.ch)
	}
	return keyNames[e.key]
}

// EventResize 调整大小事件
type EventResize struct {
	when   time.Time
	width  int
	height int
}

// When 返回事件时间
func (e *EventResize) When() time.Time {
	return e.when
}

// Size 返回尺寸
func (e *EventResize) Size() (int, int) {
	return e.width, e.height
}

// EventMouse 鼠标事件
type EventMouse struct {
	when   time.Time
	x, y   int
	btn    ButtonMask
	mod    ModMask
}

// When 返回事件时间
func (e *EventMouse) When() time.Time {
	return e.when
}

// Position 返回位置
func (e *EventMouse) Position() (int, int) {
	return e.x, e.y
}

// Buttons 返回按钮
func (e *EventMouse) Buttons() ButtonMask {
	return e.btn
}

// Modifiers 返回修饰键
func (e *EventMouse) Modifiers() ModMask {
	return e.mod
}

// Key 键码
type Key int

const (
	KeyNone Key = iota
	KeyRune
	KeyEnter
	KeyBackspace
	KeyTab
	KeyEscape
	KeyBacktab
	KeyInsert
	KeyDelete
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyUpLeft
	KeyUpRight
	KeyDownLeft
	KeyDownRight
	KeyCenter
	KeyPgUp
	KeyPgDn
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyF26
	KeyF27
	KeyF28
	KeyF29
	KeyF30
	KeyF31
	KeyF32
	KeyF33
	KeyF34
	KeyF35
	KeyF36
	KeyF37
	KeyF38
	KeyF39
	KeyF40
	KeyF41
	KeyF42
	KeyF43
	KeyF44
	KeyF45
	KeyF46
	KeyF47
	KeyF48
	KeyF49
	KeyF50
	KeyF51
	KeyF52
	KeyF53
	KeyF54
	KeyF55
	KeyF56
	KeyF57
	KeyF58
	KeyF59
	KeyF60
	KeyF61
	KeyF62
	KeyF63
	KeyF64
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
	KeyCtrlSpace
	KeyCtrlUnderscore
	KeyCtrlRightSq
	KeyCtrlBackslash
	KeyCtrlCarat
)

var keyNames = map[Key]string{
	KeyEnter:      "Enter",
	KeyBackspace:  "Backspace",
	KeyTab:        "Tab",
	KeyEscape:     "Escape",
	KeyInsert:     "Insert",
	KeyDelete:     "Delete",
	KeyUp:         "Up",
	KeyDown:       "Down",
	KeyLeft:       "Left",
	KeyRight:      "Right",
	KeyHome:       "Home",
	KeyEnd:        "End",
	KeyPgUp:       "PgUp",
	KeyPgDn:       "PgDn",
	KeyF1:         "F1",
	KeyF2:         "F2",
	KeyF3:         "F3",
	KeyF4:         "F4",
	KeyF5:         "F5",
	KeyF6:         "F6",
	KeyF7:         "F7",
	KeyF8:         "F8",
	KeyF9:         "F9",
	KeyF10:        "F10",
	KeyF11:        "F11",
	KeyF12:        "F12",
	KeyCtrlA:      "Ctrl+A",
	KeyCtrlB:      "Ctrl+B",
	KeyCtrlC:      "Ctrl+C",
	KeyCtrlD:      "Ctrl+D",
	KeyCtrlE:      "Ctrl+E",
	KeyCtrlF:      "Ctrl+F",
	KeyCtrlG:      "Ctrl+G",
	KeyCtrlH:      "Ctrl+H",
	KeyCtrlI:      "Ctrl+I",
	KeyCtrlJ:      "Ctrl+J",
	KeyCtrlK:      "Ctrl+K",
	KeyCtrlL:      "Ctrl+L",
	KeyCtrlM:      "Ctrl+M",
	KeyCtrlN:      "Ctrl+N",
	KeyCtrlO:      "Ctrl+O",
	KeyCtrlP:      "Ctrl+P",
	KeyCtrlQ:      "Ctrl+Q",
	KeyCtrlR:      "Ctrl+R",
	KeyCtrlS:      "Ctrl+S",
	KeyCtrlT:      "Ctrl+T",
	KeyCtrlU:      "Ctrl+U",
	KeyCtrlV:      "Ctrl+V",
	KeyCtrlW:      "Ctrl+W",
	KeyCtrlX:      "Ctrl+X",
	KeyCtrlY:      "Ctrl+Y",
	KeyCtrlZ:      "Ctrl+Z",
}

// ModMask 修饰键掩码
type ModMask int

const (
	ModNone ModMask = iota
	ModAlt
	ModCtrl
	ModShift
	ModMeta
)

// ButtonMask 鼠标按钮掩码
type ButtonMask int

const (
	ButtonNone ButtonMask = iota
	ButtonLeft
	ButtonMiddle
	ButtonRight
	ButtonWheelUp
	ButtonWheelDown
)