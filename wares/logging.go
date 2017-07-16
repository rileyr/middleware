package wares

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Logging(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reqID, _ := r.Context().Value(RequestIDContextKey).(string)
		method := r.Method
		log.Printf("START %s - %s\n", method, reqID)
		fn(w, r, p)
		log.Printf("END %s - %s\n", method, reqID)
	}
}
