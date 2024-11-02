package pretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

const timeFormat = "[15:05:05.000]"

type PrettyHandler struct {
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(w, opts),
		l:       stdLog.New(w, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format(timeFormat)

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.GreenString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	h.l.Println(
		color.BlackString(timeStr),
		level,
		color.CyanString(r.Message),
		string(b),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
