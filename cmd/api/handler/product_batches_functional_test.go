package handler

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_batches"
	"github.com/stretchr/testify/assert"
)

func createServerProductBatches(db *sql.DB) *gin.Engine {
	repo := product_batches.NewRepository(db)
	service := product_batches.NewService(repo)
	handler := NewProductBatches(service)

	eng := gin.Default()

	productBatches := eng.Group("/api/v1/productBatches")
	{
		productBatches.POST("/", handler.Create())
	}
	return eng
}

func createRequestProductBatches(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func Test_Create_Product_Batches_Functional(t *testing.T) {

	query := "INSERT INTO products_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 10, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ShouldBindJSON", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": error, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate required JSON fields", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: DueDate", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-92-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: ManufacturingDate", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-91-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: ManufacturingHour", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "93:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate unique batch_number", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Product not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`products`"})

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 4, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 999, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`sections`"})

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 5, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 999}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerProductBatches(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		// act
		request, response := createRequestProductBatches(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 20, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}
