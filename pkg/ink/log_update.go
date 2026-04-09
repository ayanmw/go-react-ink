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
	initialized bool
}

// NewLogUpdate 创建新的输出管理器
func NewLogUpdate(stdout io.Writer) *LogUpdate {
	return &LogUpdate{
		stdout: stdout,
	}
}

// Initialize 初始化终端 (启用备用屏幕缓冲)
func (l *LogUpdate) Initialize() error {
	if l.initialized {
		return nil
	}

	// 启用备用屏幕缓冲区
	l.stdout.Write([]byte(ANSIEnableAlternateScreen))
	// 隐藏光标
	l.stdout.Write([]byte(ANSIHideCursor))
	// 清屏
	l.stdout.Write([]byte(ANSIClearScreen))
	// 移动光标到起始位置
	l.stdout.Write([]byte(ANSIMoveCursorHome))

	l.initialized = true
	return nil
}

// Done 完成输出，恢复终端状态
func (l *LogUpdate) Done() {
	if !l.initialized {
		return
	}

	// 显示光标
	l.stdout.Write([]byte(ANSIShowCursor))
	// 禁用备用屏幕缓冲区 (恢复原屏幕内容)
	l.stdout.Write([]byte(ANSIDisableAlternateScreen))

	l.lastOutput = ""
	l.lastHeight = 0
	l.initialized = false
}

// Write 写入新输出
func (l *LogUpdate) Write(output string) error {
	// 移动光标到起始位置
	l.stdout.Write([]byte(ANSIMoveCursorHome))
	// 清除从光标到屏幕末尾的内容
	l.stdout.Write([]byte("\x1b[0J"))

	// 写入新输出
	_, err := l.stdout.Write([]byte(output))
	if err != nil {
		return err
	}

	l.lastOutput = output
	l.lastHeight = countLines(output)
	return nil
}

// Clear 清除所有输出
func (l *LogUpdate) Clear() error {
	// 移动光标到起始位置
	l.stdout.Write([]byte(ANSIMoveCursorHome))
	// 清除从光标到屏幕末尾的内容
	l.stdout.Write([]byte("\x1b[0J"))

	l.lastOutput = ""
	l.lastHeight = 0
	return nil
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
