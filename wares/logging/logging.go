package logging

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware/wares/requestID"
)

func Logging(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reqID, _ := r.Context().Value(requestID.ContextKey).(string)
		log.Printf("START %s\n", reqID)
		fn(w, r, p)
		log.Printf("END %s\n", reqID)
	}
}
