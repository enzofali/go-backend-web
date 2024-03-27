package warehouse

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryTest struct {
	mock.Mock
}

//constructor

func NewWarehouseRepository() *repositoryTest {
	return &repositoryTest{}
}

func (r *repositoryTest) Get(ctx context.Context, id int) (domain.Warehouse, error) {

	args := r.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *repositoryTest) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (r *repositoryTest) Exists(ctx context.Context, wareHCode string) bool {
	args := r.Called(ctx, wareHCode)
	return args.Get(0).(bool)
}
func (r *repositoryTest) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	args := r.Called(ctx, w)
	return args.Get(0).(int), args.Error(1)
}
func (r *repositoryTest) Update(ctx context.Context, w domain.Warehouse) error {
	args := r.Called(ctx, w)
	return args.Error(0)
}
func (r *repositoryTest) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func TestGetAllWService(t *testing.T) {

	//preparo mis datos
	ctx := context.Background()

	data := []domain.Warehouse{
		{ID: 4, Address: "Calle 33 # 34-25", Telephone: "2243567", WarehouseCode: "AB201", MinimumCapacity: 10, MinimumTemperature: 15},
		{ID: 5, Address: "Calle 32 # 33-26", Telephone: "2243568", WarehouseCode: "AB202", MinimumCapacity: 15, MinimumTemperature: 20},
		{ID: 7, Address: "Calle 31 # 33-27", Telephone: "2243569", WarehouseCode: "AB203", MinimumCapacity: 20, MinimumTemperature: 25},
	}

	//inicio casos test
	t.Run("find_all", func(t *testing.T) {

		// arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("GetAll", ctx).Return(data, nil)

		// act
		wareH, err := s.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, wareH)
		assert.Equal(t, 3, len(wareH))
		assert.True(t, r.AssertExpectations(t))
	})

}

func TestGetWService(t *testing.T) {

	//preparo mis datos
	ctx := context.Background()
	data := domain.Warehouse{
		ID:                 8,
		Address:            "Calle 33 # 34-28",
		Telephone:          "2243570",
		WarehouseCode:      "AB210",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	//inicio casos test
	t.Run("find_by_id_existent", func(t *testing.T) {
		// arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Get", ctx, data.ID).Return(data, nil)
		id := 8

		// act
		wareH, err := s.Get(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, wareH)
		assert.Equal(t, 8, data.ID)
		assert.True(t, r.AssertExpectations(t))

	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		id := 1
		r.On("Get", ctx, id).Return(domain.Warehouse{}, ErrNotFound)

		// act
		wareH, err := s.Get(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Empty(t, wareH)

	})

}

func TestCreateWService(t *testing.T) {

	//preparo mis datos
	ctx := context.Background()
	data := domain.Warehouse{
		ID:                 8,
		Address:            "Calle 33 # 34-28",
		Telephone:          "2243570",
		WarehouseCode:      "AB210",
		MinimumCapacity:    10,
		MinimumTemperature: 15}

	//inicio casos test

	t.Run("create_ok", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Exists", ctx, data.WarehouseCode).Return(false)
		r.On("Save", ctx, data).Return(data.ID, nil)

		//act
		wareH, err := s.Create(ctx, data)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, data, wareH)
		assert.Equal(t, 8, wareH.ID)
		assert.True(t, r.AssertExpectations(t))

	})

	t.Run("create_conflict", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Exists", ctx, data.WarehouseCode).Return(true)
		//r.On("Save", ctx, data).Return(0, ErrBD)

		//act
		wareH, err := s.Create(ctx, data)

		//assert
		assert.Error(t, err)
		assert.Empty(t, wareH)
		assert.True(t, r.AssertExpectations(t))

	})

	t.Run("create_fail", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Exists", ctx, data.WarehouseCode).Return(false)
		r.On("Save", ctx, data).Return(0, ErrBD)

		//act
		wareH, err := s.Create(ctx, data)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Empty(t, wareH)
		assert.True(t, r.AssertExpectations(t))

	})
}

func TestUpdateWService(t *testing.T) {
	//preparo
	ctx := context.Background()
	data := domain.Warehouse{
		ID:                 8,
		Address:            "Calle 33 # 3-28",
		Telephone:          "2243570",
		WarehouseCode:      "AB210",
		MinimumCapacity:    10,
		MinimumTemperature: 15}

	t.Run("update_ok", func(t *testing.T) {
		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Get", ctx, data.ID).Return(data, nil)
		//verifica si existe en otro lado y es diferente al que tenia anteriormente
		r.On("Exists", ctx, data.WarehouseCode).Return(false)
		r.On("Update", ctx, data).Return(nil)

		//act
		wareH, err := s.Update(ctx, data)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, data, wareH)
		assert.True(t, r.AssertExpectations(t))

	})

	t.Run("update_non_existent", func(t *testing.T) {
		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Get", ctx, data.ID).Return(data, ErrNotFound)
		//verifica si existe en otro lado y es diferente al que tenia anteriormente

		//act
		wareH, err := s.Update(ctx, data)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Empty(t, wareH)
		assert.True(t, r.AssertExpectations(t))

	})

}

func TestDeleteWService(t *testing.T) {
	//preparo

	ctx := context.Background()
	id := 8

	t.Run("delete_ok", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Delete", ctx, id).Return(nil)

		//act
		err := s.Delete(ctx, id)

		//assert
		assert.NoError(t, err)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("delete_non_existent", func(t *testing.T) {

		//arrange
		r := NewWarehouseRepository()
		s := NewService(r)
		r.On("Delete", ctx, id).Return(ErrNotFound)

		//act
		err := s.Delete(ctx, id)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.True(t, r.AssertExpectations(t))
	})

}
