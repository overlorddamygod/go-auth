package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	db     *gorm.DB
	mailer *mailer.Mailer
}

func NewAuthController(db *gorm.DB, mailer *mailer.Mailer) AuthController {
	return AuthController{
		db:     db,
		mailer: mailer,
	}
}

// var userModel = new(models.User)

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": result.Error.Error(),
		})
		return
	}
	err := a.mailer.SendConfirmationMail(user.Email, user.Name, "http://localhost:8080/api/v1/auth/confirm?token="+user.ConfirmationToken)
	fmt.Println(err)

	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"user":  user.SanitizeUser(),
	})
}

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

func (a *AuthController) RefreshToken(c *gin.Context) {
	var refreshToken string = c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "refresh token required",
		})
		return
	}

	token, err := utils.JwtRefreshTokenVerify(refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token invalid",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token invalid",
		})
		return
	}

	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)

	var refreshTokenModel models.RefreshToken
	result := a.db.First(&refreshTokenModel, "token = ?", refreshToken)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token expired",
		})
		return
	}

	if refreshTokenModel.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "refresh token revoked",
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

func (a *AuthController) GetMe(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "user id is required",
		})
		return
	}

	var user models.User

	result := a.db.First(&user, "id = ?", userId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
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

// confirm account
func (a *AuthController) ConfirmAccount(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "token required",
		})
		return
	}

	var dbUser models.User

	result := a.db.First(&dbUser, "confirmation_token = ?", token)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid token",
		})
		return
	}

	if err := dbUser.ConfirmAccount(a.db); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "failed to confirm account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
