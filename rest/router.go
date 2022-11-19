package rest

import (
	"time"

	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter defines a new gin router and routes
func NewRouter(config *entities.ConfigRest, logger *zap.Logger, twitchAdapter adapters.TwitchRequests, analyticsAdapter *adapters.AnalyticsAdapter) *gin.Engine {
	gin.SetMode(config.GinMode)
	router := gin.New()

	// use injected logger for logging and recovery in gin router
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.GET("/videos/:userID", func(c *gin.Context) {
		ViewHandler(c, logger, twitchAdapter, analyticsAdapter)
	})

	return router
}
