package hooks

import "errors"

// AppContext 应用上下文，由 Ink 实例提供
type AppContext struct {
	Exit                  func(...any)
	WaitUntilRenderFlush func() error
}

type AppContextKey struct{}

// UseApp 获取应用上下文
func UseApp(ctx *HookContext) (*AppContext, error) {
	if ctx.appContext == nil {
		return nil, ErrNotInInkContext
	}
	return ctx.appContext, nil
}

// ErrNotInInkContext 非 Ink 上下文错误
var ErrNotInInkContext = errors.New("UseApp must be used within an Ink app context")
