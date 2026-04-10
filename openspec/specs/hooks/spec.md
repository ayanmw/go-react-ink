# hooks Specification

## Purpose
定义 Go-Ink React Hooks 的设计和实现，包括 useState、useEffect、useInput、useApp、useFocus、useAnimation 等。

## Requirements

### Requirement: useState 管理状态
useState SHALL 返回当前值和 setter 函数。

#### Scenario: 初始值
- **WHEN** 组件首次调用 `useState(0)`
- **THEN** 返回 (0, setter)

#### Scenario: 更新值
- **WHEN** 调用 `setter(1)`
- **THEN** 下次渲染时返回 (1, setter)

#### Scenario: 函数式更新
- **WHEN** 调用 `setter(prev => prev + 1)`
- **THEN** 使用前值计算新值

### Requirement: useEffect 处理副作用
useEffect SHALL 在渲染后执行副作用函数。

#### Scenario: 无依赖
- **WHEN** 使用 `useEffect(fn, nil)`
- **THEN** 每次渲染后执行

#### Scenario: 空依赖
- **WHEN** 使用 `useEffect(fn, [])`
- **THEN** 仅挂载时执行一次

#### Scenario: 依赖变化
- **WHEN** 使用 `useEffect(fn, [count])` 且 count 变化
- **THEN** 重新执行副作用

#### Scenario: 清理函数
- **WHEN** useEffect 返回清理函数
- **THEN** 下次执行前或卸载时调用清理

### Requirement: useInput 处理键盘输入
useInput SHALL 注册键盘输入处理器。

#### Scenario: 按键处理
- **WHEN** 用户按下键
- **THEN** 调用注册的处理器

#### Scenario: Ctrl+C
- **WHEN** 用户按下 Ctrl+C
- **THEN** 可在处理器中检测并退出

### Requirement: useApp 访问应用实例
useApp SHALL 返回应用控制器。

#### Scenario: Exit
- **WHEN** 调用 `app.Exit()`
- **THEN** 应用退出

#### Scenario: ExitWithError
- **WHEN** 调用 `app.ExitWithError(err)`
- **THEN** 应用带错误退出

### Requirement: useFocus 管理焦点
useFocus SHALL 返回焦点状态和控制方法。

#### Scenario: isFocused
- **WHEN** 组件获得焦点
- **THEN** isFocused 为 true

#### Scenario: focus/blur
- **WHEN** 调用 `focus()`
- **THEN** 组件获得焦点

### Requirement: useAnimation 管理动画
useAnimation SHALL 返回动画控制器。

#### Scenario: Play/Pause
- **WHEN** 调用 `anim.Play()`
- **THEN** 动画开始

#### Scenario: Frame
- **WHEN** 调用 `anim.Frame()`
- **THEN** 返回当前帧号

### Requirement: useMemo 缓存计算
useMemo SHALL 缓存计算结果，依赖不变时返回缓存值。

#### Scenario: 首次计算
- **WHEN** 首次调用 `useMemo(fn, deps)`
- **THEN** 执行 fn 并缓存结果

#### Scenario: 依赖不变
- **WHEN** 依赖未变化
- **THEN** 返回缓存值

### Requirement: useRef 创建可变引用
useRef SHALL 创建可变引用，修改不触发重渲染。

#### Scenario: 访问 current
- **WHEN** 使用 `ref.current`
- **THEN** 返回当前值

### Requirement: useContext 访问上下文
useContext SHALL 访问最近的 Provider 提供的值。

#### Scenario: Provider 值
- **WHEN** 祖先组件提供 Context.Provider
- **THEN** 返回 Provider 的值

#### Scenario: 默认值
- **WHEN** 无 Provider
- **THEN** 返回 Context 默认值

## Hooks Context

```go
// Hooks 上下文 (每个组件实例一个)
type HooksContext struct {
    // 当前 Fiber
    fiber *Fiber
    
    // Hook 索引
    hookIndex int
    
    // Hook 列表
    hooks []Hook
    
    // 状态更新队列
    pendingUpdates []StateUpdate
    
    // Effect 队列
    effects []EffectHook
}

// Hook 接口
type Hook interface {
    hookType() HookType
}

type HookType int

const (
    HookState HookType = iota
    HookEffect
    HookMemo
    HookRef
    HookContext
)
```

## useState

