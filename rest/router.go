package rest

import (
	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/gin-gonic/gin"
)

// NewRouter defines a new gin router and routes
func NewRouter(twitchAdapter *adapters.TwitchAdapter) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/videos/:channelID", func(c *gin.Context) {
		ViewHandler(c, twitchAdapter)
	})
	
	return router
}
