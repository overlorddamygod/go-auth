package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (ts *AuthTestSuite) TestInvalidRefreshTokenForSignout() {
	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signout", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Refresh-Token", "sadasdwqeqwewq")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response Response

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))

	assert.Equal(ts.T(), http.StatusBadRequest, w.Code)
}

func (ts *AuthTestSuite) TestSignoutSuccessfull() {
	_, _, res := ts.Signin(ts.user.Email, ts.user.Password, "email")

	var buffer bytes.Buffer

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signout", &buffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Refresh-Token", res.RefreshToken)

	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response RefreshResponse

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))

	assert.Equal(ts.T(), http.StatusOK, w.Code)
}
