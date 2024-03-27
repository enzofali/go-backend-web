package carry

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_GetAllC(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)
	t.Run("Ok", func(t *testing.T) {

		//arrange

		expected := []domain.Carrie{
			{Id: 1, Cid: "ABC34", Company_name: "DHL", Address: "Cra 40 # 34-56", Telephone: "2245678", Locality_id: "L001"},
			{Id: 2, Cid: "ABC35", Company_name: "Servientrega", Address: "Cra 32 # 3-26", Telephone: "2235678", Locality_id: "L002"},
		}

		rows := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

		for _, f := range expected {
			rows.AddRow(f.Id, f.Cid, f.Company_name, f.Address, f.Telephone, f.Locality_id)
		}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)

		// act
		quantityCarries, err := service.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(quantityCarries))
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange

		expect := []domain.Carrie{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrBD)

		// act
		quantityCarries, err := service.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_Integration_GetAllCLocality(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)
	t.Run("Ok", func(t *testing.T) {

		//arrange

		expected := []domain.CarrieLocality{
			{Locality_id: "L001", Locality_name: "Tuluá", Cant_carries: 2},
			{Locality_id: "L002", Locality_name: "Cali", Cant_carries: 4},
		}

		rows := mock.NewRows([]string{"locality_id", "local_name", "carries_count"})

		for _, f := range expected {
			rows.AddRow(f.Locality_id, f.Locality_name, f.Cant_carries)
		}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocality)).WillReturnRows(rows)

		// act
		quantityCarries, err := service.GetByLocality(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(quantityCarries))
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange

		expect := []domain.CarrieLocality{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocality)).WillReturnError(ErrBD)

		// act
		quantityCarries, err := service.GetByLocality(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_Integration_GetCLocalityID(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	expected := domain.CarrieLocality{Locality_id: "L001", Locality_name: "Tuluá", Cant_carries: 2}
	row := mock.NewRows([]string{"locality_id", "local_name", "carries_count"})

	t.Run("Ok", func(t *testing.T) {
		// arrange

		row.AddRow(expected.Locality_id, expected.Locality_name, expected.Cant_carries)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocalityID)).WithArgs("L001").WillReturnRows(row)

		// act
		quantityCarries, err := service.GetByLocalityID(ctx, "L001")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange

		expect := domain.CarrieLocality{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocalityID)).WithArgs("L001").WillReturnError(ErrNotExist)

		// act
		quantityCarries, err := service.GetByLocalityID(ctx, "L001")

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrNotExist)
		assert.Equal(t, expect, quantityCarries)

	})

}

func Test_Integration_CreateC(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	carry := domain.Carrie{Id: 1, Cid: "ABC34", Company_name: "DHL", Address: "Cra 40 # 34-56", Telephone: "2245678", Locality_id: "L001"}

	t.Run("Ok", func(t *testing.T) {
		// arrange

		row := mock.NewRows([]string{"cid"})
		row.AddRow("ABC34")

		row2 := mock.NewRows([]string{"Id"})
		row2.AddRow("L001")

		expected := 1
		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnError(ErrBD)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExistFK)).WillReturnRows(row2)
		mock.ExpectPrepare(regexp.QuoteMeta(QueryCreate)).ExpectExec().WithArgs(carry.Cid, carry.Company_name, carry.Address, carry.Telephone, carry.Locality_id).WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		carryObt, err := service.Crear(ctx, carry)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, carryObt.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Conflict error", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"cid"})
		row.AddRow("ABC34")

		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnRows(row) //al enviar algo el error del repo = nil y genera un true en al service

		// act
		carryObt, err := service.Crear(ctx, carry)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrExist)
		assert.Equal(t, 0, carryObt.Id)
	})

	t.Run("Internal Error", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"cid"})
		row.AddRow("ABC34")

		row2 := mock.NewRows([]string{"Id"})
		row2.AddRow("L001")

		expected := 0
		mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WillReturnError(ErrBD)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExistFK)).WillReturnRows(row2)
		mock.ExpectPrepare(regexp.QuoteMeta(QueryCreate)).WillReturnError(ErrBD)
		// act
		carryObt, err := service.Crear(ctx, carry)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Equal(t, expected, carryObt.Id)

	})

}
