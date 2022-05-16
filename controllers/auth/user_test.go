package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (ts *AuthTestSuite) TestInvalidToken() {
	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/me", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", "sadasdwqeqwewq")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	assert.NotEqual(ts.T(), http.StatusOK, w.Code)
}

func (ts *AuthTestSuite) TestMe() {
	_, _, res := ts.Signin(ts.user.Email, ts.user.Password, "email")

	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", res.AccessToken)

	w := httptest.NewRecorder()
	ts.T().Log(res, res.AccessToken, w.Body.String())

	ts.router.ServeHTTP(w, req)

	assert.Equal(ts.T(), http.StatusOK, w.Code)
}
