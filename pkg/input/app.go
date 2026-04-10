// Package input 实现输入处理 Hooks
package input

import (
	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

// IInstance 应用实例接口 (避免循环依赖)
// 定义最小化的应用操作接口
type IInstance interface {
	// Exit 退出应用
	Exit(...any)
}

// AppController 应用控制器
// 提供对应用实例的控制方法
type AppController struct {
	instance IInstance
}

// NewAppController 创建应用控制器
func NewAppController(instance IInstance) *AppController {
	return &AppController{
		instance: instance,
	}
}

// Exit 退出应用
func (a *AppController) Exit(result ...any) {
	if a.instance != nil {
		a.instance.Exit(result...)
	}
}

// Rerender 重新渲染
func (a *AppController) Rerender(element any) {
	if a.instance != nil {
		// 注意：Rerender 需要完整的接口
		// 这里简化处理 - 实际使用时可以使用类型断言
	}
}

// WaitUntilExit 等待退出
func (a *AppController) WaitUntilExit() any {
	// 需要完整接口支持，这里返回 nil
	return nil
}

// Cleanup 清理资源
func (a *AppController) Cleanup() {
	// 需要完整接口支持，这里简化处理
}

// Clear 清除输出
func (a *AppController) Clear() {
	// 需要完整接口支持，这里简化处理
}

// UseApp 应用 Hook
// 返回应用控制器，用于控制应用实例（简化版本，不依赖 ink 包）
func UseAppSimple(ctx *hooks.HookContext) *AppController {
	// 使用 UseMemo 确保控制器只创建一次
	controller := hooks.UseMemo(ctx, func() any {
		// 从全局获取实例 - 避免直接依赖 ink 包
		return NewAppController(nil)
	}, nil).(*AppController)

	return controller
}
