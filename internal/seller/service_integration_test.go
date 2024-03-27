package seller

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	t.Run("OK", func(t *testing.T) {
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

		// act
		sellers, err := service.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(sellers))
		assert.Equal(t, expected, sellers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnError(ErrIntern)

		// act
		sellers, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
		assert.Nil(t, sellers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Intregation_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	sellerExpected := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

	t.Run("OK", func(t *testing.T) {
		//arrange
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(sellerExpected.CID).WillReturnRows(row)

		//act
		seller, err := service.GetByID(ctx, sellerExpected.ID)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, sellerExpected, seller)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		expect := domain.Seller{}
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(sellerExpected.ID).WillReturnError(ErrIntern)

		// act
		quantitySeller, err := service.GetByID(ctx, sellerExpected.ID)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrIntern)
		assert.Equal(t, expect, quantitySeller)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Not found error", func(t *testing.T) {
		//arrange
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetById)).WithArgs(sellerExpected.ID).WillReturnRows(row)

		//act
		seller, err := service.GetByID(ctx, sellerExpected.ID)

		//assert
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Empty(t, domain.Seller{}, seller)
	})
}

func Test_Intregation_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	sellerToCreate := domain.Seller{CID: 3, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}

	t.Run("OK", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WithArgs(sellerToCreate.CID, sellerToCreate.CompanyName, sellerToCreate.Address, sellerToCreate.Telephone, sellerToCreate.Locality_id).WillReturnResult(sqlmock.NewResult(3, 1))

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Conflict error", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"cid"})
		row.AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExistsCid)).WillReturnRows(row)

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrConflict)
		assert.Equal(t, 0, id)
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryInsert)).ExpectExec().WillReturnError(ErrIntern)

		// act
		id, err := service.Create(ctx, sellerToCreate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err, ErrIntern)
		assert.Equal(t, 0, id)
	})
}

func Test_Intregation_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repoMock := NewRepository(db)
	service := NewService(repoMock)

	sellerToUpdate := domain.Seller{ID: 2, CID: 3, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error conflict", func(t *testing.T) {
		// arrange
		row := mock.NewRows([]string{"cid"})
		row.AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(QueryExistsCid)).WillReturnRows(row)

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrConflict, err)
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryUpdate)).WillReturnError(ErrIntern)

		// act
		err := service.Update(ctx, sellerToUpdate)

		// assert
		assert.Error(t, err)
		assert.Equal(t, ErrIntern, err)
	})
}

func Test_integration_Delete(t *testing.T) {
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
		err := service.Delete(ctx, 2)

		// assert
		assert.NoError(t, err)
	})

	t.Run("Not found error", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(QueryDelete)).ExpectExec().WillReturnError(ErrIntern)

		// act
		err := service.Delete(ctx, 2)

		// assert
		assert.Error(t, err)
		assert.EqualError(t, ErrIntern, err.Error())
	})
}
