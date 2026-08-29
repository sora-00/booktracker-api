package response

type BookThumbnailPost struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func NewBookThumbnailPost(id string, url string) *BookThumbnailPost {
	return &BookThumbnailPost{ID: id, URL: url}
}

type BookThumbnailGet struct {
	ContentType string
	Body        []byte
}

func NewBookThumbnailGet(contentType string, body []byte) *BookThumbnailGet {
	return &BookThumbnailGet{ContentType: contentType, Body: body}
}
