package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
)

func (a *AuthController) RefreshToken(c *gin.Context) {
	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "refresh token required",
		})
		return
	}

	token, err := utils.JwtRefreshTokenVerify(refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token invalid",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token invalid",
		})
		return
	}

	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)

	var refreshTokenModel models.RefreshToken
	result := a.db.First(&refreshTokenModel, "token = ?", refreshToken)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token expired",
		})
		return
	}

	if refreshTokenModel.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token revoked",
		})
		return
	}

	accessToken, aTerr := utils.JwtAccessToken(utils.CustomClaims{
		UserID: userID,
		Email:  email,
	})

	if aTerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "failed to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":        false,
		"access-token": accessToken,
	})
}
