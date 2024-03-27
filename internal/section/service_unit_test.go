package section

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Controller
type repositoryTest struct {
	mock.Mock
}

// Constructor
func NewRepositoryTest() *repositoryTest {
	return &repositoryTest{}
}

func (r *repositoryTest) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}
func (r *repositoryTest) GetByID(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *repositoryTest) GetAllReportProducts(ctx context.Context) ([]domain.SectionReportProducts, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.SectionReportProducts), args.Error(1)
}

func (r *repositoryTest) GetReportProductsByID(ctx context.Context, id int) ([]domain.SectionReportProducts, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]domain.SectionReportProducts), args.Error(1)
}
func (r *repositoryTest) Create(ctx context.Context, section domain.Section) (int, error) {
	args := r.Called(ctx, section)
	return args.Get(0).(int), args.Error(1)
}
func (r *repositoryTest) Update(ctx context.Context, section domain.Section) error {
	args := r.Called(ctx, section)
	return args.Error(0)
}
func (r *repositoryTest) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

// ------------------------------- READ ---------------------------------

func Test_GetAll_Service(t *testing.T) {
	ctx := context.Background()

	data := []domain.Section{
		{ID: 1, SectionNumber: 1, CurrentTemperature: 15, MinimumTemperature: -20, CurrentCapacity: 20, MinimumCapacity: 5, MaximumCapacity: 50, WarehouseID: 1, ProductTypeID: 1},
		{ID: 2, SectionNumber: 3, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1},
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("GetAll", ctx).Return(data, nil)

		// act
		sections, err := s.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, sections)
		assert.True(t, r.AssertExpectations(t))
	})
	t.Run("GetAll: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("GetAll", ctx).Return([]domain.Section{}, ErrInternal)

		// act
		sections, err := s.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Empty(t, sections)
		assert.True(t, r.AssertExpectations(t))
	})
}

func Test_GetByID_Service(t *testing.T) {
	ctx := context.Background()

	data := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 15,
		MinimumTemperature: -20,
		CurrentCapacity:    20,
		MinimumCapacity:    5,
		MaximumCapacity:    50,
		WarehouseID:        1,
		ProductTypeID:      1}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 1
		r.On("GetByID", ctx, id).Return(data, nil)

		// act
		section, err := s.GetByID(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, section)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("GetByID: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 1
		r.On("GetByID", ctx, id).Return(domain.Section{}, ErrInternal)

		// act
		section, err := s.GetByID(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Empty(t, section)
		assert.True(t, r.AssertExpectations(t))
	})

}

func Test_GetReportProducts_Service(t *testing.T) {
	ctx := context.Background()

	reports := []domain.SectionReportProducts{
		{ID: 1, SectionNumber: 1, ProductCount: 150},
		{ID: 2, SectionNumber: 3, ProductCount: 250},
	}

	report := []domain.SectionReportProducts{
		{ID: 1, SectionNumber: 1, ProductCount: 150},
	}

	t.Run("Ok: GetAllReportProducts", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 0
		r.On("GetAllReportProducts", ctx).Return(reports, nil)

		// act
		result, err := s.GetReportProducts(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, reports, result)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("GetAllReportProducts: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 0
		r.On("GetAllReportProducts", ctx).Return([]domain.SectionReportProducts{}, ErrInternal)

		// act
		result, err := s.GetReportProducts(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("Ok: GetReportProductsByID", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 1
		r.On("GetReportProductsByID", ctx, id).Return(report, nil)

		// act
		result, err := s.GetReportProducts(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, report, result)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("GetReportProductsByID: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		id := 1
		r.On("GetReportProductsByID", ctx, id).Return([]domain.SectionReportProducts{}, ErrInternal)

		// act
		result, err := s.GetReportProducts(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.True(t, r.AssertExpectations(t))
	})
}

// -------------------------------- WRITE --------------------------------

func Test_Create_Service(t *testing.T) {
	ctx := context.Background()

	data := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 15,
		MinimumTemperature: -20,
		CurrentCapacity:    20,
		MinimumCapacity:    5,
		MaximumCapacity:    50,
		WarehouseID:        1,
		ProductTypeID:      1}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Create", ctx, data).Return(1, nil)

		// act
		section, err := s.Create(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, section)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("Create: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Create", ctx, data).Return(0, ErrInternal)

		// act
		section, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Empty(t, section)
		assert.True(t, r.AssertExpectations(t))
	})

}

func Test_Update_Service(t *testing.T) {
	ctx := context.Background()

	data := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 15,
		MinimumTemperature: -20,
		CurrentCapacity:    20,
		MinimumCapacity:    5,
		MaximumCapacity:    50,
		WarehouseID:        1,
		ProductTypeID:      1}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Update", ctx, data).Return(nil)

		// act
		err := s.Update(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("Update: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Update", ctx, data).Return(ErrInternal)

		// act
		err := s.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.True(t, r.AssertExpectations(t))
	})

}

func Test_Delete_Service(t *testing.T) {
	ctx := context.Background()
	id := 1
	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Delete", ctx, id).Return(nil)

		// act
		err := s.Delete(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("Delete: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Delete", ctx, id).Return(ErrInternal)

		// act
		err := s.Delete(ctx, id)

		// assert
		assert.Error(t, err)
		assert.True(t, r.AssertExpectations(t))
	})

}
