package response

import (
	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type BookGet struct {
	Books []*entity.Book `json:"books"`
}

func NewBookGet(books []entity.Book) *BookGet {
	bs := make([]*entity.Book, 0, len(books))
	for i := range books {
		bs = append(bs, &books[i])
	}
	return &BookGet{Books: bs}
}

type BookGetByID struct {
	*entity.Book
}

func NewBookGetByID(book *entity.Book) *BookGetByID {
	return &BookGetByID{book}
}

type BookCreate struct {
	*entity.Book
}

func NewBookCreate(book *entity.Book) *BookCreate {
	return &BookCreate{book}
}

type BookUpdate struct {
	*entity.Book
}

func NewBookUpdate(book *entity.Book) *BookUpdate {
	return &BookUpdate{book}
}

type BookDelete struct {
	BookID int `json:"bookId"`
}

func NewBookDelete(bookID int) *BookDelete {
	return &BookDelete{BookID: bookID}
}
