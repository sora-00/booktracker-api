package controller

import (
	"net/http"

	"github.com/sora-00/booktracker-api/app/controller/respond"
	httpx "github.com/sora-00/booktracker-api/app/infra/http"
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
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.Get(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *BookController) GetBooksByStatus(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGetByStatus(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.GetByStatus(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *BookController) GetLogsByBookStatus(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGetByStatus(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.GetLogsByBookStatus(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookGetByID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.GetByID(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "book not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookCreate(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.Create(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusCreated, res)
}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookUpdate(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Book.Update(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "book not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookDelete(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = c.Book.Delete(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "book not found")
		respond.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
