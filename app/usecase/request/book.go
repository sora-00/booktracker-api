package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/service/validation"
)

type BookGet struct{}

func NewBookGet(_ *http.Request) (*BookGet, error) {
	return &BookGet{}, nil
}

type BookGetByStatus struct {
	Status string
}

func NewBookGetByStatus(req *http.Request) (*BookGetByStatus, error) {
	status := strings.TrimSpace(chi.URLParam(req, "status"))
	if status == "" {
		return nil, errors.New("status is required")
	}
	if err := validation.ValidateBookStatus(status); err != nil {
		return nil, err
	}
	return &BookGetByStatus{Status: status}, nil
}

type BookGetByID struct {
	BookID int
}

func NewBookGetByID(req *http.Request) (*BookGetByID, error) {
	id, err := parseIDParam(req, "id", "book id")
	if err != nil {
		return nil, err
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
	BookID int
}

func NewBookDelete(req *http.Request) (*BookDelete, error) {
	id, err := parseIDParam(req, "id", "book id")
	if err != nil {
		return nil, err
	}
	return &BookDelete{BookID: id}, nil
}

type BookUpdate struct {
	BookID int
	BookUpdateForm
}

func NewBookUpdate(req *http.Request) (*BookUpdate, error) {
	id, err := parseIDParam(req, "id", "book id")
	if err != nil {
		return nil, err
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

// NormalizedDate は targetCompleteDate 用。YYYY-MM-DD または RFC3339 を受け付け、00:00:00 UTC に正規化する。
type NormalizedDate time.Time

func (t *NormalizedDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return errors.New("targetCompleteDate is required or invalid format")
	}
	var parsed time.Time
	var err error
	if strings.Contains(s, "T") || len(s) > 10 {
		parsed, err = time.Parse(time.RFC3339, s)
	} else {
		parsed, err = time.Parse("2006-01-02", s)
	}
	if err != nil {
		return errors.New("targetCompleteDate must be YYYY-MM-DD or RFC3339")
	}
	*t = NormalizedDate(time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC))
	return nil
}

func (t NormalizedDate) Time() time.Time { return time.Time(t) }

// BookCreateForm は POST /books のリクエストボディ。
type BookCreateForm struct {
	Title              string         `json:"title"`
	Author             string         `json:"author"`
	Publisher          string         `json:"publisher"`
	ThumbnailUrl       string         `json:"thumbnailUrl"`
	TotalPages         int            `json:"totalPages"`
	TargetCompleteDate NormalizedDate `json:"targetCompleteDate"`
	TargetPagesPerDay  int            `json:"targetPagesPerDay"`
	Status             entity.Status  `json:"status"`
	ReadPages          int            `json:"readPages"`
	EncounterNote      string         `json:"encounterNote"`
}

func (f BookCreateForm) ValidateBookCreateForm() error {
	if f.Title == "" {
		return errors.New("title is required")
	}
	if len(f.Title) > validation.MaxLenTitle {
		return errors.New("title must be 30 characters or less")
	}
	if f.Author == "" {
		return errors.New("author is required")
	}
	if len(f.Author) > validation.MaxLenAuthor {
		return errors.New("author must be 30 characters or less")
	}
	if f.Publisher == "" {
		return errors.New("publisher is required")
	}
	if len(f.Publisher) > validation.MaxLenPublish {
		return errors.New("publisher must be 20 characters or less")
	}
	if f.ThumbnailUrl == "" {
		return errors.New("thumbnailUrl is required")
	}
	if err := validation.ValidatePagesAll(f.TotalPages); err != nil {
		return err
	}
	if err := validation.ValidateBookStatus(string(f.Status)); err != nil {
		return err
	}
	if f.TargetCompleteDate.Time().Before(time.Now().UTC().Truncate(24 * time.Hour)) {
		return errors.New("targetCompleteDate must be after now")
	}
	if err := validation.ValidateTargetReadPagesPerDay(f.TargetPagesPerDay, f.TotalPages); err != nil {
		return err
	}
	if f.Status != entity.StatusUnread {
		if f.ReadPages < 0 {
			return errors.New("readPages is required when status is not unread")
		}
		if f.ReadPages > f.TotalPages {
			return errors.New("readPages must not exceed totalPages")
		}
		if f.ReadPages > validation.MaxPagesDigits {
			return errors.New("readPages must be 4 digits or less")
		}
	}
	if len(f.EncounterNote) > validation.MaxLenBackground {
		return errors.New("encounterNote must be 200 characters or less")
	}
	return nil
}

// BookUpdateForm は PUT /books/:id のリクエストボディ（全フィールドがオプション）。
type BookUpdateForm struct {
	Title              *string         `json:"title"`
	Author             *string         `json:"author"`
	Publisher          *string         `json:"publisher"`
	ThumbnailUrl       *string         `json:"thumbnailUrl"`
	TotalPages         *int            `json:"totalPages"`
	TargetCompleteDate *NormalizedDate `json:"targetCompleteDate"`
	TargetPagesPerDay  *int            `json:"targetPagesPerDay"`
	Status             *string         `json:"status"`
	EncounterNote      *string         `json:"encounterNote"`
}

func (f BookUpdateForm) ValidateBookUpdateForm() error {
	if f.Title != nil {
		if *f.Title == "" {
			return errors.New("title cannot be empty")
		}
		if len(*f.Title) > validation.MaxLenTitle {
			return errors.New("title must be 30 characters or less")
		}
	}
	if f.Author != nil && len(*f.Author) > validation.MaxLenAuthor {
		return errors.New("author must be 30 characters or less")
	}
	if f.Publisher != nil && len(*f.Publisher) > validation.MaxLenPublish {
		return errors.New("publisher must be 20 characters or less")
	}
	if f.ThumbnailUrl != nil && *f.ThumbnailUrl == "" {
		return errors.New("thumbnailUrl cannot be empty")
	}
	if f.TotalPages != nil {
		if err := validation.ValidatePagesAll(*f.TotalPages); err != nil {
			return err
		}
	}
	if f.TargetPagesPerDay != nil {
		pagesAll := 1
		if f.TotalPages != nil {
			pagesAll = *f.TotalPages
		}
		if err := validation.ValidateTargetReadPagesPerDay(*f.TargetPagesPerDay, pagesAll); err != nil {
			return err
		}
	}
	if f.Status != nil {
		if err := validation.ValidateBookStatus(*f.Status); err != nil {
			return err
		}
	}
	if f.EncounterNote != nil && len(*f.EncounterNote) > validation.MaxLenBackground {
		return errors.New("encounterNote must be 200 characters or less")
	}
	return nil
}

func parseIDParam(req *http.Request, param, displayName string) (int, error) {
	s := chi.URLParam(req, param)
	if s == "" {
		return 0, errors.New(displayName + " is required")
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid " + displayName)
	}
	return id, nil
}
