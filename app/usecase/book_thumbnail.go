package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"

	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/usecase/request"
	"github.com/sora-00/booktracker-api/app/usecase/response"
)

type BookThumbnail struct {
	repo repository.BookThumbnailRepo
}

func NewBookThumbnail(repo repository.BookThumbnailRepo) *BookThumbnail {
	return &BookThumbnail{repo: repo}
}

func (u BookThumbnail) Create(ctx context.Context, r *request.BookThumbnailPost, host string, isTLS bool) (*response.BookThumbnailPost, error) {
	ext := normalizedImageExt(r.Filename)
	id, err := randomID()
	if err != nil {
		return nil, errors.New("failed to generate id")
	}
	name := id + ext
	if err := u.repo.Save(ctx, name, r.Data); err != nil {
		return nil, errors.New("failed to save file")
	}

	scheme := "http"
	if isTLS {
		scheme = "https"
	}
	url := scheme + "://" + host + "/api/books/thumbnails/" + name
	return response.NewBookThumbnailPost(id, url), nil
}

func (u BookThumbnail) Get(ctx context.Context, r *request.BookThumbnailGet) (*response.BookThumbnailGet, error) {
	data, err := u.repo.Load(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	contentType := contentTypeByExt(r.ID)
	return response.NewBookThumbnailGet(contentType, data), nil
}

func normalizedImageExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return ext
	default:
		return ".jpg"
	}
}

func contentTypeByExt(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
