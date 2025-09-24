package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/model"
)

func GetBlueskyFeed() (*model.BskyFeedResp, error) {
	url := fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=5&filter=posts_with_media", config.Conf.BskyHandle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var feed model.BskyFeedResp
	err = json.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
