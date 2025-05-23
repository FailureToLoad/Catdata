package cat_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/failuretoload/catdata/cat"
	"github.com/failuretoload/catdata/cat/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Cats(t *testing.T) {
	t.Run("limit can't be negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoFake := NewMockRepo(ctrl)
		target := cat.NewService(repoFake)
		records, err := target.Cats(context.Background(), 0, -1)
		require.ErrorIs(t, err, cat.ErrLimit)
		assert.Empty(t, records)
	})

	t.Run("offset can't be negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoFake := NewMockRepo(ctrl)
		target := cat.NewService(repoFake)
		records, err := target.Cats(context.Background(), -1, 0)
		require.ErrorIs(t, err, cat.ErrOffset)
		assert.Empty(t, records)
	})

	t.Run("reports repo errors", func(t *testing.T) {
		limit := 1
		offset := 1
		ctrl := gomock.NewController(t)
		repoFake := NewMockRepo(ctrl)
		repoFake.EXPECT().
			Query(gomock.Any(), domain.QueryInput{
				Limit:  limit,
				Offset: offset,
			}).
			Return(
				nil,
				errors.New("repo error"),
			)

		target := cat.NewService(repoFake)
		records, err := target.Cats(context.Background(), offset, limit)
		require.Error(t, err)
		assert.ErrorContains(t, err, "repo error")
		assert.Empty(t, records)
	})

	t.Run("return records on success", func(t *testing.T) {
		limit := 3
		offset := 0
		ctrl := gomock.NewController(t)
		repoFake := NewMockRepo(ctrl)
		nimbusNote := "acquiring mass"
		yetiNote := "wants an egg"
		romNote := "stop eating carpet"
		timeformat := "01/02/2006 3:04 PM"
		expected := []domain.CatRecord{
			{
				Cat:       "nimbus",
				Weight:    13,
				Timestamp: time.Now().Local().Format(timeformat),
				Notes:     &nimbusNote,
			},
			{
				Cat:       "yeti",
				Weight:    10,
				Timestamp: time.Now().Local().Format(timeformat),
				Notes:     &yetiNote,
			},
			{
				Cat:       "rom",
				Weight:    11,
				Timestamp: time.Now().Local().Format(timeformat),
				Notes:     &romNote,
			},
		}
		repoFake.EXPECT().
			Query(gomock.Any(), domain.QueryInput{
				Limit:  limit,
				Offset: offset,
			}).
			Return(
				expected,
				nil,
			)

		target := cat.NewService(repoFake)
		records, err := target.Cats(context.Background(), offset, limit)
		require.NoError(t, err)
		assert.NotEmpty(t, records)
		assert.ElementsMatch(t, expected, records)
	})

}
