// Flexbox Layout 示例 - 演示 Flex 布局属性
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

func main() {
	// 水平布局示例
	rowExample := components.Box(core.Props{
		"flexDirection":  "row",
		"justifyContent": "space-between",
		"padding":        1,
		"width":          40,
	}, []core.Element{
		components.Text(core.Props{"children": "Left", "color": "red"}, nil),
		components.Text(core.Props{"children": "Center", "color": "green"}, nil),
		components.Text(core.Props{"children": "Right", "color": "blue"}, nil),
	})

	// 垂直布局示例
	columnExample := components.Box(core.Props{
		"flexDirection":  "column",
		"alignItems":     "center",
		"padding":        1,
	}, []core.Element{
		components.Text(core.Props{"children": "Top", "color": "red"}, nil),
		components.Text(core.Props{"children": "Middle", "color": "green"}, nil),
		components.Text(core.Props{"children": "Bottom", "color": "blue"}, nil),
	})

	// 嵌套布局示例
	nestedExample := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "single",
	}, []core.Element{
		components.Text(core.Props{"children": "Header", "bold": true}, nil),
		components.Box(core.Props{
			"flexDirection": "row",
			"flexGrow":      1,
		}, []core.Element{
			components.Box(core.Props{
				"width":   10,
				"padding": 1,
			}, []core.Element{
				components.Text(core.Props{"children": "Sidebar"}, nil),
			}),
			components.Box(core.Props{
				"flexGrow": 1,
				"padding":  1,
			}, []core.Element{
				components.Text(core.Props{"children": "Content"}, nil),
			}),
		}),
		components.Text(core.Props{"children": "Footer", "dim": true}, nil),
	})

	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
	}, []core.Element{
		components.Text(core.Props{
			"children": "Flexbox Layout Examples",
			"color":    "cyan",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Row Layout:", "bold": true}, nil),
		rowExample,
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Column Layout:", "bold": true}, nil),
		columnExample,
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Nested Layout:", "bold": true}, nil),
		nestedExample,
	})

	fmt.Println(app.Render())
}
