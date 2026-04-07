# Spec: Go-Ink TUI Framework

## Overview

Go-Ink 是 React Ink 的 Golang 完整复刻，提供声明式终端 UI 开发体验。

## Capabilities

### CAP-001: JSX 语法支持

**描述**: 支持 JSX 语法编写终端 UI 组件

**功能**:
- .gox 文件包含 Go 代码 + JSX
- 编译时转换为纯 Go 代码
- 支持所有 JSX 特性 (嵌套、属性、表达式)

**示例**:
```go
func Counter() Element {
    count, setCount := useState(0)
    return (
        <Box flexDirection="column">
            <Text color="green">Count: {count}</Text>
        </Box>
    )
}
```

### CAP-002: Flexbox 布局

**描述**: 支持 CSS Flexbox 布局算法

**功能**:
- flexDirection (row/column)
- justifyContent (5 种对齐)
- alignItems (5 种对齐)
- flexGrow / flexShrink
- width / height
- padding / margin

### CAP-003: 增量渲染

**描述**: 只更新变化的部分，最小化终端输出

**功能**:
- 双缓冲 Diff
- ANSI 序列优化
- 时间切片渲染

### CAP-004: React Hooks

**描述**: 支持 React Hooks 模式

**功能**:
- useState: 状态管理
- useEffect: 副作用处理
- useInput: 键盘输入
- useApp: 应用实例
- useFocus: 焦点管理

### CAP-005: 跨平台支持

**描述**: 支持 Windows/macOS/Linux

**功能**:
- 统一终端抽象
- ANSI 兼容处理
- 信号处理

## API Specification

### 组件 API

#### Box

```go
func Box(props Props, children ...Element) Element
```

**Props**:
- flexDirection: "row" | "column"
- justifyContent: "flex-start" | "center" | "flex-end" | "space-between" | "space-around"
- alignItems: "stretch" | "flex-start" | "center" | "flex-end"
- width: int
- height: int
- padding: int | EdgeInsets
- margin: int | EdgeInsets

#### Text

```go
func Text(props Props, children ...Element) Element
```

**Props**:
- color: Color
- backgroundColor: Color
- bold: bool
- dim: bool
- underline: bool

### Hooks API

#### useState

```go
func useState[T any](initial T) (T, func(T))
```

#### useEffect

```go
func useEffect(create func() func(), deps []any)
```

#### useInput

```go
func useInput(handler func(InputEvent))
```

## Compatibility

### 与 React Ink 兼容性

| 特性 | React Ink | Go-Ink | 状态 |
|-----|-----------|--------|------|
| Box 组件 | ✅ | ✅ | 兼容 |
| Text 组件 | ✅ | ✅ | 兼容 |
| Static 组件 | ✅ | ✅ | 兼容 |
| useState | ✅ | ✅ | 兼容 |
| useEffect | ✅ | ✅ | 兼容 |
| useInput | ✅ | ✅ | 兼容 |
| useApp | ✅ | ✅ | 兼容 |
| Flexbox | ✅ | ✅ | 兼容 |
| 增量渲染 | ✅ | ✅ | 兼容 |

## Performance

### 目标

- 渲染延迟: < 16ms (60fps)
- 内存占用: < 10MB (简单应用)
- 启动时间: < 100ms

### 测量方法

```bash
go test -bench=. -benchmem
```

## Testing

### 测试策略

1. **单元测试**: 每个模块
2. **集成测试**: 编译器 + 运行时
3. **对比测试**: 与 React Ink 结果对比

### 测试覆盖

- 目标覆盖率: 80%
- 关键模块: 100%

## References

- [React Ink Documentation](https://github.com/vadimdemedes/ink)
- [React Reconciler](https://github.com/facebook/react/tree/main/packages/react-reconciler)
- [Yoga Layout](https://github.com/facebook/yoga)
