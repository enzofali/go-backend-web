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
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/employee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type serviceEmployeeMock struct {
	mock.Mock
}

func NewServiceEmployeeMock() *serviceEmployeeMock {
	return &serviceEmployeeMock{}
}

func (sm *serviceEmployeeMock) GetAll(ctx context.Context) ([]domain.Employee, error) {
	args := sm.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (sm *serviceEmployeeMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := sm.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (sm *serviceEmployeeMock) Create(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	args := sm.Called(ctx, e)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (sm *serviceEmployeeMock) Update(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	args := sm.Called(ctx, e)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (sm *serviceEmployeeMock) Delete(ctx context.Context, id int) error {
	args := sm.Called(ctx, id)
	return args.Error(0)
}

func (sm *serviceEmployeeMock) GetAllInoundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	args := sm.Called(ctx)
	return args.Get(0).([]domain.EmployeeWithInboundOrders), args.Error(1)
}

func (sm *serviceEmployeeMock) GetWithInboundOrder(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	args := sm.Called(ctx, id)
	return args.Get(0).(domain.EmployeeWithInboundOrders), args.Error(1)
}

func createServerEmployeeUnit(service *serviceEmployeeMock) *gin.Engine {
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

func createRequestEmployeeUnit(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func Test_Employee_GetAll(t *testing.T) {
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

	t.Run("GetAll Success 200", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("GetAll", mock.Anything).Return(employees, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("GetAll Error 500", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("GetAll", mock.Anything).Return([]domain.Employee{}, employee.ErrDatabase)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: employee.ErrDatabase.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func Test_Employee_Get(t *testing.T) {
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

	t.Run("Get Success 200", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employee, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Get Error 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Get Error 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(domain.Employee{}, ErrNotFound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: ErrNotFound.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func Test_Employee_Create(t *testing.T) {
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

	employeeReq := domain.Employee{
		CardNumberID: "A12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	data := response{
		Data: employeeResponse,
	}

	t.Run("Create Success 200", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Create", mock.Anything, employeeReq).Return(employeeResponse, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error ShouldBind 422", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"first_name": "Juan", "last_name": "Perez", "warehouse_i`)
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
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusUnprocessableEntity)), " ", "_"),
			Message: "CardNumberID-required,",
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Not Warehouse 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Create", mock.Anything, employeeReq).Return(domain.Employee{}, employee.ErrWarehouseNotfound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrNotWareHouse.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Exists CardID 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Create", mock.Anything, employeeReq).Return(domain.Employee{}, employee.ErrExistsCardId)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrBadRequest.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Create Error Internal 500", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Create", mock.Anything, employeeReq).Return(domain.Employee{}, errors.New("error database"))
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPost, "/api/v1/employees", `{"card_number_id": "A12", "first_name": "Juan", "last_name": "Perez", "warehouse_id": 1}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: ErrInternalServer.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func Test_Employee_Update(t *testing.T) {
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
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employeeBefore, nil)
		service.On("Update", mock.Anything, employeeAfter).Return(employeeAfter, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Error Id 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/abc", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Error Employee Not Found 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(domain.Employee{}, employee.ErrNotFound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: employee.ErrNotFound.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Error Json 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employeeBefore, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrBadRequest.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Validator 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employeeBefore, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": ""}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: "FirstName-required,",
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Error Warehouse Not Found 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employeeBefore, nil)
		service.On("Update", mock.Anything, employeeAfter).Return(domain.Employee{}, employee.ErrWarehouseNotfound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrNotWareHouse.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Update Error DB 500", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Get", mock.Anything, 1).Return(employeeBefore, nil)
		service.On("Update", mock.Anything, employeeAfter).Return(domain.Employee{}, errors.New("error database"))
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Juan Jose"}`)
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

func Test_Employee_Delete(t *testing.T) {
	t.Run("Delete OK 204", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Delete", mock.Anything, 1).Return(nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Delete Error Id 400", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/abc", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Delete Error Not Found 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Delete", mock.Anything, 1).Return(employee.ErrNotFound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: employee.ErrNotFound.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("Delete Error Database 404", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("Delete", mock.Anything, 1).Return(employee.ErrDatabase)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodDelete, "/api/v1/employees/1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusInternalServerError)), " ", "_"),
			Message: employee.ErrDatabase.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})
}

func Test_Employee_GetAllWithInboundOrders_WithoutID(t *testing.T) {
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
		service := NewServiceEmployeeMock()
		service.On("GetAllInoundOrders", mock.Anything).Return(employees, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("GetAllWithInboundOrders Without id Error", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("GetAllInoundOrders", mock.Anything).Return([]domain.EmployeeWithInboundOrders{}, employee.ErrDatabase)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
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

func Test_Employee_GetAllWithInboundOrders_WithID(t *testing.T) {
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
		service := NewServiceEmployeeMock()
		service.On("GetWithInboundOrder", mock.Anything, 1).Return(employe, nil)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=1", "")
		server.ServeHTTP(resp, req)

		var result response
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, data, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("GetAllWithInboundOrders With id Error id format", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=abd", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusBadRequest)), " ", "_"),
			Message: ErrInvalidId.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

	t.Run("GetAllWithInboundOrders With id Error Not Found", func(t *testing.T) {
		service := NewServiceEmployeeMock()
		service.On("GetWithInboundOrder", mock.Anything, 1).Return(domain.EmployeeWithInboundOrders{}, employee.ErrNotFound)
		server := createServerEmployeeUnit(service)

		req, resp := createRequestEmployeeUnit(http.MethodGet, "/api/v1/employees/reportInboundOrders?id=1", "")
		server.ServeHTTP(resp, req)

		errResp := errorResponse{
			Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(http.StatusNotFound)), " ", "_"),
			Message: ErrNotFound.Error(),
		}

		var result errorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Equal(t, errResp, result)
		assert.True(t, service.AssertExpectations(t))
	})

}
