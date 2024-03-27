package carry

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryTestCase struct {
	mock.Mock
}

//constructor

func NewCarryRepositoryTestCase() *repositoryTestCase {

	return &repositoryTestCase{}
}

// All methods simply return the struct's initially defined members
func (r repositoryTestCase) GetAll(ctx context.Context) ([]domain.Carrie, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Carrie), args.Error(1)
}
func (r repositoryTestCase) GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.CarrieLocality), args.Error(1)
}

func (r repositoryTestCase) GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.CarrieLocality), args.Error(1)
}

func (r repositoryTestCase) Crear(ctx context.Context, c domain.Carrie) (int, error) {
	args := r.Called(ctx, c)
	return args.Get(0).(int), args.Error(1)
}

func (r repositoryTestCase) Exists(ctx context.Context, carrieCode string) bool {
	args := r.Called(ctx, carrieCode)
	return args.Get(0).(bool)
}
func (r repositoryTestCase) ExistsFK(ctx context.Context, localityCode string) bool {
	args := r.Called(ctx, localityCode)
	return args.Get(0).(bool)
}

func TestGetAllCarry(t *testing.T) {

	//Prepacion datos de testes

	ctx := context.Background()
	data := []domain.Carrie{
		{Id: 7, Cid: "ABC23", Company_name: "Servientrega", Address: "cra 40 # 38-45", Telephone: "2245678", Locality_id: "LA02"},
		{Id: 8, Cid: "ABC24", Company_name: "DHL", Address: "cra 41 # 39-46", Telephone: "2245679", Locality_id: "LA03"},
		{Id: 9, Cid: "ABC25", Company_name: "Envia", Address: "cra 42 # 40-47", Telephone: "2245680", Locality_id: "LA04"},
	}

	t.Run("find_all", func(t *testing.T) {

		//arrange
		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("GetAll", ctx).Return(data, nil)

		//act

		carryM, err := s.GetAll(ctx)

		//assert

		assert.NoError(t, err)
		assert.Equal(t, data, carryM)
		assert.Equal(t, 3, len(carryM))
		assert.True(t, r.AssertExpectations(t))

	})
}

func TestGetLocalityCarry(t *testing.T) {

	//Prepacion datos de testes

	ctx := context.Background()
	data := []domain.CarrieLocality{
		{Locality_id: "B23", Locality_name: "Armenia", Cant_carries: 2},
		{Locality_id: "B23", Locality_name: "Shangai", Cant_carries: 2},
		{Locality_id: "B22", Locality_name: "Armenia", Cant_carries: 1},
	}

	t.Run("find_all", func(t *testing.T) {

		//arrange
		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("GetByLocality", ctx).Return(data, nil)

		//act

		carryM, err := s.GetByLocality(ctx)

		//assert

		assert.NoError(t, err)
		assert.Equal(t, data, carryM)
		assert.Equal(t, 3, len(carryM))
		assert.True(t, r.AssertExpectations(t))

	})
}

func TestGetByCLocalityID(t *testing.T) {

	//preparo datos

	ctx := context.Background()
	data := domain.CarrieLocality{

		Locality_id:   "LA05",
		Locality_name: "Andalucia",
		Cant_carries:  5,
	}

	t.Run("find_by_id_existent", func(t *testing.T) {

		//arrange

		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("GetByLocalityID", ctx, data.Locality_id).Return(data, nil)
		id := "LA05"

		//act
		carryM, err := s.GetByLocalityID(ctx, id)

		//assert

		assert.NoError(t, err)
		assert.Equal(t, data, carryM)
		assert.Equal(t, "LA05", data.Locality_id)
		assert.Equal(t, 5, data.Cant_carries)
		assert.True(t, r.AssertExpectations(t))
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {

		//arrange

		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("GetByLocalityID", ctx, data.Locality_id).Return(domain.CarrieLocality{}, ErrNotExist)
		id := "LA05"

		//act
		carryM, err := s.GetByLocalityID(ctx, id)

		//assert

		assert.Error(t, err)
		assert.Equal(t, ErrNotExist, err)
		assert.Empty(t, carryM)

	})
}

func TestCreateWService(t *testing.T) {

	//preparo mis datos
	ctx := context.Background()

	data := domain.Carrie{
		Id:           7,
		Cid:          "ABC23",
		Company_name: "Servientrega",
		Address:      "cra 40 # 38-45",
		Telephone:    "2245678",
		Locality_id:  "LA02"}

	//inicio casos test

	t.Run("create_ok", func(t *testing.T) {

		//arrange
		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("Exists", ctx, data.Cid).Return(false)
		r.On("ExistsFK", ctx, data.Locality_id).Return(true)
		r.On("Crear", ctx, data).Return(data.Id, nil)

		//act
		carryM, err := s.Crear(ctx, data)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, data, carryM)
		assert.Equal(t, 7, carryM.Id)
		assert.True(t, r.AssertExpectations(t))

	})

	t.Run("create_conflict", func(t *testing.T) {

		//arrange
		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("Exists", ctx, data.Cid).Return(false)
		r.On("ExistsFK", ctx, data.Locality_id).Return(false)

		//act
		carryM, err := s.Crear(ctx, data)

		//assert
		assert.Error(t, err)
		assert.Empty(t, carryM)
		assert.True(t, r.AssertExpectations(t))

	})

	//revisar por que falla y en warehouse si pasa

	t.Run("create_fail", func(t *testing.T) {

		//arrange
		r := NewCarryRepositoryTestCase()
		s := NewService(r)
		r.On("Exists", ctx, data.Cid).Return(true)

		//act
		carryM, err := s.Crear(ctx, data)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrExist, err)
		assert.Empty(t, carryM)
		assert.True(t, r.AssertExpectations(t))

	})
}
