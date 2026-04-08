package renderer

import (
	"strings"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer(10, 5)

	if buf.Width != 10 {
		t.Errorf("Expected width 10, got %d", buf.Width)
	}

	if buf.Height != 5 {
		t.Errorf("Expected height 5, got %d", buf.Height)
	}

	if len(buf.Cells) != 5 {
		t.Errorf("Expected 5 rows, got %d", len(buf.Cells))
	}

	if len(buf.Cells[0]) != 10 {
		t.Errorf("Expected 10 columns, got %d", len(buf.Cells[0]))
	}
}

func TestSetCell(t *testing.T) {
	buf := NewBuffer(10, 5)

	style := Style{FgColor: ColorGreen}
	buf.SetCell(2, 1, 'A', style)

	cell := buf.GetCell(2, 1)
	if cell.Char != 'A' {
		t.Errorf("Expected 'A', got %c", cell.Char)
	}

	if cell.Style.FgColor != ColorGreen {
		t.Error("Expected green color")
	}
}

func TestSetCellBounds(t *testing.T) {
	buf := NewBuffer(10, 5)

	// 越界设置应该被忽略
	buf.SetCell(-1, 0, 'A', Style{})
	buf.SetCell(10, 0, 'A', Style{})
	buf.SetCell(0, -1, 'A', Style{})
	buf.SetCell(0, 5, 'A', Style{})

	// 不应该崩溃
}

func TestClear(t *testing.T) {
	buf := NewBuffer(10, 5)
	buf.SetCell(0, 0, 'A', Style{FgColor: ColorGreen})

	buf.Clear()

	cell := buf.GetCell(0, 0)
	if cell.Char != ' ' {
		t.Error("Cell should be cleared to space")
	}

	if cell.Style.FgColor != ColorDefault {
		t.Error("Style should be reset")
	}
}

func TestDiff(t *testing.T) {
	old := NewBuffer(10, 5)
	new := NewBuffer(10, 5)

	// 设置不同的单元格
	old.SetCell(0, 0, 'A', Style{})
	new.SetCell(0, 0, 'B', Style{})

	old.SetCell(1, 1, 'X', Style{FgColor: ColorRed})
	new.SetCell(1, 1, 'X', Style{FgColor: ColorGreen})

	changes := old.Diff(new)

	// 应该有 2 个变化
	if len(changes) != 2 {
		t.Errorf("Expected 2 changes, got %d", len(changes))
	}
}

func TestDiffNoChange(t *testing.T) {
	buf1 := NewBuffer(10, 5)
	buf2 := NewBuffer(10, 5)

	buf1.SetCell(0, 0, 'A', Style{FgColor: ColorGreen})
	buf2.SetCell(0, 0, 'A', Style{FgColor: ColorGreen})

	changes := buf1.Diff(buf2)

	if len(changes) != 0 {
		t.Errorf("Expected 0 changes, got %d", len(changes))
	}
}

func TestNewRenderer(t *testing.T) {
	r := NewRenderer(80, 24)

	if r.width != 80 {
		t.Errorf("Expected width 80, got %d", r.width)
	}

	if r.height != 24 {
		t.Errorf("Expected height 24, got %d", r.height)
	}

	if r.buffer == nil {
		t.Error("Buffer should not be nil")
	}

	if r.previous == nil {
		t.Error("Previous buffer should not be nil")
	}
}

func TestResize(t *testing.T) {
	r := NewRenderer(80, 24)
	r.Resize(100, 30)

	if r.width != 100 {
		t.Errorf("Expected width 100, got %d", r.width)
	}

	if r.height != 30 {
		t.Errorf("Expected height 30, got %d", r.height)
	}
}

func TestRenderNoChange(t *testing.T) {
	r := NewRenderer(10, 5)

	// 初始渲染
	output := r.Render()

	// 没有变化，应该返回空
	if output != "" {
		t.Error("Expected empty output for no changes")
	}
}

func TestRenderWithChange(t *testing.T) {
	r := NewRenderer(10, 5)

	// 设置单元格
	r.GetBuffer().SetCell(0, 0, 'A', Style{})

	// 渲染
	output := r.Render()

	// 应该有输出
	if output == "" {
		t.Error("Expected output for changes")
	}

	// 应该包含 ANSI 转义序列
	if !strings.Contains(output, "\x1b[") {
		t.Error("Output should contain ANSI escape sequences")
	}
}

func TestRenderFull(t *testing.T) {
	r := NewRenderer(10, 5)

	// 设置一些内容
	r.GetBuffer().SetCell(0, 0, 'A', Style{})
	r.GetBuffer().SetCell(1, 0, 'B', Style{})

	// 全量渲染
	output := r.RenderFull()

	// 应该包含清屏
	if !strings.Contains(output, "\x1b[2J") {
		t.Error("Output should contain clear screen")
	}

	// 应该包含字符
	if !strings.Contains(output, "A") || !strings.Contains(output, "B") {
		t.Error("Output should contain characters")
	}
}

func TestRGBColor(t *testing.T) {
	color := RGBColor(255, 128, 64)

	if !color.IsRGB {
		t.Error("Should be RGB color")
	}

	if color.R != 255 || color.G != 128 || color.B != 64 {
		t.Error("RGB values should match")
	}
}

func TestIndexColor(t *testing.T) {
	color := IndexColor(42)

	if color.IsRGB {
		t.Error("Should not be RGB color")
	}

	if color.Index != 42 {
		t.Errorf("Expected index 42, got %d", color.Index)
	}
}

func TestStyleToANSI(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		contains string
	}{
		{
			name:     "bold",
			style:    Style{Bold: true},
			contains: "\x1b[1m",
		},
		{
			name:     "italic",
			style:    Style{Italic: true},
			contains: "\x1b[3m",
		},
		{
			name:     "underline",
			style:    Style{Underline: true},
			contains: "\x1b[4m",
		},
		{
			name:     "dim",
			style:    Style{Dim: true},
			contains: "\x1b[2m",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := styleToANSI(test.style)
			if !strings.Contains(output, test.contains) {
				t.Errorf("Output should contain %q, got %q", test.contains, output)
			}
		})
	}
}

func TestStyleEqual(t *testing.T) {
	a := Style{FgColor: ColorGreen, Bold: true}
	b := Style{FgColor: ColorGreen, Bold: true}
	c := Style{FgColor: ColorRed, Bold: true}

	if !styleEqual(a, b) {
		t.Error("a and b should be equal")
	}

	if styleEqual(a, c) {
		t.Error("a and c should not be equal")
	}
}
