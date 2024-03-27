package inboundorder

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Service_Create(t *testing.T) {
	ctx := context.Background()

	query := `INSERT INTO inbound_orders(order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?,?,?,?,?)`

	data := domain.InboundOrder{
		ID:             1,
		OrderDate:      "2006-01-02",
		OrderNumber:    "12",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}

	t.Run("Create OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))

		inboundOrder := domain.InboundOrder{
			OrderDate:      "2006-01-02",
			OrderNumber:    "12",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		repo := NewRepository(db)
		service := NewService(repo)
		inboundOrderDB, err := service.Create(ctx, inboundOrder)
		assert.NoError(t, err)
		assert.Equal(t, data, inboundOrderDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "employees"})

		inboundOrder := domain.InboundOrder{
			OrderDate:      "2006-01-02",
			OrderNumber:    "12",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		repo := NewRepository(db)
		service := NewService(repo)
		inboundOrderDB, err := service.Create(ctx, inboundOrder)
		assert.Error(t, err)
		assert.Empty(t, inboundOrderDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
