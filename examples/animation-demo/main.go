// Animation Demo 示例 - 演示动画和帧更新
// 使用 TUI 渲染实现真正的动画效果
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Spinner 动画帧
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Progress bar 组件
func ProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	bar := ""
	for range filled {
		bar += "█"
	}
	for range empty {
		bar += "░"
	}
	return bar
}

// renderFrame 渲染一帧
func renderFrame(frame int, progress float64) string {
	spinnerFrame := spinnerFrames[frame%len(spinnerFrames)]

	output := ""

	// 边框
	output += "╭──────────────────────────────────────╮\n"
	output += "│  🎬 Animation Demo (Real Animation)  │\n"
	output += "╰──────────────────────────────────────╯\n"
	output += "\n"

	// Spinner 动画 - 实时更新
	output += "Spinner:\n"
	output += fmt.Sprintf("  \x1b[36m%s\x1b[0m Loading...\n\n", spinnerFrame)

	// 进度条 - 实时更新
	output += "Progress Bar:\n"
	output += fmt.Sprintf("  \x1b[32m%s\x1b[0m %.0f%%\n\n", ProgressBar(progress, 20), progress*100)

	// 动画状态
	output += "Animation State:\n"
	output += "  \x1b[33mFPS:\x1b[0m 10\n"
	output += fmt.Sprintf("  \x1b[33mFrame:\x1b[0m %d\n", frame)
	output += fmt.Sprintf("  \x1b[33mProgress:\x1b[0m \x1b[32m%.2f\x1b[0m\n\n", progress)

	// 帧序列
	output += "Frame Sequence:\n"
	output += fmt.Sprintf("  \x1b[2m%v\x1b[0m\n\n", spinnerFrames)

	// 退出提示
	output += "\x1b[2mPress Ctrl+C to exit\x1b[0m"

	return output
}

// Terminal ANSI 控制序列
const (
	hideCursor              = "\x1b[?25l"
	showCursor              = "\x1b[?25h"
	alternateScreenOn       = "\x1b[?1049h"
	alternateScreenOff      = "\x1b[?1049l"
	clearScreen             = "\x1b[2J"
	moveCursorHome          = "\x1b[H"
	clearFromCursorToEnd    = "\x1b[0J"
	saveCursorPosition      = "\x1b[s"
	restoreCursorPosition   = "\x1b[u"
)

func main() {
	// 启用备用屏幕缓冲区 (退出时恢复原屏幕)
	fmt.Print(alternateScreenOn)
	defer fmt.Print(alternateScreenOff)

	// 隐藏光标
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	// 清屏
	fmt.Print(clearScreen + moveCursorHome)

	// 保存光标位置
	fmt.Print(saveCursorPosition)

	// 处理 Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 动画状态
	frame := 0
	progress := 0.0

	// 动画循环 (10 FPS)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// 初始渲染
	output := renderFrame(frame, progress)
	fmt.Print(output)

	for {
		select {
		case <-ticker.C:
			// 更新状态
			frame++
			progress += 0.02
			if progress > 1.0 {
				progress = 0.0
			}

			// 恢复光标位置并清除屏幕
			fmt.Print(restoreCursorPosition + clearFromCursorToEnd)

			// 渲染新帧
			output := renderFrame(frame, progress)
			fmt.Print(output)

		case <-sigChan:
			// 退出 (defer 会处理恢复终端状态)
			return
		}
	}
}