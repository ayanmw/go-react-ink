# API 文档

## 核心包

### pkg/core

核心类型和函数。

#### Element

元素接口:

```go
type Element interface {
    Render() string
}
```

#### Props

属性映射:

```go
type Props map[string]any
```

#### CreateElement

创建元素 (类似 React.createElement):

```go
func CreateElement(component Component, props Props, children []Element) Element
```

---

## 组件包

### pkg/components

内置组件。

#### Box

Flexbox 容器组件。

**属性:**

| 属性 | 类型 | 说明 |
|------|------|------|
| flexDirection | string | row, column, row-reverse, column-reverse |
| justifyContent | string | flex-start, flex-end, center, space-between, space-around |
| alignItems | string | flex-start, flex-end, center, stretch, baseline |
| alignSelf | string | 覆盖父级 alignItems |
| flexWrap | string | nowrap, wrap, wrap-reverse |
| flexGrow | int | Flex 增长比例 |
| flexShrink | int | Flex 收缩比例 |
| flexBasis | string | Flex 基础值 |
| width | int | 固定宽度 |
| height | int | 固定高度 |
| minWidth | int | 最小宽度 |
| minHeight | int | 最小高度 |
| maxWidth | int | 最大宽度 |
| maxHeight | int | 最大高度 |
| padding | int | 内边距 |
| paddingTop | int | 上内边距 |
| paddingRight | int | 右内边距 |
| paddingBottom | int | 下内边距 |
| paddingLeft | int | 左内边距 |
| margin | int | 外边距 |
| marginTop | int | 上外边距 |
| marginRight | int | 右外边距 |
| marginBottom | int | 下外边距 |
| marginLeft | int | 左外边距 |
| borderStyle | string | 边框样式 |
| borderColor | string | 边框颜色 |
| display | string | flex, none |

#### Text

文本组件。

**属性:**

| 属性 | 类型 | 说明 |
|------|------|------|
| children | string | 文本内容 |
| color | string | 文本颜色 |
| background | string | 背景颜色 |
| bold | bool | 粗体 |
| italic | bool | 斜体 |
| underline | bool | 下划线 |
| dim | bool | 变暗 |
| blink | bool | 闪烁 |
| reverse | bool | 反转颜色 |
| wrap | string | truncate, wrap, wrap-middle, wrap-end |
| overflow | string | truncate, ellipsis |

#### Spacer

空白填充组件，自动填充剩余空间。

#### Newline

换行组件。

#### Static

静态内容组件，不参与重渲染。

**属性:**

| 属性 | 类型 | 说明 |
|------|------|------|
| items | []string | 静态内容项 |

#### Transform

文本转换组件。

**属性:**

| 属性 | 类型 | 说明 |
|------|------|------|
| transform | func(string) string | 转换函数 |

---

## Hooks 包

### pkg/hooks

React-like Hooks。

#### UseState

状态管理 Hook:

```go
func UseState(ctx *HookContext, initial any) (any, func(any))
```

**示例:**

```go
count, setCount := hooks.UseState(ctx, 0)
setCount(42)
```

#### UseEffect

副作用 Hook:

```go
func UseEffect(ctx *HookContext, setup func() func(), deps []any)
```

**示例:**

```go
hooks.UseEffect(ctx, func() func() {
    // 初始化逻辑
    return func() {
        // 清理逻辑
    }
}, []any{dependency})
```

#### UseLayoutEffect

布局副作用 Hook (同步执行):

```go
func UseLayoutEffect(ctx *HookContext, setup func() func(), deps []any)
```

#### UseRef

引用 Hook:

```go
func UseRef(ctx *HookContext, initial any) *RefHook
```

**示例:**

```go
ref := hooks.UseRef(ctx, 0)
ref.Current = 42
```

#### UseMemo

记忆化 Hook:

```go
func UseMemo(ctx *HookContext, factory func() any, deps []any) any
```

**示例:**

```go
result := hooks.UseMemo(ctx, func() any {
    return expensiveComputation(a, b)
}, []any{a, b})
```

#### UseCallback

回调记忆化 Hook:

```go
func UseCallback(ctx *HookContext, callback func(), deps []any) func()
```

---

## 输入包

### pkg/input

输入处理 Hooks。

#### UseInput

键盘输入 Hook:

```go
func UseInput(ctx *hooks.HookContext, handler func(Key)) *InputHook
```

**Key 结构:**

```go
type Key struct {
    Name     string  // 键名: up, down, left, right, enter, escape, etc.
    Sequence string  // 原始序列
    Shift    bool
    Ctrl     bool
    Alt      bool
    Meta     bool
}
```

**示例:**

```go
input := input.UseInput(ctx, func(key input.Key) {
    switch key.Name {
    case "up":
        // 处理上箭头
    case "down":
        // 处理下箭头
    case "enter":
        // 处理回车
    }
})
```

#### UseApp

应用控制 Hook:

```go
func UseApp(ctx *hooks.HookContext) *AppHook
```

**AppHook 方法:**

- `Exit()` - 退出应用

#### UseFocus

焦点 Hook:

```go
func UseFocus(ctx *hooks.HookContext) (isFocused bool, focus func())
```

#### UseFocusWithEvents

带事件的焦点 Hook:

```go
func UseFocusWithEvents(ctx *hooks.HookContext, onFocus func(), onBlur func()) (isFocused bool, focus func())
```

#### UseFocusManager

焦点管理器 Hook:

