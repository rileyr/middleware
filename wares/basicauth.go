package wares

import (
	"crypto/subtle"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func BasicAuth(u, p string) func(httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		if u == "" || p == "" {
			return fn
		}
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			user, pass, ok := r.BasicAuth()
			if !ok || !auth(user, u, pass, p) {
				w.Header().Set("WWW-Authenticate", `Basic`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				return
			}

			fn(w, r, params)
		}
	}
}

func auth(u1, u2, p1, p2 string) bool {
	return subtle.ConstantTimeCompare([]byte(u1), []byte(u2)) == 1 && subtle.ConstantTimeCompare([]byte(p1), []byte(p2)) == 1
}
