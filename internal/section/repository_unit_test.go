package section

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

// ------------------------------- READ ---------------------------------

func Test_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.Section{
		{ID: 1, SectionNumber: 1, CurrentTemperature: 15, MinimumTemperature: -20, CurrentCapacity: 20, MinimumCapacity: 5, MaximumCapacity: 50, WarehouseID: 1, ProductTypeID: 1},
		{ID: 2, SectionNumber: 3, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1},
	}
	rows := mock.NewRows([]string{"id", "section_number", " current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "id_product_type"})
	for _, d := range expected {
		rows.AddRow(d.ID, d.SectionNumber, d.CurrentTemperature, d.MinimumTemperature, d.CurrentCapacity, d.MinimumCapacity, d.MaximumCapacity, d.WarehouseID, d.ProductTypeID)
	}

	r := NewRepository(db)
	ctx := context.Background()

	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		sections, err := r.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		sections, err := r.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrInternal", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expected := []domain.Section{
			{ID: 1, CurrentTemperature: 15, MinimumTemperature: -20, CurrentCapacity: 20, MinimumCapacity: 5, MaximumCapacity: 50, WarehouseID: 1, ProductTypeID: 1},
			{ID: 2, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1},
		}
		rows := mock.NewRows([]string{"id", "current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "id_product_type"})
		for _, d := range expected {
			rows.AddRow(d.ID, d.CurrentTemperature, d.MinimumTemperature, d.CurrentCapacity, d.MinimumCapacity, d.MaximumCapacity, d.WarehouseID, d.ProductTypeID)
		}

		r := NewRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		sections, err := r.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := domain.Section{ID: 1, SectionNumber: 3, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1}
	row := mock.NewRows([]string{"id", "section_number", " current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "id_product_type"})
	row.AddRow(expected.ID, expected.SectionNumber, expected.CurrentTemperature, expected.MinimumTemperature, expected.CurrentCapacity, expected.MinimumCapacity, expected.MaximumCapacity, expected.WarehouseID, expected.ProductTypeID)

	r := NewRepository(db)
	ctx := context.Background()

	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections WHERE id=?;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)

		// act
		sections, err := r.GetByID(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrSectionNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrNoRows)

		// act
		sections, err := r.GetByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrSectionNotFound, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		sections, err := r.GetByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_GetAllReportProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.SectionReportProducts{
		{ID: 1, SectionNumber: 1, ProductCount: 150},
		{ID: 2, SectionNumber: 3, ProductCount: 250},
	}
	rows := mock.NewRows([]string{"id", "section_number", " product_count"})
	for _, d := range expected {
		rows.AddRow(d.ID, d.SectionNumber, d.ProductCount)
	}

	r := NewRepository(db)
	ctx := context.Background()

	query := "SELECT s.id, s.section_number, COALESCE(sum(pb.current_quantity),0) FROM sections as s LEFT JOIN products_batches as pb ON s.id = pb.section_id GROUP BY s.id, s.section_number;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		sections, err := r.GetAllReportProducts(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		sections, err := r.GetAllReportProducts(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrInternal", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expected := []domain.SectionReportProducts{
			{ID: 1, SectionNumber: 1},
			{ID: 2, SectionNumber: 3},
		}
		rows := mock.NewRows([]string{"id", "section_number"})
		for _, d := range expected {
			rows.AddRow(d.ID, d.SectionNumber)
		}

		r := NewRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		sections, err := r.GetAllReportProducts(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_GetReportProductsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.SectionReportProducts{{ID: 1, SectionNumber: 1, ProductCount: 150}}
	rows := mock.NewRows([]string{"id", "section_number", " product_count"})
	for _, d := range expected {
		rows.AddRow(d.ID, d.SectionNumber, d.ProductCount)
	}

	r := NewRepository(db)
	ctx := context.Background()

	query := "SELECT s.id, s.section_number, COALESCE(sum(pb.current_quantity),0) FROM sections as s LEFT JOIN products_batches as pb ON s.id = pb.section_id WHERE s.id = ? GROUP BY s.id, s.section_number;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)
		// act
		sections, err := r.GetReportProductsByID(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrSectionNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrNoRows)

		// act
		sections, err := r.GetReportProductsByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Empty(t, sections, sections)
		assert.Equal(t, ErrSectionNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		sections, err := r.GetReportProductsByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Empty(t, sections, sections)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// -------------------------------- WRITE --------------------------------

func Test_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	data := domain.Section{
		ID:                 0,
		SectionNumber:      1234,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	r := NewRepository(db)
	ctx := context.Background()

	query := "INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 1, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrWareHouseNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`warehouses`"})

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrWareHouseNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrProductTypeNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`product_types`"})

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrProductTypeNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrExistsSectionNumber", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrExistsSectionNumber)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RowsAffected: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 0))

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LastInsertId: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		// act
		lastId, err := r.Create(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	data := domain.Section{
		ID:                 1,
		SectionNumber:      1234,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	r := NewRepository(db)
	ctx := context.Background()

	query := "UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, id_product_type=? WHERE id=?;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		err = r.Update(ctx, data)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrWareHouseNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`warehouses`"})

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrWareHouseNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrProductTypeNotFound", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`product_types`"})

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrProductTypeNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrExistsSectionNumber", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrExistsSectionNumber)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RowsAffected: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

		// act
		err = r.Update(ctx, data)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id := 1

	r := NewRepository(db)
	ctx := context.Background()

	query := "DELETE FROM sections WHERE id=?;"
	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		err = r.Delete(ctx, id)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		err = r.Delete(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(id).WillReturnError(ErrInternal)

		// act
		err = r.Delete(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RowsAffected: ErrInternal", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 0))

		// act
		err = r.Delete(ctx, id)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
