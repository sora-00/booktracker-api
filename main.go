package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sora-00/booktracker-api/app/controller"
	"github.com/sora-00/booktracker-api/app/domain/service"
	dsrepo "github.com/sora-00/booktracker-api/app/infra/repository/datastore"
	"github.com/sora-00/booktracker-api/app/usecase"
	dsclient "github.com/sora-00/booktracker-api/pkg/datastore"
)

func main() {
	ctx := context.Background()
	// Cloud Datastore 接続
	ds, err := dsclient.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to connect datastore: %v", err)
	}
	defer ds.Close()

	// 依存関係の注入（infra層: Datastore 実装。ds は middleware で context に載せる）
	bookRepo := dsrepo.NewBookRepo()

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
	// 各リクエストの context に Datastore クライアントを入れる（repository で FromContext する前提）
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := dsclient.WithContext(r.Context(), ds)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

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
	// ブラウザが自動で叩く favicon は 204 で返して 404 ログを出さない
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
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
