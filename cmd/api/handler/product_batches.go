package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_batches"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type ProductBatches struct {
	s product_batches.Service
}

func NewProductBatches(s product_batches.Service) *ProductBatches {
	return &ProductBatches{
		s: s,
	}
}

// -------------------------------- POST Methods --------------------------------

// @Summary		Create Product Batch
// @Tags			Product Batches
// @Description	Create Product Batch
// @Accept			json
// @Produce		json
// @Param			section	body		domain.ProductBatches	true	"Product Batch to Create"
// @Success		201		{object}	web.response
// @Failure		422		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/productBatches [post]
func (s *ProductBatches) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request domain.ProductBatches

		// Bind JSON to domain.ProductBatches{}
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Validate missing JSON key:values
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Validate Date and Time fields
		DueDate, err := time.Parse("2006-01-02", request.DueDate)
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		request.DueDate = DueDate.Format("2006-01-02")

		ManufacturingDate, err := time.Parse("2006-01-02", request.ManufacturingDate)
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		request.ManufacturingDate = ManufacturingDate.Format("2006-01-02")

		ManufacturingHour, err := time.Parse("15:04:05", request.ManufacturingHour)
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		request.ManufacturingHour = ManufacturingHour.Format("15:04:05")

		// Validate unique batch_number: If the  batch_number already exists, return a 409 Conflict error
		productBatches, err := s.s.Create(ctx, request)
		if err != nil {
			switch err {
			case product_batches.ErrExistsBatchNumber:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case product_batches.ErrProductNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case product_batches.ErrSectionNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Response
		// When the data entry is successful, a 201 code will be returned along with the entered object
		web.Success(ctx, http.StatusCreated, productBatches)
	}
}
