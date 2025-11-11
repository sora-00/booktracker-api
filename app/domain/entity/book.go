package entity

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (b *Book) Rename(newTitle string) {
	b.Title = newTitle
}
