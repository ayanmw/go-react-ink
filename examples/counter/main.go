package main

import (
	"fmt"

	"github.com/anmingwei/go-ink/pkg/components"
	"github.com/anmingwei/go-ink/pkg/core"
	"github.com/anmingwei/go-ink/pkg/hooks"
	"github.com/anmingwei/go-ink/pkg/input"
)

// 模拟一个完整的终端应用

func main() {
	// 创建 Hook 上下文
	ctx := hooks.NewHookContext()

	// 创建应用
	app := createApp(ctx)

	// 渲染输出
	output := app.Render()
	fmt.Println(output)
}

func createApp(ctx *hooks.HookContext) core.Element {
	// 使用状态
	count, setCount := hooks.UseState(ctx, 0)

	// 模拟输入处理
	inputHook := input.UseInput(ctx, func(key input.Key) {
		switch key.Name {
		case "up":
			currentCount := count.(int)
			setCount(currentCount + 1)
		case "down":
			currentCount := count.(int)
			if currentCount > 0 {
				setCount(currentCount - 1)
			}
		}
	})
	_ = inputHook

	// 使用副作用
	hooks.UseEffect(ctx, func() func() {
		fmt.Println("App mounted")
		return func() {
			fmt.Println("App unmounted")
		}
	}, nil)

	// 使用记忆化
	message := hooks.UseMemo(ctx, func() any {
		c := count.(int)
		if c == 0 {
			return "Click Up to increase"
		} else if c >= 10 {
			return "You've reached 10!"
		}
		return fmt.Sprintf("Count: %d", c)
	}, []any{count})

	// 返回 UI 结构
	return core.CreateElement(
		components.Box,
		core.Props{
			"flexDirection":    "column",
			"padding":          1,
			"borderStyle":      "single",
			"borderColor":      "green",
		},
		[]core.Element{
			// 标题
			core.CreateElement(
				components.Text,
				core.Props{
					"children": "Counter App",
					"bold":     true,
					"color":    "cyan",
				},
				nil,
			),
			// 空行
			core.CreateElement(components.Newline, nil, nil),
			// 内容
			core.CreateElement(
				components.Text,
				core.Props{
					"children": message.(string),
				},
				nil,
			),
			// 空行
			core.CreateElement(components.Newline, nil, nil),
			// 提示
			core.CreateElement(
				components.Text,
				core.Props{
					"children": "Use Up/Down arrows to change",
					"dim":      true,
				},
				nil,
			),
		},
	)
}