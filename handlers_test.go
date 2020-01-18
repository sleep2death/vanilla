package vanilla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sleep2death/vanilla/core"
	"github.com/stretchr/testify/assert"
)

func getToken(r *gin.Engine) (string, error) {

	// login and get the token
	rb, _ := json.Marshal(map[string]string{
		"username": "aspirin2d",
		"password": "Passw0rd!",
	})

	req, _ := http.NewRequest("POST", "login", bytes.NewBuffer(rb))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusOK, w.Code)

	if w.Code != http.StatusOK {
		return "", errors.New("failed to login")
	}

	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		return "", err
	}
	return resp["token"], nil
}

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
	token, err := getToken(r)
	if err != nil {
		t.Error(err)
	}
	// t.Logf("response token is %s", resp["token"])

	// get api/ping with the token
	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// get api/ping with the wrong token
	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "abc"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]string

	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "token is invalid", resp["reason"])

	// get api/ping with the expired token
	time.Sleep(expire + time.Second)

	req, _ = http.NewRequest("GET", "api/ping", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "token is expired", resp["reason"])
}

func TestPlayerInfoHandler(t *testing.T) {
	r, err := setupRouter()
	if err != nil {
		t.Error(err)
	}

	// login and get the token
	token, err := getToken(r)
	if err != nil {
		t.Error(err)
	}
	// t.Logf("response token is %s", resp["token"])

	// get api/ping with the token
	req, _ := http.NewRequest("GET", "api/playerinfo?username=aspirin2d", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp core.Player
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWebsocket(t *testing.T) {
	router, err := setupRouter()
	if err != nil {
		t.Fatal(err)
	}

	go router.Run(":8082")
	time.Sleep(time.Second)

	// login and get the token
	token, err := getToken(router)
	if err != nil {
		t.Error(err)
	}

	conn, resp, err := websocket.DefaultDialer.Dial("ws://localhost:8082/ws?token="+token, nil)
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		t.Error(string(body))
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		t.Error(err)
	}

}
