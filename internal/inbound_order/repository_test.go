package inboundorder

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

func Test_Repository_Save(t *testing.T) {
	ctx := context.Background()

	query := `INSERT INTO inbound_orders(order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?,?,?,?,?)`

	inboundOrder := domain.InboundOrder{
		OrderDate:      "2006-01-02",
		OrderNumber:    "12",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}

	t.Run("Save OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.NoError(t, err)
		assert.Equal(t, 1, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Employee Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "employees"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, ErrEmployeeNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error ProductBatch Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "products_batches"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, ErrProductBatchNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Warehouse Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "warehouses"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, ErrWarehouseNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error CardNumber Exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Order number exists"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, ErrOrderNumberExtists, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1600, Message: "employees"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error LastInsertId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
