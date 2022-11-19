package rest

import (
	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter defines a new gin router and routes
func NewRouter(logger *zap.Logger, twitchAdapter adapters.TwitchRequests, analyticsAdapter *adapters.AnalyticsAdapter) *gin.Engine {
	router := gin.Default()
	router.GET("/videos/:userID", func(c *gin.Context) {
		ViewHandler(c, logger, twitchAdapter, analyticsAdapter)
	})
	
	return router
}
