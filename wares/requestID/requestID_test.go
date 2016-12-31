package requestID

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

func TestRequestID_NonePresent(t *testing.T) {
	s := middleware.NewStack()

	var reqID string
	hn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v, _ := r.Context().Value("requestID").(string)
		reqID = v
	}

	s.Use(RequestID)

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if reqID == "" {
		t.Error("failed to set request id to new uuid")
	}
}

func TestRequestID_AlreadyPresent(t *testing.T) {
	s := middleware.NewStack()

	var originalReqID = "literally almost any string"
	var reqID string
	hn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v, _ := r.Context().Value("requestID").(string)
		reqID = v
	}

	s.Use(RequestID)

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	req.Header.Set("X-Request-ID", originalReqID)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if reqID != originalReqID {
		t.Error("failed to set request id to original header")
	}
}

func plainHandler(fn httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, httprouter.Params{})
	}
}
