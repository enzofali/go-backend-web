package purchaseorder

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	query := "INSERT INTO purchase_orders(order_number, order_date, tracking_code, buyer_id, product_record_id, order_status_id) VALUES (?,?,?,?,?,?);"
	queryExist := "SELECT id FROM buyers WHERE id=?"

	purchOrder := domain.Purchase_Orders{
		Order_number:      "12345",
		Order_date:        "2023/12/12",
		Tracking_code:     "abc1234",
		Buyer_id:          1,
		Product_record_id: 2,
		Order_Status_id:   2,
	}

	t.Run("Purchase Order Created", func(t *testing.T) {

		row := mock.NewRows([]string{"buyer_id"})
		row.AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(queryExist)).WillReturnRows(row)
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		newPurchOrd, err := service.Create(ctx, purchOrder)

		assert.NoError(t, err)
		assert.NotNil(t, newPurchOrd)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Purchase Order Already Exist", func(t *testing.T) {

		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnError(ErrBuyerNotFound)

		newPurchOrd, err := service.Create(ctx, purchOrder)

		assert.Error(t, err)
		assert.Equal(t, domain.Purchase_Orders{}, newPurchOrd)
		assert.Equal(t, err, ErrBuyerNotFound)
		assert.Error(t, ErrBuyerNotFound, err)
	})

}
