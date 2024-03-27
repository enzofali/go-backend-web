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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/section"
	"github.com/stretchr/testify/assert"
)

func createServerSection(db *sql.DB) *gin.Engine {
	repo := section.NewRepository(db)
	service := section.NewService(repo)
	handler := NewSection(service)

	eng := gin.Default()

	sections := eng.Group("/api/v1/sections")
	{
		sections.GET("/", handler.GetAll())
		sections.GET("/:id", handler.Get())
		sections.GET("/reportProducts", handler.GetReportProducts())
		sections.POST("/", handler.Create())
		sections.PATCH("/:id", handler.Update())
		sections.DELETE("/:id", handler.Delete())
	}
	return eng
}

func createRequestSection(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func Test_GetAll_Section(t *testing.T) {
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

	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Get_Section(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := domain.Section{ID: 1, SectionNumber: 3, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1}
	row := mock.NewRows([]string{"id", "section_number", " current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "id_product_type"})
	row.AddRow(expected.ID, expected.SectionNumber, expected.CurrentTemperature, expected.MinimumTemperature, expected.CurrentCapacity, expected.MinimumCapacity, expected.MaximumCapacity, expected.WarehouseID, expected.ProductTypeID)

	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections WHERE id=?;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrNoRows)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_GetReportProducts_Section(t *testing.T) {
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

	query := "SELECT s.id, s.section_number, COALESCE(sum(pb.current_quantity),0) FROM sections as s LEFT JOIN products_batches as pb ON s.id = pb.section_id GROUP BY s.id, s.section_number;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/reportProducts", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/reportProducts?id=error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
	/*
		t.Run("Validate Section not found", func(t *testing.T) {
			// arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			server := createServerSection(db)

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WillReturnError(sql.ErrNoRows)

			// act
			request, response := createRequestSection(http.MethodGet, "/api/v1/sections/reportProducts?id=1", "")
			server.ServeHTTP(response, request)

			// assert
			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
		})
	*/
	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(ErrInternal)

		// act
		request, response := createRequestSection(http.MethodGet, "/api/v1/sections/reportProducts?id=1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Create_Section(t *testing.T) {

	query := "INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
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
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number": error, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
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
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate unique section_number", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate WareHouse not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`warehouses`"})

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Product Type not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`product_types`"})

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
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
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		// act
		request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Update_Section(t *testing.T) {

	query := "UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, id_product_type=? WHERE id=?;"

	/*
		t.Run("Ok", func(t *testing.T) {
			// arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			server := createServerSection(db)

			mock.ExpectPrepare(regexp.QuoteMeta(query)).
				ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

			// act
			request, response := createRequestSection(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
			server.ServeHTTP(response, request)

			// assert
			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
		})
	*/

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodPatch, "/api/v1/sections/error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(sql.ErrNoRows)

		// act
		request, response := createRequestSection(http.MethodPatch, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	/*
		t.Run("Validate json.NewDecoder(ctx.Request.Body).Decode(&section)", func(t *testing.T) {
			// arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			server := createServerSection(db)

			mock.ExpectPrepare(regexp.QuoteMeta(query)).
				ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

			// act
			request, response := createRequestSection(http.MethodPatch, "/api/v1/sections/1", `{"section_number": error, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
			server.ServeHTTP(response, request)

			// assert
			assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
			assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
		})

		// New id should not be specified in request body


			t.Run("Validate required JSON fields", func(t *testing.T) {
				// arrange
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()
				server := createServerSection(db)

				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

				// act
				request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
				server.ServeHTTP(response, request)

				// assert
				assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
				assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
			})

			t.Run("Validate unique section_number", func(t *testing.T) {
				// arrange
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()
				server := createServerSection(db)

				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

				// act
				request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
				server.ServeHTTP(response, request)

				// assert
				assert.Equal(t, http.StatusConflict, response.Code)
				assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
			})

			t.Run("", func(t *testing.T) {
				// arrange
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()
				server := createServerSection(db)

				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`warehouses`"})

				// act
				request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
				server.ServeHTTP(response, request)

				// assert
				assert.Equal(t, http.StatusConflict, response.Code)
				assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
			})

			t.Run("Validate Product Type not found", func(t *testing.T) {
				// arrange
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()
				server := createServerSection(db)

				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452, Message: "`product_types`"})

				// act
				request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
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
				server := createServerSection(db)

				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WillReturnError(&mysql.MySQLError{})

				// act
				request, response := createRequestSection(http.MethodPost, "/api/v1/sections/", `{"section_number":1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
				server.ServeHTTP(response, request)

				// assert
				assert.Equal(t, http.StatusInternalServerError, response.Code)
				assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
			})
	*/
}

func Test_Delete_Section(t *testing.T) {

	query := "DELETE FROM sections WHERE id=?;"

	t.Run("Ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodDelete, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		request, response := createRequestSection(http.MethodDelete, "/api/v1/sections/error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		server := createServerSection(db)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrNoRows)

		// act
		request, response := createRequestSection(http.MethodDelete, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

}
