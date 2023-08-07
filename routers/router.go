package routers

import (
	"github.com/gin-gonic/gin"
	"superindo_diksha_test/config"
	"superindo_diksha_test/database"
	"superindo_diksha_test/handler"
	"superindo_diksha_test/middleware"
)

func AppRouter(e *gin.Engine, config config.Config, db database.AppDatabase) {

	hdl := handler.NewHandler(config, db)

	g := e.Group("/api/v1/products", middleware.AuthMiddleware(config.UseAuth, config.AuthKey))

	{
		g.GET("/source", hdl.SourceProduct)

		g.GET("/destination", hdl.DestinationProduct)

		g.POST("/", hdl.UpdateCheckProduct)
	}
}
