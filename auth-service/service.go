package main

import (
	"auth-service/controllers"
	"auth-service/integrations/keycloak"
	"auth-service/middlewares"
	"auth-service/services"
	"auth-service/utils"
)

type service struct {
}

func (s *service) Start() {

	//Dependency
	identityManager := keycloak.NewIdentityManager()
	jwtService := middlewares.NewJWTService()

	//Services
	userService := services.NewUserService(identityManager)

	//Controllers
	authController := controllers.NewAuthController(userService)

	r := NewRouter(jwtService, authController)
	r.Run()

}
func (s *service) Stop() {
	defer utils.LogFile.Close()
}
