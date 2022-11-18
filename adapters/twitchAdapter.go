package adapters

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlecSmith96/fourth-floor-exercise/entities"
)

type TwitchAdapter struct {
	HTTPClient *http.Client
	Logger     *log.Logger
}

func (adapter *TwitchAdapter) ObtainAccessToken() *entities.AccessToken {

	// curl -X POST 'https://id.twitch.tv/oauth2/token' \
	// -H 'Content-Type: application/x-www-form-urlencoded' \
	// -d 'client_id=<your client id goes here>&client_secret=<your client secret goes here>&grant_type=client_credentials'

	return &entities.AccessToken{
		Token:    "",
		ExpiresIn: 0,
		TokenType: "",
	}
}

// https://dev.twitch.tv/docs/api/reference#get-videos
func (adapter *TwitchAdapter) GetVideosForUser(channelID string, limit int) (interface{}, error) {
	accessToken := adapter.ObtainAccessToken()
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s", channelID), nil /* body */)
	if err != nil {
		fmt.Printf("error creating /videos request: %v", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", accessToken.TokenType, accessToken.Token))
	request.Header.Set("Client-Id", ""/**/)

	response, err := adapter.HTTPClient.Do(request)
	if err != nil {
		fmt.Printf("error sending /videos request: %v", err)
	}

	log.Println(response)

	return nil, nil
}

func NewTwitchAdapter(client *http.Client, logger *log.Logger) *TwitchAdapter {
	return &TwitchAdapter{
		Logger: logger,
	}
}
