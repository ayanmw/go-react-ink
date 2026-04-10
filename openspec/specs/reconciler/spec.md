# reconciler Specification

## Purpose
定义 Go-Ink Fiber 协调器的设计和实现，实现 React 风格的可中断渲染。

## Requirements

### Requirement: 协调器使用 Fiber 架构
协调器 SHALL 使用 Fiber 树结构进行增量渲染。

#### Scenario: Fiber 树结构
- **WHEN** 协调器创建组件树
- **THEN** 每个 Fiber 节点包含 child, sibling, return 指针

#### Scenario: 双缓冲
- **WHEN** 协调器渲染
- **THEN** 使用 currentTree 和 workInProgressTree 双缓冲

### Requirement: 协调器实现调度器
协调器 SHALL 使用调度器管理渲染任务优先级。

#### Scenario: 任务调度
- **WHEN** 状态更新触发渲染
- **THEN** 调度器根据优先级安排任务

#### Scenario: 时间切片
- **WHEN** 渲染任务执行时间超过切片
- **THEN** 让出控制权，等待下一次调度

### Requirement: 协调器实现 Diff 算法
协调器 SHALL 比较新旧 Fiber 树，生成更新副作用。

#### Scenario: 元素类型相同
- **WHEN** 新旧元素类型相同
- **THEN** 复用 Fiber 节点，更新属性

#### Scenario: 元素类型不同
- **WHEN** 新旧元素类型不同
- **THEN** 标记 Placement/Deletion 副作用

#### Scenario: 子元素列表
- **WHEN** 子元素列表变化
- **THEN** 使用 key 进行最小化更新

### Requirement: 协调器实现工作循环
协调器 SHALL 通过工作循环处理 Fiber 树。

#### Scenario: beginWork
- **WHEN** 处理 Fiber 节点
- **THEN** 执行 beginWork 处理子节点

#### Scenario: completeWork
- **WHEN** Fiber 节点没有子节点
- **THEN** 执行 completeWork 处理兄弟节点

#### Scenario: commitRoot
- **WHEN** 所有 Fiber 处理完成
- **THEN** 提交副作用到宿主

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Reconciler 架构                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   调度流程:                                                  │
│                                                             │
│   状态更新                                                   │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │  Scheduler  │  任务调度                                 │
│   └─────────────┘                                          │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │ Work Loop   │  工作循环                                 │
│   └─────────────┘                                          │
│       │                                                     │
│       ├─── beginWork() ───▶ 处理子节点                       │
│       │                                                     │
│       └─── completeWork() ─▶ 处理兄弟节点                    │
│                                                             │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │ commitRoot │  提交副作用                                 │
│   └─────────────┘                                          │
│       │                                                     │
│       ▼                                                     │
│   ┌─────────────┐                                          │
│   │ Host Config │  宿主配置                                  │
│   └─────────────┘                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Fiber Structure

```go
// Fiber 节点
type Fiber struct {
    // 类型
    Type        any      // 组件类型
    Key         string   // 列表 key
    
    // 树结构
    Child       *Fiber   // 第一个子节点
    Sibling     *Fiber   // 下一个兄弟节点
    Return      *Fiber   // 父节点
    
    // 双缓冲
    Alternate   *Fiber   // 对应的另一棵树节点
    
    // 状态
    PendingProps Props   // 待处理属性
    MemoizedProps Props  // 已记忆属性
    StateNode    any     // 宿主节点
    
    // 副作用
    Flags       FiberFlags // 副作用标记
    Deletions   []*Fiber   // 待删除节点
    
    // Hooks
    Hooks       []Hook   // Hook 列表
    HookIndex   int      // Hook 索引
}

// 副作用标记
type FiberFlags int

const (
    NoFlags       FiberFlags = 0
    Placement     FiberFlags = 1 << 0 // 插入
    Update        FiberFlags = 1 << 1 // 更新
    Deletion      FiberFlags = 1 << 2 // 删除
    ChildDeletion FiberFlags = 1 << 3 // 子节点删除
)
```

## Scheduler

