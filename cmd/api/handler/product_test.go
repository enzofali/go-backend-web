package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product"
	"github.com/stretchr/testify/assert"
)

// stubProductService is a double of the product package's Service for the purpose of testing the handler
type stubProductService struct {
	Products  []domain.Product
	Product   domain.Product
	ID        int
	Err       error
	Valid     bool
	Reports   []domain.Report
	Updated   domain.Product
	ErrUpdate error
}

// implementing the product.Service interface
func (s stubProductService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.Products, s.Err
}
func (s stubProductService) GetByID(ctx context.Context, id int) (domain.Product, error) {
	return s.Product, s.Err
}
func (s stubProductService) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	return s.Product, s.Err
}
func (s stubProductService) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	return s.Updated, s.ErrUpdate
}
func (s stubProductService) Delete(ctx context.Context, id int) error {
	return s.Err
}
func (s stubProductService) ValidateProductID(ctx context.Context, pid int) bool {
	return s.Valid
}
func (s stubProductService) GetOneReport(ctx context.Context, id int) ([]domain.Report, error) {
	return s.Reports, s.Err
}
func (s stubProductService) GetAllReports(ctx context.Context) ([]domain.Report, error) {
	return s.Reports, s.Err
}
func (s stubProductService) CreateType(ctx context.Context, name string) (int, error) {
	return s.ID, s.Err
}

// example product
var exampleProduct = domain.Product{
	ID:             1,
	Description:    "pepe",
	ExpirationRate: 10,
	FreezingRate:   12,
	Height:         11.2,
	Length:         10.4,
	Netweight:      120,
	ProductCode:    "Unique",
	RecomFreezTemp: -5,
	Width:          10.3,
	ProductTypeID:  1,
	SellerID:       1,
}

// testing utility functions
func createTestProductHandler(stub stubProductService) Product {
	return Product{productService: stub}
}
func createTestGinContextAndRecorder(method string) (*httptest.ResponseRecorder, *gin.Context) {
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = &http.Request{Method: method, Header: make(http.Header)}
	return rr, c
}
func mockRequestBody(c *gin.Context, contents interface{}) {
	body, _ := json.Marshal(contents)
	c.Request.Header.Set("Content-Type", "application/json")
	// converting []byte body into an io.ReadCloser
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}

// tests

func TestProductGetAll_Ok(t *testing.T) {
	// Arrange
	expected := []domain.Product{
		{ID: 1, Description: "pepe", ProductCode: "Unique"},
		{ID: 2, Description: "other", ProductCode: "Also"},
	}

	rr, c := createTestGinContextAndRecorder("GET")

	handler := createTestProductHandler(stubProductService{
		Products: expected,
		Err:      nil,
	})

	// Act
	handler.GetAll()(c)

	var res struct {
		Data []domain.Product
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, res.Data)
}

func TestProductGetAll_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Message: ErrInternal.Error(),
		Code:    "internal_server_error",
	}

	rr, c := createTestGinContextAndRecorder("GET")

	handler := createTestProductHandler(stubProductService{
		Err: errors.New("test service error"),
	})

	// Act
	handler.GetAll()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGet_Ok(t *testing.T) {
	// Arrange
	rr, c := createTestGinContextAndRecorder("GET")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product: exampleProduct,
		Err:     nil,
	})

	// Act
	handler.Get()(c)

	var res struct {
		Data domain.Product
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, exampleProduct, res.Data)
}

func TestProductGet_BadRequest(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: ErrInvalidId.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Params = []gin.Param{{Key: "id", Value: "invalid id"}}

	handler := createTestProductHandler(stubProductService{
		Product: exampleProduct,
		Err:     nil,
	})

	// Act
	handler.Get()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGet_NotFound(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "not_found",
		Message: product.ErrNotFound.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Err: product.ErrNotFound,
	})

	// Act
	handler.Get()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductCreate_Created(t *testing.T) {
	// Arrange
	expected := domain.Product{ID: 1, Description: "different"}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, exampleProduct)

	handler := createTestProductHandler(stubProductService{
		Product: expected,
		Err:     nil,
	})

	// Act
	handler.Create()(c)

	var res struct {
		Data domain.Product
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expected, res.Data)
}

func TestProductCreate_UnprocessableEntity(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "unprocessable_entity",
		Message: ErrField.Error(),
	}
	unprocessable := struct{ Description int }{Description: 11}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, unprocessable)

	handler := createTestProductHandler(stubProductService{})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductCreate_BadRequest(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: product.ErrExists.Error(),
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, exampleProduct)

	handler := createTestProductHandler(stubProductService{Err: product.ErrExists})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductCreate_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, exampleProduct)

	handler := createTestProductHandler(stubProductService{Err: product.ErrDatabase})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_Ok(t *testing.T) {
	// Arrange
	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{ID: 1})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product:   domain.Product{},
		Updated:   exampleProduct,
		Err:       nil,
		ErrUpdate: nil,
	})

	// Act
	handler.Update()(c)

	var updated struct {
		Data domain.Product
	}
	err := json.Unmarshal(rr.Body.Bytes(), &updated)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, exampleProduct, updated.Data)
}

