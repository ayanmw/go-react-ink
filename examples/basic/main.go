package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

//go:generate go run github.com/ayanmw/go-react-ink/cmd/gox -o app_generated.go app.gox

func main() {
	// 手动创建 UI 结构 (模拟编译器输出)
	app := createApp()

	// 渲染输出
	output := app.Render()
	fmt.Println(output)
}

func createApp() core.Element {
	// 等价于:
	// <Box flexDirection="column">
	//   <Text color="green">Hello, World!</Text>
	//   <Text>Count: 42</Text>
	// </Box>

	return core.CreateElement(
		components.Box,
		core.Props{
			"flexDirection": "column",
		},
		[]core.Element{
			core.CreateElement(
				components.Text,
				core.Props{
					"children": "Hello, World!",
					"color":    "green",
				},
				nil,
			),
			core.CreateElement(
				components.Text,
				core.Props{
					"children": "Count: 42",
				},
				nil,
			),
		},
	)
}
