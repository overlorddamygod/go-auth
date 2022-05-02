package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	jwt_access_token_secret  = []byte("my_access_token_secret_key")
	jwt_refresh_token_secret = []byte("my_refresh_token_secret_key")

	jwt_access_token_expiration  = time.Hour * 1
	jwt_refresh_token_expiration = time.Hour * 24
)

type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
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
	claims.StandardClaims.ExpiresAt = time.Now().Add(jwt_access_token_expiration).Unix()
	return jwtSign(claims, jwt_access_token_secret)
}
func JwtRefreshToken(claims CustomClaims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(jwt_refresh_token_expiration).Unix()
	return jwtSign(claims, jwt_refresh_token_secret)
}

func JwtAccessTokenVerify(token string) (*jwt.Token, error) {
	return jwtVerify(token, jwt_access_token_secret)
}

func JwtRefreshTokenVerify(token string) (*jwt.Token, error) {
	return jwtVerify(token, jwt_refresh_token_secret)
}
