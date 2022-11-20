package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/rest"
	"github.com/AlecSmith96/fourth-floor-exercise/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestViewHandler_HappyPath(t *testing.T) {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=3", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}

	returnedVideos := []entities.VideoData{
		{
			Title: "Video #1",
			Duration: "1m20s",
			ViewCount: 12340,
		},
		{
			Title: "Video #2",
			Duration: "40m15s",
			ViewCount: 150067,
		},
		{
			Title: "Video #3",
			Duration: "3h12m56s",
			ViewCount: 130589,
		},
	}

	returnedAnalytics := &entities.VideoAnalytics{
		SumOfVideoViews: (12340 + 150067 + 130589),
		AverageViewsPerVideo: float64(13665.33),
		SumOfVideoLengths: "3h54m31s",
		AverageViewsPerMinute: float64(4563.45),
		MostViewedVideo: entities.MostViewedVideo{
			Title: "Video #2",
			Views: 150067,
		},
	}

	twitchAdapter.On("GetVideosForUser", "1234", "3").Return(returnedVideos, nil)
	analyticsAdapter.On("GetVideoAnalytics", returnedVideos). Return(returnedAnalytics, nil)

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"SumOfVideoViews\":292996,\"AverageViewsPerVideo\":13665.33,\"SumOfVideoLengths\":\"3h54m31s\",\"AverageViewsPerMinute\":4563.45,\"MostViewedVideo\":{\"Title\":\"Video #2\",\"Views\":150067}}", w.Body.String())
}

func TestViewHandler_InputValidationFailed(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=3", nil)

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"code\":400,\"message\":\"a bad request error occurred\"}", w.Body.String())
	assert.Equal(t, "inputs missing in request", observedLogs.All()[0].Message)
}

func TestViewHandler_InputValidation_NonIntegerLimitValue(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=d", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"code\":400,\"message\":\"a bad request error occurred\"}", w.Body.String())
	assert.Equal(t, "limit input invalid", observedLogs.All()[0].Message)
}

func TestViewHandler_InputValidation_MissingLimitValue(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"code\":400,\"message\":\"expected limit query parameter, but received none\"}", w.Body.String())
	assert.Equal(t, "missing limit query parameter", observedLogs.All()[0].Message)
}

func TestViewHandler_TwitchAdapterReturnsError(t *testing.T) {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=3", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}
	twitchAdapter.On("GetVideosForUser", "1234", "3").Return(nil, errors.New("test error"))

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"code\":500,\"message\":\"an internal server error occurred\"}", w.Body.String())
}

func TestViewHandler_NoVideosReturnedForUser(t *testing.T) {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=3", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}
	twitchAdapter.On("GetVideosForUser", "1234", "3").Return([]entities.VideoData{}, nil)

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"user has no videos to analyse\"}", w.Body.String())
}

func TestViewHandler_AnalyticsAdapterReturnsError(t *testing.T) {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	twitchAdapter := &mocks.TwitchRequests{}
	analyticsAdapter := &mocks.AnalyticsCalls{}
	gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Request = httptest.NewRequest("GET", "locslhost:8080/videos/1234?limit=3", nil)
	ctx.Params = []gin.Param{
		{
			Key: "userID",
			Value: "1234",
		},
	}

	returnedVideos := []entities.VideoData{
		{
			Title: "Video #1",
			Duration: "1m20s",
			ViewCount: 12340,
		},
		{
			Title: "Video #2",
			Duration: "40m15s",
			ViewCount: 150067,
		},
		{
			Title: "Video #3",
			Duration: "3h12m56s",
			ViewCount: 130589,
		},
	}

	twitchAdapter.On("GetVideosForUser", "1234", "3").Return(returnedVideos, nil)
	analyticsAdapter.On("GetVideoAnalytics", returnedVideos). Return(nil, entities.NewNotFoundError())

	rest.ViewHandler(ctx, observedLogger, twitchAdapter, analyticsAdapter)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"code\":404,\"message\":\"a not found errror occured\"}", w.Body.String())
}
