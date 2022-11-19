package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/gin-gonic/gin"
)

// NewHTTPServer creates the server instance from the router
func NewHTTPServer(config entities.ConfigRest, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
		Handler: router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}
