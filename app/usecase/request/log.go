package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/service/validation"
)

type LogGet struct{}

func NewLogGet(_ *http.Request) (*LogGet, error) {
	return &LogGet{}, nil
}

type LogGetByID struct {
	LogID int
}

func NewLogGetByID(req *http.Request) (*LogGetByID, error) {
	id, err := parseIDParam(req, "id", "log id")
	if err != nil {
		return nil, err
	}
	return &LogGetByID{LogID: id}, nil
}

type LogGetByBookID struct {
	BookID int
}

func NewLogGetByBookID(req *http.Request) (*LogGetByBookID, error) {
	id, err := parseIDParam(req, "bookId", "book id")
	if err != nil {
		return nil, err
	}
	return &LogGetByBookID{BookID: id}, nil
}

type LogCreate struct {
	LogCreateForm
}

type LogCreateForm struct {
	BookID    int    `json:"bookId"`
	ReadDate  string `json:"readDate"`
	StartPage int    `json:"startPage"`
	EndPage   int    `json:"endPage"`
	Memo      string `json:"memo"`
}

func NewLogCreate(req *http.Request) (*LogCreate, error) {
	r := &LogCreate{}
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return nil, err
	}
	if err := r.ValidateLogCreateForm(); err != nil {
		return nil, err
	}
	return r, nil
}

func (f LogCreateForm) ValidateLogCreateForm() error {
	if f.BookID <= 0 {
		return errors.New("bookId is required")
	}
	if f.ReadDate == "" {
		return errors.New("readDate is required")
	}
	if f.StartPage <= 0 {
		return errors.New("startPage must be greater than 0")
	}
	if f.EndPage <= 1 {
		return errors.New("endPage must be greater than 1")
	}
	if f.StartPage >= f.EndPage {
		return errors.New("startPage must be less than endPage")
	}
	if f.StartPage > validation.MaxPagesDigits || f.EndPage > validation.MaxPagesDigits {
		return errors.New("startPage and endPage must be 4 digits or less")
	}
	if len(f.Memo) > validation.MaxLenLogNote {
		return errors.New("memo must be 800 characters or less")
	}
	return nil
}

type LogUpdate struct {
	LogID int
	LogUpdateForm
}

type LogUpdateForm struct {
	ReadDate  *string `json:"readDate"`
	StartPage *int    `json:"startPage"`
	EndPage   *int    `json:"endPage"`
	Memo      *string `json:"memo"`
}

func NewLogUpdate(req *http.Request) (*LogUpdate, error) {
	id, err := parseIDParam(req, "id", "log id")
	if err != nil {
		return nil, err
	}
	r := &LogUpdate{LogID: id}
	if err := json.NewDecoder(req.Body).Decode(&r.LogUpdateForm); err != nil {
		return nil, err
	}
	if err := r.ValidateLogUpdateForm(); err != nil {
		return nil, err
	}
	return r, nil
}

func (f LogUpdateForm) ValidateLogUpdateForm() error {
	if f.StartPage != nil && *f.StartPage <= 0 {
		return errors.New("startPage must be greater than 0")
	}
	if f.EndPage != nil && *f.EndPage <= 1 {
		return errors.New("endPage must be greater than 1")
	}
	if f.StartPage != nil && f.EndPage != nil && *f.StartPage >= *f.EndPage {
		return errors.New("startPage must be less than endPage")
	}
	if f.Memo != nil && len(*f.Memo) > validation.MaxLenLogNote {
		return errors.New("memo must be 800 characters or less")
	}
	return nil
}

type LogDelete struct {
	LogID int
}

func NewLogDelete(req *http.Request) (*LogDelete, error) {
	id, err := parseIDParam(req, "id", "log id")
	if err != nil {
		return nil, err
	}
	return &LogDelete{LogID: id}, nil
}

type LogDeleteByBookID struct {
	BookID int
}

func NewLogDeleteByBookID(req *http.Request) (*LogDeleteByBookID, error) {
	id, err := parseIDParam(req, "bookId", "book id")
	if err != nil {
		return nil, err
	}
	return &LogDeleteByBookID{BookID: id}, nil
}

// ParseReadDate は readDate 文字列を time.Time にパースする。RFC3339 または YYYY-MM-DD。
func ParseReadDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("readDate is required")
	}
	if strings.Contains(s, "T") || len(s) > 10 {
		return time.Parse(time.RFC3339, s)
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
}
