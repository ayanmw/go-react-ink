# Change: Go-Ink Implementation

## Summary

使用 Golang 完整复刻 React Ink 框架，实现声明式终端 UI 开发体验。

## Status

- [x] Design Complete
- [ ] Implementation In Progress
- [ ] Testing Complete
- [ ] Released

## Scope

### In Scope

1. **编译器 (gox-compiler)**
   - JSX 解析器 (Scanner/Lexer/Parser)
   - 表达式转换 (JS → Go)
   - 代码生成器
   - CLI 工具链

2. **运行时 (go-ink)**
   - Fiber 协调器
   - Flexbox 布局引擎
   - 终端渲染器
   - 内置组件 (Box/Text/Spacer/Newline/Static/Transform)
   - Hooks (useState/useEffect/useInput/useApp/useFocus...)

3. **工具链**
   - VSCode 插件
   - LSP Server
   - 语法高亮

4. **测试与验证**
   - 功能对比测试 (React Ink 测试用例)
   - API 兼容性验证

### Out of Scope

- React DevTools 集成 (后续版本)
- 服务端渲染 (SSR)
- 原生移动端支持

## Design Decisions

| 决策 | 选择 | 理由 |
|-----|------|------|
| JSX 支持 | 预处理器 (.gox → .go) | 零运行时开销，编译时类型检查 |
| 时间切片 | 固定 16ms | 简单可靠，与 React 一致 |
| Lane 系统 | 简化版 3 级 | 终端场景足够，降低复杂度 |
| Flexbox | 倒序实现：CGO → Go子集 → 完整 | 降低风险，快速验证 |
| 终端抽象 | tcell | 跨平台兼容性 |
| 协调器 | Fiber 架构 | 与 React 一致，支持可中断渲染 |

## Implementation Phases

### Phase 1: 编译器 (预计 2 周)

- [ ] Scanner (源码扫描)
- [ ] Lexer (词法分析)
- [ ] Parser (语法分析)
- [ ] Code Generator
- [ ] CLI Tool

### Phase 2: 核心运行时 (预计 3 周)

- [ ] Fiber 协调器
- [ ] Flexbox 布局 (CGO + yoga-layout)
- [ ] 终端渲染器
- [ ] Host Config

### Phase 3: 组件与 Hooks (预计 2 周)

- [ ] Box/Text 组件
- [ ] Spacer/Newline 组件
- [ ] Static/Transform 组件
- [ ] useState/useEffect
- [ ] useInput/useApp
- [ ] useFocus/useCursor

### Phase 4: Flexbox Go 实现 (预计 2 周)

- [ ] 替换 CGO 为纯 Go 实现
- [ ] 核心属性支持
- [ ] 测试验证

### Phase 5: 工具链 (预计 1 周)

- [ ] VSCode 插件
- [ ] LSP Server
- [ ] 语法高亮

### Phase 6: 测试与文档 (预计 1 周)

- [ ] 测试用例移植
- [ ] 功能对比报告
- [ ] API 文档

## Success Criteria

1. **功能完整性**: React Ink 所有测试用例通过
2. **API 兼容性**: 组件/Hooks API 与 React Ink 一致
3. **跨平台**: Windows/macOS/Linux 支持
4. **文档完整**: README + API 文档 + 示例

## Dependencies

### 编译时
- Go 1.21+
- golang.org/x/tools (代码格式化)

### 运行时
- github.com/gdamore/tcell/v2 (终端抽象)
- github.com/facebook/yoga (Flexbox, CGO, Phase 1)
- golang.org/x/sys (系统调用)

### 开发时
- github.com/fsnotify/fsnotify (文件监听)

## Risks

| 风险 | 影响 | 缓解措施 |
|-----|------|---------|
| CGO 跨平台编译 | 高 | Phase 2 后替换为纯 Go |
| Fiber 复杂度 | 中 | 简化 Lane 系统 |
| 性能不达标 | 低 | 增量渲染优化 |

## References

- [React Ink](https://github.com/vadimdemedes/ink) - 原始框架
- [React Reconciler](https://github.com/facebook/react/tree/main/packages/react-reconciler)
- [Yoga Layout](https://github.com/facebook/yoga)
- [Templ](https://github.com/a-h/templ) - Go 模板引擎参考
