// Text Styling 示例 - 演示文本样式属性
package main

import (
	"fmt"

	"github.com/ayanmw/go-react-ink/pkg/components"
	"github.com/ayanmw/go-react-ink/pkg/core"
)

func main() {
	// 颜色示例
	colors := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	colorExamples := []core.Element{
		components.Text(core.Props{"children": "Colors:", "bold": true}, nil),
	}

	for _, color := range colors {
		colorExamples = append(colorExamples, components.Box(core.Props{
			"flexDirection": "row",
			"margin":        1,
		}, []core.Element{
			components.Text(core.Props{
				"children":        "████",
				"color":           color,
				"backgroundColor": "black",
			}, nil),
			components.Text(core.Props{"children": color}, nil),
		}))
	}

	// 样式属性示例
	styleExamples := []core.Element{
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Styles:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Normal", "color": "white"}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Bold", "bold": true}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Italic", "italic": true}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Underline", "underline": true}, nil),
		}),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{"children": "Dim", "dim": true}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Blink", "blink": true}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{"children": "Reverse", "reverse": true}, nil),
		}),
	}

	// 组合样式示例
	combinedExamples := []core.Element{
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Combined:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{
				"children": "Bold+Red",
				"bold":     true,
				"color":    "red",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children":  "Italic+Green",
				"italic":    true,
				"color":     "green",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children":   "Underline+Blue+Bold",
				"underline":  true,
				"color":      "blue",
				"bold":       true,
			}, nil),
		}),
	}

	// 背景色示例
	bgExamples := []core.Element{
		components.Newline(core.Props{}, nil),
		components.Text(core.Props{"children": "Background:", "bold": true}, nil),
		components.Box(core.Props{"flexDirection": "row", "margin": 1}, []core.Element{
			components.Text(core.Props{
				"children":        " Red BG ",
				"backgroundColor": "red",
				"color":           "white",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children":        " Green BG ",
				"backgroundColor": "green",
				"color":           "black",
			}, nil),
			components.Spacer(core.Props{}, nil),
			components.Text(core.Props{
				"children":        " Blue BG ",
				"backgroundColor": "blue",
				"color":           "white",
			}, nil),
		}),
	}

	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "double",
		"borderColor":   "magenta",
	}, append(
		append(
			append(colorExamples, styleExamples...),
			combinedExamples...,
		),
		bgExamples...,
	))

	fmt.Println(app.Render())
}
