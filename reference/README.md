# React Ink 原版参考

> 本目录包含 React Ink 原版代码及测试用例，用于实现对比验证

## 目录结构

```
reference/react-ink/
├── src/                    # 源代码
│   ├── index.ts           # 入口 (导出所有 API)
│   ├── reconciler.ts      # React Reconciler 配置
│   ├── render.ts          # 渲染入口
│   ├── renderer.ts        # 渲染器
│   ├── dom.ts             # DOM 操作
│   ├── output.ts          # 输出处理
│   ├── components/        # 内置组件
│   │   ├── Box.ts
│   │   ├── Text.ts
│   │   ├── Static.ts
│   │   ├── Transform.ts
│   │   ├── Newline.ts
│   │   ├── Spacer.ts
│   │   └── *Context.ts
│   └── hooks/             # Hooks
│       ├── use-input.ts
│       ├── use-app.ts
│       ├── use-stdin.ts
│       ├── use-stdout.ts
│       ├── use-focus.ts
│       ├── use-cursor.ts
│       └── use-animation.ts
│
├── test/                   # 测试用例 (60+ 文件)
│   ├── flex*.tsx          # Flexbox 布局测试
│   ├── hooks*.tsx         # Hooks 测试
│   ├── components.tsx     # 组件测试
│   ├── render.tsx         # 渲染测试
│   ├── reconciler.tsx     # 协调器测试
│   └── fixtures/          # 测试 fixtures
│
├── examples/               # 示例 (30+ 个)
│   ├── counter/           # 计数器
│   ├── chat/              # 聊天 UI
│   ├── table/             # 表格
│   ├── use-input/         # 输入处理
│   └── ...
│
└── benchmark/              # 性能基准测试
    ├── simple/            # 简单渲染
    └── static/            # 静态内容
```

## 核心依赖

| 依赖 | 版本 | 用途 |
|-----|------|------|
| react | ^19.2.4 | React 核心 |
| react-reconciler | ^0.33.0 | 协调器 |
| yoga-layout | ~3.2.1 | Flexbox 布局引擎 |
| scheduler | ^0.27.0 | 调度器 |

## API 导出

### 组件
- `Box` - Flexbox 容器
- `Text` - 文本组件
- `Static` - 静态内容
- `Transform` - 变换组件
- `Newline` - 换行
- `Spacer` - 占位

### Hooks
- `useInput` - 键盘输入
- `usePaste` - 粘贴事件
- `useApp` - 应用实例
- `useStdin` - 标准输入
- `useStdout` - 标准输出
- `useStderr` - 标准错误
- `useFocus` - 焦点管理
- `useFocusManager` - 焦点管理器
- `useCursor` - 光标控制
- `useAnimation` - 动画
- `useWindowSize` - 窗口尺寸
- `useBoxMetrics` - Box 尺寸

### API
- `render` - 渲染入口
- `renderToString` - 渲染为字符串
- `measureElement` - 测量元素

## 测试用例分类

### 布局测试
- `flex.tsx` - Flex 基础
- `flex-direction.tsx` - 方向
- `flex-wrap.tsx` - 换行
- `flex-justify-content.tsx` - 主轴对齐
- `flex-align-items.tsx` - 交叉轴对齐
- `flex-align-self.tsx` - 自身对齐
- `flex-align-content.tsx` - 多行对齐
- `gap.tsx` - 间距
- `padding.tsx` - 内边距
- `margin.tsx` - 外边距
- `overflow.tsx` - 溢出处理

### 组件测试
- `components.tsx` - 组件基础
- `text.tsx` - 文本组件
- `borders.tsx` - 边框
- `background.tsx` - 背景
- `display.tsx` - 显示模式
- `position.tsx` - 定位

### Hooks 测试
- `hooks.tsx` - Hooks 基础
- `hooks-use-input.tsx` - 输入 Hook
- `hooks-use-paste.tsx` - 粘贴 Hook
- `focus.tsx` - 焦点

### 渲染测试
- `render.tsx` - 渲染基础
- `render-to-string.tsx` - 字符串渲染
- `reconciler.tsx` - 协调器
- `log-update.tsx` - 日志更新
- `terminal-resize.tsx` - 终端尺寸变化

### 特性测试
- `exit.tsx` - 退出处理
- `errors.tsx` - 错误处理
- `cursor.tsx` - 光标
- `screen-reader.tsx` - 屏幕阅读器

## 运行原版测试

```bash
cd reference/react-ink
npm install
npm test
```

## 运行示例

```bash
cd reference/react-ink
npm run example examples/counter/counter.tsx
```

## 实现对比计划

### Phase 1: 功能对比
1. 运行所有测试用例
2. 记录通过/失败情况
3. 分析失败原因

### Phase 2: 性能对比
1. 运行 benchmark
2. 对比渲染时间
3. 对比内存使用
4. 生成对比报告

### Phase 3: 兼容性对比
1. API 兼容性检查
2. 属性支持对比
3. 边界情况处理

## 对比报告模板

```markdown
# Go-Ink vs React Ink 对比报告

## 功能对比

| 功能 | React Ink | Go-Ink | 状态 |
|-----|-----------|--------|------|
| Box 组件 | ✅ | ✅ | 完成 |
| Text 组件 | ✅ | ✅ | 完成 |
| ... | ... | ... | ... |

## 测试通过率

- React Ink: 100% (baseline)
- Go-Ink: XX%

## 性能对比

| 指标 | React Ink | Go-Ink | 提升 |
|-----|-----------|--------|------|
| 渲染时间 | X ms | Y ms | Z% |
| 内存使用 | X MB | Y MB | Z% |

## 差异说明

1. ...
2. ...
```
