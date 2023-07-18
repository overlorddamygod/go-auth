package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/overlorddamygod/go-auth/utils"
	"github.com/stretchr/testify/assert"
)

type RefreshResponse struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

func (ts *AuthTestSuite) TestInvalidRefreshToken() {
	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Refresh-Token", "sadasdwqeqwewq")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response SigninResponse

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))

	assert.Equal(ts.T(), http.StatusUnauthorized, w.Code)
}

func (ts *AuthTestSuite) TestRefreshToken() {
	_, _, res := ts.Signin(ts.user.Email, ts.user.Password, "email")

	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Refresh-Token", res.RefreshToken)

	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response RefreshResponse

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))
	ts.T().Log(response, res)

	assert.Equal(ts.T(), http.StatusOK, w.Code)

	assert.NotEmpty(ts.T(), response.AccessToken)

	_, err := utils.JwtAccessTokenVerify(response.AccessToken)
	assert.NoError(ts.T(), err)
}
