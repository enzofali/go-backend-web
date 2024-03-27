package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_records"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type ProductRecords struct {
	productRecordsService product_records.Service
}

func NewProductRecords(p product_records.Service) *ProductRecords {
	return &ProductRecords{productRecordsService: p}
}

var (
	ErrField = errors.New("missing or incorrect fields")
)

// @summary		Create product record
// @tags			Product Records
// @Description	Creates and returns a single product record
// @Accept			json
// @Produce		json
// @Param			request	body		domain.ProductRecordRequest	true	"Product Record parameters"
// @Success		201		{object}	web.response{data=domain.ProductRecord}
// @Failure		409		{object}	web.errorResponse
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/productRecords/ [post]
func (pr *ProductRecords) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productRecordToCreate domain.ProductRecord

		// bind json object to ProductRecord instance, return with error status 422 if body is malformed
		err := c.ShouldBindJSON(&productRecordToCreate)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}

		validator := validator.New()
		if err := validator.Struct(&productRecordToCreate); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}

		// validates that the given product_id associated with the product record corresponds
		// to an existing product in the database, returns with error status 404 otherwise
		if !pr.productRecordsService.ValidateProductID(context.Background(), productRecordToCreate.ProductID) {
			web.Error(c, http.StatusConflict, "product id does not exist")
			return
		}

		// attempts to create the product in the database, returns with error status 500 if it fails
		productRecord, err := pr.productRecordsService.Create(context.Background(), productRecordToCreate)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}

		web.Success(c, http.StatusCreated, productRecord)
	}
}
