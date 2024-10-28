package http

import (
	"api-ticket/internal/entity"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// initBannerURLRoutes initializes banner routes with Swagger documentation
func (r *Router) initBannerURLRoutes(app *gin.RouterGroup) {
	banners := app.Group("banners")
	{
		banners.POST("/", r.CreateBanner)
		banners.GET("/", r.GetAllBanner)
		banners.GET("/:id", r.GetBannerById)
		banners.PATCH("/:id", r.UpdateBanner)
		banners.DELETE("/:id", r.DeleteBanner)
	}
}

// @Summary Get a banner by ID
// @Description Retrieve a banner by its ID
// @Tags banners
// @Produce  json
// @Param id path string true "Banner ID"
// @Success 200 {object} entity.Banner
// @Router /banners/{id} [get]
func (r *Router) GetBannerById(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan banner berdasarkan ID
	banner, err := r.bannerService.FindByID(c, id)
	if err != nil {
		SendError(c, "record not found", err.Error(), 404)
		return
	}

	var scheme string
	if c.Request.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}
	serverAddress := scheme + c.Request.Host
	banner.Img = serverAddress + banner.Img
	SendResponse(c, banner, "success")
}

// @Summary Delete a banner
// @Description Delete a banner by its ID
// @Tags banners
// @Param id path string true "Banner ID"
// @Success 200 {object} map[string]string
// @Router /banners/{id} [delete]
func (r *Router) DeleteBanner(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan banner berdasarkan ID
	banner, err := r.bannerService.Delete(c, id)
	if err != nil {
		SendError(c, "error", err.Error(), 404)
		return
	}

	var scheme string
	if c.Request.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}
	serverAddress := scheme + c.Request.Host
	banner.Img = serverAddress + banner.Img
	SendResponse(c, banner, "success")
}

// @Summary Create a new banner
// @Description Create a new banner with the given information
// @Tags banners
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "title of the banner"
// @Param slug formData string true "slug of the banner"
// @Param desc formData string true "desc of the banner"
// @Param img formData file true "Image file to upload"
// @Success 200 {object} entity.Banner
// @Router /banners/ [post]
func (r *Router) CreateBanner(c *gin.Context) {
	var req entity.BannerInput
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	if err := r.bannerService.Create(c, req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	SendResponse(c, nil, "success")
}

// @Summary Get all banners
// @Description Retrieve a list of all banners
// @Tags banners
// @Produce  json
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param order_by query string false "order_by"
// @Param sort query string false "sort"
// @Param title query string false "Filter by title"
// @Success 200 {array} entity.Banner
// @Router /banners/ [get]
func (r *Router) GetAllBanner(c *gin.Context) {
	var filter entity.RequestGetBanner
	if err := c.ShouldBindQuery(&filter); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	banner, totalCount, err := r.bannerService.FindAll(c, filter)
	if err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	//Pagination
	var pagination Paginations
	//pagination := utils.CalculatePagination(int(totalCount), filter.Limit, filter.Offset)

	if filter.Limit != nil {
		limit := *filter.Limit
		pagination.PerPage = int(limit)
		pagination.TotalPages = int(math.Ceil(float64(totalCount) / float64(limit)))

		if filter.Offset != nil {
			offset := *filter.Offset
			pagination.CurrentPage = int(offset/limit) + 1

			if pagination.CurrentPage < pagination.TotalPages {
				pagination.NextPage = pagination.CurrentPage + 1
			}

			if pagination.CurrentPage > 1 {
				pagination.PrevPage = pagination.CurrentPage - 1
			}
		} else {
			pagination.CurrentPage = 1
			if pagination.TotalPages > 1 {
				pagination.NextPage = 2
			}
		}

		pagination.LastPage = pagination.TotalPages
	}

	//
	var scheme string
	if c.Request.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}
	serverAddress := scheme + c.Request.Host
	for i := range banner {
		banner[i].Img = serverAddress + banner[i].Img
	}

	//Final Response
	SendResponseGetAll(c, banner, "Get all banner", pagination)
}

// @Summary Update a banner
// @Description Update a banner by its ID
// @Tags banners
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Banner ID"
// @Param title formData string true "title of the banner"
// @Param slug formData string true "slug of the banner"
// @Param desc formData string true "desc of the banner"
// @Param img formData file true "Image file to upload"
// @Success 200 {object} entity.Banner
// @Router /banners/{id} [patch]
func (r *Router) UpdateBanner(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Bind input dari request JSON ke struct input
	var input entity.BannerInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	// Panggil service untuk memperbarui banner berdasarkan ID
	banner, err := r.bannerService.Update(c, id, input)
	if err != nil {
		// Jika banner tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			SendError(c, "error", err.Error(), 404)
			return
		}

		SendError(c, "error", err.Error(), 400)
	}

	// Jika berhasil, kembalikan respon sukses
	SendResponse(c, banner, "success")
}
