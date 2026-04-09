// Animation Demo 示例 - 演示动画和帧更新
// 使用 ink.Render 实现真正的动画效果
package main

import (
	"fmt"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/ink"
	"github.com/ayanmw/go-react-ink/pkg/input"
)

// Spinner 动画帧
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Progress bar 组件
func ProgressBar(progress float64, width int) core.Element {
	filled := int(progress * float64(width))
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return components.Text(core.Props{
		"children": bar,
		"color":    "green",
	}, nil)
}

// App 主应用组件
func App(ctx *hooks.HookContext) core.Element {
	// 动画 Hook - 10 FPS
	anim := input.UseAnimationController(ctx, 10)
	anim.Play()

	// 进度状态 - 从 0 开始，每帧增加
	progress, setProgress := hooks.UseState(ctx, 0.0)

	// 帧计数
	frame, setFrame := hooks.UseState(ctx, 0)

	// 使用 UseEffect 监听动画帧变化
	hooks.UseEffect(ctx, func() func() {
		// 每 100ms 检查动画帧并更新进度
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					currentFrame := anim.Frame()
					setFrame(currentFrame)

					// 每帧增加进度
					currentProgress := progress.(float64)
					newProgress := currentProgress + 0.02
					if newProgress > 1.0 {
						newProgress = 0.0 // 循环进度条
					}
					setProgress(newProgress)
				}
			}
		}()
		return func() {}
	}, nil)

	// 应用退出控制
	app := input.UseApp(ctx)

	// 输入处理 - q 或 escape 退出
	inputHook := input.UseInput(ctx, func(key input.Key) {
		if key.Name == "q" || key.Name == "escape" {
			app.Exit()
		}
	})
	_ = inputHook // 使用 hook

	// 渲染 UI
	spinnerFrame := spinnerFrames[frame.(int)%len(spinnerFrames)]

	return components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "round",
		"borderColor":   "yellow",
	}, []core.Element{
		components.Text(core.Props{
			"children": "🎬 Animation Demo (Real Animation)",
			"color":    "yellow",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),

		// Spinner 动画 - 实时更新
		components.Text(core.Props{"children": "Spinner:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{
				"children": spinnerFrame,
				"color":    "cyan",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Loading..."}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 进度条 - 实时更新
		components.Text(core.Props{"children": "Progress Bar:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			ProgressBar(progress.(float64), 20),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children": fmt.Sprintf("%.0f%%", progress.(float64)*100),
				"color":    "green",
			}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 动画状态
		components.Text(core.Props{"children": "Animation State:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "FPS: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": "10"}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Frame: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%d", frame.(int))}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Progress: ", "color": "yellow"}, nil),
			components.Text(core.Props{
				"children": fmt.Sprintf("%.2f", progress.(float64)),
				"color":    "green",
			}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 帧序列
		components.Text(core.Props{"children": "Frame Sequence:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{
				"children": fmt.Sprintf("%v", spinnerFrames),
				"dim":      true,
			}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 退出提示
		components.Text(core.Props{
			"children": "Press 'q' or 'Escape' to exit",
			"dim":      true,
		}, nil),
	})
}

func main() {
	// 创建 Hook 上下文
	ctx := hooks.NewHookContext()

	// 渲染应用并启动动画循环
	instance := ink.Render(App(ctx), &ink.RenderOptions{
		MaxFps:      30,
		Interactive: true,
		ExitOnCtrlC: true,
	})

	// 等待退出
	instance.WaitUntilExit()

	fmt.Println("Animation demo exited cleanly")
}