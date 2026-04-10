// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"io"
	"strings"
)

// LogUpdate 终端输出管理器 (类似 log-update)
// 负责管理终端输出，支持清除之前的输出并写入新内容
type LogUpdate struct {
	stdout          io.Writer
	previousOutput  string
	previousLines   []string
	initialized     bool
	incremental     bool // 增量渲染模式
	alternateScreen bool // 备用屏幕缓冲模式
}

// NewLogUpdate 创建新的输出管理器
func NewLogUpdate(stdout io.Writer, incremental bool, alternateScreen bool) *LogUpdate {
	return &LogUpdate{
		stdout:          stdout,
		previousLines:   make([]string, 0),
		incremental:     incremental,
		alternateScreen: alternateScreen,
	}
}

// SetIncremental 设置增量渲染模式
func (l *LogUpdate) SetIncremental(enabled bool) {
	l.incremental = enabled
}

// Initialize 初始化终端
func (l *LogUpdate) Initialize() error {
	if l.initialized {
		return nil
	}

	// 启用备用屏幕缓冲
	if l.alternateScreen {
		l.stdout.Write([]byte(ANSIEnableAlternateScreen))
	}

	// 隐藏光标
	l.stdout.Write([]byte(ANSIHideCursor))
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

	// 禁用备用屏幕缓冲
	if l.alternateScreen {
		l.stdout.Write([]byte(ANSIDisableAlternateScreen))
	}

	l.previousOutput = ""
	l.previousLines = make([]string, 0)
	l.initialized = false
}

// Write 写入新输出
func (l *LogUpdate) Write(output string) error {
	// 检查内容是否变化
	if output == l.previousOutput {
		return nil // 内容未变化，跳过渲染
	}

	if l.incremental {
		return l.writeIncremental(output)
	}
	return l.writeStandard(output)
}

// writeStandard 标准渲染模式
func (l *LogUpdate) writeStandard(output string) error {
	// 计算之前的行数
	previousLineCount := len(l.previousLines)

	// 移动光标到起始位置
	l.stdout.Write([]byte(ANSIMoveCursorHome))

	// 清除之前的行
	if previousLineCount > 0 {
		l.stdout.Write([]byte(EraseLines(previousLineCount)))
	}

	// 写入新内容
	l.stdout.Write([]byte(output))

	// 更新状态
	l.previousOutput = output
	l.previousLines = strings.Split(output, "\n")

	return nil
}

// writeIncremental 增量渲染模式 - 只更新变化的行
func (l *LogUpdate) writeIncremental(output string) error {
	nextLines := strings.Split(output, "\n")
	previousLineCount := len(l.previousLines)

	// 如果是第一次渲染，直接写入
	if previousLineCount == 0 {
		l.stdout.Write([]byte(output))
		l.previousOutput = output
		l.previousLines = nextLines
		return nil
	}

	// 移动光标到第一行
	l.stdout.Write([]byte(ANSIMoveCursorHome))

	// 计算最大行数
	maxLines := previousLineCount
	if len(nextLines) > maxLines {
		maxLines = len(nextLines)
	}

	// 逐行比较和更新
	for i := 0; i < maxLines; i++ {
		// 如果新内容行数少于旧行数，清除多余行
		if i >= len(nextLines) {
			l.stdout.Write([]byte(ANSIClearLine))
			if i < maxLines-1 {
				l.stdout.Write([]byte(MoveCursorDown(1)))
			}
			continue
		}

		// 如果旧行数不够，直接写入新行
		if i >= len(l.previousLines) {
			l.stdout.Write([]byte(nextLines[i]))
			if i < len(nextLines)-1 || strings.HasSuffix(output, "\n") {
				l.stdout.Write([]byte("\n"))
			}
			continue
		}

		// 比较行内容
		if nextLines[i] == l.previousLines[i] {
			// 内容相同，跳过但移动到下一行
			if i < maxLines-1 {
				l.stdout.Write([]byte(MoveCursorDown(1)))
			}
			continue
		}

		// 内容不同，清除并重写
		l.stdout.Write([]byte(ANSIClearLine))
		l.stdout.Write([]byte("\x1b[0G")) // 移动到行首
		l.stdout.Write([]byte(nextLines[i]))
		if i < len(nextLines)-1 || strings.HasSuffix(output, "\n") {
			l.stdout.Write([]byte("\n"))
		}
	}

	// 更新状态
	l.previousOutput = output
	l.previousLines = nextLines

	return nil
}

// Clear 清除所有输出
func (l *LogUpdate) Clear() error {
	lineCount := len(l.previousLines)
	if lineCount > 0 {
		l.stdout.Write([]byte(ANSIMoveCursorHome))
		l.stdout.Write([]byte(EraseLines(lineCount)))
	}

	l.previousOutput = ""
	l.previousLines = make([]string, 0)
	return nil
}

// Reset 重置状态（不改变终端）
func (l *LogUpdate) Reset() {
	l.previousOutput = ""
	l.previousLines = make([]string, 0)
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

// EraseLines 清除 N 行 (从当前位置开始)
func EraseLines(n int) string {
	if n <= 0 {
		return ""
	}
	// 使用 \x1b[NM 清除 N 行
	return "\x1b[" + intToStr(n) + "M"
}