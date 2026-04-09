// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"io"
	"strings"
)

// LogUpdate 终端输出管理器 (类似 log-update)
// 负责管理终端输出，支持清除之前的输出并写入新内容
type LogUpdate struct {
	stdout     io.Writer
	lastOutput string
	lastHeight int
}

// NewLogUpdate 创建新的输出管理器
func NewLogUpdate(stdout io.Writer) *LogUpdate {
	return &LogUpdate{
		stdout: stdout,
	}
}

// Write 写入新输出，自动清除之前的输出
func (l *LogUpdate) Write(output string) error {
	// 计算新输出的行数
	newHeight := countLines(output)

	// 如果之前有输出，先清除
	if l.lastHeight > 0 {
		l.clearPrevious()
	}

	// 写入新输出
	_, err := l.stdout.Write([]byte(output))
	if err != nil {
		return err
	}

	l.lastOutput = output
	l.lastHeight = newHeight
	return nil
}

// Clear 清除所有输出
func (l *LogUpdate) Clear() error {
	if l.lastHeight == 0 {
		return nil
	}

	l.clearPrevious()
	l.lastOutput = ""
	l.lastHeight = 0
	return nil
}

// Done 完成输出，保留最后一帧
func (l *LogUpdate) Done() {
	// 显示光标
	l.stdout.Write([]byte("\x1b[?25h"))
	l.lastOutput = ""
	l.lastHeight = 0
}

// clearPrevious 清除之前的输出
func (l *LogUpdate) clearPrevious() {
	if l.lastHeight <= 0 {
		return
	}

	// 移动光标到之前输出的开头
	// 先移动到行首
	l.stdout.Write([]byte("\x1b[0G"))
	// 向上移动 lastHeight-1 行
	if l.lastHeight > 1 {
		l.stdout.Write([]byte("\x1b[" + intToStr(l.lastHeight-1) + "A"))
	}

	// 清除每一行
	for i := 0; i < l.lastHeight; i++ {
		if i > 0 {
			// 移动到下一行
			l.stdout.Write([]byte("\x1b[1B"))
		}
		// 清除当前行
		l.stdout.Write([]byte("\x1b[2K"))
		// 移动到行首
		l.stdout.Write([]byte("\x1b[0G"))
	}

	// 移动回第一行
	if l.lastHeight > 1 {
		l.stdout.Write([]byte("\x1b[" + intToStr(l.lastHeight-1) + "A"))
	}
}

// countLines 计算字符串的行数
func countLines(s string) int {
	if s == "" {
		return 0
	}
	lines := strings.Count(s, "\n")
	// 如果字符串不以换行结尾，需要加1
	if !strings.HasSuffix(s, "\n") {
		lines++
	}
	return lines
}

// intToStr 整数转字符串
func intToStr(n int) string {
	if n <= 0 {
		return "0"
	}
	// 简单的整数转字符串
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

// ANSI 转义序列常量
const (
	// 光标控制
	ANSIHideCursor      = "\x1b[?25l"
	ANSIShowCursor      = "\x1b[?25h"
	ANSIClearScreen     = "\x1b[2J"
	ANSIMoveCursorHome  = "\x1b[H"
	ANSIClearLine       = "\x1b[2K"
	ANSIMoveCursorStart = "\x1b[0G"
	ANSISaveCursor      = "\x1b[s"
	ANSIRestoreCursor   = "\x1b[u"

	// 屏幕缓冲
	ANSIEnableAlternateScreen  = "\x1b[?1049h"
	ANSIDisableAlternateScreen = "\x1b[?1049l"

	// 样式
	ANSIReset = "\x1b[0m"
	ANSIBold  = "\x1b[1m"
	ANSIDim   = "\x1b[2m"
)

// MoveCursorUp 移动光标向上 N 行
func MoveCursorUp(n int) string {
	if n <= 0 {
		return ""
	}
	return "\x1b[" + intToStr(n) + "A"
}

// MoveCursorDown 移动光标向下 N 行
func MoveCursorDown(n int) string {
	if n <= 0 {
		return ""
	}
	return "\x1b[" + intToStr(n) + "B"
}

// MoveCursorForward 移动光标向前 N 列
func MoveCursorForward(n int) string {
	if n <= 0 {
		return ""
	}
	return "\x1b[" + intToStr(n) + "C"
}

// MoveCursorBackward 移动光标向后 N 列
func MoveCursorBackward(n int) string {
	if n <= 0 {
		return ""
	}
	return "\x1b[" + intToStr(n) + "D"
}

// MoveCursorTo 移动光标到指定位置 (1-indexed)
func MoveCursorTo(row, col int) string {
	return "\x1b[" + intToStr(row) + ";" + intToStr(col) + "H"
}

// ClearLines 清除 N 行
func ClearLines(n int) string {
	var result string
	for i := 0; i < n; i++ {
		if i > 0 {
			result += "\x1b[1B" // 移动到下一行
		}
		result += "\x1b[2K" // 清除行
		result += "\x1b[0G" // 移动到行首
	}
	return result
}
