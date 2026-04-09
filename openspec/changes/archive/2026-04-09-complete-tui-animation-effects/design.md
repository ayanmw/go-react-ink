# Design: Complete TUI Animation Effects

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          User Application                            │
│   render(<App />) ──► {unmount, waitUntilExit, rerender, ...}       │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          pkg/ink/ink.go                              │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                      Ink Instance                              │  │
│  │  - rootNode: *dom.Element                                      │  │
│  │  - scheduler: *fiber.Scheduler                                  │  │
│  │  - renderer: *renderer.Renderer                                 │  │
│  │  - isRunning: bool                                              │  │
│  │  - exitChan: chan any                                           │  │
│  │                                                                 │  │
│  │  Methods:                                                       │  │
│  │  - Render(node)                                                 │  │
│  │  - Unmount(error?)                                              │  │
│  │  - WaitUntilExit()                                              │  │
│  │  - onRender()                                                   │  │
│  │  - startRenderLoop()                                            │  │
│  └──────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          ▼                         ▼                         ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│  pkg/fiber/      │    │  pkg/renderer/   │    │  pkg/input/      │
│  scheduler.go    │    │  renderer.go     │    │  input.go        │
│                  │    │                  │    │                  │
│  - ScheduleUpdate│    │  - Render()      │    │  - UseInput      │
│  - PerformWork   │    │  - RenderFull()  │    │  - UseAnimation  │
│  - CommitRoot    │    │  - Diff/ANSI     │    │  - RawMode       │
└──────────────────┘    └──────────────────┘    └──────────────────┘
```

## Core Components

### 1. Ink Instance (`pkg/ink/ink.go`)

The main orchestrator that manages the render loop.

```go
type Ink struct {
    // Core components
    rootNode   *dom.Element
    scheduler  *fiber.Scheduler
    renderer   *renderer.Renderer
    log        *logUpdate
    
    // State
    isRunning     bool
    isUnmounted   bool
    isUnmounting  bool
    lastOutput    string
    lastHeight    int
    
    // Channels
    exitChan      chan any
    renderChan    chan struct{}
    
    // Options
    options       *RenderOptions
}

type RenderOptions struct {
    Stdout          io.Writer
    Stderr          io.Writer
    ExitOnCtrlC     bool
    MaxFps          int     // default: 30
    PatchConsole    bool
    Interactive     bool
    AlternateScreen bool
}

type Instance struct {
    Rerender             func(core.Element)
    Unmount              func()
    WaitUntilExit        func() (any, error)
    WaitUntilRenderFlush func() error
    Cleanup              func()
    Clear                func()
}
```

### 2. Render Function (`pkg/ink/render.go`)

Entry point that creates an Ink instance and starts the render loop.

```go
func Render(element core.Element, opts *RenderOptions) *Instance {
    // Apply defaults
    opts = applyDefaults(opts)
    
    // Create instance
    ink := NewInk(opts)
    
    // Initial render
    ink.Render(element)
    
    // Start render loop
    go ink.startRenderLoop()
    
    // Handle signals
    go ink.handleSignals()
    
    return &Instance{
        Unmount:       ink.Unmount,
        WaitUntilExit: ink.WaitUntilExit,
        Rerender:      ink.Render,
        Cleanup:       ink.Cleanup,
        Clear:         ink.Clear,
    }
}
```

### 3. Render Loop

```go
func (ink *Ink) startRenderLoop() {
    // Throttle ticker
    throttle := time.NewTicker(time.Second / time.Duration(ink.options.MaxFps))
    defer throttle.Stop()
    
    for {
        select {
        case <-ink.renderChan:
            // Scheduled update
            ink.onRender()
            
        case <-throttle.C:
            // Animation tick
            if ink.hasAnimationSubscribers() {
                ink.tickAnimations()
                ink.onRender()
            }
            
        case result := <-ink.exitChan:
            // Exit requested
            ink.finishUnmount(result)
            return
            
        case <-ink.doneChan:
            return
        }
    }
}

func (ink *Ink) onRender() {
    // Calculate layout
    ink.calculateLayout()
    
    // Render to string
    output := render(ink.rootNode)
    
    // Diff with previous output
    if output == ink.lastOutput {
        return
    }
    
    // Write to stdout using ANSI escapes
    ink.writeOutput(output)
    
    ink.lastOutput = output
}
```

### 4. Animation System

Integrate `useAnimation` with the render loop.

```go
// pkg/input/animation.go
type AnimationController struct {
    frame      int
    isPlaying  bool
    fps        int
    lastTick   time.Time
    subscribers []AnimationCallback
}

type AnimationCallback func(frame int)

func (a *AnimationController) Play() {
    a.isPlaying = true
    a.lastTick = time.Now()
}

func (a *AnimationController) Pause() {
    a.isPlaying = false
}

