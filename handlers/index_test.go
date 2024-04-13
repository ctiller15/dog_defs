package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "favicon.ico")
	r.LoadHTMLGlob("../templates/*")
	router := SetupIndex(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Welcome to DogDefs")
}

func TestNotFoundRoute(t *testing.T) {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "favicon.ico")
	r.LoadHTMLGlob("../templates/*")
	router := SetupIndex(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/gewgw", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "404 NOT FOUND")

}
