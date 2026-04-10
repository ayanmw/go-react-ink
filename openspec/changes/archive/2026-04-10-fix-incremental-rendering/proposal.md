# 完整差异分析与迁移路径

## 1. 架构差异矩阵

### 1.1 核心选项

| 选项字段 | React Ink | Go 实现当前状态 | 优先级 | 影响 |
|---------|-----------|----------------|--------|------|
| `IncrementalRendering` | ✅ 存在 | ❌ 不存在 | **P0** | 渲染性能、用户体验 |
| `Concurrent` | ✅ 支持 Legacy/Concurrent Root | ⚠️ 仅 Legacy | P1 | 高优先级更新中断支持 |
| `AlternateScreen` | ✅ 存在 | ✅ 已实现 | - | 正常 |
| `Interactive` | ✅ 自动检测 | ⚠️ 简化检测 | P2 | 非交互环境支持 |
| `KingKeyboard` | ✅ 支持 | ❌ 不存在 | P2 | 高级键盘功能 |
| `IsScreenReaderEnabled` | ✅ 支持 | ❌ 不存在 | P3 | 无障碍 |
| `OnRender` | ✅ 回调 | ❌ 不存在 | P2 | 渲染指标监控 |
| `PatchConsole` | ✅ 拦截 console | ⚠️ 简化 | P2 | 日志混合问题 |

### 1.2 Hook 差异

| Hook | React Ink | Go 实现当前状态 | 优先级 |
|------|-----------|----------------|--------|
| `UseState` | ✅ | ✅ | - |
| `UseEffect` | ✅ | ✅ | - |
| `UseMemo` | ✅ | ✅ | - |
| `UseCallback` | ✅ | ✅ | - |
| `UseRef` | ✅ | ✅ | - |
| **UseApp** | ✅ 返回 `{ exit, waitUntilRenderFlush }` | ❌ 不存在 | **P0** |
| **UseInput** | ✅ 完整键盘解析 | ❌ 不存在 | **P0** |
| `UseFocus` | ✅ 焦点管理 | ❌ 不存在 | P1 |
| `UseAnimation` | ✅ 动画驱动 | ⚠️ 基础实现 | P1 |
| `UseCursor` | ✅ 光标位置 | ❌ 不存在 | P2 |
| `UseStdin` / `UseStdout` | ✅ 流访问 | ❌ 不存在 | P2 |
| `UseWindowSize` | ✅ 窗口尺寸 | ❌ 不存在 | P2 |
| `UseTransition` | ✅ 并发模式专用 | ❌ 不存在 | P1（需 Concurrency） |

### 1.3 内置组件差异

| 组件 | React Ink | Go 实现当前状态 | 优先级 |
|------|-----------|----------------|--------|
| `Box` | ✅ | ✅ | - |
| `Text` | ✅ | ✅ | - |
| **Static** | ✅ 渐进式静态内容 | ❌ 不存在 | **P1** |
| `Transform` | ✅ 输出转换 | ❌ 不存在 | P3 |
| `Newline` | ✅ | ✅ | - |
| `Spacer` | ✅ | ✅ | - |

### 1.4 渲染流程差异

```
React Ink 渲染流程:
┌─────────────────────────────────────────────────────────────┐
│ React Component Tree                                         │
│        ↓                                                     │
│ Fiber Reconciler (LegacyRoot / ConcurrentRoot)               │
│        ↓                                                     │
│ Yoga Layout Engine → DOM Nodes                              │
│        ↓                                                     │
│ Renderer.render() → {output, outputHeight, staticOutput}    │
│        ↓                                                     │
│ Ink.onRender():                                             │
│   - 渲染器调用 render()                                     │
│   - 检查 outputHeight, staticOutput                        │
│   - 处理 <Static> 内容                                      │
│   - logUpdate 用户选项 (incremental)                        │
│        ↓                                                     │
│ LogUpdate.write()                                           │
│   - createIncremental() 或 createStandard()                 │
│   - 逐行 diff 或全屏刷新                                    │
│   - 光标位置管理 (useCursor)                                │
│   - 同步输出模式 (write-synchronized)                       │
└─────────────────────────────────────────────────────────────┘

Go 实现渲染流程:
┌─────────────────────────────────────────────────────────────┐
│ Go Component Tree (simplified)                              │
│        ↓                                                     │
│ Scheduler.Render()                                         │
│        ↓                                                     │
│ (简单布局，Yoga 集成不完整)                                │
│        ↓                                                     │
│ Ink.onRender():                                            │
│   - 调度器执行 Scheduler.Render()                           │
│   - 直接 log.Write(output)                                  │
│        ↓                                                     │
│ LogUpdate.Write()                                          │
│   - writeIncremental() (硬编码启用)                        │
│   - 无 <Static> 处理                                        │
│   - 无光标位置管理                                         │
│   - 无同步输出模式                                         │
└─────────────────────────────────────────────────────────────┘
```

