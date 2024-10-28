package http

import (
	"api-ticket/internal/entity"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// initTalentURLRoutes initializes event routes with Swagger documentation
func (r *Router) initEventURLRoutes(app *gin.RouterGroup) {
	events := app.Group("events")
	{
		events.POST("/", r.CreateEvent)
		events.GET("/", r.GetAllEvent)
		events.GET("/:id", r.GetEventById)
		events.PATCH("/:id", r.UpdateEvent)
		events.DELETE("/:id", r.DeleteEvent)
	}
}

// @Summary Get a event by ID
// @Description Retrieve a event by its ID
// @Tags events
// @Produce  json
// @Param id path string true "Event ID"
// @Success 200 {object} entity.Event
// @Router /events/{id} [get]
func (r *Router) GetEventById(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan event berdasarkan ID
	event, err := r.eventService.FindByID(c, id)
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
	event.Img_layout = serverAddress + event.Img_layout
	SendResponse(c, event, "success")
}

// @Summary Delete a event
// @Description Delete a event by its ID
// @Tags events
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Router /events/{id} [delete]
func (r *Router) DeleteEvent(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")
	// Panggil service untuk menemukan event berdasarkan ID
	event, err := r.eventService.Delete(c, id)
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
	event.Img_layout = serverAddress + event.Img_layout
	SendResponse(c, event, "success")
}

// @Summary Create a new event
// @Description Create a new event with the given information
// @Tags events
// @Accept multipart/form-data
// @Produce json
// @Param id_talent formData string true "id_talent of the event"
// @Param title formData string true "title of the event"
// @Param desc formData string true "desc of the event"
// @Param date formData string true "date of the event" format(date)
// @Param location formData string true "location of the event"
// @Param sk formData string true "sk of the event"
// @Param tag formData string true "tag of the event"
// @Param id_promotor_created formData string true "id_promotor_created of the event"
// @Param img_layout formData file true "img_layout file to upload"
// @Success 200 {object} entity.Event
// @Router /events/ [post]
func (r *Router) CreateEvent(c *gin.Context) {
	var req entity.EventInput
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	if err := r.eventService.Create(c, req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	SendResponse(c, nil, "success")
}

// @Summary Get all events
// @Description Retrieve a list of all events
// @Tags events
// @Produce  json
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param order_by query string false "order_by"
// @Param sort query string false "sort"
// @Param id_talent query string false "Filter by id_talent"
// @Param title query string false "Filter by title"
// @Param desc query string false "Filter by desc"
// @Param location query string false "Filter by location"
// @Param sk query string false "Filter by sk"
// @Param tag query string false "Filter by tag"
// @Param id_promotor_created query integer false "Filter by id_promotor_created"
// @Success 200 {array} entity.Event
// @Router /events/ [get]
func (r *Router) GetAllEvent(c *gin.Context) {
	var filter entity.RequestGetEvent
	if err := c.ShouldBindQuery(&filter); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	event, totalCount, err := r.eventService.FindAll(c, filter)
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
	for i := range event {
		event[i].Img_layout = serverAddress + event[i].Img_layout
	}

	//Final Response
	SendResponseGetAll(c, event, "Get all event", pagination)
}

// @Summary Update a event
// @Description Update a event by its ID
// @Tags events
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Event ID"
// @Param id_talent formData string true "id_talent of the event"
// @Param title formData string true "title of the event"
// @Param desc formData string true "desc of the event"
// @Param date formData string true "date of the event" format(date)
// @Param location formData string true "location of the event"
// @Param sk formData string true "sk of the event"
// @Param tag formData string true "tag of the event"
// @Param id_promotor_created formData string true "id_promotor_created of the event"
// @Param img_layout formData file true "img_layout file to upload"
// @Success 200 {object} entity.Event
// @Router /events/{id} [patch]
func (r *Router) UpdateEvent(c *gin.Context) {
	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Bind input dari request JSON ke struct input
	var input entity.EventInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	// Panggil service untuk memperbarui event berdasarkan ID
	event, err := r.eventService.Update(c, id, input)
	if err != nil {
		// Jika event tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			SendError(c, "error", err.Error(), 404)
			return
		}

		SendError(c, "error", err.Error(), 400)
	}

	// Jika berhasil, kembalikan respon sukses
	SendResponse(c, event, "success")
}
