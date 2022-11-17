package main

import (
	"fmt"
	"log"
	"net/http"
)

type TwitchAdapter struct {
	Logger *log.Logger
}

// https://dev.twitch.tv/docs/api/reference#get-videos
func (adapter *TwitchAdapter) GetVideosForUser(channelID string, numberOfVideos int) (interface{}, error) {
	response, err := http.Get(fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s", channelID))
	if err != nil {
		fmt.Printf("error in GetVideos request: %v", err)
	}

	log.Panicln(response)

	return nil, nil
}

func NewTwitchAdapter(logger *log.Logger) *TwitchAdapter {
	return &TwitchAdapter{
		Logger: logger,
	}
}
