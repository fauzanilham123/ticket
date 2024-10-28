package entity

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	Talent struct {
		Id                  uint      `gorm:"primary_key" json:"id"`
		Name                string    `json:"name"`
		Photo               string    `json:"photo" gorm:"text"`
		Id_promotor_created int       `json:"id_promotor_created"`
		CreatedAt           time.Time `json:"created_at"`
		UpdatedAt           time.Time `json:"updated_at"`
	}

	ITalentService interface {
		Create(c *gin.Context, req TalentInput) (err error)
		Update(c *gin.Context, id string, req TalentInput) (talent Talent, err error)
		FindByID(c *gin.Context, id string) (talent Talent, err error)
		Delete(c *gin.Context, id string) (talent Talent, err error)
		FindAll(c *gin.Context, filter RequestGetTalent) ([]Talent, int64, error)
	}

	ITalentRepository interface {
		Create(talent Talent) error
		Update(id string, talent Talent) (Talent, error)
		FindByID(id string) (Talent, error)
		Delete(id string) (Talent, error)
		FindAll(filter RequestGetTalent) ([]Talent, int64, error)
	}

	TalentInput struct {
		Name                string `json:"name" form:"name"`
		Id_promotor_created int    `json:"id_promotor_created" form:"id_promotor_created" `
	}

	RequestGetTalent struct {
		Limit               *int     `json:"limit" form:"limit"`
		Offset              *int     `json:"offset" form:"offset"`
		OrderBy             *OrderBy `json:"order_by" form:"order_by"`
		Sort                *string  `json:"sort" form:"sort"`
		Name                *string  `json:"name" form:"name"`
		Id_promotor_created *int     `json:"id_promotor_created" form:"id_promotor_created"`
	}
)

func TableTalent() string {
	return "talents"
}
