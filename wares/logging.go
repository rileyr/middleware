package wares

import (
	"log"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/julienschmidt/httprouter"
)

func Logging(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reqID, _ := r.Context().Value(RequestIDContextKey).(string)
		var code int

		log.Printf("START id=%s method=%s url=%s\n", reqID, r.Method, r.URL)
		defer func(ts time.Time) {
			duration := time.Since(ts).String()
			if code == 0 {
				code = 500
			}

			log.Printf("END id=%s method=%s code=%d duration=%s\n", reqID, r.Method, code, duration)
		}(time.Now())

		wrapper := func(w2 http.ResponseWriter, r2 *http.Request) {
			fn(w2, r2, p)
		}

		m := httpsnoop.CaptureMetrics(http.HandlerFunc(wrapper), w, r)
		code = m.Code
	}
}
