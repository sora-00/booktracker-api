package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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
	r := &BookCreate{}
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return nil, err
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

type BookUpdate struct {
	BookID int
	BookUpdateForm
}

func NewBookUpdate(req *http.Request) (*BookUpdate, error) {
	idStr := chi.URLParam(req, "id")
	if idStr == "" {
		return nil, errors.New("book id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid book id")
	}
	r := &BookUpdate{BookID: id}
	if err := json.NewDecoder(req.Body).Decode(&r.BookUpdateForm); err != nil {
		return nil, err
	}
	if err := r.ValidateBookUpdateForm(); err != nil {
		return nil, err
	}
	return r, nil
}

// ---

// NormalizedDate は targetCompleteDate 用。YYYY-MM-DD のみ受け付け、その日の 00:00:00Z に正規化してから DB に保存する。
type NormalizedDate time.Time

func (t *NormalizedDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return errors.New("targetCompleteDate is required or invalid format")
	}
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return errors.New("targetCompleteDate must be YYYY-MM-DD")
	}
	*t = NormalizedDate(time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC))
	return nil
}

func (t NormalizedDate) Time() time.Time { return time.Time(t) }

// BookCreateForm の targetCompleteDate は YYYY-MM-DD。正規化後 00:00:00Z で保存する。
type BookCreateForm struct {
	Title              string         `json:"title"`
	Author             string         `json:"author"`
	TotalPages         int            `json:"totalPages"`
	Publisher          string         `json:"publisher"`
	ThumbnailUrl       string         `json:"thumbnailUrl"`
	Status             string         `json:"status"`
	TargetCompleteDate NormalizedDate `json:"targetCompleteDate"`
	EncounterNote      string         `json:"encounterNote"`      // この本に出会った経緯
	ReadPages          int            `json:"readPages"`          // 読み終わったページ数
	TargetPagesPerDay  int            `json:"targetPagesPerDay"`  // 目標ページ数/日
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
	case f.TargetCompleteDate.Time().IsZero():
		return errors.New("targetCompleteDate is required or invalid format (use YYYY-MM-DD)")
	case f.ReadPages < 0:
		return errors.New("readPages must be 0 or greater")
	case f.ReadPages > f.TotalPages:
		return errors.New("readPages must not exceed totalPages")
	case f.TargetPagesPerDay < 0:
		return errors.New("targetPagesPerDay must be 0 or greater")
	}
	return nil
}

// BookUpdateForm は更新可能な項目のみ。送った項目だけ更新する（nil の項目は既存のまま）。
// targetCompleteDate は YYYY-MM-DD。正規化後 00:00:00Z で保存する。
type BookUpdateForm struct {
	ThumbnailUrl       *string         `json:"thumbnailUrl"`
	TargetCompleteDate *NormalizedDate `json:"targetCompleteDate"`
	EncounterNote      *string         `json:"encounterNote"`
	TargetPagesPerDay  *int            `json:"targetPagesPerDay"`
}

func (f BookUpdateForm) ValidateBookUpdateForm() error {
	if f.TargetPagesPerDay != nil && *f.TargetPagesPerDay < 0 {
		return errors.New("targetPagesPerDay must be 0 or greater")
	}
	if f.TargetCompleteDate != nil && f.TargetCompleteDate.Time().IsZero() {
		return errors.New("targetCompleteDate invalid format (use YYYY-MM-DD)")
	}
	return nil
}
