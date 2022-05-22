package auth_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/controllers/auth"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/middlewares"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/server"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type AuthTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	user   User
}

type User struct {
	Name         string
	Email        string
	Password     string
	AccessToken  string
	RefreshToken string
}

func TestAuth(t *testing.T) {
	var dbConn *gorm.DB
	var router *gin.Engine
	_ = fx.New(
		fx.Provide(
			configs.NewConfig("../../.env"),
			db.NewDB,
			mailer.NewMailer,
			models.NewLogger,
			middlewares.NewLimiter,
			auth.NewAuthController,
			server.NewRouter,
		),
		fx.Populate(&dbConn),
		fx.Populate(&router),
		fx.Populate(&configs.MainConfig),
		fx.Invoke(server.RegisterServer),
	)

	s := &AuthTestSuite{
		db:     dbConn,
		router: router,
		user: User{
			Name:         "test123",
			Email:        "test123@gmail.com",
			Password:     "test123",
			AccessToken:  "",
			RefreshToken: "",
		},
	}

	suite.Run(t, s)
}

func (suite *AuthTestSuite) SetupSuite() {
	res := suite.db.Unscoped().Where("email = ?", suite.user.Email).Delete(&models.User{})
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	suite.SignUp(suite.user.Name, suite.user.Email, suite.user.Password)
	// Temporary Test User
	res = suite.db.Unscoped().Where("email = ?", "test12345@gmail.com").Delete(&models.User{})
	if res.Error != nil {
		fmt.Println(res.Error)
	}
}

func (suite *AuthTestSuite) TearDownSuite() {
	res := suite.db.Unscoped().Where("email = ?", suite.user.Email).Delete(&models.User{})
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	// Temporary Test User
	res = suite.db.Unscoped().Where("email = ?", "test12345@gmail.com").Delete(&models.User{})
	if res.Error != nil {
		fmt.Println(res.Error)
	}
}
