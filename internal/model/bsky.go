package model

type BskyFeedResp struct {
	Posts []BskyPost `json:"feed"`
}

type BskyPost struct {
	Post struct {
		Uri    string
		Author BskyAuthor
		Embed  BskyEmbed
		Record BskyRecord `json:"record"`
	}
	Reason *BskyReason `json:"reason,omitempty"`
}

type BskyReason struct {
	Type string `json:"$type"`
}

type BskyRecord struct {
	Text string `json:"text"`
}

type BskyAuthor struct {
	did    string
	handle string
}

type BskyEmbed struct {
	Images []BskyImage
}

type BskyImage struct {
	Thumb    string
	Fullsize string
}
