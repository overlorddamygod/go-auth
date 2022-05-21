package configs

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
	PORT                     string
	RequireEmailConfirmation bool
	Database                 DBConfig
	AccessJwt                JwtConfig
	RefreshJwt               JwtConfig
	Mail                     SMTPConfig
	TokenSecret1             []byte
	TokenSecret2             []byte
	AllowOrigins             []string
}

type DBConfig struct {
	PostgresDSN string
}

var MainConfig *Config

func NewConfig(envPath string) func() *Config {
	return func() *Config {
		err := godotenv.Load(envPath)

		if err != nil {
			log.Println("Error loading .env file")
		}

		access, err := loadJWTConfig("JWT_ACCESS")

		if err != nil {
			log.Println("Error loading access jwt config using default")
		}

		refresh, err := loadJWTConfig("JWT_REFRESH")

		if err != nil {
			log.Println("Error loading refresh jwt config using default")
		}

		smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

		if err != nil {
			log.Println("SMTP_PORT not set, using default")
		}

		allowOrigins := getEnv("ALLOW_ORIGINS", configMap["AllowOrigins"])
		allowOriginsArray := strings.Split(allowOrigins, " ")

		config := &Config{
			PORT:                     getEnv("PORT", defaultConfig.PORT),
			RequireEmailConfirmation: getEnv("MAIL_CONFIRMATION", "0") == "1",
			AccessJwt:                access,
			RefreshJwt:               refresh,
			Database: DBConfig{
				PostgresDSN: getEnv("POSTGRES_DSN", defaultConfig.Database.PostgresDSN),
			},
			Mail: SMTPConfig{
				Host:     getEnv("SMTP_HOST", defaultConfig.Mail.Host),
				Port:     smtpPort,
				Username: getEnv("SMTP_USERNAME", defaultConfig.Mail.Username),
				Password: getEnv("SMTP_PASSWORD", defaultConfig.Mail.Password),
			},
			TokenSecret1: []byte(getEnv("TOKEN_SECRET1", configMap["TokenSecret1"])),
			TokenSecret2: []byte(getEnv("TOKEN_SECRET2", configMap["TokenSecret2"])),
			AllowOrigins: allowOriginsArray,
		}
		// fmt.Println(config)
		return config
	}
}

func GetConfig() *Config {
	return MainConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Println("Error loading", key, "... Using default")

	return fallback
}

func loadJWTConfig(prefix string) (c JwtConfig, e error) {
	secret := getEnv(prefix+"_SECRET", configMap[prefix+"_SECRET"])
	expiration, err := strconv.Atoi(getEnv(prefix+"_EXPIRATION_HOURS", configMap[prefix+"_EXPIRATION_HOURS"]))
	c = JwtConfig{
		Secret: []byte(secret),
	}
	e = err
	c.Expiration = time.Hour * time.Duration(expiration)
	return c, e
}

var configMap = map[string]string{
	"JWT_ACCESS_SECRET":            "my_access_token_secret_key",
	"JWT_ACCESS_EXPIRATION_HOURS":  "1",
	"JWT_REFRESH_SECRET":           "my_refresh_token_secret_key",
	"JWT_REFRESH_EXPIRATION_HOURS": "24",
	"TokenSecret1":                 "PQNFuUjXBfOBbDcc8IlJlqL4",
	"TokenSecret2":                 "Qjb2MwC5aPTA26gc",
	"AllowOrigins":                 "http://localhost:3000",
}

var defaultConfig = Config{
	PORT:                     "8080",
	RequireEmailConfirmation: false,
	AccessJwt: JwtConfig{
		Secret:     []byte(configMap["JWT_ACCESS_SECRET"]),
		Expiration: time.Hour * time.Duration(1),
	},
	RefreshJwt: JwtConfig{
		Secret:     []byte(configMap["JWT_REFRESH_SECRET"]),
		Expiration: time.Hour * time.Duration(24),
	},
	Database: DBConfig{
		PostgresDSN: "user=postgres password=sdfsdfwer3qe3e3edwdeqe host=db.sfasdasdasdwe.domain.co port=5432 dbname=postgres",
	},
	Mail: SMTPConfig{
		Host:     "in-v3.smtphost.com",
		Port:     587,
		Username: "qwr23423twggggdf8",
		Password: "Sasdasdasd",
	},
	TokenSecret1: []byte(configMap["TokenSecret1"]),
	TokenSecret2: []byte(configMap["TokenSecret2"]),
	AllowOrigins: []string{configMap["AllowOrigins"]},
}
