package entity

import "time"

type Book struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	Author    string      `json:"author"`
	CreatedAt time.Time   `json:"createdAt"`
}

func (b *Book) Rename(newTitle string) {
	b.Title = newTitle
}
