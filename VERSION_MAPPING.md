# Version Mapping

本文档记录 Go-React-Ink 与 React Ink 的版本对应关系，以及同步更新流程。

## 当前版本

| 项目 | 版本 | Commit Hash | 发布日期 |
|------|------|-------------|----------|
| **Go-React-Ink** | v0.1.0 | - | 2026-04-08 |
| **React Ink** | v6.8.0 | `be1b1bb6ec65056e2ed60ef3c5ae642704b82d31` | 2025-01 |
| **Go** | 1.22 | - | - |
| **tcell** | v2.7.4 | - | 2024-02 |

## Go 与 tcell 版本要求

### 当前配置

- **Go**: 1.22 (兼容性最佳)
- **tcell**: v2.7.4 (支持 Go 1.12+)

### tcell 版本与 Go 要求对照

| tcell 版本 | 最低 Go 版本 | 推荐 |
|-----------|-------------|------|
| v2.7.4 及以下 | go 1.12 | ✅ 当前使用 |
| v2.8.0 - v2.8.1 | go 1.12 | ✅ 可选升级 |
| v2.9.0 | go 1.23 | ⚠️ 需升级 Go |
| v2.10.0+ | go 1.24 | ⚠️ 需升级 Go |
| v2.13.8 (最新) | go 1.24 | ⚠️ 需升级 Go |

### 升级 tcell 的步骤

如需升级 tcell 到 v2.9.0+，需要同步升级项目 Go 版本：

1. **升级 Go 版本**
   ```bash
   # 方式1: 使用 go.mod 工具链指令
   # go.mod 中设置: go 1.24

   # 方式2: 安装 Go 1.24
   # Windows: https://go.dev/dl/
   # 或使用: go install golang.org/toolchain@go1.24.0
   ```

2. **更新 go.mod**
   ```go
   go 1.24

   require (
       github.com/gdamore/tcell/v2 v2.13.8
       // ...
   )
   ```

3. **更新依赖**
   ```bash
   go mod tidy
   go build ./...
   go test ./...
   ```

4. **更新 CI 配置**
   - `.github/workflows/*.yml` 中的 Go 版本
   - `VERSION_MAPPING.md` 版本记录

### 注意事项

- 项目通过 `pkg/tcell/tcell.go` 定义了 tcell 兼容的抽象接口
- tcell 作为可选依赖，升级不影响核心功能
- 升级 Go 版本需考虑用户兼容性

## React Ink 版本历史

| 版本 | Commit Hash | 发布日期 | 同步状态 |
|------|-------------|----------|----------|
| v6.8.0 | `be1b1bb6ec65056e2ed60ef3c5ae642704b82d31` | 2025-01 | ✅ 已同步 |
| v6.7.0 | `135cb23ae3b7ca94918b1cd913682f6356f12c5c` | - | ⬜ 待同步 |
| v6.6.0 | `43a913c41d13b4e68c2932adeeccf2b6406036fa` | - | ⬜ 待同步 |
| v6.5.1 | `45ed972ddfd3b508ed6397f9f844d076454a891c` | - | ⬜ 待同步 |
| v6.5.0 | `69813b4b5bdb17589d7c690262d7a474a58fbc21` | - | ⬜ 待同步 |

## 同步更新流程

当 React Ink 发布新版本时，按以下流程进行同步：

### 1. 检测新版本

```bash
# 检查 npm 上的最新版本
curl -s https://registry.npmjs.org/ink/latest | jq '.version'

# 检查 GitHub tags
curl -s https://api.github.com/repos/vadimdemedes/ink/tags | jq '.[0].name'
```

### 2. 分析变更

```bash
# 查看新版本的 commit
curl -s https://api.github.com/repos/vadimdemedes/ink/compare/{old_commit}...{new_commit}

# 查看 changelog
# https://github.com/vadimdemedes/ink/releases
```

### 3. 更新代码

1. 创建新的 feature branch
2. 同步 React Ink 的 API 变更
3. 更新测试用例
4. 更新文档

### 4. 更新版本映射

更新以下文件：
- `VERSION_MAPPING.md` - 版本对应表
- `CHANGELOG.md` - 变更日志
- `README.md` - 对标版本信息

