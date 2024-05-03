package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() any {
		b := bytes.NewBuffer(make([]byte, 0, 1024))
		return b
	},
}

func NewBuffer() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

type PlainTextHandler struct {
	*slog.TextHandler
	output     io.Writer
	extraAttrs []slog.Attr
}

func (p PlainTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return p.TextHandler.Enabled(ctx, level)
}

func (p PlainTextHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := NewBuffer()
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.WriteString(r.Time.Format(time.RFC1123))
	buf.WriteByte(' ')
	buf.WriteString(r.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(r.Message)
	buf.WriteByte(' ')

	for _, attr := range p.extraAttrs {
		buf.WriteString(attr.Key)
		buf.WriteByte('=')
		buf.WriteString(attr.Value.String())
		buf.WriteByte(' ')
	}

	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool {
			buf.WriteString(attr.Key)
			buf.WriteByte('=')
			buf.WriteString(attr.Value.String())
			buf.WriteByte(' ')
			return true
		})
	}

	buf.WriteByte('\n')
	_, err := p.output.Write(buf.Bytes())
	return err
}

func (p PlainTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	p.extraAttrs = append(p.extraAttrs, attrs...)
	return p
}

func (p PlainTextHandler) WithGroup(name string) slog.Handler {
	return p.TextHandler.WithGroup(name)
}
