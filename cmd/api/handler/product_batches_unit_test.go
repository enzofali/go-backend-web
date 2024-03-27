package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_batches"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Controller
type serviceProductBatchesTest struct {
	mock.Mock
}

// Constructor
func NewServiceTestProductBatches() *serviceProductBatchesTest {
	return &serviceProductBatchesTest{}
}

func (r *serviceProductBatchesTest) Create(ctx context.Context, productBatches domain.ProductBatches) (domain.ProductBatches, error) {
	args := r.Called(ctx, productBatches)
	return args.Get(0).(domain.ProductBatches), args.Error(1)
}

func createServerProductBatchesUnit(service *serviceProductBatchesTest) *gin.Engine {
	handler := NewProductBatches(service)
	eng := gin.Default()
	productBatches := eng.Group("/api/v1/productBatches")
	{
		productBatches.POST("/", handler.Create())
	}
	return eng
}

func createRequestProductBatchesUnit(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	return request, httptest.NewRecorder()
}

func Test_Create_Product_Batches_Unit(t *testing.T) {

	type resp struct {
		Data domain.ProductBatches
	}

	productBatch := domain.ProductBatches{
		BatchNumber:        1234,
		CurrentQuantity:    10,
		CurrentTemperature: 10,
		DueDate:            "2023-02-01",
		InitialQuantity:    5,
		ManufacturingDate:  "2023-01-01",
		ManufacturingHour:  "13:01:06",
		MinumumTemperature: 5,
		ProductID:          1,
		SectionID:          1,
	}

	expected := resp{
		Data: productBatch,
	}

	t.Run("Ok", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(productBatch, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1234, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		var result resp
		err := json.Unmarshal(response.Body.Bytes(), &result)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expected, result)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate ShouldBindJSON", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": error, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate required JSON fields", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: DueDate", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-92-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: ManufacturingDate", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-91-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Date and Time fields: ManufacturingHour", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, nil)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 3, "current_quantity": 50, "current_temperature": 15, "due_date": "2023-02-01", "initial_quantity": 50, "manufacturing_date": "2023-01-01", "manufacturing_hour": "93:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 3}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate unique batch_number", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, product_batches.ErrExistsBatchNumber)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1234, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Product not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, product_batches.ErrProductNotFound)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1234, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Section not found", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, product_batches.ErrSectionNotFound)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1234, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

	t.Run("Validate Default error", func(t *testing.T) {
		// arrange
		service := NewServiceTestProductBatches()
		service.On("Create", mock.Anything, productBatch).Return(domain.ProductBatches{}, ErrInternal)
		server := createServerProductBatchesUnit(service)

		// act
		request, response := createRequestProductBatchesUnit(http.MethodPost, "/api/v1/productBatches/", `{"batch_number": 1234, "current_quantity": 10, "current_temperature": 10, "due_date": "2023-02-01", "initial_quantity": 5, "manufacturing_date": "2023-01-01", "manufacturing_hour": "13:01:06", "minumum_temperature": 5, "product_id": 1, "section_id": 1}`)
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json; charset=utf-8")
	})

}
