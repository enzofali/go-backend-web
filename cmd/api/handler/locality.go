package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/locality"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Locality struct {
	localityService locality.Service
}

func NewLocality(localityService locality.Service) *Locality {
	return &Locality{localityService: localityService}
}

// @Summary		Create locality
// @Tags			Localities
// @Description	Create locality
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Locality	true	"Locality parameters"
// @Success		201		{object}	web.response{data=domain.Locality}
// @Failure		422		{object}	web.errorResponse
// @Failure		409		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/localities [post]
func (l *Locality) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var localityRequest domain.Locality
		if err := ctx.ShouldBind(&localityRequest); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, "error bad request")
			return
		}

		// Validate the JSON
		// If the JSON object does not contain the necessary fields, a 422 code will be returned.
		validator := validator.New()
		if err := validator.Struct(&localityRequest); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		err := l.localityService.Create(ctx, localityRequest)
		if err != nil {
			switch err {
			case locality.ErrIntern:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				//log.Println(fmt.Sprintf("FATAL ERROR >  error: %s", err))
				return
			case locality.ErrDuplicated:
				web.Error(ctx, http.StatusConflict, err.Error())
				return
			default:
				web.Error(ctx, http.StatusInternalServerError, "internal error")
				return
			}
		}

		// Response
		// When the data entry is successful, a 201 code will be returned along with the entered object
		web.Success(ctx, http.StatusCreated, localityRequest)
	}
}

// @Summary		ReportSellers
// @Tags			Localities
// @Description	Returns a list of all reports
// @Param			id	query	string	false	"locality Id"
// @Produce		json
// @Success		200	{object}	web.response{domain.QuantitySellerByLocality}
// @Success		200	{object}	web.response{data=[]domain.QuantitySellerByLocality}
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/localities/reportSellers [get]
func (l *Locality) GetQuantitySellerByLocality() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := ctx.GetQuery("id")

		if !ok {
			result, err := l.localityService.GetSellerAll(ctx)
			if err != nil {
				web.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}

			web.Success(ctx, http.StatusOK, result)
			return
		}
		result, err := l.localityService.GetSellerByLocality(ctx, id)
		if err != nil {
			switch err {
			case locality.ErrLocalityNotFound:
				web.Error(ctx, http.StatusNotFound, err.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusOK, result)
	}
}
