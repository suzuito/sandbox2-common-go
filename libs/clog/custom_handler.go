package clog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

const TraceInfosKey = "traceInfos"

type CustomHandler struct {
	slog.Handler
}

func (t *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			terr, ok := a.Value.Any().(terrors.TraceableError)
			if ok {
				traceInfos := []string{}
				for _, st := range terr.StackTrace() {
					traceInfos = append(traceInfos, fmt.Sprintf("%s:%d", st.Filename, st.Line))
				}
				r.AddAttrs(slog.Attr{Key: TraceInfosKey, Value: slog.AnyValue(traceInfos)})
			}
			return false
		}
		return true
	})
	return t.Handler.Handle(ctx, r)
}
