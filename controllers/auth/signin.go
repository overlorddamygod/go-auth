package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"gorm.io/gorm"
)

type SignInParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

func (a *AuthController) SignIn(c *gin.Context) {
	var params SignInParams
	c.Bind(&params)

	loginType := c.Query("type")

	if params.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "email is required",
		})
		return
	}

	if loginType == "email" && params.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "password is required",
		})
	}

	var dbUser models.User
	result := dbUser.GetUserByEmail(params.Email, a.db)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "email doesnot exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "server error",
		})
		return
	}

	switch loginType {
	case "email":
		res, code, err := dbUser.SignInWithEmail(params.Password, a.db, c)

		if err != nil {
			c.JSON(code, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, res)
	case "magiclink":
		res, code, err := dbUser.GenerateMagicLink(c, a.db, a.mailer)

		if err != nil {
			c.JSON(code, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, res)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "invalid login type",
		})
	}
}

func (a *AuthController) VerifyLogin(c *gin.Context) {
	loginType := c.Query("type")
	token := c.Query("token")
	redirectTo := c.Query("redirect_to")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "token required",
		})
		return
	}

	switch loginType {
	case "magiclink":
		token, err := utils.Decrypt(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "invalid token",
			})
			return
		}
		var dbUser models.User
		result := a.db.First(&dbUser, "token = ?", token)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "invalid token",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "server error",
			})
			return
		}

		if dbUser.IsConfirmed() {
			if time.Since(dbUser.TokenSentAt).Hours() > 1 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "token expired",
				})
				return
			}
			tokenMap, err := dbUser.GenerateAccessRefreshToken(c, a.db)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": "server error",
				})
				return
			}

			if redirectTo != "" {
				dbUser.Token = ""
				dbUser.TokenSentAt = time.Time{}
				result := a.db.Save(&dbUser)

				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   true,
						"message": "server error",
					})
					return
				}

				redirectTo = fmt.Sprintf("%s?type=magiclink&access_token=%s&refresh_token=%s", redirectTo, tokenMap["accessToken"], tokenMap["refreshToken"])
				c.Redirect(http.StatusFound, redirectTo)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "redirect url not defined",
				})
				return
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "invalid login type",
		})
		return
	}

	// accesstoken := c.GetHeader("X-Access-Token")

	// if accesstoken == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error":   true,
	// 		"message": "access token required",
	// 	})
	// 	return
	// }

	// _, err := utils.JwtAccessTokenVerify(accesstoken)

	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error":   true,
	// 		"message": "access token invalid",
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"error": false,
	// })
}
