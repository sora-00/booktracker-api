package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/usecase"
	"github.com/sora-00/booktracker-api/app/usecase/request"
)

type BookController struct {
	Book *usecase.Book
}

func NewBookController(b *usecase.Book) *BookController {
	return &BookController{Book: b}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGet(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := c.Book.Get(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGetByID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := c.Book.GetByID(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookCreate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := c.Book.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookUpdate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := c.Book.Update(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookDelete(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = c.Book.Delete(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
