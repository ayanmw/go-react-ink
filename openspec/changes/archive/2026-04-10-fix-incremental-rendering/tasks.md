# Tasks

- [x] 添加 IncrementalRendering 字段到 RenderOptions
- [x] 修改 NewLogUpdate() 签名接受 incremental 参数
- [x] 修改 NewInk() 传递 IncrementalRendering 到 NewLogUpdate
- [x] 更新 DefaultRenderOptions() 设置默认值
- [x] 更新测试中的 NewLogUpdate 调用
- [x] 创建 pkg/hooks/use_app.go
- [x] 添加 AppContext 字段到 HookContext
- [x] 添加 SetAppContext 方法到 HookContext
- [x] 在 Ink.Render() 中创建并设置 AppContext
- [x] 实现 Ink.ScheduleUpdate() 作为 Dispatcher
- [x] 创建 pkg/input/input_parser.go
- [x] 创建 pkg/hooks/use_input.go
- [x] 在 Ink 中添加 stdin 输入处理
- [x] 简化 animation-demo 示例
