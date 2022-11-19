package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"go.uber.org/zap"
)

const (
	tokenString = `{"client_id": "%s", "client_secret": "%s", "grant_type": "client_credentials"}`
	tokenURL    = "https://id.twitch.tv/oauth2/token"
	videosURL   = "https://api.twitch.tv/helix/videos?user_id=%s&first=%s"
)

type TwitchAdapter struct {
	Auth       *entities.ConfigAuth
	HTTPClient *http.Client // ideally this would be swapped out for a separate client that would make the http calls
	Logger     *zap.Logger
}

// TwitchRequests interface used for easily mocking functionality
type TwitchRequests interface {
	ObtainAccessToken(clientID, clientSecret string) (*entities.AccessToken, error)
	GetVideosForUser(userID string, limit string) ([]entities.VideoData, error)
}

// ObtainAccessToken sends oauth request to twitch API to obtain the access token for further requests
func (adapter *TwitchAdapter) ObtainAccessToken(clientID, clientSecret string) (*entities.AccessToken, error) {
	// set up request
	body := fmt.Sprintf(tokenString, clientID, clientSecret)
	requestBody := []byte(body)
	request, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(requestBody))
	if err != nil {
		adapter.Logger.Error("obtaining access token", zap.Error(err))
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := adapter.HTTPClient.Do(request)
	if err != nil {
		adapter.Logger.Error("sending request", zap.Error(err))
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, handleUnsuccessfulStatus(response)
	}

	defer response.Body.Close()

	// unmarshal response body into struct and return it
	accessToken := &entities.AccessToken{}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		adapter.Logger.Error("response body", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(responseBody, accessToken)
	if err != nil {
		adapter.Logger.Error("unmarshalling response body", zap.Error(err))
		return nil, err
	}

	return accessToken, nil
}

// GetVideosForUser queries the Twitch API for the last number of videos specified by the limit parameter
func (adapter *TwitchAdapter) GetVideosForUser(userID string, limit string) ([]entities.VideoData, error) {
	accessToken, err := adapter.ObtainAccessToken(adapter.Auth.ClientID, adapter.Auth.ClientSecret)
	if err != nil {
		adapter.Logger.Error("obtaining access token for request", zap.Error(err))
		return nil, err
	}

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(videosURL, userID, limit), nil /* body */)
	if err != nil {
		adapter.Logger.Error("creating request", zap.Error(err))
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Token))
	request.Header.Set("Client-Id", adapter.Auth.ClientID)

	response, err := adapter.HTTPClient.Do(request)
	if err != nil {
		adapter.Logger.Error("sending request", zap.Error(err))
		return nil, err
	}
	defer response.Body.Close()

	// if unsuccessful call, return error back to handler
	if response.Status != "200 OK" {
		return nil, handleUnsuccessfulStatus(response)
	}

	// marshal response into struct
	videosInResponse := &entities.Video{}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		adapter.Logger.Error("response body", zap.Error(err))
		return nil, nil
	}
	err = json.Unmarshal(responseBody, videosInResponse)
	if err != nil {
		adapter.Logger.Error("unmarshalling response body", zap.Error(err))
		return nil, nil
	}

	return videosInResponse.Data, nil
}

// handleUnsuccessfulStatus returns correct ResponseError based on status code
func handleUnsuccessfulStatus(response *http.Response) error {
	switch response.StatusCode {
	case http.StatusBadRequest:
		return entities.NewBadRequestError()
	case http.StatusUnauthorized:
		return entities.NewUnauthorizedError()
	case http.StatusNotFound:
		return entities.NewNotFoundError()
	default:
		return errors.New("Unexpected response status occurred ")
	}
}

func NewTwitchAdapter(logger *zap.Logger, config entities.ConfigAuth) TwitchRequests {
	return &TwitchAdapter{
		HTTPClient: &http.Client{},
		Auth:       &config,
		Logger:     logger.Named("twitchAdapter"),
	}
}
