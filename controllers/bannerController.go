package controllers

import (
	"api-ticket/constanta"
	"api-ticket/models"
	"api-ticket/utils"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bannerInput struct {
	Title     string    `json:"title" form:"title" binding:"required"`
	Slug      string    `json:"slug" form:"slug" binding:"required"`
	Desc      string    `json:"desc" form:"desc" binding:"required"`
	Flag      bool      `json:"flag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllBanner(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var banner []models.Banner

	// Ambil nilai dari query parameter "sort" (asc/desc)
	sort := c.DefaultQuery("orderBy", "asc")
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}

	// Ambil nilai dari query parameter "sortBy" (default ke "id")
	sortBy := c.DefaultQuery("sortBy", "id")

	columns := GetColumns(models.Banner{})

	// Validasi kolom sortBy, jika tidak valid maka default ke "id"
	if _, valid := columns[sortBy]; !valid {
		sortBy = "id"
	}

	// Ambil pagination
	pagination := ExtractPagination(c)
	query := db.Where("flag = 1")

	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page', 'perPage', 'sort', dan 'sortBy' dari queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "orderBy")
	delete(queryParams, "sortBy")
	for column, values := range queryParams {
		value := values[0] // In case there are multiple values, we take the first one

		// Apply filtering condition if the value is not empty
		if value != "" {
			query = query.Where(column+" LIKE ?", "%"+value+"%")
		}
	}

	var totalCount int64
	query.Model(&banner).Where("flag = 1").Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting (sorting by column and order)
	err := query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&banner).Error
	if err != nil {
		SendError(c, "internal server error", err.Error(), 400)
		return
	}

	// Calculate "last_page" based on total pages
	lastPage := totalPages

	// Calculate "nextPage" and "prevPage"
	nextPage := pagination.Page + 1
	if nextPage > totalPages {
		nextPage = 1
	}

	prevPage := pagination.Page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	// Mendapatkan alamat server dari permintaan
	var scheme string
	if c.Request.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	// Gabungkan dengan host
	serverAddress := scheme + c.Request.Host

	// Mengubah setiap entri dalam data banner untuk menambahkan URL lengkap
	for i := range banner {
		banner[i].Img = serverAddress + banner[i].Img
	}

	// Buat paginasi
	paginations := Paginations{
		pagination.Page,
		pagination.PerPage,
		lastPage,
		nextPage,
		prevPage,
		totalPages,
		totalCount,
	}

	// Kirim response
	SendResponseGetAll(c, banner, "Get all banner", paginations)
}

func CreateBanner(c *gin.Context) {
	// Validate input
	var input bannerInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	filename, errfilename := HandleUploadFile(c, "img")
	if errfilename != nil {
		SendError(c, "Upload error", errfilename.Error(), 400)
		return
	}

	img := constanta.DIR_FILE + filename

	// Create
	banner := models.Banner{
		Title:     input.Title,
		Slug:      input.Slug,
		Desc:      input.Desc,
		Img:       img,
		Flag:      true,
		CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&banner)

	SendResponse(c, banner, "success")
}

func GetBannerById(c *gin.Context) { // Get model if exist
	var banner models.Banner

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
		SendError(c, "Record not found", err.Error(), 404)
		return
	}

	// Mendapatkan alamat server dari permintaan
	var scheme string
	// Cek apakah request menggunakan HTTPS (TLS)
	if c.Request.TLS != nil {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	// Gabungkan dengan host
	serverAddress := scheme + c.Request.Host

	// Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
	banner.Img = serverAddress + banner.Img

	SendResponse(c, banner, "Get banner by id "+c.Param("id"))
}

func UpdateBanner(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var banner models.Banner
	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
		SendError(c, "Record not found", err.Error(), 404)
		return
	}

	// Validate input
	var input bannerInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	filename, errfilename := HandleUploadFile(c, "img")
	if errfilename != nil {
		SendError(c, "File upload error", errfilename.Error(), 400)
		return
	}

	// Cek apakah ada file yang diunggah
	if filename != "" {
		// Hapus gambar lama jika ada
		if banner.Img != "" {
			oldImage := "." + banner.Img
			utils.RemoveFile(oldImage)
		}
	}

	// Jika ada file yang diunggah, set nama file yang baru
	if filename != "" {
		banner.Img = constanta.DIR_FILE + filename
	}

	var updatedInput models.Banner
	updatedInput.Title = input.Title
	updatedInput.Slug = input.Slug
	updatedInput.Desc = input.Desc
	updatedInput.Img = banner.Img
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&banner).Updates(updatedInput)
	SendResponse(c, banner, "success")

}

func DeleteBanner(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var banner models.Banner
	if err := db.Where("id = ?", c.Param("id")).First(&banner).Error; err != nil {
		SendError(c, "Record not found", err.Error(), 404)
		return
	}

	// Set the flag to 0
	if err := db.Model(&banner).Update("flag", false).Error; err != nil {
		SendError(c, "Failed to delete", err.Error(), 500)
		return
	}

	// Return success response
	SendResponse(c, banner, "success")
}
