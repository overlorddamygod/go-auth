package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/overlorddamygod/go-auth/utils"
	"github.com/stretchr/testify/assert"
)

type SigninResponse struct {
	Error        bool   `json:"error"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (ts *AuthTestSuite) TestUserNotExists() {
	var buffer bytes.Buffer
	assert.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    "test@gmail.com",
		"password": ts.user.Password,
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin?type=email", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response SigninResponse

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))
	fmt.Println(response)
	ts.T().Log(response)
	assert.Equal(ts.T(), http.StatusNotFound, w.Code)
}

func (ts *AuthTestSuite) TestUserSignin() {
	w, _, response := ts.Signin(ts.user.Email, ts.user.Password, "email")

	assert.Equal(ts.T(), http.StatusOK, w.Code)
	assert.NotEmpty(ts.T(), response.AccessToken)
	assert.NotEmpty(ts.T(), response.RefreshToken)

	_, err := utils.JwtAccessTokenVerify(response.AccessToken)
	assert.NoError(ts.T(), err)

	_, err = utils.JwtRefreshTokenVerify(response.RefreshToken)
	assert.NoError(ts.T(), err)
}

func (ts *AuthTestSuite) Signin(email string, password string, _type string) (*httptest.ResponseRecorder, *http.Request, SigninResponse) {
	var buffer bytes.Buffer
	assert.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    email,
		"password": password,
	}))

	if _type == "" {
		_type = "email"
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin?type="+_type, &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response SigninResponse

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))
	return w, req, response
}
