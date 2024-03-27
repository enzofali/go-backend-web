package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Controller
type serviceSectionTest struct {
	mock.Mock
}

// Constructor
func NewServiceTestSection() *serviceSectionTest {
	return &serviceSectionTest{}
}

func (r *serviceSectionTest) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}

func (r *serviceSectionTest) GetByID(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *serviceSectionTest) GetReportProducts(ctx context.Context, id int) ([]domain.SectionReportProducts, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]domain.SectionReportProducts), args.Error(1)
}

func (r *serviceSectionTest) Create(ctx context.Context, section domain.Section) (domain.Section, error) {
	args := r.Called(ctx, section)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *serviceSectionTest) Update(ctx context.Context, section domain.Section) error {
	args := r.Called(ctx, section)
	return args.Error(0)
}

func (r *serviceSectionTest) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func createServerSectionUnit(service *serviceSectionTest) *gin.Engine {
	handler := NewSection(service)
	eng := gin.Default()

	sections := eng.Group("/api/v1/sections")
	{
		sections.GET("/", handler.GetAll())
		sections.GET("/:id", handler.Get())
		sections.GET("/reportProducts", handler.GetReportProducts())
		sections.POST("/", handler.Create())
		sections.PATCH("/:id", handler.Update())
		sections.DELETE("/:id", handler.Delete())
	}
	return eng
}

func createRequestSectionUnit(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

// ------------------------------- READ ---------------------------------

func Test_GetAll_Section_Unit(t *testing.T) {
	type resp struct {
		Data []domain.Section
	}

	sections := []domain.Section{
		{ID: 1, SectionNumber: 1, CurrentTemperature: 15, MinimumTemperature: -20, CurrentCapacity: 20, MinimumCapacity: 5, MaximumCapacity: 50, WarehouseID: 1, ProductTypeID: 1},
		{ID: 2, SectionNumber: 3, CurrentTemperature: 25, MinimumTemperature: -10, CurrentCapacity: 10, MinimumCapacity: 2, MaximumCapacity: 20, WarehouseID: 1, ProductTypeID: 1},
	}

	expected := resp{
		Data: sections,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetAll", mock.Anything).Return(sections, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/", "")
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("GetAll: ErrInternal", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetAll", mock.Anything).Return([]domain.Section{}, ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/", "")
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Get_Section_Unit(t *testing.T) {
	type resp struct {
		Data domain.Section
	}

	sect := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 15,
		MinimumTemperature: -20,
		CurrentCapacity:    20,
		MinimumCapacity:    5,
		MaximumCapacity:    50,
		WarehouseID:        1,
		ProductTypeID:      1}

	expected := resp{
		Data: sect,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 9999
		service.On("GetByID", mock.Anything, id).Return(domain.Section{}, section.ErrSectionNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/9999", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(domain.Section{}, ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_GetReportProducts_Section_Unit(t *testing.T) {

	type resp struct {
		Data []domain.SectionReportProducts
	}

	reports := []domain.SectionReportProducts{
		{ID: 1, SectionNumber: 1, ProductCount: 150},
		{ID: 2, SectionNumber: 3, ProductCount: 250},
	}

	expected := resp{
		Data: reports,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 0
		service.On("GetReportProducts", mock.Anything, id).Return(reports, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/reportProducts", "")
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetReportProducts", mock.Anything).Return([]domain.SectionReportProducts{}, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/reportProducts?id=error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetReportProducts", mock.Anything, id).Return([]domain.SectionReportProducts{}, section.ErrSectionNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/reportProducts?id=1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetReportProducts", mock.Anything, id).Return([]domain.SectionReportProducts{}, ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodGet, "/api/v1/sections/reportProducts?id=1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Create_Section_Unit(t *testing.T) {

	type resp struct {
		Data domain.Section
	}

	sect := domain.Section{
		SectionNumber:      1234,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1}

	expected := resp{
		Data: sect,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ShouldBindJSON", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": error, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate required JSON fields", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate unique section_number", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, section.ErrExistsSectionNumber)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate WareHouse not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, section.ErrWareHouseNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Product Type not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, section.ErrProductTypeNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("Create", mock.Anything, sect).Return(domain.Section{}, ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPost, "/api/v1/sections/", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Create_Update_Unit(t *testing.T) {
	type resp struct {
		Data domain.Section
	}

	sect := domain.Section{
		ID:                 1,
		SectionNumber:      1234,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1}

	expected := resp{
		Data: sect,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(sect, nil)
		service.On("Update", mock.Anything, sect).Return(nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"id":1, "section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/error", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(domain.Section{}, ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate json.NewDecoder(ctx.Request.Body).Decode(&section)", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("GetByID", mock.Anything, id).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": error, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Update ID", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 9999
		service.On("GetByID", mock.Anything, id).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/9999", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.NotEqual(t, fmt.Sprint(sect.ID), request.RequestURI[17:])
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate required JSON fields", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetByID", mock.Anything, 1).Return(sect, nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 0, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate unique section_number", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetByID", mock.Anything, 1).Return(sect, nil)
		service.On("Update", mock.Anything, sect).Return(section.ErrExistsSectionNumber)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate WareHouse not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetByID", mock.Anything, 1).Return(sect, nil)
		service.On("Update", mock.Anything, sect).Return(section.ErrWareHouseNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Product Type not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetByID", mock.Anything, 1).Return(sect, nil)
		service.On("Update", mock.Anything, sect).Return(section.ErrProductTypeNotFound)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		service.On("GetByID", mock.Anything, 1).Return(sect, nil)
		service.On("Update", mock.Anything, sect).Return(ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodPatch, "/api/v1/sections/1", `{"section_number": 1234, "current_temperature": 10, "minimum_temperature": 5, "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}

func Test_Create_Delete_Unit(t *testing.T) {

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("Delete", mock.Anything, id).Return(nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodDelete, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ID type", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("Delete", mock.Anything, id).Return(nil)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodDelete, "/api/v1/sections/error", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestSection()
		id := 1
		service.On("Delete", mock.Anything, id).Return(ErrInternal)
		server := createServerSectionUnit(service)

		// act
		request, response := createRequestSectionUnit(http.MethodDelete, "/api/v1/sections/1", "")
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})
}
