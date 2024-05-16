package middleware

import (
	"be-capstone-project/src/internal/core/web/context"
	assert "github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type dummyTestRequestContextHandler struct {
	writer         http.ResponseWriter
	request        *http.Request
	responseStatus int
}

func (d *dummyTestRequestContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(d.responseStatus)
	d.writer = w
	d.request = r
}

type mockResponseWriter struct {
}

func (d mockResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (d mockResponseWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (d mockResponseWriter) WriteHeader(statusCode int) {
}

func TestAdvancedResponseWriter_ShouldReplaceDefaultWriter(t *testing.T) {
	next := dummyTestRequestContextHandler{responseStatus: http.StatusOK}
	handler := AdvancedResponseWriter()
	assert.NotNil(t, handler)

	internalHandler := handler(&next)
	assert.NotNil(t, internalHandler)

	handlerFunc, ok := internalHandler.(http.HandlerFunc)
	assert.True(t, ok)

	r, _ := http.NewRequest("GET", "/test", nil)
	handlerFunc(&mockResponseWriter{}, r)

	assert.IsType(t, &context.AdvancedResponseWriter{}, next.writer)
}
