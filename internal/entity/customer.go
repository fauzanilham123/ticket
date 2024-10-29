package entity

import (
	"time"
)

type (
	Customer struct {
		Id        uint      `gorm:"primary_key" json:"id"`
		Name      string    `json:"name"`
		Gender    string    `json:"gender"`
		Birthday  string    `json:"birthday"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func TableCustomer() string {
	return "customers"
}
