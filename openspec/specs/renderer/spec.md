# renderer Specification

## Purpose
定义 Go-Ink 终端渲染器的设计和实现，包括增量渲染、备用屏幕缓冲、TTY 检测、信号处理和动画系统。

## Requirements

### Requirement: 渲染器实现增量渲染
渲染器 SHALL 只更新变化的部分，最小化终端输出。

#### Scenario: 双缓冲 Diff
- **WHEN** 组件状态变化触发重新渲染
- **THEN** 系统比较新旧缓冲区，只输出变化的单元格

#### Scenario: 输出未变化
- **WHEN** 渲染输出与上一帧完全相同
- **THEN** 系统跳过写入 stdout

#### Scenario: 输出部分变化
- **WHEN** 渲染输出部分变化
- **THEN** 系统只写入变化的行

### Requirement: 渲染器支持备用屏幕缓冲
渲染器 SHALL 使用 ANSI `\x1b[?1049h` 启用备用屏幕，`\x1b[?1049l` 恢复原屏幕。

#### Scenario: 启用备用屏幕
- **WHEN** 应用启动且 AlternateScreen 为 true
- **THEN** 系统写入 ANSI `\x1b[?1049h`

#### Scenario: 禁用备用屏幕
- **WHEN** 应用卸载且使用了备用屏幕
- **THEN** 系统写入 ANSI `\x1b[?1049l`

#### Scenario: TUI 不影响原终端
- **WHEN** 应用启用 AlternateScreen 选项
- **THEN** TUI 在独立屏幕渲染，退出后恢复原终端内容

### Requirement: 渲染器自动检测 TTY
渲染器 SHALL 自动检测 stdout 是否为终端，非终端时禁用交互模式。

#### Scenario: 终端环境
- **WHEN** stdout 是 TTY
- **THEN** 启用交互模式和增量渲染

#### Scenario: 管道/重定向
- **WHEN** stdout 不是 TTY (管道或重定向)
- **THEN** 禁用交互模式，直接输出

#### Scenario: CI 环境
- **WHEN** 在 CI 环境运行 (stdout 不是 TTY)
- **THEN** 禁用交互模式

### Requirement: 渲染器处理信号
渲染器 SHALL 内置 SIGINT/SIGTERM 处理，支持 Ctrl+C 可靠退出。

#### Scenario: Ctrl+C 退出
- **WHEN** 用户按下 Ctrl+C 且 ExitOnCtrlC 为 true
- **THEN** 调用 Unmount，清理资源，退出

#### Scenario: SIGTERM 退出
- **WHEN** 进程收到 SIGTERM 信号
- **THEN** 调用 Unmount，清理资源，退出

### Requirement: 渲染器隐藏/显示光标
渲染器 SHALL 在交互模式下隐藏光标，退出时恢复。

#### Scenario: 隐藏光标
- **WHEN** 应用启动且 Interactive 为 true
- **THEN** 系统写入 ANSI `\x1b[?25l`

#### Scenario: 显示光标
- **WHEN** 应用卸载
- **THEN** 系统写入 ANSI `\x1b[?25h`

### Requirement: 渲染器支持动画系统
渲染器 SHALL 提供 AnimationManager 管理动画订阅和帧推进。

#### Scenario: 动画订阅
- **WHEN** 组件调用 UseAnimation
- **THEN** 系统在 AnimationManager 中注册订阅者

#### Scenario: 动画帧推进
- **WHEN** 渲染循环 tick 且有活跃动画
- **THEN** 系统推进动画帧并触发重新渲染

#### Scenario: 动画 FPS 控制
- **WHEN** 动画设置为 10 FPS
- **THEN** 系统每秒推进约 10 帧

### Requirement: 渲染器提供 Render 函数
渲染器 SHALL 提供 `Render(element, options)` 函数创建应用实例。

#### Scenario: 基本渲染
- **WHEN** 用户调用 `ink.Render(<App />, nil)`
- **THEN** 系统返回 Instance，包含 Unmount, WaitUntilExit, Rerender, Cleanup 方法

