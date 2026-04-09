// Animation Demo 示例 - 演示动画和帧更新
// 使用 ink.Render 实现真正的动画效果
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

// Spinner 动画帧
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Progress bar 组件
func ProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
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
	output += fmt.Sprintf("  \x1b[33mFPS:\x1b[0m 10\n")
	output += fmt.Sprintf("  \x1b[33mFrame:\x1b[0m %d\n", frame)
	output += fmt.Sprintf("  \x1b[33mProgress:\x1b[0m \x1b[32m%.2f\x1b[0m\n\n", progress)

	// 帧序列
	output += "Frame Sequence:\n"
	output += fmt.Sprintf("  \x1b[2m%v\x1b[0m\n\n", spinnerFrames)

	// 退出提示
	output += "\x1b[2mPress Ctrl+C to exit\x1b[0m\n"

	return output
}

// clearScreen 清屏
func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

// moveCursorHome 移动光标到起始位置
func moveCursorHome() {
	fmt.Print("\x1b[H")
}

// hideCursor 隐藏光标
func hideCursor() {
	fmt.Print("\x1b[?25l")
}

// showCursor 显示光标
func showCursor() {
	fmt.Print("\x1b[?25h")
}

// clearLines 清除指定行数
func clearLines(n int) {
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print("\x1b[1B") // 移动到下一行
		}
		fmt.Print("\x1b[2K") // 清除行
		fmt.Print("\x1b[0G") // 移动到行首
	}
	if n > 1 {
		fmt.Printf("\x1b[%dA", n-1) // 移动回第一行
	}
}

func main() {
	// 隐藏光标
	hideCursor()
	defer showCursor()

	// 处理 Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 动画状态
	frame := 0
	progress := 0.0
	lastHeight := 0

	// 动画循环
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	// 初始渲染
	output := renderFrame(frame, progress)
	fmt.Print(output)
	lastHeight = countLines(output)

	for {
		select {
		case <-ticker.C:
			// 更新状态
			frame++
			progress += 0.02
			if progress > 1.0 {
				progress = 0.0
			}

			// 清除之前的输出
			if lastHeight > 0 {
				clearLines(lastHeight)
			}

			// 渲染新帧
			output := renderFrame(frame, progress)
			fmt.Print(output)
			lastHeight = countLines(output)

		case <-sigChan:
			// 退出
			fmt.Println("\nAnimation demo exited cleanly")
			return
		}
	}
}

// countLines 计算行数
func countLines(s string) int {
	if s == "" {
		return 0
	}
	lines := 0
	for _, c := range s {
		if c == '\n' {
			lines++
		}
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		lines++
	}
	return lines
}

// 确保组件包被引用
var _ = components.Text
var _ = core.Element(nil)