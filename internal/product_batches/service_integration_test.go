package product_batches

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

func Benchmark_Create_Service_Integration(b *testing.B) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)
	ctx := context.Background()

	data := domain.ProductBatches{
		ID:                 1,
		BatchNumber:        1234,
		CurrentQuantity:    10,
		CurrentTemperature: 10,
		DueDate:            "2023-02-01",
		InitialQuantity:    5,
		ManufacturingDate:  "2023-01-01",
		ManufacturingHour:  "13:01:06",
		MinumumTemperature: 5,
		ProductID:          1,
		SectionID:          1,
	}

	for i := 0; i < b.N; i++ {
		s.Create(ctx, data)
	}
}

func Test_Create_Service_Integration(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)
	ctx := context.Background()

	data := domain.ProductBatches{
		ID:                 1,
		BatchNumber:        1234,
		CurrentQuantity:    10,
		CurrentTemperature: 10,
		DueDate:            "2023-02-01",
		InitialQuantity:    5,
		ManufacturingDate:  "2023-01-01",
		ManufacturingHour:  "13:01:06",
		MinumumTemperature: 5,
		ProductID:          1,
		SectionID:          1,
	}

	query := "INSERT INTO products_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, data, productBatches)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrProductNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`products`"})

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrProductNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrSectionNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`sections`"})

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrSectionNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrExistsBatchNumber", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrExistsBatchNumber, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RowsAffected: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 0))

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LastInsertId: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		// act
		productBatches, err := s.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, domain.ProductBatches{}, productBatches)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
