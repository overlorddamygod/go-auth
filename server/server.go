package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/controllers/auth"
	"github.com/overlorddamygod/go-auth/middlewares"
)

func RegisterServer(config *configs.Config, router *gin.Engine, authC *auth.AuthController) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowOrigins,
		AllowHeaders: []string{"content-type", "x-access-token", "x-refresh-token"},
	}))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("api/v1")
	{
		authGroup := v1.Group("auth")
		{
			authGroup.Use(func(c *gin.Context) {
				c.Writer.Header().Set("Content-Type", "application/json")
				c.Next()
			})

			authGroup.POST("signup", authC.SignUp)
			authGroup.POST("signin", authC.SignIn)
			authGroup.POST("signout", authC.SignOut)
			authGroup.POST("refresh", authC.RefreshToken)
			authGroup.GET("verify", authC.VerifyLogin)
			authGroup.POST("verify", authC.VerifyLogin)
			authGroup.POST("request-password-reset", authC.RequestPasswordRecovery)
			authGroup.POST("reset-password", authC.PasswordReset)
			authGroup.GET("confirm", authC.ConfirmAccount)
			authGroup.POST("confirm", authC.ConfirmAccount)

			authGroup.Use(middlewares.IsLoggedIn())
			authGroup.GET("me", authC.GetMe)
		}
	}
}
