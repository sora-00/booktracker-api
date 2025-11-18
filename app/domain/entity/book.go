package entity

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	TotalPages int    `json:"total_pages"`
	Publisher  string `json:"publisher"`
}

func (b *Book) Rename(newTitle string) {
	b.Title = newTitle
}
