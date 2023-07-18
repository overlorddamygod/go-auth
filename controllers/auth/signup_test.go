package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (ts *AuthTestSuite) TestRequiresName() {
	var buffer bytes.Buffer
	assert.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    ts.user.Email,
		"password": ts.user.Password,
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response Response

	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))
	assert.Equal(ts.T(), http.StatusBadRequest, w.Code)
	assert.Equal(ts.T(), "invalid params", response.Message)
}

func (ts *AuthTestSuite) TestNewUser() {
	w, _, _ := ts.SignUp("Test12345", "test12345@gmail.com", "test123456")

	assert.Equal(ts.T(), http.StatusCreated, w.Code)
}

func (ts *AuthTestSuite) TestUserAlreadyRegistered() {
	w, _, _ := ts.SignUp(ts.user.Name, ts.user.Email, ts.user.Password)

	assert.NotEqual(ts.T(), http.StatusCreated, w.Code)
}

func (ts *AuthTestSuite) SignUp(name string, email string, password string) (*httptest.ResponseRecorder, *http.Request, Response) {
	var buffer bytes.Buffer
	assert.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": password,
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.router.ServeHTTP(w, req)

	var response Response
	assert.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&response))
	return w, req, response
}
