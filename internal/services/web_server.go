package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services/socials/instagram"
)

var fileserver http.Handler
var server *http.Server

func StartFileServer() {
	fileserver = http.FileServer(http.Dir("./data/images"))
	http.Handle("/", fileserver)
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", config.Conf.WebServerPort),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
}

func StartInstagramLoginServer() {
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", config.Conf.WebServerPort),
	}

	// Handle /login with query parameters
	http.HandleFunc("/instagram_login", func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		query := r.URL.Query()
		code := query.Get("code")

		logger.Log("Webserver: Received", code)
		err := instagram.GetTokenFromCode(code)
		if err != nil {
			logger.Error(err)
		}

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
