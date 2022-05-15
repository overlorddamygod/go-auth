package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

func (a *AuthController) RefreshToken(c *gin.Context) {
	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		response.BadRequest(c, "refresh token required")
		return
	}

	token, err := utils.JwtRefreshTokenVerify(refreshToken)

	if err != nil {
		response.Unauthorized(c, "refresh token invalid")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		response.Unauthorized(c, "refresh token invalid")
		return
	}

	userID := claims["user_id"].(uuid.UUID)
	email := claims["email"].(string)

	var dbUser models.User
	result := a.db.First(&dbUser, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.ServerError(c, "failed to refresh token")
		return
	}

	var refreshTokenModel models.RefreshToken
	result = a.db.First(&refreshTokenModel, "token = ?", refreshToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.NotFound(c, "refresh token invalid")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	if refreshTokenModel.Revoked {
		response.Unauthorized(c, "refresh token revoked")
		return
	}

	accessToken, aTerr := utils.JwtAccessToken(utils.CustomClaims{
		UserID: userID,
		Email:  email,
	})

	if aTerr != nil {
		response.ServerError(c, "failed to refresh token")
		return
	}

	response.WithCustomStatusAndMessage(c, http.StatusOK, gin.H{
		"error":        false,
		"access-token": accessToken,
	})
}
