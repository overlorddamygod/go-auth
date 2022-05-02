package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
)

type SignInParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) SignIn(c *gin.Context) {
	var params SignInParams
	c.Bind(&params)

	if params.Email == "" || params.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "email and password are required",
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

	if !utils.CheckPasswordHash(params.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid password",
		})
		return
	}

	// check if user is confirmed
	if !dbUser.Confirmed {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "user is not confirmed",
		})
		return
	}

	accessToken, aTerr := utils.JwtAccessToken(utils.CustomClaims{
		UserID: dbUser.ID,
		Email:  dbUser.Email,
	})

	if aTerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "failed to sign in",
		})
		return
	}

	refreshToken, rTerr := utils.JwtRefreshToken(utils.CustomClaims{
		UserID: dbUser.ID,
		Email:  dbUser.Email,
	})

	if rTerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "failed to sign in",
		})
		return
	}
	// get user agent header
	userAgent := c.GetHeader("User-Agent")

	// get user ip
	ip := c.ClientIP()

	result = a.db.Create(&models.RefreshToken{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		UserAgent: userAgent,
		IP:        ip,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":         false,
		"access-token":  accessToken,
		"refresh-token": refreshToken,
	})
}

func (a *AuthController) VerifyLogin(c *gin.Context) {
	accesstoken := c.GetHeader("X-Access-Token")

	if accesstoken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "access token required",
		})
		return
	}

	_, err := utils.JwtAccessTokenVerify(accesstoken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "access token invalid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
