package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sora-00/booktracker-api/app/controller"
	"github.com/sora-00/booktracker-api/app/domain/service"
	"github.com/sora-00/booktracker-api/app/infra/repository/postgres"
	"github.com/sora-00/booktracker-api/app/usecase"
	"github.com/sora-00/booktracker-api/pkg/db"
)

func main() {
	// DB接続
	conn, err := db.NewPostgres()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer conn.Close()

	// 依存関係の注入
	// infra層（PostgreSQL 実装）→ domain の repository インターフェースを満たす
	bookRepo := postgres.NewBookRepo(conn)

	// domain層（ビジネスロジック）
	bookService := service.NewService(bookRepo)

	// usecase層（アプリケーションロジック）
	bookUsecase := usecase.NewUsecase(bookRepo, bookService)

	// controller層（HTTPハンドラ）
	bookController := controller.NewController(bookUsecase)

	// ルーティング設定
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 404/405を可視化
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s %s", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("405 Method Not Allowed: %s %s", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// /api/books（末尾なし）も直に受ける
	r.Route("/api", func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", bookController.GetBooks)
			r.Get("/{id}", bookController.GetBookByID)
			r.Post("/", bookController.CreateBook)
			r.Delete("/{id}", bookController.DeleteBook)
		})
	})

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	addr := ":" + port
	log.Printf("Listening on %s 🚀\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
