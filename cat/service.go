package cat

import (
	"context"
	"errors"

	"github.com/failuretoload/catdata/cat/domain"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=cat_test
type (
	Repo interface {
		Query(ctx context.Context, input domain.QueryInput) ([]domain.CatRecord, error)
	}
	Service struct {
		repo Repo
	}
)

var (
	ErrOffset = errors.New("offset cannot be negative")
	ErrLimit  = errors.New("limit cannot be negative")
)

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) Cats(ctx context.Context, offset, limit int) ([]domain.CatRecord, error) {
	if offset < 0 {
		return nil, ErrOffset
	}
	if limit < 0 {
		return nil, ErrLimit
	}

	input := domain.QueryInput{
		Offset: offset,
		Limit:  limit,
	}

	return s.repo.Query(ctx, input)
}