func TestProductUpdate_BadRequestInvalidID(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "not_found",
		Message: product.ErrNotFound.Error(),
	}

	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Err: product.ErrNotFound,
	})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_NotFound(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: ErrInvalidId.Error(),
	}

	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{})
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	handler := createTestProductHandler(stubProductService{})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_BadRequestCannotUpdateID(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: "cannot update product id",
	}
	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{ID: 2})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product: domain.Product{},
		Err:     nil,
	})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_BadRequestCannotDecode(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: ErrBadRequest.Error(),
	}
	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, struct{ Description int }{Description: 1})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product: domain.Product{},
		Err:     nil,
	})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_BadRequestCodeAlreadyExists(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: product.ErrExists.Error(),
	}
	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{ID: 1})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product:   domain.Product{},
		Err:       nil,
		Updated:   domain.Product{},
		ErrUpdate: product.ErrExists,
	})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductUpdate_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}
	rr, c := createTestGinContextAndRecorder("PATCH")
	mockRequestBody(c, domain.Product{ID: 1})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Product:   domain.Product{},
		Err:       nil,
		Updated:   domain.Product{},
		ErrUpdate: product.ErrDatabase,
	})

	// Act
	handler.Update()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductDelete_NoContent(t *testing.T) {
	// Arrange
	rr, c := createTestGinContextAndRecorder("DELETE")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Err: nil,
	})

	// Act
	handler.Delete()(c)

	// Assert

	assert.Equal(t, http.StatusNoContent, rr.Code)

}

func TestProductDelete_BadRequestInvalidID(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: ErrInvalidId.Error(),
	}

	rr, c := createTestGinContextAndRecorder("DELETE")
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	handler := createTestProductHandler(stubProductService{
		Err: nil,
	})

	// Act
	handler.Delete()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductDelete_NotFound(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "not_found",
		Message: product.ErrNotFound.Error(),
	}

	rr, c := createTestGinContextAndRecorder("DELETE")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Err: product.ErrNotFound,
	})

	// Act
	handler.Delete()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductDelete_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("DELETE")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := createTestProductHandler(stubProductService{
		Err: product.ErrDatabase,
	})

	// Act
	handler.Delete()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductCreateType_Created(t *testing.T) {
	// Arrange
	expected := domain.ProductType{
		ID:   1,
		Name: "pepe",
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, domain.ProductTypeRequest{Name: "pepe"})

	handler := createTestProductHandler(stubProductService{
		ID:  1,
		Err: nil,
	})

	// Act
	handler.CreateType()(c)

	var res struct {
		Data domain.ProductType
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expected, res.Data)
}

func TestProductCreateType_UnprocessableEntity(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "unprocessable_entity",
		Message: ErrField.Error(),
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, struct{ Name int }{Name: 1})

	handler := createTestProductHandler(stubProductService{})

	// Act
	handler.CreateType()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductCreateType_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, domain.ProductTypeRequest{Name: "pepe"})

	handler := createTestProductHandler(stubProductService{
		ID:  1,
		Err: product.ErrDatabase,
	})

	// Act
	handler.CreateType()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGetAllReports_Ok(t *testing.T) {
	// Arrange
	expected := []domain.Report{
		{ProductID: 1, Description: "pepe", Count: 1},
		{ProductID: 2, Description: "sanchez", Count: 0},
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL = &url.URL{}

	handler := createTestProductHandler(stubProductService{
		Reports: expected,
		Err:     nil,
	})

	// Act
	handler.GetReport()(c)

	var res struct {
		Data []domain.Report
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, res.Data)
}

func TestProductGetAllReports_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL = &url.URL{}

	handler := createTestProductHandler(stubProductService{
		Err: product.ErrDatabase,
	})

	// Act
	handler.GetReport()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGetOneReport_Ok(t *testing.T) {
	// Arrange
	expected := []domain.Report{
		{ProductID: 1, Description: "pepe", Count: 1},
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL, _ = url.Parse("?id=1")

	handler := createTestProductHandler(stubProductService{
		Reports: expected,
		Err:     nil,
		Valid:   true,
	})

	// Act
	handler.GetReport()(c)

	var res struct {
		Data []domain.Report
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, res.Data)
}

func TestProductGetOneReport_BadRequestInvalidID(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "bad_request",
		Message: ErrInvalidId.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL, _ = url.Parse("?id=invalid")

	handler := createTestProductHandler(stubProductService{
		Reports: []domain.Report{},
		Err:     nil,
		Valid:   true,
	})

	// Act
	handler.GetReport()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGetOneReport_NotFound(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "not_found",
		Message: ErrNotFound.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL, _ = url.Parse("?id=1")

	handler := createTestProductHandler(stubProductService{
		Reports: []domain.Report{},
		Err:     nil,
		Valid:   false,
	})

	// Act
	handler.GetReport()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestProductGetOneReport_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("GET")
	c.Request.URL, _ = url.Parse("?id=1")

	handler := createTestProductHandler(stubProductService{
		Reports: []domain.Report{},
		Err:     product.ErrDatabase,
		Valid:   true,
	})

	// Act
	handler.GetReport()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}
