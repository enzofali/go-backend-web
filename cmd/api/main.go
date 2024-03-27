package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/cmd/api/routes"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/docs"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/pkg/database"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			MeLi Bootcamp API
// @version		1.0
// @description	This API was created as part of Mercado Libre's Go backend bootcamp for training purposes.
// @termsOfService	https://developers.mercadolibre.com.ar/es-ar/terminos-y-condiciones
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @description start with Bearer
// @name x-tiger-token
func main() {
	db, err := database.NewDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	eng := gin.Default()

	eng.GET("/ping", func(c *gin.Context) { c.JSON(200, "pong") })

	docs.SwaggerInfo.Host = "test--bootcamp-go-w7-s4-8-3.furyapps.io"
	eng.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := routes.NewRouter(eng, db)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}
