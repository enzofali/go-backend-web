package carry

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	QueryGetAll      = "SELECT id, cid, company_name, address, telephone, locality_id FROM carries;"
	QueryGetLocality = "SELECT  localities.id , localities.local_name, COUNT(carries.cid) as carries_count FROM carries " +
		"INNER JOIN localities ON carries.locality_id = localities.id GROUP BY carries.locality_id"
	QueryGetLocalityID = "SELECT  localities.id , localities.local_name, COUNT(carries.cid) as carries_count FROM carries " +
		"INNER JOIN localities ON carries.locality_id = localities.id GROUP BY carries.locality_id HAVING carries.locality_id =?;"
	QueryExist   = "SELECT cid FROM carries WHERE cid=?;"
	QueryExistFK = "SELECT id FROM localities WHERE id=?;"
	QueryCreate  = "INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
)

func Test_GetAllC(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
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

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityCarries, err := rp.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange

		var expect []domain.Carrie
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrBD)

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityCarries, err := rp.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_GetAllCLocality(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
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

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityCarries, err := rp.GetByLocality(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal Error", func(t *testing.T) {

		//arrange

		var expect []domain.CarrieLocality
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocality)).WillReturnError(ErrBD)

		rp := NewRepository(db)
		ctx := context.Background()

		// act
		quantityCarries, err := rp.GetByLocality(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_GetCLocalityID(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		expected := domain.CarrieLocality{Locality_id: "L001", Locality_name: "Tuluá", Cant_carries: 2}

		row := mock.NewRows([]string{"locality_id", "local_name", "carries_count"})
		row.AddRow(expected.Locality_id, expected.Locality_name, expected.Cant_carries)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocalityID)).WithArgs("L001").WillReturnRows(row)
		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityCarries, err := rp.GetByLocalityID(ctx, "L001")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange

		expect := domain.CarrieLocality{}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocalityID)).WithArgs("L001").WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityCarries, err := rp.GetByLocalityID(ctx, "L001")

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrBD)
		assert.Equal(t, expect, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expected := domain.CarrieLocality{}

		row := mock.NewRows([]string{"locality_id", "local_name", "carries_count"})

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetLocalityID)).WithArgs("L001").WillReturnRows(row)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantityCarries, err := rp.GetByLocalityID(ctx, "L001")

		// assert
		assert.Error(t, err)
		//assert.Equal(t, err, ErrNotFound) ??????
		assert.Equal(t, expected, quantityCarries)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_Exist_True(t *testing.T) {
	//arrage
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := true
	row := mock.NewRows([]string{"cid"})
	row.AddRow("ABC34")

	mock.ExpectQuery(regexp.QuoteMeta(QueryExist)).WithArgs("ABC34").WillReturnRows(row)

	rp := NewRepository(db)

	ctx := context.Background()

	//act
	result := rp.Exists(ctx, "ABC34")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_ExistFK_True(t *testing.T) {
	//arrage
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := true
	row := mock.NewRows([]string{"id"})
	row.AddRow("L001")

	mock.ExpectQuery(regexp.QuoteMeta(QueryExistFK)).WithArgs("L001").WillReturnRows(row)

	rp := NewRepository(db)

	ctx := context.Background()

	//act
	result := rp.ExistsFK(ctx, "L001")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// falla
func Test_CreateC(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		carry := domain.Carrie{Id: 1, Cid: "ABC34", Company_name: "DHL", Address: "Cra 40 # 34-56", Telephone: "2245678", Locality_id: "L001"}
		expected := 1

		mock.ExpectPrepare(regexp.QuoteMeta(QueryCreate)).ExpectExec().WithArgs(carry.Cid, carry.Company_name, carry.Address, carry.Telephone, carry.Locality_id).WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Crear(ctx, carry)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Internal Error", func(t *testing.T) {
		// arrange
		carry := domain.Carrie{Id: 1, Cid: "ABC34", Company_name: "DHL", Address: "Cra 40 # 34-56", Telephone: "2245678", Locality_id: "L001"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryCreate)).ExpectExec().WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Crear(ctx, carry)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error Prepare", func(t *testing.T) {
		// arrange
		carry := domain.Carrie{Id: 1, Cid: "ABC34", Company_name: "DHL", Address: "Cra 40 # 34-56", Telephone: "2245678", Locality_id: "L001"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryCreate)).WillReturnError(ErrBD)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Crear(ctx, carry)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrBD, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
