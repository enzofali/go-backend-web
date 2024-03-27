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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/carry"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorResponseCarry struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type serviceCarryTest struct {
	mock.Mock
}

//constructor

func NewServiceCarryTest() *serviceCarryTest {
	return &serviceCarryTest{}
}

func (ca *serviceCarryTest) GetAll(ctx context.Context) ([]domain.Carrie, error) {
	args := ca.Called(ctx)
	return args.Get(0).([]domain.Carrie), args.Error(1)
}
func (ca *serviceCarryTest) GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error) {

	args := ca.Called(ctx)
	return args.Get(0).([]domain.CarrieLocality), args.Error(1)
}
func (ca *serviceCarryTest) GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error) {
	args := ca.Called(ctx, id)
	return args.Get(0).(domain.CarrieLocality), args.Error(1)
}
func (ca *serviceCarryTest) Crear(ctx context.Context, c domain.Carrie) (carriesG domain.Carrie, err error) {
	args := ca.Called(ctx, c)
	return args.Get(0).(domain.Carrie), args.Error(1)
}

/*func (ca *serviceCarryTest) GetQuery(idd string) (id string, ok bool) {
	args := ca.Called(idd)
	return args.Get(0).(string), args.Get(1).(bool)
}*/

func CreateServerCarries(service *serviceCarryTest) *gin.Engine {
	handler := NewCarry(service)

	eng := gin.Default()

	carryG := eng.Group("/api/v1/carries")
	{
		carryG.GET("", handler.GetAll())  //http://localhost:8080/api/v1/carries/
		carryG.POST("", handler.Create()) //http://localhost:8080/api/v1/carries/

	}

	eng.GET("/api/v1/localities/reportCarries", handler.GetAllByLocality()) //http://localhost:8080/api/v1/localities/reportCarries?id=2001
	//eng.GET("/api/v1/localities/reportCarries?id=C05", handler.GetAllByLocality())
	return eng
}

