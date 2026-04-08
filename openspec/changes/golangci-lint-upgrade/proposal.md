---
change: golangci-lint-upgrade
status: proposed
created: 2026-04-08
---

# Change: 升级 golangci-lint 和增加测试覆盖率

## Summary

升级 golangci-lint 到 go1.24 兼容版本，并增加 fiber/compiler 模块的测试覆盖率，提升项目代码质量。

## Status

- [ ] Design Complete
- [ ] Implementation In Progress
- [ ] Testing Complete
- [ ] Released

## Scope

### In Scope

1. **golangci-lint 升级**
   - 安装 go1.24 兼容的 golangci-lint 版本
   - 更新 .golangci.yml 配置
   - 修复 lint 报告的问题

2. **测试覆盖率提升**
   - fiber 模块测试 (当前 47.2% → 目标 70%+)
   - compiler 模块测试 (当前 48.6% → 目标 70%+)
   - input 模块测试 (当前 56.2% → 目标 70%+)

3. **CI/CD 配置更新**
   - 更新 GitHub Actions 中的 golangci-lint 版本

### Out of Scope

- React Ink 新版本同步 (另有任务)
- 核心架构改动

## Motivation

当前 golangci-lint 使用 go1.23 编译，无法分析 go1.24.10 项目。测试覆盖率偏低，影响代码质量信心。

## Risks

| 风险 | 影响 | 缓解措施 |
|-----|------|---------|
| golangci-lint 新版本配置变化 | 可能需要调整配置 | 查阅官方迁移文档 |
| 新增测试可能发现现有 bug | 需要修复 | 预留修复时间 |

## Success Criteria

- [ ] golangci-lint 运行无错误
- [ ] fiber 测试覆盖率 ≥ 70%
- [ ] compiler 测试覆盖率 ≥ 70%
- [ ] 所有测试通过