package servicebus

import "context"

type (
	// Handler exposes the functionality required to process a Service Bus message.
	Handler interface {
		Handle(context.Context, *Message) DispositionAction
	}

	// HandlerFunc is a type converter that allows a func to be used as a `Handler`
	HandlerFunc func(context.Context, *Message) DispositionAction
)

// Handle redirects this call to the func that was provided.
func (hf HandlerFunc) Handle(ctx context.Context, msg *Message) DispositionAction {
	return hf(ctx, msg)
}
