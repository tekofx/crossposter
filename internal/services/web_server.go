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

func StartFileServer() {
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

func StartInstagramLoginServer() {
	fileserver = http.FileServer(http.Dir("./data/images"))
	http.Handle("/", fileserver)
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", config.Conf.FileServerPort),
	}

	// Handle /login with query parameters
	http.HandleFunc("/instagram_login", func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		query := r.URL.Query()
		code := query.Get("code")

		logger.Log("Webserver: Received", code)
	})
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
}

func StopWebServer() {
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
