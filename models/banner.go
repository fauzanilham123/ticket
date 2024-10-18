package models

import "time"

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
)
