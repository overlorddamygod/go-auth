package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

type RecoveryParams struct {
	Email string `json:"email" binding:"required"`
}

func (a *AuthController) RequestPasswordRecovery(c *gin.Context) {
	var params RecoveryParams
	c.Bind(&params)

	if strings.TrimSpace(params.Email) == "" {
		response.BadRequest(c, "email address required")
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "email = ?", params.Email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.NotFound(c, "email address not found")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	resetCode, err := dbUser.GeneratePasswordRecoveryToken(a.db)

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	a.mailer.SendPasswordRecoveryMail(dbUser.Email, dbUser.Name, resetCode)
	response.Ok(c, "recovery email sent")
}

func (a *AuthController) PasswordReset(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "invalid token")
		return
	}
	token, err := utils.Decrypt(token)

	if err != nil {
		response.Unauthorized(c, "invalid token")
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "password_reset_token = ?", token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.Unauthorized(c, "invalid token")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	// check if reset token is between 1 day
	if time.Since(dbUser.PasswordResetTokenAt).Hours() > 24 {
		response.Unauthorized(c, "token expired")
		return
	}

	// get password from body
	var params SignInParams
	c.Bind(&params)

	err = dbUser.ResetPasswordWithToken(a.db, params.Password)

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Ok(c, "password reset successfully")
}
