package product_records

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrDatabase = errors.New("internal server error")
)

type Service interface {
	ValidateProductID(ctx context.Context, pid int) bool
	Create(ctx context.Context, pr domain.ProductRecord) (domain.ProductRecord, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// returns true if a product with id equal to pid exists in the database, false otherwise
func (s *service) ValidateProductID(ctx context.Context, pid int) bool {
	return s.r.ValidateProductID(ctx, pid)
}

// creates a product record and returns it
func (s *service) Create(ctx context.Context, pr domain.ProductRecord) (domain.ProductRecord, error) {
	id, err := s.r.Store(ctx, pr)
	if err != nil {
		return domain.ProductRecord{}, ErrDatabase
	}

	pr.ID = id

	return pr, nil
}
