//package main
//
//import (
//	"auth-service/config"
//	"auth-service/controllers"
//	"auth-service/middlewares"
//	"github.com/gin-gonic/gin"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//)
//
//type router struct {
//	authController controllers.IAuthController
//	jwtService     middlewares.IJWTService
//}
//
//func NewRouter(JWTService middlewares.IJWTService, AuthController controllers.IAuthController) *router {
//	return &router{
//		jwtService:     JWTService,
//		authController: AuthController,
//	}
//}
//
//func (r *router) Run() {
//	gin.SetMode(gin.ReleaseMode)
//	gr := gin.Default()
//	middlewares.RegisterPrometheusMetrics()
//	authRouterGroup := gr.Group("api/v1/auth", middlewares.RecordRequestLatency())
//	{
//		authRouterGroup.POST("/user-registration-by-admin", r.jwtService.AuthorizeByRole("admin"), r.authController.RegisterUserAsAdmin)
//	}
//	gr.GET("/metrics", prometheusHandler())
//	gr.Run(config.Configuration.Server.HTTPAddr)
//
//}
//func prometheusHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		h := promhttp.Handler()
//		h.ServeHTTP(c.Writer, c.Request)
//	}
//}

package main

import (
	"auth-service/config"
	"auth-service/controllers"
	"auth-service/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type router struct {
	authController controllers.IAuthController
	jwtService     middlewares.IJWTService
}

func NewRouter(jwtService middlewares.IJWTService, authController controllers.IAuthController) *router {
	return &router{
		jwtService:     jwtService,
		authController: authController,
	}
}

func (r *router) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	middlewares.RegisterPrometheusMetrics()
	router.Use(middlewares.RecordRequestLatency())

	authGroup := router.Group("/api/v1/auth")
	authGroup.Use(middlewares.RecordRequestLatency())
	authGroup.POST("/user-registration-by-admin", r.jwtService.AuthorizeByRole("admin"), r.authController.RegisterUserAsAdmin)

	router.GET("/metrics", prometheusHandler())

	// Note: Consider extracting the address to a parameter for flexibility.
	router.Run(config.Configuration.Server.HTTPAddr)
}

func prometheusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}
