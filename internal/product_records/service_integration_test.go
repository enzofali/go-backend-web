package product_records

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

var testPR = domain.ProductRecord{
	LastUpdateDate: "12-12-2012",
	PurchasePrice:  123.123,
	SalePrice:      321.321,
	ProductID:      12,
}

func TestIntegration_ValidateProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)
	t.Run("True", func(t *testing.T) {
		// Arrange
		row := mock.NewRows([]string{"COUNT(id)"})
		row = row.AddRow(1) // must be exactly 1, so this should return true
		mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).ExpectQuery().WithArgs(12).WillReturnRows(row)

		// Act
		valid := s.ValidateProductID(context.Background(), 12)

		// Assert
		assert.True(t, valid)
	})
	t.Run("False", func(t *testing.T) {
		// Arrange
		row := mock.NewRows([]string{"COUNT(id)"})
		row = row.AddRow(0) // must be exactly 1, so this should return false
		mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).ExpectQuery().WithArgs(12).WillReturnRows(row)

		// Act
		valid := s.ValidateProductID(context.Background(), 12)

		// Assert
		assert.False(t, valid)
	})
}

func TestIntegration_Create(t *testing.T) {
	// Arrange
	expected := domain.ProductRecord{
		ID:             5,
		LastUpdateDate: "12-12-2012",
		PurchasePrice:  123.123,
		SalePrice:      321.321,
		ProductID:      12,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(5, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE)).ExpectExec().WillReturnResult(res)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	// should return expected, nil
	pr, err := s.Create(context.Background(), testPR)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, pr)
}
