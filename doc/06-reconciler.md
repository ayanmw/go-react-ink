# React 协调器设计

> 复刻 React Fiber 架构，实现声明式 UI 更新

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  Fiber Reconciler 架构                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   用户组件树                                                 │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │  Scheduler  │  调度更新任务                             │
│   └─────────────┘                                          │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │  Work Loop  │  协调阶段 (可中断)                        │
│   └─────────────┘                                          │
│       │                                                     │
│       ├──▶ beginWork()    创建/复用 Fiber                   │
│       ├──▶ completeWork() 完成子树                          │
│       └──▶ reconcileChildren() Diff 算法                   │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │  Commit     │  提交阶段 (不可中断)                      │
│   │  Phase      │                                          │
│   └─────────────┘                                          │
│       │                                                     │
│       ├──▶ commitPlacement()    插入                        │
│       ├──▶ commitUpdate()       更新                        │
│       └──▶ commitDeletion()     删除                        │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │ Host Config │  宿主环境适配                             │
│   └─────────────┘                                          │
│       │                                                     │
│       └──▶ Layout + Renderer                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Fiber 节点结构

```go
// Fiber 节点
type Fiber struct {
    // 标识
    ID        uint64          // 唯一 ID
    Tag       FiberTag        // 节点类型
    
    // 类型信息
    Type      any             // 组件函数 或 元素标签名
    Key       any             // 列表 key
    
    // 状态
    Props     any             // 输入属性
    State     any             // 组件状态 (hooks)
    EffectTag EffectTag       // 副作用标记
    
    // 树结构
    Return    *Fiber          // 父节点
    Child     *Fiber          // 第一个子节点
    Sibling   *Fiber          // 下一个兄弟节点
    Index     int             // 在父节点子列表中的索引
    
    // 替代节点 (用于 diff)
    Alternate *Fiber          // 对应的上一次渲染节点
    
    // 子节点列表 (用于 diff)
    PendingChildren []*Fiber  // 新的子节点列表
    
    // 输出
    StateNode any             // 宿主实例 (布局节点/渲染节点)
    
    // 副作用链表
    FirstEffect *Fiber        // 第一个有副作用的子节点
    LastEffect  *Fiber        // 最后一个有副作用的子节点
    NextEffect  *Fiber        // 下一个有副作用的节点
    
    // 调度优先级
    Lane       Lane           // 优先级车道
    EventTime  time.Time      // 事件时间
}

// Fiber 类型标签
type FiberTag int

const (
    FunctionComponent FiberTag = iota  // 函数组件
    HostElement                        // 宿主元素 (Box, Text)
    HostText                           // 文本节点
    Fragment                           // Fragment
    Root                               // 根节点
)

// 副作用标记
type EffectTag uint16

const (
    NoEffect      EffectTag = 0
    Placement     EffectTag = 1 << 0  // 插入
    Update        EffectTag = 1 << 1  // 更新
    Deletion      EffectTag = 1 << 2  // 删除
    ChildDeletion EffectTag = 1 << 3  // 子节点删除
)
```

---

## 调度器 (Scheduler)

```go
// 调度器
type Scheduler struct {
    root       *Fiber          // 当前根节点
    workRoot   *Fiber          // 正在处理的工作
    next       *Fiber          // 下一个工作单元
    
    lanes      LaneSet         // 待处理车道
    callback   func()          // 完成回调
    
    // 时间切片
    frameInterval time.Duration // 默认 16ms
    deadline      time.Time
    
    // 优先级
    normalPriority Priority
    highPriority   Priority
}

// 调度更新
func (s *Scheduler) ScheduleUpdate(root *Fiber, lane Lane) {
    // 1. 标记车道
    s.lanes |= lane
    
    // 2. 设置工作根
    if s.workRoot == nil {
        s.workRoot = root
    }
    
    // 3. 请求工作
    s.requestWork()
}

// 请求工作 (异步)
func (s *Scheduler) requestWork() {
    // 使用 goroutine 模拟异步调度
    go s.performWork()
}

// 执行工作
func (s *Scheduler) performWork() {
    defer s.workLoop()
    
    for s.workRoot != nil && s.lanes != 0 {
        s.workLoop()
    }
}
```

