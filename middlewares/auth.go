package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/utils"
)

func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string = c.GetHeader("X-Access-Token")

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "access token is required",
			})
			c.Abort()

			return
		}

		token, err := utils.JwtAccessTokenVerify(accessToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "access token is invalid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "access token is invalid",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
