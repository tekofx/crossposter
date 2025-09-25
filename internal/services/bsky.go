package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/model"
)

func GetBlueskyPosts() ([]model.Post, error) {
	url := fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=5&filter=posts_no_replies", config.Conf.BskyHandle)
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

	var posts []model.Post
	for _, bskyPost := range feed.Posts {
		var images []string
		if len(bskyPost.Post.Embed.Images) > 0 {
			for _, image := range bskyPost.Post.Embed.Images {
				images = append(images, image.Fullsize)
			}
		}

		posts = append(posts, model.Post{
			BskyId:   bskyPost.Post.Uri,
			Text:     bskyPost.Post.Record.Text,
			Images:   images,
			Date:     bskyPost.Post.CreatedAt,
			IsQuote:  bskyPost.IsQuote(),
			IsRepost: bskyPost.IsRepost(),
			IsReply:  bskyPost.IsReply(),
		})

	}

	return posts, nil
}
