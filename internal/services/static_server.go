package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
)

var fileserver http.Handler
var server *http.Server

func StartServer() {
	fileserver = http.FileServer(http.Dir("./data/images"))
	http.Handle("/", fileserver)
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", config.Conf.FileServerPort),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
}

func StopServer() {
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