```go
// 调度器
type Scheduler struct {
    // 任务队列
    taskQueue   []*Task
    currentTask *Task
    
    // 配置
    yieldInterval time.Duration
    
    // 状态
    isPaused bool
}

// 任务
type Task struct {
    Priority   int
    Callback   func() bool // 返回 true 表示完成
    Expiration time.Time
}

// 调度更新
func (s *Scheduler) scheduleUpdateOnFiber(fiber *Fiber) {
    // 创建任务
    task := &Task{
        Priority:   NormalPriority,
        Callback:   s.performConcurrentWorkOnRoot(fiber),
        Expiration: time.Now().Add(5 * time.Second),
    }
    
    s.taskQueue = append(s.taskQueue, task)
    s.requestWork()
}

// 工作循环
func (s *Scheduler) workLoop() {
    for {
        if s.currentTask == nil && len(s.taskQueue) > 0 {
            s.currentTask = s.taskQueue[0]
            s.taskQueue = s.taskQueue[1:]
        }
        
        if s.currentTask != nil {
            // 执行任务
            finished := s.currentTask.Callback()
            
            if finished {
                s.currentTask = nil
            } else {
                // 让出控制权
                time.Sleep(s.yieldInterval)
            }
        }
    }
}
```

## Diff Algorithm

```go
// 协调子节点
func reconcileChildren(current *Fiber, newChildren []Element) *Fiber {
    var firstChild *Fiber
    var prevSibling *Fiber
    
    oldFiber := current.Child
    newIdx := 0
    
    for newIdx < len(newChildren) || oldFiber != nil {
        var newFiber *Fiber
        
        // 比较 key
        if oldFiber != nil && newIdx < len(newChildren) {
            if oldFiber.Key == newChildren[newIdx].Key &&
               oldFiber.Type == newChildren[newIdx].Type {
                // 复用
                newFiber = useFiber(oldFiber, newChildren[newIdx].Props)
            } else {
                // 替换
                newFiber = createFiber(newChildren[newIdx])
                if oldFiber != nil {
                    oldFiber.Flags |= Deletion
                    current.Deletions = append(current.Deletions, oldFiber)
                }
            }
        } else if oldFiber != nil {
            // 删除
            oldFiber.Flags |= Deletion
            current.Deletions = append(current.Deletions, oldFiber)
        } else if newIdx < len(newChildren) {
            // 新增
            newFiber = createFiber(newChildren[newIdx])
        }
        
        // 建立兄弟关系
        if firstChild == nil {
            firstChild = newFiber
        } else {
            prevSibling.Sibling = newFiber
        }
        prevSibling = newFiber
        
        if oldFiber != nil {
            oldFiber = oldFiber.Sibling
        }
        newIdx++
    }
    
    return firstChild
}
```

## Host Config

```go
// 宿主配置 (HostConfig)
type HostConfig struct {
    // 创建实例
    CreateInstance  func(type string, props Props) Instance
    
    // 追加子节点
    AppendChild     func(parent, child Instance)
    
    // 插入子节点
    InsertBefore    func(parent, child, before Instance)
    
    // 移除子节点
    RemoveChild     func(parent, child Instance)
    
    // 更新实例
    UpdateInstance  func(instance Instance, props Props)
    
    // 提交更新
    CommitUpdate    func(instance Instance, oldProps, newProps Props)
}
```

## Module Structure

```
fiber/
├── fiber.go            # Fiber 结构
├── scheduler.go        # 调度器
├── work_loop.go        # 工作循环
├── begin_work.go       # beginWork
├── complete_work.go    # completeWork
├── commit_work.go      # 提交副作用
├── diff.go             # Diff 算法
├── host_config.go      # 宿主配置
└── flags.go            # 副作用标记
```

## React Ink 对齐

| 特性 | React Ink | Go-Ink | 说明 |
|-----|-----------|--------|------|
| Fiber 架构 | ✅ react-reconciler | ✅ | 可中断渲染 |
| 双缓冲 | ✅ | ✅ | current/workInProgress |
| 调度器 | ✅ | ✅ | 优先级调度 |
| 时间切片 | ✅ | ✅ | 让出控制权 |
| Diff 算法 | ✅ | ✅ | key 匹配 |
| 副作用标记 | ✅ | ✅ | Placement/Update/Deletion |
| Host Config | ✅ | ✅ | 宿主抽象 |