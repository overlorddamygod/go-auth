package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	authController "github.com/overlorddamygod/go-auth/controllers/auth"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: configs.GetConfig().AllowOrigins,
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

			mailer := mailer.NewMailer()
			auth := authController.NewAuthController(db.GetDB(), mailer)

			authGroup.POST("signup", auth.SignUp)
			authGroup.POST("signin", auth.SignIn)
			authGroup.POST("signout", auth.SignOut)
			authGroup.POST("refresh", auth.RefreshToken)
			authGroup.GET("verify", auth.VerifyLogin)
			authGroup.POST("verify", auth.VerifyLogin)
			authGroup.POST("request-password-reset", auth.RequestPasswordRecovery)
			authGroup.POST("reset-password", auth.PasswordReset)
			authGroup.GET("confirm", auth.ConfirmAccount)
			authGroup.POST("confirm", auth.ConfirmAccount)
			authGroup.Use(middlewares.IsLoggedIn())
			authGroup.GET("me", auth.GetMe)
		}
	}
	return router
}
