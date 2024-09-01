package logger

// 参考: 将 logrus 和 slog 结合起来使用
// https://github.com/thangchung/go-coffeeshop logger 实现

import (
	"context"
	"log/slog"

	"github.com/fzzp/gotk"
	"github.com/sirupsen/logrus"
)

type LogrusHandler struct {
	logger *logrus.Logger
}

func NewLogrusHandler(logger *logrus.Logger) *LogrusHandler {
	return &LogrusHandler{
		logger: logger,
	}
}

/*
LogrusHandler 实现slog Handler 下面四个接口，就能扩展slog了。
这样就能无论何时调用slog写log时，日志记录仍将通过 Logrus 处理，但我们不再需要依赖库中的 Logrus API。

type Handler interface {
    Enabled(Level) bool
    Handle(r Record) error
    WithAttrs(attrs []Attr) Handler
    WithGroup(name string) Handler
}
*/

func (h *LogrusHandler) Enabled(ctx context.Context, _ slog.Level) bool {
	// 支持所有级别logger
	return true
}

// Handle 所有日志都会经过这里，在这里交由 logrus 处理输出
func (h *LogrusHandler) Handle(ctx context.Context, rec slog.Record) error {
	requestID := ctx.Value(gotk.RequestIDCtxKey)
	if requestID != nil {
		rec.AddAttrs(slog.Any("request_id", requestID))
	}
	fields := make(map[string]interface{}, rec.NumAttrs())
	// fields["request_id"] = requestID
	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	entry := h.logger.WithFields(fields)

	switch rec.Level {
	case slog.LevelDebug:
		entry.Debug(rec.Message)
	case slog.LevelInfo.Level():
		entry.Info(rec.Message)
	case slog.LevelWarn:
		entry.Warn(rec.Message)
	case slog.LevelError:
		entry.Error(rec.Message)
	}

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// not implemented for brevity
	return h
}

func (h *LogrusHandler) WithGroup(name string) slog.Handler {
	// not implemented for brevity
	return h
}
