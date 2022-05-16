package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type JwtConfig struct {
	Secret     []byte
	Expiration time.Duration
}

type Config struct {
	RequireEmailConfirmation bool
	Database                 DBConfig
	AccessJwt                JwtConfig
	RefreshJwt               JwtConfig
	Mail                     SMTPConfig
}

type DBConfig struct {
	PostgresDSN string
}

var config Config

func Load(envPath string) {
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	access, err := loadJWTConfig("JWT_ACCESS")

	if err != nil {
		log.Fatalf("Error loading config")
	}

	refresh, err := loadJWTConfig("JWT_REFRESH")

	if err != nil {
		log.Fatalf("Error loading config")
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if err != nil {
		log.Println("SMTP_PORT not set, using default")
		smtpPort = 587
	}

	config = Config{
		RequireEmailConfirmation: false,
		AccessJwt:                access,
		RefreshJwt:               refresh,
		Database: DBConfig{
			PostgresDSN: os.Getenv("POSTGRES_DSN"),
		},
		Mail: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}
}

func LoadConfig() {
	Load(".env")
}

func GetConfig() Config {
	return config
}

func loadJWTConfig(prefix string) (c JwtConfig, e error) {
	secret := os.Getenv(prefix + "_SECRET")
	expiration, err := strconv.Atoi(os.Getenv(prefix + "_EXPIRATION_HOURS"))
	c = JwtConfig{
		Secret: []byte(secret),
	}
	e = err
	c.Expiration = time.Hour * time.Duration(expiration)
	return c, e
}

func (d DBConfig) GetDialector() gorm.Dialector {
	return postgres.Open(d.PostgresDSN)
}
