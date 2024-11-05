package entity

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	Event struct {
		Id                  uint      `gorm:"primary_key" json:"id"`
		Id_talent           string    `json:"id_talent"`
		Title               string    `json:"title"`
		Desc                string    `json:"desc" gorm:"text"`
		Date                string    `json:"date"`
		Location            string    `json:"location"`
		Sk                  string    `json:"sk" gorm:"text"`
		Tag                 string    `json:"tag"`
		Id_promotor_created int       `json:"id_promotor_created"`
		Img_layout          string    `json:"img_layout" gorm:"text"`
		CreatedAt           time.Time `json:"created_at"`
		UpdatedAt           time.Time `json:"updated_at"`
	}

	IEventService interface {
		Create(c *gin.Context, req EventInput) (err error)
		Update(c *gin.Context, id string, req EventInput) (event Event, err error)
		FindByID(c *gin.Context, id string) (event Event, err error)
		Delete(c *gin.Context, id string) (event Event, err error)
		FindAll(c *gin.Context, filter RequestGetEvent) ([]Event, int64, error)
	}

	IEventRepository interface {
		Create(event Event) error
		Update(id string, event Event) (Event, error)
		FindByID(id string) (Event, error)
		Delete(id string) (Event, error)
		FindAll(filter RequestGetEvent) ([]Event, int64, error)
	}

	EventInput struct {
		Id_talent           []string `json:"id_talent" form:"id_talent"`
		Title               string   `json:"title" form:"title"`
		Desc                string   `json:"desc" form:"desc"`
		Date                string   `json:"date" form:"date"`
		Location            string   `json:"location" form:"location"`
		Sk                  string   `json:"sk" form:"sk"`
		Tag                 []string `json:"tag" form:"tag"`
		Id_promotor_created int      `json:"id_promotor_created" form:"id_promotor_created"`
	}

	RequestGetEvent struct {
		Limit               *int     `json:"limit" form:"limit"`
		Offset              *int     `json:"offset" form:"offset"`
		OrderBy             *OrderBy `json:"order_by" form:"order_by"`
		Sort                *string  `json:"sort" form:"sort"`
		Id_talent           *string  `json:"id_talent" form:"id_talent"`
		Title               *string  `json:"title" form:"title"`
		Desc                *string  `json:"desc" form:"desc"`
		Date                *string  `json:"date" form:"date"`
		Location            *string  `json:"location" form:"location"`
		Sk                  *string  `json:"sk" form:"sk"`
		Tag                 *string  `json:"tag" form:"tag"`
		Id_promotor_created *int     `json:"id_promotor_created" form:"id_promotor_created"`
	}
)

func TableEvent() string {
	return "events"
}
