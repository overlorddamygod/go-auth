package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
)

func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string = c.GetHeader("X-Access-Token")

		if accessToken == "" {
			response.Unauthorized(c, "access token is required")
			c.Abort()
			return
		}

		token, err := utils.JwtAccessTokenVerify(accessToken)

		if err != nil {
			response.Unauthorized(c, "access token is invalid")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			response.Unauthorized(c, "access token is invalid")
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
