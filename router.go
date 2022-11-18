package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NewRouter defines a new gin router and routes
func NewRouter(twitchAdapter *TwitchAdapter) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/videos/:channelID", func(c *gin.Context) {
		ViewHandler(c, twitchAdapter)
	})
	
	return router
}

// NewHTTPServer creates the server instance from the router
func NewHTTPServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}
