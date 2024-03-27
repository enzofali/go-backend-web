package inboundorder

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *repositoryMock {
	return &repositoryMock{}
}

func (rm *repositoryMock) Save(ctx context.Context, i domain.InboundOrder) (int, error) {
	args := rm.Called(ctx, i)
	return args.Int(0), args.Error(1)
}

func Test_Service_Create(t *testing.T) {
	ctx := context.Background()

	data := domain.InboundOrder{
		ID:             1,
		OrderDate:      "2006-01-02",
		OrderNumber:    "order1",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}

	t.Run("Create OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		inboundOrder := domain.InboundOrder{
			OrderDate:      "2006-01-02",
			OrderNumber:    "order1",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		repo.On("Save", ctx, inboundOrder).Return(1, nil)

		inboundOrderDB, err := service.Create(ctx, inboundOrder)
		assert.NoError(t, err)
		assert.Equal(t, data, inboundOrderDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Create Error", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		inboundOrder := domain.InboundOrder{
			OrderDate:      "order1",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		repo.On("Save", ctx, inboundOrder).Return(0, ErrEmployeeNotFound)

		inboundOrderDB, err := service.Create(ctx, inboundOrder)
		assert.Error(t, err)
		assert.EqualError(t, ErrEmployeeNotFound, err.Error())
		assert.Empty(t, inboundOrderDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}
