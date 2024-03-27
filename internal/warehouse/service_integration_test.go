package warehouse

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_IntegrationGetAllWh(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	t.Run("Ok", func(t *testing.T) {

		//arrange

		expected := []domain.Warehouse{
			{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22},
			{ID: 2, Address: "Calle 12 #32-21", Telephone: "2254586,", WarehouseCode: "ABC456", MinimumCapacity: 10, MinimumTemperature: 22},
		}

		rows := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})

		for _, f := range expected {
			rows.AddRow(f.ID, f.Address, f.Telephone, f.WarehouseCode, f.MinimumCapacity, f.MinimumTemperature)
		}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)

		// act
		quantityWarehouses, err := service.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(quantityWarehouses))
		assert.Equal(t, expected, quantityWarehouses)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrBD)

		// act
		quantityWarehouses, err := service.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, []domain.Warehouse{}, quantityWarehouses)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_Integration_GetWh(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	expected := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}
	row := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})

	t.Run("Ok", func(t *testing.T) {
		// arrange

		row.AddRow(expected.ID, expected.Address, expected.Telephone, expected.WarehouseCode, expected.MinimumCapacity, expected.MinimumTemperature)
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WithArgs(1).WillReturnRows(row)

		// act
		quantityWarehouse, err := service.Get(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityWarehouse)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Error_Not_Found", func(t *testing.T) {
		// arrange

		expect := domain.Warehouse{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WithArgs(1).WillReturnError(ErrNotFound)

		// act
		quantityWarehouse, err := service.Get(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
		assert.Equal(t, expect, quantityWarehouse)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_Integration_CreateWh(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

	t.Run("Ok", func(t *testing.T) {
		// arrange

		//expected := 1

		mock.ExpectPrepare(regexp.QuoteMeta(QuerySave)).ExpectExec().WithArgs(wareH.Address, wareH.Telephone, wareH.WarehouseCode, wareH.MinimumCapacity, wareH.MinimumTemperature).WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		wareHObt, err := service.Create(ctx, wareH)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, wareHObt, wareH)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Conflict error", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"warehouse_code"})
		row.AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnRows(row)

		// act
		wareHObt, err := service.Create(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrExist)
		assert.Equal(t, domain.Warehouse{}, wareHObt)
	})

	t.Run("Internal Error", func(t *testing.T) {
		// arrange

		expected := domain.Warehouse{}

		mock.ExpectPrepare(regexp.QuoteMeta(QuerySave)).ExpectExec().WillReturnError(ErrBD)

		// act
		lastId, err := service.Create(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_Integration_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

	t.Run("OK", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"warehouse_code"})
		row.AddRow(1)

		row2 := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})
		row2.AddRow(1, "Calle 23 #4-45", "2245678,", "ABC123", 10, 22)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WillReturnRows(row2)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnRows(row)

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		_, err = service.Update(ctx, wareH)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error_conflict", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"warehouse_code"})
		row.AddRow(1)

		row2 := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})
		row2.AddRow(1, "Calle 23 #4-45", "2245678,", "ABC1234", 10, 22)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WillReturnRows(row2)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnRows(row)

		// act
		_, err = service.Update(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrExist, err)

	})

	t.Run("Internal_Error", func(t *testing.T) {
		// arrange
		//row := mock.NewRows([]string{"warehouse_code"})
		//row.AddRow(1)

		row2 := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})
		row2.AddRow(1, "Calle 23 #4-45", "2245678,", "ABC1234", 10, 22)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WillReturnRows(row2)
		//mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnRows(row)
		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).WillReturnError(ErrBD)

		// act
		_, err = service.Update(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)

	})
}

func Test_Integration_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	t.Run("OK", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
		// act
		err = service.Delete(ctx, 1)

		// assert
		assert.NoError(t, err)

	})

	t.Run("Error not found", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnError(ErrNotFound)

		ctx := context.Background()

		rp := NewRepository(db)
		err = rp.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
