package carry

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrNotFound   = errors.New("carries not found")
	ErrBD         = errors.New("carries is empty")
	ErrNotExist   = errors.New("locality not exist")
	ErrExist      = errors.New("carries already exist")
	ErrForeignKey = errors.New("foreign key does not exist")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Carrie, error)
	GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error)
	GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error)
	Crear(ctx context.Context, c domain.Carrie) (carriesG domain.Carrie, err error)
	//Update(ctx context.Context, c domain.Carrie) (carriesG domain.Carrie, err error)
	//Delete(ctx context.Context, id int) (err error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Carrie, error) {
	l, err := s.r.GetAll(ctx)
	//valite if is empty

	if err != nil {
		return []domain.Carrie{}, ErrBD
	}
	return l, nil
}

func (s *service) GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error) {
	l, err := s.r.GetByLocality(ctx)
	//valite if is empty

	if err != nil {
		return []domain.CarrieLocality{}, ErrBD
	}
	return l, nil
}
func (s *service) GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error) {

	c, err := s.r.GetByLocalityID(ctx, id)
	if err != nil {
		return domain.CarrieLocality{}, ErrNotExist
	}
	return c, nil

}

func (s *service) Crear(ctx context.Context, c domain.Carrie) (carriesG domain.Carrie, err error) {

	// returns an error if cid code is already in use
	if s.r.Exists(ctx, c.Cid) {
		return domain.Carrie{}, ErrExist
	}

	// validar que la localidad exista ////
	if !s.r.ExistsFK(ctx, c.Locality_id) {

		return domain.Carrie{}, ErrForeignKey
	}

	// create carry in database and save its id
	id, er := s.r.Crear(ctx, c)

	if er != nil {

		return domain.Carrie{}, ErrBD

	}
	// assign carry its database-defined  id and return it
	c.Id = id
	return c, nil
}

/*func (s *service) Update(ctx context.Context, c domain.Carrie) (carriesG domain.Carrie, err error) {
	wPre, _ := s.r.Get(ctx, c.Id)
	// Validate code carry
	if s.r.Exists(ctx, c.Cid) && wPre.Cid != c.Cid {
		return domain.Carrie{}, ErrExist
	}

	//update carry in database and return an error
	er := s.r.Update(ctx, c)

	if er != nil {
		return domain.Carrie{}, er
	}

	return c, nil
}
func (s *service) Delete(ctx context.Context, id int) (err error) {
	er := s.r.Delete(ctx, id)
	if er != nil {
		return er
	}
	return nil
}*/
