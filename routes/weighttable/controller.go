package weighttable

import (
	"math"
	"net/http"
	"strconv"

	"github.com/failuretoload/catdata/cat"
	"github.com/failuretoload/catdata/cat/domain"
	"github.com/failuretoload/catdata/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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
	dtoSlice := toDTOSlice(rows)
	sse := datastar.NewSSE(w, r)
	sse.MergeFragmentTempl(tableRows(dtoSlice))
}

func toDTOSlice(rows []domain.CatRecord) []CatRowDTO {
	result := []CatRowDTO{}
	for _, row := range rows {
		pounds, ounces := kgToPoundsOunces(row.Weight)

		dto := CatRowDTO{
			ID:        row.ID,
			Cat:       row.Cat,
			Timestamp: row.Timestamp,
			WeightKG:  row.Weight,
			WeightLB:  pounds,
			WeightOZ:  ounces,
			Notes:     row.Notes,
		}
		result = append(result, dto)
	}
	return result
}

func kgToPoundsOunces(kg float32) (int, float32) {
	totalPounds := kg * 2.20462

	pounds := int(totalPounds)

	remainingPounds := totalPounds - float32(pounds)
	ounces := remainingPounds * 16

	ounces = float32(math.Round(float64(ounces)*10) / 10)

	if ounces >= 16.0 {
		pounds++
		ounces = 0.0
	}

	return pounds, ounces
}

type CatRowDTO struct {
	ID        uuid.UUID
	Cat       string
	Timestamp string
	WeightKG  float32
	WeightLB  int
	WeightOZ  float32
	Notes     *string
}
