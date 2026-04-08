# Contributing to Go-React-Ink

感谢您有兴趣为 Go-React-Ink 做出贡献！

## 开发环境设置

### 前置要求

- Go 1.22 或更高版本
- Git
- Make (可选，用于使用 Makefile 命令)

### 克隆项目

```bash
git clone https://github.com/ayanmw/go-react-ink.git
cd go-react-ink
```

### 安装依赖

```bash
go mod download
```

### 验证安装

```bash
go build ./...
go test ./...
```

## 开发工作流

### 1. 创建分支

```bash
git checkout -b feature/your-feature-name
```

分支命名规范：
- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 2. 编写代码

#### 代码风格

- 遵循 [Go 代码规范](https://golang.org/doc/effective_go)
- 使用 `gofmt` 格式化代码
- 添加必要的注释

```bash
# 格式化代码
gofmt -s -w .

# 检查代码
go vet ./...
```

#### 测试

- 为新功能编写测试
- 确保测试覆盖率不低于 70%
- 使用表驱动测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./pkg/layout/... -v

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. 提交代码

#### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

类型：
- `feat` - 新功能
- `fix` - Bug 修复
- `docs` - 文档更新
- `style` - 代码格式（不影响功能）
- `refactor` - 代码重构
- `test` - 测试相关
- `chore` - 构建/工具相关

示例：
```
feat(layout): add flex-wrap support

Add support for flex-wrap property in layout engine.
Supports wrap, nowrap, and wrap-reverse values.

Closes #123
```

#### 提交前检查

```bash
# 格式化
gofmt -s -w .

# 静态检查
go vet ./...

# 运行测试
go test ./...

# 构建
go build ./...
```

### 4. 推送并创建 PR

```bash
git push origin feature/your-feature-name
```

在 GitHub 上创建 Pull Request。

## Pull Request 指南

### PR 标题

使用与提交信息相同的格式：
```
feat(component): add new component
```

### PR 描述

```markdown
## 变更说明
简要描述变更内容

## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 其他

## 测试
- [ ] 已添加测试
- [ ] 所有测试通过

## 相关 Issue
Closes #xxx
```

### 代码审查

- 响应审查意见
- 及时更新代码
- 保持讨论专业和礼貌

## 项目结构

```
go-react-ink/
├── cmd/                    # 命令行工具
│   └── gox/               # gox 编译器
├── pkg/                    # 公共包
│   ├── components/        # UI 组件
│   ├── core/              # 核心类型
│   ├── fiber/             # Fiber 协调器
│   ├── hooks/             # React Hooks
│   ├── hostconfig/        # Host 配置
│   ├── input/             # 输入处理
│   ├── layout/            # 布局引擎
│   ├── renderer/          # 渲染器
│   └── tcell/             # 终端抽象
├── internal/              # 内部包
│   ├── codegen/           # 代码生成
│   ├── compiler/          # 编译器
│   ├── lexer/             # 词法分析
│   ├── parser/            # 语法分析
│   └── scanner/           # 扫描器
├── examples/              # 示例程序
├── test/                  # 测试用例
├── tools/                 # 工具
│   └── goland-gox/        # GoLand 插件
├── doc/                   # 文档
└── openspec/              # OpenSpec 变更管理
```

## 添加新组件

1. 在 `pkg/components/` 创建组件文件
2. 在 `pkg/core/` 定义组件类型（如需要）
3. 添加组件测试
4. 更新文档
5. 添加示例

示例：

```go
// pkg/components/mycomponent.go
package components

import "github.com/ayanmw/go-react-ink/pkg/core"

// MyComponent 创建自定义组件
func MyComponent(props core.Props, children []core.Element) core.Element {
    // 组件实现
    return &MyComponentElement{
        Props:    props,
        Children: children,
    }
}
```

## 添加新 Hook

1. 在 `pkg/hooks/` 或 `pkg/input/` 添加 Hook
2. 遵循 React Hooks 规范
3. 添加 Hook 测试
4. 更新文档

示例：

```go
// pkg/hooks/myhook.go
func UseMyHook(ctx *HookContext, initial any) (any, func(any)) {
    idx := ctx.nextHook()
    
    ctx.mu.Lock()
    var hook *StateHook
    if idx >= len(ctx.hooks) {
        hook = &StateHook{
            Value:   initial,
            context: ctx,
        }
        hook.Setter = func(newValue any) {
            hook.Value = newValue
            ctx.scheduleUpdate()
        }
        ctx.hooks = append(ctx.hooks, hook)
    } else {
        hook = ctx.hooks[idx].(*StateHook)
    }
    ctx.mu.Unlock()
    
    return hook.Value, hook.Setter
}
```

## 发布流程

1. 更新 `VERSION_MAPPING.md` 版本信息
2. 更新 `CHANGELOG.md`
3. 创建 Git tag
4. 推送 tag 到 GitHub
5. GitHub Actions 自动构建发布

## 获取帮助

- 提交 Issue: https://github.com/ayanmw/go-react-ink/issues
- 查看 Wiki: https://github.com/ayanmw/go-react-ink/wiki
- 阅读文档: `doc/` 目录

## 许可证

贡献的代码将采用 MIT 许可证。