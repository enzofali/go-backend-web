package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/buyer"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type buyerHandlerMock struct {
	mock.Mock
}

func NewServiceBuyerHandlerMock() *buyerHandlerMock {
	return &buyerHandlerMock{}
}

type errResponseBuyer struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (buhamock *buyerHandlerMock) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	args := buhamock.Called(ctx)
	return args.Get(0).([]domain.Buyer), args.Error(1)
}

func (buhamock *buyerHandlerMock) Get(ctx context.Context, id int) (domain.Buyer, error) {
	args := buhamock.Called(ctx, id)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (buhamock *buyerHandlerMock) Create(ctx context.Context, b domain.Buyer) (domain.Buyer, error) {
	args := buhamock.Called(ctx, b)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (buhamock *buyerHandlerMock) Update(ctx context.Context, b domain.Buyer, id int) error {
	args := buhamock.Called(ctx, b, id)
	return args.Error(0)
}

func (buhamock *buyerHandlerMock) Delete(ctx context.Context, id int) error {
	args := buhamock.Called(ctx, id)
	return args.Error(0)
}

func (buhamock *buyerHandlerMock) GetReports(ctx context.Context, id int) ([]domain.ReportBuyersPurchases, error) {
	args := buhamock.Called(ctx, id)
	return args.Get(0).([]domain.ReportBuyersPurchases), args.Error(1)
}

func CreateServerBuyer(buhamock *buyerHandlerMock) (engine *gin.Engine) {
	handler := NewBuyer(buhamock)

	engine = gin.Default()

	routerBuyer := engine.Group("/api/v1/buyers")
	{
		routerBuyer.GET("", handler.GetAll())
		routerBuyer.GET("/:id", handler.Get())
		routerBuyer.POST("", handler.Create())
		routerBuyer.PATCH(("/:id"), handler.Update())
		routerBuyer.DELETE(("/:id"), handler.Delete())
		routerBuyer.GET("/reportPurchaseOrders", handler.GetReport())
	}

	return engine
}

func CreateReqBuyer(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {

	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	return req, httptest.NewRecorder()

}

//---------- TESTS ---------

func TestGetAllBuyers(t *testing.T) {

	type responseBuyer struct {
		Data []domain.Buyer
	}

	buyers := []domain.Buyer{
		{
			ID:           1,
			CardNumberID: "123456",
			FirstName:    "Leandro",
			LastName:     "Villalba",
		},
		{
			ID:           2,
			CardNumberID: "456789",
			FirstName:    "Gonzalo",
			LastName:     "Alvarez",
		},
	}

	buyersData := responseBuyer{
		Data: buyers,
	}

	t.Run("find_all", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("GetAll", mock.Anything).Return(buyers, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusOK

		//action
		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers", "")
		serv.ServeHTTP(resp, req)

		var response responseBuyer

		err := json.NewDecoder(resp.Body).Decode(&response)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, buyersData, response)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Internal Error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("GetAll", mock.Anything).Return([]domain.Buyer{}, buyer.ErrDatabase)
		serv := CreateServerBuyer(service)
		code := http.StatusInternalServerError

		//action
		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers", "")
		serv.ServeHTTP(resp, req)

		var response responseBuyer

		err := json.NewDecoder(resp.Body).Decode(&response)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestGetBuyerById(t *testing.T) {

	type responseBuyer struct {
		Data domain.Buyer `json:"data"`
	}

	buyerRep := domain.Buyer{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "Leandro",
		LastName:     "Villalba",
	}

	buyerDB := responseBuyer{
		Data: buyerRep,
	}

	t.Run("find_by_id_existent", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 1).Return(buyerRep, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusOK

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/1", "")
		serv.ServeHTTP(resp, req)

		//act
		var response responseBuyer
		err := json.NewDecoder(resp.Body).Decode(&response)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, buyerDB, response)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(domain.Buyer{}, buyer.ErrNotFound)
		server := CreateServerBuyer(service)
		code := http.StatusNotFound

		//act
		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/2", "")
		server.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: buyer.ErrNotFound.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("ErrFormat Buyer Id Bad Request", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		server := CreateServerBuyer(service)
		code := http.StatusBadRequest

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/1a", "")
		server.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Get By ID Database Error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(domain.Buyer{}, buyer.ErrDatabase)
		server := CreateServerBuyer(service)
		code := http.StatusInternalServerError

		//act
		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/2", "")
		server.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: buyer.ErrDatabase.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestCreateBuyer(t *testing.T) {

	type responseBuyer struct {
		Data domain.Buyer
	}

	buyerResp := domain.Buyer{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "Leandro",
		LastName:     "Villalba",
	}

	buyerReq := domain.Buyer{
		CardNumberID: "123456",
		FirstName:    "Leandro",
		LastName:     "Villalba",
	}

	buyerDb := responseBuyer{
		Data: buyerResp,
	}

	t.Run("create_ok", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Create", mock.Anything, buyerReq).Return(buyerResp, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusCreated

		//act
		req, resp := CreateReqBuyer(http.MethodPost, "/api/v1/buyers",
			`{"card_number_id": "123456","first_name": "Leandro","last_name": "Villalba"}`)

		serv.ServeHTTP(resp, req)

		var buyerResponse responseBuyer
		err := json.NewDecoder(resp.Body).Decode(&buyerResponse)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, buyerDb, buyerResponse)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("create buyer format error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		serv := CreateServerBuyer(service)
		code := http.StatusBadRequest

		//act
		req, resp := CreateReqBuyer(http.MethodPost, "/api/v1/buyers", `{"card_number_id": "123456","first_name": "Leandro","last_name":`)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))

	})

	t.Run("create_fail", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		serv := CreateServerBuyer(service)
		code := http.StatusUnprocessableEntity

		//act
		req, resp := CreateReqBuyer(http.MethodPost, "/api/v1/buyers", `{"first_name": "Leandro","last_name": "Villalba"}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("create_conflict", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Create", mock.Anything, buyerReq).Return(domain.Buyer{}, buyer.ErrAlreadyExists)
		serv := CreateServerBuyer(service)
		code := http.StatusConflict

		//act
		req, resp := CreateReqBuyer(http.MethodPost, "/api/v1/buyers", `{"card_number_id": "123456","first_name": "Leandro","last_name": "Villalba"}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: buyer.ErrAlreadyExists.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Buyer Database Error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Create", mock.Anything, buyerReq).Return(domain.Buyer{}, buyer.ErrDatabase)
		serv := CreateServerBuyer(service)
		code := http.StatusInternalServerError

		//act
		req, resp := CreateReqBuyer(http.MethodPost, "/api/v1/buyers", `{"card_number_id": "123456","first_name": "Leandro","last_name": "Villalba"}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: buyer.ErrDatabase.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestDeleteBuyer(t *testing.T) {

	t.Run("delete_ok", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Delete", mock.Anything, 1).Return(nil)
		serv := CreateServerBuyer(service)
		code := http.StatusNoContent

		//act
		req, resp := CreateReqBuyer(http.MethodDelete, "/api/v1/buyers/1", "	")
		serv.ServeHTTP(resp, req)

		assert.Equal(t, code, resp.Code)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("delete_non_existent", func(t *testing.T) {

		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Delete", mock.Anything, 1).Return(buyer.ErrNotFound)
		serv := CreateServerBuyer(service)
		code := http.StatusNotFound

		//act
		req, resp := CreateReqBuyer(http.MethodDelete, "/api/v1/buyers/1", "")
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: buyer.ErrNotFound.Error(),
		}

		var eResp errResponseBuyer

		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Delete method with wrong id", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		server := CreateServerBuyer(service)
		code := http.StatusBadRequest

		req, resp := CreateReqBuyer(http.MethodDelete, "/api/v1/buyers/1a", "")
		server.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Delete method database error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Delete", mock.Anything, 1).Return(buyer.ErrDatabase)
		server := CreateServerBuyer(service)
		code := http.StatusInternalServerError

		//act
		req, resp := CreateReqBuyer(http.MethodDelete, "/api/v1/buyers/1", "")
		server.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: buyer.ErrDatabase.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestReportPurchaseOrdersForAllBuyers(t *testing.T) {

	type responsePurchaseOrdersForAllBuyers struct {
		Data []domain.ReportBuyersPurchases
	}

	buyersWithPurchases := []domain.ReportBuyersPurchases{
		{
			ID:                   1,
			Card_Number_ID:       "123456",
			First_Name:           "Leandro",
			Last_Name:            "Villalba",
			Purchase_Order_Count: 2,
		},
		{
			ID:                   2,
			Card_Number_ID:       "456",
			First_Name:           "Gonzalo",
			Last_Name:            "Alvarez",
			Purchase_Order_Count: 3,
		},
	}

	buyersResponse := responsePurchaseOrdersForAllBuyers{
		Data: buyersWithPurchases,
	}

	t.Run("find_all", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("GetReports", mock.Anything, 0).Return(buyersWithPurchases, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusOK

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders", "")
		serv.ServeHTTP(resp, req)

		//act
		var expectedResponse responsePurchaseOrdersForAllBuyers
		err := json.Unmarshal(resp.Body.Bytes(), &expectedResponse)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, buyersResponse, expectedResponse)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("buyer id format error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		serv := CreateServerBuyer(service)
		code := http.StatusBadRequest

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders?id=3a", "")
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))

	})

	t.Run("purchase by buyer not found", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("GetReports", mock.Anything, 3).Return([]domain.ReportBuyersPurchases{}, buyer.ErrPurchaseNotFound)
		serv := CreateServerBuyer(service)
		code := http.StatusNotFound

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders?id=3", "")
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: buyer.ErrPurchaseNotFound.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))

	})

	t.Run("purchase by buyers not found", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("GetReports", mock.Anything, 0).Return([]domain.ReportBuyersPurchases{}, buyer.ErrPurchasesNotFound)
		serv := CreateServerBuyer(service)
		code := http.StatusNotFound

		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders", "")
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: buyer.ErrPurchasesNotFound.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))

	})

	t.Run("Database error", func(t *testing.T) {
		service := NewServiceBuyerHandlerMock()
		service.On("GetReports", mock.Anything, 0).Return([]domain.ReportBuyersPurchases{}, buyer.ErrDatabase)
		serv := CreateServerBuyer(service)
		code := http.StatusInternalServerError

		//action
		req, resp := CreateReqBuyer(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders", "")
		serv.ServeHTTP(resp, req)

		var response responsePurchaseOrdersForAllBuyers

		err := json.NewDecoder(resp.Body).Decode(&response)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestUpdateBuyer(t *testing.T) {

	type buyerResponse struct {
		Data domain.Buyer `json:"data"`
	}

	buyerNow := domain.Buyer{
		ID:           1,
		CardNumberID: "abc123",
		FirstName:    "Leandro",
		LastName:     "Villalba",
	}

	buyerAfter := domain.Buyer{
		ID:           1,
		CardNumberID: "abc123",
		FirstName:    "Alejo",
		LastName:     "Gonzalez",
	}

	buyerBD := buyerResponse{
		Data: buyerAfter,
	}

	t.Run("update_ok", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 1).Return(buyerNow, nil)
		service.On("Update", mock.Anything, buyerAfter, 1).Return(nil)
		serv := CreateServerBuyer(service)
		code := http.StatusOK

		//act
		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/1", `{"card_number_id":"abc123","first_name": "Alejo", "last_name": "Gonzalez"}`)
		serv.ServeHTTP(resp, req)

		var expBuy buyerResponse
		fmt.Println("response: ", resp.Body.Bytes())
		err := json.Unmarshal(resp.Body.Bytes(), &expBuy)
		//err := json.NewDecoder(resp.Body).Decode(&expBuy)
		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, buyerBD, expBuy)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("id format error", func(t *testing.T) {
		service := NewServiceBuyerHandlerMock()
		serv := CreateServerBuyer(service)
		code := http.StatusBadRequest

		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/reportPurchaseOrders?id=3a", "")
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: buyer.ErrFormat.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("update_non_existent", func(t *testing.T) {

		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(domain.Buyer{}, buyer.ErrNotFound)
		serv := CreateServerBuyer(service)
		code := http.StatusNotFound
		//act
		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/2", `{"card_number_id":"abc12","first_name": "Alejo", "last_name": "Gonzalez"}`)
		serv.ServeHTTP(resp, req)

		errResponse := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: buyer.ErrNotFound.Error(),
		}

		var errResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&errResp)

		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResponse, errResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("method update with json error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(buyerNow, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusConflict

		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/2", `{"card_number_id":"123456","first_name": "Alejo", "last_name": `)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: buyer.ErrCantChange.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("method update with json error", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(buyerNow, nil)
		serv := CreateServerBuyer(service)
		code := http.StatusConflict

		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/2", `{"card_number_id":"1234","first_name": "Alejo", "last_name": `)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: buyer.ErrCantChange.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("update error cant change", func(t *testing.T) {
		//arrange
		service := NewServiceBuyerHandlerMock()
		service.On("Get", mock.Anything, 2).Return(buyerNow, nil)
		//service.On("Update", mock.Anything, buyerAfter,2).Return(buyer.ErrCantChange)
		serv := CreateServerBuyer(service)
		code := http.StatusConflict

		//act
		req, resp := CreateReqBuyer(http.MethodPatch, "/api/v1/buyers/2", `{"card_number_id":"","first_name": "Alejo","last_name": "Gonzalez"}`)
		serv.ServeHTTP(resp, req)

		errResp := errResponseBuyer{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: buyer.ErrCantChange.Error(),
		}

		var eResp errResponseBuyer
		err := json.NewDecoder(resp.Body).Decode(&eResp)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, code, resp.Code)
		assert.Equal(t, errResp, eResp)
		assert.True(t, service.AssertExpectations(t))
	})

}