## 2. 渐进式迁移计划

### Phase 1: 核心渲染修复 (P0)

**目标**: 修复 incremental 集成问题，实现基础 Hook

#### 1.1 Incremental Rendering 修复
- ✅ 设计已完成
- 实现：
  ```go
  // 修改 files:
  // - pkg/ink/options.go: 添加 IncrementalRendering
  // - pkg/ink/log_update.go: NewLogUpdate 添加参数
  // - pkg/ink/ink.go: 传递配置到 NewLogUpdate
  ```

#### 1.2 UseApp Hook
- 新建 `pkg/hooks/use_app.go`
- 集成到 Hook 上下文系统
- 返回组件：
  ```go
  type AppContext struct {
      Exit func(...any)
      WaitUntilRenderFlush func() error
  }
  ```

#### 1.3 UseInput Hook
- 新建 `pkg/hooks/use_input.go`
- 新建 `pkg/input/input_parser.go` (解析 ANSI 键盘序列)
- 新建 `pkg/input/parse_keypress.go` (键盘事件映射)
- 集成到 Xi）
- 暂不支持 Kitty keyboard (后期)

**产物**:
- 基础示例可运行：counter, use-input, use-animation

### Phase 2: 增强功能 (P1)

**目标**: 支持更多交互功能

#### 2.1 Static 组件
- 新建 `pkg/components/static.go`
- 实现 DOM 层面的 `staticNode` 支持
- 渲染逻辑修改 (`renderer.ts` 中的 staticOutput 分支)

#### 2.2 UseFocus Hook
- 新建 `pkg/hooks/use_focus.go`
- 焦点栈管理
- `useFocusManager`

#### 2.3 UseAnimation 增强
- 增加 `time`, `delta`, `reset` 返回值
- 支持 `pause/resume` (使用 `isActive` 选项)

#### 2.4 OnRender 回调
- 添加到 `RenderOptions`
- 用于渲染性能监控

**产物**:
- 高级示例可运行：static, use-focus, 增强版 use-animation

### Phase 3: 高级功能 (P2)

**目标**: 完善功能对齐

#### 3.1 UseCursor Hook
- 光标位置管理
- 集成 `cursor-helpers`

#### 3.2 UseWindowSize
- 监听终端尺寸变化
- 触发重新布局

#### 3.3 Transform 组件
- 输出转换管道

#### 3.4 同步输出模式
- 实现 `write-synchronized.ts` (bsu/esu)

**产物**:
- 完整功能对齐，所有 React Ink 基础用例可运行

### Phase 4: 高级扩展 (P3)

**目标**: 可选高级功能

#### 4.1 Keyboard Protocol
- Kitty keyboard 支持
- UseMeta, UseSuper 等

#### 4.2 Concurrency
- Concurrent Root 支持
- UseTransition 实现

#### 4.3 Screen Reader
- 无障碍模式

## 3. 示例/测试移植清单

### Counter 示例
**React Ink**: `examples/counter/counter.tsx`
```tsx
function Counter() {
  const [counter, setCounter] = React.useState(0);
  React.useEffect(() => {
    const timer = setInterval(() => {
      setCounter(prevCounter => prevCounter + 1);
    }, 100);
    return () => {
      clearInterval(timer);
    };
  }, []);
  return <Text color="green">{counter} tests passed</Text>;
}
render(<Counter />);
```

**Go Porting 要求**:
- ✅ UseState 支持
- ✅ UseEffect 支持
- ✅ Text/color 属性支持
- ⚠️ animation-demo 已有类似功能

优先级: **示例已实现同等功能，可跳过**

---

