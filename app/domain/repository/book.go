package repository

import "github.com/sora-00/booktracker-api/app/domain/entity"

type Book interface {
	FindAll() ([]entity.Book, error)
	FindByID(id int) (*entity.Book, error)
	Save(book *entity.Book) error
	Delete(id int) error
}
