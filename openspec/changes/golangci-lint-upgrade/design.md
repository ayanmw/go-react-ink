---
change: golangci-lint-upgrade
---

# Design: 升级 golangci-lint 和增加测试覆盖率

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                  CI/CD Pipeline                      │
├─────────────────────────────────────────────────────┤
│  golangci-lint v1.64+ (go1.24 compatible)            │
│  ├── .golangci.yml (updated config)                  │
│  └── GitHub Actions workflow                         │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│                  Test Coverage                       │
├─────────────────────────────────────────────────────┤
│  pkg/fiber/fiber_test.go     (47% → 70%+)           │
│  internal/compiler/compiler_test.go (48% → 70%+)    │
│  pkg/input/input_test.go     (56% → 70%+)           │
└─────────────────────────────────────────────────────┘
```

## Technical Decisions

### 1. golangci-lint 版本选择

**选择**: v1.64.6+ (支持 go1.24)

**安装方式**:
```bash
# Windows (PowerShell)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 或使用 scoop
scoop install golangci-lint
```

### 2. 配置更新

`.golangci.yml` 需要更新:
- 移除废弃的 `run.skip-files` → `issues.exclude-files`
- 移除废弃的 `run.skip-dirs` → `issues.exclude-dirs`
- 移除废弃的 `output.format` → `output.formats`

### 3. 测试策略

**fiber 模块测试重点**:
- Scheduler 调度逻辑
- Fiber 节点创建/更新
- 工作循环 (workLoop)
- 提交阶段 (commitRoot)

**compiler 模块测试重点**:
- JSX → Go 转换
- 表达式处理
- Import 管理
- 错误处理

**input 模块测试重点**:
- Hook 注册/注销
- 事件处理
- 状态管理

## Implementation Approach

### Phase 1: golangci-lint 升级

1. 安装新版本 golangci-lint
2. 更新配置文件
3. 运行 lint 并修复问题
4. 更新 CI workflow

### Phase 2: 测试补充

按模块逐个增加测试:
1. 分析未覆盖代码
2. 编写测试用例
3. 验证覆盖率达标

## Dependencies

- Go 1.24.10 (已安装)
- golangci-lint v1.64.6+
- 现有测试框架 (testing package)