### Static 示例
**React Ink**: `examples/static/static.tsx`
```tsx
function Example() {
  const [tests, setTests] = useState<Array<{id: number; title: string}>>([]);

  React.useEffect(() => {
    let completedTests = 0;
    let timer: NodeJS.Timeout | undefined;

    const run = () => {
      if (completedTests++ < 10) {
        setTests(previousTests => [
          ...previousTests,
          {id: previousTests.length, title: `Test #${previousTests.length + 1}`}
        ]);
        timer = setTimeout(run, 100);
      }
    };
    run();
    return () => {
      clearTimeout(timer);
    };
  }, []);

  return (
    <>
      <Static items={tests}>
        {test => (
          <Box key={test.id}>
            <Text color="green">✔ {test.title}</Text>
          </Box>
        )}
      </Static>
      <Box marginTop={1}>
        <Text dimColor>Completed tests: {tests.length}</Text>
      </Box>
    </>
  );
}
```

**Go Porting 要求**:
- ❌ Static 组件未实现
- ✅ UseState 支持
- ✅ UseEffect 支持
- ✅ Box/Text 属性支持

优先级: **P1 - Phase 2 实现 Static 后**

---

### UseInput 示例
**React Ink**: `examples/use-input/use-input.tsx`
```tsx
function Robot() {
  const {exit} = useApp();
  const [x, setX] = React.useState(1);
  const [y, setY] = React.useState(1);

  useInput((input, key) => {
    if (input === 'q') { exit(); }
    if (key.leftArrow) { setX(Math.max(1, x - 1)); }
    if (key.rightArrow) { setX(Math.min(20, x + 1)); }
    if (key.upArrow) { setY(Math.max(1, y - 1)); }
    if (key.downArrow) { setY(Math.min(10, y + 1)); }
  });

  return (
    <Box flexDirection="column">
      <Text>Use arrow keys to move the face. Press "q" to exit.</Text>
      <Box height={12} paddingLeft={x} paddingTop={y}>
        <Text>^_^</Text>
      </Box>
    </Box>
  );
}
```

**Go Porting 要求**:
- ❌ UseApp 未实现
- ❌ UseInput 未实现
- ❌ 键盘解析器未实现
- ✅ UseState 支持
- ✅ Box/Text 属性支持

优先级: **P0 - Phase 1 必须实现**

---

### UseAnimation 示例
**React Ink**: `examples/use-animation/use-animation.tsx`
- 已在 animation-demo 中有类似实现
- Need to port exact behavior

优先级: **P1 - Phase 2**

---

### Incremental Rendering 示例
**React Ink**: `examples/incremental-rendering/incremental-rendering.tsx`
- 用于展示 incremental 渲染效果
- Complex UI with multiple rapidly updating sections

优先级: **P1 - Phase 2 完成 incremental 修复后**

---

### UseFocus 示例
**React Ink**: `examples/use-focus/use-focus.tsx`
```tsx
function Focus() {
  return (
    <Box flexDirection="column" padding={1}>
      <Box marginBottom={1}>
        <Text>Press Tab to focus next element...</Text>
      </Box>
      <Item label="First" />
      <Item label="Second" />
      <Item label="Third" />
    </Box>
  );
}

function Item({label}: {readonly label: string}) {
  const {isFocused} = useFocus();
  return (
    <Text>
      {label} {isFocused ? <Text color="green">(focused)</Text> : null}
    </Text>
  );
}
```

优先级: **P1 - Phase 2 实现 UseFocus**

---

### UseTransition 示例
**React Ink**: `examples/use-transition/use-transition.tsx`
- 需要并发模式支持
- 需要实现 useTransition

优先级: **P1 - 但需先完成 P2 Concurrency Phase 4**

---

### 测试 Fixture 移植

React Ink 测试使用 `test/fixtures/` 中的独立程序。移植策略：

#### Initial Test Set (Phase 1):
- `exit-on-exit.tsx` → 测试 `UseApp().exit()`
- `use-input.tsx` → 测试键盘输入解析

####中期测试集 (Phase 2):
- `fullscreen-no-extra-newline.tsx` → 测试全屏渲染
- `static-*.tsx` → 测试 Static 组件

#### 后期测试集 (Phase 3+):
- 所有 `use-input-*.tsx` fixture
- 光标处理 fixture

## 4. 输出一致性保证

### 4.1 渲染输出对比架构

```
┌─────────────────────────────────────────────────────────────┐
│                     对比测试架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  React Ink Output  ───────────┐                             │
│  (node + .tsx)                  │                             │
│                                 ↓                             │
│ ┌─────────────────────────┐  ┌─────────────────────────┐   │
│ │    Snapshot A1          │  │    Snapshot A2          │   │
 │ └─────────────────────────┘  └─────────────────────────┘   │
