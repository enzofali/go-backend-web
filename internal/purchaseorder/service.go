package purchaseorder

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	//"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internl/purchaseorder"
)

var (
	ErrNotFound      = errors.New("purchase order not found")
	ErrDatabase      = errors.New("database error")
	ErrExists        = errors.New("product order already exists")
	ErrBuyerNotFound = errors.New("buyer not found")
	ErrIdNotExist    = errors.New("id does not exist")
	ErrFieldNotExist = errors.New("field does not exist")
)

type Service interface {
	Create(ctx context.Context, purchOrd domain.Purchase_Orders) (domain.Purchase_Orders, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) Create(ctx context.Context, purchOrd domain.Purchase_Orders) (domain.Purchase_Orders, error) {

	buyer_id := s.r.ExistsBuyer(ctx, purchOrd.Buyer_id)

	if !buyer_id {
		return domain.Purchase_Orders{}, ErrBuyerNotFound
	}

	id, err := s.r.Save(ctx, purchOrd)

	if err != nil {
		return domain.Purchase_Orders{}, ErrDatabase
	}

	purchOrd.ID = id

	return purchOrd, nil
}
