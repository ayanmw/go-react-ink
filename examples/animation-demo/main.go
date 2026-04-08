// Animation Demo 示例 - 演示动画和帧更新
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
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

func main() {
	ctx := hooks.NewHookContext()

	// 动画 Hook
	anim := input.UseAnimation(ctx, 10) // 10 FPS
	anim.Play()

	// 模拟几帧
	for i := 0; i < 5; i++ {
		anim.NextFrame()
	}

	// 进度状态
	progress, setProgress := hooks.UseState(ctx, 0.0)
	setProgress(0.65)

	// 渲染动画 UI
	frame := spinnerFrames[anim.Frame()%len(spinnerFrames)]

	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "round",
		"borderColor":   "yellow",
	}, []core.Element{
		components.Text(core.Props{
			"children": "🎬 Animation Demo",
			"color":    "yellow",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),

		// Spinner 动画
		components.Text(core.Props{"children": "Spinner:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{
				"children": frame,
				"color":    "cyan",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Loading..."}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 进度条
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
			components.Text(core.Props{"children": fmt.Sprintf("%d", anim.Frame())}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Playing: ", "color": "yellow"}, nil),
			components.Text(core.Props{
				"children": fmt.Sprintf("%v", anim.IsPlaying()),
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
	})

	fmt.Println(app.Render())
}
