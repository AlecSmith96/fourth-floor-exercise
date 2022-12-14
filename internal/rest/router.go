package rest

import (
	"net/http"
	"time"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter defines a new gin router and routes
func NewRouter(config *entities.ConfigRest, logger *zap.Logger, twitchAdapter adapters.TwitchRequests, analyticsAdapter adapters.AnalyticsCalls) *gin.Engine {
	gin.SetMode(config.GinMode)
	router := gin.New()

	// use injected logger for logging and recovery in gin router
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	
	// default routes for NotFound and MethodNotAllowed
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"message": "method not allowed",
		})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "page not found",
		})
	})

	router.GET("/videos/:userID", func(c *gin.Context) {
		ViewHandler(c, logger, twitchAdapter, analyticsAdapter)
	})

	return router
}
