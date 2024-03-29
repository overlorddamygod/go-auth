package models

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/utils"
	"gorm.io/gorm"
)

type User struct {
	Basic
	Name     string `validate:"required,min=3" binding:"required"`
	Email    string `gorm:"index:email,unique" validate:"required,email"`
	Password string `validate:"required,min=6,max=20"`

	IdentityType string `gorm:"default:'email'"`
	Identities   JSONMap

	PasswordResetToken   string
	PasswordResetTokenAt time.Time

	// magic link token
	Token       string
	TokenSentAt time.Time

	ConfirmationToken   string
	ConfirmationTokenAt time.Time
	Confirmed           bool `gorm:"default:false"`
	ConfirmedAt         time.Time
	RefreshToken        []RefreshToken `gorm:"one2many;constraint:OnDelete:CASCADE"`
	Roles               []UserRole     `gorm:"one2many:user_roles"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if len(strings.TrimSpace(u.Name)) < 3 {
		return errors.New("name must be at least 3 characters")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !utils.IsEmailValid(u.Email) {
		return errors.New("invalid email")
	}
	if !configs.GetConfig().RequireEmailConfirmation {
		u.Confirmed = true
		u.ConfirmedAt = time.Now()
	} else {
		u.ConfirmationToken, err = utils.GenerateRandomString(15)
		u.ConfirmationTokenAt = time.Now()

		if err != nil {
			return errors.New("server error")
		}
	}

	return nil
}

func (u *User) GetUserByEmail(email string, db *gorm.DB) (tx *gorm.DB) {
	return db.First(u, "email = ?", email)
}

func (u *User) GeneratePasswordRecoveryToken(db *gorm.DB) (token string, err error) {
	randomString, err := utils.GenerateRandomString(9)

	if err != nil {
		return "", errors.New("error while password recovery")
	}

	u.PasswordResetToken = randomString

	encryptedToken, err := utils.Encrypt(randomString)

	if err != nil {
		return "", errors.New("error while password recovery")
	}

	u.PasswordResetTokenAt = time.Now()
	result := db.Save(u)

	if result.Error != nil {
		return "", errors.New("error saving to the db")
	}
	return encryptedToken, nil
}

func (u *User) ResetPasswordWithToken(db *gorm.DB, password string) (err error) {
	u.Password, err = utils.HashPassword(password)

	if err != nil {
		return err
	}

	u.PasswordResetToken = ""
	u.PasswordResetTokenAt = time.Time{}
	// result := .Save(u)

	// using transaction delete all refresh token and then save user to db
	err = db.Transaction(func(tx *gorm.DB) error {
		var refreshTokens []RefreshToken

		if err := tx.Save(u).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", u.ID).Delete(&refreshTokens).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.New("error saving to the db")
	}

	return nil
}

func (u *User) IsConfirmed() bool {
	return u.Confirmed
}

func (u *User) ConfirmAccount(db *gorm.DB) error {
	if u.IsConfirmed() {
		return errors.New("account already confirmed")
	}
	u.ConfirmationToken = ""
	u.Confirmed = true
	u.ConfirmedAt = time.Now()
	return db.Save(u).Error
}

func (u *User) GetRoles() []int {
	var roles []int

	for _, role := range u.Roles {
		roles = append(roles, role.Type)
	}

	return roles
}

func (u *User) GetAccessToken() (string, error) {
	return utils.JwtAccessToken(utils.CustomClaims{
		IdentityType: u.IdentityType,
		UserID:       u.ID,
		Email:        u.Email,
		Roles:        u.GetRoles(),
	})
}
func (u *User) GetRefreshToken() (string, error) {
	return utils.JwtRefreshToken(utils.CustomClaims{
		IdentityType: u.IdentityType,
		UserID:       u.ID,
		Email:        u.Email,
	})
}

func (u *User) GenerateAccessRefreshToken(c *gin.Context, db *gorm.DB) (tokenMap map[string]string, err error) {
	tokenMap = make(map[string]string)

	accessToken, aTerr := u.GetAccessToken()

	if aTerr != nil {
		return nil, errors.New("failed to sign in")
	}

	refreshToken, rTerr := u.GetRefreshToken()

	if rTerr != nil {
		return nil, errors.New("failed to sign in")
	}

	userAgent := c.GetHeader("User-Agent")

	// get user ip
	ip := c.ClientIP()

	// create refresh token
	refreshTokenModel := RefreshToken{
		Token:     refreshToken,
		UserID:    u.ID,
		UserAgent: userAgent,
		IP:        ip,
	}
	result := db.Create(&refreshTokenModel)

	if result.Error != nil {
		return nil, result.Error
	}

	tokenMap["accessToken"] = accessToken
	tokenMap["refreshToken"] = refreshToken

	return tokenMap, nil
}

func (u *User) SignInWithEmail(password string, db *gorm.DB, c *gin.Context) (obj interface{}, code int, err error) {
	fmt.Println(utils.CheckPasswordHash(password, u.Password), password, u.Password)

	// sign in with email
	if !utils.CheckPasswordHash(password, u.Password) {
		return nil, http.StatusUnauthorized, errors.New("invalid password")
	}
	fmt.Println(password)

	if !u.IsConfirmed() {
		return nil, http.StatusUnauthorized, errors.New("user is not confirmed")
	}

	token, err := u.GenerateAccessRefreshToken(c, db)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to sign in")
	}

	return gin.H{
		"error":         false,
		"access_token":  token["accessToken"],
		"refresh_token": token["refreshToken"],
	}, http.StatusOK, nil
}

func (u *User) GenerateMagicLink(c *gin.Context, db *gorm.DB, mailer *mailer.Mailer) (obj interface{}, code int, err error) {
	redirect_to := c.Query("redirect_to")

	if redirect_to == "" {
		return nil, http.StatusBadRequest, errors.New("redirect url is required")
	}

	// sign in with magic link
	randomString, err := utils.GenerateRandomString(9)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	u.Token = randomString

	encryptedToken, err := utils.Encrypt(randomString)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	u.TokenSentAt = time.Now()

	magicLink := fmt.Sprintf("%s/api/v1/auth/verify?type=%s&token=%s&redirect_to=%s", configs.MainConfig.ApiUrl, "magiclink", encryptedToken, redirect_to)
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(u).Error; err != nil {
			return err
		}

		println(magicLink)

		if err := mailer.SendMagicLink(u.Email, u.Name, magicLink); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	return gin.H{
		"error":   false,
		"message": "Succesfully sent magic link to your email",
	}, http.StatusOK, nil
}

func NewUser(name string, email string, password string) (User, error) {
	if len(password) < 6 {
		return User{}, errors.New("password must be at least 6 characters")
	}

	pass, err := utils.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	return User{
		Name:     name,
		Email:    email,
		Password: pass,
	}, nil
}

type SanitizedUser struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (u *User) SanitizeUser() SanitizedUser {
	return SanitizedUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
