package purchaseorder

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	createdPurchaseOrder := domain.Purchase_Orders{
		ID:                0,
		Order_number:      "123456",
		Order_date:        "2023/12/12",
		Tracking_code:     "abc2344",
		Buyer_id:          1,
		Product_record_id: 2,
		Order_Status_id:   2,
	}

	rep := NewRepository(db)
	ctx := context.Background()

	query := "INSERT INTO purchase_orders(order_number, order_date, tracking_code, buyer_id, product_record_id, order_status_id) VALUES (?,?,?,?,?,?);"
	t.Run("Save OK", func(t *testing.T) {
		//arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		//act
		id, err := rep.Save(ctx, createdPurchaseOrder)

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrInternal Prepare", func(t *testing.T) {
		//arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(ErrDatabase)

		//act
		id, err := rep.Save(ctx, createdPurchaseOrder)

		assert.Error(t, err)
		assert.Equal(t, ErrDatabase, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrInternal Exec", func(t *testing.T) {
		//arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnError(&mysql.MySQLError{})

		//act
		id, err := rep.Save(ctx, createdPurchaseOrder)

		//assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, ErrDatabase, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Err Internal Exec", func(t *testing.T) {

		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnError(ErrDatabase)

		id, err := rep.Save(ctx, createdPurchaseOrder)

		assert.Error(t, err)
		assert.Equal(t, ErrDatabase, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("ErrExists 1062", func(t *testing.T) {
		//arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		//act
		id, err := rep.Save(ctx, createdPurchaseOrder)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrExists, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrInternal LastInsertId", func(t *testing.T) {
		//arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		//act
		id, err := rep.Save(ctx, createdPurchaseOrder)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, ErrDatabase, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPurchaseExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id"})
	row.AddRow(1)

	query := "SELECT id FROM purchase_orders WHERE id=?"

	rep := NewRepository(db)
	ctx := context.Background()

	//arrange
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(row)

	//act
	value := rep.Exists(ctx, 1)

	//assert
	assert.NoError(t, err)
	assert.Equal(t, true, value)

}

func TestBuyerExists(t *testing.T) {
	//arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id"})
	row.AddRow(1)

	query := "SELECT id FROM buyers WHERE id=?"

	rep := NewRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(row)

	//act
	value := rep.ExistsBuyer(ctx, 1)

	//assert
	assert.NoError(t, err)
	assert.Equal(t, true, value)
	assert.NoError(t, mock.ExpectationsWereMet())
}
