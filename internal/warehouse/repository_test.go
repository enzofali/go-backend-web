package warehouse

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	QueryGetAll  = "SELECT * FROM warehouses"
	QueryGetByID = "SELECT * FROM warehouses WHERE id=?;"
	QueryExist   = "SELECT warehouse_code FROM warehouses WHERE warehouse_code=?;"
	QuerySave    = "INSERT INTO warehouses (address, telephone, warehouse_code, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?, ?)"
	QueryUpdate  = "UPDATE warehouses SET address=?, telephone=?, warehouse_code=?, minimum_capacity=?, minimum_temperature=? WHERE id=?"
	QueryDelete  = "DELETE FROM warehouses WHERE id=?"
)

func Test_GetAllWh(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
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

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityWarehouses, err := rp.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityWarehouses)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange

		var expect []domain.Warehouse
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrBD)

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityWarehouses, err := rp.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityWarehouses)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_GetWh(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		expected := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

		row := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})
		row.AddRow(expected.ID, expected.Address, expected.Telephone, expected.WarehouseCode, expected.MinimumCapacity, expected.MinimumTemperature)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WithArgs(1).WillReturnRows(row)
		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityWarehouse, err := rp.Get(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityWarehouse)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange

		expect := domain.Warehouse{}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WithArgs(1).WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityWarehouse, err := rp.Get(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityWarehouse)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expected := domain.Warehouse{}

		row := mock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"})

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetByID)).WithArgs(1).WillReturnRows(row)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityWarehouse, err := rp.Get(ctx, 1)

		// assert
		assert.Error(t, err)
		//assert.Equal(t, err, ErrNotFound)
		assert.Equal(t, expected, quantityWarehouse)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_Exist_True(t *testing.T) {
	//arrage
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := true
	row := mock.NewRows([]string{"warehouse_code"})
	row.AddRow("ABC123")

	mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WithArgs("ABC123").WillReturnRows(row)

	rp := NewRepository(db)

	ctx := context.Background()

	//act
	result := rp.Exists(ctx, "ABC123")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// falla
func Test_SaveWh(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 0, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}
		expected := 1

		mock.ExpectPrepare(regexp.QuoteMeta(QuerySave)).ExpectExec().WithArgs(wareH.Address, wareH.Telephone, wareH.WarehouseCode, wareH.MinimumCapacity, wareH.MinimumTemperature).WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, wareH)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Internal Error", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QuerySave)).ExpectExec().WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error Prepare", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QuerySave)).WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	t.Run("OK", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, wareH)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error prepare", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error exec", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error result", func(t *testing.T) {
		// arrange
		wareH := domain.Warehouse{ID: 1, Address: "Calle 23 #4-45", Telephone: "2245678,", WarehouseCode: "ABC123", MinimumCapacity: 10, MinimumTemperature: 22}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrBD))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, wareH)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	t.Run("OK", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error prepare", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Error exec", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error result", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrBD))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error not found", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WithArgs(1).WillReturnResult(driver.RowsAffected(0))

		ctx := context.Background()

		rp := NewRepository(db)
		err = rp.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
