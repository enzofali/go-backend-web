package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/purchaseorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// controller
type purchase_orderMock struct {
	mock.Mock
}

// constructor
func NewServicePurchaseOrderMock() *purchase_orderMock {
	return &purchase_orderMock{}
}

type errResponsePurchaseOrder struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (puOrdMock *purchase_orderMock) Create(ctx context.Context, purchord domain.Purchase_Orders) (domain.Purchase_Orders, error) {
	args := puOrdMock.Called(ctx, purchord)
	return args.Get(0).(domain.Purchase_Orders), args.Error(1)
}

func CreateServerPurchaseOrder(puOrdMock *purchase_orderMock) (engine *gin.Engine) {

	handler := NewPurchaseOrder(puOrdMock)

	engine = gin.Default()

	routerPurchaseOrder := engine.Group("/api/v1/purchaseorders")
	{
		routerPurchaseOrder.POST("", handler.Create())
	}

	return engine
}

func CreateReqPurchOrder(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {

	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	return req, httptest.NewRecorder()
}

//-----------TESTS-----------

func TestCreatePurchaseOrder(t *testing.T) {

	type responsePurchaseOrder struct {
		Data domain.Purchase_Orders
	}

	purchOrdResp := domain.Purchase_Orders{
		ID:                1,
		Order_number:      "abc123",
		Order_date:        "2023/12/12",
		Tracking_code:     "asd4321",
		Buyer_id:          1,
		Product_record_id: 12,
		Order_Status_id:   2,
	}

	purchOrdReq := domain.Purchase_Orders{
		Order_number:      "abc123",
		Order_date:        "2023/12/12",
		Tracking_code:     "asd4321",
		Buyer_id:          1,
		Product_record_id: 12,
		Order_Status_id:   2,
	}

	purchOrdDB := responsePurchaseOrder{
		Data: purchOrdResp,
	}

	t.Run("create purchase order success", func(t *testing.T) {
		//arrange
		service := NewServicePurchaseOrderMock()
		service.On("Create", mock.Anything, purchOrdReq).Return(purchOrdResp, nil)
		serv := CreateServerPurchaseOrder(service)
		code := http.StatusCreated

		//act
		req, resp := CreateReqPurchOrder(http.MethodPost, "/api/v1/purchaseorders",
			`{"order_number":"abc123","order_date":"2023/12/12","tracking_code":"asd4321","buyer_id":1,"product_record_id":12,"order_status_id":2}`)

		serv.ServeHTTP(resp, req)

		var purchaseOrder responsePurchaseOrder

		err := json.NewDecoder(resp.Body).Decode(&purchaseOrder)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, purchOrdDB, purchaseOrder)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("create purchase order format error", func(t *testing.T) {
		//arrange
		service := NewServicePurchaseOrderMock()
		serv := CreateServerPurchaseOrder(service)
		code := http.StatusUnprocessableEntity

		req, resp := CreateReqPurchOrder(http.MethodPost, "/api/v1/purchaseorders", `{"order_number":"abc123","order_date":"2023/12/12","tracking}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponsePurchaseOrder{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: purchaseorder.ErrFieldNotExist.Error(),
		}

		var eResp errResponsePurchaseOrder
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("created purchase order fields incomplete", func(t *testing.T) {
		//arrange
		service := NewServicePurchaseOrderMock()
		serv := CreateServerPurchaseOrder(service)
		code := http.StatusUnprocessableEntity

		//act
		req, resp := CreateReqPurchOrder(http.MethodPost, "/api/v1/purchaseorders",
			`{"order_date":"2023/12/12","tracking_code":"asd4321","buyer_id":1,"product_record_id":12,"order_status_id":2}`)

		serv.ServeHTTP(resp, req)

		errResp := errResponsePurchaseOrder{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: purchaseorder.ErrFieldNotExist.Error(),
		}

		var eResp errResponsePurchaseOrder
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("create purchase order conflict error", func(t *testing.T) {
		//arrange
		service := NewServicePurchaseOrderMock()
		service.On("Create", mock.Anything, purchOrdReq).Return(domain.Purchase_Orders{}, purchaseorder.ErrBuyerNotFound)
		serv := CreateServerPurchaseOrder(service)
		code := http.StatusConflict

		//act
		req, resp := CreateReqPurchOrder(http.MethodPost, "/api/v1/purchaseorders",
			`{"order_number":"abc123","order_date":"2023/12/12","tracking_code":"asd4321","buyer_id":1,"product_record_id":12,"order_status_id":2}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponsePurchaseOrder{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: purchaseorder.ErrBuyerNotFound.Error(),
		}

		var eResp errResponsePurchaseOrder

		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("create purchase order database error", func(t *testing.T) {
		//arrange
		service := NewServicePurchaseOrderMock()
		service.On("Create", mock.Anything, purchOrdReq).Return(domain.Purchase_Orders{}, purchaseorder.ErrDatabase)
		serv := CreateServerPurchaseOrder(service)
		code := http.StatusInternalServerError

		//act
		req, resp := CreateReqPurchOrder(http.MethodPost, "/api/v1/purchaseorders",
			`{"order_number":"abc123","order_date":"2023/12/12","tracking_code":"asd4321","buyer_id":1,"product_record_id":12,"order_status_id":2}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponsePurchaseOrder{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: purchaseorder.ErrDatabase.Error(),
		}

		var eResp errResponsePurchaseOrder

		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})
}
