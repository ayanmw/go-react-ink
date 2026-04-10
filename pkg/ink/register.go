// Package ink 提供输入事件注册接口
package ink

// InputEventHandler 输入事件处理器类型
type InputEventHandler func(input string)

// InputEventEmitter 输入事件发射器接口
// 由 Ink 实现以避免循环依赖
type InputEventEmitter interface {
	// EmitInput 发射输入事件
	EmitInput(input string)
}

// InputRegistry 输入注册器接口
// 用于非循环依赖地注册输入处理器
type InputRegistry interface {
	// RegisterInputHandler 注册输入处理器
	RegisterInputHandler(handler InputEventHandler) func()
}

// inputEventEmitterHalts 输入事件发射器的默认空实现
type inputEventEmitterHalts struct{}

// EmitInput 默认不发射事件
func (h *inputEventEmitterHalts) EmitInput(input string) {}

// inputRegistryNil 输入注册器的默认空实现
type inputRegistryNil struct{}

// RegisterInputHandler 默认不注册，返回无操作的清理函数
func (r *inputRegistryNil) RegisterInputHandler(handler InputEventHandler) func() {
	// 返回一个无操作的清理函数
	return func() {}
}

// DefaultInputEventEmitter 默认输入事件发射器
var DefaultInputEventEmitter InputEventEmitter = &inputEventEmitterHalts{}

// DefaultInputRegistry 默认输入注册器
var DefaultInputRegistry InputRegistry = &inputRegistryNil{}

// SetInputEventEmitter 设置输入事件发射器（由 Ink 实例在构造时调用）
func SetInputEventEmitter(emitter InputEventEmitter) {
	DefaultInputEventEmitter = emitter
}

// SetInputRegistry 设置输入注册器（由 Ink 实例在构造时调用）
func SetInputRegistry(registry InputRegistry) {
	DefaultInputRegistry = registry
}

// EmitInput 便利函数 - 发射输入事件
func EmitInput(input string) {
	if DefaultInputEventEmitter != nil {
		DefaultInputEventEmitter.EmitInput(input)
	}
}
