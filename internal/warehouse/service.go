package warehouse

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrNotFound   = errors.New("warehouse not found")
	ErrBD         = errors.New("warehouse is empty")
	ErrExist      = errors.New("Warehouse already exist")
	ErrBadRequest = errors.New("bad request")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Create(ctx context.Context, w domain.Warehouse) (wareH domain.Warehouse, err error)
	//Save(ctx context.Context, w domain.Warehouse) (int, error)
	Update(ctx context.Context, w domain.Warehouse) (wareH domain.Warehouse, err error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

// return all warehouse
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	l, err := s.r.GetAll(ctx)
	//valite if is empty

	if err != nil {
		return []domain.Warehouse{}, ErrBD
	}
	return l, nil
}

// returns product specified by id
func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	w, err := s.r.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}
	return w, nil

}

// adds one warehouse to warehouse list and returns it
func (s *service) Create(ctx context.Context, w domain.Warehouse) (wareH domain.Warehouse, err error) {

	// returns an error if warehouse code is already in use
	if s.r.Exists(ctx, w.WarehouseCode) {
		return domain.Warehouse{}, ErrExist
	}
	// create warehouse in database and save its id
	id, er := s.r.Save(ctx, w)

	if er != nil {
		return domain.Warehouse{}, ErrBD

	}
	// assign warehouse its database-defined  id and return it
	w.ID = id
	return w, nil
}

// updates one warehouse specified by parameter with new values
func (s *service) Update(ctx context.Context, w domain.Warehouse) (wareH domain.Warehouse, err error) {
	wPre, err := s.r.Get(ctx, w.ID)

	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}
	// Validate code warehouse si existe en otro lado y es diferente al que tenia anteriormente
	if s.r.Exists(ctx, w.WarehouseCode) && wPre.WarehouseCode != w.WarehouseCode {
		return domain.Warehouse{}, ErrExist
	}

	//update warehouse in database and return an error
	er := s.r.Update(ctx, w)

	if er != nil {
		return domain.Warehouse{}, ErrBD
	}

	return w, nil
}

// // deletes warehouse  with specified id
func (s *service) Delete(ctx context.Context, id int) (err error) {
	er := s.r.Delete(ctx, id)

	if er == ErrNotFound {
		return ErrNotFound
	}

	if er != nil {
		return ErrBD
	}
	return nil
}
