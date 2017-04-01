## middleware

middleware library for [github.com/julienschmidt/httprouter](https://www.github.com/julienschmidt/httprouter).

--

**Usage**

A middleware is any function that takes and returns a `httprouter.Handle`:

```go
type Middleware func(httprouter.Handle) httprouter.Handle
```

A basic passthrough middleware:

```go
package main

import(
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/rileyr/middleware"
  "github.com/rileyr/middleware/wares"
)

// Define the middleware function:
func passThrough(fn httprouter.Handle) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
      fn(w, r, p)
  }
}

// Define some handler:
func handler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  w.Write([]byte("hello world!"))
}

func main() {
  router := httprouter.New()

  // Create a middleware stack:
  s := middleware.NewStack()

  // Use your middleware:
  s.Use(passThrough)
  s.Use(wares.RequestID)
  s.Use(wares.Logging)

  // Wrap Handlers:
  router.GET("/", s.Wrap(handler))

  http.ListenAndServe(":3000", router)
}
```
