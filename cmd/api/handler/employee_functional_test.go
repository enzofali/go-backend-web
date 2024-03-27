package handler

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/employee"
	"github.com/stretchr/testify/assert"
)

func createServerEmployeeFunctional(db *sql.DB) *gin.Engine {
	repo := employee.NewRepository(db)
	service := employee.NewService(repo)
	handler := NewEmployee(service)

	eng := gin.Default()

	rEmp := eng.Group("/api/v1/employees")
	{
		rEmp.GET("", handler.GetAll())
		rEmp.GET("/:id", handler.Get())
		rEmp.POST("", handler.Create())
		rEmp.PATCH("/:id", handler.Update())
		rEmp.DELETE("/:id", handler.Delete())
		rEmp.GET("/reportInboundOrders", handler.GetAllWithInboundOrders())
	}

	return eng
}

func Test_Functional_Employee_GetAll(t *testing.T) {
	query := "SELECT * FROM employees"

	type response struct {
		Data []domain.Employee `json:"data"`
	}

	employees := []domain.Employee{
		{
			ID:           1,
			CardNumberID: "A12",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		},
		{
			ID:           2,
			CardNumberID: "A13",
			FirstName:    "Jose",
			LastName:     "Gomez",
			WarehouseID:  3,
		},
	}

	data := response{
		Data: employees,
	}

	t.Run("GetAll OK 200", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		for _, d := range employees {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName, d.WarehouseID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees", "")
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAll Error 500", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrConnDone)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: employee.ErrDatabase.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_Get(t *testing.T) {
	query := "SELECT * FROM employees WHERE id=?;"

	type response struct {
		Data domain.Employee `json:"data"`
	}

	employee := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	data := response{
		Data: employee,
	}

	t.Run("Get Ok 200", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employee.ID, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)
		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get Error 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get Error 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(sql.ErrConnDone)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: ErrNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_Create(t *testing.T) {
	query := "INSERT INTO employees(card_number_id,first_name,last_name,warehouse_id) VALUES (?,?,?,?)"

	queryExists := "SELECT card_number_id FROM employees WHERE card_number_id=?;"

	type response struct {
		Data domain.Employee `json:"data"`
	}

	employeeResponse := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	data := response{
		Data: employeeResponse,
	}

	t.Run("Create 200 OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"card_number_id"})
		mock.ExpectQuery(regexp.QuoteMeta(queryExists)).WithArgs("A12").WillReturnRows(rows)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error ShouldBind 422", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"first_name": "Juan", "last_name": "Perez", "warehouse_i`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: ErrBadRequest.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Validator 422", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "CardNumberID-required,",
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Not Warehouse 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "warehouse"})

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrNotWareHouse.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Exists CardID 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"card_number_id"})
		rows.AddRow(employeeResponse.CardNumberID)

		mock.ExpectQuery(regexp.QuoteMeta(queryExists)).WithArgs("A12").WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrBadRequest.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Internal 500", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1000, Message: ""})

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: ErrInternalServer.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_Update(t *testing.T) {
	query := "UPDATE employees SET first_name=?, last_name=?, warehouse_id=?  WHERE id=?"
	queryGet := "SELECT * FROM employees WHERE id=?;"

	type response struct {
		Data domain.Employee `json:"data"`
	}

	employeeBefore := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	employeeAfter := domain.Employee{
		ID:           1,
		CardNumberID: "A12",
		FirstName:    "Juan Jose",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	data := response{
		Data: employeeAfter,
	}

	t.Run("Update Success 200", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employeeBefore.ID, employeeBefore.CardNumberID, employeeBefore.FirstName, employeeBefore.LastName, employeeBefore.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnRows(rows)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Id 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/abc", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Employee Not Found 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnError(sql.ErrConnDone)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: employee.ErrNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Json 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employeeBefore.ID, employeeBefore.CardNumberID, employeeBefore.FirstName, employeeBefore.LastName, employeeBefore.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrBadRequest.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Validator 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employeeBefore.ID, employeeBefore.CardNumberID, employeeBefore.FirstName, employeeBefore.LastName, employeeBefore.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": ""}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "FirstName-required,",
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Warehouse Not Found 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employeeBefore.ID, employeeBefore.CardNumberID, employeeBefore.FirstName, employeeBefore.LastName, employeeBefore.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnRows(rows)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452})

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrNotWareHouse.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error DB 500", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(employeeBefore.ID, employeeBefore.CardNumberID, employeeBefore.FirstName, employeeBefore.LastName, employeeBefore.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(queryGet)).WithArgs(1).WillReturnRows(rows)

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: ErrInternalServer.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_Delete(t *testing.T) {
	query := "DELETE FROM employees WHERE id=?"

	t.Run("Delete OK 204", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().
			WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Id 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Not Found 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(1).WillReturnResult(driver.RowsAffected(0))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: employee.ErrNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Database 404", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(employee.ErrDatabase))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: employee.ErrDatabase.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_GetAllWithInboundOrders_WithoutID(t *testing.T) {
	query := "SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(i.id) FROM employees e LEFT JOIN inbound_orders i ON e.id = i.employee_id GROUP BY e.id;"

	type response struct {
		Data []domain.EmployeeWithInboundOrders `json:"data"`
	}

	employees := []domain.EmployeeWithInboundOrders{
		{
			ID:                 1,
			CardNumberID:       "A12",
			FirstName:          "Juan",
			LastName:           "Perez",
			WarehouseID:        1,
			InboundOrdersCount: 3,
		},
		{
			ID:                 2,
			CardNumberID:       "A13",
			FirstName:          "Jose",
			LastName:           "Gomez",
			WarehouseID:        3,
			InboundOrdersCount: 1,
		},
	}

	data := response{
		Data: employees,
	}

	t.Run("GetAllWithInboundOrders Without id OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "COUNT(i.id)"})
		for _, d := range employees {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName, d.WarehouseID, d.InboundOrdersCount)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllWithInboundOrders Without id Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("Error data base"))

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: ErrInternalServer.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Functional_Employee_GetAllWithInboundOrders_WithID(t *testing.T) {
	query := "SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(i.id) FROM employees e LEFT JOIN inbound_orders i ON e.id = i.employee_id WHERE e.id=? GROUP BY e.id;"

	type response struct {
		Data domain.EmployeeWithInboundOrders `json:"data"`
	}

	employe := domain.EmployeeWithInboundOrders{
		ID:                 1,
		CardNumberID:       "A12",
		FirstName:          "Juan",
		LastName:           "Perez",
		WarehouseID:        1,
		InboundOrdersCount: 3,
	}

	data := response{
		Data: employe,
	}

	t.Run("GetAllWithInboundOrders With id OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "COUNT(i.id)"})
		rows.AddRow(employe.ID, employe.CardNumberID, employe.FirstName, employe.LastName, employe.WarehouseID, employe.InboundOrdersCount)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=1", "")
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllWithInboundOrders With id Error id format", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=abd", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllWithInboundOrders With id Error Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(sql.ErrConnDone)

		server := createServerEmployeeFunctional(db)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: ErrNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
