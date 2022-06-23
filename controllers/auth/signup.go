package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils/response"
)

type SignUpParams struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

func (a *AuthController) SignUp(c *gin.Context) {
	var params SignUpParams
	if err := c.Bind(&params); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	user, err := models.NewUser(params.Name, params.Email, params.Password)

	if err != nil {
		response.ServerError(c, "server error")
		return
	}

	result := a.db.Create(&user)

	if result.Error != nil {
		err, ok := result.Error.(*pgconn.PgError)

		if ok {
			if err.Code == "23505" {
				response.BadRequest(c, "email already exists")
				return
			}
		}

		response.BadRequest(c, result.Error.Error())
		return
	}

	if a.config.RequireEmailConfirmation {
		fmt.Println("Confirmation Token: ", user.ConfirmationToken)
		err := a.mailer.SendConfirmationMail(user.Email, user.Name, a.config.ApiUrl+"/api/v1/auth/confirm?token="+user.ConfirmationToken)
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

	if a.config.RequireEmailConfirmation {
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
