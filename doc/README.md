# Go-Ink 设计文档索引

> 使用 Golang 复刻 React Ink 完整框架实现

## 文档列表

| 文档 | 内容 | 状态 |
|-----|------|------|
| [01-overview.md](./01-overview.md) | 项目概览、架构设计、核心决策 | ✅ 完成 |
| [02-jsx-parser.md](./02-jsx-parser.md) | JSX 解析器设计 (Scanner/Lexer/Parser) | ✅ 完成 |
| [03-expr-transform.md](./03-expr-transform.md) | 表达式转换规则 (JS → Go) | ✅ 完成 |
| [04-codegen.md](./04-codegen.md) | 代码生成器设计 (AST → Go 代码) | ✅ 完成 |
| [05-toolchain.md](./05-toolchain.md) | 工具链集成 (CLI/IDE/LSP) | ✅ 完成 |
| [06-reconciler.md](./06-reconciler.md) | React 协调器设计 (Fiber 架构) | ✅ 完成 |
| [07-layout.md](./07-layout.md) | Flexbox 布局引擎 | ✅ 完成 |
| [08-renderer.md](./08-renderer.md) | 终端渲染器 (增量渲染) | ✅ 完成 |
| [09-components.md](./09-components.md) | 内置组件 (Box/Text/Spacer 等) | ✅ 完成 |
| [10-hooks.md](./10-hooks.md) | React Hooks 复刻 | ✅ 完成 |

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
    │   ├── 依赖 core
    │   └── 依赖 layout
    ├── components (内置组件)
    │   └── 依赖 core
    └── hooks (状态管理)
        └── 依赖 core, reconciler
```

---

## 开发优先级

### Phase 1: 编译器 (已完成设计)
1. ✅ JSX 解析器
2. ✅ 表达式转换
3. ✅ 代码生成器
4. ✅ CLI 工具

### Phase 2: 核心运行时 (已完成设计)
1. ✅ Reconciler (Fiber 架构)
2. ✅ Flexbox 布局引擎
3. ✅ 终端渲染器

### Phase 3: 组件与 API (已完成设计)
1. ✅ 内置组件
2. ✅ Hooks API
3. 🔄 事件处理 (待细化)

### Phase 4: 工具链完善
1. 🔜 VSCode 插件
2. 🔜 LSP Server
3. 🔜 调试支持

### Phase 5: 实现与测试
1. 🔜 编译器实现
2. 🔜 运行时实现
3. 🔜 测试用例对比

---

## 关键设计决策

| 决策 | 选择 | 理由 |
|-----|------|------|
| JSX 支持 | 预处理器 | 零运行时开销，编译时类型检查 |
| Flexbox | Go 核心子集实现 | 终端场景可简化 30-40% |
| 渲染策略 | 增量渲染 | 复刻 Ink 核心魔法 |
| 终端抽象 | tcell | 跨平台兼容性 |
| 协调器 | Fiber 架构 | 与 React 一致 |

---

## 原版参考

- `reference/react-ink/` - React Ink 原版代码及测试用例

## 参考项目

- [React Ink](https://github.com/vadimdemedes/ink) - 原始框架
- [React Reconciler](https://github.com/facebook/react/tree/main/packages/react-reconciler) - 协调器
- [Yoga Layout](https://github.com/facebook/yoga) - Flexbox 引擎
- [Templ](https://github.com/a-h/templ) - Go 模板引擎参考
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Go TUI 参考