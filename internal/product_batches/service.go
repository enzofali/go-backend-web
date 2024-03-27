package product_batches

import (
	"context"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

type Service interface {
	Create(ctx context.Context, productBatches domain.ProductBatches) (domain.ProductBatches, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) Create(ctx context.Context, productBatches domain.ProductBatches) (domain.ProductBatches, error) {
	id, err := s.r.Create(ctx, productBatches)
	if err != nil {
		return domain.ProductBatches{}, err
	}

	productBatches.ID = id

	return productBatches, nil
}
