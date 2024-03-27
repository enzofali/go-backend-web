package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/web"
)

type Product struct {
	productService product.Service
}

func NewProduct(s product.Service) *Product {
	return &Product{
		productService: s,
	}
}

var (
	ErrInvalidId  = errors.New("invalid id")
	ErrInternal   = errors.New("internal server error")
	ErrBadRequest = errors.New("bad request")
)

// @summary		List products
// @tags			Products
// @Description	Returns a list of all products
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Product}
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/products [get]
func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get and return all products
		products, err := p.productService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		web.Success(c, http.StatusOK, products)
	}
}

// @summary		Get product by ID
// @tags			Products
// @Description	Returns a single product specified by its ID passed as a URL parameter
// @Param			id	path	int	true	"Product ID"
// @Produce		json
// @Success		200	{object}	web.response{data=domain.Product}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Router			/api/v1/products/{id} [get]
func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id param from URL, must be integer
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}
		// get product defined by id param, return 404 if product with given id doesn't exist
		prod, err := p.productService.GetByID(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, product.ErrNotFound.Error())
			return
		}
		web.Success(c, http.StatusOK, prod)
	}
}

// @summary		Create product
// @tags			Products
// @Description	Creates and returns a single product
// @Accept			json
// @Produce		json
// @Param			request	body		domain.ProductRequest	true	"Product parameters"
// @Success		201		{object}	web.response{data=domain.Product}
// @Failure		400		{object}	web.errorResponse
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/products/ [post]
func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind body to product struct, return error 422 if product is improperly constructed
		var prodToCreate domain.Product
		err := c.ShouldBindJSON(&prodToCreate)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}
		validator := validator.New()
		if err := validator.Struct(&prodToCreate); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}
		// create product in backend and return it
		prod, err := p.productService.Create(c, prodToCreate)
		switch err {
		case product.ErrExists:
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		case product.ErrDatabase:
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		web.Success(c, http.StatusCreated, prod)
	}
}

// @summary		Update product
// @tags			Products
// @Description	Updates the product specified by URL id parameter with fields passed by request body. All object fields are optional: only the given fields will be updated. Returns updated object.
// @Param			id	path	int	true	"Product ID"
// @Accept			json
// @Produce		json
// @Param			request	body		domain.ProductRequest	true	"Product parameters"
// @Success		200		{object}	web.response{data=domain.Product}
// @Failure		400		{object}	web.errorResponse
// @Failure		404		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/products/{id} [patch]
func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id param from URL, must be integer
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}
		// get product that must be updated from id param
		productToUpdate, err := p.productService.GetByID(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, product.ErrNotFound.Error())
			return
		}
		// decode and update fetched product object with fields decoded from request body
		body := c.Request.Body
		err = json.NewDecoder(body).Decode(&productToUpdate)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrBadRequest.Error())
			return
		}
		// new id should not be specified in request body, i.e. it should not change
		if productToUpdate.ID != id {
			web.Error(c, http.StatusBadRequest, "cannot update product id")
			return
		}
		// update product in backend
		prod, err := p.productService.Update(c, productToUpdate)
		switch err {
		case product.ErrExists:
			// this error occurs when a new product code is provided but is already in use
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		case product.ErrDatabase:
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		web.Success(c, http.StatusOK, prod)
	}
}

// @summary		Delete product
// @tags			Products
// @Description	Deletes the product specified by URL id parameter.
// @Param			id	path	int	true	"Product ID"
// @Produce		json
// @Success		204	{object}	web.response
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/products/{id} [delete]
func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get id param from URL, must be integer
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidId.Error())
			return
		}
		// delete product in backend
		err = p.productService.Delete(c, id)
		switch err {
		case product.ErrNotFound:
			web.Error(c, http.StatusNotFound, err.Error())
			return
		case product.ErrDatabase:
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		web.Success(c, http.StatusNoContent, gin.H{})
	}
}

// @summary		Create product type
// @tags			Products
// @Description	Given a name, creates a product type with that name
// @Accept			json
// @Produce		json
// @Param			request	body		domain.ProductTypeRequest	true	"Product type"
// @Success		201		{object}	web.response{data=domain.ProductType}
// @Failure		422		{object}	web.errorResponse
// @Failure		500		{object}	web.errorResponse
// @Router			/api/v1/products/type [post]
func (p *Product) CreateType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.ProductTypeRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}

		validator := validator.New()
		if err := validator.Struct(&req); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrField.Error())
			return
		}

		id, err := p.productService.CreateType(context.Background(), req.Name)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
			return
		}

		web.Success(c, http.StatusCreated, domain.ProductType{Name: req.Name, ID: id})
	}
}

// @summary		Get product record report
// @tags			Products
// @Description	Given a product id as a query, it will return the amount of product records for that given product. If given no id, it will return the amount of product records for all products.
// @Param			id	query	int	false	"Product ID"
// @Produce		json
// @Success		200	{object}	web.response{data=[]domain.Report}
// @Failure		400	{object}	web.errorResponse
// @Failure		404	{object}	web.errorResponse
// @Failure		500	{object}	web.errorResponse
// @Router			/api/v1/products/reportRecords [get]
func (p *Product) GetReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get query parameter "id"; if it exists, ok will be true, otherwise ok will be false
		idQuery, ok := c.GetQuery("id")
		// dispatch appropriate handler function
		if ok {
			p.getOneReport(c, idQuery)
		} else {
			p.getAllReports(c)
		}
	}
}

// returns the number of product records associated with a single product_id passed by query
// as well as the product's description and ID
func (p *Product) getOneReport(c *gin.Context, idQuery string) {
	// converts id parameter from string to int, returns with error status 400 on failure
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	// validates that the given product_id associated with the product record corresponds
	// to an existing product in the database, returns with error status 404 otherwise
	if !p.productService.ValidateProductID(context.Background(), id) {
		web.Error(c, http.StatusNotFound, ErrNotFound.Error())
		return
	}

	// fetches the number of product records associated with given product_id and the
	// product's description, returns with error status 500 on failure
	report, err := p.productService.GetOneReport(context.Background(), id)
	if err != nil {
		web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	web.Success(c, http.StatusOK, report)
}

// returns the count of records and the product description and ID for every product in the database
func (p *Product) getAllReports(c *gin.Context) {
	// fetches the report, returns with status error 500 on failure
	report, err := p.productService.GetAllReports(context.Background())
	if err != nil {
		web.Error(c, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	web.Success(c, http.StatusOK, report)
}
