package http

import (
	"api-ticket/internal/entity"

	_ "api-ticket/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Router struct {
	bannerService entity.IBannerService
}

func RouteService(
	app *gin.RouterGroup,
	bannerService entity.IBannerService,
) {
	router := &Router{
		bannerService: bannerService,
	}

	//Apply Middleware here
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.handlers(app)
}

func (r *Router) handlers(app *gin.RouterGroup) {
	app.GET("/ping", ping)

	apiGroupV1 := app.Group("v1")
	{
		r.initBannerURLRoutes(apiGroupV1)
	}
}

func ping(c *gin.Context) {
	SendResponse(c, map[string]interface{}{"data": "ang ang ang ang"}, "Success")
}
