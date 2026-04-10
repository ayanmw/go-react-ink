package hooks

// InputHandler 输入处理函数签名
type InputHandler func(input string, key InputKey)

// InputKey 输入键（避免与 input 包的 Key 类型冲突）
type InputKey struct {
	UpArrow    bool
	DownArrow  bool
	LeftArrow  bool
	RightArrow bool
	PageDown   bool
	PageUp     bool
	Home       bool
	End        bool
	Return     bool
	Escape     bool
	Ctrl       bool
	Shift      bool
	Tab        bool
	Backspace  bool
	Delete     bool
	Meta       bool
}

type UseInputOptions struct {
	IsActive bool
}

// UseInputHook 键盘输入 Hook
type UseInputHook struct {
	appContext *AppContext
	handler    InputHandler
	isActive   bool
	dispatcher Dispatcher
}

func (h *UseInputHook) Cleanup() {
	if h.appContext != nil && h.isActive {
		// 清理 stdin raw mode (待实现)
	}
}

// UseInput 键盘输入 Hook
func UseInput(ctx *HookContext, handler InputHandler, opts *UseInputOptions) {
	isActive := true
	if opts != nil {
		isActive = opts.IsActive
	}

	appCtx, err := UseApp(ctx)
	if err != nil || appCtx == nil {
		return // 非 Ink 上下文中无法使用
	}

	hook := &UseInputHook{
		appContext: appCtx,
		handler:    handler,
		isActive:   isActive,
		dispatcher: ctx.dispatcher,
	}

	idx := ctx.nextHook()

	ctx.mu.Lock()
	if idx >= len(ctx.hooks) || ctx.hooks[idx] == nil {
		ctx.hooks = append(ctx.hooks, hook)
	} else {
		// 复用已有 Hook
		existingHook := ctx.hooks[idx].(*UseInputHook)
		existingHook.handler = handler
		existingHook.isActive = isActive
	}
	ctx.mu.Unlock()
}
