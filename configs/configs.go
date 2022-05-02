package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}

type JwtConfig struct {
	Secret     string
	Expiration time.Duration
}

type Config struct {
	RequireConfirmation bool
	AccessJwt           JwtConfig
	RefreshJwt          JwtConfig
	Mail                SMTP
}

var config Config

func LoadConfig() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	access, _ := GetJWTConfig("JWT_ACCESS")
	refresh, _ := GetJWTConfig("JWT_REFRESH")

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	config = Config{
		RequireConfirmation: true,
		AccessJwt:           access,
		RefreshJwt:          refresh,
		Mail: SMTP{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}
	fmt.Println(config)
}

func GetConfig() Config {
	return config
}

func GetJWTConfig(prefix string) (c JwtConfig, e error) {
	secret := os.Getenv(prefix + "_SECRET")
	expiration, err := strconv.Atoi(os.Getenv(prefix + "_EXPIRATION_HOURS"))
	c = JwtConfig{
		Secret: secret,
	}
	e = err
	c.Expiration = time.Hour * time.Duration(expiration)
	return c, e
}
