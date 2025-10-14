package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sora-00/booktracker-api/app/usecase"
)

type BookController struct {
	Usecase *usecase.BookUsecase
}

func NewController(u *usecase.BookUsecase) *BookController {
	return &BookController{Usecase: u}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := c.Usecase.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	book, err := c.Usecase.GetBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := c.Usecase.CreateBook(req.Title, req.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	if err := c.Usecase.DeleteBook(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}