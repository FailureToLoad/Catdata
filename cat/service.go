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
		Insert(ctx context.Context, cat string, weight float32, notes *string) error
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

func (s Service) AddRecord(ctx context.Context, cat string, weight float32, notes *string) error {
	if cat == "" {
		return errors.New("cat name cannot be empty")
	}
	if weight <= 0 {
		return errors.New("cat weight must be positive")
	}
	if notes != nil && *notes == "" {
		return errors.New("cannot submit empty notes")
	}

	return s.repo.Insert(ctx, cat, weight, notes)
}
