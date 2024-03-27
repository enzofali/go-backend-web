package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/buyer"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Buyer struct {
	// buyerService buyer.Service
	s buyer.Service
}

func NewBuyer(s buyer.Service) *Buyer {
	return &Buyer{
		s: s,
	}
}

// @Summary		Buyer by id
// @Tags			Buyers
// @Description	get buyer by id
// @Produce		json
// @Param			id	path		int	true	"buyer id"
// @Success		200	{object}	web.response{data=domain.Buyer}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/buyers/{id} [get]
func (b *Buyer) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, buyer.ErrFormat.Error())
			return
		}

		// Process
		// When the buyer does not exist, a 404 code will be returned
		getbuyer, err := b.s.Get(c, id)
		if err != nil {
			switch err {
			case buyer.ErrNotFound:
				web.Error(c, http.StatusNotFound, err.Error())
				return
			case buyer.ErrDatabase:
				web.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		//When the request is successful, the backend will return the buyer with that id
		web.Success(c, http.StatusOK, getbuyer)
	}
}

// @summary		List buyers
// @tags			Buyers
// @Description	Returns a list of all buyers
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Buyer}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/buyers [get]
func (b *Buyer) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		buyers, err := b.s.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, buyer.ErrDatabase.Error())
			return
		}
		// Response
		// When the request is successful, the backend will return a list of all existing sellers
		web.Success(c, http.StatusOK, buyers)
	}
}

// @Summary		Create buyer
// @Tags			Buyers
// @Description	Create buyer
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Buyer	true	"buyers parameters"
// @Success		201		{object}	web.response{data=domain.Buyer}
// @Failure		422		{object}	web.errorResponse
// @Failure		400		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/buyers [post]
func (b *Buyer) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		var buyerRequest domain.BuyerRequest
		err := c.ShouldBind(&buyerRequest)
		if err != nil {
			web.Error(c, http.StatusBadRequest, buyer.ErrFormat.Error())
			return
		}

		// Validate the JSON
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validateBuyer := validator.New()
		err = validateBuyer.Struct(&buyerRequest)

		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, buyer.ErrFormat.Error())
			return
		}

		// Process
		newBuyer := domain.Buyer{
			CardNumberID: buyerRequest.CardNumberID,
			FirstName:    buyerRequest.FirstName,
			LastName:     buyerRequest.LastName,
		}

		// Validate unique CNID: If the CNID already exists, return a 409 Conflict error
		newBuyer, err = b.s.Create(c, newBuyer)

		switch err {
		case buyer.ErrAlreadyExists:
			web.Error(c, http.StatusConflict, err.Error())
			return
		case buyer.ErrDatabase:
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		//newBuyer.ID = id

		// Response
		// When the data entry is successful, a 201 code will be returned along with the entered object
		web.Success(c, http.StatusCreated, newBuyer)
	}
}

// @Summary		Update buyer
// @Tags			Buyers
// @Description	Update buyer
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"buyer id"
// @Param			request	body		domain.Buyer	true	"buyer parameters"
// @Success		200		{object}	web.response{data=domain.Buyer}
// @Failure		404		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/buyers/{id} [patch]
func (b *Buyer) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, buyer.ErrFormat.Error())
			return
		}

		// If the buyer to be updated does not exist, a 404 code will be returned
		buyerDB, err := b.s.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, buyer.ErrNotFound.Error())
			return
		}

		// decode and update fetched product object with fields decoded from request body
		err = json.NewDecoder(c.Request.Body).Decode(&buyerDB)
		if err != nil {
			web.Error(c, http.StatusConflict, buyer.ErrCantChange.Error())
			return
		}

		//new id should not be specified in request body

		if buyerDB.ID != id {
			web.Error(c, http.StatusConflict, buyer.ErrCantChange.Error())
			return
		}

		updateBuyer := domain.BuyerRequest{
			CardNumberID: buyerDB.CardNumberID,
			FirstName:    buyerDB.FirstName,
			LastName:     buyerDB.LastName,
		}

		//if buyer exists, the user updates it and has an error
		validJson := validator.New()
		err = validJson.Struct(&updateBuyer)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, buyer.ErrCantChange.Error())
			return
		}

		buyerDB = domain.Buyer{
			ID:           id,
			CardNumberID: buyerDB.CardNumberID,
			FirstName:    updateBuyer.FirstName,
			LastName:     updateBuyer.LastName,
		}

		err = b.s.Update(c, buyerDB, id)
		if err != nil {
			switch err {
			case buyer.ErrCantChange:
				web.Error(c, http.StatusConflict, err.Error())
				return
			case buyer.ErrDatabase:
				web.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		//buyer actualize OK
		web.Success(c, http.StatusOK, buyerDB)
	}
}

// @Summary		Delete buyer
// @Tags			Buyers
// @Description	Delete buyer
// @Param			id	path		int	true	"buyer id"
// @Success		204	{object}	web.response
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/buyers/{id} [delete]
func (b *Buyer) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, buyer.ErrFormat.Error())
			return
		}
		err = b.s.Delete(c, id)

		// Process
		// When the buyer does not exist a 404 code will be returned
		if err != nil {
			switch err {
			case buyer.ErrNotFound:
				web.Error(c, http.StatusNotFound, err.Error())
				return
			case buyer.ErrDatabase:
				web.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		// Response
		// When the deletion is successful, a 204 code will be returned.
		web.Success(c, http.StatusNoContent, nil)
	}
}

// @Summary		Purchase orders by buyer and all
// @Tags			Buyers
// @Description	get report by id or all buyers
// @Produce		json
// @Param			id	query	int	false	"Buyer ID"
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Buyer}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/buyers/reportPurchaseOrders [get]
func (b *Buyer) GetReport() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get and corroborate query id

		var id int
		var err error

		idquery, ok := ctx.GetQuery("id")
		if ok {
			id, err = strconv.Atoi(idquery)
			if err != nil {
				web.Error(ctx, http.StatusBadRequest, buyer.ErrFormat.Error())
				return
			}
		}
		//create report by buyer id. If buyer id not exist or a database error happens, return error.
		//if id==0, return all buyers with they purchase orders
		report, err := b.s.GetReports(ctx, id)
		if err != nil {
			switch err {
			case buyer.ErrPurchaseNotFound:
				web.Error(ctx, http.StatusNotFound, err.Error())
				return
			case buyer.ErrPurchasesNotFound:
				web.Error(ctx, http.StatusNotFound, err.Error())
				return
			case buyer.ErrDatabase:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
		//when the request is successful, the backend will return the report
		web.Success(ctx, http.StatusOK, report)
	}
}
