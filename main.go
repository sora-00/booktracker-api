package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/sora-00/booktracker-api/app/controller"
	"github.com/sora-00/booktracker-api/app/controller/authn"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/domain/service"
	"github.com/sora-00/booktracker-api/app/infra/auth"
	dsclient "github.com/sora-00/booktracker-api/app/infra/datastore"
	infrarepo "github.com/sora-00/booktracker-api/app/infra/repository"
	"github.com/sora-00/booktracker-api/app/infra/storage"
	"github.com/sora-00/booktracker-api/app/usecase"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	ds, err := dsclient.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to connect datastore: %v", err)
	}
	defer ds.Close()

	// --- repos ---
	bookRepo := infrarepo.NewBookRepo()
	logRepo := infrarepo.NewLogRepo()
	meRepo := infrarepo.NewMeRepo()
	bookThumbnailRepo := initBookThumbnailRepo(ctx)

	// --- services ---
	logSvc := service.NewLogService(logRepo)

	// --- usecases ---
	bookUsecase := usecase.NewBook(bookRepo, logRepo)
	logUsecase := usecase.NewLog(logRepo, bookRepo, logSvc)
	meUsecase := usecase.NewMe(meRepo)
	thumbUsecase := usecase.NewBookThumbnail(bookThumbnailRepo)

	// --- controllers ---
	bookCtrl := controller.NewBookController(bookUsecase)
	logCtrl := controller.NewLogController(logUsecase)
	meCtrl := controller.NewMeController(meUsecase)
	thumbCtrl := controller.NewBookThumbnailController(thumbUsecase)

	// --- auth ---
	credPath := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	fbVerifier, err := auth.NewFirebaseVerifier(ctx, credPath)
	if err != nil {
		log.Printf("firebase auth disabled (init failed): %v", err)
		fbVerifier = nil
	}

	// --- router ---
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true,"message":"BookTracker API"}`))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(withDatastore(ds))
		r.Use(authn.RequireAuth(fbVerifier, meRepo))

		r.Route("/me", func(r chi.Router) {
			r.Get("/", meCtrl.GetMe)
			r.Delete("/", meCtrl.DeleteMe)
		})

		r.Route("/books", func(r chi.Router) {
			r.Get("/", bookCtrl.GetBooks)
			r.Post("/", bookCtrl.CreateBook)

			r.Route("/status/{status}", func(r chi.Router) {
				r.Get("/", bookCtrl.GetBooksByStatus)
				r.Get("/logs", bookCtrl.GetLogsByBookStatus)
			})

			r.Route("/thumbnails", func(r chi.Router) {
				r.Post("/", thumbCtrl.PostThumbnail)
				r.Get("/{id}", thumbCtrl.GetThumbnail)
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", bookCtrl.GetBookByID)
				r.Put("/", bookCtrl.UpdateBook)
				r.Delete("/", bookCtrl.DeleteBook)
			})

			r.Route("/{bookId}/logs", func(r chi.Router) {
				r.Get("/", logCtrl.GetLogsByBookID)
				r.Delete("/", logCtrl.DeleteLogsByBookID)
			})
		})

		r.Route("/logs", func(r chi.Router) {
			r.Get("/", logCtrl.GetLogs)
			r.Post("/", logCtrl.CreateLog)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", logCtrl.GetLogByID)
				r.Put("/", logCtrl.UpdateLog)
				r.Delete("/", logCtrl.DeleteLog)
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	log.Printf("Listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// withDatastore は Datastore クライアントをリクエストのコンテキストに載せるミドルウェア。
func withDatastore(ds *datastore.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := dsclient.WithContext(r.Context(), ds)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func initBookThumbnailRepo(ctx context.Context) repository.BookThumbnailRepo {
	repo, err := storage.NewBookThumbnailRepo(
		ctx,
		os.Getenv("AWS_S3_BUCKET"),
		os.Getenv("AWS_S3_PREFIX"),
		os.Getenv("AWS_REGION"),
	)
	if err != nil {
		log.Fatalf("failed to initialize S3 thumbnail storage: %v", err)
	}
	return repo
}
