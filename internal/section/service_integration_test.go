package section

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

// ------------------------------- READ ---------------------------------

func Test_GetAll_Service_Integration(t *testing.T) {
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
	s := NewService(r)
	ctx := context.Background()

	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		sections, err := s.GetAll(ctx)

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
		sections, err := s.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
		assert.Empty(t, sections, sections)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// ... To Continue ...
