package request

type CreateBook struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	TotalPages int    `json:"total_pages"`
	Publisher  string `json:"publisher"`
}
