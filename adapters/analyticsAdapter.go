package adapters

import (
	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"go.uber.org/zap"
)

type AnalyticsAdapter struct {
	Logger     *zap.Logger
}

func (adapter *AnalyticsAdapter) GetVideoAnalytics(videos []*entities.VideoData) (*entities.VideoAnalytics, error) {
	return &entities.VideoAnalytics{
		SumOfVideoViews: getSumOfVideoViews(videos),
		AverageViewsPerVideo: getAvgViewsPerVideo(videos),
	}, nil
}

func getSumOfVideoViews(videos []*entities.VideoData) int {
	views := 0

	for _, video := range videos {
		views += video.ViewCount
	}
	return views
}

func getAvgViewsPerVideo(videos []*entities.VideoData) int {
	return 0
}

func NewAnalyticsAdapter(logger *zap.Logger) *AnalyticsAdapter {
	return &AnalyticsAdapter{
		Logger: logger,
	}
}
