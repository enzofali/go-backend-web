package employee

import (
	"context"
	"errors"
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

func (rm *repositoryMock) GetAll(ctx context.Context) ([]domain.Employee, error) {
	args := rm.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (rm *repositoryMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := rm.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (rm *repositoryMock) Exists(ctx context.Context, cardNumberID string) bool {
	args := rm.Called(ctx, cardNumberID)
	return args.Get(0).(bool)
}

func (rm *repositoryMock) Save(ctx context.Context, e domain.Employee) (int, error) {
	args := rm.Called(ctx, e)
	return args.Get(0).(int), args.Error(1)
}

func (rm *repositoryMock) Update(ctx context.Context, e domain.Employee) error {
	args := rm.Called(ctx, e)
	return args.Error(0)
}

func (rm *repositoryMock) Delete(ctx context.Context, id int) error {
	args := rm.Called(ctx, id)
	return args.Error(0)
}

func (rm *repositoryMock) GetAllInoundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	args := rm.Called(ctx)
	return args.Get(0).([]domain.EmployeeWithInboundOrders), args.Error(1)
}

func (rm *repositoryMock) GetWithInboundOrder(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	args := rm.Called(ctx, id)
	return args.Get(0).(domain.EmployeeWithInboundOrders), args.Error(1)
}

func Test_Service_Create(t *testing.T) {
	ctx := context.Background()

	data := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	t.Run("Create OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		employee := domain.Employee{
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		repo.On("Exists", ctx, employee.CardNumberID).Return(false)
		repo.On("Save", ctx, employee).Return(1, nil)

		employeeDB, err := service.Create(ctx, employee)
		assert.NoError(t, err)
		assert.Equal(t, data, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Create conflict card number", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		employee := domain.Employee{
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		repo.On("Exists", ctx, employee.CardNumberID).Return(true)

		employeeDB, err := service.Create(ctx, employee)
		assert.Error(t, err)
		assert.EqualError(t, ErrExistsCardId, err.Error())
		assert.Empty(t, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Create error in repository.save", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		employee := domain.Employee{
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		repo.On("Exists", ctx, employee.CardNumberID).Return(false)
		repo.On("Save", ctx, employee).Return(0, ErrWarehouseNotfound)

		employeeDB, err := service.Create(ctx, employee)
		assert.Error(t, err)
		assert.EqualError(t, ErrWarehouseNotfound, err.Error())
		assert.Empty(t, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_GetAll(t *testing.T) {
	ctx := context.Background()

	data := []domain.Employee{
		{
			ID:           1,
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		},
		{
			ID:           2,
			CardNumberID: "A13",
			FirstName:    "Jose",
			LastName:     "Gomez",
			WarehouseID:  3,
		},
	}

	t.Run("GetAll OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		repo.On("GetAll", ctx).Return(data, nil)

		employeesDB, err := service.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, data, employeesDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("GetAll error", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := errors.New("Error in DB")
		repo.On("GetAll", ctx).Return([]domain.Employee{}, err)

		employeesDB, err := service.GetAll(ctx)
		assert.Error(t, err)
		assert.EqualError(t, ErrDatabase, err.Error())
		assert.Empty(t, employeesDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_Get(t *testing.T) {
	ctx := context.Background()

	data := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	t.Run("GetByID OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		repo.On("Get", ctx, 1).Return(data, nil)

		employeeDB, err := service.Get(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, data, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Get ID not exists", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := errors.New("error in DB")
		repo.On("Get", ctx, 1).Return(domain.Employee{}, err)

		employeeDB, err := service.Get(ctx, 1)
		assert.Error(t, err)
		assert.EqualError(t, ErrNotFound, err.Error())
		assert.Empty(t, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("Update OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		employee := domain.Employee{
			ID:           3,
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		repo.On("Update", ctx, employee).Return(nil)

		employeeDB, err := service.Update(ctx, employee)
		assert.NoError(t, err)
		assert.Equal(t, employee, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Update error not exists", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		employee := domain.Employee{
			ID:           3,
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		err := errors.New("Error not exists")
		repo.On("Update", ctx, employee).Return(err)

		employeeDB, err := service.Update(ctx, employee)
		assert.Error(t, err)
		assert.Empty(t, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("Delete OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		repo.On("Delete", ctx, 1).Return(nil)

		err := service.Delete(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("Delete Id not exists", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := ErrNotFound
		repo.On("Delete", ctx, 3).Return(err)

		err = service.Delete(ctx, 3)
		assert.Error(t, err)
		assert.EqualError(t, ErrNotFound, err.Error())
	})

	t.Run("Delete error DB", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := ErrDatabase
		repo.On("Delete", ctx, 3).Return(err)

		err = service.Delete(ctx, 3)
		assert.Error(t, err)
		assert.EqualError(t, ErrDatabase, err.Error())
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_GetAllInoundOrders(t *testing.T) {
	ctx := context.Background()

	data := []domain.EmployeeWithInboundOrders{
		{
			ID:                 1,
			CardNumberID:       "A12",
			FirstName:          "Juan",
			LastName:           "Perez",
			WarehouseID:        1,
			InboundOrdersCount: 3,
		},
		{
			ID:                 2,
			CardNumberID:       "A13",
			FirstName:          "Jose",
			LastName:           "Gomez",
			WarehouseID:        3,
			InboundOrdersCount: 1,
		},
	}

	t.Run("GetAllInboundOrders OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		repo.On("GetAllInoundOrders", ctx).Return(data, nil)

		employeesDB, err := service.GetAllInoundOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, data, employeesDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("GetAllInboundOrders Error", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := errors.New("Error in DB")
		repo.On("GetAllInoundOrders", ctx).Return([]domain.EmployeeWithInboundOrders{}, err)

		employeesDB, err := service.GetAllInoundOrders(ctx)
		assert.Error(t, err)
		assert.EqualError(t, ErrDatabase, err.Error())
		assert.Empty(t, employeesDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}

func Test_Service_GetWithInoundOrders(t *testing.T) {
	ctx := context.Background()

	data := domain.EmployeeWithInboundOrders{
		ID:                 1,
		CardNumberID:       "A12",
		FirstName:          "Juan",
		LastName:           "Perez",
		WarehouseID:        1,
		InboundOrdersCount: 3,
	}

	t.Run("GetWithInboundOrders OK", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		repo.On("GetWithInboundOrder", ctx, 1).Return(data, nil)

		employeeDB, err := service.GetWithInboundOrder(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, data, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})

	t.Run("GetWithInboundOrders Not Found", func(t *testing.T) {
		repo := NewRepositoryMock()
		service := NewService(repo)

		err := errors.New("Error not found")
		repo.On("GetWithInboundOrder", ctx, 2).Return(domain.EmployeeWithInboundOrders{}, err)

		employeeDB, err := service.GetWithInboundOrder(ctx, 2)
		assert.Error(t, err)
		assert.EqualError(t, ErrNotFound, err.Error())
		assert.Empty(t, employeeDB)
		assert.True(t, repo.AssertExpectations(t))
	})
}
