package models

import (
	"errors"
	"strings"
	"time"

	"github.com/overlorddamygod/go-auth/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                 string `validate:"required,min=3" binding:"required"`
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=6,max=20"`
	PasswordResetToken   string
	PasswordResetTokenAt time.Time
	ConfirmationToken    string
	Confirmed            bool `gorm:"default:false"`
	ConfirmedAt          time.Time
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
	u.ConfirmationToken, err = utils.GenerateRandomString(15)

	if err != nil {
		return errors.New("server error")
	}
	// validate := validator.New()
	// // validator.Validate(u)
	// err := validate.Struct(u)

	// if err != nil {
	// 	return err
	// }

	result := tx.First(&User{}, "email = ?", u.Email)

	if result.Error == nil {
		return errors.New("email already exists")
	}

	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
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

func (u *User) ConfirmAccount(db *gorm.DB) error {
	if u.Confirmed {
		return errors.New("account already confirmed")
	}
	u.ConfirmationToken = ""
	u.Confirmed = true
	u.ConfirmedAt = time.Now()
	return db.Save(u).Error
}

func NewUser(name string, email string, password string) User {
	return User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

type SanitizedUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) SanitizeUser() SanitizedUser {
	return SanitizedUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
