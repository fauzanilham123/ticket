package entity

import (
	"github.com/gin-gonic/gin"
	"time"
)

type (
	Banner struct {
		Id        uint      `gorm:"primary_key" json:"id"`
		Title     string    `json:"title"`
		Slug      string    `json:"slug" gorm:"primary_key"`
		Desc      string    `gorm:"text" json:"desc"`
		Img       string    `gorm:"text" json:"img"`
		Flag      bool      `json:"flag"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	IBannerService interface {
		Create(c *gin.Context, req BannerInput) (err error)
		FindAll(c *gin.Context, filter RequestGetBanner) ([]Banner, int64, error)
	}

	IBannerRepository interface {
		Create(banner Banner) error
		FindAll(filter RequestGetBanner) ([]Banner, int64, error)
	}

	BannerInput struct {
		Title     string    `json:"title" form:"title" binding:"required"`
		Slug      string    `json:"slug" form:"slug" binding:"required"`
		Desc      string    `json:"desc" form:"desc" binding:"required"`
		Flag      bool      `json:"flag"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	RequestGetBanner struct {
		Limit   *int     `json:"limit" form:"limit"`
		Offset  *int     `json:"offset" form:"offset"`
		OrderBy *OrderBy `json:"order_by" form:"order_by"`
		Sort    *string  `json:"sort" form:"sort"`
		Title   *string  `json:"title" form:"title"`
	}
)

func TableName() string {
	return "banners"
}