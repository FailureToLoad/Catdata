package weighttable

import (
	"net/http"

	"github.com/failuretoload/catdata/cat"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	service cat.Service
}

func NewController(s cat.Service) Controller {
	return Controller{
		service: s,
	}
}

func (c Controller) RegisterRoutes(r *chi.Mux) {
	r.Get("/query", c.fetchRows)
}

func (c Controller) fetchRows(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
