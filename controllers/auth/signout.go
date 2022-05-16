package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

func (a *AuthController) SignOut(c *gin.Context) {
	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		response.BadRequest(c, "refresh token required")
		return
	}

	_, err := utils.JwtRefreshTokenVerify(refreshToken)
	if err != nil {
		response.BadRequest(c, "invalid refresh token")
		return
	}

	// delete refresh token
	result := a.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})

	if result.Error != nil {
		// check if the error is record not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.NotFound(c, "refresh token not found")
			return
		}
		response.ServerError(c, "failed to sign out")
		return
	}
	response.Ok(c, "successfully signed out")
}
