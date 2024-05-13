package lib

import "maunium.net/go/mautrix/event"

// MContext represents the context for handling messages.
type MContext struct {
	Request *event.MessageEventContent // Request contains the message event content
	Params  map[string]any             // Params contains the parameters from the request
	Result  chan any                   // Result is a channel for the result
}

// ContextFunc is a function type for context modification.
type ContextFunc func(*MContext)

// WithParams returns a ContextFunc that sets the parameters of the context.
func WithParams(params map[string]any) ContextFunc {
	return func(ctx *MContext) {
		ctx.Params = make(map[string]any, len(params))

		for k, v := range params {
			ctx.Params[k] = v
		}
	}
}

// NewMContext creates a new context with the given options.
func NewMContext(Opts ...ContextFunc) *MContext {
	ctx := MContext{
		Params: make(map[string]any),
		Result: make(chan any),
	}

	for _, opt := range Opts {
		opt(&ctx)
	}
	return &ctx
}
