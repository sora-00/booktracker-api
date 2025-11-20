package response

import (
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

func FromBookEntity(book *entity.Book) Book {
	return Book{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
		CreatedAt: book.CreatedAt,
	}
}

func FromBookEntities(books []entity.Book) []Book {
	responses := make([]Book, 0, len(books))
	for i := range books {
		responses = append(responses, FromBookEntity(&books[i]))
	}
	return responses
}
