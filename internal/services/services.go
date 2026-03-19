package services

import (
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services/socials/bsky"
)

func Initialize() {
	if config.Conf.BskyEnabled {
		if err := bsky.Initialize(); err != nil {
			logger.Fatal(err)
		}
	}
}
