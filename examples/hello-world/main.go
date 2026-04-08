// Hello World 示例 - 最简单的 Go-Ink 应用
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

func main() {
	// 创建简单的文本组件
	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
	}, []core.Element{
		components.Text(core.Props{
			"children": "Hello, World!",
			"color":    "green",
			"bold":     true,
		}, nil),
		components.Text(core.Props{
			"children": "Welcome to Go-Ink!",
			"color":    "cyan",
		}, nil),
	})

	// 渲染输出
	fmt.Println(app.Render())
}
