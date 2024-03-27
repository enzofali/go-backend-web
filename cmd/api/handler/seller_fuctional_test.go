package handler

import (
	"database/sql"
	"database/sql/driver"
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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/seller"
	"github.com/stretchr/testify/assert"
)

// ______________________________________________________
// tools
func Create_Server_Seller_Fuctional(db *sql.DB) *gin.Engine {
	// instances
	repoSqlMock := seller.NewRepository(db)
	service := seller.NewService(repoSqlMock)
	handler := NewSeller(service)

	// server
	server := gin.Default()

	// -> routes
	routes := server.Group("/api/v1/sellers")
	{
		//sr := routes.Group("/sellers")
		routes.GET("/", handler.GetAll())
		routes.POST("/", handler.Create())
		routes.GET("/:id", handler.Get())
		routes.PATCH("/:id", handler.Update())
		routes.DELETE("/:id", handler.Delete())
	}

	return server
}

func NewRequestSellerFunctional(method, path, body string) (req *http.Request, res *httptest.ResponseRecorder) {
	// request
	req = httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// response
	res = httptest.NewRecorder()

	return
}

func Test_Functional_Create_Seller(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	type responseStruct struct {
		Data domain.Seller `json:"data"`
	}
	sellerToCreate := domain.Seller{CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	sellerCreated := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	data := responseStruct{
		Data: sellerCreated,
	}

	t.Run("OK", func(t *testing.T) {
		// arrange
		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryInsert)).ExpectExec().WithArgs(sellerToCreate.CID, sellerToCreate.CompanyName, sellerToCreate.Address, sellerToCreate.Telephone, sellerToCreate.Locality_id).WillReturnResult(sqlmock.NewResult(1, 1))

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, data, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Bad request error", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "error bad request",
		}

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Validator error", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Seller.CID' Error:Field validation for 'CID' failed on the 'required' tag",
		}

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Cid invalid error", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "invalid cid",
		}

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"cid": -1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("conflict error creating", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: seller.ErrConflict.Error(),
		}
		row := mock.NewRows([]string{"cid"})
		row.AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryExistsCid)).WillReturnRows(row)

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error creating", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryInsert)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1054})

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("not found error creating", func(t *testing.T) {
		// arrange

		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrInvalidLocality.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryInsert)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452})

		request, response := NewRequestSellerFunctional(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Functional_GetAll_Seller(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	type responseStruct struct {
		Data []domain.Seller `json:"data"`
	}
	sellersExpected := []domain.Seller{
		{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"},
		{ID: 2, CID: 2, CompanyName: "Digital House", Address: "Monroe 860", Telephone: "47470000", Locality_id: "6700"},
	}
	data := responseStruct{
		Data: sellersExpected,
	}
	t.Run("Ok", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)
		rows := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		for _, d := range sellersExpected {
			rows.AddRow(d.ID, d.CID, d.CompanyName, d.Address, d.Telephone, d.Locality_id)
		}

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetAll)).WillReturnRows(rows)

		request, response := NewRequestSellerFunctional(http.MethodGet, "/api/v1/sellers/", "")

		// act
		server.ServeHTTP(response, request)
		var sellers responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellers)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 200, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, data, sellers)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error internal", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetAll)).WillReturnError(seller.ErrIntern)

		request, response := NewRequestSellerFunctional(http.MethodGet, "/api/v1/sellers/", "")

		// act
		server.ServeHTTP(response, request)
		var result errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, errResp, result)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Functional_Get_Seller(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	type responseStruct struct {
		Data domain.Seller `json:"data"`
	}
	sellerExpected := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	data := responseStruct{
		Data: sellerExpected,
	}
	t.Run("OK", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodGet, "/api/v1/sellers/1", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, data, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrNotFound.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodGet, "/api/v1/sellers/1", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error bad request", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "strconv.Atoi: parsing \"fhdfh\": invalid syntax",
		}

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodGet, "/api/v1/sellers/fhdfh", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Functional_Update_Seller(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	type responseStruct struct {
		Data domain.Seller `json:"data"`
	}
	sellerExpected := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	data := responseStruct{
		Data: sellerExpected,
	}
	t.Run("OK", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, data, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrNotFound.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error invalid id in the request (atoi)", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "strconv.Atoi: parsing \"hola\": invalid syntax",
		}

		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/hola", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error newDecoder", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "bad request",
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("the id does not match the id of the db", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "cannot update product id",
		}

		server := Create_Server_Seller_Fuctional(db)

		sDb := domain.Seller{ID: 2, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sDb.ID, sDb.CID, sDb.CompanyName, sDb.Address, sDb.Telephone, sDb.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error: validator", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Seller.Locality_id' Error:Field validation for 'Locality_id' failed on the 'required' tag",
		}
		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": ""}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Error: cid invalid", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "invalid cid",
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"cid": -1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Error conflic", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: seller.ErrConflict.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		rowexist := mock.NewRows([]string{"cid"})
		rowexist.AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryExistsCid)).WithArgs(1).WillReturnRows(rowexist)

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error updating", func(t *testing.T) {
		// arrange
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		server := Create_Server_Seller_Fuctional(db)

		row := mock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		row.AddRow(sellerExpected.ID, sellerExpected.CID, sellerExpected.CompanyName, sellerExpected.Address, sellerExpected.Telephone, sellerExpected.Locality_id)

		mock.ExpectQuery(regexp.QuoteMeta(seller.QueryGetById)).WithArgs(1).WillReturnRows(row)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryUpdate)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(seller.ErrIntern))

		request, response := NewRequestSellerFunctional(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Functional_Delete_Seller(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	t.Run("OK", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		request, response := NewRequestSellerFunctional(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert

		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("bad request err (atoi id)", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		request, response := NewRequestSellerFunctional(http.MethodDelete, "/api/v1/sellers/hola", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Not found seller error", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryDelete)).ExpectExec().WithArgs(1).WillReturnResult(driver.RowsAffected(0))

		request, response := NewRequestSellerFunctional(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		server := Create_Server_Seller_Fuctional(db)

		mock.ExpectPrepare(regexp.QuoteMeta(seller.QueryDelete)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(seller.ErrIntern))

		request, response := NewRequestSellerFunctional(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}
