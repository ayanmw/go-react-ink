// Todo List 示例 - 演示列表渲染和状态管理
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

type Todo struct {
	ID    int
	Text  string
	Done  bool
}

func main() {
	ctx := hooks.NewHookContext()

	// 使用 useState 管理 todo 列表
	todos, _ := hooks.UseState(ctx, []Todo{
		{ID: 1, Text: "Learn Go-Ink", Done: true},
		{ID: 2, Text: "Build awesome CLI apps", Done: false},
		{ID: 3, Text: "Share with the community", Done: false},
	})

	// 渲染 todo 列表
	children := []core.Element{
		components.Text(core.Props{
			"children": "📝 Todo List",
			"color":    "blue",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),
	}

	for _, todo := range todos.([]Todo) {
		status := "○"
		color := "white"
		if todo.Done {
			status = "✓"
			color = "green"
		}

		children = append(children, components.Box(core.Props{
			"flexDirection": "row",
			"margin":        1,
		}, []core.Element{
			components.Text(core.Props{
				"children": status,
				"color":    color,
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children": todo.Text,
				"color":    color,
			}, nil),
		}))
	}

	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "single",
		"borderColor":   "blue",
	}, children)

	fmt.Println(app.Render())
}
