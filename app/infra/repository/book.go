package repository

import (
	"database/sql"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBook(db *sql.DB) repository.Book {
	return &BookRepository{DB: db}
}

func (r *BookRepository) FindAll() ([]entity.Book, error) {
	rows, err := r.DB.Query("SELECT id, title, author, total_pages, publisher FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]entity.Book, 0)
	for rows.Next() {
		var b entity.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.TotalPages, &b.Publisher); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) FindByID(id int) (*entity.Book, error) {
	var b entity.Book
	err := r.DB.QueryRow("SELECT id, title, author, total_pages, publisher FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.Author, &b.TotalPages, &b.Publisher)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepository) Save(book *entity.Book) error {
	return r.DB.QueryRow(
		"INSERT INTO books (title, author, total_pages, publisher) VALUES ($1, $2, $3, $4) RETURNING id",
		book.Title, book.Author, book.TotalPages, book.Publisher,
	).Scan(&book.ID)
}

func (r *BookRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}
