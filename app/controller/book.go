package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/usecase"
	"github.com/sora-00/booktracker-api/app/usecase/request"
	"github.com/sora-00/booktracker-api/app/usecase/response"
)

type BookController struct {
	Usecase *usecase.BookUsecase
}

func NewController(u *usecase.BookUsecase) *BookController {
	return &BookController{Usecase: u}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := c.Usecase.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewBookGet(books))
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGetByID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book, err := c.Usecase.GetBookByID(r.Context(), req.BookID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewBookGetByID(book))
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookCreate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := c.Usecase.CreateBook(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewBookCreate(book))
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookDelete(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Usecase.DeleteBook(r.Context(), req.BookID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
