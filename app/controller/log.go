package controller

import (
	"net/http"

	"github.com/sora-00/booktracker-api/app/controller/respond"
	httpx "github.com/sora-00/booktracker-api/app/infra/http"
	"github.com/sora-00/booktracker-api/app/usecase"
	"github.com/sora-00/booktracker-api/app/usecase/request"
)

type LogController struct {
	Log *usecase.Log
}

func NewLogController(l *usecase.Log) *LogController {
	return &LogController{Log: l}
}

func (c *LogController) GetLogs(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogGet(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Log.Get(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *LogController) GetLogByID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogGetByID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Log.GetByID(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "log not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *LogController) GetLogsByBookID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogGetByBookID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Log.GetByBookID(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *LogController) CreateLog(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogCreate(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Log.Create(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "book not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusCreated, res)
}

func (c *LogController) UpdateLog(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogUpdate(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Log.Update(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "log not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *LogController) DeleteLog(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogDelete(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	err = c.Log.Delete(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "log not found")
		respond.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *LogController) DeleteLogsByBookID(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewLogDeleteByBookID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.Log.DeleteByBookID(r.Context(), req); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
