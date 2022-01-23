package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RefreshToken{})

	r := gin.Default()

	auth := r.Group("/api/auth")

	auth.POST("/signup", func(c *gin.Context) {
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
	})
	auth.POST("/signin", func(c *gin.Context) {
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
	})
	auth.POST("/refresh", func(c *gin.Context) {
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
	})
	auth.GET("/verify", func(c *gin.Context) {
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
	})

	r.Use(func(c *gin.Context) {
		accesstoken := c.GetHeader("X-Access-Token")

		if accesstoken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "access token is required",
			})
			return
		}

		_, err := utils.JwtAccessTokenVerify(accesstoken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "access token is invalid",
			})
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "welcome to api",
		})
	})

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	if err = r.Run(":" + PORT); err != nil {
		log.Fatal(err)
	} else {
		log.Println("server started on port " + PORT)
	}
}
