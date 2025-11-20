package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CreateBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func NewCreateBook(r *http.Request) (*CreateBook, error) {
	req := &CreateBook{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Author == "" {
		return nil, errors.New("author is required")
	}
	return req, nil
}
