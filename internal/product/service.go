package product

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("product not found")
	ErrDatabase = errors.New("database error")
	ErrExists   = errors.New("product code already exists")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id int) (domain.Product, error)
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	Update(ctx context.Context, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id int) error
	ValidateProductID(ctx context.Context, pid int) bool
	GetOneReport(ctx context.Context, id int) ([]domain.Report, error)
	GetAllReports(ctx context.Context) ([]domain.Report, error)
	CreateType(ctx context.Context, name string) (int, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// returns all products
func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	products, err := s.r.GetAll(ctx)
	if err != nil {
		return []domain.Product{}, ErrDatabase
	}
	return products, nil
}

// returns product specified by id parameter
func (s *service) GetByID(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.r.Get(ctx, id)
	if err != nil {
		return domain.Product{}, ErrNotFound
	}
	return product, nil
}

// adds one product to product list and returns it
func (s *service) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	// returns an error if provided product code is already in use
	if s.r.Exists(ctx, p.ProductCode) {
		return domain.Product{}, ErrExists
	}
	// create product in database and save its id
	id, err := s.r.Save(ctx, p)
	if err != nil {
		return domain.Product{}, ErrDatabase
	}
	// assign product its database-defined id and return it
	p.ID = id
	return p, nil
}

// updates one product specified by parameter with new values and returns it
func (s *service) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	// compare product's old product code with its new one; if it changes, validate that it isn't already in use
	product, _ := s.GetByID(ctx, p.ID)
	if product.ProductCode != p.ProductCode && s.r.Exists(ctx, p.ProductCode) {
		return domain.Product{}, ErrExists
	}
	// update product in database and return it or an error
	err := s.r.Update(ctx, p)
	if err != nil {
		return domain.Product{}, ErrDatabase
	}
	return p, nil
}

// deletes product with specified id
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.r.Delete(ctx, id)
	switch {
	case err == ErrNotFound:
		return ErrNotFound
	case err != nil:
		return ErrDatabase
	}
	return nil
}

// creates a product type with given name and returns its id
func (s *service) CreateType(ctx context.Context, name string) (int, error) {
	id, err := s.r.StoreType(ctx, name)
	if err != nil {
		return 0, ErrDatabase
	}
	return id, nil
}

// returns true if a product with id equal to pid exists in the database, false otherwise
func (s *service) ValidateProductID(ctx context.Context, pid int) bool {
	return s.r.ValidateProductID(ctx, pid)
}

// returns the amount of records and product description associated with given product id
func (s *service) GetOneReport(ctx context.Context, id int) ([]domain.Report, error) {
	count, description, err := s.r.GetOneReport(ctx, id)
	if err != nil {
		return []domain.Report{}, ErrDatabase
	}
	record := domain.Report{Count: count, Description: description, ProductID: id}

	return []domain.Report{record}, nil
}

// returns the amount of records for each product in the database
func (s *service) GetAllReports(ctx context.Context) ([]domain.Report, error) {
	reports, err := s.r.GetAllReports(ctx)
	if err != nil {
		return []domain.Report{}, ErrDatabase
	}

	return reports, nil
}
