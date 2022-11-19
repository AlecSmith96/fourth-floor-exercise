package adapters

import (
	"math"
	"time"

	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"go.uber.org/zap"
)

type AnalyticsAdapter struct {
	Logger     *zap.Logger
}

type AnalyticsCalls interface {
	GetVideoAnalytics(videos []entities.VideoData) (*entities.VideoAnalytics, error)
}

// GetVideoAnalytics calculates the necessary information from the passed slice of videos and returns it
func (adapter *AnalyticsAdapter) GetVideoAnalytics(videos []entities.VideoData) (*entities.VideoAnalytics, error) {
	sumOfVideoViews := 0
	var totalDuration time.Duration
	var mostViewedVideo entities.VideoData

	for _, video := range videos {
		sumOfVideoViews += video.ViewCount
		videoDuration, err := adapter.parseVideoDuration(video.Duration)
		if err != nil {
			adapter.Logger.Error("parsing video duration", zap.Error(err))
		}

		totalDuration += videoDuration
		if video.ViewCount > mostViewedVideo.ViewCount {
			mostViewedVideo = video
		}
	}
	
	numMinutes := totalDuration / (1*time.Minute)
	averageViewsPerMinute := float64(sumOfVideoViews) / float64(numMinutes)
	averageViewsPerVideo := float64(sumOfVideoViews)/float64(len(videos))
	
	return &entities.VideoAnalytics{
		SumOfVideoViews: sumOfVideoViews,
		AverageViewsPerVideo: math.Round(averageViewsPerVideo*100)/100,
		SumOfVideoLengths: totalDuration.String(),
		AverageViewsPerMinute: math.Round(averageViewsPerMinute*100)/100,
		MostViewedVideo: entities.MostViewedVideo{
			Title: mostViewedVideo.Title,
			Views: mostViewedVideo.ViewCount,
		},
	}, nil
}

func (adapter *AnalyticsAdapter) parseVideoDuration(duration string) (time.Duration, error) {
	videoDuration, err := time.ParseDuration(duration)
		if err != nil {
			adapter.Logger.Error("parsing video duration", zap.Error(err))
			return time.Duration(0), err
		}
	return videoDuration, nil
}

func NewAnalyticsAdapter(logger *zap.Logger) AnalyticsCalls {
	return &AnalyticsAdapter{
		Logger: logger,
	}
}
