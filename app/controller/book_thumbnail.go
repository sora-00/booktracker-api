package controller

import (
	"net/http"

	"github.com/sora-00/booktracker-api/app/controller/respond"
	httpx "github.com/sora-00/booktracker-api/app/infra/http"
	"github.com/sora-00/booktracker-api/app/usecase"
	"github.com/sora-00/booktracker-api/app/usecase/request"
)

// BookThumbnailController は本の表紙画像アップロード用のHTTPハンドラです。
type BookThumbnailController struct {
	BookThumbnail *usecase.BookThumbnail
}

func NewBookThumbnailController(bookThumbnail *usecase.BookThumbnail) *BookThumbnailController {
	return &BookThumbnailController{BookThumbnail: bookThumbnail}
}

// PostThumbnail は本の表紙画像を multipart/form-data で受け取り保存し、{ id, url } を返す。
func (c *BookThumbnailController) PostThumbnail(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookThumbnailPost(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.BookThumbnail.Create(r.Context(), req, r.Host, r.TLS != nil)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.Data(w, http.StatusOK, res)
}

// GetThumbnail は保存した本の表紙画像を返す。
func (c *BookThumbnailController) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewBookThumbnailGet(r)
	if err != nil {
		respond.Error(w, http.StatusNotFound, "not found")
		return
	}
	res, err := c.BookThumbnail.Get(r.Context(), req)
	if err != nil {
		status, msg := httpx.StatusAndMessage(err, "not found")
		respond.Error(w, status, msg)
		return
	}

	w.Header().Set("Content-Type", res.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(res.Body)
}
