package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "server.com/docs"
)

/* 아래 항목이 swagger에 의해 문서화 된다. */
// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host metering-svc-hnc-metering.apps.blackwhale.cloud.hancom.com
// @BasePath /api
func main() {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	api := router.Group("/api")
	{
		api.GET("/:namespace/count", GetCountNamespace)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.Run(":8080")
}
