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
			var code int

			defer func(ts time.Time) {
				dur := time.Since(ts)
				if code == 0 {
					code = 500
				}

				path := r.URL.Path
				if r.URL.RawQuery != "" {
					path += "?" + r.URL.RawQuery
				}

				l.Info(
					"request",
					"method", r.Method,
					"path", path,
					"status", code,
					"duration", dur.String(),
				)
			}(time.Now())

			// Fix: ensure params is never nil/empty for routes without path parameters
			if p == nil || len(p) == 0 {
				p = httprouter.Params{{Key: "", Value: ""}}
			}

			wrapper := func(w2 http.ResponseWriter, r2 *http.Request) {
				fn(w2, r2, p)
			}

			m := httpsnoop.CaptureMetrics(http.HandlerFunc(wrapper), w, r)
			code = m.Code
		}
	}
}
