package entity

import (
	"errors"
	"time"
)

type Me struct {
	ID        string    `json:"id"        datastore:"-"`
	CreatedAt time.Time `json:"createdAt" datastore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" datastore:"updatedAt"`
}

// NewMe は初回登録用。Created / Updated で時刻を入れる。
func NewMe(id string) (*Me, error) {
	if id == "" {
		return nil, errors.New("entity: me id is required")
	}
	m := &Me{ID: id}
	m.Created()
	return m, nil
}

func (m *Me) Created() {
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *Me) Updated() {
	m.UpdatedAt = time.Now().UTC()
}
