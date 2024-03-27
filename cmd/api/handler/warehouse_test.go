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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorResponseWarehouse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type serviceWarehouseTest struct {
	mock.Mock
}

// constructor

func NewServiceWarehouseTest() *serviceWarehouseTest {
	return &serviceWarehouseTest{}
}

// Prepacion de metodos
func (s *serviceWarehouseTest) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (s *serviceWarehouseTest) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (s *serviceWarehouseTest) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	args := s.Called(ctx, w)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (s *serviceWarehouseTest) Update(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	args := s.Called(ctx, w)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (s *serviceWarehouseTest) Delete(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func CreateServerWarehouses(service *serviceWarehouseTest) *gin.Engine {
	handler := NewWarehouse(service)

	eng := gin.Default()

	rWareH := eng.Group("/api/v1/warehouses")

	{
		rWareH.GET("", handler.GetAll())
		rWareH.GET("/:id", handler.Get())
		rWareH.POST("", handler.Create())
		rWareH.PATCH("/:id", handler.Update())
		rWareH.DELETE("/:id", handler.Delete())
	}
	return eng
}

func createRequestWarehouse(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

// cases Tests
func TestGetAllWHandler(t *testing.T) {

	type response struct {
		Data []domain.Warehouse `json:"data"`
	}

	warehousess := []domain.Warehouse{
		{
			ID:                 1,
			Address:            "call3 40 # 3 -23",
			Telephone:          "2345678",
			WarehouseCode:      "ABC007",
			MinimumCapacity:    10,
			MinimumTemperature: 15,
		},

		{
			ID:                 2,
			Address:            "call3 4 # 32 -23",
			Telephone:          "2324678",
			WarehouseCode:      "ABC008",
			MinimumCapacity:    10,
			MinimumTemperature: 15,
		},
	}

	data := response{
		Data: warehousess,
	}
	//si la peticion es exitosa retorna 200 y un listado de los warehouses existentes
	t.Run("find_all", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("GetAll", mock.Anything).Return(warehousess, nil)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodGet, "/api/v1/warehouses", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//generar error 500 si algo sale mal en la BD, se envia ErrBD para forzar el error
	t.Run("find_all_err", func(t *testing.T) {

		service := NewServiceWarehouseTest()

		service.On("GetAll", mock.Anything).Return([]domain.Warehouse{}, warehouse.ErrBD)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodGet, "/api/v1/warehouses", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: warehouse.ErrBD.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func TestGetWHandler(t *testing.T) {
	type response struct {
		Data domain.Warehouse `json:"data"`
	}

	warehousess := domain.Warehouse{
		ID:                 3,
		Address:            "calle 14 # 32 -23",
		Telephone:          "2224678",
		WarehouseCode:      "ABC009",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	data := response{
		Data: warehousess,
	}

	//Si el id existe retorna un 200 y el warehouse solicitado
	t.Run("find_by_id_existent", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Get", mock.Anything, 3).Return(warehousess, nil)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodGet, "/api/v1/warehouses/3", "")
		server.ServeHTTP(resp, req)

		var result response

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))

	})
	//si el id solicitado no existe retorna un 404 y una structura warehouse vacia
	t.Run("find_by_id_non_existent", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Get", mock.Anything, 3).Return(domain.Warehouse{}, warehouse.ErrNotFound)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodGet, "/api/v1/warehouses/3", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: warehouse.ErrNotFound.Error(),
		}

		var result errorResponseWarehouse

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusNotFound)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))

	})

	//retorna un 400 si se envia un id invalido, ej un string
	t.Run("Get Error 400", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodGet, "/api/v1/warehouses/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: warehouse.ErrBadRequest.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

}

func TestCreateWHandler(t *testing.T) {

	type response struct {
		Data domain.Warehouse `json:"data"`
	}

	wareHResp := domain.Warehouse{
		ID:                 3,
		Address:            "calle 14 # 32 -23",
		Telephone:          "2224678",
		WarehouseCode:      "ABC009",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	wareHReq := domain.Warehouse{

		Address:            "calle 14 # 32 -23",
		Telephone:          "2224678",
		WarehouseCode:      "ABC009",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	data := response{
		Data: wareHResp,
	}

	// Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.
	t.Run("create_ok", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Create", mock.Anything, wareHReq).Return(wareHResp, nil)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `{"address": "calle 14 # 32 -23", "telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15 }`)

		server.ServeHTTP(resp, req)

		var result response

		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, resp.Code, http.StatusCreated)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//Si el objeto  no contiene formato Json devolverá un código 422
	t.Run("create_fail", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		//se envia json invalido
		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `"address": "calle 14 # 32 -23","telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "json: cannot unmarshal string into Go value of type domain.Warehouse",
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//Si el warehouse_code ya existe devuelve un error 409 Conflict
	t.Run("create_conflict", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		service.On("Create", mock.Anything, wareHReq).Return(domain.Warehouse{}, warehouse.ErrExist)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `{"address": "calle 14 # 32 -23", "telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15 }`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusConflict)), " ", "_"),
			Message: warehouse.ErrExist.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})
	//falla si se envia campo vacio y retorna un 400
	t.Run("create_fail_400", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `{"telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "bad request",
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//si algo sale mal en la BD retorna un 500, forzamos el error enviando ErrBD
	t.Run("create_fail_500", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		service.On("Create", mock.Anything, wareHReq).Return(domain.Warehouse{}, warehouse.ErrBD)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `{"address": "calle 14 # 32 -23", "telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15 }`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: warehouse.ErrBD.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

}

func TestUpdateWHandler(t *testing.T) {

	type response struct {
		Data domain.Warehouse `json:"data"`
	}

	wareHActual := domain.Warehouse{
		ID:                 9,
		Address:            "Calle 14 # 32 -23",
		Telephone:          "2224678",
		WarehouseCode:      "ABC009",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	wareHNuevo := domain.Warehouse{
		ID:                 9,
		Address:            "Calle 14 # 1-23",
		Telephone:          "2224678",
		WarehouseCode:      "ABC009",
		MinimumCapacity:    10,
		MinimumTemperature: 15,
	}

	data := response{
		Data: wareHNuevo,
	}
	// Cuando la actualización de datos sea exitosa se devolverá el warehouse con la información actualizada junto con un código 200
	t.Run("update_ok", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Get", mock.Anything, 9).Return(wareHActual, nil)
		service.On("Update", mock.Anything, wareHNuevo).Return(wareHNuevo, nil)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPatch, "/api/v1/warehouses/9", `{"ID":9,"Address": "Calle 14 # 1-23","Telephone": "2224678", "WarehouseCode":"ABC009","MinimumCapacity":10,"MinimumTemperature": 15}`)

		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))

	})
	// Si el warehouse que se desea actualizar no existe se devolverá un código 404
	t.Run("update_non_existent", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		service.On("Get", mock.Anything, 1).Return(domain.Warehouse{}, warehouse.ErrNotFound)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPatch, "/api/v1/warehouses/1", `{"ID":9,"Address": "Calle 14 # 1-23","Telephone": "2224678", "WarehouseCode":"ABC009","MinimumCapacity":10,"MinimumTemperature": 15}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: warehouse.ErrNotFound.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
	//si se envia campo vacio , retorna un  400
	t.Run("Update_fail_400", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPatch, "/api/v1/warehouses/abc", `{"telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15 }`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: warehouse.ErrBadRequest.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//si se envia json invalido falal y retorna un 422
	t.Run("Update_fail_422", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPost, "/api/v1/warehouses", `"address": "calle 14 # 32 -23","telephone": "2224678", "warehouse_code": "ABC009", "minimum_capacity": 10, "minimum_temperature": 15}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "json: cannot unmarshal string into Go value of type domain.Warehouse",
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	//si algo sale mal en la BD retorna un 500 . forzamos el error con la variable ErrBd
	t.Run("Update_fail_500", func(t *testing.T) {
		service := NewServiceWarehouseTest()
		service.On("Get", mock.Anything, 9).Return(wareHActual, nil)
		service.On("Update", mock.Anything, wareHNuevo).Return(domain.Warehouse{}, warehouse.ErrBD)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodPatch, "/api/v1/warehouses/9", `{"ID":9,"Address": "Calle 14 # 1-23","Telephone": "2224678", "WarehouseCode":"ABC009","MinimumCapacity":10,"MinimumTemperature": 15}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: warehouse.ErrBD.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

}

func TestDeleteWHandler(t *testing.T) {
	///Cuando la eliminación sea exitosa se devolverá un código 204
	t.Run("delete_ok", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Delete", mock.Anything, 9).Return(nil)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodDelete, "/api/v1/warehouses/9", "")
		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		assert.True(t, service.AssertExpectations(t))

	})

	//Cuando el warehouse no existe se devolverá un código 404

	t.Run("delete_non_existent", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		service.On("Delete", mock.Anything, 1).Return(warehouse.ErrNotFound)
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodDelete, "/api/v1/warehouses/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: warehouse.ErrNotFound.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))

	})
	//Cuando se evia un id invalido retorna 400 , ej un string
	t.Run("delete_id_invalid_400", func(t *testing.T) {

		service := NewServiceWarehouseTest()
		server := CreateServerWarehouses(service)

		req, resp := createRequestWarehouse(http.MethodDelete, "/api/v1/warehouses/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponseWarehouse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: warehouse.ErrBadRequest.Error(),
		}

		var result errorResponseWarehouse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))

	})
}