```go
func UseFocusManager(ctx *hooks.HookContext) *FocusManager
```

**FocusManager 方法:**

- `Focus(id string)` - 聚焦指定元素
- `Blur()` - 取消焦点
- `FocusNext()` - 聚焦下一个
- `FocusPrev()` - 聚焦上一个
- `Register(id string)` - 注册可聚焦元素
- `Unregister(id string)` - 注销可聚焦元素
- `IsFocused(id string) bool` - 检查是否聚焦

#### UseCursor

光标 Hook:

```go
func UseCursor(ctx *hooks.HookContext) *CursorHook
```

**CursorHook 方法:**

- `Show()` - 显示光标
- `Hide()` - 隐藏光标
- `SetPosition(x, y int)` - 设置位置
- `SetShape(shape string)` - 设置形状 (block, underline, bar)

#### UseAnimation

动画 Hook:

```go
func UseAnimation(ctx *hooks.HookContext, fps int) *AnimationHook
```

**AnimationHook 方法:**

- `Play()` - 播放
- `Pause()` - 暂停
- `Stop()` - 停止
- `NextFrame()` - 下一帧
- `Reset()` - 重置
- `Frame() int` - 获取当前帧
- `IsPlaying() bool` - 是否播放中

---

## 布局包

### pkg/layout

Flexbox 布局引擎。

#### Node

布局节点:

```go
node := layout.NewNode()
node.Direction = layout.DirectionRow
node.Justify = layout.JustifyCenter
node.AlignItems = layout.AlignCenter
node.Width = 100
node.Height = 20
node.FlexGrow = 1
node.Padding = [4]float64{1, 1, 1, 1} // top, right, bottom, left
```

**Direction 常量:**

- `DirectionRow`
- `DirectionRowReverse`
- `DirectionColumn`
- `DirectionColumnReverse`

**Justify 常量:**

- `JustifyFlexStart`
- `JustifyFlexEnd`
- `JustifyCenter`
- `JustifySpaceBetween`
- `JustifySpaceAround`
- `JustifySpaceEvenly`

**Align 常量:**

- `AlignAuto`
- `AlignFlexStart`
- `AlignFlexEnd`
- `AlignCenter`
- `AlignStretch`
- `AlignBaseline`

**方法:**

- `AddChild(child *Node)` - 添加子节点
- `RemoveChild(child *Node)` - 移除子节点
- `CalculateLayout(width, height float64)` - 计算布局
- `GetLayout() Layout` - 获取布局结果
- `SetMargin(top, right, bottom, left float64)` - 设置外边距
- `SetPadding(top, right, bottom, left float64)` - 设置内边距
- `SetBorder(top, right, bottom, left float64)` - 设置边框

---

## 渲染包

### pkg/renderer

终端渲染器。

#### Renderer

渲染器:

```go
r := renderer.NewRenderer(80, 24)
buf := r.GetBuffer()
buf.SetCell(0, 0, 'A', renderer.Style{FgColor: renderer.ColorGreen})
output := r.Render()
```

#### Buffer

输出缓冲:

```go
buf := renderer.NewBuffer(80, 24)
buf.SetCell(x, y int, char rune, style Style)
cell := buf.GetCell(x, y int)
changes := buf.Diff(other *Buffer) []Change
```

#### Style

样式:

```go
style := renderer.Style{
    FgColor:   renderer.ColorGreen,
    BgColor:   renderer.ColorBlack,
    Bold:      true,
    Italic:    true,
    Underline: true,
    Dim:       false,
    Blink:     false,
    Reverse:   false,
}
```

#### Color

颜色:

```go
// 预定义颜色
renderer.ColorDefault
renderer.ColorBlack
renderer.ColorRed
renderer.ColorGreen
renderer.ColorYellow
renderer.ColorBlue
renderer.ColorMagenta
renderer.ColorCyan
renderer.ColorWhite

// RGB 颜色
color := renderer.RGBColor(255, 128, 64)

// 索引颜色
color := renderer.IndexColor(42)
```

---

## Fiber 包

### pkg/fiber

Fiber 协调器 (内部 API)。

#### Fiber

Fiber 节点:

```go
fiber := fiber.NewFiber(fiber.TagHostComponent, props, "key")
fiber.AppendChild(child)
fiber.RemoveChild(child)
next := fiber.FindNextDefibr()
```

#### Scheduler

调度器:

```go
s := fiber.NewScheduler()
s.ScheduleUpdate(root)
for !s.IsWorkComplete() {
    unit := s.RequestWork()
    s.PerformUnitOfWork(unit)
}
s.CommitRoot()
```

---

## 编译器包

### internal/compiler

编译器 (内部 API)。

#### Compiler

编译器:

```go
c := compiler.New("ink")
output, err := c.Compile("app.gox", source)
err := c.CompileFile("app.gox")
err := c.CompileDir("src", "out")
```

---

## CLI 工具

### gox

命令行编译器。

**用法:**

```bash
gox [options] <input.gox...>
gox [options] <input-dir>
```

**选项:**

| 选项 | 说明 |
|------|------|
| -o | 输出文件 |
| -pkg | 组件包名 (默认: ink) |
| -watch | 监听模式 |
| -version | 显示版本 |
| -help | 显示帮助 |

**示例:**

```bash
# 编译单个文件
gox app.gox -o app.go

# 编译目录
gox ./src

# 使用自定义包名
gox -pkg myui app.gox
```