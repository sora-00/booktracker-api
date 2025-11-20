package repository

import (
	"context"
	"database/sql"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type BookRepo interface {
	Create(ctx context.Context, book *entity.Book) error
	FindAll(ctx context.Context) ([]entity.Book, error)
	FindByID(ctx context.Context, id int) (*entity.Book, error)
	Delete(ctx context.Context, id int) error
}

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepo {
	return &bookRepo{db: db}
}

const (
	bookColumns     = `id, title, author, "createdAt"`
	queryInsertBook = `
		INSERT INTO books (title, author)
		VALUES ($1, $2)
		RETURNING ` + bookColumns + `
`
	queryFindAllBooks = `
		SELECT ` + bookColumns + `
		FROM books
`
	queryFindBookByID = `
		SELECT ` + bookColumns + `
		FROM books
		WHERE id = $1
`
	queryDeleteBook = `
		DELETE FROM books
		WHERE id = $1
`
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanBook(scanner rowScanner) (*entity.Book, error) {
	var b entity.Book
	if err := scanBookInto(scanner, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func scanBookInto(scanner rowScanner, book *entity.Book) error {
	return scanner.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt)
}

func (r *bookRepo) Create(ctx context.Context, book *entity.Book) error {
	return scanBookInto(r.db.QueryRowContext(
		ctx,
		queryInsertBook,
		book.Title, book.Author,
	), book)
}

func (r *bookRepo) FindAll(ctx context.Context) ([]entity.Book, error) {
	rows, err := r.db.QueryContext(ctx, queryFindAllBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]entity.Book, 0)
	for rows.Next() {
		var book entity.Book
		if err := scanBookInto(rows, &book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepo) FindByID(ctx context.Context, id int) (*entity.Book, error) {
	return scanBook(r.db.QueryRowContext(
		ctx,
		queryFindBookByID,
		id,
	))
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queryDeleteBook, id)
	return err
}
