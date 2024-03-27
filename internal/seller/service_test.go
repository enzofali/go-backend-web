package seller

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// define a mock repository struct that implements the Repository interface for testing purposes
type RepositoryMock struct {
	mock.Mock
	Repository
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Seller, error) {
	args := r.Mock.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}
func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := r.Mock.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}
func (r *RepositoryMock) Exists(ctx context.Context, cid int) bool {
	args := r.Mock.Called(ctx, cid)
	return args.Get(0).(bool)
}
func (r *RepositoryMock) Save(ctx context.Context, s domain.Seller) (int, error) {
	args := r.Mock.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}
func (r *RepositoryMock) Update(ctx context.Context, s domain.Seller) error {
	args := r.Mock.Called(ctx, s)
	return args.Error(0)
}
func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Mock.Called(ctx, id)
	return args.Error(0)
}

func Test_GetAll_Seller(t *testing.T) {
	repoMock := NewRepositoryMock()
	service := NewService(repoMock)
	ctx := context.Background()
	sellersExpected := []domain.Seller{
		{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"},
		{ID: 2, CID: 2, CompanyName: "Digital House", Address: "Monroe 860", Telephone: "47470000", Locality_id: "6700"},
	}

	t.Run("OK", func(t *testing.T) {
		//arrange
		repoMock.On("GetAll", ctx).Return(sellersExpected, nil)

		// act
		sellers, err := service.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(sellers))
		assert.Equal(t, sellersExpected, sellers)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)
		repoMock.On("GetAll", ctx).Return([]domain.Seller{}, ErrIntern)

		// act
		sellers, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.Nil(t, sellers)
		assert.True(t, repoMock.AssertExpectations(t))
	})

}

func Test_GetByID(t *testing.T) {
	sellerexpexted := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		//arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)
		id := 1

		repoMock.On("Get", ctx, id).Return(sellerexpexted, nil)

		//act
		seller, err := service.GetByID(ctx, id)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, sellerexpexted, seller)
	})

	t.Run("Not found error", func(t *testing.T) {
		//arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)
		id := 1

		repoMock.On("Get", ctx, id).Return(domain.Seller{}, ErrNotFound)

		//act
		seller, err := service.GetByID(ctx, id)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Empty(t, domain.Seller{}, seller)
		assert.True(t, repoMock.AssertExpectations(t))
	})
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	sellerToCreate := domain.Seller{CID: 3, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}

	t.Run("OK", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Exists", ctx, sellerToCreate.CID).Return(false)
		repoMock.On("Save", ctx, sellerToCreate).Return(3, nil)

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Conflict error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Exists", ctx, sellerToCreate.CID).Return(true)

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrConflict)
		assert.Equal(t, 0, id)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Exists", ctx, sellerToCreate.CID).Return(false)
		repoMock.On("Save", ctx, sellerToCreate).Return(0, ErrIntern)

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrIntern)
		assert.Equal(t, 0, id)
		assert.True(t, repoMock.AssertExpectations(t))
	})
}

func Test_UpdateSeller(t *testing.T) {
	ctx := context.Background()
	sellerDb := domain.Seller{ID: 2, CID: 2, CompanyName: "Digital House", Address: "Monroe 860", Telephone: "47470000", Locality_id: "6700"}
	sellerToUpdate := domain.Seller{ID: 2, CID: 3, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Get", ctx, 2).Return(sellerDb, nil)
		repoMock.On("Exists", ctx, sellerToUpdate.CID).Return(false)
		repoMock.On("Update", ctx, sellerToUpdate).Return(nil)

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.NoError(t, err)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Conflict error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Get", ctx, 2).Return(sellerDb, nil)
		repoMock.On("Exists", ctx, sellerToUpdate.CID).Return(true)

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrConflict, err)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Not found error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Get", ctx, 2).Return(sellerDb, nil)
		repoMock.On("Exists", ctx, sellerToUpdate.CID).Return(false)
		repoMock.On("Update", ctx, sellerToUpdate).Return(ErrNotFound)

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.True(t, repoMock.AssertExpectations(t))
	})
}

func Test_DeleteSeller(t *testing.T) {
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Delete", ctx, 2).Return(nil)

		// act
		err := service.Delete(ctx, 2)

		// assert
		assert.NoError(t, err)
		assert.True(t, repoMock.AssertExpectations(t))
	})

	t.Run("Not found error", func(t *testing.T) {
		// arrange
		repoMock := NewRepositoryMock()
		service := NewService(repoMock)

		repoMock.On("Delete", ctx, 2).Return(ErrNotFound)

		// act
		err := service.Delete(ctx, 2)

		// assert
		assert.Error(t, err)
		assert.EqualError(t, ErrNotFound, err.Error())
		assert.True(t, repoMock.AssertExpectations(t))
	})
}
