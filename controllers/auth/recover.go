package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
)

type RecoveryParams struct {
	Email string `json:"email" binding:"required"`
}

func (a *AuthController) RequestPasswordRecovery(c *gin.Context) {
	var params RecoveryParams
	c.Bind(&params)

	if strings.TrimSpace(params.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "email address required",
		})
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "email = ?", params.Email)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "email doesnot exist",
		})
		return
	}

	resetCode, err := dbUser.GeneratePasswordRecoveryToken(a.db)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": false,
		"code":  resetCode,
	})
}

func (a *AuthController) PasswordReset(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "invalid token",
		})
		return
	}
	token, err := utils.Decrypt(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid token",
		})
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "password_reset_token = ?", token)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid token",
		})
		return
	}

	// check if reset token is between 1 day
	if time.Since(dbUser.PasswordResetTokenAt).Hours() > 24 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "token expired",
		})
		return
	}

	// get password from body
	var params SignInParams
	c.Bind(&params)

	err = dbUser.ResetPasswordWithToken(a.db, params.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
