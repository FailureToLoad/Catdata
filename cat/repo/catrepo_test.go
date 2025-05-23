package repo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/failuretoload/catdata/cat/domain"
	"github.com/failuretoload/catdata/cat/repo"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo_Query(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	limit := 3
	offset := 0
	query := fmt.Sprintf("SELECT .* FROM catstats LIMIT %d OFFSET %d", limit, offset)

	catRepo := repo.NewCatRepo(mock)
	timestamp := time.Now().Local()
	want := expected(timestamp)

	mock.ExpectQuery(query).
		WillReturnRows(pgxmock.NewRows([]string{
			"id",
			"name",
			"date",
			"weight",
			"notes",
		}).AddRows(
			[]any{
				want[0].ID,
				timestamp,
				want[0].Cat,
				want[0].Weight,
				want[0].Notes,
			},
			[]any{
				want[1].ID,
				timestamp,
				want[1].Cat,
				want[1].Weight,
				want[1].Notes,
			},
			[]any{
				want[2].ID,
				timestamp,
				want[2].Cat,
				want[2].Weight,
				want[2].Notes,
			},
		))

	got, err := catRepo.Query(context.Background(), domain.QueryInput{Limit: 3, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, got)
	assert.ElementsMatch(t, got, want)

}

func expected(timestamp time.Time) []domain.CatRecord {
	nimbusNote := "acquiring mass"
	yetiNote := "wants an egg"
	romNote := "stop eating carpet"
	timestring := timestamp.Format("01/02/2006 3:04 PM")
	return []domain.CatRecord{
		{
			ID:        uuid.New(),
			Cat:       "nimbus",
			Timestamp: timestring,
			Weight:    12,
			Notes:     &nimbusNote,
		},
		{
			ID:        uuid.New(),
			Cat:       "yeti",
			Timestamp: timestring,
			Weight:    8,
			Notes:     &yetiNote,
		},
		{
			ID:        uuid.New(),
			Cat:       "rom",
			Timestamp: timestring,
			Weight:    11,
			Notes:     &romNote,
		},
	}
}
