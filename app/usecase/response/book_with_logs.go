package response

import "github.com/sora-00/booktracker-api/app/domain/entity"

type BookWithLogs struct {
	Book *entity.Book  `json:"book"`
	Logs []*entity.Log `json:"logs"`
}

type BookLogsByStatus struct {
	Items []*BookWithLogs `json:"items"`
}

