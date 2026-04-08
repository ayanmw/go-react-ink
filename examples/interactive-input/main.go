// Interactive Input 示例 - 演示输入处理和焦点管理
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/input"
)

func main() {
	ctx := hooks.NewHookContext()

	// 输入处理
	keyLog := []string{}
	inputHook := input.UseInput(ctx, func(k input.Key) {
		keyLog = append(keyLog, k.Name)
	})

	// 注册快捷键
	inputHook.Handle("up", func(k input.Key) {
		fmt.Println("↑ Moving up")
	})
	inputHook.Handle("down", func(k input.Key) {
		fmt.Println("↓ Moving down")
	})
	inputHook.Handle("left", func(k input.Key) {
		fmt.Println("← Moving left")
	})
	inputHook.Handle("right", func(k input.Key) {
		fmt.Println("→ Moving right")
	})
	inputHook.Handle("return", func(k input.Key) {
		fmt.Println("⏎ Enter pressed")
	})
	inputHook.Handle("escape", func(k input.Key) {
		fmt.Println("⎋ Escape pressed - exiting")
	})
	inputHook.Handle("tab", func(k input.Key) {
		fmt.Println("⇥ Tab pressed")
	})

	// 模拟按键
	inputHook.ProcessKey(input.Key{Name: "up"})
	inputHook.ProcessKey(input.Key{Name: "right"})
	inputHook.ProcessKey(input.Key{Name: "return"})

	// 焦点管理
	focusManager := input.UseFocusManager(ctx)
	focusManager.Register("input1")
	focusManager.Register("input2")
	focusManager.Register("button1")
	focusManager.Focus("input1")

	// 光标
	cursor := input.UseCursor(ctx)
	cursor.SetPosition(10, 5)
	cursor.SetShape("underline")

	// 标准输入/输出
	stdin := input.UseStdin(ctx)
	stdin.SetRaw(true)

	stdout := input.UseStdout(ctx)

	// 窗口尺寸
	windowSize := input.UseWindowSize(ctx)

	// 渲染 UI
	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "single",
		"borderColor":   "green",
	}, []core.Element{
		components.Text(core.Props{
			"children": "⌨️ Interactive Input Demo",
			"color":    "green",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),

		// 键盘事件
		components.Text(core.Props{"children": "Keyboard Events:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "↑↓←→", "color": "cyan"}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Arrow keys for navigation"}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "⏎", "color": "cyan"}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Enter to confirm"}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "⎋", "color": "cyan"}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Escape to exit"}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "⇥", "color": "cyan"}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Tab to cycle focus"}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 焦点状态
		components.Text(core.Props{"children": "Focus State:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Active: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": "input1", "bold": true, "color": "green"}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 终端信息
		components.Text(core.Props{"children": "Terminal Info:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Size: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%dx%d", stdout.Width(), stdout.Height())}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Window: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%dx%d", windowSize.Columns(), windowSize.Rows())}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "TTY: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%v", stdin.IsTTY())}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Raw Mode: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%v", stdin.IsRaw())}, nil),
		}),

		components.Newline(core.Props{}, nil),

		// 按键日志
		components.Text(core.Props{"children": "Key Log:", "bold": true}, nil),
		components.Text(core.Props{
			"children": fmt.Sprintf("%v", keyLog),
			"dim":      true,
		}, nil),
	})

	fmt.Println(app.Render())
}
