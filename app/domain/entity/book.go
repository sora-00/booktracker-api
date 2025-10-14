package entity

type Book struct {
	ID     int
	Title  string
	Author string
}

func (b *Book) Rename(newTitle string) {
	b.Title = newTitle
}
