package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) (_ bool) {
	return false
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) (_ error) {
	return nil
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) (_ slog.Handler) {
	return h
}

func (h *DiscardHandler) WithGroup(name string) (_ slog.Handler) {
	return h
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
