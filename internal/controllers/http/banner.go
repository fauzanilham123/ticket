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
		// @Summary Create a new banner
		// @Description Create a new banner with the given information
		// @Tags banners
		// @Accept  json
		// @Produce  json
		// @Param banner body entity.BannerInput true "Banner input"
		// @Success 200 {object} entity.Banner
		// @Failure 400 {object} ErrorResponse
		// @Router /banners/ [post]
		banners.POST("/", r.CreateBanner)

		// @Summary Get all banners
		// @Description Retrieve a list of all banners
		// @Tags banners
		// @Produce  json
		// @Success 200 {array} entity.Banner
		// @Failure 400 {object} ErrorResponse
		// @Router /banners/ [get]
		banners.GET("/", r.GetAllBanner)

		// @Summary Get a banner by ID
		// @Description Retrieve a banner by its ID
		// @Tags banners
		// @Produce  json
		// @Param id path string true "Banner ID"
		// @Success 200 {object} entity.Banner
		// @Failure 404 {object} ErrorResponse
		// @Router /banners/{id} [get]
		banners.GET("/:id", r.GetBannerById)

		// @Summary Update a banner
		// @Description Update a banner by its ID
		// @Tags banners
		// @Accept  json
		// @Produce  json
		// @Param id path string true "Banner ID"
		// @Param banner body entity.BannerInput true "Banner input"
		// @Success 200 {object} entity.Banner
		// @Failure 400 {object} ErrorResponse
		// @Router /banners/{id} [patch]
		banners.PATCH("/:id", r.UpdateBanner)

		// @Summary Delete a banner
		// @Description Delete a banner by its ID
		// @Tags banners
		// @Param id path string true "Banner ID"
		// @Success 200 {object} map[string]string
		// @Failure 404 {object} ErrorResponse
		// @Router /banners/{id} [delete]
		banners.DELETE("/:id", r.DeleteBanner)
	}
}
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

//
//func GetBannerById(c *gin.Context) { // Get model if exist
//	var banner models.Banner
//
//	db := c.MustGet("db").(*gorm.DB)
//	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
//		SendError(c, "Record not found", err.Error(), 404)
//		return
//	}
//
//	// Mendapatkan alamat server dari permintaan
//	var scheme string
//	// Cek apakah request menggunakan HTTPS (TLS)
//	if c.Request.TLS != nil {
//		scheme = "https://"
//	} else {
//		scheme = "http://"
//	}
//
//	// Gabungkan dengan host
//	serverAddress := scheme + c.Request.Host
//
//	// Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
//	banner.Img = serverAddress + banner.Img
//
//	SendResponse(c, banner, "Get banner by id "+c.Param("id"))
//}
//
//func UpdateBanner(c *gin.Context) {
//
//	db := c.MustGet("db").(*gorm.DB)
//	// Get model if exist
//	var banner models.Banner
//	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
//		SendError(c, "Record not found", err.Error(), 404)
//		return
//	}
//
//	// Validate input
//	var input bannerInput
//	if err := c.ShouldBind(&input); err != nil {
//		SendError(c, "error", err.Error(), 400)
//		return
//	}
//
//	filename, errfilename := HandleUploadFile(c, "img")
//	if errfilename != nil {
//		SendError(c, "File upload error", errfilename.Error(), 400)
//		return
//	}
//
//	// Cek apakah ada file yang diunggah
//	if filename != "" {
//		// Hapus gambar lama jika ada
//		if banner.Img != "" {
//			oldImage := "." + banner.Img
//			utils.RemoveFile(oldImage)
//		}
//	}
//
//	// Jika ada file yang diunggah, set nama file yang baru
//	if filename != "" {
//		banner.Img = constanta.DIR_FILE + filename
//	}
//
//	var updatedInput models.Banner
//	updatedInput.Title = input.Title
//	updatedInput.Slug = input.Slug
//	updatedInput.Desc = input.Desc
//	updatedInput.Img = banner.Img
//	updatedInput.Flag = input.Flag
//	updatedInput.UpdatedAt = time.Now()
//
//	db.Model(&banner).Updates(updatedInput)
//	SendResponse(c, banner, "success")
//
//}
//
//func DeleteBanner(c *gin.Context) {
//	// Get model if exiForm
//	db := c.MustGet("db").(*gorm.DB)
//	var banner models.Banner
//	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
//		SendError(c, "Record not found", err.Error(), 404)
//		return
//	}
//
//	// Set the flag to 0
//	if err := db.Model(&banner).Update("flag", false).Error; err != nil {
//		SendError(c, "Failed to delete", err.Error(), 500)
//		return
//	}
//
//	// Return success response
//	SendResponse(c, banner, "success")
//}
