package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	inboundorder "github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/inbound_order"
	"github.com/stretchr/testify/assert"
)

func createServerInboundOrderFunctional(db *sql.DB) *gin.Engine {
	repo := inboundorder.NewRepository(db)
	service := inboundorder.NewService(repo)
	handler := NewInoudOrder(service)

	eng := gin.Default()

	rIO := eng.Group("/api/v1/inboundOrders")
	{
		rIO.POST("", handler.Create())
	}

	return eng
}

func Test_Functional_InboundOrder_Create(t *testing.T) {
	query := `INSERT INTO inbound_orders(order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?,?,?,?,?)`

	type response struct {
		Data domain.InboundOrder `json:"data"`
	}

	inboundOrderResp := domain.InboundOrder{
		ID:             1,
		OrderDate:      "2006-01-02",
		OrderNumber:    "order1",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}

	data := response{
		Data: inboundOrderResp,
	}

	t.Run("Create OK 200", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		var result response
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, data, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Json 422", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id`)
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

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "OrderNumber-required,",
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Date 400", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "otra cosa",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "Date format incorrect",
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Employee Not Found 409", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "employees"})

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: ErrEmployeeNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error ProductBatch Not Found 409", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "products_batches"})

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: ErrProductBatchNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error Warehouse Not Found 409", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "warehouses"})

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: ErrWarehouseNotFound.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error OrderNomberExists 409", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Order number exists"})

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: ErrOrderNumberExtists.Error(),
		}

		var result errorResponse
		err = json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Error DB 500", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		server := createServerInboundOrderFunctional(db)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
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