func (a *AnimationController) NextFrame() {
    if a.isPlaying {
        a.frame++
        for _, cb := range a.subscribers {
            cb(a.frame)
        }
    }
}

func UseAnimation(ctx *hooks.HookContext, fps int) *AnimationController {
    // Hook implementation that registers with global animation manager
}
```

### 5. Exit Handling

```go
func (ink *Ink) Unmount(errorOrResult ...any) {
    if ink.isUnmounted || ink.isUnmounting {
        return
    }
    
    ink.isUnmounting = true
    
    // Final render
    ink.onRender()
    
    // Restore terminal state
    if ink.options.Interactive {
        // Show cursor
        ink.options.Stdout.Write([]byte("\x1b[?25h"))
        // Exit alternate screen if used
        if ink.options.AlternateScreen {
            ink.options.Stdout.Write([]byte("\x1b[?1049l"))
        }
    }
    
    // Signal exit
    var result any
    if len(errorOrResult) > 0 {
        result = errorOrResult[0]
    }
    ink.exitChan <- result
    ink.isUnmounted = true
}
```

## File Structure

```
pkg/
├── ink/
│   ├── ink.go           # Ink instance and render loop
│   ├── render.go        # Render() function entry point
│   ├── options.go       # RenderOptions and defaults
│   └── log_update.go    # Terminal output management (like log-update)
├── fiber/
│   ├── scheduler.go     # (existing) Fiber scheduler
│   └── fiber.go         # (existing) Fiber types
├── renderer/
│   ├── renderer.go      # (existing) ANSI rendering
│   └── buffer.go        # (existing) Output buffer
├── input/
│   ├── input.go         # (existing) useInput hook
│   └── animation.go     # Animation controller
├── hooks/
│   └── hooks.go         # (existing) Hooks implementation
└── core/
    └── element.go       # (existing) Element types
```

## API Examples

### Basic Usage

```go
package main

import (
    "github.com/ayanmw/go-react-ink/pkg/ink"
    "github.com/ayanmw/go-react-ink/pkg/core"
    "github.com/ayanmw/go-react-ink/pkg/components"
)

func App() core.Element {
    return (
        <Box flexDirection="column">
            <Text>Hello, World!</Text>
        </Box>
    )
}

func main() {
    instance := ink.Render(App(), nil)
    instance.WaitUntilExit()
}
```

### Animation Example

```go
package main

import (
    "github.com/ayanmw/go-react-ink/pkg/ink"
    "github.com/ayanmw/go-react-ink/pkg/core"
    "github.com/ayanmw/go-react-ink/pkg/components"
    "github.com/ayanmw/go-react-ink/pkg/hooks"
    "github.com/ayanmw/go-react-ink/pkg/input"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func SpinnerApp(ctx *hooks.HookContext) core.Element {
    anim := input.UseAnimation(ctx, 10) // 10 FPS
    anim.Play()
    
    frame := spinnerFrames[anim.Frame() % len(spinnerFrames)]
    
    return (
        <Box>
            <Text color="cyan">{frame}</Text>
            <Text> Loading...</Text>
        </Box>
    )
}

func main() {
    ctx := hooks.NewHookContext()
    instance := ink.Render(SpinnerApp(ctx), nil)
    instance.WaitUntilExit()
}
```

### With Exit

```go
func CounterApp(ctx *hooks.HookContext) core.Element {
    count, setCount := hooks.UseState(ctx, 0)
    app := input.UseApp(ctx)
    
    input.UseInput(ctx, func(key input.Key) {
        if key.Name == "q" || key.Name == "escape" {
            app.Exit() // Clean unmount
        }
        if key.Name == "up" {
            setCount(count.(int) + 1)
        }
    })
    
    return (
        <Box flexDirection="column">
            <Text>Count: {count}</Text>
            <Text dim>Press Up to increment, Q to exit</Text>
        </Box>
    )
}
```

## Implementation Order

1. **pkg/ink/log_update.go** - Terminal output management (clear lines, move cursor)
2. **pkg/ink/options.go** - RenderOptions and defaults
3. **pkg/ink/ink.go** - Ink instance, render loop, unmount
4. **pkg/ink/render.go** - Render() entry point
5. **pkg/input/animation.go** - Animation controller integration
6. **pkg/input/app.go** - useApp hook with Exit()
7. **examples/animation-demo/main.go** - Update with real animation

## Terminal Control Sequences

| Feature | ANSI Sequence |
|---------|---------------|
| Hide cursor | `\x1b[?25l` |
| Show cursor | `\x1b[?25h` |
| Clear screen | `\x1b[2J` |
| Move cursor home | `\x1b[H` |
| Clear line | `\x1b[2K` |
| Move cursor up N | `\x1b[N A` |
| Save cursor | `\x1b[s` |
| Restore cursor | `\x1b[u` |
| Alternate screen | `\x1b[?1049h` / `\x1b[?1049l` |