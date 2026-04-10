# Incremental Rendering Integration 设计

## 问题陈述

当前 Go React-Ink 的增量渲染集成存在问题：
1. `RenderOptions` 缺少 `IncrementalRendering` 字段
2. `NewLogUpdate` 不接受 `incremental` 参数
3. 虽然 `LogUpdate` 默认启用 incremental 模式，但用户不可配置
4. animation-demo 重复实现了incremental 渲染器

## 设计目标

参考 React Ink 的架构：
- 在 `RenderOptions` 中添加 `IncrementalRendering` 选项（默认 `true`）
- `NewLogUpdate` 接受 `incremental` 参数
- `Ink` 实例根据 `RenderOptions.IncrementalRendering` 配置 `LogUpdate`
- 简化 animation-demo 使用基础库的功能

## 架构修改

### 1. RenderOptions 添加字段

```go
type RenderOptions struct {
    // ... 现有字段

    // IncrementalRendering 是否使用增量渲染（默认: true）
    // true: 只更新变化的行，减少闪烁
    // false: 全屏刷新模式
    IncrementalRendering bool
}
```

### 2. NewLogUpdate 修改签名

```go
// NewLogUpdate 创建新的输出管理器
func NewLogUpdate(stdout io.Writer, incremental bool) *LogUpdate {
    return &LogUpdate{
        stdout:      stdout,
        previousLines: make([]string, 0),
        incremental: incremental, // 使用参数而非硬编码 true
    }
}
```

### 3. NewInk 配置 LogUpdate

```go
func NewInk(opts *RenderOptions) *Ink {
    opts = applyDefaults(opts)

    return &Ink{
        options:          opts,
        scheduler:        fiber.NewScheduler(),
        log:              NewLogUpdate(opts.Stdout, opts.IncrementalRendering), // 传递选项
        // ...
    }
}
```

### 4. 默认值设置

```go
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
    // IncrementalRendering 默认为 true
    if opts.IncrementalRendering == false {
        // 保持 false（如果用户显式设置为 false）
    } else {
        // 隐式保留 true
    }

    return opts
}

// DefaultRenderOptions 更新
func DefaultRenderOptions() *RenderOptions {
    return &RenderOptions{
        Stdout:               os.Stdout,
        Stderr:               os.Stderr,
        ExitOnCtrlC:          true,
        MaxFps:               30,
        PatchConsole:         false,
        Interactive:          true,
        AlternateScreen:      false,
        IncrementalRendering: true,  // 新增默认值 true
    }
}
```

### 5. 简化 animation-demo

移除 `IncrementalRenderer`，直接使用 `Ink` 的功能：

```go
// 之前：需要自己实现增量渲染器
func main() {
    renderer := NewIncrementalRenderer()  // ❌ 冗余

    for {
        lines := renderFrame(frame, progress)
        renderer.Render(lines)  // ❌ 自己管理渲染
    }
}

// 之后：使用 Ink 实例
func main() {
    opts := &ink.RenderOptions{
        IncrementalRendering: true,  // ✅ 配置选项
    }

    instance := ink.Render(Appcomponent{}, opts)
    // Ink 内部处理所有渲染逻辑
}
```

## 向后兼容性

- `IncrementalRendering` 默认值为 `true`，保持现有行为
- `NewLogUpdate` 修改签名需要更新所有调用点
- 简单的 migration：所有 `NewLogUpdate(stdout)` 改为 `NewLogUpdate(stdout, true)`

## 实现步骤

1. ✅ 修改 `RenderOptions` 添加 `IncrementalRendering` 字段
2. ✅ 修改 `NewLogUpdate` 签名接受 `incremental` 参数
3. ✅ 修改 `NewInk` 传递 `IncrementalRendering` 到 `NewLogUpdate`
4. ✅ 更新 `DefaultRenderOptions()`
5. ✅ 更新测试中的 `NewLogUpdate` 调用
6. ✅ 简化 animation-demo 使用 Ink 基础功能（可选， 示例简化）

## 替代方案考虑

### 方案 A：保持 NewLogUpdate 签名不变
改为 `SetIncremental` setter：
```go
log := NewLogUpdate(stdout)
log.SetIncremental(opts.IncrementalRendering)
```
**优点**：不破坏现有签名
**缺点**：多一行代码，不够直截了当

### 方案 B：禁用 writeStandard 模式
保持 `incremental=true` 硬编码，但移除 `writeStandard` 分支
**优点**：最简化
**缺点**：失去用户选择全屏刷新模式的能力

### 方案 C（选择）：修改 NewLogUpdate 签名
与 React Ink 一致，更清晰
**优点**：架构清晰，用户明确控制，最易维护

## 参考 React Ink

```typescript
// React Ink 的 create 函数
const logUpdate = {
  create: (stream, {showCursor = false, incremental = false} = {}): LogUpdate => {
    if (incremental) {
      return createIncremental(stream, {showCursor});
    }
    return createStandard(stream, {showCursor});
  }
}

// Ink 构造时使用
this.log = logUpdate.create(options.stdout, {
  incremental: options.incrementalRendering,
});
```

## 影响

- **用户代码**：只需在 `RenderOptions` 中配置 `IncrementalRendering`
- **内部代码**：最小改动，主要是参数传递
- **性能**：无变化，逻辑相同
- **维护性**：显著提高，消除重复代码

## 后续扩展

此改动为以下功能铺路：
- `cursorPosition` 支持与 `buildCursorSuffix` 辅助函数
- `onRender` 回调暴露给用户
- 允许运行时切换模式（如需）
