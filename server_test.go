package vanilla

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	Start(":8282")
	time.Sleep(1 * time.Second)

	req := httptest.NewRequest("GET", "localhost:8282/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome Gin Server", w.Body.String())

	Stop()
}
