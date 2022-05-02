package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
)

func (a *AuthController) SignOut(c *gin.Context) {
	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "refresh token required",
		})
		return
	}

	// delete refresh token
	result := a.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	fmt.Println(result.Error)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "failed to sign out",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "successfully signed out",
	})
}