## API 兼容性矩阵

### 组件

| React Ink Component | Version Added | Go-React-Ink Status |
|---------------------|---------------|---------------------|
| `<Box>` | v0.1.0 | ✅ 实现 |
| `<Text>` | v0.1.0 | ✅ 实现 |
| `<Spacer>` | v0.1.0 | ✅ 实现 |
| `<Newline>` | v0.1.0 | ✅ 实现 |
| `<Static>` | v0.1.0 | ✅ 实现 |
| `<Transform>` | v0.1.0 | ✅ 实现 |
| `<Fragment>` | v0.1.0 | ✅ 实现 |

### Hooks

| React Ink Hook | Version Added | Go-React-Ink Status |
|----------------|---------------|---------------------|
| `useState` | v0.1.0 | ✅ 实现 |
| `useEffect` | v0.1.0 | ✅ 实现 |
| `useLayoutEffect` | v0.1.0 | ✅ 实现 |
| `useRef` | v0.1.0 | ✅ 实现 |
| `useMemo` | v0.1.0 | ✅ 实现 |
| `useCallback` | v0.1.0 | ✅ 实现 |
| `useInput` | v0.1.0 | ✅ 实现 |
| `useApp` | v0.1.0 | ✅ 实现 |
| `useFocus` | v0.1.0 | ✅ 实现 |
| `useFocusManager` | v0.1.0 | ✅ 实现 |
| `useCursor` | v0.1.0 | ✅ 实现 |
| `useAnimation` | v0.1.0 | ✅ 实现 |
| `useStdout` | v0.1.0 | ✅ 实现 |
| `useStdin` | v0.1.0 | ✅ 实现 |
| `useStderr` | v0.1.0 | ✅ 实现 |
| `useWindowSize` | v0.1.0 | ✅ 实现 |
| `useBoxMetrics` | v0.1.0 | ✅ 实现 |
| `usePaste` | v0.1.0 | ✅ 实现 |
| `useIsScreenReaderEnabled` | v0.1.0 | ✅ 实现 |

### 布局属性

| React Ink Prop | Version Added | Go-React-Ink Status |
|----------------|---------------|---------------------|
| `flexDirection` | v0.1.0 | ✅ 实现 |
| `justifyContent` | v0.1.0 | ✅ 实现 |
| `alignItems` | v0.1.0 | ✅ 实现 |
| `flexGrow` | v0.1.0 | ✅ 实现 |
| `flexShrink` | v0.1.0 | ✅ 实现 |
| `flexBasis` | v0.1.0 | ✅ 实现 |
| `padding` | v0.1.0 | ✅ 实现 |
| `margin` | v0.1.0 | ✅ 实现 |
| `width` | v0.1.0 | ✅ 实现 |
| `height` | v0.1.0 | ✅ 实现 |

## 差异说明

### 不支持的特性

以下 React Ink 特性暂未实现：

| 特性 | 原因 | 计划版本 |
|------|------|----------|
| tcell 实际集成 | 已定义抽象接口，待集成实现 | v0.2.0 |

## 待办事项

| 任务 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| tcell 完整集成 | 中 | ⏳ 待办 | 需升级 Go 到 1.24，集成跨平台终端特性 |
| compiler 测试覆盖率提升 | 低 | ⏳ 待办 | 当前 48.6%，目标 70%+ |
| React Ink 新版本同步 | 低 | ⏳ 待办 | v6.7.0 ~ v6.5.0 待同步 |

### 实现差异

| 特性 | React Ink | Go-React-Ink | 说明 |
|------|-----------|--------------|------|
| 布局引擎 | yoga-layout (C++) | 纯 Go 实现 | 避免 CGO 依赖 |
| Fiber 优先级 | Lane 系统 | 简化 3 级 | 简化实现 |
| 时间切片 | 可配置 | 固定 16ms | 简化实现 |

## 相关链接

- [React Ink GitHub](https://github.com/vadimdemedes/ink)
- [React Ink NPM](https://www.npmjs.com/package/ink)
- [React Ink 文档](https://github.com/vadimdemedes/ink#readme)
- [React Ink Changelog](https://github.com/vadimdemedes/ink/releases)