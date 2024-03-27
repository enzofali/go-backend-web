package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/cmd/api/handler"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/buyer"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/carry"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/employee"
	inboundorder "github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/inbound_order"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/locality"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_batches"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/product_records"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/purchaseorder"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/section"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/seller"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/warehouse"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildSellerRoutes()
	r.buildProductRoutes()
	r.buildSectionRoutes()
	r.buildProductBatchesRoutes()
	r.buildWarehouseRoutes()
	r.buildEmployeeRoutes()
	r.buildBuyerRoutes()
	r.buildCarriesRoutes()
	r.buildPurchasOrderRoutes()
	r.buildInoundOrderRoutes()
	r.builLocalityRoutes()
	r.buildProductRecordRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	sr := r.rg.Group("/sellers")
	{
		sr.GET("/", handler.GetAll())
		sr.POST("/", handler.Create())
		sr.GET("/:id", handler.Get())
		sr.PATCH("/:id", handler.Update())
		sr.DELETE("/:id", handler.Delete())
	}
}

func (r *router) buildProductRoutes() {
	repo := product.NewRepository(r.db)
	service := product.NewService(repo)
	handler := handler.NewProduct(service)

	pr := r.rg.Group("/products")
	{
		pr.GET("/", handler.GetAll())
		pr.GET("/:id", handler.Get())
		pr.POST("/", handler.Create())
		pr.PATCH("/:id", handler.Update())
		pr.DELETE("/:id", handler.Delete())
		pr.GET("/reportRecords", handler.GetReport())
		pr.POST("/type", handler.CreateType())
	}
}

func (r *router) buildSectionRoutes() {
	repo := section.NewRepository(r.db)
	service := section.NewService(repo)
	handler := handler.NewSection(service)

	sections := r.rg.Group("/sections")
	{
		sections.GET("/", handler.GetAll())
		sections.GET("/:id", handler.Get())
		sections.GET("/reportProducts", handler.GetReportProducts())
		sections.POST("/", handler.Create())
		sections.PATCH("/:id", handler.Update())
		sections.DELETE("/:id", handler.Delete())
	}
}

func (r *router) buildProductBatchesRoutes() {
	repo := product_batches.NewRepository(r.db)
	service := product_batches.NewService(repo)
	handler := handler.NewProductBatches(service)

	productBatches := r.rg.Group("/productBatches")
	{
		productBatches.POST("/", handler.Create())
	}
}

func (r *router) buildWarehouseRoutes() {

	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	handler := handler.NewWarehouse(service)

	//r.eng.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	wareH := r.rg.Group("/warehouses")
	{
		wareH.GET("", handler.GetAll()) //http://localhost:8080/api/v1/warehouses
		wareH.GET(":id", handler.Get()) //http://localhost:8080/api/v1/warehouses/2
		wareH.POST("", handler.Create())
		wareH.PATCH(":id", handler.Update())
		wareH.DELETE(":id", handler.Delete())
	}

}

func (r *router) buildCarriesRoutes() {

	repo := carry.NewRepository(r.db)
	service := carry.NewService(repo)
	handler := handler.NewCarry(service)

	//localities/reportCarries
	carryG := r.rg.Group("/carries")
	{
		carryG.GET("", handler.GetAll())  //http://localhost:8080/api/v1/carries/
		carryG.POST("", handler.Create()) //http://localhost:8080/api/v1/carries/

	}

	r.rg.GET("/localities/reportCarries", handler.GetAllByLocality()) //http://localhost:8080/api/v1/localities/reportCarries?id=2001

}

func (r *router) buildEmployeeRoutes() {
	repo := employee.NewRepository(r.db)
	service := employee.NewService(repo)
	handler := handler.NewEmployee(service)

	rEmp := r.rg.Group("/employees")
	rEmp.GET("", handler.GetAll())
	rEmp.GET("/:id", handler.Get())
	rEmp.POST("", handler.Create())
	rEmp.PATCH("/:id", handler.Update())
	rEmp.DELETE("/:id", handler.Delete())
	rEmp.GET("/reportInboundOrders", handler.GetAllWithInboundOrders())

}

func (r *router) buildBuyerRoutes() {
	// Example
	//instances
	repo := buyer.NewRepository(r.db)
	service := buyer.NewService(repo)
	handler := handler.NewBuyer(service)

	//endpoints
	r.rg.GET("/buyers", handler.GetAll())
	r.rg.GET("/buyers/:id", handler.Get())
	r.rg.GET("/buyers/reportPurchaseOrders", handler.GetReport())
	r.rg.POST("/buyers", handler.Create())
	r.rg.PATCH("/buyers/:id", handler.Update())
	r.rg.DELETE("/buyers/:id", handler.Delete())
}

func (r *router) buildPurchasOrderRoutes() {
	repo := purchaseorder.NewRepository(r.db)
	service := purchaseorder.NewService(repo)
	handler := handler.NewPurchaseOrder(service)

	r.rg.POST("/purchaseorders", handler.Create())
}

func (r *router) buildInoundOrderRoutes() {
	repo := inboundorder.NewRepository(r.db)
	service := inboundorder.NewService(repo)
	handler := handler.NewInoudOrder(service)

	rEmp := r.rg.Group("/inboundOrders")
	rEmp.POST("", handler.Create())
}

func (r *router) builLocalityRoutes() {
	repo := locality.NewRepository(r.db)
	service := locality.NewService(repo)
	handler := handler.NewLocality(service)
	sr := r.rg.Group("/localities")

	//endpoints
	sr.POST("", handler.Create())
	sr.GET("/reportSellers", handler.GetQuantitySellerByLocality())
}

func (r *router) buildProductRecordRoutes() {
	repo := product_records.NewRepository(r.db)
	service := product_records.NewService(repo)
	handler := handler.NewProductRecords(service)

	pr := r.rg.Group("/productRecords")
	{
		pr.POST("/", handler.Create())
	}
}
