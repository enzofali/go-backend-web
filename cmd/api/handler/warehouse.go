package handler

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/warehouse"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Warehouse struct {
	s warehouse.Service
}

func NewWarehouse(s warehouse.Service) *Warehouse {
	return &Warehouse{
		s: s,
	}
}

// @summary		list warehouse
// @tags			Warehouse
// @Description	Returns a list of all warehouse
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Warehouse}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/warehouses/ [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get and return all warehouse
		wareH, err := w.s.GetAll(c)
		if err != nil {
			web.Error(c, 500, err.Error())
			return

		}
		web.Success(c, 200, wareH)
	}
}

// @summary		warehouse
// @tags			Warehouse
// @Description	Get warehouse by id
// @Produce		json
// @Param			id	path		int	true	"Werehouse Id"
// @Success		200	{object}	web.response{data=domain.Warehouse}
// @Failure		400	{object}	web.errorResponse
// @failure		404	{object}	web.errorResponse
// @Router			/api/v1/warehouses/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get id param fro  URL , and conver to integer
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)

		if err != nil {
			web.Error(c, 400, ErrBadRequest.Error())
			return
		}
		//get warehouse defined by id param, return 404 if warehouse doesn't exist
		found, err := w.s.Get(c, id)

		if err != nil {
			web.Error(c, 404, err.Error())
			return
		}
		web.Success(c, 200, found)

	}
}

// @summary		Create warehouse
// @tags			Warehouse
// @Description	Create warehouse
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Warehouse	true	"query params"
// @Success		201		{object}	web.response{data=domain.Warehouse}
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Failure		400		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Router			/api/v1/warehouses [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Validate format JSON, and return 422 if warehouse is
		//improperly constructed
		var wareHRequest domain.Warehouse
		err := c.ShouldBindJSON(&wareHRequest)

		if err != nil {
			web.Error(c, 422, err.Error())
			return
		}

		//validate field empty
		validate := validator.New()
		err = validate.Struct(wareHRequest)

		if err != nil {
			web.Error(c, 400, ErrBadRequest.Error())
			return
		}

		//Create warehouse and Validate type errors

		ware, err := w.s.Create(c, wareHRequest)

		if err == warehouse.ErrExist {
			web.Error(c, 409, err.Error())
			return

		}

		if err == warehouse.ErrBD {
			web.Error(c, 500, err.Error())
			return

		}

		web.Success(c, 201, ware)

	}
}

// @summary		Update warehouse
// @tags			Warehouse
// @Description	update warehouse
// @Accept			json
// @Produce		json
// @Param			request	body		domain.WarehouseRequest	true	"query params"
// @Param			id		path		int						true	"Warehouse Id"
// @Success		200		{object}	web.response{data=domain.Warehouse}
// @Failure		400		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Router			/api/v1/warehouses/{id} [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id param from URL, must be integer
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			web.Error(c, 400, ErrBadRequest.Error())
			return
		}

		// get warehouse by id
		wareH, err := w.s.Get(c, id)

		if err != nil {
			web.Error(c, 404, err.Error())
			return
		}

		//get the body of what I want to update
		// indicate in which structure I want to save it
		//Validate that it comes in correct JSON format
		body := c.Request.Body
		err = json.NewDecoder(body).Decode(&wareH)

		if err != nil {
			web.Error(c, 422, err.Error())
			return
		}

		//validate empty
		validate := validator.New()
		err = validate.Struct(wareH)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			msgFields := ""

			//build message to know in which field is the error
			for _, ve := range validateErr {
				msgFields += ve.Field() + "-" + ve.Tag() + ","
			}
			web.Error(c, 400, msgFields)
			return
		}

		// I pass the evaluated request struct to the
		// initial struct with id

		wareHBD := domain.Warehouse{
			ID:                 id,
			Address:            wareH.Address,
			Telephone:          wareH.Telephone,
			WarehouseCode:      wareH.WarehouseCode,
			MinimumCapacity:    wareH.MinimumCapacity,
			MinimumTemperature: wareH.MinimumTemperature,
		}

		//update data in the BD
		wareHUpdate, er := w.s.Update(c, wareHBD)

		if er != nil {
			web.Error(c, 500, er.Error())
			return
		}

		web.Success(c, 200, wareHUpdate)

	}
}

// @summary		Delete warehouse
// @tags			Warehouse
// @Description	delete warehouse by id
// @Param			id	path	int	true	"Warehouse Id"
// @Success		204
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Router			/api/v1/warehouses/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id param from URL
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(c, 400, ErrBadRequest.Error())
			return
		}
		// delete product in BD
		err = w.s.Delete(c, id)
		if err != nil {
			web.Error(c, 404, err.Error())
			return
		}
		web.Success(c, 204, nil)
	}
}
