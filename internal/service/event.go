package service

import (
	"api-ticket/constanta"
	"api-ticket/internal/entity"
	"api-ticket/utils"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EventService struct {
	eventRepository entity.IEventRepository
}

func NewEventService(eventRepository entity.IEventRepository) entity.IEventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (service *EventService) FindByID(c *gin.Context, id string) (event entity.Event, err error) {
	return service.eventRepository.FindByID(id)
}

func (service *EventService) Delete(c *gin.Context, id string) (event entity.Event, err error) {
	event, err = service.eventRepository.FindByID(id)
	if err != nil {
		return event, err
	}

	// Step 2: Hapus gambar jika ada
	if event.Img_layout != "" {
		oldImage := "." + event.Img_layout
		if removeErr := utils.RemoveFile(oldImage); removeErr != nil {
			log.Println("Failed to delete image: ", removeErr)
			return event, removeErr
		}
	}

	return service.eventRepository.Delete(id)
}

func (service EventService) Create(c *gin.Context, req entity.EventInput) (err error) {
	filename, err := utils.HandleUploadFile(c, "img_layout")
	if err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	img_layout := constanta.DIR_FILE + filename

	tagWithDelimiter := strings.Join(req.Tag, ",")
	id_talentWithDelimiter := strings.Join(req.Id_talent, ",")

	// Create
	event := entity.Event{
		Id_talent:           id_talentWithDelimiter,
		Title:               req.Title,
		Desc:                req.Desc,
		Date:                req.Date,
		Location:            req.Location,
		Sk:                  req.Sk,
		Tag:                 tagWithDelimiter,
		Id_promotor_created: req.Id_promotor_created,
		Img_layout:          img_layout,
		CreatedAt:           time.Now(),
	}

	if err = service.eventRepository.Create(event); err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	return
}

func (service EventService) FindAll(c *gin.Context, filter entity.RequestGetEvent) ([]entity.Event, int64, error) {
	return service.eventRepository.FindAll(filter)
}

func (service EventService) Update(c *gin.Context, id string, req entity.EventInput) (event entity.Event, err error) {
	eventOld, err := service.eventRepository.FindByID(id)
	if err != nil {
		return eventOld, err
	}
	// Map EventInput ke entity.Event
	filename, err := utils.HandleUploadFile(c, "img_layout")
	if err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	// Cek apakah ada file yang diunggah
	if filename != "" {
		// Hapus gambar lama jika ada
		if eventOld.Img_layout != "" {
			oldImage := "." + eventOld.Img_layout
			utils.RemoveFile(oldImage)
		}
	}

	// Jika ada file yang diunggah, set nama file yang baru
	if filename != "" {
		event.Img_layout = constanta.DIR_FILE + filename
	}

	tagWithDelimiter := strings.Join(req.Tag, ",")
	id_talentWithDelimiter := strings.Join(req.Id_talent, ",")

	updatedEvent := entity.Event{
		Id_talent:           id_talentWithDelimiter,
		Title:               req.Title,
		Desc:                req.Desc,
		Date:                req.Date,
		Location:            req.Location,
		Sk:                  req.Sk,
		Tag:                 tagWithDelimiter,
		Id_promotor_created: req.Id_promotor_created,
		Img_layout:          event.Img_layout,
		UpdatedAt:           time.Now(),
	}

	return service.eventRepository.Update(id, updatedEvent)
}

// func (service EventService) Update(c *gin.Context, id string, req entity.EventInput) (err error) {

// }
