package postgres

import (
	"context"
	"database/sql"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

// BookRepo の PostgreSQL 実装（domain のインターフェースに依存するだけ）
type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) repository.BookRepo {
	return &bookRepo{db: db}
}

const (
	bookColumns     = `"id", "title", "author", "totalPages", "publisher", "thumbnailUrl", "status", "targetCompleteDate", "createdAt", "updatedAt"`
	queryInsertBook = `
		INSERT INTO books (title, author, totalPages, publisher, thumbnailUrl, status, targetCompleteDate)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
	return scanner.Scan(&book.ID, &book.Title, &book.Author, &book.TotalPages, &book.Publisher, &book.ThumbnailUrl, &book.Status, &book.TargetCompleteDate, &book.CreatedAt, &book.UpdatedAt)
}

func (r *bookRepo) Create(ctx context.Context, book *entity.Book) error {
	return scanBookInto(r.db.QueryRowContext(
		ctx,
		queryInsertBook,
		book.Title, book.Author, book.TotalPages, book.Publisher, book.ThumbnailUrl, book.Status, book.TargetCompleteDate,
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
	book, err := scanBook(r.db.QueryRowContext(ctx, queryFindBookByID, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return book, nil
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, queryDeleteBook, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