---

## 工作循环 (Work Loop)

```go
// 工作循环
func (s *Scheduler) workLoop() {
    // 设置截止时间
    s.deadline = time.Now().Add(s.frameInterval)
    
    // 协调阶段 (可中断)
    for s.next != nil && !s.shouldYield() {
        s.next = s.performUnitOfWork(s.next)
    }
    
    // 提交阶段 (不可中断)
    if s.next == nil && s.workRoot != nil {
        s.commitRoot()
    }
}

// 是否应该让出控制权
func (s *Scheduler) shouldYield() bool {
    return time.Now().After(s.deadline)
}

// 执行一个工作单元
func (s *Scheduler) performUnitOfWork(unit *Fiber) *Fiber {
    // 1. 开始工作 (向下遍历)
    next := s.beginWork(unit)
    
    if next != nil {
        return next
    }
    
    // 2. 完成工作 (向上回溯)
    return s.completeWork(unit)
}
```

---

## 开始工作 (beginWork)

```go
// 开始工作
func (s *Scheduler) beginWork(fiber *Fiber) *Fiber {
    switch fiber.Tag {
    case FunctionComponent:
        return s.updateFunctionComponent(fiber)
    case HostElement:
        return s.updateHostElement(fiber)
    case HostText:
        return s.updateHostText(fiber)
    case Fragment:
        return s.updateFragment(fiber)
    }
    return nil
}

// 更新函数组件
func (s *Scheduler) updateFunctionComponent(fiber *Fiber) *Fiber {
    // 1. 准备 hooks 上下文
    s.prepareHooks(fiber)
    
    // 2. 执行组件函数
    component := fiber.Type.(func(any) Element)
    children := component(fiber.Props)
    
    // 3. 处理 hooks
    s.finishHooks(fiber)
    
    // 4. 协调子节点
    return s.reconcileChildren(fiber, children)
}

// 更新宿主元素
func (s *Scheduler) updateHostElement(fiber *Fiber) *Fiber {
    // 1. 协调子节点
    return s.reconcileChildren(fiber, fiber.Props.(Props).Children)
}
```

---

## 子节点协调 (Diff 算法)

```go
// 协调子节点
func (s *Scheduler) reconcileChildren(returnFiber *Fiber, newChildren []Element) *Fiber {
    // 获取旧子节点
    oldFiber := returnFiber.Child
    
    var prevNewFiber *Fiber
    var newFiber *Fiber
    
    newIndex := 0
    
    // 第一轮: 对比旧节点和新节点
    for oldFiber != nil && newIndex < len(newChildren) {
        newChild := newChildren[newIndex]
        
        // 尝试复用
        if s.canReuse(oldFiber, newChild) {
            // 更新现有节点
            newFiber = s.updateFiber(oldFiber, newChild, newIndex)
        } else {
            // 不能复用，创建新节点
            newFiber = s.createFiber(newChild, newIndex)
            
            // 标记旧节点删除
            oldFiber.EffectTag |= Deletion
        }
        
        // 链接兄弟节点
        if prevNewFiber == nil {
            returnFiber.Child = newFiber
        } else {
            prevNewFiber.Sibling = newFiber
        }
        
        prevNewFiber = newFiber
        oldFiber = oldFiber.Sibling
        newIndex++
    }
    
    // 旧节点多余，标记删除
    if oldFiber != nil {
        for ; oldFiber != nil; oldFiber = oldFiber.Sibling {
            oldFiber.EffectTag |= Deletion
            s.trackEffect(oldFiber)
        }
    }
    
    // 新节点多余，创建
    for newIndex < len(newChildren) {
        newChild := newChildren[newIndex]
        newFiber = s.createFiber(newChild, newIndex)
        
        if prevNewFiber == nil {
            returnFiber.Child = newFiber
        } else {
            prevNewFiber.Sibling = newFiber
        }
        
        prevNewFiber = newFiber
        newIndex++
    }
    
    // 返回第一个子节点
    return returnFiber.Child
}

// 判断是否可复用
func (s *Scheduler) canReuse(oldFiber *Fiber, newChild Element) bool {
    // 类型相同
    if oldFiber.Type != newChild.Type() {
        return false
    }
    
    // key 相同
    oldKey := oldFiber.Key
    newKey := newChild.Key()
    return oldKey == newKey
}
```

