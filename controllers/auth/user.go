package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

func (a *AuthController) GetMe(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		response.BadRequest(c, "user id is required")
		return
	}

	var user models.User

	result := a.db.First(&user, "id = ?", userId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.BadRequest(c, "user not found")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	if !user.IsConfirmed() {
		response.BadRequest(c, "user not confirmed")
		return
	}

	response.WithCustomStatusAndMessage(c, http.StatusOK, gin.H{
		"error": false,
		"user":  user.SanitizeUser(),
	})
}
