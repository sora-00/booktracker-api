package request

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const maxThumbnailSize = 10 << 20

type BookThumbnailPost struct {
	Filename string
	Data     []byte
}

func NewBookThumbnailPost(req *http.Request) (*BookThumbnailPost, error) {
	if err := req.ParseMultipartForm(maxThumbnailSize); err != nil {
		return nil, errors.New("failed to parse multipart form")
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		return nil, errors.New("file is required")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("failed to read file")
	}
	return &BookThumbnailPost{Filename: header.Filename, Data: data}, nil
}

type BookThumbnailGet struct {
	ID string
}

func NewBookThumbnailGet(req *http.Request) (*BookThumbnailGet, error) {
	id := strings.TrimSpace(chi.URLParam(req, "id"))
	if id == "" || strings.Contains(id, "/") || strings.Contains(id, "..") {
		return nil, errors.New("not found")
	}
	return &BookThumbnailGet{ID: id}, nil
}
