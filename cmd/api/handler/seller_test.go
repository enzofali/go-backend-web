package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/seller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorResponseSeller struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ______________________________________________________
// tools
type serviceMockSeller struct {
	mock.Mock
}

func NewserviceMockSeller() *serviceMockSeller {
	return &serviceMockSeller{}
}

func (r *serviceMockSeller) GetAll(ctx context.Context) ([]domain.Seller, error) {
	args := r.Mock.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}
func (r *serviceMockSeller) GetByID(ctx context.Context, id int) (domain.Seller, error) {
	args := r.Mock.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (r *serviceMockSeller) Create(ctx context.Context, s domain.Seller) (int, error) {
	args := r.Mock.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}
func (r *serviceMockSeller) Update(ctx context.Context, s domain.Seller) error {
	args := r.Mock.Called(ctx, s)
	return args.Error(0)
}
func (r *serviceMockSeller) Delete(ctx context.Context, id int) error {
	args := r.Mock.Called(ctx, id)
	return args.Error(0)
}
func (r *serviceMockSeller) Exists(ctx context.Context, cid int) bool {
	args := r.Mock.Called(ctx, cid)
	return args.Get(0).(bool)
}

// ______________________________________________________
// tools
func CreateServerSeller(service seller.Service) *gin.Engine {
	// instances
	handler := NewSeller(service)

	// server
	server := gin.Default()

	// -> routes
	routes := server.Group("/api/v1/sellers")
	{
		routes.GET("/", handler.GetAll())
		routes.POST("/", handler.Create())
		routes.GET("/:id", handler.Get())
		routes.PATCH("/:id", handler.Update())
		routes.DELETE("/:id", handler.Delete())
	}

	return server
}
func NewRequestSeller(method, path, body string) (req *http.Request, res *httptest.ResponseRecorder) {
	// request
	req = httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// response
	res = httptest.NewRecorder()

	return
}

func Test_Create_Seller(t *testing.T) {
	// arrange
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
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("Create", mock.Anything, sellerToCreate).Return(1, nil)

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, data, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Bad request error", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "error bad request",
		}

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Validator error", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Seller.CID' Error:Field validation for 'CID' failed on the 'required' tag",
		}

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Cid invalid error", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "invalid cid",
		}

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"cid": -1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("conflict error creating", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: seller.ErrConflict.Error(),
		}

		service.On("Create", mock.Anything, sellerToCreate).Return(0, seller.ErrConflict)

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error creating", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		service.On("Create", mock.Anything, sellerToCreate).Return(0, seller.ErrIntern)

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("not found error creating", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrInvalidLocality.Error(),
		}

		service.On("Create", mock.Anything, sellerToCreate).Return(0, seller.ErrInvalidLocality)

		request, response := NewRequestSeller(http.MethodPost, "/api/v1/sellers/", `{"cid": 1,"company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_GetAll_Seller(t *testing.T) {
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
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("GetAll", mock.Anything).Return(sellersExpected, nil)

		request, response := NewRequestSeller(http.MethodGet, "/api/v1/sellers/", "")

		// act
		server.ServeHTTP(response, request)
		var sellers responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellers)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 200, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, data, sellers)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error internal", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		service.On("GetAll", mock.Anything).Return([]domain.Seller{}, seller.ErrIntern)

		request, response := NewRequestSeller(http.MethodGet, "/api/v1/sellers/", "")

		// act
		server.ServeHTTP(response, request)
		var result errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, result)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Get_Seller(t *testing.T) {
	type responseStruct struct {
		Data domain.Seller `json:"data"`
	}
	sellerExpected := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	data := responseStruct{
		Data: sellerExpected,
	}
	t.Run("OK", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		id := 1

		service.On("GetByID", mock.Anything, id).Return(sellerExpected, nil)

		request, response := NewRequestSeller(http.MethodGet, "/api/v1/sellers/1", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, data, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		id := 1
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrNotFound.Error(),
		}

		service.On("GetByID", mock.Anything, id).Return(domain.Seller{}, seller.ErrNotFound)

		request, response := NewRequestSeller(http.MethodGet, "/api/v1/sellers/1", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Error bad request", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		//id := 1
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "strconv.Atoi: parsing \"fhdfh\": invalid syntax",
		}

		request, response := NewRequestSeller(http.MethodGet, "/api/v1/sellers/fhdfh", "")

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

}

func Test_Update_Seller(t *testing.T) {
	type responseStruct struct {
		Data domain.Seller `json:"data"`
	}
	sellerDB := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	sellerUpdated := domain.Seller{ID: 1, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
	data := responseStruct{
		Data: sellerUpdated,
	}
	t.Run("OK", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)
		service.On("Update", mock.Anything, sellerUpdated).Return(nil)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, data, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error not found", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: seller.ErrNotFound.Error(),
		}

		service.On("GetByID", mock.Anything, 1).Return(domain.Seller{}, seller.ErrNotFound)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error invalid id in the request (atoi)", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "strconv.Atoi: parsing \"hola\": invalid syntax",
		}

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/hola", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Error newDecoder", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "bad request",
		}

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("the id does not match the id of the db", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "cannot update product id",
		}
		sDb := domain.Seller{ID: 2, CID: 1, CompanyName: "Mercado Libre", Address: "Ramallo 6023", Telephone: "48557589", Locality_id: "6700"}
		service.On("GetByID", mock.Anything, 1).Return(sDb, nil)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error: validator", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Seller.Locality_id' Error:Field validation for 'Locality_id' failed on the 'required' tag",
		}

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": ""}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Error: cid invalid", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "invalid cid",
		}

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"cid": -1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	t.Run("Error conflic", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: seller.ErrConflict.Error(),
		}

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)
		service.On("Update", mock.Anything, sellerUpdated).Return(seller.ErrConflict)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error updating", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)
		errResp := errorResponseSeller{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: seller.ErrIntern.Error(),
		}

		service.On("GetByID", mock.Anything, 1).Return(sellerDB, nil)
		service.On("Update", mock.Anything, sellerUpdated).Return(seller.ErrIntern)

		request, response := NewRequestSeller(http.MethodPatch, "/api/v1/sellers/1", `{"id": 1, "cid": 1, "company_name": "Mercado Libre", "address": "Ramallo 6023", "telephone": "48557589", "locality_id": "6700"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseSeller
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, errResp, sellerResult)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_Delete_Seller(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("Delete", mock.Anything, 1).Return(nil)

		request, response := NewRequestSeller(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("bad request err (atoi id)", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		request, response := NewRequestSeller(http.MethodDelete, "/api/v1/sellers/hola", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Not found seller error", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("Delete", mock.Anything, 1).Return(seller.ErrNotFound)

		request, response := NewRequestSeller(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Internal error", func(t *testing.T) {
		// arrange
		service := NewserviceMockSeller()
		server := CreateServerSeller(service)

		service.On("Delete", mock.Anything, 1).Return(seller.ErrIntern)

		request, response := NewRequestSeller(http.MethodDelete, "/api/v1/sellers/1", "")
		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}
