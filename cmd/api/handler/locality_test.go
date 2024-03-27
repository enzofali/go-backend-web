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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/locality"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorResponseLocality struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ______________________________________________________
// tools
type serviceMockLocality struct {
	mock.Mock
}

func NewServiceMockLocality() *serviceMockLocality {
	return &serviceMockLocality{}
}

// All methods simply return the struct's initially defined members
func (r *serviceMockLocality) Create(ctx context.Context, l domain.Locality) error {
	args := r.Called(ctx, l)
	return args.Error(0)
}
func (r *serviceMockLocality) GetSellerAll(ctx context.Context) ([]domain.QuantitySellerByLocality, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.QuantitySellerByLocality), args.Error(1)
}
func (r *serviceMockLocality) GetSellerByLocality(ctx context.Context, id string) (domain.QuantitySellerByLocality, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.QuantitySellerByLocality), args.Error(1)
}

func CreateServerLocality(service locality.Service) *gin.Engine {
	// instances
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
func NewRequestLocality(method, path, body string) (req *http.Request, res *httptest.ResponseRecorder) {
	// request
	req = httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// response
	res = httptest.NewRecorder()

	return
}

func Test_Create_Locality(t *testing.T) {
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
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)

		service.On("Create", mock.Anything, localityToCreate).Return(nil)

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, data, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("Bad request error", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "error bad request",
		}

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseLocality
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
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "Key: 'Locality.Id' Error:Field validation for 'Id' failed on the 'required' tag",
		}

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `{"locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var sellerResult errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &sellerResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, sellerResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("internal error creating", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}

		service.On("Create", mock.Anything, localityToCreate).Return(locality.ErrIntern)

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("id duplicate error creating", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: locality.ErrDuplicated.Error(),
		}

		service.On("Create", mock.Anything, localityToCreate).Return(locality.ErrDuplicated)

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("default error creating", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: "internal error",
		}

		service.On("Create", mock.Anything, localityToCreate).Return(locality.ErrLocalityNotFound)

		request, response := NewRequestLocality(http.MethodPost, "/api/v1/localities", `{"id": "6701", "locality_name": "Villa Crespo", "province_name": "Buenos Aires", "country_name": "Argentina"}`)

		// act
		server.ServeHTTP(response, request)
		var localityResult errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &localityResult)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, localityResult)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}

func Test_GetQuantitySellerByLocality(t *testing.T) {
	// arrange
	type responseStruct struct {
		Data []domain.QuantitySellerByLocality `json:"data"`
	}

	t.Run("OK report the number of sellers of all locations", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		resultExpected := []domain.QuantitySellerByLocality{
			{Locality_id: "6701", Locality_name: "Villa Crespo", Sellers_count: 6},
			{Locality_id: "6702", Locality_name: "Nuñez", Sellers_count: 3},
		}
		data := responseStruct{Data: resultExpected}
		service.On("GetSellerAll", mock.Anything).Return(resultExpected, nil)

		request, response := NewRequestLocality(http.MethodGet, "/api/v1/localities/reportSellers", "")

		// act
		server.ServeHTTP(response, request)
		var report responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, data, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
	//NO FUNCIONE
	t.Run("internal error getting number of sellers from all locations", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}
		//var expected []domain.QuantitySellerByLocality
		service.On("GetSellerAll", mock.Anything).Return([]domain.QuantitySellerByLocality{}, locality.ErrIntern)

		request, response := NewRequestLocality(http.MethodGet, "/api/v1/localities/reportSellers", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("OK report of the number of sellers by location", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		resultExpected := domain.QuantitySellerByLocality{Locality_id: "6701", Locality_name: "Villa Crespo", Sellers_count: 6}
		type responseStruct struct {
			Data domain.QuantitySellerByLocality `json:"data"`
		}
		data := responseStruct{Data: resultExpected}
		service.On("GetSellerByLocality", mock.Anything, "6701").Return(resultExpected, nil)

		request, response := NewRequestLocality(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report responseStruct
		err := json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, data, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("not found error when obtaining the number of sellers in a location", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)

		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: locality.ErrIntern.Error(),
		}
		service.On("GetSellerByLocality", mock.Anything, "6701").Return(domain.QuantitySellerByLocality{}, locality.ErrIntern)

		request, response := NewRequestLocality(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})

	t.Run("not found error when obtaining the number of sellers in a location", func(t *testing.T) {
		// arrange
		service := NewServiceMockLocality()
		server := CreateServerLocality(service)
		//resultExpected := domain.QuantitySellerByLocality{Locality_id: "6701", Locality_name: "Villa Crespo", Sellers_count: 6}
		errResp := errorResponseLocality{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: locality.ErrLocalityNotFound.Error(),
		}
		service.On("GetSellerByLocality", mock.Anything, "6701").Return(domain.QuantitySellerByLocality{}, locality.ErrLocalityNotFound)

		request, response := NewRequestLocality(http.MethodGet, "/api/v1/localities/reportSellers?id=6701", "")

		// act
		server.ServeHTTP(response, request)
		var report errorResponseLocality
		err := json.Unmarshal(response.Body.Bytes(), &report)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.True(t, service.AssertExpectations(t))
		assert.Equal(t, errResp, report)
		assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	})
}
