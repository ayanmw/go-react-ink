# Go-Ink 开发文档

> 使用 Golang 复刻 React Ink 完整框架实现

## 文档结构

- **用户文档**: 项目根目录 `README.md`
- **开发文档**: `openspec/specs/` 目录下的规格文档

## Specs 规格文档

| Spec | 内容 | React Ink 对齐 |
|-----|------|---------------|
| [compiler](./specs/compiler/spec.md) | JSX 编译器 (Scanner/Lexer/Parser/CodeGen) | ✅ 完全对齐 |
| [renderer](./specs/renderer/spec.md) | 终端渲染器 (增量渲染/TTY/信号/动画) | ✅ 完全对齐 |
| [reconciler](./specs/reconciler/spec.md) | Fiber 协调器 (调度/Diff/工作循环) | ✅ 完全对齐 |
| [layout](./specs/layout/spec.md) | Flexbox 布局引擎 | ✅ 完全对齐 |
| [components](./specs/components/spec.md) | 内置组件 (Box/Text/Static 等) | ✅ 完全对齐 |
| [hooks](./specs/hooks/spec.md) | React Hooks (useState/useEffect 等) | ✅ 完全对齐 |

---

## 核心架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Go-Ink 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   .gox 文件 (JSX) ──▶ gox-compiler ──▶ .go 文件             │
│                                                             │
│   ┌─────────────────────────────────────────────────────┐  │
│   │                   运行时架构                          │  │
│   ├─────────────────────────────────────────────────────┤  │
│   │                                                      │  │
│   │   用户组件 ──▶ Reconciler ──▶ Host Config            │  │
│   │                                │                     │  │
│   │                     ┌──────────┴──────────┐          │  │
│   │                     ▼                     ▼          │  │
│   │               Flexbox Layout      Terminal Renderer   │  │
│   │                     │                     │          │  │
│   │                     └──────────┬──────────┘          │  │
│   │                                ▼                     │  │
│   │                          ANSI Output                  │  │
│   │                                                      │  │
│   └─────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## React Ink 核心技术对齐

### 编译器 (compiler)
| 特性 | React Ink | Go-Ink |
|-----|-----------|--------|
| JSX 语法 | ✅ Babel | ✅ gox-compiler |
| 嵌套元素 | ✅ | ✅ |
| 表达式 `{}` | ✅ | ✅ |
| 属性展开 | ✅ | ✅ |
| 条件渲染 | ✅ | ✅ |
| 列表渲染 | ✅ | ✅ |

### 渲染器 (renderer)
| 特性 | React Ink | Go-Ink |
|-----|-----------|--------|
| 增量渲染 | ✅ log-update | ✅ LogUpdate |
| 备用屏幕缓冲 | ✅ | ✅ |
| TTY 检测 | ✅ | ✅ |
| 信号处理 | ✅ | ✅ |
| FPS 节流 | ✅ | ✅ |
| 动画系统 | ✅ | ✅ |

### 协调器 (reconciler)
| 特性 | React Ink | Go-Ink |
|-----|-----------|--------|
| Fiber 架构 | ✅ react-reconciler | ✅ |
| 双缓冲 | ✅ | ✅ |
| 调度器 | ✅ | ✅ |
| Diff 算法 | ✅ | ✅ |

### 布局 (layout)
| 特性 | React Ink | Go-Ink |
|-----|-----------|--------|
| Flexbox | ✅ Yoga | ✅ |
| flexDirection | ✅ | ✅ |
| justifyContent | ✅ | ✅ |
| alignItems | ✅ | ✅ |
| flexGrow/shrink | ✅ | ✅ |

### 组件 (components)
| 组件 | React Ink | Go-Ink |
|-----|-----------|--------|
| Box | ✅ | ✅ |
| Text | ✅ | ✅ |
| Static | ✅ | ✅ |
| Newline | ✅ | ✅ |
| Spacer | ✅ | ✅ |
| Transform | ✅ | ✅ |

### Hooks
| Hook | React Ink | Go-Ink |
|-----|-----------|--------|
| useState | ✅ | ✅ |
| useEffect | ✅ | ✅ |
| useInput | ✅ | ✅ |
| useApp | ✅ | ✅ |
| useFocus | ✅ | ✅ |
| useAnimation | ✅ | ✅ |
| useMemo | ✅ | ✅ |
| useRef | ✅ | ✅ |
| useContext | ✅ | ✅ |

---

## 模块依赖关系

```
gox-compiler (预编译时)
    ├── scanner (源码扫描)
    ├── lexer (词法分析)
    ├── parser (语法分析)
    └── codegen (代码生成)

go-ink (运行时)
    ├── core (核心类型)
    ├── reconciler (协调器)
    │   └── 依赖 core, hooks
    ├── layout (Flexbox)
    │   └── 依赖 core
    ├── renderer (渲染器)
    │   └── 依赖 core, layout
    ├── components (内置组件)
    │   └── 依赖 core
    └── hooks (状态管理)
        └── 依赖 core, reconciler
```

---

## 开发状态

| 模块 | 设计 | 实现 | 测试 |
|-----|------|------|------|
| gox-compiler | ✅ | ✅ | ✅ |
| reconciler | ✅ | ✅ | ✅ |
| layout | ✅ | ✅ | ✅ |
| renderer | ✅ | ✅ | ✅ |
| components | ✅ | ✅ | ✅ |
| hooks | ✅ | ✅ | ✅ |

---

## 参考项目

- [React Ink](https://github.com/vadimdemedes/ink) - 原始框架
- [React Reconciler](https://github.com/facebook/react/tree/main/packages/react-reconciler) - 协调器
- [Yoga Layout](https://github.com/facebook/yoga) - Flexbox 引擎