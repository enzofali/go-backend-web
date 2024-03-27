package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_records"
	"github.com/stretchr/testify/assert"
)

// stubPRService is a double of the product package's Service for the purpose of testing the handler
type stubPRService struct {
	PR    domain.ProductRecord
	Valid bool
	Err   error
}

// implementing the product_records.Service interface
func (s stubPRService) Create(ctx context.Context, pr domain.ProductRecord) (domain.ProductRecord, error) {
	return s.PR, s.Err
}
func (s stubPRService) ValidateProductID(ctx context.Context, id int) bool {
	return s.Valid
}

// testing utility functions
func createTestPRHandler(stub stubPRService) ProductRecords {
	return ProductRecords{productRecordsService: stub}
}

// example ProductRecord
var examplePR = domain.ProductRecord{
	ID:             1,
	LastUpdateDate: "12/12/2011",
	PurchasePrice:  123.123,
	SalePrice:      200.0,
	ProductID:      1,
}

// tests

func TestPRCreate_Created(t *testing.T) {
	// Arrange
	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, examplePR)

	handler := createTestPRHandler(stubPRService{
		PR:    examplePR,
		Valid: true,
		Err:   nil,
	})

	// Act
	handler.Create()(c)

	var res struct {
		Data domain.ProductRecord
	}
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, examplePR, res.Data)
}

func TestPRCreate_UnprocessableEntity(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "unprocessable_entity",
		Message: "missing or incorrect fields",
	}
	unprocessable := struct{}{}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, unprocessable)
	fmt.Println(c.Request.Body)

	handler := createTestPRHandler(stubPRService{
		PR:    examplePR,
		Valid: true,
		Err:   nil,
	})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestPRCreate_Conflict(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "conflict",
		Message: "product id does not exist",
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, examplePR)

	handler := createTestPRHandler(stubPRService{
		PR:    examplePR,
		Valid: false,
		Err:   nil,
	})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, expectedRes, res)
}

func TestPRCreate_InternalServerError(t *testing.T) {
	// Arrange
	expectedRes := errorResponse{
		Code:    "internal_server_error",
		Message: ErrInternal.Error(),
	}

	rr, c := createTestGinContextAndRecorder("POST")
	mockRequestBody(c, examplePR)

	handler := createTestPRHandler(stubPRService{
		PR:    domain.ProductRecord{},
		Valid: true,
		Err:   product_records.ErrDatabase,
	})

	// Act
	handler.Create()(c)

	var res errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedRes, res)
}
