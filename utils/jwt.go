package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/overlorddamygod/go-auth/configs"
)

type CustomClaims struct {
	IdentityType string    `json:"identity_type"`
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Roles        []int     `json:"roles"`
	jwt.StandardClaims
}

func jwtSign(claims CustomClaims, secret interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func jwtVerify(token string, secret interface{}) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func JwtAccessToken(claims CustomClaims) (string, error) {
	accessJwt := configs.GetConfig().AccessJwt

	claims.StandardClaims.ExpiresAt = time.Now().Add(accessJwt.Expiration).Unix()
	return jwtSign(claims, accessJwt.Secret)
}
func JwtRefreshToken(claims CustomClaims) (string, error) {
	refreshJwt := configs.GetConfig().RefreshJwt

	claims.StandardClaims.ExpiresAt = time.Now().Add(refreshJwt.Expiration).Unix()
	return jwtSign(claims, refreshJwt.Secret)
}

func JwtAccessTokenVerify(token string) (*jwt.Token, error) {
	accessJwt := configs.GetConfig().AccessJwt

	return jwtVerify(token, accessJwt.Secret)
}

func JwtRefreshTokenVerify(token string) (*jwt.Token, error) {
	refreshJwt := configs.GetConfig().RefreshJwt
	return jwtVerify(token, refreshJwt.Secret)
}