func createRequestCarry(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func TestGetallCarryH(t *testing.T) {

	type response struct {
		Data []domain.Carrie `json:"data"`
	}

	carriess := []domain.Carrie{

		{
			Id:           2,
			Cid:          "ABC34",
			Company_name: "Servientrega",
			Address:      "Call 40 # 23 -50",
			Telephone:    "2256789",
			Locality_id:  "C05",
		},

		{
			Id:           3,
			Cid:          "ABC35",
			Company_name: "DHL",
			Address:      "Call 22 # 21 -50",
			Telephone:    "2247791",
			Locality_id:  "C05",
		},
	}

	data := response{
		Data: carriess,
	}

	//si la peticion es exitosa retorna 200 y un listado de los carry existentes
	t.Run("find_all", func(t *testing.T) {

		service := NewServiceCarryTest()
		service.On("GetAll", mock.Anything).Return(carriess, nil)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/carries", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))

	})

	//genera error 500 si algo sale mal en la BD, forzamos el error con ErrBD
	t.Run("find_all_err_500", func(t *testing.T) {

		service := NewServiceCarryTest()

		service.On("GetAll", mock.Anything).Return([]domain.Carrie{}, carry.ErrBD)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/carries", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: carry.ErrBD.Error(),
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

}

func TestGetallLocalityCarryH(t *testing.T) {

	type response struct {
		Data []domain.CarrieLocality `json:"data"`
	}

	carriess := []domain.CarrieLocality{

		{
			Locality_id:   "C05",
			Locality_name: "Cartagena",
			Cant_carries:  5,
		},

		{
			Locality_id:   "C07",
			Locality_name: "San Andres",
			Cant_carries:  2,
		},
	}

	data := response{
		Data: carriess,
	}

	//caso no se envia id,retorna un 200 y todos las localidades con cant carry
	t.Run("find_all", func(t *testing.T) {

		service := NewServiceCarryTest()
		service.On("GetByLocality", mock.Anything).Return(carriess, nil)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/localities/reportCarries", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))

	})

}

func TestGetLocalityIdlCarryH(t *testing.T) {

	type response struct {
		Data domain.CarrieLocality `json:"data"`
	}

	carryM := domain.CarrieLocality{
		Locality_id:   "C05",
		Locality_name: "Cartagena",
		Cant_carries:  5,
	}

	data := response{
		Data: carryM,
	}

	//caso se envia id,retorna un 200 y una localidad determinada con cant carry
	t.Run("find_by_id_existent", func(t *testing.T) {

		service := NewServiceCarryTest()
		//service.On("GetQuery", "C05").Return("C05", true)
		service.On("GetByLocalityID", mock.Anything, "C05").Return(carryM, nil)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/localities/reportCarries?id=C05", "")
		server.ServeHTTP(resp, req)

		var result response

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))

	})

	//si el id solicitado no existe retorna un 404 y una structura Carry vacia
	t.Run("find_by_id_non_existent", func(t *testing.T) {

		service := NewServiceCarryTest()
		service.On("GetByLocalityID", mock.Anything, "C07").Return(domain.CarrieLocality{}, carry.ErrNotFound)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/localities/reportCarries?id=C07", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: carry.ErrNotFound.Error(),
		}

		var result errorResponseCarry

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusNotFound)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))

	})

	//genera error 500 si algo sale mal en la BD, Forzamos el error con ErrBD
	t.Run("find_by_id_err_500", func(t *testing.T) {

		service := NewServiceCarryTest()

		service.On("GetAll", mock.Anything).Return([]domain.Carrie{}, carry.ErrBD)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodGet, "/api/v1/carries", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: carry.ErrBD.Error(),
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestCrearCarryH(t *testing.T) {

	type response struct {
		Data domain.Carrie `json:"data"`
	}

	carriesResp := domain.Carrie{

		Id:           2,
		Cid:          "ABC34",
		Company_name: "Servientrega",
		Address:      "Call 40 # 23 -50",
		Telephone:    "2256789",
		Locality_id:  "C05",
	}

	carriesReq := domain.Carrie{

		Cid:          "ABC34",
		Company_name: "Servientrega",
		Address:      "Call 40 # 23 -50",
		Telephone:    "2256789",
		Locality_id:  "C05",
	}

	data := response{
		Data: carriesResp,
	}

	//Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.

	t.Run("create_ok", func(t *testing.T) {

		service := NewServiceCarryTest()
		service.On("Crear", mock.Anything, carriesReq).Return(carriesResp, nil)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `{"cid": "ABC34","company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)

		server.ServeHTTP(resp, req)

		var result response

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusCreated)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//Si el objeto no contiene formato Json retona 422
	t.Run("create_fail", func(t *testing.T) {
		service := NewServiceCarryTest()
		server := CreateServerCarries(service)

		//se envia json invalido
		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `"cid": "ABC34","company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "json: cannot unmarshal string into Go value of type domain.Carrie",
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//Si el cid de carry ya existe retorna un error 409 Conflict , forzamos este error con ErrExist
	t.Run("create_conflict", func(t *testing.T) {
		service := NewServiceCarryTest()
		service.On("Crear", mock.Anything, carriesReq).Return(domain.Carrie{}, carry.ErrExist)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `{"cid": "ABC34","company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: carry.ErrExist.Error(),
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//si se envia campo vacio retorna un 400
	t.Run("create_fail_400", func(t *testing.T) {
		service := NewServiceCarryTest()
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `{"company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "bad request",
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
	//si algo sale mal en la BD retorna un 500, forzamos el error enviando ErrBD
	t.Run("create_fail_500", func(t *testing.T) {

		service := NewServiceCarryTest()
		service.On("Crear", mock.Anything, carriesReq).Return(domain.Carrie{}, carry.ErrBD)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `{"cid": "ABC34","company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: carry.ErrBD.Error(),
		}

		var result errorResponseCarry

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//Si el la FK (Locality_id) no existe retorna un error 409 Conflict

	t.Run("create_conflict_fk_409", func(t *testing.T) {
		service := NewServiceCarryTest()
		service.On("Crear", mock.Anything, carriesReq).Return(domain.Carrie{}, carry.ErrForeignKey)
		server := CreateServerCarries(service)

		req, resp := createRequestCarry(http.MethodPost, "/api/v1/carries", `{"cid": "ABC34","company_name":"Servientrega","address": "Call 40 # 23 -50","telephone":"2256789","locality_id":  "C05"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseCarry{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: carry.ErrForeignKey.Error(),
		}

		var result errorResponseCarry
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

}
