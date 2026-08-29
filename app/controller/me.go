package controller

import (
	"net/http"

	"github.com/sora-00/booktracker-api/app/controller/respond"
	httpx "github.com/sora-00/booktracker-api/app/infra/http"
	"github.com/sora-00/booktracker-api/app/usecase"
)

type MeController struct {
	Me *usecase.Me
}

func NewMeController(m *usecase.Me) *MeController {
	return &MeController{Me: m}
}

func (c *MeController) GetMe(w http.ResponseWriter, r *http.Request) {
	res, err := c.Me.Get(r.Context())
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "me not found")
		respond.Error(w, status, msg)
		return
	}
	respond.Data(w, http.StatusOK, res)
}

func (c *MeController) DeleteMe(w http.ResponseWriter, r *http.Request) {
	err := c.Me.Delete(r.Context())
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "me not found")
		respond.Error(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
