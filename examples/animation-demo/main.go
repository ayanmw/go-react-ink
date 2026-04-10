// Animation Demo 示例 - 演示动画和帧更新
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/ink"
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

// AnimationComponent 动画组件
type AnimationComponent struct {
	frame    int
	progress float64
}

// Render 实现渲染方法
func (a *AnimationComponent) Render() string {
	spinnerFrame := spinnerFrames[a.frame%len(spinnerFrames)]

	var builder strings.Builder

	// 边框
	builder.WriteString("╭──────────────────────────────────────╮\n")
	builder.WriteString("│  🎬 Animation Demo (Real Animation)  │\n")
	builder.WriteString("╰──────────────────────────────────────╯\n")
	builder.WriteString("\n")

	// Spinner 动画
	builder.WriteString("Spinner:\n")
	builder.WriteString(fmt.Sprintf("  \x1b[36m%s\x1b[0m Loading...", spinnerFrame))
	builder.WriteString("\n")

	// 进度条
	builder.WriteString("Progress Bar:\n")
	builder.WriteString(fmt.Sprintf("  \x1b[32m%s\x1b[0m %.0f%%", ProgressBar(a.progress, 20), a.progress*100))
	builder.WriteString("\n")

	// 动画状态
	builder.WriteString("Animation State:\n")
	builder.WriteString("  \x1b[33mFPS:\x1b[0m 10\n")
	builder.WriteString(fmt.Sprintf("  \x1b[33mFrame:\x1b[0m %d", a.frame))
	builder.WriteString(fmt.Sprintf("  \x1b[33mProgress:\x1b[0m \x1b[32m%.2f\x1b[0m", a.progress))
	builder.WriteString("\n")

	// 帧序列
	builder.WriteString("Frame Sequence:\n")
	builder.WriteString(fmt.Sprintf("  \x1b[2m%v\x1b[0m", spinnerFrames))
	builder.WriteString("\n")

	// 退出提示
	builder.WriteString("\x1b[2mPress Ctrl+C to exit\x1b[0m")

	return builder.String()
}

func main() {
	// 使用 Ink，配置增量渲染和备用屏幕缓冲
	opts := &ink.RenderOptions{
		IncrementalRendering: true, // 启用增量渲染
		AlternateScreen:      true, // 使用备用屏幕缓冲
		ExitOnCtrlC:          true, // 启用 Ctrl+C 退出
		Interactive:          true, // 交互模式
	}

	// 创建动画组件
	app := &AnimationComponent{
		frame:    0,
		progress: 0.0,
	}

	// 使用 Ink 渲染应用
	instance := ink.Render(app, opts)

	// 启动动画更新循环（独立 goroutine）
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 更新动画状态
				app.frame++
				app.progress += 0.01
				if app.progress > 1.0 {
					app.progress = 0.0
				}
				// 触发重新渲染
				instance.Rerender(app)
			}
		}
	}()

	// 等待用户退出（Ctrl+C 或程序结束）- 库内部处理信号
	result := instance.WaitUntilExit()

	// 清理资源
	instance.Cleanup()
	if result != nil {
		fmt.Println("\n程序已退出:", result)
	}
}

// 确保 strings 包被使用
var _ = strings.TrimSpace("")