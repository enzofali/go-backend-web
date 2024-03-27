package employee

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrNotFound     = errors.New("employee not found")
	ErrDatabase     = errors.New("Database error.")
	ErrExistsCardId = errors.New("Exists card number id.")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Get(ctx context.Context, id int) (domain.Employee, error)
	Create(ctx context.Context, e domain.Employee) (domain.Employee, error)
	Update(ctx context.Context, e domain.Employee) (domain.Employee, error)
	Delete(ctx context.Context, id int) error
	GetAllInoundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error)
	GetWithInboundOrder(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error)
}

type service struct {
	repository Repository
}

func NewService(employee Repository) Service {
	return &service{
		repository: employee,
	}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Employee, error) {
	employees, err := s.repository.GetAll(ctx)
	if err != nil {
		return employees, ErrDatabase
	}
	return employees, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Employee, error) {
	employee, err := s.repository.Get(ctx, id)
	if err != nil {
		return employee, ErrNotFound
	}

	return employee, nil
}

func (s *service) Create(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	if s.repository.Exists(ctx, e.CardNumberID) {
		return domain.Employee{}, ErrExistsCardId
	}

	id, err := s.repository.Save(ctx, e)
	if err != nil {
		return domain.Employee{}, err
	}

	e.ID = id
	return e, nil
}

func (s *service) Update(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	err := s.repository.Update(ctx, e)
	if err != nil {
		return domain.Employee{}, err
	}

	return e, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (s *service) GetAllInoundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	employees, err := s.repository.GetAllInoundOrders(ctx)
	if err != nil {
		return employees, ErrDatabase
	}
	return employees, nil
}

func (s *service) GetWithInboundOrder(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	employee, err := s.repository.GetWithInboundOrder(ctx, id)
	if err != nil {
		return employee, ErrNotFound
	}

	return employee, nil
}
