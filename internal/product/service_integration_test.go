package product

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_GetAll(t *testing.T) {
	// Arrange
	expected := []domain.Product{
		{ID: 12, Description: "pepe", ExpirationRate: 10, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "unique", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2},
		{ID: 14, Description: "lele", ExpirationRate: 10, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "other", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2},
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	row.AddRow(12, "pepe", 10, 11, 9.1, 0.5, 100.2, "unique", -10, 1, 1, 2)
	row.AddRow(14, "lele", 10, 11, 9.1, 0.5, 100.2, "other", -10, 1, 1, 2)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL)).ExpectQuery().WillReturnRows(row)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	products, err := s.GetAll(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, products)
}

func TestIntegration_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)

	t.Run("OK", func(t *testing.T) {
		// Arrange
		expected := domain.Product{ID: 12, Description: "pepe", ExpirationRate: 10, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "unique", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2}

		row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
		row.AddRow(12, "pepe", 10, 11, 9.1, 0.5, 100.2, "unique", -10, 1, 1, 2)
		mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery().WithArgs(12).WillReturnRows(row)

		// Act
		product, err := s.GetByID(context.Background(), 12)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, product)
	})
	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery()

		// Act
		product, err := s.GetByID(context.Background(), 12)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, domain.Product{}, product)
	})
}

func TestIntegration_Create(t *testing.T) {
	// Arrange
	expected := domain.Product{ID: 12, Description: "pepe", ExpirationRate: 10, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "unique", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2}
	create := domain.Product{Description: "pepe", ExpirationRate: 10, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "unique", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(SAVE)).ExpectExec().WillReturnResult(res)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	product, err := s.Create(context.Background(), create)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, product)
}

func TestIntegration_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)

	t.Run("OK", func(t *testing.T) {
		// Arrange
		expected := domain.Product{ID: 12, Description: "pepe", ExpirationRate: 11, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "unique", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2}

		row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
		row.AddRow(12, "pepe", 10, 11, 9.1, 0.5, 100.2, "unique", -10, 1, 1, 2)
		mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery().WithArgs(12).WillReturnRows(row)

		res := sqlmock.NewResult(12, 1)
		mock.ExpectPrepare(regexp.QuoteMeta(UPDATE)).ExpectExec().WillReturnResult(res)

		// Act
		product, err := s.Update(context.Background(), expected)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, product)
	})
	t.Run("Exists", func(t *testing.T) {
		// Arrange
		differentProductCode := domain.Product{ID: 12, Description: "pepe", ExpirationRate: 11, FreezingRate: 11, Height: 9.1, Length: 0.5, Netweight: 100.2, ProductCode: "different", RecomFreezTemp: -10, Width: 1, ProductTypeID: 1, SellerID: 2}

		rowGet := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
		rowGet.AddRow(12, "pepe", 10, 11, 9.1, 0.5, 100.2, "unique", -10, 1, 1, 2)
		mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery().WithArgs(12).WillReturnRows(rowGet)

		rowExists := mock.NewRows([]string{"product_code"})
		rowExists.AddRow("different")
		mock.ExpectQuery(regexp.QuoteMeta(EXISTS)).WithArgs("different").WillReturnRows(rowExists)

		// Act
		product, err := s.Update(context.Background(), differentProductCode)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrExists, err)
		assert.Equal(t, domain.Product{}, product)
	})
}

func TestIntegration_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := NewRepository(db)
	s := NewService(r)

	t.Run("OK", func(t *testing.T) {
		// Arrange
		res := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnResult(res)

		// Act
		err := s.Delete(context.Background(), 12)

		// Assert
		assert.NoError(t, err)
	})
	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		res := sqlmock.NewResult(0, 0)
		mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnResult(res)

		// Act
		err := s.Delete(context.Background(), 12)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
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

func TestIntegration_CreateType(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(9, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE_TYPE)).ExpectExec().WillReturnResult(res)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	id, err := s.CreateType(context.Background(), "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 9, id)
}

func TestIntegration_GetAllReports(t *testing.T) {
	// Arrange
	expected := []domain.Report{
		{Count: 15, Description: "desc", ProductID: 11},
		{Count: 13, Description: "pepe", ProductID: 12},
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description", "id"})
	row.AddRow(15, "desc", 11)
	row.AddRow(13, "pepe", 12)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL_REPORTS)).ExpectQuery().WillReturnRows(row)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	products, err := s.GetAllReports(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, products)
}

func TestIntegration_GetOneReport(t *testing.T) {
	// Arrange
	expected := []domain.Report{
		{Count: 1, Description: "desc", ProductID: 10},
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description"})
	row.AddRow(1, "desc")
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE_REPORT)).ExpectQuery().WillReturnRows(row)

	r := NewRepository(db)
	s := NewService(r)

	// Act
	pr, err := s.GetOneReport(context.Background(), 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, pr)
}
