// Package renderer 实现终端渲染器
package renderer

import (
	"strconv"
	"strings"
)

// Cell 单元格
type Cell struct {
	Char  rune
	Style Style
}

// Style 样式
type Style struct {
	FgColor Color
	BgColor Color
	Bold      bool
	Italic    bool
	Underline bool
	Dim       bool
	Blink     bool
	Reverse   bool
}

// Color 颜色
type Color struct {
	IsRGB  bool
	R, G, B uint8
	Index   uint8 // 256色索引
}

// PredefinedColors 预定义颜色
var (
	ColorDefault = Color{}
	ColorBlack   = Color{Index: 0}
	ColorRed     = Color{Index: 1}
	ColorGreen   = Color{Index: 2}
	ColorYellow  = Color{Index: 3}
	ColorBlue    = Color{Index: 4}
	ColorMagenta = Color{Index: 5}
	ColorCyan    = Color{Index: 6}
	ColorWhite   = Color{Index: 7}
)

// RGBColor 创建 RGB 颜色
func RGBColor(r, g, b uint8) Color {
	return Color{IsRGB: true, R: r, G: g, B: b}
}

// IndexColor 创建索引颜色
func IndexColor(index uint8) Color {
	return Color{Index: index}
}

// Buffer 输出缓冲
type Buffer struct {
	Width  int
	Height int
	Cells  [][]Cell
}

// NewBuffer 创建新缓冲
func NewBuffer(width, height int) *Buffer {
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
		for j := range cells[i] {
			cells[i][j] = Cell{Char: ' ', Style: Style{}}
		}
	}

	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// SetCell 设置单元格
func (b *Buffer) SetCell(x, y int, char rune, style Style) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	b.Cells[y][x] = Cell{Char: char, Style: style}
}

// GetCell 获取单元格
func (b *Buffer) GetCell(x, y int) Cell {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return Cell{}
	}
	return b.Cells[y][x]
}

// Clear 清空缓冲
func (b *Buffer) Clear() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.Cells[y][x] = Cell{Char: ' ', Style: Style{}}
		}
	}
}

// Diff 计算与另一个缓冲的差异
func (b *Buffer) Diff(other *Buffer) []Change {
	var changes []Change

	minWidth := b.Width
	if other.Width < minWidth {
		minWidth = other.Width
	}

	minHeight := b.Height
	if other.Height < minHeight {
		minHeight = other.Height
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			old := b.GetCell(x, y)
			new := other.GetCell(x, y)

			if old.Char != new.Char || !styleEqual(old.Style, new.Style) {
				changes = append(changes, Change{
					X:     x,
					Y:     y,
					Cell:  new,
				})
			}
		}
	}

	return changes
}

// Change 单元格变化
type Change struct {
	X, Y int
	Cell Cell
}

// styleEqual 比较样式是否相等
func styleEqual(a, b Style) bool {
	return a.FgColor == b.FgColor &&
		a.BgColor == b.BgColor &&
		a.Bold == b.Bold &&
		a.Italic == b.Italic &&
		a.Underline == b.Underline &&
		a.Dim == b.Dim &&
		a.Blink == b.Blink &&
		a.Reverse == b.Reverse
}

// Renderer 渲染器
type Renderer struct {
	width    int
	height   int
	buffer   *Buffer
	previous *Buffer
}

// NewRenderer 创建渲染器
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:    width,
		height:   height,
		buffer:   NewBuffer(width, height),
		previous: NewBuffer(width, height),
	}
}

// Resize 调整大小
func (r *Renderer) Resize(width, height int) {
	r.width = width
	r.height = height
	r.buffer = NewBuffer(width, height)
	r.previous = NewBuffer(width, height)
}

// GetBuffer 获取当前缓冲
func (r *Renderer) GetBuffer() *Buffer {
	return r.buffer
}

