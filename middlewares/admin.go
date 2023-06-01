package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/utils/response"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string = c.GetHeader("x-api-key")

		if accessToken == "" {
			response.Unauthorized(c, "access token is required")
			c.Abort()
			return
		}

		if accessToken != configs.MainConfig.AdminSecret {
			response.Unauthorized(c, "access token is invalid")
			c.Abort()
			return
		}

		c.Next()
	}
}
