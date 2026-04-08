// Hooks Demo 示例 - 演示所有 React Hooks
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

	// useState - 状态管理
	count, setCount := hooks.UseState(ctx, 0)
	setCount(count.(int) + 1)

	// useRef - 引用
	ref := hooks.UseRef(ctx, "initial")
	ref.Current = "updated"

	// useMemo - 记忆化计算
	expensiveValue := hooks.UseMemo(ctx, func() any {
		// 复杂计算
		result := 0
		for i := 0; i < 100; i++ {
			result += i
		}
		return result
	}, []any{"dep1", "dep2"})

	// useCallback - 回调记忆化
	callback := hooks.UseCallback(ctx, func() {
		fmt.Println("Callback called")
	}, []any{"callback-dep"})

	// useInput - 输入处理
	inputHook := input.UseInput(ctx, func(k input.Key) {
		fmt.Printf("Key pressed: %s\n", k.Name)
	})
	inputHook.Handle("escape", func(k input.Key) {
		fmt.Println("Escape pressed - exit")
	})

	// useFocus - 焦点管理
	isFocused, focus := input.UseFocus(ctx)
	_ = isFocused
	focus()

	// useAnimation - 动画
	anim := input.UseAnimation(ctx, 30)
	anim.Play()
	anim.NextFrame()

	// useWindowSize - 窗口尺寸
	size := input.UseWindowSize(ctx)

	// 渲染演示 UI
	app := components.Box(core.Props{
		"flexDirection": "column",
		"padding":       1,
		"borderStyle":   "round",
		"borderColor":   "cyan",
	}, []core.Element{
		components.Text(core.Props{
			"children": "🎯 React Hooks Demo",
			"color":    "cyan",
			"bold":     true,
		}, nil),
		components.Newline(core.Props{}, nil),

		// useState
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useState: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("count = %d", count)}, nil),
		}),

		// useRef
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useRef: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("ref.Current = %s", ref.Current)}, nil),
		}),

		// useMemo
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useMemo: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("computed = %d", expensiveValue)}, nil),
		}),

		// useCallback
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useCallback: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("callback != nil: %v", callback != nil)}, nil),
		}),

		// useInput
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useInput: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": "listening for key events"}, nil),
		}),

		// useFocus
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useFocus: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("focused: %v", isFocused)}, nil),
		}),

		// useAnimation
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useAnimation: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("frame: %d, playing: %v", anim.Frame(), anim.IsPlaying())}, nil),
		}),

		// useWindowSize
		components.Box(core.Props{"flexDirection": "row"}, []core.Element{
			components.Text(core.Props{"children": "useWindowSize: ", "color": "yellow"}, nil),
			components.Text(core.Props{"children": fmt.Sprintf("%dx%d", size.Columns(), size.Rows())}, nil),
		}),
	})

	fmt.Println(app.Render())
}
