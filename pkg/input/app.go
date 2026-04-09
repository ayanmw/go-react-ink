// Package input 实现输入处理 Hooks
package input

import (
	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/ink"
)

// AppController 应用控制器
// 提供对应用实例的控制方法
type AppController struct {
	instance *ink.Instance
}

// NewAppController 创建应用控制器
func NewAppController(instance *ink.Instance) *AppController {
	return &AppController{
		instance: instance,
	}
}

// Exit 退出应用
func (a *AppController) Exit(result ...any) {
	if a.instance != nil {
		a.instance.Unmount(result...)
	}
}

// Rerender 重新渲染
func (a *AppController) Rerender(element any) {
	if a.instance != nil {
		// 注意：Rerender 需要 core.Element 类型
		// 这里简化处理
	}
}

// WaitUntilExit 等待退出
func (a *AppController) WaitUntilExit() any {
	if a.instance != nil {
		return a.instance.WaitUntilExit()
	}
	return nil
}

// Cleanup 清理资源
func (a *AppController) Cleanup() {
	if a.instance != nil {
		a.instance.Cleanup()
	}
}

// Clear 清除输出
func (a *AppController) Clear() {
	if a.instance != nil {
		a.instance.Clear()
	}
}

// UseApp 应用 Hook
// 返回应用控制器，用于控制应用实例
func UseApp(ctx *hooks.HookContext) *AppController {
	// 使用 UseMemo 确保控制器只创建一次
	controller := hooks.UseMemo(ctx, func() any {
		instance := ink.GetCurrentInstance()
		return NewAppController(instance)
	}, nil).(*AppController)

	return controller
}

// UseAppWithInstance 使用指定实例的应用 Hook
func UseAppWithInstance(ctx *hooks.HookContext, instance *ink.Instance) *AppController {
	controller := hooks.UseMemo(ctx, func() any {
		return NewAppController(instance)
	}, []any{instance}).(*AppController)

	return controller
}
