package adapters

import (
	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"go.uber.org/zap"
)

type AnalyticsAdapter struct {
	Logger     *zap.Logger
}

func (adapter *AnalyticsAdapter) GetVideoAnalytics(videos []*entities.VideoData) (*entities.VideoAnalytics, error) {
	return &entities.VideoAnalytics{}, nil
}

func NewAnalyticsAdapter(logger *zap.Logger) *AnalyticsAdapter {
	return &AnalyticsAdapter{
		Logger: logger,
	}
}
