package response

import "github.com/sora-00/booktracker-api/app/domain/entity"

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	TotalPages int    `json:"total_pages"`
	Publisher  string `json:"publisher"`
}

func FromBookEntity(book *entity.Book) Book {
	return Book{
		ID:         book.ID,
		Title:      book.Title,
		Author:     book.Author,
		TotalPages: book.TotalPages,
		Publisher:  book.Publisher,
	}
}

func FromBookEntities(books []entity.Book) []Book {
	responses := make([]Book, 0, len(books))
	for i := range books {
		responses = append(responses, FromBookEntity(&books[i]))
	}
	return responses
}
