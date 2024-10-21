package http

import (
	"api-ticket/internal/entity"
	"github.com/gin-gonic/gin"
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
