package product

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

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

func TestRepoGetOneReport_Ok(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description"})
	row.AddRow(1, "desc")
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE_REPORT)).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return 1, "desc", nil
	count, description, err := repo.GetOneReport(context.Background(), 12)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, "desc", description)
}

func TestRepoGetOneReport_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description"})
	row.AddRow(15, "desc")
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE_REPORT)).WillReturnError(ErrDatabase).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return 0, "", ErrDatabase
	count, description, err := repo.GetOneReport(context.Background(), 12)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, "", description)
}

func TestRepoGetOneReport_FailsQuery(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description"})
	row.AddRow(15, "desc")
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE_REPORT)).ExpectQuery().WillReturnError(ErrDatabase).WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return 0, "", ErrDatabase
	count, description, err := repo.GetOneReport(context.Background(), 12)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, "", description)
}

func TestRepoGetAllReports_Ok(t *testing.T) {
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

	repo := NewRepository(db)

	// Act
	// should return expected, nil
	reports, err := repo.GetAllReports(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, reports)
}

func TestRepoGetAllReports_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description", "id"})
	row.AddRow(15, "desc", 11)
	row.AddRow(13, "pepe", 12)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL_REPORTS)).WillReturnError(ErrDatabase).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return []domain.Report{}, ErrDatabase
	reports, err := repo.GetAllReports(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, []domain.Report{}, reports)
}

func TestRepoGetAllReports_FailsQuery(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"COUNT(pr.id)", "description", "id"})
	row.AddRow(15, "desc", 11)
	row.AddRow(13, "pepe", 12)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL_REPORTS)).ExpectQuery().WillReturnError(ErrDatabase).WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return []domain.Report{}, ErrDatabase
	reports, err := repo.GetAllReports(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, []domain.Report{}, reports)
}

func TestRepoStoreType_Ok(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(1, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE_TYPE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 1, nil
	typeId, err := repo.StoreType(context.Background(), "pepe")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, typeId)
}

func TestRepoStoreType_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(1, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE_TYPE)).WillReturnError(ErrDatabase).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	typeId, err := repo.StoreType(context.Background(), "pepe")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, typeId)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoStoreType_FailsExec(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(1, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE_TYPE)).ExpectExec().WillReturnError(ErrDatabase).WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	typeId, err := repo.StoreType(context.Background(), "pepe")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, typeId)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoStoreType_FailsResult(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewErrorResult(ErrDatabase)
	mock.ExpectPrepare(regexp.QuoteMeta(STORE_TYPE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	typeId, err := repo.StoreType(context.Background(), "pepe")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, typeId)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoGetAll_Ok(t *testing.T) {
	// Arrange
	expected := []domain.Product{
		{ID: 1, Description: "pepe", ExpirationRate: 12, FreezingRate: 13, Height: 12.1, Length: 10.9, Netweight: 1111.11, ProductCode: "UNIQUE", RecomFreezTemp: -10.2, Width: 172, ProductTypeID: 2, SellerID: 1},
		{ID: 2, Description: "not_pepe", ExpirationRate: 212, FreezingRate: 321, Height: 0.18, Length: 332.1, Netweight: 12312.1323, ProductCode: "OTHER", RecomFreezTemp: -50, Width: 172, ProductTypeID: 2, SellerID: 1},
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	rows.AddRow(1, "pepe", 12, 13, 12.1, 10.9, 1111.11, "UNIQUE", -10.2, 172, 2, 1)
	rows.AddRow(2, "not_pepe", 212, 321, 0.18, 332.1, 12312.1323, "OTHER", -50, 172, 2, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL)).ExpectQuery().WillReturnRows(rows)

	repo := NewRepository(db)

	// Act
	// should return expected, nil
	products, err := repo.GetAll(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, products)
}

func TestRepoGetAll_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	rows.AddRow(1, "pepe", 12, 13, 12.1, 10.9, 1111.11, "UNIQUE", -10.2, 172, 2, 1)
	rows.AddRow(2, "not_pepe", 212, 321, 0.18, 332.1, 12312.1323, "OTHER", -50, 172, 2, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL)).WillReturnError(ErrDatabase).ExpectQuery().WillReturnRows(rows)

	repo := NewRepository(db)

	// Act
	// should return []domain.Product{}, ErrDatabase
	products, err := repo.GetAll(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, []domain.Product{}, products)
}

