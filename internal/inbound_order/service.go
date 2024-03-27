package inboundorder

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

var (
	ErrDatabase = errors.New("Database error.")
)

type Service interface {
	Create(ctx context.Context, i domain.InboundOrder) (domain.InboundOrder, error)
}

type service struct {
	repository Repository
}

func NewService(inboundOrder Repository) Service {
	return &service{
		repository: inboundOrder,
	}
}

func (s *service) Create(ctx context.Context, i domain.InboundOrder) (domain.InboundOrder, error) {
	id, err := s.repository.Save(ctx, i)
	if err != nil {
		return domain.InboundOrder{}, err
	}

	i.ID = id
	return i, nil
}
