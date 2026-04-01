package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
)

var fileserver http.Handler
var server *http.Server

func startFileServer() {
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

func StartLoginServer() {
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", config.Conf.WebServerPort),
	}

	// Handle /login with query parameters
	http.HandleFunc("/instagram_login", func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		query := r.URL.Query()
		code := query.Get("code")

		logger.Log("Webserver: Received", code)

		resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", url.Values{
			"client_id":     {config.Conf.InstagramClientId},
			"client_secret": {config.Conf.InstagramClientSecret},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {config.Conf.InstagramLoginRedirectUrl},
			"code":          {code},
		})
		if err != nil {
			logger.Error(merrors.New(merrors.InstagramInvalidGetTokenFromCodeErrorCode, err.Error()))
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &resp)
		if err != nil {
			logger.Error(merrors.New(merrors.ParseJSONErrorCode, err.Error()))
		}

		merr := processCode(code)
		if merr != nil {
			logger.Error(merr)
		}
		stopWebServer()

	})
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
}

func stopWebServer() {
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
