package services

import (
	"auth-service/dto/request"
	"auth-service/integrations/keycloak"
	"auth-service/utils"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserService interface {
	RegisterUserByAdmin(context *gin.Context)
}

type userService struct {
	identityManager keycloak.IIdentityManager
}

func NewUserService(IdentityManager keycloak.IIdentityManager) IUserService {
	return &userService{
		identityManager: IdentityManager,
	}
}

func (u userService) RegisterUserByAdmin(context *gin.Context) {
	utils.Log("INFO", "Start Create User...")

	var userRequest request.RegisterRequest
	err := context.ShouldBind(&userRequest)
	if err != nil {
		utils.Log("ERROR", err.Error())
		context.AbortWithStatusJSON(400, utils.BuildErrorResponse("unknown request body", "400", err.Error()))
		return
	}

	var validate = validator.New()
	err = validate.Struct(userRequest)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Error registering custom validator: %s", err))
		context.AbortWithStatusJSON(400, utils.BuildErrorResponse("invalid user ", "400", err.Error()))
		return
	}
	isValid := utils.ValidateUserRoles(userRequest.UserRoles)
	if !isValid {
		utils.Log("ERROR", fmt.Sprintf("Invalid roles %s,error:%s", userRequest.UserRoles, err))
		context.AbortWithStatusJSON(400, utils.BuildErrorResponse("invalid user ", "400", err.Error()))
		return
	}

	var user = gocloak.User{
		Username:      gocloak.StringP(userRequest.Username),
		FirstName:     gocloak.StringP(userRequest.FirstName),
		LastName:      gocloak.StringP(userRequest.LastName),
		Email:         gocloak.StringP(userRequest.Email),
		EmailVerified: gocloak.BoolP(true),
		Enabled:       gocloak.BoolP(true),
		Attributes:    &map[string][]string{},
	}
	if strings.TrimSpace(userRequest.MobileNumber) != "" {
		(*user.Attributes)["mobile"] = []string{userRequest.MobileNumber}
	}

	userResponse, err := u.identityManager.CreateUser(context, user, userRequest.Password, userRequest.UserRoles)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("create user error: %s", err))
		context.AbortWithStatusJSON(400, utils.BuildErrorResponse("user create error ", "400", err.Error()))
		return
	}
	utils.Log("INFO", fmt.Sprintf("Create User Success Username: %s", userRequest.Username))
	context.AbortWithStatusJSON(200, utils.BuildSuccessResponse("Login Successfully", "200", userResponse))
}
