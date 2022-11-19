package adapters_test

import (
	"testing"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type AnalyticsAdapterTestSuite struct {
	suite.Suite
}

func TestAnalyticsAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsAdapterTestSuite))
}

func (s *AnalyticsAdapterTestSuite) TestGetVideoAnalytics_HappyPath() {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	adapter := adapters.NewAnalyticsAdapter(observedLogger)

	videos := &entities.Video{
		Data: []entities.VideoData{
			{
				Title: "Video #1",
				Duration: "1m20s",
				ViewCount: 156,
			},
			{
				Title: "Video #2",
				Duration: "1h12m20s",
				ViewCount: 15432,
			},
			{
				Title: "Video #3",
				Duration: "42m30s",
				ViewCount: 15432,
			},
		},
	}

	analytics, err := adapter.GetVideoAnalytics(videos.Data)
	s.Assert().NoError(err)
	s.Assert().Equal(31020, analytics.SumOfVideoViews)
	s.Assert().Equal(float64(10340), analytics.AverageViewsPerVideo)
	s.Assert().Equal("1h56m10s", analytics.SumOfVideoLengths)
	s.Assert().Equal(float64(267.41), analytics.AverageViewsPerMinute)
	s.Assert().Equal("Video #2", analytics.MostViewedVideo.Title)
	s.Assert().Equal(15432, analytics.MostViewedVideo.Views)
}

func (s *AnalyticsAdapterTestSuite) TestParseVideoDuration_ValidDuration() {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	adapter := adapters.NewAnalyticsAdapter(observedLogger)

	videos := &entities.Video{
		Data: []entities.VideoData{
			{
				Duration: "1m20s",
			},
		},
	}

	// as parseVideoDuration is private, it can't be called directly here
	analytics, err := adapter.GetVideoAnalytics(videos.Data)
	s.Assert().NoError(err)
	s.Assert().Equal("1m20s", analytics.SumOfVideoLengths)
}

func (s *AnalyticsAdapterTestSuite) TestParseVideoDuration_InvalidDuration() {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	adapter := adapters.NewAnalyticsAdapter(observedLogger)

	videos := &entities.Video{
		Data: []entities.VideoData{
			{
				Duration: "invalid",
			},
		},
	}

	// as parseVideoDuration is private, it can't be called directly here
	analytics, err := adapter.GetVideoAnalytics(videos.Data)
	s.Assert().EqualError(err, "time: invalid duration \"invalid\"")
	s.Assert().Equal("error", observedLogs.All()[0].Context[0].Key)
	s.Assert().Equal("parsing video duration", observedLogs.All()[0].Message)
	s.Assert().Nil(analytics)
}

func (s *AnalyticsAdapterTestSuite) TestNewAnalyticsAdapter() {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	adapter := adapters.NewAnalyticsAdapter(observedLogger)

	s.Assert().IsType(&adapters.AnalyticsAdapter{}, adapter)
}