package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type nextHandler struct{}

// Pull Gin's context from the request context
// and call the next item in the chain.
func (h *nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.Context().Value(h).(*ginState)
	defer func(r *http.Request) { state.ctx.Request = r }(state.ctx.Request)
	state.load(w, r)
	state.ctx.Next()
}

// Modify the Gin's context using raw Golang Http Context
func (m *ginState) load(w http.ResponseWriter, r *http.Request) {
	m.childCalled = true
	m.ctx.Request = r
	m.ctx.Writer = &wrappedResponseWriter{m.ctx.Writer, w}
}

// Keep the state of gin context
type ginState struct {
	ctx         *gin.Context
	childCalled bool
}

// WrapAll allow wrapping multiple http.Handler,
// returns a slice of gin.HandlerFunc
func WrapAll(hh []func(h http.Handler) http.Handler) []gin.HandlerFunc {
	functions := make([]gin.HandlerFunc, 0)
	for _, h := range hh {
		functions = append(functions, Wrap(h))
	}
	return functions
}

// Wrap takes the common HTTP middleware function signature,
// calls it to generate a handler,
// and wraps it into a Gin middleware handler.
// This is just a convenience wrapper around New.
func Wrap(f func(h http.Handler) http.Handler) gin.HandlerFunc {
	next, adapter := New()
	return adapter(f(next))
}

func New() (http.Handler, func(h http.Handler) gin.HandlerFunc) {
	next := new(nextHandler)
	return next, func(h http.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			state := &ginState{ctx: c}
			ctx := context.WithValue(c.Request.Context(), next, state)
			h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
			if !state.childCalled {
				c.Abort()
			}
		}
	}
}
