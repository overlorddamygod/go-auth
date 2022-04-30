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
	Name                 string `validate:"required,min=3"`
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=6,max=20"`
	PasswordResetToken   string
	PasswordResetTokenAt time.Time
	Confirmed            bool `gorm:"default:false"`
	ConfirmedAt          time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(strings.TrimSpace(u.Name)) < 3 {
		return errors.New("name must be at least 3 characters")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if !utils.IsEmailValid(u.Email) {
		return errors.New("invalid email")
	}
	// validate := validator.New()
	// // validator.Validate(u)
	// err := validate.Struct(u)

	// if err != nil {
	// 	return err
	// }

	result := tx.First(&User{}, "email = ?", u.Email)

	if result.Error == nil {
		return errors.New("user of that email already exist")
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
	result := db.Save(u)

	if result.Error != nil {
		return errors.New("error saving to the db")
	}

	return nil
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
