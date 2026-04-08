# Go-Ink Examples

本目录包含多个示例程序，演示 Go-Ink 的各种功能特性。

## 示例列表

### 1. hello-world
最简单的 Go-Ink 应用，展示基本的组件创建和渲染。

```bash
cd hello-world && go run main.go
```

**演示特性：**
- Box 和 Text 组件
- 基本属性设置（color, bold, padding）
- 组件嵌套

---

### 2. counter-app
计数器应用，演示 useState Hook 的使用。

```bash
cd counter-app && go run main.go
```

**演示特性：**
- useState Hook
- 状态管理
- 状态更新函数

---

### 3. todo-list
待办事项列表，演示列表渲染和复杂数据结构状态管理。

```bash
cd todo-list && go run main.go
```

**演示特性：**
- 列表渲染
- 结构体状态
- 条件渲染
- 多样式组合

---

### 4. flexbox-layout
Flexbox 布局示例，演示各种 Flex 布局属性。

```bash
cd flexbox-layout && go run main.go
```

**演示特性：**
- flexDirection (row/column)
- justifyContent (space-between/center)
- alignItems
- flexGrow
- 嵌套布局

---

### 5. hooks-demo
React Hooks 完整演示，展示所有可用的 Hooks。

```bash
cd hooks-demo && go run main.go
```

**演示特性：**
- useState - 状态管理
- useRef - 引用
- useMemo - 记忆化计算
- useCallback - 回调记忆化
- useInput - 输入处理
- useFocus - 焦点管理
- useAnimation - 动画
- useWindowSize - 窗口尺寸

---

### 6. text-styling
文本样式示例，演示所有文本样式属性。

```bash
cd text-styling && go run main.go
```

**演示特性：**
- 颜色 (color)
- 背景色 (backgroundColor)
- 粗体 (bold)
- 斜体 (italic)
- 下划线 (underline)
- 变暗 (dim)
- 闪烁 (blink)
- 反转 (reverse)
- 样式组合

---

### 7. interactive-input
交互式输入示例，演示输入处理和焦点管理。

```bash
cd interactive-input && go run main.go
```

**演示特性：**
- 键盘事件处理
- 快捷键注册
- 焦点管理器
- 光标控制
- 标准输入/输出
- 终端信息

---

### 8. animation-demo
动画示例，演示动画和帧更新。

```bash
cd animation-demo && go run main.go
```

**演示特性：**
- useAnimation Hook
- Spinner 动画
- 进度条
- 帧更新
- 动画控制 (play/pause/stop)

---

## 快速运行

使用 Makefile 运行所有示例：

```bash
# 运行所有示例
make run-all

# 运行特定示例
make run-hello
make run-counter
make run-todo
make run-flexbox
make run-hooks
make run-styling
make run-input
make run-animation

# 构建所有示例
make build-all

# 清理构建产物
make clean
```

## 功能覆盖矩阵

| 示例 | 组件 | Hooks | 布局 | 样式 | 输入 | 动画 |
|------|------|-------|------|------|------|------|
| hello-world | ✓ | | | ✓ | | |
| counter-app | ✓ | ✓ | | | | |
| todo-list | ✓ | ✓ | | ✓ | | |
| flexbox-layout | ✓ | | ✓ | | | |
| hooks-demo | ✓ | ✓ | | | ✓ | ✓ |
| text-styling | ✓ | | | ✓ | | |
| interactive-input | ✓ | ✓ | | | ✓ | |
| animation-demo | ✓ | ✓ | | | | ✓ |

## 组件覆盖

- [x] Box - 容器组件
- [x] Text - 文本组件
- [x] Spacer - 间隔组件
- [x] Newline - 换行组件
- [x] Static - 静态内容
- [x] Transform - 变换组件
- [x] Fragment - 片段组件

## Hooks 覆盖

- [x] useState - 状态管理
- [x] useEffect - 副作用
- [x] useLayoutEffect - 布局副作用
- [x] useRef - 引用
- [x] useMemo - 记忆化
- [x] useCallback - 回调记忆化
- [x] useInput - 输入处理
- [x] useFocus - 焦点
- [x] useFocusManager - 焦点管理
- [x] useApp - 应用
- [x] useStdin - 标准输入
- [x] useStdout - 标准输出
- [x] useStderr - 标准错误
- [x] useCursor - 光标
- [x] useAnimation - 动画
- [x] useWindowSize - 窗口尺寸
- [x] useBoxMetrics - 盒子尺寸
- [x] usePaste - 粘贴
- [x] useIsScreenReaderEnabled - 屏幕阅读器

## 布局属性覆盖

- [x] flexDirection (row/column/row-reverse/column-reverse)
- [x] justifyContent (flex-start/flex-end/center/space-between/space-around/space-evenly)
- [x] alignItems (flex-start/flex-end/center/stretch/baseline)
- [x] flexWrap (nowrap/wrap/wrap-reverse)
- [x] flexGrow
- [x] flexShrink
- [x] flexBasis
- [x] width/height
- [x] minWidth/maxWidth/minHeight/maxHeight
- [x] padding/paddingX/paddingY/paddingTop/paddingRight/paddingBottom/paddingLeft
- [x] margin/marginX/marginY/marginTop/marginRight/marginBottom/marginLeft
- [x] borderStyle/borderColor

## 文本样式覆盖

- [x] color
- [x] backgroundColor
- [x] bold
- [x] italic
- [x] underline
- [x] dim
- [x] blink
- [x] reverse
- [x] wrap

## 学习路径

推荐的学习顺序：

1. **hello-world** - 了解基本概念
2. **text-styling** - 学习文本样式
3. **flexbox-layout** - 掌握布局系统
4. **counter-app** - 学习状态管理
5. **todo-list** - 实践列表渲染
6. **hooks-demo** - 深入理解所有 Hooks
7. **interactive-input** - 学习输入处理
8. **animation-demo** - 学习动画效果
