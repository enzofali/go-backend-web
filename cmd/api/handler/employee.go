package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/employee"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

var (
	ErrNotFound       = errors.New("Employee not found.")
	ErrNotWareHouse   = errors.New("WareHouse not found.")
	ErrInternalServer = errors.New("Internal server error.")
)

type Employee struct {
	employeeService employee.Service
}

func NewEmployee(e employee.Service) *Employee {
	return &Employee{
		employeeService: e,
	}
}

// @summary		Get employee
// @tags			Employees
// @Description	get employee by id
// @Produce		json
// @Param			id	path		int	true	"Employee Id"
// @Success		200	{object}	web.response{data=domain.Employee}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Router			/api/v1/employees/{id} [get]
func (e *Employee) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}

		employeeDB, err := e.employeeService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, ErrNotFound.Error())
			return
		}

		web.Success(c, http.StatusOK, employeeDB)
	}
}

// @summary		List employees
// @tags			Employees
// @Description	get employees
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Employee}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/employees [get]
func (e *Employee) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := e.employeeService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, employee.ErrDatabase.Error())
			return
		}

		web.Success(c, http.StatusOK, employees)
	}
}

// @summary		Create employee
// @tags			Employees
// @Description	create employee
// @Accept			json
// @Param			request	body	domain.EmployeeRequest	true	"query params"
// @Produce		json
// @Success		201	{object}	web.response{data=domain.Employee}
// @Failure		422	{object}	web.errorResponse
// @Failure		400	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/employees [post]
func (e *Employee) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employeeRequest domain.EmployeeRequest
		err := c.ShouldBind(&employeeRequest)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrBadRequest.Error())
			return
		}

		validate := validator.New()
		err = validate.Struct(employeeRequest)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			msgFields := ""
			for _, ve := range validateErr {
				msgFields += ve.Field() + "-" + ve.Tag() + ","
			}
			web.Error(c, http.StatusUnprocessableEntity, msgFields)
			return
		}

		employeeSave := domain.Employee{
			CardNumberID: employeeRequest.CardNumberID,
			FirstName:    employeeRequest.FirstName,
			LastName:     employeeRequest.LastName,
			WarehouseID:  employeeRequest.WarehouseID,
		}
		employeeDB, err := e.employeeService.Create(c, employeeSave)
		if err != nil {
			switch err {
			case employee.ErrWarehouseNotfound:
				web.Error(c, http.StatusBadRequest, ErrNotWareHouse.Error())
				return
			case employee.ErrExistsCardId:
				web.Error(c, http.StatusBadRequest, ErrBadRequest.Error())
				return
			default:
				web.Error(c, http.StatusInternalServerError, ErrInternalServer.Error())
				return
			}
		}

		web.Success(c, http.StatusCreated, employeeDB)
	}
}

// @summary		Update employee
// @tags			Employees
// @Description	update employee
// @Accept			json
// @Produce		json
// @Param			request	body		domain.EmployeeRequest	true	"query params"
// @Param			id		path		int						true	"Employee Id"
// @Success		200		{object}	web.response{data=domain.Employee}
// @Failure		400		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/employees/{id} [patch]
func (e *Employee) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}

		employeeDB, err := e.employeeService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		employeeRequestDB := domain.EmployeeRequest{
			CardNumberID: employeeDB.CardNumberID,
			FirstName:    employeeDB.FirstName,
			LastName:     employeeDB.LastName,
			WarehouseID:  employeeDB.WarehouseID,
		}

		err = json.NewDecoder(c.Request.Body).Decode(&employeeRequestDB)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrBadRequest.Error())
			return
		}

		validate := validator.New()
		err = validate.Struct(employeeRequestDB)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			msgFields := ""
			for _, ve := range validateErr {
				msgFields += ve.Field() + "-" + ve.Tag() + ","
			}
			web.Error(c, http.StatusBadRequest, msgFields)
			return
		}

		employeeDB = domain.Employee{
			ID:           id,
			CardNumberID: employeeRequestDB.CardNumberID,
			FirstName:    employeeRequestDB.FirstName,
			LastName:     employeeRequestDB.LastName,
			WarehouseID:  employeeRequestDB.WarehouseID,
		}

		employeeNew, err := e.employeeService.Update(c, employeeDB)
		if err != nil {
			switch err {
			case employee.ErrWarehouseNotfound:
				web.Error(c, http.StatusBadRequest, ErrNotWareHouse.Error())
				return
			default:
				web.Error(c, http.StatusInternalServerError, ErrInternalServer.Error())
				return
			}

		}

		web.Success(c, http.StatusOK, employeeNew)
	}
}

// @summary		Delete employee
// @tags			Employees
// @Description	delete employee by id
// @Param			id	path	int	true	"Employee Id"
// @Success		204
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/employees/{id} [delete]
func (e *Employee) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}

		err = e.employeeService.Delete(c, id)
		if errors.Is(err, employee.ErrNotFound) {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		if err != nil {
			web.Error(c, http.StatusInternalServerError, employee.ErrDatabase.Error())
			return
		}

		web.Success(c, http.StatusNoContent, "")
	}
}

// @summary		Employee with inbound orders count
// @tags			Employees
// @Description	get employee with inbound orders count
// @Produce		json
// @Param			id	query		int	false	"Employee Id"
// @Success		200	{object}	web.response{data=[]domain.EmployeeWithInboundOrders}
// @Success		200	{object}	web.response{data=domain.EmployeeWithInboundOrders}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/employees/reportInboundOrders [get]
func (e *Employee) GetAllWithInboundOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		idQuery := c.Query("id")
		if idQuery == "" {
			employees, err := e.employeeService.GetAllInoundOrders(c)
			if err != nil {
				web.Error(c, http.StatusInternalServerError, ErrInternalServer.Error())
				return
			}

			web.Success(c, http.StatusOK, employees)
			return
		}

		id, err := strconv.Atoi(idQuery)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}

		employeeDB, err := e.employeeService.GetWithInboundOrder(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, ErrNotFound.Error())
			return
		}

		web.Success(c, http.StatusOK, employeeDB)
	}
}
