package log

import (
	"context"
	"log/slog"
)

type ctxKey string

const slogFields ctxKey = "slog_fields"

// ContextHandler is a struct that wraps a slog.Handler, allowing additional context to be added to log records.
type ContextHandler struct {
	Handler slog.Handler
}

// Enabled checks if the underlying slog Handler is enabled for a given log level.
func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// WithAttrs creates a new ContextHandler with additional attributes.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	return &ContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

// WithGroup creates a new ContextHandler with a group name.
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

// Handle adds contextual attributes to the Record before calling the underlying
// Handler
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds a slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, args ...any) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	attr := argsToAttrSlice(args)

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr...)
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr...)
	return context.WithValue(parent, slogFields, v)
}

const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}

func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}