│                                         │                     │
│                                        ↓                      │
│                         ╔═══════════════════════════╗        │
│                         ║  diff 对比工具              ║        │
│                         ║  - 完全匹配 ✓               ║        │
│                         ║  - 预期差异解释              ║        │
│                         ╚═══════════════════════════╝        │
│                                                              │
│  Go React-Ink Output  ──────────────┐                        │
│  (go + .go → .gox → .go)           │                        │
│                                   ↓                         │
│ ┌─────────────────────────┐       │                          │
│ │    Snapshot B1          │       │                          │
│ └─────────────────────────┘       │                          │
│                                   │                          │
│           ┌───────────────────────┴───────────────────────┐ │
│           │                                              │ │
│           ↓                                              ↓ │
│   ┌──────────────────────┐                   ┌──────────────────┐ │
│   │  String 消除 ANSI     │                   │  结构化对比      │ │   │
│   │  StripAnsi 对比       │                   │  (metadata)      │ │   │   │
│   └──────────────────────┘                   └──────────────────┘ │   │
└───────────────────────────────────────────────────────────────┘
```

### 4.2 预期差异

#### 4.2.1 小差异 (可接受)
**React Ink**: 控制输出在 Coordinator 中
```
reconciler.discreteUpdates(() => {
  inputHandler(input, key);
});
```

**Go 方案**: 直接在 Input 处理器中调用
```go
go func() {
  for input := range inputChan {
    inputHandler(input, key)
  }
}()
```

#### 4.2.2 功能延迟 (需记录)
- Kitty keyboard: 后期实现
- Concurrency: React Concurrent 特性

### 4.3 对比工具链

**工具 1**: 简单文本对比
```go
func CompareOutput(expected, actual string) error {
    expected = stripAnsi(expected)
    actual = stripAnsi(actual)
    if expected != actual {
        return fmt.Errorf("output mismatch")
    }
    return nil
}
```

**工具 2**: 结构化对比 (for advanced tests)
```go
type RenderedFrame struct {
    Output       string
    OutputHeight int
    Static       []string  // <Static> sections
    CursorPos    *Position // 光标位置
}
```

## 5. 实现优先级矩阵

| 功能 | React Ink | Go 状态 | 实现成本 | 用户价值 | 优先级 | Phase |
|-----|-----------|--------|---------|---------|--------|-------|
| Incremental Config | ✅ | ❌ 断链 | 低 | 高 | P0 | 1 |
| UseApp | ✅ | ❌ | 低 | 高 | P0 | 1 |
| UseInput (核心) | ✅ | ❌ | 中 | 高 | P0 | 1 |
| UseAnimation (基础) | ✅ | ⚠️ 部分 | 低 | 中 | P1 | 2 |
| Static 组件 | ✅ | ❌ | 中 | 高 | P1 | 2 |
| UseFocus | ✅ | ❌ | 中 | 中 | P1 | 2 |
| OnRender 回调 | ✅ | ❌ | 低 | 低-中 | P2 | 3 |
| Sync Output Mode | ✅ | ❌ | 中 | 中 | P2 | 3 |
| UseCursor Hook | ✅ | ❌ | 低 | 中 | P2 | 3 |
| UseWindowSize | ✅ | ❌ | 低 | 低 | P2 | 3 |
| kitty keyboard | ✅ | ❌ | 高 | 低 | P3 | 4 |
| Concurrency | ✅ | 部分 | 高 | 中 | P2 | 4 |
| UseTransition | ✅ | ❌ | 中 | 中 | P2 (需 Concurrency) | 4 |
| Screen Reader | ✅ | ❌ | 中 | 低 | P3 | N/A |

## 6. 里程碑

### 里程碑 1: 使基础示例可运行
**状态**: 设计完成，进入实施
**包含**:
- ✅ Incremental rendering 配置集成
- ✅ UseApp Hook
- ✅ UseInput Hook (核心 keys: 方向, 字母, Ctrl, Shift)
- ✅ 输入解析器 (parseKeypress 核心子集)

**示例可运行**:
- counter (已有类似)
- use-input
- use-animation (增强版)

---

### 里程碑 2: 增强交互功能
**包含**:
- ✅ Static 组件
- ✅ UseFocus Hook
- ✅ UseAnimation 增强
- ✅ OnRender 回调

**示例可运行**:
- static
- use-focus
- incremental-rendering

---

### 里程碑 3: 完整对齐
**包含**:
- ✅ 光标位置管理
- ✅ 同步输出模式
- ✅ 窗口尺寸监听
- ✅ Transform 组件

**测试覆盖**: 全部 React Ink 基础测试

---

### 里程碑 4: 高级功能
**包含**:
- ❌ Kitty keyboard (可选)
- ❌ 并发模式
- ❌ UseTransition
- ❌ Screen reader (可选)

## 7. 当前上下文

**已完成的工作**:
- ✅ 基础 JSX 编译器
- ✅ Fiber 协调器架构
- ✅ Flexbox 布局 (简化)
- ✅ 渲染循环
- ✅ Alternate screen buffer
- ✅ 增量渲染 (基础版本)
- ✅ 动画框架

**为什么当前状态不够?**
1. **不可配置**: incremental hard-coded，用户无法切换
2. **缺少关键 Hook**: UseApp, UseInput 是交互式应用的核心
3. **示例孤立**: animation-demo 自行实现 incremental

**下一步行动**:
1. 修复 incremental 渲染配置传递 (已设计方案)
2. 实现 UseApp Hook
3. 实现 UseInput Hook (+ 键盘解析)
4. 移植示例

这个分析完整覆盖了 React Ink 的架构、功能、策略和移植计划。
