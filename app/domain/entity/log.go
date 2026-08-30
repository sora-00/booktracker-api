package entity

import "time"

type Log struct {
	ID        int       `json:"id"         datastore:"-"`
	UserID    string    `json:"userId"     datastore:"userId"`
	BookID    int       `json:"bookId"     datastore:"bookId"`
	ReadDate  time.Time `json:"readDate"   datastore:"readDate"`
	StartPage int       `json:"startPage"  datastore:"startPage"`
	EndPage   int       `json:"endPage"    datastore:"endPage"`
	PagesRead int       `json:"pagesRead"  datastore:"pagesRead"`
	Memo      string    `json:"memo"       datastore:"note"`
	CreatedAt time.Time `json:"createdAt"  datastore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"  datastore:"updatedAt"`
}

func NewLog(userID string, bookID int, readDate time.Time, startPage, endPage int, memo string) *Log {
	l := &Log{
		UserID:    userID,
		BookID:    bookID,
		ReadDate:  readDate,
		StartPage: startPage,
		EndPage:   endPage,
		Memo:      memo,
	}
	l.ComputePagesRead()
	l.Created()
	return l
}

func (l *Log) Created() {
	now := time.Now().UTC()
	l.CreatedAt = now
	l.UpdatedAt = now
}

func (l *Log) Updated() {
	l.UpdatedAt = time.Now().UTC()
}

// ComputePagesRead は StartPage・EndPage から PagesRead を再計算する。
func (l *Log) ComputePagesRead() {
	l.PagesRead = l.EndPage - l.StartPage + 1
}

type LogPatch struct {
	ReadDate  *time.Time
	StartPage *int
	EndPage   *int
	Memo      *string
}

func (l *Log) ApplyPatch(p LogPatch) {
	if p.ReadDate != nil {
		l.ReadDate = *p.ReadDate
	}
	if p.StartPage != nil {
		l.StartPage = *p.StartPage
	}
	if p.EndPage != nil {
		l.EndPage = *p.EndPage
	}
	if p.Memo != nil {
		l.Memo = *p.Memo
	}
}