#### Scenario: 带选项渲染
- **WHEN** 用户调用 `ink.Render(<App />, &ink.RenderOptions{Stdout: customWriter})`
- **THEN** 系统使用提供的选项进行渲染

### Requirement: 渲染器启动渲染循环
渲染器 SHALL 在独立 goroutine 中启动渲染循环。

#### Scenario: 渲染循环自动启动
- **WHEN** 用户调用 `ink.Render(<App />, nil)`
- **THEN** 系统启动 goroutine 持续渲染应用

#### Scenario: FPS 节流
- **WHEN** 用户未指定 MaxFps
- **THEN** 系统节流渲染到最大 30 FPS

### Requirement: 渲染器提供 Instance 方法
渲染器 SHALL 在 Instance 上提供 Unmount, WaitUntilExit, Rerender, Cleanup, Clear 方法。

#### Scenario: Unmount
- **WHEN** 用户调用 `instance.Unmount()`
- **THEN** 系统停止渲染循环并恢复终端状态

#### Scenario: WaitUntilExit
- **WHEN** 用户调用 `instance.WaitUntilExit()`
- **THEN** 系统阻塞直到 Unmount 被调用

#### Scenario: Rerender
- **WHEN** 用户调用 `instance.Rerender(<NewApp />)`
- **THEN** 系统更新渲染输出

#### Scenario: Cleanup
- **WHEN** 用户调用 `instance.Cleanup()`
- **THEN** 系统释放所有资源

#### Scenario: Clear
- **WHEN** 用户调用 `instance.Clear()`
- **THEN** 系统清除终端输出

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Renderer 架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: 布局树 (来自 Layout Engine)                           │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │ Layout Tree │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Build      │  构建输出缓冲                             │
│   │  Buffer     │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐    ┌─────────────┐                      │
│   │  Previous   │───▶│    Diff     │                      │
│   │  Buffer     │    │  (差异计算)  │                      │
│   └─────────────┘    └─────────────┘                      │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Patches    │                     │
│                        │  (补丁列表)  │                     │
│                        └─────────────┘                     │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Generate   │                     │
│                        │  ANSI       │                     │
│                        └─────────────┘                     │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Output     │                     │
│                        │  to Stdout  │                     │
│                        └─────────────┘                     │
│                                                             │
│   核心: 最小化输出，只更新变化的部分                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Output Buffer

```go
// 输出缓冲 (二维字符网格)
type OutputBuffer struct {
    Width  int
    Height int
    Lines  [][]Cell
    CursorX int
    CursorY int
}

// 单元格
type Cell struct {
    Char    rune
    Style   CellStyle
    ZIndex  int
    Dirty   bool
}

// 单元格样式
type CellStyle struct {
    Foreground    Color
    Background    Color
    FgBold        bool
    FgDim         bool
    Underline     bool
    Strikethrough bool
    Inverse       bool
}
```

## LogUpdate (增量渲染)

```go
type LogUpdate struct {
    stdout          io.Writer
    previousOutput  string
    previousLines   []string
    initialized     bool
    incremental     bool // 增量渲染模式
    alternateScreen bool // 备用屏幕缓冲模式
}

// Initialize 初始化终端
func (l *LogUpdate) Initialize() error {
    if l.initialized {
        return nil
    }

    // 启用备用屏幕缓冲
    if l.alternateScreen {
        l.stdout.Write([]byte(ANSIEnableAlternateScreen))
    }

    // 隐藏光标
    l.stdout.Write([]byte(ANSIHideCursor))
    // 移动光标到起始位置
    l.stdout.Write([]byte(ANSIMoveCursorHome))

    l.initialized = true
    return nil
}

// Done 完成输出，恢复终端状态
func (l *LogUpdate) Done() {
    if !l.initialized {
        return
    }

    // 显示光标
    l.stdout.Write([]byte(ANSIShowCursor))

    // 禁用备用屏幕缓冲
    if l.alternateScreen {
        l.stdout.Write([]byte(ANSIDisableAlternateScreen))
    }

    l.previousOutput = ""
    l.previousLines = make([]string, 0)
    l.initialized = false
}
```

