package http

import (
	"api-ticket/internal/entity"
	"github.com/gin-gonic/gin"
	"math"
)

func (r *Router) initBannerURLRoutes(app *gin.RouterGroup) {
	banners := app.Group("banners")
	{
		banners.POST("/add", r.CreateBanner)
		banners.GET("/", r.GetAllBanner)
	}
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