```go
// State Hook
type StateHook struct {
    value    any
    dispatch func(any)
    
    hasUpdater bool
    updater    func(any) any
}

// useState 实现
func useState[T any](initial T) (T, func(T)) {
    ctx := getCurrentHooksContext()
    
    // 首次渲染: 创建 Hook
    if ctx.hookIndex >= len(ctx.hooks) {
        hook := &StateHook{
            value: initial,
        }
        hook.dispatch = func(newValue T) {
            ctx.scheduleStateUpdate(ctx.fiber, ctx.hookIndex, newValue)
        }
        ctx.hooks = append(ctx.hooks, hook)
    }
    
    // 获取 Hook
    hook := ctx.hooks[ctx.hookIndex].(*StateHook)
    ctx.hookIndex++
    
    return hook.value.(T), hook.dispatch.(func(T))
}
```

## useEffect

```go
// Effect Hook
type EffectHook struct {
    create  func() func()  // 返回 cleanup 函数
    cleanup func()
    deps    []any
    
    hasChanged bool
}

// useEffect 实现
func useEffect(create func() func(), deps []any) {
    ctx := getCurrentHooksContext()
    
    // 首次渲染: 创建 Hook
    if ctx.hookIndex >= len(ctx.hooks) {
        hook := &EffectHook{
            create: create,
            deps:   deps,
        }
        ctx.hooks = append(ctx.hooks, hook)
        ctx.effects = append(ctx.effects, hook)
        ctx.hookIndex++
        return
    }
    
    // 获取现有 Hook
    hook := ctx.hooks[ctx.hookIndex].(*EffectHook)
    ctx.hookIndex++
    
    // 检查依赖变化
    hasChanged := depsChanged(hook.deps, deps)
    
    if hasChanged {
        hook.create = create
        hook.deps = deps
        hook.hasChanged = true
        ctx.effects = append(ctx.effects, hook)
    }
}

// 执行 Effects (提交阶段)
func flushEffects(effects []EffectHook) {
    for _, effect := range effects {
        // 执行旧 cleanup
        if effect.cleanup != nil {
            effect.cleanup()
        }
        
        // 执行新 create
        effect.cleanup = effect.create()
    }
}
```

## useInput

```go
// Input Hook - 键盘输入
type InputHook struct {
    handler func(InputEvent)
    active  bool
}

// useInput 实现
func useInput(handler func(InputEvent)) {
    ctx := getCurrentHooksContext()
    
    useEffect(func() func() {
        // 注册输入处理器
        app := getAppContext()
        app.RegisterInputHandler(handler)
        
        return func() {
            app.UnregisterInputHandler(handler)
        }
    }, []any{handler})
}

// InputEvent
type InputEvent struct {
    Key   Key
    Rune  rune
    Ctrl  bool
    Alt   bool
    Shift bool
}
```

## useApp

```go
// App Hook - 应用实例访问
type AppHook struct {
    app *App
}

// useApp 实现
func useApp() *App {
    ctx := getCurrentHooksContext()
    return ctx.fiber.Root.StateNode.(*App)
}

// App 接口
type App interface {
    Exit()
    ExitWithError(error)
    Size() (width, height int)
    ForceUpdate()
}
```

## useFocus

```go
// Focus Hook - 焦点管理
type FocusHook struct {
    isFocused bool
    focus     func()
    blur      func()
}

// useFocus 实现
func useFocus() (isFocused bool, focus func(), blur func()) {
    ctx := getCurrentHooksContext()
    
    isFocused, setIsFocused := useState(false)
    
    focus = func() {
        setIsFocused(true)
        getFocusManager().SetFocus(ctx.fiber)
    }
    
    blur = func() {
        setIsFocused(false)
        getFocusManager().ClearFocus(ctx.fiber)
    }
    
    return isFocused, focus, blur
}
```

## useAnimation

