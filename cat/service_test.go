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
func TestService_AddRecord(t *testing.T) {
	type args struct {
		cat    string
		weight float32
		notes  *string
	}
	tests := []struct {
		name       string
		args       args
		setup      func(repo *MockRepo)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success with notes",
			args: args{
				cat:    "nimbus",
				weight: 12.5,
				notes:  ptr("hungry"),
			},
			setup: func(repo *MockRepo) {
				repo.EXPECT().
					Insert(gomock.Any(), "nimbus", float32(12.5), ptr("hungry")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with nil notes",
			args: args{
				cat:    "yeti",
				weight: 10,
				notes:  nil,
			},
			setup: func(repo *MockRepo) {
				repo.EXPECT().
					Insert(gomock.Any(), "yeti", float32(10), nil).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty cat name",
			args: args{
				cat:    "",
				weight: 10,
				notes:  nil,
			},
			setup:      func(repo *MockRepo) {},
			wantErr:    true,
			wantErrMsg: "cat name cannot be empty",
		},
		{
			name: "zero weight",
			args: args{
				cat:    "rom",
				weight: 0,
				notes:  nil,
			},
			setup:      func(repo *MockRepo) {},
			wantErr:    true,
			wantErrMsg: "cat weight must be positive",
		},
		{
			name: "negative weight",
			args: args{
				cat:    "rom",
				weight: -5,
				notes:  nil,
			},
			setup:      func(repo *MockRepo) {},
			wantErr:    true,
			wantErrMsg: "cat weight must be positive",
		},
		{
			name: "empty notes string",
			args: args{
				cat:    "nimbus",
				weight: 12,
				notes:  ptr(""),
			},
			setup:      func(repo *MockRepo) {},
			wantErr:    true,
			wantErrMsg: "cannot submit empty notes",
		},
		{
			name: "repo returns error",
			args: args{
				cat:    "nimbus",
				weight: 12,
				notes:  ptr("hungry"),
			},
			setup: func(repo *MockRepo) {
				repo.EXPECT().
					Insert(gomock.Any(), "nimbus", float32(12), ptr("hungry")).
					Return(errors.New("repo error"))
			},
			wantErr:    true,
			wantErrMsg: "repo error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoFake := NewMockRepo(ctrl)
			tt.setup(repoFake)
			target := cat.NewService(repoFake)
			err := target.AddRecord(context.Background(), tt.args.cat, tt.args.weight, tt.args.notes)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.ErrorContains(t, err, tt.wantErrMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
