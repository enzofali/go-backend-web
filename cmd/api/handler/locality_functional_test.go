package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/locality"
	"github.com/stretchr/testify/assert"
)

func CreateServerLocalityFunctional(db *sql.DB) *gin.Engine {
	// instances
	repoSqlMock := locality.NewRepository(db)
	service := locality.NewService(repoSqlMock)
	handler := NewLocality(service)

	// server
	server := gin.Default()

	// -> routes
	routes := server.Group("/api/v1/localities")
	{
		routes.POST("", handler.Create())
		routes.GET("/reportSellers", handler.GetQuantitySellerByLocality())
	}

	return server
}

func NewRequestLocalityFunctional(method, path, body string) (req *http.Request, res *httptest.ResponseRecorder) {
	// request
	req = httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// response
	res = httptest.NewRecorder()

	return
}

func Test_Functional_Create_Locality(t *testing.T) {
	// arrange
	type responseStruct struct {
		Data domain.Locality `json:"data"`
	}
	localityToCreate := domain.Locality{Id: "6701", Locality_name: "Villa Crespo", Province_name: "Buenos Aires", Country_name: "Argentina"}
	localityCreated := domain.Locality{Id: "6701", Locality_name: "Villa Crespo", Province_name: "Buenos Aires", Country_name: "Argentina"}
	data := responseStruct{
		Data: localityCreated,
	}

	t.Run("OK", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(locality.QueryInsert)).ExpectExec().WithArgs(localityToCreate.Id, localityToCreate.Locality_name, localityToCreate.Province_name, localityToCreate.Country_name).WillReturnResult(sqlmock.NewResult(6701, 1))

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult responseStruct
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, data, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Bad request error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "error bad request",
		}

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Validator error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Locality.Id' Error:Field validation for 'Id' failed on the 'required' tag",
		}

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `{"locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error creating", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(locality.QueryInsert)).ExpectExec().WithArgs(localityToCreate.Id, localityToCreate.Locality_name, localityToCreate.Province_name, localityToCreate.Country_name).WillReturnError(locality.ErrIntern)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("id duplicate error creating", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(locality.QueryInsert)).ExpectExec().WithArgs(localityToCreate.Id, localityToCreate.Locality_name, localityToCreate.Province_name, localityToCreate.Country_name).WillReturnError(&mysql.MySQLError{Number: 1062})

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: locality.ErrDuplicated.Error(),
		}

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	/* REVISAR */
	t.Run("default error creating", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(locality.QueryInsert)).WillReturnError(sql.ErrConnDone)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: "internal error",
		}

		request, response := NewRequestLocalityFunctional(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

}

func Test_Functional_GetQuantitySellerByLocality(t *testing.T) {
	// arrange
	type responseStruct struct {
		Data []domain.QuantitySellerByLocality `json:"data"`
	}

	t.Run("OK report the number of sellers of all locations", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)
		expected := []domain.QuantitySellerByLocality{
			{Locality_id: "6701", Locality_name: "Villa Crespo", Sellers_count: 6},
			{Locality_id: "6702", Locality_name: "Almagro", Sellers_count: 4},
			{Locality_id: "6703", Locality_name: "Nuñez", Sellers_count: 5},
		}
		rows := mock.NewRows([]string{"locality_id", "locality_name", " sellers_count"})
		for _, d := range expected {
			rows.AddRow(d.Locality_id, d.Locality_name, d.Sellers_count)
		}
		mock.ExpectPrepare(regexp.QuoteMeta(locality.QuerySellerAll)).ExpectQuery().WillReturnRows(rows)

		data := responseStruct{Data: expected}

		request, response := NewRequestLocalityFunctional(http.MethodGet, "/api/v1/localities/reportSellers", "")

		// act
		server.ServeHTTP(response, request)
		var report responseStruct
		err = json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, data, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error getting number of sellers from all locations", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(locality.QuerySellerAll)).ExpectQuery().WillReturnError(locality.ErrIntern)

		request, response := NewRequestLocalityFunctional(http.MethodGet, "/api/v1/localities/reportSellers", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("OK report of the number of sellers by location", func(t *testing.T) {
		type responseStruct struct {
			Data domain.QuantitySellerByLocality `json:"data"`
		}
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		expected := domain.QuantitySellerByLocality{Locality_id: "6701", Locality_name: "Villa Crespo", Sellers_count: 6}
		row := mock.NewRows([]string{"locality_id", "locality_name", " sellers_count"})
		row.AddRow(expected.Locality_id, expected.Locality_name, expected.Sellers_count)

		mock.ExpectQuery(regexp.QuoteMeta(locality.QuerySellerByLocality)).WithArgs("6701").WillReturnRows(row)

		data := responseStruct{Data: expected}

		request, response := NewRequestLocalityFunctional(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report responseStruct
		err = json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, data, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error when obtaining the number of sellers in a location", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		mock.ExpectQuery(regexp.QuoteMeta(locality.QuerySellerByLocality)).WithArgs("6701").WillReturnError(locality.ErrIntern)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}

		request, response := NewRequestLocalityFunctional(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("not found error when obtaining the number of sellers in a location", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := CreateServerLocalityFunctional(db)

		row := mock.NewRows([]string{"locality_id", "locality_name", " sellers_count"})

		mock.ExpectQuery(regexp.QuoteMeta(locality.QuerySellerByLocality)).WithArgs("6701").WillReturnRows(row)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: locality.ErrLocalityNotFound.Error(),
		}

		request, response := NewRequestLocalityFunctional(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err = json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}
