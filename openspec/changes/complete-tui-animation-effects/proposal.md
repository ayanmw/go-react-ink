# Proposal: Complete TUI Animation Effects and UI Design

## Summary

补全与 React Ink 相同的 TUI 动画效果和 UI 设计，实现完整的终端 UI 渲染循环、动画系统、以及与 React Ink 对齐的 API。

## Motivation

当前 go-react-ink 项目存在以下问题：

1. **缺少渲染循环** - 示例只是调用 `app.Render()` 打印一次，不是真正的动画
2. **缺少 render/unmount API** - React Ink 的核心 API 未实现
3. **动画系统不完整** - `useAnimation` hook 存在但未与渲染循环集成
4. **缺少终端交互** - 无 stdin 处理、光标控制、退出机制

用户反馈：
> example 中测试都是 `fmt.Println(app.Render())`，没有动画，这应该是不对的吧。比如看 animation-demo 这个 demo，只有输出，没有进度条变化。不符合动画预期。

## Goals

1. 实现 `render()` 函数 - 启动渲染循环，返回 `{unmount, waitUntilExit, rerender, waitUntilRenderFlush}`
2. 实现渲染循环 - 持续更新终端显示，支持节流
3. 集成动画系统 - `useAnimation` 与渲染循环配合
4. 实现终端控制 - 光标显示/隐藏、raw mode、退出处理
5. 对齐 React Ink API - `useApp()`, `exit()`, `unmount()`

## Non-Goals

- 不实现 React Ink 的所有内部细节（如 reconciler 的并发模式）
- 不实现测试工具（ink-testing-library）
- 不实现非交互模式的所有特性

## Success Criteria

1. `animation-demo` 示例能显示真正的动画效果（spinner 旋转、进度条变化）
2. 用户可以调用 `render()` 启动应用，调用 `unmount()` 退出
3. `useAnimation` hook 能按指定 FPS 触发更新
4. Ctrl+C 能正确退出应用

## Dependencies

- 现有的 `pkg/renderer/renderer.go` - ANSI 渲染
- 现有的 `pkg/fiber/scheduler.go` - Fiber 调度器
- 现有的 `pkg/hooks/hooks.go` - Hooks 实现
- 现有的 `pkg/input/input.go` - 输入处理

## Risks

1. **终端兼容性** - 不同终端的 ANSI 支持程度不同
   - 缓解：使用广泛支持的 ANSI 序列，提供降级方案

2. **性能** - 高 FPS 动画可能导致 CPU 占用高
   - 缓解：默认节流 60fps，提供配置选项

3. **跨平台** - Windows 终端行为可能不同
   - 缓解：使用 go-isatty 和适当的 Windows API
