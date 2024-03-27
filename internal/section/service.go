package section

import (
	"context"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	GetByID(ctx context.Context, id int) (domain.Section, error)
	GetReportProducts(ctx context.Context, id int) ([]domain.SectionReportProducts, error)
	Create(ctx context.Context, section domain.Section) (domain.Section, error)
	Update(ctx context.Context, section domain.Section) error
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

// ------------------------------- READ ---------------------------------

func (s *service) GetAll(ctx context.Context) ([]domain.Section, error) {
	sections, err := s.r.GetAll(ctx)
	if err != nil {
		return []domain.Section{}, err
	}
	return sections, nil
}

func (s *service) GetByID(ctx context.Context, id int) (domain.Section, error) {
	section, err := s.r.GetByID(ctx, id)
	if err != nil {
		return domain.Section{}, err
	}
	return section, nil
}

func (s *service) GetReportProducts(ctx context.Context, id int) ([]domain.SectionReportProducts, error) {
	var report []domain.SectionReportProducts
	var err error
	if id == 0 {
		report, err = s.r.GetAllReportProducts(ctx)
		if err != nil {
			return []domain.SectionReportProducts{}, err
		}
	} else {
		report, err = s.r.GetReportProductsByID(ctx, id)
		if err != nil {
			return []domain.SectionReportProducts{}, err
		}
	}
	return report, nil
}

// -------------------------------- WRITE --------------------------------

func (s *service) Create(ctx context.Context, section domain.Section) (domain.Section, error) {
	id, err := s.r.Create(ctx, section)
	if err != nil {
		return domain.Section{}, err
	}

	section.ID = id

	return section, nil
}

func (s *service) Update(ctx context.Context, section domain.Section) error {
	// Validate unique section_numeber
	err := s.r.Update(ctx, section)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.r.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
