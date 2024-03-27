package product_batches

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

func (r *repositoryTest) Create(ctx context.Context, productBatches domain.ProductBatches) (int, error) {
	args := r.Called(ctx, productBatches)
	return args.Get(0).(int), args.Error(1)
}

func Test_Create_Service(t *testing.T) {
	ctx := context.Background()

	data := domain.ProductBatches{
		ID:                 1,
		BatchNumber:        1234,
		CurrentQuantity:    10,
		CurrentTemperature: 10,
		DueDate:            "2023-02-01",
		InitialQuantity:    5,
		ManufacturingDate:  "2023-01-01",
		ManufacturingHour:  "13:01:06",
		MinumumTemperature: 5,
		ProductID:          1,
		SectionID:          1,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Create", ctx, data).Return(1, nil)

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, productBatches)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("Create: ErrInternal", func(t *testing.T) {
		// arrange
		r := NewRepositoryTest()
		s := NewService(r)
		r.On("Create", ctx, data).Return(0, ErrInternal)

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Empty(t, productBatches)
		assert.True(t, r.AssertExpectations(t))
	})
}
