package product_records

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

var dummyPR = domain.ProductRecord{
	LastUpdateDate: "12-12-2012",
	PurchasePrice:  123.123,
	SalePrice:      321.321,
	ProductID:      12,
}

func TestRepoValidateProductID_true(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(id)"})
	row = row.AddRow(1) // must be exactly 1, so this should return true
	mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	valid := repo.ValidateProductID(context.Background(), 12)

	// Assert
	assert.True(t, valid)
}

func TestRepoValidateProductID_false(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(id)"})
	row = row.AddRow(3) // must be exactly one because it is unique, so this should return false
	mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	valid := repo.ValidateProductID(context.Background(), 12)

	// Assert
	assert.False(t, valid)
}

func TestRepoValidateProductID_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(id)"})
	row = row.AddRow(1) // returns true if it does not fail
	mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).WillReturnError(ErrDatabase).ExpectQuery().WillReturnRows(row)
	// should fail

	repo := NewRepository(db)

	// Act
	valid := repo.ValidateProductID(context.Background(), 12)

	// Assert
	assert.False(t, valid)
}

func TestRepoValidateProductID_FailsQuery(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(id)"})
	row = row.AddRow(1) // returns true if it does not fail
	mock.ExpectPrepare(regexp.QuoteMeta(VALIDATE)).ExpectQuery().WillReturnError(ErrDatabase).WillReturnRows(row)
	// should fail

	repo := NewRepository(db)

	// Act
	valid := repo.ValidateProductID(context.Background(), 12)

	// Assert
	assert.False(t, valid)
}

func TestStore_Ok(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 12, nil
	newID, err := repo.Store(context.Background(), dummyPR)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 12, newID)
}

func TestStore_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE)).WillReturnError(ErrDatabase).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Store(context.Background(), dummyPR)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, newID)
	assert.Equal(t, ErrDatabase, err)
}

func TestStore_FailsExec(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE)).ExpectExec().WillReturnError(ErrDatabase).WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Store(context.Background(), dummyPR)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, newID)
	assert.Equal(t, ErrDatabase, err)
}

func TestStore_FailsResult(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewErrorResult(ErrDatabase)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Store(context.Background(), dummyPR)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, newID)
	assert.Equal(t, ErrDatabase, err)
}
