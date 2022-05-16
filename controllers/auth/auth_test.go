package auth_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/server"
	"github.com/stretchr/testify/suite"
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
	db := db.GetDB()
	router := server.InitForTest()

	s := &AuthTestSuite{
		db:     db,
		router: router,
	}

	suite.Run(t, s)
}

func (suite *AuthTestSuite) SetupSuite() {
	suite.db = db.GetDB()
	suite.user = User{
		Name:         "test123",
		Email:        "test123@gmail.com",
		Password:     "test123",
		AccessToken:  "",
		RefreshToken: "",
	}
	// User for whole test suite
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
