package clog

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

type LoggerType uint8

const (
	LoggerTypeJSON LoggerType = iota + 1
	LoggerTypeDevSlog
	LoggerTypeE2E
	LoggerTypeText
)

func NewLogger(
	level slog.Level,
	loggerType LoggerType,
) *slog.Logger {
	var slogHandler slog.Handler
	switch loggerType {
	case LoggerTypeDevSlog:
		slogHandler = devslog.NewHandler(os.Stdout, &devslog.Options{
			HandlerOptions: &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		})
	case LoggerTypeJSON:
		slogHandler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.LevelKey {
						a.Key = "severity"
					}
					return a
				},
			},
		)
	case LoggerTypeE2E:
		slogHandler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

					switch a.Key {
					case slog.LevelKey:
						a.Key = "severity"
					case slog.TimeKey, slog.SourceKey, "traceInfos":
						a = slog.Attr{}
					}

					return a
				},
			},
		)
	case LoggerTypeText:
		fallthrough
	default:
		slogHandler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		)
	}

	slogCustomHandler := CustomHandler{
		Handler: slogHandler,
	}

	return slog.New(&slogCustomHandler)
}

func SetDefaultLogger(
	level slog.Level,
	loggerType LoggerType,
) {
	logger := NewLogger(level, loggerType)
	slog.SetDefault(logger)
}
