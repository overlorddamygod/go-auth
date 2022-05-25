package models

import (
	"gorm.io/gorm"
)

type EventType string

const (
	SIGNUP                 EventType = "SIGNUP"
	SIGNIN_EMAIL           EventType = "SIGNIN_EMAIL"
	SIGNIN_GITHUB          EventType = "SIGNIN_GITHUB"
	SIGNIN_MAGICLINK       EventType = "SIGNIN_MAGICLINK"
	SIGNOUT                EventType = "SIGNOUT"
	MAIL_CONFIRMATION_SENT EventType = "MAIL_CONFIRMATION_SENT"
	MAIL_CONFIRMED         EventType = "MAIL_CONFIRMED"
	PASSWORD_RESET_REQUEST EventType = "PASSWORD_RESET_REQUEST"
	PASSWORD_RESET         EventType = "PASSWORD_RESET"
	TOKEN_REFRESH          EventType = "TOKEN_REFRESH"
)

var Events = map[EventType]string{
	"SIGNUP":                 "SIGNUP",
	"SIGNIN_EMAIL":           "SIGNIN_EMAIL",
	"SIGNIN_MAGICLINK":       "SIGNIN_MAGICLINK",
	"SIGNOUT":                "SIGNOUT",
	"MAIL_CONFIRMATION_SENT": "MAIL_CONFIRMATION_SENT",
	"MAIL_CONFIRMED":         "MAIL_CONFIRMED",
	"PASSWORD_RESET_REQUEST": "PASSWORD_RESET_REQUEST",
	"PASSWORD_RESET":         "PASSWORD_RESET",
	"TOKEN_REFRESH":          "TOKEN_REFRESH",
}

type Log struct {
	EventType EventType
	Email     string
	MetaData  JSONMap
	Basic
}

type Logger struct {
	db *gorm.DB
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

func (l *Logger) Log(eventType EventType, email string) (tx *gorm.DB) {
	log := &Log{
		EventType: eventType,
		Email:     email,
	}
	return l.db.Create(log)
}
