package model

type BskyFeedResp struct {
	Posts []BskyPost `json:"feed"`
}

type BskyPost struct {
	Post struct {
		Uri    string
		Author BskyAuthor
		Embed  BskyEmbed
	}
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
	fullsize string
}
