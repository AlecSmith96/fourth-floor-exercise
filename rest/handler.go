package rest

import (
	"net/http"
	"strconv"

	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ViewHandler obtains the data parsed in the request and calls the necessary adapter logic to
// return the needed data in the response
func ViewHandler(c *gin.Context, logger *zap.Logger, twitchAdapter adapters.TwitchRequests, analyticsAdapter *adapters.AnalyticsAdapter) {
	channelID := c.Param("userID")
	limit := c.Query("limit")

	numberOfVideos, err := strconv.Atoi(limit)
	if err != nil {
		logger.Warn("invalid query param: %v", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"message": "received request with invalid parameters",
		})
		return
	}

	videos, err := twitchAdapter.GetVideosForUser(channelID, numberOfVideos)
	if err != nil {
		handlerError(c, err)
	}

	analytics, err := analyticsAdapter.GetVideoAnalytics(videos)
	if err != nil {
		handlerError(c, err)
	}

	c.JSON(http.StatusOK, analytics)
}

// handleError takes the error returned by adapters and returns correct http response
func handlerError(c *gin.Context, err error) {
	switch e := err.(type) {
	case entities.ResponseError:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Code,
			"message": e.PresentableError,
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"message": "an internal server error occurred",
		})
	}
}
