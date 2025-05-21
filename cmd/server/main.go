package main

import (
	"log/slog"
	"net/http"

	"github.com/failuretoload/catdata/routes/index"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	index.RegisterRoutes(r)

	slog.Info("server started on :8080")
	http.ListenAndServe(":8080", r)
}
