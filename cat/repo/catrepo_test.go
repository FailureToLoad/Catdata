package repo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/failuretoload/catdata/cat/domain"
	"github.com/failuretoload/catdata/cat/repo"
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
	want := expected()

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
				want[0].Timestamp,
				want[0].Cat,
				want[0].Weight,
				want[0].Notes,
			},
			[]any{
				want[1].ID,
				want[1].Timestamp,
				want[1].Cat,
				want[1].Weight,
				want[1].Notes,
			},
			[]any{
				want[2].ID,
				want[2].Timestamp,
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

func expected() []domain.CatRecord {
	nimbusNote := "acquiring mass"
	yetiNote := "wants an egg"
	romNote := "stop eating carpet"
	return []domain.CatRecord{
		{
			ID:        1,
			Cat:       "nimbus",
			Timestamp: time.Now().Unix(),
			Weight:    12,
			Notes:     &nimbusNote,
		},
		{
			ID:        2,
			Cat:       "yeti",
			Timestamp: time.Now().Unix(),
			Weight:    8,
			Notes:     &yetiNote,
		},
		{
			ID:        3,
			Cat:       "rom",
			Timestamp: time.Now().Unix(),
			Weight:    11,
			Notes:     &romNote,
		},
	}
}
