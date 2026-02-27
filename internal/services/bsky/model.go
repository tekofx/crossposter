package bsky

type BlueskyClient struct {
	Handle   string
	Password string
	JWT      string
	DID      string
}

type BlobResponse struct {
	Blob Blob `json:"blob"`
}

type Post struct {
	Type      string       `json:"$type"`
	Text      string       `json:"text"`
	CreatedAt string       `json:"createdAt"`
	Embed     *EmbedImages `json:"embed,omitempty"`
}

type EmbedImages struct {
	Type   string      `json:"$type"`
	Images []ImageItem `json:"images"`
}

type ImageItem struct {
	Alt   string `json:"alt"`
	Image Blob   `json:"image"`
}

type Blob struct {
	Type     string `json:"$type"`
	Ref      Ref    `json:"ref"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

type Ref struct {
	Link string `json:"$link"`
}

type PostRequest struct {
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

// PostRecord represents the actual post content (the "record")
type PostRecord struct {
	Type      string       `json:"$type"`
	Text      string       `json:"text"`
	CreatedAt string       `json:"createdAt"`
	Embed     *EmbedImages `json:"embed,omitempty"` // Optional: only if embedding images
}

type PublishResponse struct {
	Uri string `json:"uri"`
}
