package seller

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrConflict = errors.New("error conflic")
)

type Service interface {
	GetAll(context.Context) ([]domain.Seller, error)
	GetByID(context.Context, int) (domain.Seller, error)
	Create(context.Context, domain.Seller) (int, error)
	Update(context.Context, domain.Seller) error
	Delete(context.Context, int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Returns all sellers
func (service service) GetAll(ctx context.Context) (sellers []domain.Seller, err error) {
	sellers, err = service.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return
}

// returns seller specified by id parameter
func (service service) GetByID(ctx context.Context, id int) (seller domain.Seller, err error) {
	seller, err = service.repo.Get(ctx, id)
	if err != nil {
		return seller, err
	}
	return
}

// adds one seller to seller list and returns it
func (service service) Create(ctx context.Context, newSeller domain.Seller) (id int, er error) {
	if service.repo.Exists(ctx, newSeller.CID) {
		return 0, ErrConflict
	}
	id, err := service.repo.Save(ctx, newSeller)
	if err != nil {
		return 0, err
	}
	return
}

// updates one seller specified by parameter with new values and returns it
func (service service) Update(ctx context.Context, sel domain.Seller) error {
	seller, _ := service.GetByID(ctx, sel.ID)

	if service.repo.Exists(ctx, sel.CID) && seller.CID != sel.CID {
		return ErrConflict
	}
	err := service.repo.Update(ctx, sel)
	if err != nil {
		return err
	}
	return nil
}

// deletes seller with specified id
func (service service) Delete(ctx context.Context, id int) error {
	err := service.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