## TTY Detection

```go
// IsTerminal 检查文件描述符是否为终端 (TTY)
func IsTerminal(file *os.File) bool {
    if file == nil {
        return false
    }

    fi, err := file.Stat()
    if err != nil {
        return false
    }

    return (fi.Mode() & os.ModeCharDevice) != 0
}

// IsStdoutTerminal 检查 stdout 是否为终端
func IsStdoutTerminal() bool {
    return IsTerminal(os.Stdout)
}
```

### TTY Behavior

| 环境 | TTY 检测 | 交互模式 | 渲染模式 |
|-----|---------|---------|---------|
| 终端 | true | 启用 | 增量渲染 |
| 管道 `|` | false | 禁用 | 直接输出 |
| 重定向 `>` | false | 禁用 | 直接输出 |
| CI 环境 | false | 禁用 | 直接输出 |

## Signal Handling

```go
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
```

## Animation System

```go
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

// Tick 推进动画帧
func (m *AnimationManager) Tick() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for _, sub := range m.subscribers {
        if sub.IsPlaying {
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
```

## Render Loop

```go
// StartRenderLoop 启动渲染循环
func (ink *Ink) StartRenderLoop() {
    ink.mu.Lock()
    if ink.isRunning {
        ink.mu.Unlock()
        return
    }
    ink.isRunning = true
    ink.mu.Unlock()

    // 初始化终端
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
            ink.onRender()

        case <-throttle.C:
            if ink.animationManager.HasActiveAnimations() {
                ink.animationManager.Tick()
                ink.scheduleRender()
            }

        case result := <-ink.exitChan:
            ink.finishUnmount(result)
            return

        case <-ink.doneChan:
            return
        }
    }
}
```

## ANSI Sequences

```go
const (
    // 光标控制
    ANSIHideCursor      = "\x1b[?25l"
    ANSIShowCursor      = "\x1b[?25h"
    ANSIClearScreen     = "\x1b[2J"
    ANSIMoveCursorHome  = "\x1b[H"
    ANSIClearLine       = "\x1b[2K"
    ANSIMoveCursorStart = "\x1b[0G"

    // 屏幕缓冲
    ANSIEnableAlternateScreen  = "\x1b[?1049h"
    ANSIDisableAlternateScreen = "\x1b[?1049l"

    // 样式
    ANSIReset = "\x1b[0m"
    ANSIBold  = "\x1b[1m"
    ANSIDim   = "\x1b[2m"
)
```

## RenderOptions

```go
type RenderOptions struct {
    Stdout              io.Writer
    Stderr              io.Writer
    ExitOnCtrlC         bool
    MaxFps              int
    PatchConsole        bool
    Interactive         bool
    AlternateScreen     bool
    IncrementalRendering bool
}
```

## Module Structure

```
ink/
├── ink.go              # Ink 实例、渲染循环
├── log_update.go       # 增量渲染、备用屏幕
├── render.go           # Render 函数、Instance
├── options.go          # RenderOptions
├── tty.go              # TTY 检测
└── register.go         # 接口注册
```

## React Ink 对齐

| 特性 | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| 增量渲染 | ✅ log-update | ✅ LogUpdate | 双缓冲 Diff |
| 备用屏幕缓冲 | ✅ | ✅ | ANSI `\x1b[?1049h/l` |
| TTY 检测 | ✅ | ✅ | `isatty` 检测 |
| 信号处理 | ✅ | ✅ | SIGINT/SIGTERM |
| 光标隐藏 | ✅ | ✅ | ANSI `\x1b[?25l/h` |
| FPS 节流 | ✅ | ✅ | time.Ticker |
| 动画系统 | ✅ useAnimation | ✅ UseAnimation | AnimationManager |
| 交互模式自动检测 | ✅ | ✅ | stdout TTY 检测 |
| Instance API | ✅ | ✅ | Unmount/WaitUntilExit/Rerender |
| ExitOnCtrlC | ✅ | ✅ | 默认启用 |