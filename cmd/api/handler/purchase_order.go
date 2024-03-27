package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/purchaseorder"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Purchase_Order struct {
	purchOrdService purchaseorder.Service
}

func NewPurchaseOrder(s purchaseorder.Service) *Purchase_Order {
	return &Purchase_Order{
		purchOrdService: s,
	}
}

// @Summary		Create purchase order
// @Tags			Purchase Order
// @Description	Create purchase order
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Purchase_Orders	true	"buyers parameters"
// @Success		201		{object}	web.response{data=domain.Purchase_Orders}
// @Failure		409		{object}	web.errorResponse
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/purchaseorders [post]
func (PurchOrder *Purchase_Order) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Request
		var NewPurchaseOrder domain.Purchase_Orders
		// Validate the JSON
		err := ctx.ShouldBind(&NewPurchaseOrder)

		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, purchaseorder.ErrFieldNotExist.Error())
			return
		}

		validator := validator.New()
		if err := validator.Struct(&NewPurchaseOrder); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, purchaseorder.ErrFieldNotExist.Error())
			return
		}
		//Validate product_record_id if exists

		//then create the purchase order
		purchOrder, err := PurchOrder.purchOrdService.Create(ctx, NewPurchaseOrder)
		//control possible errors while creating purchase order
		if err != nil {
			switch err {
			case purchaseorder.ErrBuyerNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case purchaseorder.ErrDatabase:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
		web.Success(ctx, http.StatusCreated, purchOrder)
	}
}
