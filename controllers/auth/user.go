package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
)

func (a *AuthController) GetMe(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "user id is required",
		})
		return
	}

	var user models.User

	result := a.db.First(&user, "id = ?", userId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user.SanitizeUser(),
	})
}
