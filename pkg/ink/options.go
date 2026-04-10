// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"io"
	"os"
	"time"
)

// RenderOptions 渲染选项
type RenderOptions struct {
	// Stdout 标准输出
	Stdout io.Writer

	// Stderr 标准错误输出
	Stderr io.Writer

	// ExitOnCtrlC 是否在 Ctrl+C 时退出
	ExitOnCtrlC bool

	// MaxFps 最大帧率 (默认: 30)
	MaxFps int

	// PatchConsole 是否修补控制台
	PatchConsole bool

	// Interactive 是否交互模式
	Interactive bool

	// AlternateScreen 是否使用备用屏幕缓冲
	AlternateScreen bool

	// IncrementalRendering 是否启用增量渲染 (默认: true)
	IncrementalRendering bool
}

// DefaultRenderOptions 默认渲染选项
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		ExitOnCtrlC:     true,
		MaxFps:          30,
		PatchConsole:    false,
		Interactive:           true,
		AlternateScreen:       false,
		IncrementalRendering:  true,
	}
}

// applyDefaults 应用默认值
func applyDefaults(opts *RenderOptions) *RenderOptions {
	if opts == nil {
		return DefaultRenderOptions()
	}

	// 应用默认值
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.MaxFps <= 0 {
		opts.MaxFps = 30
	}

	// 自动检测 TTY：如果 stdout 是终端且未显式设置 Interactive，则启用交互模式
	// 注意：这里无法区分"未设置"和"设置为false"，所以只在默认情况下检测
	// 用户可以通过显式设置 Interactive: false 来禁用

	return opts
}

// FrameDuration 返回每帧的持续时间
func (opts *RenderOptions) FrameDuration() time.Duration {
	if opts.MaxFps <= 0 {
		return time.Second / 30
	}
	return time.Second / time.Duration(opts.MaxFps)
}

// Merge 合并两个选项 (other 覆盖 opts)
func (opts *RenderOptions) Merge(other *RenderOptions) *RenderOptions {
	if other == nil {
		return opts
	}

	result := &RenderOptions{}

	// 复制 opts 的值
	if opts != nil {
		result.Stdout = opts.Stdout
		result.Stderr = opts.Stderr
		result.ExitOnCtrlC = opts.ExitOnCtrlC
		result.MaxFps = opts.MaxFps
		result.PatchConsole = opts.PatchConsole
		result.Interactive = opts.Interactive
		result.AlternateScreen = opts.AlternateScreen
	}

	// 用 other 的值覆盖
	if other.Stdout != nil {
		result.Stdout = other.Stdout
	}
	if other.Stderr != nil {
		result.Stderr = other.Stderr
	}
	if other.MaxFps > 0 {
		result.MaxFps = other.MaxFps
	}
	// 布尔值直接覆盖
	result.ExitOnCtrlC = other.ExitOnCtrlC
	result.PatchConsole = other.PatchConsole
	result.Interactive = other.Interactive
	result.AlternateScreen = other.AlternateScreen

	return result
}