// Render 渲染并返回输出
func (r *Renderer) Render() string {
	// 计算差异
	changes := r.previous.Diff(r.buffer)

	if len(changes) == 0 {
		return ""
	}

	// 生成 ANSI 输出
	var output strings.Builder

	// 按行分组变化
	rowChanges := make(map[int][]Change)
	for _, change := range changes {
		rowChanges[change.Y] = append(rowChanges[change.Y], change)
	}

	// 生成输出
	for y := 0; y < r.height; y++ {
		if rowChanges[y] == nil {
			continue
		}

		// 移动光标到行首
		output.WriteString("\x1b[0G")
		// 移动光标到正确行
		output.WriteString("\x1b[")
		output.WriteString(intToStr(y + 1))
		output.WriteString(";1H")

		// 写入变化
		for _, change := range rowChanges[y] {
			// 移动到正确列
			if change.X > 0 {
				output.WriteString("\x1b[")
				output.WriteString(intToStr(change.X + 1))
				output.WriteString("G")
			}

			// 应用样式
			output.WriteString(styleToANSI(change.Cell.Style))

			// 写入字符
			output.WriteRune(change.Cell.Char)
		}

		// 重置样式
		output.WriteString("\x1b[0m")
	}

	// 更新 previous
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.previous.Cells[y][x] = r.buffer.Cells[y][x]
		}
	}

	return output.String()
}

// RenderFull 全量渲染
func (r *Renderer) RenderFull() string {
	var output strings.Builder

	// 清屏
	output.WriteString("\x1b[2J")
	output.WriteString("\x1b[H")

	// 渲染所有内容
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			cell := r.buffer.Cells[y][x]

			// 应用样式
			output.WriteString(styleToANSI(cell.Style))

			// 写入字符
			output.WriteRune(cell.Char)
		}

		// 重置样式并换行
		output.WriteString("\x1b[0m\n")
	}

	// 更新 previous
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.previous.Cells[y][x] = r.buffer.Cells[y][x]
		}
	}

	return output.String()
}

// styleToANSI 将样式转换为 ANSI 转义序列
func styleToANSI(style Style) string {
	var codes []string

	// 前景色
	if style.FgColor.IsRGB {
		codes = append(codes,
			"\x1b[38;2;"+intToStr(int(style.FgColor.R))+";"+
			intToStr(int(style.FgColor.G))+";"+
			intToStr(int(style.FgColor.B))+"m")
	} else if style.FgColor.Index > 0 {
		if style.FgColor.Index < 8 {
			codes = append(codes, "\x1b[" + intToStr(30 + int(style.FgColor.Index)) + "m")
		} else {
			codes = append(codes, "\x1b[38;5;" + intToStr(int(style.FgColor.Index)) + "m")
		}
	}

	// 背景色
	if style.BgColor.IsRGB {
		codes = append(codes,
			"\x1b[48;2;"+intToStr(int(style.BgColor.R))+";"+
			intToStr(int(style.BgColor.G))+";"+
			intToStr(int(style.BgColor.B))+"m")
	} else if style.BgColor.Index > 0 {
		if style.BgColor.Index < 8 {
			codes = append(codes, "\x1b[" + intToStr(40 + int(style.BgColor.Index)) + "m")
		} else {
			codes = append(codes, "\x1b[48;5;" + intToStr(int(style.BgColor.Index)) + "m")
		}
	}

	// 样式属性
	if style.Bold {
		codes = append(codes, "\x1b[1m")
	}
	if style.Italic {
		codes = append(codes, "\x1b[3m")
	}
	if style.Underline {
		codes = append(codes, "\x1b[4m")
	}
	if style.Dim {
		codes = append(codes, "\x1b[2m")
	}
	if style.Blink {
		codes = append(codes, "\x1b[5m")
	}
	if style.Reverse {
		codes = append(codes, "\x1b[7m")
	}

	return strings.Join(codes, "")
}

// intToStr 整数转字符串
func intToStr(n int) string {
	return strconv.Itoa(n)
}