package entity

import "time"

type Status string

const (
	StatusUnread Status = "unread"
	StatusReading Status = "reading"
	StatusCompleted Status = "completed"
)
type Book struct {
	ID                  int         `json:"id"`
	Title               string      `json:"title"`
	Author              string      `json:"author"`
	TotalPages          int         `json:"totalPages"`
	Publisher           string      `json:"publisher"`
	ThumbnailUrl        string      `json:"thumbnailUrl"`
	Status              Status      `json:"status"`
	TargetCompleteDate  time.Time   `json:"targetCompleteDate"`
	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           time.Time   `json:"updatedAt"`
}
