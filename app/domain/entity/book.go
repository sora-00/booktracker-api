package entity

import "time"

type Status string

const (
	StatusUnread  Status = "unread"
	StatusReading Status = "reading"
	StatusRead    Status = "read"
)

type Book struct {
	ID                 int       `json:"id"                  datastore:"-"`
	UserID             string    `json:"userId"              datastore:"userId"`
	Title              string    `json:"title"               datastore:"title"`
	Author             string    `json:"author"              datastore:"author"`
	Publisher          string    `json:"publisher"           datastore:"publish"`
	ThumbnailUrl       string    `json:"thumbnailUrl"        datastore:"coverUrl"`
	TotalPages         int       `json:"totalPages"          datastore:"pagesAll"`
	TargetCompleteDate time.Time `json:"targetCompleteDate"  datastore:"targetReadDate"`
	TargetPagesPerDay  int       `json:"targetPagesPerDay"   datastore:"targetReadPagesPerDay"`
	Status             Status    `json:"status"              datastore:"status"`
	ReadPages          int       `json:"readPages"           datastore:"readPages"`
	EncounterNote      string    `json:"encounterNote"       datastore:"background"`
	CreatedAt          time.Time `json:"createdAt"           datastore:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"           datastore:"updatedAt"`
}

func NewBook(userID, title, author, publisher, thumbnailUrl string, totalPages int, targetCompleteDate time.Time, targetPagesPerDay int, status Status, readPages int, encounterNote string) *Book {
	b := &Book{
		UserID:             userID,
		Title:              title,
		Author:             author,
		Publisher:          publisher,
		ThumbnailUrl:       thumbnailUrl,
		TotalPages:         totalPages,
		TargetCompleteDate: targetCompleteDate,
		TargetPagesPerDay:  targetPagesPerDay,
		Status:             status,
		ReadPages:          readPages,
		EncounterNote:      encounterNote,
	}
	b.Created()
	return b
}

func (b *Book) Created() {
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now
}

func (b *Book) Updated() {
	b.UpdatedAt = time.Now().UTC()
}

type BookPatch struct {
	Title              *string
	Author             *string
	Publisher          *string
	ThumbnailUrl       *string
	TotalPages         *int
	TargetCompleteDate *time.Time
	TargetPagesPerDay  *int
	Status             *Status
	EncounterNote      *string
}

func (b *Book) ApplyPatch(p BookPatch) {
	if p.Title != nil {
		b.Title = *p.Title
	}
	if p.Author != nil {
		b.Author = *p.Author
	}
	if p.Publisher != nil {
		b.Publisher = *p.Publisher
	}
	if p.ThumbnailUrl != nil {
		b.ThumbnailUrl = *p.ThumbnailUrl
	}
	if p.TotalPages != nil {
		b.TotalPages = *p.TotalPages
	}
	if p.TargetCompleteDate != nil {
		b.TargetCompleteDate = *p.TargetCompleteDate
	}
	if p.TargetPagesPerDay != nil {
		b.TargetPagesPerDay = *p.TargetPagesPerDay
	}
	if p.Status != nil {
		b.Status = *p.Status
	}
	if p.EncounterNote != nil {
		b.EncounterNote = *p.EncounterNote
	}
}
