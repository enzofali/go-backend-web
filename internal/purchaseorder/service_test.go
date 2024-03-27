package purchaseorder

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// controller
type repositoryPurchaseOrderTestUnit struct {
	mock.Mock
	Repository
}

// constructor
func NewRepositoryPurchaseOrderTestUnit() *repositoryPurchaseOrderTestUnit {
	return &repositoryPurchaseOrderTestUnit{}
}

func (r *repositoryPurchaseOrderTestUnit) Save(ctx context.Context, purchOrd domain.Purchase_Orders) (int, error) {
	args := r.Mock.Called(ctx, purchOrd)
	return args.Get(0).(int), args.Error(1)
}

func (r *repositoryPurchaseOrderTestUnit) Exist(ctx context.Context, id int) bool {
	args := r.Mock.Called(ctx, id)
	return args.Bool(0)
}

func (r *repositoryPurchaseOrderTestUnit) ExistsBuyer(ctx context.Context, id int) bool {
	args := r.Mock.Called(ctx, id)
	return args.Bool(0)
}

func TestCreate(t *testing.T) {

	newPurchaseOrder := domain.Purchase_Orders{
		ID:                1,
		Order_number:      "abc123",
		Order_date:        "2023/12/12",
		Tracking_code:     "asd4321",
		Buyer_id:          1,
		Product_record_id: 12,
		Order_Status_id:   2,
	}

	ctx := context.Background()

	t.Run("Purchase Order created successfully", func(t *testing.T) {
		//arrange
		repoMockPurchOrd := NewRepositoryPurchaseOrderTestUnit()
		serv := NewService(repoMockPurchOrd)

		repoMockPurchOrd.On("ExistsBuyer", ctx, newPurchaseOrder.Buyer_id).Return(true)
		repoMockPurchOrd.On("Save", ctx, newPurchaseOrder).Return(1, nil)

		//act
		id, err := serv.Create(ctx, newPurchaseOrder)

		//arrange
		assert.NoError(t, err)
		assert.Equal(t, newPurchaseOrder, id)
		assert.True(t, repoMockPurchOrd.AssertExpectations(t))
	})

	t.Run("Buyer not exists", func(t *testing.T) {
		//arrange
		repoMockPurchOrd := NewRepositoryPurchaseOrderTestUnit()
		serv := NewService(repoMockPurchOrd)

		repoMockPurchOrd.On("ExistsBuyer", ctx, newPurchaseOrder.Buyer_id).Return(false)

		//act
		id, err := serv.Create(ctx, newPurchaseOrder)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrBuyerNotFound, err)
		assert.NotEqual(t, newPurchaseOrder.Buyer_id, id)
		assert.True(t, repoMockPurchOrd.AssertExpectations(t))
	})

	t.Run("Database error", func(t *testing.T) {
		repoMockPurchOrd := NewRepositoryPurchaseOrderTestUnit()
		serv := NewService(repoMockPurchOrd)

		repoMockPurchOrd.On("ExistsBuyer", ctx, newPurchaseOrder.Buyer_id).Return(true)
		repoMockPurchOrd.On("Save", ctx, newPurchaseOrder).Return(0, ErrDatabase)

		//act
		id, err := serv.Create(ctx, newPurchaseOrder)

		//arrange
		assert.Error(t, err)
		assert.Equal(t, ErrDatabase, err)
		assert.Equal(t, domain.Purchase_Orders{}, id)
		assert.True(t, repoMockPurchOrd.AssertExpectations(t))
	})
}
