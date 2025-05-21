package weighttable

import (
	"net/http"

	"github.com/failuretoload/catdata/cat"
	"github.com/failuretoload/catdata/response"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"
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
	ctx := r.Context()
	limit := 50
	offset := 0

	rows, err := c.service.Cats(ctx, offset, limit)
	if err != nil {
		response.InternalServerError(ctx, w, "unable to request records", err)
	}

	sse := datastar.NewSSE(w, r)
	sse.MergeFragmentTempl(Rows(rows))
}
