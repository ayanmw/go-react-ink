# CLAUDE.md

本文件提供了 React Ink 代码库的指导和说明，供 Claude Code 在此代码库中工作时参考。

## 项目概述

React Ink 是一个基于 React 的终端 UI 开发库，用于构建 CLI 应用。它使用 Yoga 布局引擎实现 Flexbox 布局，提供与浏览器中类似的组件化开发体验。

## 常用命令

### 构建
- `npm run build` - 使用 TypeScript 编译源代码到 build/ 目录
- `npm run dev` - 监听模式，自动重新编译

### 测试
- `npm test` - 运行完整测试套件（包括类型检查、lint、ava 测试）
- `npm run typecheck` - 仅运行 TypeScript 类型检查
- `npm run lint` - 仅运行 XO lint 检查

### 示例
- `npm run example examples/<name>/<name>.tsx` - 运行特定示例
  - 例如: `npm run example examples/counter/counter.tsx`
  - 例如: `npm run example examples/static/static.tsx`
  - 例如: `npm run example examples/borders/borders.tsx`

### 格式化
- `npm run format` - 使用 Prettier 格式化代码

## 架构概览

### 核心文件结构
- `src/` - 主要源代码目录
  - `components/` - 内置组件 (Box, Text, Static, Newline, Spacer, Transform, AppContext 等)
  - `hooks/` - React Hooks (useInput, useApp, useStdin, useStdout, useFocus 等)
  - `render.ts` - 主渲染器，协调器入口
  - `reconciler.ts` - 基于 react-reconciler 的自定义协调器 (Fiber 架构)
  - `renderer.ts` - 输出渲染器，处理终端输出
  - `dom.ts` - DOM 节点抽象 (DOMElement, TextNode)
  - `render-to-string.ts` - 同步字符串渲染（无副作用，用于测试）
  - `render-node-to-output.ts` - 将 DOM 节点转换为终端输出
  - `styles.ts` - Flexbox 样式处理
  - `wrap-text.ts` - 文本换行处理
  - `measure-element.ts` - 元素尺寸测量功能
  - `kitty-keyboard.ts` - Kitty 终端键盘协议支持
  - `log-update.ts` - 终端输出同步更新逻辑

### 渲染流程
1. `render()` 创建渲染器实例
2. Reconciler 使用 Fiber 架构调度 React 组件更新
3. DOM 节点通过 Yoga 布局引擎计算布局（Flexbox）
4. Renderer 将布局转换为 ANSI 终端输出
5. 输出经过 diff 计算最小化更新，控制闪烁

### React 协调器配置
Ink 基于官方 `react-reconciler` 包实现自定义协调器，支持:
- Concurrent Mode (`concurrent: true` 选项)
- 可中断渲染
- Fiber 节点优先级调度

## 测试与示例

### Test Helpers
测试目录 `test/helpers/` 包含实用工具:
- `create-stdout.js` - 创建模拟 stdout 流用于测试
- `create-stdin.js` - 创建模拟 stdin 流用于测试
- `render-to-string.js` - 同步渲染到字符串（无终端副作用）

### 测试关键点
- 使用 AVA 框架运行测试
- `renderToString()` 用于静态检查，不启动终端会话
- `renderAsync()` 用于需要终端状态的测试
- 套件包括组件测试、Flexbox 布局测试、Hooks 测试、聚焦管理等

### 示例说明
- 用于参考完整 UI 实现
- 运行时会产生动态终端输出，需 Ctrl+C 退出
- `examples/counter/` - 简单计数器，展示实时更新
- `examples/static/` - <Static> 组件用法，固定已完成内容
- `examples/use-focus/` - 焦点管理
- `examples/use-input/` - 用户输入处理
- `examples/borders/` - 边框样式
- `examples/box-backgrounds/` - 背景色

## 内建组件

### `<Text>`
- 显示文本，应用颜色/样式（bold, italic, underline 等）
- 支持颜色/背景色、文本换行（wrap, truncate-*）

### `<Box>`
- 核心布局容器，实现 Flexbox
- 支持所有 Flexbox 属性（flexDirection, justifyContent, alignItems 等）
- 支持 padding, margin, gap, width/height
- 支持边框（borderStyle, borderColor）
- 支持背景色（backgroundColor）

### `<Static>`
- 将子项永久渲染在动态内容上方
- 用于显示已完成任务、日志等一旦渲染不再更改的内容

### `<Spacer>`
- 弹性空间元素，沿主轴展开填充

### `<Transform>`
- 在写入输出前对渲染字符串进行转换（如大小写、添加链接等）

### `<Newline>`
- 插入换行符

## 关键 Hooks

### `useInput(inputHandler, options?)`
- 监听用户键盘输入
- 回调参数: `input` (字符串), `key` (对象，包含 leftArrow, return, ctrl 等)

### `useApp()`
- 返回应用生命周期方法:
  - `exit(errorOrResult?)` - 退出应用
  - `waitUntilRenderFlush()` - 等待渲染刷新到 stdout

### `useStdin() / useStdout() / useStderr()`
- 访问标准流及工具

### `useFocus(options?)`
- 焦点状态管理
- 选项: autoFocus, isActive, id

### `useFocusManager()`
- 焦点管理（focusNext, focusPrevious, focus, activeId）

### `useCursor()`
- 控制光标位置 (setCursorPosition)

### `useWindowSize()`
- 终端尺寸监听，调整时重渲染

### `useAnimation(options?)`
- 动画驱动，返回 frame, time, delta, reset

### `useBoxMetrics(ref)`
- 跟踪 `<Box>` 元素的布局度量值（width, height, left, top）

## 渲染器配置

`render(tree, options)` 支持的关键选项:
- `stdout/stdin/stderr` - 流配置，默认使用 process.*
- `exitOnCtrlC` - 自动 Ctrl+C 退出（默认 true）
- `patchConsole` - 拦截 console.* 输出避免与 Ink 输出冲突（默认 true）
- `maxFps` - 最大渲染帧率（默认 30）
- `concurrent` - 启用 React Concurrent Mode（默认 false）
- `interactive` - 覆盖自动交互模式检测（默认基于 CI 和 TTY 检测）
- `alternateScreen` - 使用终端备用屏幕缓冲区（默认 false）
- `kittyKeyboard` - 启用 Kitty 键盘协议（支持更多键盘信息）
- `debug` - 调试模式，不替换前一次渲染（默认 false）
- `incrementalRendering` - 增量渲染模式，只更新变更行（默认 false）

## 开发注意事项

### TypeScript 配置
- 源代码在 `src/` 目录，编译输出到 `build/`
- 使用 `@sindresorhus/tsconfig` 基础配置
- JSX 编译为 `react`（JIT 需要运行时）

### 依赖的关键包
- `react-reconciler` - 官方 React 协调器底座
- `scheduler` - React 调度器
- `yoga-layout` - Yoga Flexbox 布局引擎

### 行为注意
- 在终端环境中，应用会保持存活直到 unmount 或 exit() 调用
- <Static> 只渲染 items 中的新增项，修改已存在项不会重新渲染
- 长时间运行的示例/测试需手动终止（终端环境模拟）

### 测试编写
- 使用 `renderToString()` 进行无副作用验证
- 使用测试 helpers 中的模拟流进行复杂场景测试
- 使用 `ava` 断言，关注 `t.is()` 和 `t.deepEqual()`
