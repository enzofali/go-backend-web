package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/seller"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Seller struct {
	sellerService seller.Service
}

func NewSeller(s seller.Service) *Seller {
	return &Seller{
		sellerService: s,
	}
}

// @summary		List sellers
// @tags			Sellers
// @Description	Returns a list of all sellers
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Seller}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/sellers [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process
		sellers, err := s.sellerService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		// Response
		// When the request is successful, the backend will return a list of all existing sellers
		web.Success(c, http.StatusOK, sellers)
	}
}

// @Summary		Seller by id
// @Tags			Sellers
// @Description	get seller by id
// @Produce		json
// @Param			id	path		int	true	"seller id"
// @Success		200	{object}	web.response{data=domain.Seller}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Router			/api/v1/sellers/{id} [get]
func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Process
		// When the seller does not exist, a 404 code will be returned
		sel, err := s.sellerService.GetByID(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		//When the request is successful, the backend will return the vendor with that id
		web.Success(c, http.StatusOK, sel)
	}
}

// @Summary		Create seller
// @Tags			Sellers
// @Description	Create seller
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Seller	true	"Sellers parameters"
// @Success		201		{object}	web.response{data=domain.Seller}
// @Failure		422		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/sellers [post]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		var request domain.Seller
		if err := c.ShouldBind(&request); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "error bad request")
			return
		}

		// Validate the JSON
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validator := validator.New()
		if err := validator.Struct(&request); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if request.CID < 0 {
			web.Error(c, http.StatusUnprocessableEntity, "invalid cid")
			return
		}

		// Validate unique CID: If the CID already exists, return a 409 Conflict error
		id, err := s.sellerService.Create(c, request)
		switch err {
		case seller.ErrConflict:
			web.Error(c, http.StatusConflict, err.Error())
			return
		case seller.ErrIntern:
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		case seller.ErrInvalidLocality:
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		request.ID = id

		// Response
		// When the data entry is successful, a 201 code will be returned along with the entered object
		web.Success(c, http.StatusCreated, request)
	}
}

// @Summary		Update seller
// @Tags			Sellers
// @Description	Update seller
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"section id"
// @Param			request	body		domain.Seller	true	"Seller parameters"
// @Success		200		{object}	web.response{data=domain.Seller}
// @Failure		400		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Failure		422		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/sellers/{id} [patch]
func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// If the seller to be updated does not exist, a 404 code will be returned
		sellerDB, err := s.sellerService.GetByID(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		// decode and update fetched product object with fields decoded from request body
		err = json.NewDecoder(c.Request.Body).Decode(&sellerDB)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "bad request")
			return
		}
		// new id should not be specified in request body, i.e. it should not change
		if sellerDB.ID != id {
			web.Error(c, http.StatusBadRequest, "cannot update product id")
			return
		}

		// Validate the JSON
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validator := validator.New()
		if err := validator.Struct(&sellerDB); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if sellerDB.CID < 0 {
			web.Error(c, http.StatusUnprocessableEntity, "invalid cid")
			return
		}
		// Validate unique CID: If the CID already exists, return a 409 Conflict error
		err = s.sellerService.Update(c, sellerDB)
		switch err {
		case seller.ErrConflict:
			web.Error(c, http.StatusConflict, err.Error())
			return
		case seller.ErrIntern:
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		// When the data update is successful, the section with the updated information will be returned along with a code 200
		web.Success(c, http.StatusOK, sellerDB)
	}
}

// @Summary		Delete seller
// @Tags			Sellers
// @Description	Delete seller
// @Param			id	path		int	true	"seller id"
// @Success		204	{object}	web.response
// @Failure		404	{object}	web.errorResponse
// @Failure		400	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/sellers/{id} [delete]
func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Process
		// When the section does not exist a 404 code will be returned
		err = s.sellerService.Delete(c, id)
		switch err {
		case seller.ErrNotFound:
			web.Error(c, http.StatusNotFound, err.Error())
			return
		case seller.ErrIntern:
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		// When the deletion is successful, a 204 code will be returned.
		web.Success(c, http.StatusNoContent, gin.H{})
	}
}
