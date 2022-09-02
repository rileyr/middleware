package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
)

func TestUse(t *testing.T) {
	s := NewStack()
	mw := func(fn httprouter.Handle) httprouter.Handle {
		return fn
	}
	c := len(s.middlewares)

	s.Use(mw)

	if len(s.middlewares) != c+1 {
		t.Error("expected Use() to increase the number of items in the stack")
	}
}

func TestWrap(t *testing.T) {
	s := NewStack()

	var middlewareCalled bool
	mw := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			middlewareCalled = true
			fn(w, r, p)
		}
	}
	s.Use(mw)

	var handlerCalled bool
	hn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		handlerCalled = true
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("expected handler to have been called")
	}

	if !middlewareCalled {
		t.Error("expected middleware to have been called")
	}
}

func TestWrap_Ordering(t *testing.T) {
	s := NewStack()

	var handlerCallAt *time.Time

	first, firstCallAt := timingHandler()
	second, secondCallAt := timingHandler()
	third, thirdCallAt := timingHandler()
	fourth, fourthCallAt := timingHandler()

	s.Use(first)
	s.Use(second)
	s.Use(third)
	s.Use(fourth)

	hn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ts := time.Now()
		handlerCallAt = &ts
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(plainHandler(wrapped))
	handler.ServeHTTP(w, req)

	if firstCallAt.IsZero() || secondCallAt.IsZero() || thirdCallAt.IsZero() || fourthCallAt.IsZero() || handlerCallAt.IsZero() {
		t.Fatal("failed to call one or more functions")
	}

	if firstCallAt.After(*secondCallAt) || firstCallAt.After(*thirdCallAt) || firstCallAt.After(*fourthCallAt) || firstCallAt.After(*handlerCallAt) {
		t.Error("failed to call first middleware first")
	}

	if fourthCallAt.Before(*thirdCallAt) || fourthCallAt.Before(*secondCallAt) || fourthCallAt.After(*handlerCallAt) {
		t.Error("failed to call fourth middleware last before the handler")
	}

	if secondCallAt.After(*thirdCallAt) {
		t.Error("expected second middleware to come before the third")
	}
}

func timingHandler() (func(fn httprouter.Handle) httprouter.Handle, *time.Time) {
	tmp := time.Time{}
	var t *time.Time
	t = &tmp

	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ts := time.Now()
			*t = ts
			fn(w, r, p)
		}
	}, t
}

func TestWrap_WhenEmpty(t *testing.T) {
	s := NewStack()
	hn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}
	w := s.Wrap(hn)

	if reflect.ValueOf(hn).Pointer() != reflect.ValueOf(w).Pointer() {
		t.Error("expected that Wrap() would return the given function when stack is empty")
	}
}

func plainHandler(fn httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, httprouter.Params{})
	}
}
