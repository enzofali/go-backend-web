package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/carry"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Carry struct {
	s carry.Service
}

func NewCarry(s carry.Service) *Carry {
	return &Carry{
		s: s,
	}
}

/*
// @summary		list carries
// @tags		Carry
// @Description	Returns a list of all carries
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Carrie}
// @Failure		500	{object}	web.errorResponse
// @Router		/api/v1/carries/ [get]
*/
func (ca *Carry) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get and return all carries
		carryG, err := ca.s.GetAll(c)
		if err != nil {
			web.Error(c, 500, err.Error())
			return

		}
		web.Success(c, 200, carryG)
	}
}

// @summary		count carries by locality
// @tags			Carry
// @Description	Returns a list of all carries by locality
// @Produce		json
// @Param			id	query		string	false	"locality Id"
// @Success		200	{object}	web.response{data=[]domain.CarrieLocality}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/localities/reportCarries [get]
func (ca *Carry) GetAllByLocality() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, ok := c.GetQuery("id")

		if ok {

			found, err := ca.s.GetByLocalityID(c, id)

			if err != nil {
				web.Error(c, 404, err.Error())
				return
			}
			web.Success(c, 200, found)
			return

		}

		//get and return all carries
		carryG, err := ca.s.GetByLocality(c)
		if err != nil {
			web.Error(c, 500, err.Error())
			return

		}
		web.Success(c, 200, carryG)
	}
}

// @summary		Create carry
// @tags			Carry
// @Description	Create carry
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Carrie	true	"query params"
// @Success		201		{object}	web.response{data=domain.Carrie}
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Router			/api/v1/carries [post]
func (ca *Carry) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Validate format JSON, and return 422 if carry is
		//improperly constructed
		var carryRequest domain.Carrie
		err := c.ShouldBindJSON(&carryRequest)

		if err != nil {
			web.Error(c, 422, err.Error())
			return
		}

		//validate field empty
		validate := validator.New()
		err = validate.Struct(carryRequest)

		if err != nil {
			web.Error(c, 400, ErrBadRequest.Error())
			return
		}

		//Create carry and Validate type errors

		carryG, err := ca.s.Crear(c, carryRequest)

		if err == carry.ErrBD {
			web.Error(c, 500, err.Error())
			return

		}
		if err == carry.ErrExist {
			web.Error(c, 409, err.Error())
			return

		}

		if err == carry.ErrForeignKey {
			web.Error(c, 409, err.Error())
			return

		}

		web.Success(c, 201, carryG)

	}
}
