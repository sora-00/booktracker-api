package response

import (
	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type LogGet struct {
	Logs []*entity.Log `json:"logs"`
}

func NewLogGet(logs []entity.Log) *LogGet {
	items := make([]*entity.Log, 0, len(logs))
	for i := range logs {
		items = append(items, &logs[i])
	}
	return &LogGet{Logs: items}
}

type LogGetByID struct {
	*entity.Log
}

func NewLogGetByID(log *entity.Log) *LogGetByID {
	return &LogGetByID{Log: log}
}

type LogGetByBookID struct {
	Logs []*entity.Log `json:"logs"`
}

func NewLogGetByBookID(logs []entity.Log) *LogGetByBookID {
	items := make([]*entity.Log, 0, len(logs))
	for i := range logs {
		items = append(items, &logs[i])
	}
	return &LogGetByBookID{Logs: items}
}

type LogCreate struct {
	*entity.Log
}

func NewLogCreate(log *entity.Log) *LogCreate {
	return &LogCreate{Log: log}
}

type LogUpdate struct {
	*entity.Log
}

func NewLogUpdate(log *entity.Log) *LogUpdate {
	return &LogUpdate{Log: log}
}
