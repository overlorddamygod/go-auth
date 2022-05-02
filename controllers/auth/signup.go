package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
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
