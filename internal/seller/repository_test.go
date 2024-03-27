package seller

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	t.Run("Ok", func(t *testing.T) {
		// arrange
		expected := []domain.Seller{
			{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"},
			{ID: 2, CID: 2, CompanyName: "Coto", Address: "Av Corrientes 980", Telephone: "48792311", Locality_id: "2"},
		}
		rows := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		for _, d := range expected {
			rows.AddRow(d.ID, d.CID, d.CompanyName, d.Address, d.Telephone, d.Locality_id)
		}

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantitySeller, err := rp.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Internal Error", func(t *testing.T) {
		// arrange
		var expect []domain.Seller

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantitySeller, err := rp.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrIntern)
		assert.Equal(t, expect, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		expected := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(expected.ID, expected.CID, expected.CompanyName, expected.Address, expected.Telephone, expected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(1).WillReturnRows(row)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantitySeller, err := rp.Get(ctx, 1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Internal error", func(t *testing.T) {
		// arrange
		expect := domain.Seller{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(1).WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantitySeller, err := rp.Get(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrIntern)
		assert.Equal(t, expect, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Error not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expected := domain.Seller{}
		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(1).WillReturnRows(row)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		quantitySeller, err := rp.Get(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
		assert.Equal(t, expected, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func Test_Exist_True(t *testing.T) {
	//arage
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := true
	row := mock.NewRows([]string{"cid"})
	row.AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(QueryExistsCid)).WithArgs(1).WillReturnRows(row)

	rp := NewRepository(db)

	ctx := context.Background()

	//act
	result := rp.Exists(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_Save(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Ok", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 1

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id).WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error Prepare", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Internal Error", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error: code mysql invalid FK", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452})

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidLocality, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error: code mysql internal error", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1054})

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.Equal(t, expected, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error: Rows affected", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}
		expected := 0

		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id).WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		lastId, err := rp.Save(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
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
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, seller)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error prepare", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error exec", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error result", func(t *testing.T) {
		// arrange
		seller := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Av Tronador 980", Telephone: "48792311", Locality_id: "4"}

		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrIntern))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Update(ctx, seller)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
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
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Error exec", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnError(ErrIntern)

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error result", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrIntern))

		rp := NewRepository(db)

		ctx := context.Background()
		// act
		err = rp.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
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
