package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

type OauthProvider interface {
	Do(code string) (string, *models.User, error)
	GetOauthUrl(redirect_to string) (string, error)
}

func (a *AuthController) OAuthGithub(c *gin.Context) {
	oauthProvider := c.Query("oauth_provider")
	redirect_to := c.Query("redirect_to")

	if redirect_to == "" {
		response.BadRequest(c, "redirect url required")
		return
	}

	ok := IsProviderValid(a.config, oauthProvider)
	if !ok {
		response.BadRequest(c, "invalid oauth provider")
		return
	}

	var oauthUrl string

	switch oauthProvider {
	case "github":
		githubOauth, err := NewGithubOauth()

		if err != nil {
			response.ServerError(c, "server error")
			return
		}
		oauthUrl, err = githubOauth.GetOauthUrl(redirect_to)

		if err != nil {
			response.ServerError(c, "server error")
			return
		}
	default:
		response.BadRequest(c, "invalid oauth provider")
		return
	}
	c.Redirect(302, oauthUrl)
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (a *AuthController) OAuthAuthorize(c *gin.Context) {
	oauthProvider := c.Query("oauth_provider")
	redirect_to := c.Query("redirect_to")
	code := c.Query("code")

	fmt.Println(c.Request.URL.String())

	ok := IsProviderValid(a.config, oauthProvider)
	if !ok {
		response.BadRequest(c, "invalid oauth provider")
		return
	}

	switch oauthProvider {
	case "github":
		githubOauth, err := NewGithubOauth()

		if err != nil {
			response.ServerError(c, err.Error())
			return
		}

		_, user, err := githubOauth.Do(code)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		var dbUser models.User
		result := a.db.First(&dbUser, "email = ?", user.Email)

		if result.Error == nil {
			identities, ok := dbUser.Identities["identities"].([]interface{})

			if !ok {
				// respond error
				response.ServerError(c, "server error")
				return
			}

			if dbUser.IdentityType == oauthProvider || Contains(identities, oauthProvider) {

				tokenMap, err := dbUser.GenerateAccessRefreshToken(c, a.db)

				if err != nil {
					response.ServerError(c, "server error")
					return
				}
				result = a.logger.Log(models.SIGNIN_GITHUB, dbUser.Email)

				if result.Error != nil {
					fmt.Println("Error Logging: ", models.SIGNIN_GITHUB, result.Error)
				}
				redirect_to = fmt.Sprintf("%s?type=githuboauth&access_token=%s&refresh_token=%s", redirect_to, tokenMap["accessToken"], tokenMap["refreshToken"])
				c.Redirect(http.StatusFound, redirect_to)
				return
			}

			redirect_to = fmt.Sprintf("%s?error=not_found", redirect_to)
			c.Redirect(http.StatusNotFound, redirect_to)
			return
		}

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.BadRequest(c, "server error")
			return
		}

		randomPassword, err := utils.GenerateRandomString(10)

		if err != nil {
			response.ServerError(c, "server error")
			return
		}

		newUser, err := models.NewUser(user.Name, user.Email, randomPassword)

		if err != nil {
			response.ServerError(c, "server error")
			return
		}

		newUser.IdentityType = "github"

		if err := a.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&newUser).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			response.BadRequest(c, "server error")
			return
		}

		tokenMap, err := newUser.GenerateAccessRefreshToken(c, a.db)

		if err != nil {
			response.ServerError(c, "server error")
			return
		}

		result = a.logger.Log(models.SIGNIN_GITHUB, newUser.Email)

		if result.Error != nil {
			fmt.Println("Error Logging: ", models.SIGNIN_GITHUB, result.Error)
		}

		redirect_to = fmt.Sprintf("%s?type=githuboauth&access_token=%s&refresh_token=%s", redirect_to, tokenMap["accessToken"], tokenMap["refreshToken"])

		c.Redirect(http.StatusFound, redirect_to)
	default:
		response.BadRequest(c, "invalid oauth provider")
		return
	}
}

func IsProviderValid(config *configs.Config, providerName string) bool {
	provider, ok := config.Oauth[providerName]

	if !ok || !provider.AllowLogin {
		return false
	}
	return true
}

func Contains(slice []interface{}, match string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == match {
			return true
		}
	}
	return false
}
