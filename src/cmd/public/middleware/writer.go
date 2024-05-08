package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// A wrapper that turns a http.ResponseWriter into a gin.ResponseWriter, given an existing gin.ResponseWriter
// Needed if the middleware you are using modifies the writer it passes downstream
// Wrap more methods: https://golang.org/pkg/net/http/#ResponseWriter
type wrappedResponseWriter struct {
	gin.ResponseWriter
	writer http.ResponseWriter
}

func (w *wrappedResponseWriter) Writer() http.ResponseWriter {
	return w.writer
}

func (w *wrappedResponseWriter) WriteString(s string) (n int, err error) {
	return w.writer.Write([]byte(s))
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}
