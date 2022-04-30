package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
)

type AuthController struct{}

var userModel = new(models.User)

func (a AuthController) SignUp(c *gin.Context) {
	var db = db.GetDB()

	var user models.User
	c.Bind(&user)
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"user":  user.SanitizeUser(),
	})
}

func (a AuthController) SignIn(c *gin.Context) {
	var db = db.GetDB()
	var user models.User
	c.Bind(&user)

	var dbUser models.User

	result := db.First(&dbUser, "email = ?", user.Email)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "email doesnot exist",
		})
		return
	}

	if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid password",
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

	result = db.Create(&models.RefreshToken{
		Token:  refreshToken,
		UserID: dbUser.ID,
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

func (a AuthController) RefreshToken(c *gin.Context) {
	var db = db.GetDB()

	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token is required",
		})
		return
	}

	token, err := utils.JwtRefreshTokenVerify(refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token is invalid",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token is invalid",
		})
		return
	}

	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)

	var refreshTokenModel models.RefreshToken
	result := db.First(&refreshTokenModel, "token = ?", refreshToken)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token is invalid",
		})
		return
	}

	if refreshTokenModel.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token is revoked",
		})
		return
	}

	accessToken, aTerr := utils.JwtAccessToken(utils.CustomClaims{
		UserID: userID,
		Email:  email,
	})

	if aTerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "failed to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":        false,
		"access-token": accessToken,
	})
}

func (a AuthController) VerifyLogin(c *gin.Context) {
	accesstoken := c.GetHeader("X-Access-Token")

	if accesstoken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "access token is required",
		})
		return
	}

	_, err := utils.JwtAccessTokenVerify(accesstoken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "access token is invalid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a AuthController) RequestPasswordRecovery(c *gin.Context) {
	var db = db.GetDB()
	var user models.User
	c.Bind(&user)

	if strings.TrimSpace(user.Email) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "email address is required",
		})
		return
	}

	var dbUser models.User

	result := db.First(&dbUser, "email = ?", user.Email)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "email doesnot exist",
		})
		return
	}

	resetCode, err := dbUser.GeneratePasswordRecoveryToken(db)

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

func (a AuthController) PasswordReset(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
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

	var db = db.GetDB()

	var dbUser models.User

	result := db.First(&dbUser, "password_reset_token = ?", token)

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
	var user models.User
	c.Bind(&user)

	err = dbUser.ResetPasswordWithToken(db, user.Password)

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

func (a AuthController) GetMe(c *gin.Context) {

	var db = db.GetDB()
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "user id is required",
		})
		return
	}

	var user models.User

	result := db.First(&user, "id = ?", userId)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
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
