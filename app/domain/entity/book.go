package entity

import "time"

type Status string

const (
	StatusUnread Status = "unread"
	StatusReading Status = "reading"
	StatusCompleted Status = "completed"
)
// Book は本のドメインエンティティ。JSON と Datastore の両方で使う。
// ID は Datastore の Key で持つため datastore:"-" で保存しない。
type Book struct {
	ID                  int       `json:"id"                 datastore:"-"`
	Title               string    `json:"title"              datastore:"title"`
	Author              string    `json:"author"             datastore:"author"`
	TotalPages          int       `json:"totalPages"         datastore:"totalPages"`
	Publisher           string    `json:"publisher"          datastore:"publisher"`
	ThumbnailUrl        string    `json:"thumbnailUrl"       datastore:"thumbnailUrl"`
	Status              Status    `json:"status"             datastore:"status"`
	TargetCompleteDate  time.Time `json:"targetCompleteDate" datastore:"targetCompleteDate"`
	EncounterNote       string    `json:"encounterNote"      datastore:"encounterNote"`
	ReadPages           int       `json:"readPages"          datastore:"readPages"`
	TargetPagesPerDay   int       `json:"targetPagesPerDay"  datastore:"targetPagesPerDay"`
	CreatedAt           time.Time `json:"createdAt"          datastore:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"          datastore:"updatedAt"`
}
