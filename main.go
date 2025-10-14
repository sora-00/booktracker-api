package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS設定
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // 開発中は全部許可でOK
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	// ヘルスチェック
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// APIルート
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/books", func(w http.ResponseWriter, r *http.Request) {
			books := []Book{
				{ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"},
				{ID: 2, Title: "1984", Author: "George Orwell"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(books)
		})
	})

	// サーバ起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Println("Server starting on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
