// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ayanmw/go-react-ink/pkg/core"
	"github.com/ayanmw/go-react-ink/pkg/fiber"
	"github.com/ayanmw/go-react-ink/pkg/hooks"
	_ "github.com/ayanmw/go-react-ink/pkg/input"
)

// Ink 应用实例
type Ink struct {
	// 核心组件
	rootNode  core.Element
	scheduler *fiber.Scheduler
	log       *LogUpdate

	// 状态
	isRunning    bool
	isUnmounted  bool
	isUnmounting bool
	lastOutput   string
	lastHeight   int

	// 渠道
	exitChan   chan any
	renderChan chan struct{}
	doneChan   chan struct{}

	// 选项
	options *RenderOptions

	// 动画管理
	animationManager *AnimationManager

	// 互斥锁
	mu sync.Mutex
}

// AnimationManager 动画管理器
type AnimationManager struct {
	subscribers []*AnimationSubscriber
	mu          sync.Mutex
}

// AnimationSubscriber 动画订阅者
type AnimationSubscriber struct {
	Fps       int
	LastTick  time.Time
	Callback  func(int)
	Frame     int
	IsPlaying bool
}

// NewAnimationManager 创建动画管理器
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		subscribers: make([]*AnimationSubscriber, 0),
	}
}

// Subscribe 注册动画订阅
func (m *AnimationManager) Subscribe(subscriber *AnimationSubscriber) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers = append(m.subscribers, subscriber)
}

// Unsubscribe 移除动画订阅
func (m *AnimationManager) Unsubscribe(subscriber *AnimationSubscriber) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, s := range m.subscribers {
		if s == subscriber {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
			break
		}
	}
}

// Tick 推进动画帧
func (m *AnimationManager) Tick() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for _, sub := range m.subscribers {
		if sub.IsPlaying {
			// 检查是否需要推进帧
			frameDuration := time.Second / time.Duration(sub.Fps)
			if now.Sub(sub.LastTick) >= frameDuration {
				sub.Frame++
				sub.LastTick = now
				if sub.Callback != nil {
					sub.Callback(sub.Frame)
				}
			}
		}
	}
}

// HasActiveAnimations 检查是否有活跃的动画
func (m *AnimationManager) HasActiveAnimations() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sub := range m.subscribers {
		if sub.IsPlaying {
			return true
		}
	}
	return false
}

// NewInk 创建新的 Ink 实例
func NewInk(opts *RenderOptions) *Ink {
	opts = applyDefaults(opts)

	// 自动检测 TTY：如果 stdout 不是终端且 Interactive 为默认值，则禁用交互模式
	if opts.Interactive && opts.Stdout != nil {
		// 检查是否为 TTY
		if f, ok := opts.Stdout.(*os.File); ok {
			if !IsTerminal(f) {
				opts.Interactive = false
			}
		}
	}

	return &Ink{
		options:          opts,
		scheduler:        fiber.NewScheduler(),
		log:              NewLogUpdate(opts.Stdout, opts.IncrementalRendering, opts.AlternateScreen),
		animationManager: NewAnimationManager(),
		exitChan:         make(chan any, 1),
		renderChan:       make(chan struct{}, 1),
		doneChan:         make(chan struct{}),
	}
}

// Render 渲染元素
func (ink *Ink) Render(element core.Element) {
	ink.mu.Lock()
	defer ink.mu.Unlock()

	// 存储 element 作为 any 类型
	ink.rootNode = element
	ink.scheduleRender()
}

// scheduleRender 调度渲染
func (ink *Ink) scheduleRender() {
	select {
	case ink.renderChan <- struct{}{}:
	default:
		// 已经有渲染请求在队列中，跳过
	}
}

