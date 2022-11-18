package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer creates the server instance from the router
func NewHTTPServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}