```go
// Animation Hook
type AnimationHook struct {
    frame    int
    fps      int
    isPlaying bool
    subscriber *AnimationSubscriber
}

// useAnimation 实现
func useAnimation(fps int) *AnimationController {
    ctx := getCurrentHooksContext()
    
    frame, setFrame := useState(0)
    isPlaying, setIsPlaying := useState(false)
    
    controller := &AnimationController{
        frame:    frame,
        fps:      fps,
        isPlaying: isPlaying,
    }
    
    useEffect(func() func() {
        subscriber := &AnimationSubscriber{
            Fps:       fps,
            Callback:  func(f int) { setFrame(f) },
            IsPlaying: isPlaying,
        }
        
        getAnimationManager().Subscribe(subscriber)
        
        return func() {
            getAnimationManager().Unsubscribe(subscriber)
        }
    }, []any{fps, isPlaying})
    
    return controller
}

// AnimationController
type AnimationController struct {
    frame    int
    fps      int
    isPlaying bool
}

func (c *AnimationController) Play() {
    c.isPlaying = true
}

func (c *AnimationController) Pause() {
    c.isPlaying = false
}

func (c *AnimationController) Frame() int {
    return c.frame
}
```

## useMemo

```go
// Memo Hook - 计算缓存
type MemoHook struct {
    value  any
    deps   []any
    create func() any
}

// useMemo 实现
func useMemo[T any](create func() T, deps []any) T {
    ctx := getCurrentHooksContext()
    
    // 首次渲染
    if ctx.hookIndex >= len(ctx.hooks) {
        hook := &MemoHook{
            value:  create(),
            deps:   deps,
            create: create,
        }
        ctx.hooks = append(ctx.hooks, hook)
        ctx.hookIndex++
        return hook.value.(T)
    }
    
    // 获取现有 Hook
    hook := ctx.hooks[ctx.hookIndex].(*MemoHook)
    ctx.hookIndex++
    
    // 检查依赖变化
    if depsChanged(hook.deps, deps) {
        hook.value = create()
        hook.deps = deps
    }
    
    return hook.value.(T)
}
```

## useRef

```go
// Ref Hook - 可变引用
type RefHook struct {
    current any
}

// useRef 实现
func useRef[T any](initial T) *RefHook {
    ctx := getCurrentHooksContext()
    
    // 首次渲染
    if ctx.hookIndex >= len(ctx.hooks) {
        hook := &RefHook{
            current: initial,
        }
        ctx.hooks = append(ctx.hooks, hook)
        ctx.hookIndex++
        return hook
    }
    
    // 获取现有 Hook
    hook := ctx.hooks[ctx.hookIndex].(*RefHook)
    ctx.hookIndex++
    
    return hook
}
```

## useContext

```go
// Context Hook - 上下文访问
type ContextHook struct {
    value   any
    context *Context
}

// Context 定义
type Context struct {
    name         string
    defaultValue any
    Provider     func(any, Element) Element
}

// createContext
func createContext[T any](name string, defaultValue T) *Context {
    return &Context{
        name:         name,
        defaultValue: defaultValue,
    }
}

// useContext 实现
func useContext[T any](ctx *Context) T {
    hooksCtx := getCurrentHooksContext()
    
    // 向上查找 Provider
    value := findContextValue(hooksCtx.fiber, ctx)
    if value == nil {
        value = ctx.defaultValue
    }
    
    return value.(T)
}
```

## Module Structure

```
hooks/
├── context.go        # Hooks 上下文
├── state.go          # useState
├── effect.go         # useEffect
├── input.go          # useInput
├── app.go            # useApp
├── memo.go           # useMemo
├── ref.go            # useRef
├── focus.go          # useFocus
├── animation.go      # useAnimation
├── context_hook.go   # useContext
├── types.go          # Hook 类型定义
└── utils.go          # 辅助函数
```

## React Ink 对齐

| Hook | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| useState | ✅ | ✅ | 状态管理 |
| useEffect | ✅ | ✅ | 副作用 |
| useInput | ✅ | ✅ | 键盘输入 |
| useApp | ✅ | ✅ | 应用实例 |
| useFocus | ✅ | ✅ | 焦点管理 |
| useAnimation | ✅ | ✅ | 动画控制 |
| useMemo | ✅ | ✅ | 计算缓存 |
| useRef | ✅ | ✅ | 可变引用 |
| useContext | ✅ | ✅ | 上下文访问 |
| useStdin | ✅ | ✅ | stdin 访问 |
| useStdout | ✅ | ✅ | stdout 访问 |
| useWindowSize | ✅ | ✅ | 窗口尺寸 |

### Hook 规则对齐

| 规则 | React Ink | Go-Ink |
|-----|-----------|--------|
| 顶层调用 | ✅ | ✅ |
| 条件中禁止 | ✅ | ✅ |
| 循环中禁止 | ✅ | ✅ |
| 顺序一致 | ✅ | ✅ |