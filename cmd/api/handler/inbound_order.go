package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	inboundorder "github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/inbound_order"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

var (
	ErrEmployeeNotFound     = errors.New("Employee not found")
	ErrProductBatchNotFound = errors.New("Product Batch not found")
	ErrWarehouseNotFound    = errors.New("Warehouse not found")
	ErrOrderNumberExtists   = errors.New("Order number exists")
)

type InboudOrder struct {
	service inboundorder.Service
}

func NewInoudOrder(i inboundorder.Service) *InboudOrder {
	return &InboudOrder{
		service: i,
	}
}

// @summary		Create inbound order
// @tags			Inbound Order
// @Description	create inbound order
// @Accept			json
// @Param			request	body	domain.InboundOrderRequest	true	"query params"
// @Produce		json
// @Success		201	{object}	web.response{data=domain.InboundOrder}
// @Failure		422	{object}	web.errorResponse
// @Failure		400	{object}	web.errorResponse
// @Failure		409	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/inboundOrders [post]
func (i *InboudOrder) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var inboundOrderRequest domain.InboundOrderRequest
		err := c.ShouldBind(&inboundOrderRequest)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrBadRequest.Error())
			return
		}

		validate := validator.New()
		err = validate.Struct(inboundOrderRequest)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			msgFields := ""
			for _, ve := range validateErr {
				msgFields += ve.Field() + "-" + ve.Tag() + ","
			}
			web.Error(c, http.StatusUnprocessableEntity, msgFields)
			return
		}

		date, err := time.Parse("2006-01-02", inboundOrderRequest.OrderDate)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "Date format incorrect")
			return
		}

		inboudOrderSave := domain.InboundOrder{
			OrderDate:      date.Format("2006-01-02"),
			OrderNumber:    inboundOrderRequest.OrderNumber,
			EmployeeID:     inboundOrderRequest.EmployeeID,
			ProductBatchID: inboundOrderRequest.ProductBatchID,
			WarehouseID:    inboundOrderRequest.WarehouseID,
		}

		inboudOrdeDB, err := i.service.Create(c, inboudOrderSave)
		if err != nil {
			switch err {
			case inboundorder.ErrEmployeeNotFound:
				web.Error(c, http.StatusConflict, ErrEmployeeNotFound.Error())
				return
			case inboundorder.ErrProductBatchNotFound:
				web.Error(c, http.StatusConflict, ErrProductBatchNotFound.Error())
				return
			case inboundorder.ErrWarehouseNotFound:
				web.Error(c, http.StatusConflict, ErrWarehouseNotFound.Error())
				return
			case inboundorder.ErrOrderNumberExtists:
				web.Error(c, http.StatusConflict, ErrOrderNumberExtists.Error())
				return
			default:
				web.Error(c, http.StatusInternalServerError, ErrInternalServer.Error())
				return
			}
		}

		web.Success(c, http.StatusCreated, inboudOrdeDB)
	}
}
