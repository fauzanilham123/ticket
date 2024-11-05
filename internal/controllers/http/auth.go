package http

import (
	"api-ticket/internal/entity"

	"github.com/gin-gonic/gin"
)

// initAuthURLRoutes initializes Auth routes with Swagger documentation
func (r *Router) initAuthURLRoutes(app *gin.RouterGroup) {
	auth := app.Group("auth")
	{
		auth.POST("/customer/register", r.Register)
		auth.POST("/customer/login", r.Login)
	}
}

// @Summary Register Customer
// @Description Register Customer with the given information
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Param gender formData string true "gender"
// @Param birthday formData string true "birthday"
// @Success 200 {object} entity.User
// @Router /auth/customer/register [post]
func (r *Router) Register(c *gin.Context) {
	var req entity.RegisterInputCustomer
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	user, err := r.authService.Register(c, req)
	if err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	SendResponse(c, user, "success")
}

// @Summary Login Customer
// @Description Login Customer
// @Tags auth
// @Produce  json
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Success 200 {array} entity.User
// @Router /auth/customer/login [post]
func (r *Router) Login(c *gin.Context) {
	var req entity.LoginInput
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	user, err := r.authService.LoginCustomer(c, req)
	if err != nil {
		SendError(c, "error", err.Error(), 400)
		return
	}

	SendResponse(c, user, "success")
}
