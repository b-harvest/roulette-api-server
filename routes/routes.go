package routes

import (
	"roulette-api-server/controllers"
	"roulette-api-server/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))

	route.POST("auth/signin", middlewares.IsClientAuthenticated, controllers.AuthSignin)
	route.DELETE("auth/signout", middlewares.IsUserAuthenticated, controllers.AuthSignout)

	route.GET("users", middlewares.IsUserAuthenticated, controllers.UserFetchAll)
	route.GET("users/:id", middlewares.IsUserAuthenticated, controllers.UserFetchSingle)
	route.POST("users", middlewares.IsClientAuthenticated, controllers.UserCreate)
	route.PUT("users/:id", middlewares.IsUserAuthenticated, controllers.UserUpdate)
	route.DELETE("users/:id", middlewares.IsUserAuthenticated, controllers.UserDelete)

	return route
}