---

## 完成工作 (completeWork)

```go
// 完成工作
func (s *Scheduler) completeWork(fiber *Fiber) *Fiber {
    switch fiber.Tag {
    case HostElement:
        // 创建/更新宿主实例
        if fiber.EffectTag&Placement != 0 {
            fiber.StateNode = s.hostConfig.CreateInstance(fiber.Type, fiber.Props)
        } else if fiber.EffectTag&Update != 0 {
            s.hostConfig.UpdateInstance(fiber.StateNode, fiber.Props)
        }
        
    case HostText:
        // 创建/更新文本节点
        text := fiber.Props.(string)
        if fiber.StateNode == nil {
            fiber.StateNode = s.hostConfig.CreateTextInstance(text)
        } else {
            s.hostConfig.UpdateTextInstance(fiber.StateNode, text)
        }
    }
    
    // 收集副作用
    s.collectEffects(fiber)
    
    // 返回兄弟节点或父节点
    if fiber.Sibling != nil {
        return fiber.Sibling
    }
    return fiber.Return
}

// 收集副作用
func (s *Scheduler) collectEffects(fiber *Fiber) {
    // 将有副作用的子节点链接到父节点
    returnFiber := fiber.Return
    if returnFiber == nil {
        return
    }
    
    // 首次添加
    if returnFiber.FirstEffect == nil {
        returnFiber.FirstEffect = fiber.FirstEffect
        returnFiber.LastEffect = fiber.LastEffect
    } else {
        // 追加到链表
        returnFiber.LastEffect.NextEffect = fiber.FirstEffect
        returnFiber.LastEffect = fiber.LastEffect
    }
    
    // 节点自身有副作用
    if fiber.EffectTag != NoEffect {
        if returnFiber.LastEffect != nil {
            returnFiber.LastEffect.NextEffect = fiber
        } else {
            returnFiber.FirstEffect = fiber
        }
        returnFiber.LastEffect = fiber
    }
}
```

---

## 提交阶段 (Commit Phase)

```go
// 提交根节点
func (s *Scheduler) commitRoot() {
    root := s.workRoot
    
    // 处理删除
    s.commitDeletion(root)
    
    // 处理插入和更新
    s.commitMutation(root)
    
    // 处理布局
    s.commitLayout(root)
    
    // 清理
    s.workRoot = nil
    s.next = nil
}

// 提交变更
func (s *Scheduler) commitMutation(root *Fiber) {
    for fiber := root.FirstEffect; fiber != nil; fiber = fiber.NextEffect {
        switch {
        case fiber.EffectTag&Placement != 0:
            s.commitPlacement(fiber)
        case fiber.EffectTag&Update != 0:
            s.commitUpdate(fiber)
        }
    }
}

// 提交插入
func (s *Scheduler) commitPlacement(fiber *Fiber) {
    parent := fiber.Return.StateNode
    node := fiber.StateNode
    
    s.hostConfig.AppendChild(parent, node)
}

// 提交更新
func (s *Scheduler) commitUpdate(fiber *Fiber) {
    node := fiber.StateNode
    props := fiber.Props
    
    s.hostConfig.UpdateInstance(node, props)
}

// 提交删除
func (s *Scheduler) commitDeletion(root *Fiber) {
    for fiber := root.FirstEffect; fiber != nil; fiber = fiber.NextEffect {
        if fiber.EffectTag&Deletion != 0 {
            s.commitDelete(fiber)
        }
    }
}

func (s *Scheduler) commitDelete(fiber *Fiber) {
    parent := fiber.Return.StateNode
    node := fiber.StateNode
    
    s.hostConfig.RemoveChild(parent, node)
}
```

---

## Host Config 接口

