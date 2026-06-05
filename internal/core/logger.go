package core

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// PrettyHandler is a zero-dependency slog.Handler that writes
// human-readable colored log lines to the given writer.
//
// Output format:
//
//	12:04:05  INFO   server starting          addr=:8000
//	12:04:05  WARN   module disabled          module=legacy
//	12:04:05  ERROR  bootstrap failed         error=...
type PrettyHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
}

func NewPrettyHandler(level slog.Level) *PrettyHandler {
	return &PrettyHandler{w: os.Stdout, level: level}
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGrey   = "\033[90m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// Time
	ts := colorGrey + r.Time.Format("15:04:05") + colorReset

	// Level
	var levelStr string
	switch {
	case r.Level >= slog.LevelError:
		levelStr = colorRed + colorBold + "ERROR" + colorReset
	case r.Level >= slog.LevelWarn:
		levelStr = colorYellow + colorBold + "WARN " + colorReset
	case r.Level >= slog.LevelInfo:
		levelStr = colorGreen + colorBold + "INFO " + colorReset
	default:
		levelStr = colorGrey + colorBold + "DEBUG" + colorReset
	}

	// Message
	msg := colorBold + fmt.Sprintf("%-35s", r.Message) + colorReset

	// Attributes
	var attrs []string
	for _, a := range h.attrs {
		attrs = append(attrs, colorCyan+a.Key+colorReset+"="+fmt.Sprintf("%v", a.Value.Any()))
	}
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, colorCyan+a.Key+colorReset+"="+fmt.Sprintf("%v", a.Value.Any()))
		return true
	})

	attrStr := ""
	if len(attrs) > 0 {
		attrStr = "  " + strings.Join(attrs, "  ")
	}

	fmt.Fprintf(h.w, "%s  %s  %s%s\n", ts, levelStr, msg, attrStr)
	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		w:     h.w,
		level: h.level,
		attrs: append(h.attrs, attrs...),
	}
}

func (h *PrettyHandler) WithGroup(_ string) slog.Handler { return h }