// StartRenderLoop 启动渲染循环
func (ink *Ink) StartRenderLoop() {
	ink.mu.Lock()
	if ink.isRunning {
		ink.mu.Unlock()
		return
	}
	ink.isRunning = true
	ink.mu.Unlock()

	// 初始化终端 (使用备用屏幕缓冲区)
	if ink.options.Interactive {
		ink.log.Initialize()
	}

	// 创建节流定时器
	throttle := time.NewTicker(ink.options.FrameDuration())
	defer throttle.Stop()

	// 处理信号
	if ink.options.ExitOnCtrlC {
		go ink.handleSignals()
	}

	for {
		select {
		case <-ink.renderChan:
			// 调度的更新
			ink.onRender()

		case <-throttle.C:
			// 动画帧
			if ink.animationManager.HasActiveAnimations() {
				ink.animationManager.Tick()
				ink.scheduleRender()
			}

		case result := <-ink.exitChan:
			// 退出请求
			ink.finishUnmount(result)
			return

		case <-ink.doneChan:
			// 完成信号
			return
		}
	}
}

// onRender 执行渲染
func (ink *Ink) onRender() {
	ink.mu.Lock()
	defer ink.mu.Unlock()

	if ink.isUnmounted || ink.isUnmounting || ink.rootNode == nil {
		return
	}

	// 调度器执行工作循环
	output := ink.scheduler.Render(ink.rootNode)

	// 检查输出是否变化
	if output == ink.lastOutput {
		return
	}

	// 写入输出
	ink.log.Write(output)
	ink.lastOutput = output
	ink.lastHeight = countLines(output)
}

// Unmount 卸载应用
func (ink *Ink) Unmount(result ...any) {
	ink.mu.Lock()
	if ink.isUnmounted || ink.isUnmounting {
		ink.mu.Unlock()
		return
	}
	ink.isUnmounting = true
	ink.mu.Unlock()

	// 最终渲染
	ink.onRender()

	// 恢复终端状态
	ink.log.Done()

	// 发送退出信号
	var exitResult any
	if len(result) > 0 {
		exitResult = result[0]
	}
	ink.exitChan <- exitResult

	ink.mu.Lock()
	ink.isUnmounted = true
	ink.mu.Unlock()
}

// WaitUntilExit 等待退出
func (ink *Ink) WaitUntilExit() any {
	result := <-ink.exitChan
	return result
}

// finishUnmount 完成卸载
func (ink *Ink) finishUnmount(result any) {
	ink.mu.Lock()
	ink.isUnmounted = true
	ink.isRunning = false
	ink.mu.Unlock()

	// 关闭完成渠道
	close(ink.doneChan)
}

// handleSignals 处理信号
func (ink *Ink) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				ink.Unmount()
				return
			}
		case <-ink.doneChan:
			return
		}
	}
}

// GetAnimationManager 获取动画管理器
func (ink *Ink) GetAnimationManager() *AnimationManager {
	return ink.animationManager
}

// Clear 清除输出
func (ink *Ink) Clear() {
	ink.mu.Lock()
	defer ink.mu.Unlock()
	ink.log.Clear()
	ink.lastOutput = ""
	ink.lastHeight = 0
}

// Cleanup 清理资源
func (ink *Ink) Cleanup() {
	ink.mu.Lock()
	defer ink.mu.Unlock()

	// 关闭所有渠道
	close(ink.doneChan)
	close(ink.renderChan)
}

// Instance 应用实例接口
type Instance struct {
	// Rerender 重新渲染
	Rerender func(element core.Element)

	// Unmount 卸载
	Unmount func(result ...any)

	// WaitUntilExit 等待退出
	WaitUntilExit func() any

	// WaitUntilRenderFlush 等待渲染完成
	WaitUntilRenderFlush func() error

	// Cleanup 清理
	Cleanup func()

	// Clear 清除
	Clear func()

	// HookContext Hook 上下文
	HookContext *hooks.HookContext
}

// WaitUntilRenderFlush 等待渲染完成
func (ink *Ink) WaitUntilRenderFlush() error {
	// 触发一次渲染请求
	ink.scheduleRender()

	// 等待渲染完成 (简化实现)
	time.Sleep(ink.options.FrameDuration())
	return nil
}