```go
// 宿主环境配置
type HostConfig interface {
    // 创建实例
    CreateInstance(typeName string, props Props) Instance
    CreateTextInstance(text string) Instance
    
    // 更新实例
    UpdateInstance(instance Instance, props Props)
    UpdateTextInstance(instance Instance, text string)
    
    // 树操作
    AppendChild(parent, child Instance)
    InsertBefore(parent, child, before Instance)
    RemoveChild(parent, child Instance)
    
    // 布局
    ScheduleLayout(instance Instance)
    
    // 渲染
    ScheduleRender()
    
    // 属性准备
    PrepareUpdate(instance Instance, oldProps, newProps Props) UpdatePayload
}

// 实例接口
type Instance interface {
    // 布局节点
    LayoutNode() *LayoutNode
    
    // 渲染节点
    RenderNode() *RenderNode
}
```

---

## Hooks 集成

```go
// Hooks 上下文
type HooksContext struct {
    fiber      *Fiber
    hookIndex  int
    hooks      []Hook
}

// Hook 接口
type Hook interface {
    hasEffect() bool
}

// State Hook
type StateHook struct {
    value    any
    dispatch func(any)
}

// Effect Hook
type EffectHook struct {
    create   func() func()
    destroy  func()
    deps     []any
}

// 准备 Hooks
func (s *Scheduler) prepareHooks(fiber *Fiber) {
    s.hooksContext = &HooksContext{
        fiber:     fiber,
        hookIndex: 0,
    }
    
    // 复用上一次的 hooks
    if fiber.Alternate != nil {
        s.hooksContext.hooks = fiber.Alternate.State.([]Hook)
    } else {
        s.hooksContext.hooks = []Hook{}
    }
}

// 完成 Hooks
func (s *Scheduler) finishHooks(fiber *Fiber) {
    fiber.State = s.hooksContext.hooks
}
```

---

## 模块结构

```
reconciler/
├── reconciler.go      # 协调器主入口
├── fiber.go           # Fiber 节点定义
├── scheduler.go       # 调度器
├── workloop.go        # 工作循环
├── beginwork.go       # 开始工作
├── completework.go    # 完成工作
├── commitwork.go      # 提交工作
├── childfiber.go      # 子节点协调 (Diff)
├── hostconfig.go      # Host Config 接口
├── hooks.go           # Hooks 集成
├── lane.go            # 优先级车道
└── effect.go          # 副作用处理
```

---

## 实现决策

### 时间切片策略

**决策：固定 16ms 切片，后续可扩展自适应**

```
固定切片 (当前实现):
├── 默认 16ms (60fps)
├── 简单可靠
└── 与 React 保持一致

自适应切片 (可选扩展):
├── 监控实际渲染时间
├── 动态调整切片长度
└── 策略:
    ├── 简单渲染 (<5ms)  → 切片 20ms
    ├── 中等渲染 (5-15ms) → 切片 16ms
    └── 复杂渲染 (>15ms) → 切片 10ms
```

### 优先级系统 (Lane)

**决策：简化版 3 级优先级**

```
React 完整 Lane 系统 (31 个车道，非常复杂)

Go-Ink 简化版:
├── High   - 用户输入、同步更新
├── Normal - 普通渲染
└── Low    - 空闲时处理

理由:
├── 终端场景不需要完整的 31 级优先级
├── 简化实现降低复杂度
└── 足够覆盖终端 TUI 使用场景
```

---

## 总结

### 核心要点

1. **Fiber 架构**
   - 树形结构 + 链表结构
   - 支持可中断渲染
   - 副作用链表优化提交

2. **两阶段提交**
   - 协调阶段 (Render Phase): 可中断，计算 Diff
   - 提交阶段 (Commit Phase): 不可中断，应用变更

3. **Diff 算法**
   - 同级比较
   - Key 匹配复用
   - 最小化操作

4. **Host Config**
   - 解耦协调器与渲染器
   - 支持不同宿主环境
   - 统一接口抽象

5. **Go 特性适配**
   - Goroutine 替代 event loop
   - Channel 替代 callback
   - 时间切片保证响应性
