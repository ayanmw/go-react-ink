package main

import (
	"fmt"

	"github.com/anmingwei/go-ink/pkg/components"
	"github.com/anmingwei/go-ink/pkg/core"
	"github.com/anmingwei/go-ink/pkg/layout"
)

// 布局示例：展示 Flexbox 布局功能

func main() {
	// 创建布局
	root := createLayout()

	// 计算布局
	root.CalculateLayout(80, 24)

	// 打印布局结果
	printLayout(root, 0)
}

func createLayout() *layout.Node {
	// 根容器
	root := layout.NewNode()
	root.Direction = layout.DirectionColumn
	root.SetPadding(1, 2, 1, 2)

	// 顶部栏
	header := layout.NewNode()
	header.Height = 1
	header.Direction = layout.DirectionRow
	header.Justify = layout.JustifyCenter

	headerTitle := layout.NewNode()
	headerTitle.Width = 20
	headerTitle.Height = 1
	header.AddChild(headerTitle)

	// 主内容区
	content := layout.NewNode()
	content.FlexGrow = 1
	content.Direction = layout.DirectionRow

	// 侧边栏
	sidebar := layout.NewNode()
	sidebar.Width = 15
	sidebar.SetPadding(0, 1, 0, 0)

	sidebarItem1 := layout.NewNode()
	sidebarItem1.Height = 1
	sidebar.AddChild(sidebarItem1)

	sidebarItem2 := layout.NewNode()
	sidebarItem2.Height = 1
	sidebar.AddChild(sidebarItem2)

	sidebarItem3 := layout.NewNode()
	sidebarItem3.Height = 1
	sidebar.AddChild(sidebarItem3)

	// 主区域
	main := layout.NewNode()
	main.FlexGrow = 1
	main.Direction = layout.DirectionColumn
	main.SetPadding(0, 0, 0, 1)

	mainContent := layout.NewNode()
	mainContent.FlexGrow = 1
	main.AddChild(mainContent)

	mainFooter := layout.NewNode()
	mainFooter.Height = 1
	main.AddChild(mainFooter)

	content.AddChild(sidebar)
	content.AddChild(main)

	// 底部栏
	footer := layout.NewNode()
	footer.Height = 1
	footer.Direction = layout.DirectionRow
	footer.Justify = layout.JustifySpaceBetween

	footerLeft := layout.NewNode()
	footerLeft.Width = 20
	footerLeft.Height = 1
	footer.AddChild(footerLeft)

	footerRight := layout.NewNode()
	footerRight.Width = 20
	footerRight.Height = 1
	footer.AddChild(footerRight)

	root.AddChild(header)
	root.AddChild(content)
	root.AddChild(footer)

	return root
}

func printLayout(node *layout.Node, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	fmt.Printf("%sNode: (%.0f, %.0f) %.0fx%.0f\n",
		indent,
		node.Layout.X,
		node.Layout.Y,
		node.Layout.Width,
		node.Layout.Height,
	)

	for _, child := range node.Children {
		printLayout(child, depth+1)
	}
}

// 使用组件 API 创建 UI
func createUI() core.Element {
	return core.CreateElement(
		components.Box,
		core.Props{
			"flexDirection": "column",
			"padding":       1,
		},
		[]core.Element{
			// Header
			core.CreateElement(
				components.Box,
				core.Props{
					"justifyContent": "center",
					"paddingBottom":  1,
				},
				[]core.Element{
					core.CreateElement(
						components.Text,
						core.Props{
							"children": "Layout Demo",
							"bold":     true,
							"color":    "cyan",
						},
						nil,
					),
				},
			),
			// Content
			core.CreateElement(
				components.Box,
				core.Props{
					"flexGrow": 1,
				},
				[]core.Element{
					// Sidebar
					core.CreateElement(
						components.Box,
						core.Props{
							"width":        15,
							"paddingRight": 1,
							"borderRight":  true,
						},
						[]core.Element{
							core.CreateElement(components.Text, core.Props{"children": "Menu Item 1"}, nil),
							core.CreateElement(components.Text, core.Props{"children": "Menu Item 2"}, nil),
							core.CreateElement(components.Text, core.Props{"children": "Menu Item 3"}, nil),
						},
					),
					// Main
					core.CreateElement(
						components.Box,
						core.Props{
							"flexGrow":     1,
							"paddingLeft":  1,
							"flexDirection": "column",
						},
						[]core.Element{
							core.CreateElement(
								components.Box,
								core.Props{"flexGrow": 1},
								[]core.Element{
									core.CreateElement(
										components.Text,
										core.Props{"children": "Main Content Area"},
										nil,
									),
								},
							),
							core.CreateElement(
								components.Text,
								core.Props{"children": "Status: Ready", "dim": true},
								nil,
							),
						},
					),
				},
			),
			// Footer
			core.CreateElement(
				components.Box,
				core.Props{
					"justifyContent": "space-between",
					"paddingTop":     1,
				},
				[]core.Element{
					core.CreateElement(components.Text, core.Props{"children": "Left Footer"}, nil),
					core.CreateElement(components.Text, core.Props{"children": "Right Footer"}, nil),
				},
			),
		},
	)
}