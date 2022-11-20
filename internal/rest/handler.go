package rest

import (
	"net/http"
	"strconv"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ViewHandler acts as the controller for the videos use case, it obtains the data parsed in the request and calls the necessary adapter logic to
// return the needed data in the response
func ViewHandler(c *gin.Context, logger *zap.Logger, twitchAdapter adapters.TwitchRequests, analyticsAdapter adapters.AnalyticsCalls) {
	channelID := c.Param("userID")
	limit := c.Query("limit")

	if err := validateInputs(channelID, limit, logger); err != nil {
		handlerError(c, err)
		return
	}

	videos, err := twitchAdapter.GetVideosForUser(channelID, limit)
	if err != nil {
		handlerError(c, err)
		return
	}

	if len(videos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "user has no videos to analyse",
		})
		return
	}

	analytics, err := analyticsAdapter.GetVideoAnalytics(videos)
	if err != nil {
		handlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// validateInputs ensures the parameters aren't missing and that limit is an integer, otherwise a 400 response is returned
func validateInputs(channelID, limit string, logger *zap.Logger) error {
	// ideally would add a string length constraint here for channel id
	if channelID == "" {
		logger.Error("inputs missing in request")
		return entities.NewBadRequestError()
	}

	if limit == "" {
		logger.Error("missing limit query parameter")
		return &entities.ResponseError{
			Code: http.StatusBadRequest,
			PresentableError: "expected limit query parameter, but received none",
		}
	}

	_, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error("limit input invalid", zap.Error(err))
		return entities.NewBadRequestError()
	}

	return nil
}

// handleError takes the error returned by adapters and returns correct http response
func handlerError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *entities.ResponseError:
		c.JSON(e.Code, gin.H{
			"code":    e.Code,
			"message": e.PresentableError,
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "an internal server error occurred",
		})
	}
}
