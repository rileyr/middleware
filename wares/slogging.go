package wares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/julienschmidt/httprouter"
)

func Slogger(l *slog.Logger) func(httprouter.Handle) httprouter.Handle {
	if l == nil {
		l = slog.Default()
	}

	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			reqID, _ := r.Context().Value(RequestIDContextKey{}).(string)
			var code int

			l.Info(
				"START",
				"id", reqID,
				"method", r.Method,
				"url", r.URL,
			)

			defer func(ts time.Time) {
				dur := time.Since(ts).String()
				if code == 0 {
					code = 500
				}

				l.Info(
					"END",
					"id", reqID,
					"method", r.Method,
					"url", r.URL,
					"code", code,
					"duration", dur,
				)
			}(time.Now())

			wrapper := func(w2 http.ResponseWriter, r2 *http.Request) {
				fn(w2, r2, p)
			}

			m := httpsnoop.CaptureMetrics(http.HandlerFunc(wrapper), w, r)
			code = m.Code
		}
	}
}
