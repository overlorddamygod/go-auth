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

	config = Config{
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
	}
}

func LoadConfig() {
	Load(".env")
}

func GetConfig() Config {
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Println("Error loading", key, "... Using default")

	return fallback
}

func loadJWTConfig(prefix string) (c JwtConfig, e error) {
	secret := getEnv(prefix+"_SECRET", jwtMap[prefix+"_SECRET"])
	expiration, err := strconv.Atoi(getEnv(prefix+"_EXPIRATION_HOURS", jwtMap[prefix+"_EXPIRATION_HOURS"]))
	c = JwtConfig{
		Secret: []byte(secret),
	}
	e = err
	c.Expiration = time.Hour * time.Duration(expiration)
	return c, e
}

var jwtMap = map[string]string{
	"JWT_ACCESS_SECRET":            "my_access_token_secret_key",
	"JWT_ACCESS_EXPIRATION_HOURS":  "1",
	"JWT_REFRESH_SECRET":           "my_refresh_token_secret_key",
	"JWT_REFRESH_EXPIRATION_HOURS": "24",
}

var defaultConfig = Config{
	RequireEmailConfirmation: false,
	AccessJwt: JwtConfig{
		Secret:     []byte(jwtMap["JWT_ACCESS_SECRET"]),
		Expiration: time.Hour * time.Duration(1),
	},
	RefreshJwt: JwtConfig{
		Secret:     []byte(jwtMap["JWT_REFRESH_SECRET"]),
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
}

func (d DBConfig) GetDialector() gorm.Dialector {
	return postgres.Open(d.PostgresDSN)
}
