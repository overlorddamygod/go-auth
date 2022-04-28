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
	Name                string
	Email               string
	Password            string
	PasswordResetCode   string
	PasswordResetCodeAt time.Time
	Confirmed           bool `gorm:"default:false"`
	ConfirmedAt         time.Time
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

	result := tx.First(&User{}, "email = ?", u.Email)

	if result.Error == nil {
		return errors.New("user of that email already exist")
	}

	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return
}

func (u *User) GeneratePasswordRecoveryCode(db *gorm.DB) (code string, err error) {
	randomString, err := utils.GenerateRandomString(12)

	if err != nil {
		return "", errors.New("error while password recovery")
	}

	u.PasswordResetCode = randomString
	u.PasswordResetCodeAt = time.Now()
	result := db.Save(u)

	if result.Error != nil {
		return "", errors.New("error saving to the db")
	}
	return u.PasswordResetCode, nil
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
