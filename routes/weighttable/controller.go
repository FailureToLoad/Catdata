package weighttable

import (
	"net/http"
	"strconv"

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
	r.Get("/insert", c.addRecord)
}

func (c Controller) fetchRows(w http.ResponseWriter, r *http.Request) {
	c.updateTable(w, r, 50, 0)
}

func (c Controller) addRecord(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	formData := r.Form
	weightStr := formData.Get("weight")
	f, err := strconv.ParseFloat(weightStr, 32)
	if err != nil {
		response.BadRequest(ctx, w, "unable to parse weight", err)
	}
	weight := float32(f)
	cat := formData.Get("cat")
	notesStr := formData.Get("notes")
	var notes *string
	if notesStr != "" {
		notes = &notesStr
	}

	err = c.service.AddRecord(ctx, cat, weight, notes)
	if err != nil {
		response.BadRequest(ctx, w, "unable to add record", err)
	}

	c.updateTable(w, r, 50, 0)
}

func (c Controller) updateTable(w http.ResponseWriter, r *http.Request, limit, offset int) {
	ctx := r.Context()

	rows, err := c.service.Cats(ctx, offset, limit)
	if err != nil {
		response.InternalServerError(ctx, w, "unable to request records", err)
	}

	sse := datastar.NewSSE(w, r)
	sse.MergeFragmentTempl(tableRows(rows))
}
