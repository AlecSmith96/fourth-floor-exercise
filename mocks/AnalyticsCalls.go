// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entities "github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
import mock "github.com/stretchr/testify/mock"

// AnalyticsCalls is an autogenerated mock type for the AnalyticsCalls type
type AnalyticsCalls struct {
	mock.Mock
}

// GetVideoAnalytics provides a mock function with given fields: videos
func (_m *AnalyticsCalls) GetVideoAnalytics(videos []entities.VideoData) (*entities.VideoAnalytics, error) {
	ret := _m.Called(videos)

	var r0 *entities.VideoAnalytics
	if rf, ok := ret.Get(0).(func([]entities.VideoData) *entities.VideoAnalytics); ok {
		r0 = rf(videos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.VideoAnalytics)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]entities.VideoData) error); ok {
		r1 = rf(videos)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}