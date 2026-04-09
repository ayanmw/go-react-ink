// Animation Demo 示例 - 演示动画和帧更新
// 使用增量渲染避免闪烁
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Spinner 动画帧
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ProgressBar 生成进度条
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
func renderFrame(frame int, progress float64) []string {
	spinnerFrame := spinnerFrames[frame%len(spinnerFrames)]

	lines := make([]string, 0, 15)

	// 边框
	lines = append(lines, "╭──────────────────────────────────────╮")
	lines = append(lines, "│  🎬 Animation Demo (Real Animation)  │")
	lines = append(lines, "╰──────────────────────────────────────╯")
	lines = append(lines, "")

	// Spinner 动画
	lines = append(lines, "Spinner:")
	lines = append(lines, fmt.Sprintf("  \x1b[36m%s\x1b[0m Loading...", spinnerFrame))
	lines = append(lines, "")

	// 进度条
	lines = append(lines, "Progress Bar:")
	lines = append(lines, fmt.Sprintf("  \x1b[32m%s\x1b[0m %.0f%%", ProgressBar(progress, 20), progress*100))
	lines = append(lines, "")

	// 动画状态
	lines = append(lines, "Animation State:")
	lines = append(lines, "  \x1b[33mFPS:\x1b[0m 10")
	lines = append(lines, fmt.Sprintf("  \x1b[33mFrame:\x1b[0m %d", frame))
	lines = append(lines, fmt.Sprintf("  \x1b[33mProgress:\x1b[0m \x1b[32m%.2f\x1b[0m", progress))
	lines = append(lines, "")

	// 帧序列
	lines = append(lines, "Frame Sequence:")
	lines = append(lines, fmt.Sprintf("  \x1b[2m%v\x1b[0m", spinnerFrames))
	lines = append(lines, "")

	// 退出提示
	lines = append(lines, "\x1b[2mPress Ctrl+C to exit\x1b[0m")

	return lines
}

// IncrementalRenderer 增量渲染器
type IncrementalRenderer struct {
	previousLines []string
}

// NewIncrementalRenderer 创建增量渲染器
func NewIncrementalRenderer() *IncrementalRenderer {
	return &IncrementalRenderer{
		previousLines: make([]string, 0),
	}
}

// Render 增量渲染 - 只更新变化的行
func (r *IncrementalRenderer) Render(lines []string) {
	// 第一次渲染
	if len(r.previousLines) == 0 {
		for _, line := range lines {
			fmt.Println(line)
		}
		r.previousLines = lines
		return
	}

	// 移动光标到第一行
	fmt.Print("\x1b[H")

	maxLines := len(r.previousLines)
	if len(lines) > maxLines {
		maxLines = len(lines)
	}

	// 逐行比较
	for i := 0; i < maxLines; i++ {
		// 新内容行数少于旧行数，清除多余行
		if i >= len(lines) {
			fmt.Print("\x1b[2K") // 清除当前行
			if i < maxLines-1 {
				fmt.Print("\x1b[1B") // 移动到下一行
			}
			continue
		}

		// 旧行数不够，直接写入新行
		if i >= len(r.previousLines) {
			fmt.Print("\x1b[2K") // 清除当前行
			fmt.Print("\x1b[0G") // 移动到行首
			fmt.Print(lines[i])
			if i < maxLines-1 {
				fmt.Print("\x1b[1B") // 移动到下一行
			}
			continue
		}

		// 比较行内容
		if lines[i] == r.previousLines[i] {
			// 内容相同，跳过
			if i < maxLines-1 {
				fmt.Print("\x1b[1B") // 移动到下一行
			}
			continue
		}

		// 内容不同，清除并重写
		fmt.Print("\x1b[2K") // 清除当前行
		fmt.Print("\x1b[0G") // 移动到行首
		fmt.Print(lines[i])
		if i < maxLines-1 {
			fmt.Print("\x1b[1B") // 移动到下一行
		}
	}

	r.previousLines = lines
}

// Clear 清除所有内容
func (r *IncrementalRenderer) Clear() {
	if len(r.previousLines) > 0 {
		fmt.Print("\x1b[H")
		fmt.Printf("\x1b[%dM", len(r.previousLines))
	}
	r.previousLines = make([]string, 0)
}

func main() {
	// 启用备用屏幕缓冲区
	fmt.Print("\x1b[?1049h")
	defer fmt.Print("\x1b[?1049l")

	// 隐藏光标
	fmt.Print("\x1b[?25l")
	defer fmt.Print("\x1b[?25h")

	// 清屏
	fmt.Print("\x1b[2J\x1b[H")

	// 处理 Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 动画状态
	frame := 0
	progress := 0.0

	// 创建增量渲染器
	renderer := NewIncrementalRenderer()

	// 动画循环 (10 FPS)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// 初始渲染
	lines := renderFrame(frame, progress)
	renderer.Render(lines)

	for {
		select {
		case <-ticker.C:
			// 更新状态
			frame++
			progress += 0.02
			if progress > 1.0 {
				progress = 0.0
			}

			// 渲染新帧（增量更新）
			lines := renderFrame(frame, progress)
			renderer.Render(lines)

		case <-sigChan:
			// 退出
			return
		}
	}
}

// 确保 strings 包被使用
var _ = strings.TrimSpace("")