// Counter 应用示例 - 演示 useState Hook
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

func main() {
	ctx := hooks.NewHookContext()

	// 使用 useState 管理计数器状态
	count, setCount := hooks.UseState(ctx, 0)

	// 渲染计数器 UI
	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
	}, []core.Element{
		components.Text(core.Props{
			"children": "Counter App",
			"color":    "blue",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{
			"children": fmt.Sprintf("Count: %d", count),
			"color":    "green",
		}, nil),
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{
			"children": "Press + to increment, - to decrement",
			"dim":      true,
		}, nil),
	})

	fmt.Println(app.Render())

	// 演示状态更新
	setCount(1)
	setCount(2)
	setCount(3)
}
