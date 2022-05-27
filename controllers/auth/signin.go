package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

type SignInParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

func (a *AuthController) SignIn(c *gin.Context) {
	var params SignInParams
	if err := c.Bind(&params); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	loginType := c.Query("type")

	if params.Email == "" {
		response.BadRequest(c, "email is required")
		return
	}

	if loginType == "email" && params.Password == "" {
		response.BadRequest(c, "password is required")
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "email = ? AND identity_type = ?", params.Email, "email")

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.NotFound(c, "email doesnot exist")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	if !dbUser.IsConfirmed() {
		response.Unauthorized(c, "user not confirmed")
		return
	}

	switch loginType {
	case "email":
		res, code, err := dbUser.SignInWithEmail(params.Password, a.db, c)

		if err != nil {
			response.WithCustomStatusAndMessage(c, code, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		result = a.logger.Log(models.SIGNIN_EMAIL, dbUser.Email)

		if result.Error != nil {
			fmt.Println("Error Logging: ", models.SIGNIN_EMAIL, result.Error)
		}

		c.JSON(http.StatusOK, res)
	case "magiclink":
		res, code, err := dbUser.GenerateMagicLink(c, a.db, a.mailer)

		if err != nil {
			response.WithCustomStatusAndMessage(c, code, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	default:
		response.BadRequest(c, "invalid login type")
	}
}

func (a *AuthController) VerifyLogin(c *gin.Context) {
	loginType := c.Query("type")
	token := c.Query("token")
	redirectTo := c.Query("redirect_to")

	if token == "" {
		response.BadRequest(c, "token required")
		return
	}

	switch loginType {
	case "magiclink":
		if redirectTo == "" {
			response.BadRequest(c, "redirect url required")
		}
		token, err := utils.Decrypt(token)

		if err != nil {
			response.Unauthorized(c, "invalid token")
			return
		}
		var dbUser models.User
		result := a.db.First(&dbUser, "token = ? AND identity_type = ?", token, "email")

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				response.Unauthorized(c, "invalid token")
				return
			}
			response.ServerError(c, "server error")
			return
		}

		if dbUser.IsConfirmed() {
			if time.Since(dbUser.TokenSentAt).Hours() > 1 {
				response.Unauthorized(c, "token expired")
				return
			}
			tokenMap, err := dbUser.GenerateAccessRefreshToken(c, a.db)

			if err != nil {
				response.ServerError(c, "server error")
				return
			}

			dbUser.Token = ""
			dbUser.TokenSentAt = time.Time{}
			result := a.db.Save(&dbUser)

			if result.Error != nil {
				response.ServerError(c, "server error")
				return
			}

			result = a.logger.Log(models.SIGNIN_MAGICLINK, dbUser.Email)

			if result.Error != nil {
				fmt.Println("Error Logging: ", models.SIGNIN_MAGICLINK, result.Error)
			}

			redirectTo = fmt.Sprintf("%s?type=magiclink&access_token=%s&refresh_token=%s", redirectTo, tokenMap["accessToken"], tokenMap["refreshToken"])
			c.Redirect(http.StatusFound, redirectTo)
		}
		return
	default:
		response.BadRequest(c, "invalid login type")
		return
	}
}
