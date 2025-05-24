package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/failuretoload/catdata/cat"
	catrepo "github.com/failuretoload/catdata/cat/repo"
	"github.com/failuretoload/catdata/routes/index"
	"github.com/failuretoload/catdata/routes/weighttable"
	"github.com/go-chi/chi/v5"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	r := chi.NewRouter()
	serveCSS(r)
	serveFavicon(r)
	index.RegisterRoutes(r)

	ctx := context.Background()
	pool := initConnPool(ctx)
	wtc := makeWeightTableController(pool)
	wtc.RegisterRoutes(r)

	slog.Info("server started on :8080")
	http.ListenAndServe(":8080", r)
}

func makeWeightTableController(pool *pgxpool.Pool) weighttable.Controller {
	catRepo := catrepo.NewCatRepo(pool)
	catService := cat.NewService(catRepo)
	return weighttable.NewController(catService)
}

func initConnPool(ctx context.Context) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to parse db config: %v\n", err)
	}

	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return pool
}

func serveCSS(r chi.Router) {
	r.Get("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./static/styles.css")
	})
}

func serveFavicon(r chi.Router) {
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		http.ServeFile(w, r, "./static/favicon.ico")
	})
}
