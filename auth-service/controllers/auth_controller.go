package controllers

import (
	"auth-service/services"
	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	RegisterUserAsAdmin(context *gin.Context)
}
type authController struct {
	userService services.IUserService
}

func NewAuthController(UserService services.IUserService) IAuthController {
	return &authController{
		userService: UserService,
	}
}

func (a authController) RegisterUserAsAdmin(context *gin.Context) {
	a.userService.RegisterUserByAdmin(context)
}
