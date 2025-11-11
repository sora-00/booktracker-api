package request

type CreateBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
