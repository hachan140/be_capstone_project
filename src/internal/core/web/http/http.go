package http

import (
	"be-capstone-project/src/cmd/public/config"
	"fmt"
	"net/http"
)

func NewHttpServer(handler http.Handler, config *config.Config) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%v", config.App.Port),
		Handler: handler,
	}
}
