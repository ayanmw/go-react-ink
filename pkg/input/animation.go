// Package input 实现输入处理 Hooks
package input

import (
	"sync"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/hooks"
	"github.com/ayanmw/go-react-ink/pkg/ink"
)

// AnimationController 动画控制器
// 与渲染循环集成的动画系统
type AnimationController struct {
	frame     int
	isPlaying bool
	fps       int
	lastTick  time.Time
	subscriber *ink.AnimationSubscriber
	mu        sync.Mutex
}

// animationRegistry 动画注册表
var animationRegistry = struct {
	subscribers []*AnimationController
	mu          sync.Mutex
}{
	subscribers: make([]*AnimationController, 0),
}

// NewAnimationController 创建动画控制器
func NewAnimationController(fps int) *AnimationController {
	return &AnimationController{
		fps: fps,
	}
}

// Play 开始播放
func (a *AnimationController) Play() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isPlaying {
		return
	}

	a.isPlaying = true
	a.lastTick = time.Now()

	// 注册到动画管理器
	a.registerWithManager()
}

// Pause 暂停播放
func (a *AnimationController) Pause() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.isPlaying = false
}

// Stop 停止播放并重置帧
func (a *AnimationController) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.isPlaying = false
	a.frame = 0
}

// Reset 重置帧计数
func (a *AnimationController) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.frame = 0
}

// Frame 获取当前帧
func (a *AnimationController) Frame() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.frame
}

// IsPlaying 检查是否正在播放
func (a *AnimationController) IsPlaying() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isPlaying
}

// SetFrame 设置帧
func (a *AnimationController) SetFrame(frame int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.frame = frame
}

// NextFrame 推进到下一帧
func (a *AnimationController) NextFrame() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.frame++
}

// registerWithManager 注册到动画管理器
func (a *AnimationController) registerWithManager() {
	instance := ink.GetCurrentInstance()
	if instance == nil {
		return
	}

	// 创建订阅者
	a.subscriber = &ink.AnimationSubscriber{
		Fps:       a.fps,
		LastTick:  time.Now(),
		IsPlaying: true,
		Callback: func(frame int) {
			a.mu.Lock()
			a.frame = frame
			a.mu.Unlock()
		},
	}

	// 获取动画管理器并注册
	// 注意：这里需要 Ink 实例暴露 AnimationManager
	// 简化实现：使用全局注册表
	animationRegistry.mu.Lock()
	animationRegistry.subscribers = append(animationRegistry.subscribers, a)
	animationRegistry.mu.Unlock()
}

// unregisterFromManager 从动画管理器注销
func (a *AnimationController) unregisterFromManager() {
	animationRegistry.mu.Lock()
	defer animationRegistry.mu.Unlock()

	for i, sub := range animationRegistry.subscribers {
		if sub == a {
			animationRegistry.subscribers = append(
				animationRegistry.subscribers[:i],
				animationRegistry.subscribers[i+1:]...,
			)
			break
		}
	}
}

// UseAnimationController 动画控制器 Hook
// 返回一个与渲染循环集成的动画控制器
func UseAnimationController(ctx *hooks.HookContext, fps int) *AnimationController {
	// 使用 UseMemo 确保控制器只创建一次
	controller := hooks.UseMemo(ctx, func() any {
		return NewAnimationController(fps)
	}, []any{fps}).(*AnimationController)

	return controller
}

// TickAnimations 推进所有动画
// 由渲染循环调用
func TickAnimations() {
	animationRegistry.mu.Lock()
	defer animationRegistry.mu.Unlock()

	now := time.Now()
	for _, controller := range animationRegistry.subscribers {
		controller.mu.Lock()
		if controller.isPlaying {
			frameDuration := time.Second / time.Duration(controller.fps)
			if now.Sub(controller.lastTick) >= frameDuration {
				controller.frame++
				controller.lastTick = now
			}
		}
		controller.mu.Unlock()
	}
}

// HasActiveAnimations 检查是否有活跃的动画
func HasActiveAnimations() bool {
	animationRegistry.mu.Lock()
	defer animationRegistry.mu.Unlock()

	for _, controller := range animationRegistry.subscribers {
		if controller.isPlaying {
			return true
		}
	}
	return false
}
