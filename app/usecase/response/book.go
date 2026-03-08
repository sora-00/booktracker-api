package response

import (
	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/service"
)

// BookListItem は一覧用。API が計算した remainingDays を含む。
type BookListItem struct {
	*entity.Book
	RemainingDays int `json:"remainingDays"`
}

type BookGet struct {
	Books []*BookListItem `json:"books"`
}

func NewBookGet(books []entity.Book) *BookGet {
	items := make([]*BookListItem, 0, len(books))
	for i := range books {
		b := &books[i]
		items = append(items, &BookListItem{
			Book:          b,
			RemainingDays: service.RemainingDays(b.TargetCompleteDate),
		})
	}
	return &BookGet{Books: items}
}

// BookGetByID は 1 件取得用。remainingDays を含む。
type BookGetByID struct {
	*entity.Book
	RemainingDays int `json:"remainingDays"`
}

func NewBookGetByID(book *entity.Book) *BookGetByID {
	return &BookGetByID{
		Book:          book,
		RemainingDays: service.RemainingDays(book.TargetCompleteDate),
	}
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
