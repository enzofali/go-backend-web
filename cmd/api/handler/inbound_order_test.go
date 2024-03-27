package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	inboundorder "github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/inbound_order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceInboundOrderMock struct {
	mock.Mock
}

func NewServiceInboundOrderMock() *serviceInboundOrderMock {
	return &serviceInboundOrderMock{}
}

func (sm *serviceInboundOrderMock) Create(ctx context.Context, i domain.InboundOrder) (domain.InboundOrder, error) {
	args := sm.Called(ctx, i)
	return args.Get(0).(domain.InboundOrder), args.Error(1)
}

func createServerInboundOrderUnit(service *serviceInboundOrderMock) *gin.Engine {
	handler := NewInoudOrder(service)

	eng := gin.Default()

	rIO := eng.Group("/api/v1/inboundOrders")
	{
		rIO.POST("", handler.Create())
	}

	return eng
}

func createRequestInboundOrderUnit(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func Test_InboundOrder_Create(t *testing.T) {
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

	inboundOrderReq := domain.InboundOrder{
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
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(inboundOrderResp, nil)
		server := createServerInboundOrderUnit(service)

		req, resp := createRequestInboundOrderUnit(http.MethodPost, "/api/v1/inboundOrders",
			`{"order_date": "2006-01-02",
		"order_number": "order1",
		"employee_id": 1,
		"product_batch_id": 1,
		"warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Json 422", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Validator 422", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Date 400", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Employee Not Found 409", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(domain.InboundOrder{}, inboundorder.ErrEmployeeNotFound)
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error ProductBatch Not Found 409", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(domain.InboundOrder{}, inboundorder.ErrProductBatchNotFound)
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Warehouse Not Found 409", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(domain.InboundOrder{}, inboundorder.ErrWarehouseNotFound)
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error OrderNomberExists 409", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(domain.InboundOrder{}, inboundorder.ErrOrderNumberExtists)
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error DB 500", func(t *testing.T) {
		service := NewServiceInboundOrderMock()
		service.On("Create", mock.Anything, inboundOrderReq).Return(domain.InboundOrder{}, errors.New("error database"))
		server := createServerInboundOrderUnit(service)

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
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}
