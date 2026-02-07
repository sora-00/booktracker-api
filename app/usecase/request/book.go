package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type BookGet struct{}

func NewBookGet(req *http.Request) (*BookGet, error) {
	return &BookGet{}, nil
}

type BookGetByID struct {
	BookID int `json:"bookId"`
}

func NewBookGetByID(req *http.Request) (*BookGetByID, error) {
	idStr := chi.URLParam(req, "id")
	if idStr == "" {
		return nil, errors.New("book id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid book id")
	}
	return &BookGetByID{BookID: id}, nil
}

type BookCreate struct {
	BookCreateForm
}

func NewBookCreate(req *http.Request) (*BookCreate, error) {
	var raw struct {
		Title              string `json:"title"`
		Author             string `json:"author"`
		TotalPages         int    `json:"totalPages"`
		Publisher          string `json:"publisher"`
		ThumbnailUrl       string `json:"thumbnailUrl"`
		Status             string `json:"status"`
		TargetCompleteDate string `json:"targetCompleteDate"`
	}
	if err := json.NewDecoder(req.Body).Decode(&raw); err != nil {
		return nil, err
	}
	targetDate, err := parseTargetCompleteDate(raw.TargetCompleteDate)
	if err != nil {
		return nil, err
	}
	r := &BookCreate{
		BookCreateForm: BookCreateForm{
			Title:              raw.Title,
			Author:             raw.Author,
			TotalPages:         raw.TotalPages,
			Publisher:          raw.Publisher,
			ThumbnailUrl:       raw.ThumbnailUrl,
			Status:             raw.Status,
			TargetCompleteDate: targetDate,
		},
	}
	if err := r.ValidateBookCreateForm(); err != nil {
		return nil, err
	}
	return r, nil
}

type BookDelete struct {
	BookID int `json:"bookId"`
}

func NewBookDelete(req *http.Request) (*BookDelete, error) {
	idStr := chi.URLParam(req, "id")
	if idStr == "" {
		return nil, errors.New("book id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid book id")
	}
	return &BookDelete{BookID: id}, nil
}

// ---

type BookCreateForm struct {
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	TotalPages         int       `json:"totalPages"`
	Publisher          string    `json:"publisher"`
	ThumbnailUrl       string    `json:"thumbnailUrl"`
	Status             string    `json:"status"`
	TargetCompleteDate time.Time `json:"targetCompleteDate"`
}

func (f BookCreateForm) ValidateBookCreateForm() error {
	switch {
	case f.Title == "":
		return errors.New("title is required")
	case f.Author == "":
		return errors.New("author is required")
	case f.TotalPages == 0:
		return errors.New("totalPages is required")
	case f.Publisher == "":
		return errors.New("publisher is required")
	case f.ThumbnailUrl == "":
		return errors.New("thumbnailUrl is required")
	case f.Status == "":
		return errors.New("status is required")
	case f.Status != "unread" && f.Status != "reading" && f.Status != "completed":
		return errors.New("status must be unread, reading, or completed")
	// targetCompleteDate は NewBookCreate でパース済み（ゼロ値なら不正）
	if f.TargetCompleteDate.IsZero() {
		return errors.New("targetCompleteDate is required or invalid format")
	}
	}
	return nil
}

func parseTargetCompleteDate(value string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid targetCompleteDate format: %s", value)
}
