# Tasks: Go-Ink Implementation

## Phase 1: 编译器

### 1.1 Scanner (源码扫描)
- [x] 实现 JSX 区域识别
- [x] 处理嵌套括号
- [x] 生成代码段列表
- [x] 单元测试

### 1.2 Lexer (词法分析)
- [x] Token 类型定义
- [x] 状态机实现
- [x] 处理字符串/表达式
- [x] 单元测试

### 1.3 Parser (语法分析)
- [ ] AST 节点定义
- [ ] 递归下降解析器
- [ ] 属性解析
- [ ] 子节点解析
- [ ] 错误恢复
- [ ] 单元测试

### 1.4 Code Generator
- [ ] 元素生成
- [ ] 属性生成
- [ ] 子节点生成
- [ ] 表达式转换
- [ ] 导入管理
- [ ] 单元测试

### 1.5 CLI Tool
- [ ] 命令行参数解析
- [ ] 文件编译
- [ ] 目录编译
- [ ] 监听模式
- [ ] 错误报告

## Phase 2: 核心运行时

### 2.1 Fiber 协调器
- [ ] Fiber 节点定义
- [ ] 调度器实现
- [ ] 工作循环
- [ ] beginWork/completeWork
- [ ] Diff 算法
- [ ] 提交阶段
- [ ] 单元测试

### 2.2 Flexbox 布局 (CGO)
- [ ] yoga-layout CGO 绑定
- [ ] 布局节点包装
- [ ] 样式应用
- [ ] 布局计算
- [ ] 单元测试

### 2.3 终端渲染器
- [ ] 输出缓冲
- [ ] 单元格样式
- [ ] Diff 计算
- [ ] ANSI 生成
- [ ] tcell 集成
- [ ] 单元测试

### 2.4 Host Config
- [ ] 创建实例
- [ ] 更新实例
- [ ] 树操作
- [ ] 布局调度
- [ ] 渲染调度

## Phase 3: 组件与 Hooks

### 3.1 基础组件
- [ ] Box 组件
- [ ] Text 组件
- [ ] Spacer 组件
- [ ] Newline 组件

### 3.2 高级组件
- [ ] Static 组件
- [ ] Transform 组件

### 3.3 核心 Hooks
- [ ] useState
- [ ] useEffect
- [ ] useRef
- [ ] useMemo

### 3.4 专用 Hooks
- [ ] useInput
- [ ] useApp
- [ ] useFocus
- [ ] useFocusManager
- [ ] useCursor
- [ ] useAnimation

## Phase 4: Flexbox Go 实现

### 4.1 核心算法
- [ ] 布局节点定义
- [ ] 测量阶段
- [ ] Flex 解析
- [ ] 子节点定位
- [ ] 对齐计算

### 4.2 文本处理
- [ ] 文本测量
- [ ] 宽字符处理
- [ ] 自动换行

### 4.3 验证
- [ ] 替换 CGO
- [ ] 测试通过
- [ ] 性能验证

## Phase 5: 工具链

### 5.1 VSCode 插件
- [ ] package.json
- [ ] 语法高亮 (tmLanguage)
- [ ] 自动编译
- [ ] 错误提示

### 5.2 LSP Server
- [ ] 服务器实现
- [ ] 诊断
- [ ] 补全
- [ ] 悬停信息

## Phase 6: 测试与文档

### 6.1 测试
- [ ] 移植 React Ink 测试用例
- [ ] 功能对比测试
- [ ] 生成对比报告

### 6.2 文档
- [ ] README
- [ ] API 文档
- [ ] 示例代码

### 6.3 发布
- [ ] 版本号确定
- [ ] CHANGELOG
- [ ] GitHub Release

## Progress Tracking

| Phase | 任务数 | 完成 | 进度 |
|-------|--------|------|------|
| Phase 1 | 20 | 8 | 40% |
| Phase 2 | 25 | 0 | 0% |
| Phase 3 | 15 | 0 | 0% |
| Phase 4 | 10 | 0 | 0% |
| Phase 5 | 10 | 0 | 0% |
| Phase 6 | 10 | 0 | 0% |
| **Total** | **90** | **8** | **8.9%** |
