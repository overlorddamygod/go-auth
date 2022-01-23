package utils

import "net/mail"

func IsEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
