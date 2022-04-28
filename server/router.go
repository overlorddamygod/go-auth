package server

import (
	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("api/v1")
	{
		authGroup := v1.Group("auth")
		{
			auth := new(controllers.AuthController)
			authGroup.POST("signup", auth.SignUp)
			authGroup.POST("signin", auth.SignIn)
			authGroup.POST("refresh", auth.RefreshToken)
			authGroup.POST("verify", auth.VerifyLogin)
			authGroup.POST("requestpasswordrecovery", auth.RequestPasswordRecovery)
			authGroup.POST("passwordRecovery", auth.RequestPasswordRecovery)
		}
	}
	return router
}
