package entity

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		Id          uint      `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email" gorm:"unique"`
		Id_type     uint      `json:"id_type"  gorm:"default:1"`
		Id_Promotor uint      `json:"id_promotor"`
		Id_Customer uint      `json:"id_customer"`
		Password    string    `gorm:"text" json:"-"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Token       string    `json:"token,omitempty"` // field token
	}

	IAuthService interface {
		Register(c *gin.Context, req RegisterInputCustomer) (user User, err error)
		LoginCustomer(c *gin.Context, req LoginInput) (user User, err error)
	}

	IAuthRepository interface {
		RegisterCustomer(user User, customer Customer) (User, error)
		LoginCustomer(user LoginInput) (User, error)
	}

	RegisterInputCustomer struct {
		Name        string `json:"name" form:"name" binding:"required"`
		Email       string `json:"email" form:"email" binding:"required"`
		Id_type     uint   `json:"id_type" form:"id_type"`
		Id_Promotor uint   `json:"id_promotor" form:"id_promotor"`
		Password    string `json:"password" form:"password" binding:"required"`
		Gender      string `json:"gender" form:"gender"`
		Birthday    string `json:"birthday" form:"birthday"`
		PhoneNumber string `json:"phone_number" form:"phone_number"`
	}

	LoginInput struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)

func TableUser() string {
	return "users"
}
