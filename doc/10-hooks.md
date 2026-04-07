# React Hooks 设计

> 在 Go 中实现 React Hooks 模式

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Hooks 架构                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   调用时机:                                                  │
│                                                             │
│   组件渲染 ──▶ beginWork() ──▶ 执行组件函数                   │
│         │                    │                              │
│         │                    ▼                              │
│         │              ┌─────────────┐                     │
│         │              │  Hooks      │                     │
│         │              │  Context    │                     │
│         │              └─────────────┘                     │
│         │                    │                              │
│         │                    ├── useState()                │
│         │                    ├── useEffect()               │
│         │                    ├── useInput()                │
│         │                    ├── useApp()                  │
│         │                    └── ...                       │
│         │                    │                              │
│         │                    ▼                              │
│         │              ┌─────────────┐                     │
│         │              │  Hook       │                     │
│         │              │  Queue      │                     │
│         │              └─────────────┘                     │
│         │                                                   │
│         └──▶ finishWork() ──▶ 处理 Effects                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Hooks 上下文

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

---

## useState

```go
// State Hook
type StateHook struct {
    value    any
    dispatch func(any)
    
    // 计算新值 (用于函数式更新)
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

// 调度状态更新
func (ctx *HooksContext) scheduleStateUpdate(fiber *Fiber, hookIndex int, newValue any) {
    update := StateUpdate{
        fiber:     fiber,
        hookIndex: hookIndex,
        value:     newValue,
    }
    
    ctx.pendingUpdates = append(ctx.pendingUpdates, update)
    
    // 触发重新渲染
    scheduleUpdateOnFiber(fiber)
}

// 应用状态更新
func applyStateUpdates(ctx *HooksContext) {
    for _, update := range ctx.pendingUpdates {
        hook := ctx.hooks[update.hookIndex].(*StateHook)
        
        if hook.hasUpdater {
            hook.value = hook.updater(hook.value)
        } else {
            hook.value = update.value
        }
    }
    
    ctx.pendingUpdates = nil
}
```

### useState 使用示例

```go
func Counter() Element {
    count, setCount := useState(0)
    
    return Box(Props{},
        Text(Props{}, fmt.Sprintf("Count: %d", count)),
        Button(Props{"onClick": func() { setCount(count + 1) }},
            Text(Props{}, "Increment"),
        ),
    )
}
```

---

## useEffect

```go
// Effect Hook
type EffectHook struct {
    create  func() func()  // 返回 cleanup 函数
    cleanup func()
    deps    []any
    
    // 标记
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
    hook := ctx.hooks[ctx.hookIndex].(*StateHook)
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

// 检查依赖变化
func depsChanged(oldDeps, newDeps []any) bool {
    if len(oldDeps) != len(newDeps) {
        return true
    }
    
    for i := range oldDeps {
        if oldDeps[i] != newDeps[i] {
            return true
        }
    }
    
    return false
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

### useEffect 使用示例

```go
func Timer() Element {
    count, setCount := useState(0)
    
    useEffect(func() func() {
        ticker := time.NewTicker(1 * time.Second)
        go func() {
            for range ticker.C {
                setCount(count + 1)
            }
        }()
        
        return func() {
            ticker.Stop()
        }
    }, nil)  // 空依赖: 只在挂载时执行
    
    return Text(Props{}, fmt.Sprintf("Time: %d", count))
}
```

---

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

### useInput 使用示例

```go
func KeyboardDemo() Element {
    lastKey, setLastKey := useState("")
    
    useInput(func(ev InputEvent) {
        if ev.Ctrl && ev.Key == KeyC {
            // Ctrl+C 退出
            os.Exit(0)
        }
        setLastKey(string(ev.Key))
    })
    
    return Text(Props{}, fmt.Sprintf("Last key: %s", lastKey))
}
```

---

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
    // 退出
    Exit()
    ExitWithError(error)
    
    // 尺寸
    Size() (width, height int)
    
    // 重渲染
    ForceUpdate()
}
```

### useApp 使用示例

```go
func AppDemo() Element {
    app := useApp()
    width, height := app.Size()
    
    useInput(func(ev InputEvent) {
        if ev.Key == KeyEscape {
            app.Exit()
        }
    })
    
    return Text(Props{}, fmt.Sprintf("Size: %dx%d", width, height))
}
```

---

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

### useMemo 使用示例

```go
func ExpensiveComponent() Element {
    items, _ := useState([]Item{...})
    
    // 只在 items 变化时重新计算
    total := useMemo(func() int {
        result := 0
        for _, item := range items {
            result += item.Value
        }
        return result
    }, []any{items})
    
    return Text(Props{}, fmt.Sprintf("Total: %d", total))
}
```

---

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

### useRef 使用示例

```go
func InputWithRef() Element {
    inputRef := useRef("")
    
    useInput(func(ev InputEvent) {
        if ev.Key == KeyEnter {
            // 访问 ref 值
            fmt.Println("Input:", inputRef.current)
        }
    })
    
    return TextInput(Props{
        "ref": inputRef,
    })
}
```

---

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
        // 通知焦点管理器
        getFocusManager().SetFocus(ctx.fiber)
    }
    
    blur = func() {
        setIsFocused(false)
        getFocusManager().ClearFocus(ctx.fiber)
    }
    
    return isFocused, focus, blur
}
```

### useFocus 使用示例

```go
func FocusableItem() Element {
    isFocused, focus, _ := useFocus()
    
    color := "white"
    if isFocused {
        color = "green"
    }
    
    useInput(func(ev InputEvent) {
        if ev.Key == KeyTab {
            // Tab 切换焦点
            getFocusManager().Next()
        }
    })
    
    return Box(Props{"onFocus": focus},
        Text(Props{"color": color}, "Item"),
    )
}
```

---

## useContext

```go
// Context Hook - 上下文访问
type ContextHook struct {
    value  any
    context *Context
}

// Context 定义
type Context struct {
    name     string
    defaultValue any
    Provider func(any, Element) Element
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

// 查找 Context 值
func findContextValue(fiber *Fiber, ctx *Context) any {
    for f := fiber; f != nil; f = f.Return {
        if f.Contexts != nil {
            if v, ok := f.Contexts[ctx]; ok {
                return v
            }
        }
    }
    return nil
}
```

### useContext 使用示例

```go
// 创建 Context
ThemeContext := createContext("theme", Theme{Color: "white"})

// Provider
func App() Element {
    return ThemeContext.Provider(Theme{Color: "green"},
        Child(),
    )
}

// Consumer
func Child() Element {
    theme := useContext(ThemeContext)
    return Text(Props{"color": theme.Color}, "Themed Text")
}
```

---

## 模块结构

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
├── context_hook.go   # useContext
├── types.go          # Hook 类型定义
└── utils.go          # 辅助函数
```

---

## 总结

### 核心要点

1. **调用规则**
   - 只在组件函数顶层调用
   - 不能在条件/循环中调用
   - 顺序必须保持一致

2. **状态管理**
   - useState: 状态存储
   - useRef: 可变引用
   - useMemo: 计算缓存

3. **副作用**
   - useEffect: 副作用处理
   - 依赖数组控制执行时机
   - cleanup 函数清理资源

4. **应用集成**
   - useInput: 键盘输入
   - useApp: 应用实例
   - useFocus: 焦点管理

5. **上下文**
   - createContext: 创建上下文
   - useContext: 消费上下文
   - Provider: 提供上下文

6. **Go 适配**
   - 泛型支持 (useState[T])
   - 闭包实现 dispatch
   - Fiber 关联存储
