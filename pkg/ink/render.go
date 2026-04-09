// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
)

// Render 渲染应用并返回实例
func Render(element core.Element, opts *RenderOptions) *Instance {
	// 应用默认值
	opts = applyDefaults(opts)

	// 创建 Ink 实例
	ink := NewInk(opts)

	// 设置 Hook 上下文
	ctx := hooks.NewHookContext()

	// 初始渲染
	ink.Render(element)

	// 启动渲染循环
	go ink.StartRenderLoop()

	// 创建实例接口
	instance := &Instance{
		Rerender: func(el core.Element) {
			ink.Render(el)
		},
		Unmount: func(result ...any) {
			ink.Unmount(result...)
		},
		WaitUntilExit: func() any {
			return ink.WaitUntilExit()
		},
		WaitUntilRenderFlush: func() error {
			return ink.WaitUntilRenderFlush()
		},
		Cleanup: func() {
			ink.Cleanup()
		},
		Clear: func() {
			ink.Clear()
		},
		HookContext: ctx,
	}

	// 存储实例到全局 (用于 UseApp 等 Hooks)
	SetCurrentInstance(instance)

	return instance
}

// currentInstance 当前实例 (全局)
var currentInstance *Instance

// instanceMutex 实例互斥锁
var instanceMutex = hooksMutex{}

// hooksMutex 简化的互斥锁类型
type hooksMutex struct{}

// SetCurrentInstance 设置当前实例
func SetCurrentInstance(instance *Instance) {
	currentInstance = instance
}

// GetCurrentInstance 获取当前实例
func GetCurrentInstance() *Instance {
	return currentInstance
}

// Exit 退出当前应用
func Exit(result ...any) {
	if currentInstance != nil {
		currentInstance.Unmount(result...)
	}
}