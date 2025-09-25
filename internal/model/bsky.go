package model

import (
	"time"
)

type BskyFeedResp struct {
	Posts []BskyPost `json:"feed"`
}

type BskyPost struct {
	Post struct {
		Uri       string
		Author    BskyAuthor
		Embed     BskyEmbed
		Record    BskyRecord `json:"record"`
		CreatedAt time.Time  `json:"createdAt"`
	}
	Reason *BskyReason `json:"reason"`
	Reply  *BskyReply  `json:"reply,omitempty"`
}

func (post *BskyPost) IsQuote() bool {
	if post.Reason == nil {
		return false
	}

	return post.Reason.Type == "app.bsky.embed.record#view"
}

func (post *BskyPost) IsReply() bool {
	return post.Reply != nil
}

func (post *BskyPost) IsRepost() bool {
	if post.Reason == nil {
		return false
	}
	return post.Reason.Type == "app.bsky.feed.defs#reasonRepost"
}

type BskyReply struct {
	Root struct {
		Cid string
		Uri string
	}
}

type BskyReason struct {
	/* Type of post. Can be:
	 * - Repost: app.bsky.feed.defs#reasonRepost
	 */
	Type string `json:"$type"`
	Uri  string `json:"uri"`
}

type BskyRecord struct {
	Text string `json:"text"`
}

type BskyAuthor struct {
	did    string
	handle string
}

type BskyEmbed struct {
	/* Type of embed. Can be:
	 * - Quote: app.bsky.embed.record#view
	 */
	Type   string `json:"$type"`
	Images []BskyImage
}

type BskyImage struct {
	Thumb    string
	Fullsize string
}
