package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils/response"
)

type SignUpParams struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) SignUp(c *gin.Context) {
	var params SignUpParams
	c.Bind(&params)

	var user models.User = models.NewUser(params.Name, params.Email, params.Password)

	result := a.db.Create(&user)

	if result.Error != nil {
		response.BadRequest(c, result.Error.Error())
		return
	}
	if configs.GetConfig().RequireEmailConfirmation {
		fmt.Println("Confirmation Token: ", user.ConfirmationToken)
		err := a.mailer.SendConfirmationMail(user.Email, user.Name, "http://localhost:8080/api/v1/auth/confirm?token="+user.ConfirmationToken)
		fmt.Println("MAIL: ", err)

		result = a.logger.Log(models.MAIL_CONFIRMATION_SENT, user.Email)

		if result.Error != nil {
			fmt.Println("Error Logging: ", models.MAIL_CONFIRMATION_SENT, result.Error)
		}
	}

	result = a.logger.Log(models.SIGNUP, user.Email)

	if result.Error != nil {
		fmt.Println("Error Logging: ", models.SIGNUP, result.Error)
	}

	msg := "account created"

	if configs.GetConfig().RequireEmailConfirmation {
		msg = "account created, please check your email"
	}

	response.Created(c, msg)
}

// confirm account
func (a *AuthController) ConfirmAccount(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		response.BadRequest(c, "token required")
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "confirmation_token = ?", token)

	if result.Error != nil {
		response.Unauthorized(c, "invalid token")
		return
	}

	if err := dbUser.ConfirmAccount(a.db); err != nil {
		response.Unauthorized(c, "failed to confirm account")
		return
	}

	result = a.logger.Log(models.MAIL_CONFIRMED, dbUser.Email)

	if result.Error != nil {
		fmt.Println("Error Logging: ", models.MAIL_CONFIRMED, result.Error)
	}

	response.Ok(c, "account confirmed")
}
