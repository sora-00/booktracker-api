package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// 本の表紙画像の保存先（本の表紙専用であることが分かるように）
const bookThumbnailDir = "uploads/thumbnails"

// BookThumbnailController は本の表紙画像アップロード用のHTTPハンドラです。
// 当面は認証なし。のちにログイン必須に変更する。
type BookThumbnailController struct{}

func NewBookThumbnailController() *BookThumbnailController {
	return &BookThumbnailController{}
}

// PostThumbnail は本の表紙画像を multipart/form-data で受け取り保存し、{ id, url } を返す。
func (c *BookThumbnailController) PostThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 10MB 制限
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	ext = strings.ToLower(ext)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		ext = ".jpg"
	}

	id, err := randomID()
	if err != nil {
		http.Error(w, "failed to generate id", http.StatusInternalServerError)
		return
	}
	name := id + ext

	if err := os.MkdirAll(bookThumbnailDir, 0755); err != nil {
		log.Printf("book_thumbnail: mkdir %s: %v", bookThumbnailDir, err)
		http.Error(w, "failed to create upload dir", http.StatusInternalServerError)
		return
	}

	path := filepath.Join(bookThumbnailDir, name)
	dst, err := os.Create(path)
	if err != nil {
		log.Printf("book_thumbnail: create %s: %v", path, err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(path)
		log.Printf("book_thumbnail: copy: %v", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + r.Host + "/api/books/thumbnails/" + name

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id, "url": url})
}

// GetThumbnail は保存した本の表紙画像を返す。
func (c *BookThumbnailController) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" || strings.Contains(id, "/") || strings.Contains(id, "..") {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	path := filepath.Join(bookThumbnailDir, id)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil || info.IsDir() {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	ext := strings.ToLower(filepath.Ext(id))
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	default:
		w.Header().Set("Content-Type", "image/jpeg")
	}
	w.Header().Set("Cache-Control", "public, max-age=86400")
	io.Copy(w, f)
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
