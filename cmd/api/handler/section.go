package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/section"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

var (
	ErrSectionId = errors.New("error: cannot update section id")
)

type Section struct {
	s section.Service
}

func NewSection(s section.Service) *Section {
	return &Section{
		s: s,
	}
}

// -------------------------------- GET Methods --------------------------------

// @Summary		List sections
// @Tags			Sections
// @Description	Get All Sections
// @Produce		json
// @Success		200	{object}	web.response
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/sections [get]
func (s *Section) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request

		// Process
		sections, err := s.s.GetAll(ctx)
		if err != nil {
			web.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		// When the request is successful, the backend will return a list of all existing sections
		web.Success(ctx, http.StatusOK, sections)
	}
}

// @Summary		Section by id
// @Tags			Sections
// @Description	Get Section by Id
// @Produce		json
// @Param			id	path		int	true	"section id"
// @Success		200	{object}	web.response
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Router			/api/v1/sections/{id} [get]
func (s *Section) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Response(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Process
		// When the seller does not exist, a 404 code will be returned
		sec, err := s.s.GetByID(ctx, id)
		if err != nil {
			switch err {
			case section.ErrSectionNotFound:
				web.Error(ctx, http.StatusNotFound, err.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		// Response
		// When the request is successful, the backend will return the information of the requested section
		web.Success(ctx, http.StatusOK, sec)
	}
}

// @Summary		Report Products
// @Tags			Sections
// @Description	Get the quantity of products of each section or the quantity of products for a determined section
// @Produce		json
// @Param			id	query		int	false	"Section id"
// @Success		200	{object}	web.response{data=[]domain.SectionReportProducts}
// @Success		200	{object}	web.response{data=domain.SectionReportProducts}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/sections/reportProducts [get]
func (s *Section) GetReportProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var id int
		var err error
		// Defaul id == 0, if id request is valid set id to the given parameter
		stringId, ok := ctx.GetQuery("id")
		if ok {
			id, err = strconv.Atoi(stringId)
			if err != nil {
				web.Response(ctx, http.StatusBadRequest, ErrInvalidId.Error())
				return
			}
		}

		// Process
		report, err := s.s.GetReportProducts(ctx, id)
		if err != nil {
			switch err {
			case section.ErrSectionNotFound:
				web.Error(ctx, http.StatusNotFound, err.Error())
				return
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Response
		// When the request is successful, the backend will return the report list
		web.Success(ctx, http.StatusOK, report)
	}
}

// -------------------------------- POST Methods --------------------------------

// @Summary		Create section
// @Tags			Sections
// @Description	Create section
// @Accept			json
// @Produce		json
// @Param			section	body		domain.Section	true	"Section to create"
// @Success		201		{object}	web.response
// @Failure		422		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/sections [post]
func (s *Section) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request domain.Section

		// Bind JSON to domain.Section{}
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

		// Validate unique section_number: If the section_number already exists, return a 409 Conflict error
		sec, err := s.s.Create(ctx, request)
		if err != nil {
			switch err {
			case section.ErrExistsSectionNumber:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case section.ErrWareHouseNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case section.ErrProductTypeNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Response
		// When the data entry is successful, a 201 code will be returned along with the entered object
		web.Success(ctx, http.StatusCreated, sec)
	}
}

// -------------------------------- PATCH Methods --------------------------------

// @Summary		Update section
// @Tags			Sections
// @Description	Update section
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"section id"
// @Param			section	body		domain.Section	true	"Section to update"
// @Success		200		{object}	web.response
// @Failure		400		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/sections/{id} [patch]
func (s *Section) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// If the section to be updated does not exist, a 404 code will be returned
		sectionDB, err := s.s.GetByID(ctx, id)
		if err != nil {
			web.Error(ctx, http.StatusNotFound, err.Error())
			return
		}

		// Bind the given JSON key:value to the required section
		err = json.NewDecoder(ctx.Request.Body).Decode(&sectionDB)
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// New id should not be specified in request body
		if sectionDB.ID != id {
			web.Error(ctx, http.StatusBadRequest, ErrSectionId.Error())
			return
		}

		// Validate the JSON
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validator := validator.New()
		if err := validator.Struct(&sectionDB); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Validate unique section_number: If the section_number already exists, return a 409 Conflict error
		err = s.s.Update(ctx, sectionDB)
		if err != nil {
			switch err {
			case section.ErrExistsSectionNumber:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case section.ErrWareHouseNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			case section.ErrProductTypeNotFound:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Response
		// When the data update is successful, the section with the updated information will be returned along with a code 200
		web.Success(ctx, http.StatusOK, sectionDB)
	}
}

// -------------------------------- DELETE Methods --------------------------------

// @Summary		Delete section
// @Tags			Sections
// @Description	Delete section
// @Param			id	path		int	true	"section id"
// @Success		204	{object}	web.response
// @Failure		404	{object}	web.errorResponse
// @Failure		400	{object}	web.errorResponse
// @Router			/api/v1/sections/{id} [delete]
func (s *Section) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		// Get the id and validate it
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Process
		// When the section does not exist a 404 code will be returned
		err = s.s.Delete(ctx, id)
		if err != nil {
			web.Error(ctx, http.StatusNotFound, err.Error())
			return
		}

		// Response
		// When the deletion is successful, a 204 code will be returned.
		web.Success(ctx, http.StatusNoContent, nil)
	}
}