func TestRepoGetAll_FailsQuery(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	rows.AddRow(1, "pepe", 12, 13, 12.1, 10.9, 1111.11, "UNIQUE", -10.2, 172, 2, 1)
	rows.AddRow(2, "not_pepe", 212, 321, 0.18, 332.1, 12312.1323, "OTHER", -50, 172, 2, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ALL)).ExpectQuery().WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return []domain.Product{}, ErrDatabase
	products, err := repo.GetAll(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, []domain.Product{}, products)
}

func TestRepoGet_Ok(t *testing.T) {
	// Arrange
	expected := domain.Product{ID: 1, Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	row.AddRow(1, "desc", 1, 1, 1.1, 1.0, 12.2, "code", -1, 10, 10, 3)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return expected, nil
	product, err := repo.Get(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, product)
}

func TestRepoGet_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	row.AddRow(1, "desc", 1, 1, 1.1, 1.0, 12.2, "code", -1, 10, 10, 3)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).WillReturnError(ErrDatabase).ExpectQuery().WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return domain.Product{}, ErrDatabase
	product, err := repo.Get(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, domain.Product{}, product)
}

func TestRepoGet_FailsQuery(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "lenght", "netweight", "product_code", "recommended_freezing_temperature", "width", "id_product_type", "id_seller"})
	row.AddRow(1, "desc", 1, 1, 1.1, 1.0, 12.2, "code", -1, 10, 10, 3)
	mock.ExpectPrepare(regexp.QuoteMeta(GET_ONE)).ExpectQuery().WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return domain.Product{}, ErrDatabase
	product, err := repo.Get(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, domain.Product{}, product)
}

func TestRepoExists_True(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"product_code"})
	row.AddRow("code")
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS)).WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return true
	exists := repo.Exists(context.Background(), "code")

	// Assert
	assert.True(t, exists)
}

func TestRepoExists_False(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	row := mock.NewRows([]string{"product_code"})
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS)).WillReturnRows(row)

	repo := NewRepository(db)

	// Act
	// should return false
	exists := repo.Exists(context.Background(), "code")

	// Assert
	assert.False(t, exists)
}

func TestRepoExists_Fails(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(EXISTS)).WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return false
	exists := repo.Exists(context.Background(), "code")

	// Assert
	assert.False(t, exists)
}

func TestRepoSave_Ok(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(SAVE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 12, nil
	newID, err := repo.Save(context.Background(), product)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 12, newID)
}

func TestRepoSave_FailsPrepare(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(SAVE)).WillReturnError(ErrDatabase).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Save(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, 0, newID)
}

func TestRepoSave_FailsExec(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE)).ExpectExec().WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Save(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, 0, newID)
}

func TestRepoSave_ErrorResult(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	resErr := sqlmock.NewErrorResult(ErrDatabase)
	mock.ExpectPrepare(regexp.QuoteMeta(SAVE)).ExpectExec().WillReturnResult(resErr)

	repo := NewRepository(db)

	// Act
	// should return 0, ErrDatabase
	newID, err := repo.Save(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, 0, newID)
}

func TestRepoUpdate_Ok(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return nil
	err = repo.Update(context.Background(), product)

	// Assert
	assert.NoError(t, err)
}

func TestRepoUpdate_FailsPrepare(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(12, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE)).WillReturnError(ErrDatabase).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Update(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoUpdate_FailsExec(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE)).ExpectExec().WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Update(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoUpdate_ErrorResult(t *testing.T) {
	// Arrange
	product := domain.Product{Description: "desc", ExpirationRate: 1, FreezingRate: 1, Height: 1.1, Length: 1.0, Netweight: 12.2, ProductCode: "code", RecomFreezTemp: -1, Width: 10, ProductTypeID: 10, SellerID: 3}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	errRes := sqlmock.NewErrorResult(ErrDatabase)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE)).ExpectExec().WillReturnResult(errRes)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Update(context.Background(), product)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoDelete_Ok(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(0, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return nil
	err = repo.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
}

func TestRepoDelete_FailsPrepare(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(0, 1)
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).WillReturnError(ErrDatabase).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Delete(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoDelete_FailsExec(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnError(ErrDatabase)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Delete(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoDelete_ErrorResult(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	errRes := sqlmock.NewErrorResult(ErrDatabase)
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnResult(errRes)

	repo := NewRepository(db)

	// Act
	// should return ErrDatabase
	err = repo.Delete(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestRepoDelete_ErrNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	res := sqlmock.NewResult(0, 0)
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE)).ExpectExec().WillReturnResult(res)

	repo := NewRepository(db)

	// Act
	// should return ErrNotFound
	err = repo.Delete(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}
