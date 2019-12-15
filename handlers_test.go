package vanilla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	r, err := setupRouter()
	if err != nil {
		t.Error(err)
	}

	// set expire time to 5s
	expire = time.Second * 3

	// get api/ping with fail if not login
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "api/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// login and get the token
	rb, _ := json.Marshal(map[string]string{
		"username": "aspirin2d",
		"password": "Passw0rd!",
	})

	req, _ = http.NewRequest("POST", "login", bytes.NewBuffer(rb))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "ok", resp["status"])
	// t.Logf("response token is %s", resp["token"])

	// get api/ping with the token
	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp["token"]))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// get api/ping with the wrong token
	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "abc"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "token is invalid", resp["reason"])

	// get api/ping with the expired token
	time.Sleep(expire + time.Second)

	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp["token"]))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "token is expired", resp["reason"])
}
