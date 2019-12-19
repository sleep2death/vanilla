package vanilla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
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

func TestSocketIO(t *testing.T) {
	router, err := setupRouter()
	if err != nil {
		t.Fatal(err)
	}

	io, err := setupIO(router)
	if err != nil {
		t.Fatal(err)
	}

	defer io.Close()
	go io.Serve()

	// login and get the token
	token, err := getToken(router)
	if err != nil {
		t.Error(err)
	}

	go router.Run(":8082")
	time.Sleep(time.Second)

	header := make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	dialer := engineio.Dialer{
		Transports: []transport.Transport{polling.Default, websocket.Default},
	}
	conn, err := dialer.Dial("http://localhost:8082/io/", header)
	if err != nil {
		log.Fatalln("dial error:", err)
	}
	defer conn.Close()
	// log.Println(conn.ID(), conn.LocalAddr(), "->", conn.RemoteAddr(), "with", conn.RemoteHeader())
}
