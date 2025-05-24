package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	index.RegisterRoutes(r)

	ctx := context.Background()
	pool := initConnPool(ctx)
	wtc := makeWeightTableController(pool)
	wtc.RegisterRoutes(r)

	serveStatic(r)

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

func serveStatic(r chi.Router) {
	path := "/static"
	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, "static"))
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
