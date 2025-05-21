package cat

import (
	"context"

	"github.com/failuretoload/catdata/cat/domain"
)

type (
	Repo interface {
		Query(ctx context.Context, input domain.QueryInput) ([]domain.CatRecord, error)
	}
	Service struct {
		repo Repo
	}
)

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}
