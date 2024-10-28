package http

import (
	"api-ticket/internal/entity"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// initTalentURLRoutes initializes talent routes with Swagger documentation
func (r *Router) initTalentURLRoutes(app *gin.RouterGroup) {
	talents := app.Group("talents")
	{
		talents.POST("/", r.CreateTalent)
		talents.GET("/", r.GetAllTalent)
		talents.GET("/:id", r.GetTalentById)
		talents.PATCH("/:id", r.UpdateTalent)
		talents.DELETE("/:id", r.DeleteTalent)
	}
}

// @Summary Get a talent by ID
// @Description Retrieve a talent by its ID
// @Tags talents
// @Produce  json
// @Param id path string true "Talents ID"
// @Success 200 {object} entity.Talent
// @Router /talents/{id} [get]
func (r *Router) GetTalentById(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan talent berdasarkan ID
	talent, err := r.talentService.FindByID(c, id)
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
	talent.Photo = serverAddress + talent.Photo
	SendResponse(c, talent, "success")
}

// @Summary Delete a talent
// @Description Delete a talent by its ID
// @Tags talents
// @Param id path string true "Talent ID"
// @Success 200 {object} map[string]string
// @Router /talents/{id} [delete]
func (r *Router) DeleteTalent(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan talent berdasarkan ID
	talent, err := r.talentService.Delete(c, id)
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
	talent.Photo = serverAddress + talent.Photo
	SendResponse(c, talent, "success")
}

// @Summary Create a new talent
// @Description Create a new talent with the given information
// @Tags talents
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name of the talent"
// @Param id_promotor_created formData string true "id_promotor_created of the talent"
// @Param photo formData file true "poto file to upload"
// @Success 200 {object} entity.Talent
// @Router /talents/ [post]
func (r *Router) CreateTalent(c *gin.Context) {
	var req entity.TalentInput
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	if err := r.talentService.Create(c, req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	SendResponse(c, nil, "success")
}

// @Summary Get all talents
// @Description Retrieve a list of all talents
// @Tags talents
// @Produce  json
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param order_by query string false "order_by"
// @Param sort query string false "sort"
// @Param name query string false "Filter by name"
// @Param id_promotor_created query string false "Filter by id_promotor_created"
// @Success 200 {array} entity.Talent
// @Router /talents/ [get]
func (r *Router) GetAllTalent(c *gin.Context) {
	var filter entity.RequestGetTalent
	if err := c.ShouldBindQuery(&filter); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	talent, totalCount, err := r.talentService.FindAll(c, filter)
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
	for i := range talent {
		talent[i].Photo = serverAddress + talent[i].Photo
	}

	//Final Response
	SendResponseGetAll(c, talent, "Get all talent", pagination)
}

// @Summary Update a talent
// @Description Update a talent by its ID
// @Tags talents
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Talent ID"
// @Param name formData string true "name of the talent"
// @Param id_promotor_created formData integer true "id_promotor_created of the talent"
// @Param photo formData file true "photo file to upload"
// @Success 200 {object} entity.Talent
// @Router /talents/{id} [patch]
func (r *Router) UpdateTalent(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Bind input dari request JSON ke struct input
	var input entity.TalentInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	// Panggil service untuk memperbarui talent berdasarkan ID
	talent, err := r.talentService.Update(c, id, input)
	if err != nil {
		// Jika talent tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			SendError(c, "error", err.Error(), 404)
			return
		}

		SendError(c, "error", err.Error(), 400)
	}

	// Jika berhasil, kembalikan respon sukses
	SendResponse(c, talent, "success")
}
