package adapters_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/AlecSmith96/fourth-floor-exercise/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

const (
	testAccessTokenString = "{\"access_token\":\"qq61yk4102dgxaz29feuuugrq44xqp\",\"expires_in\":4901701,\"token_type\":\"bearer\"}\n"
	videosResponseString = "{\"data\":[{\"id\":\"1622391775\",\"stream_id\":\"47293140269\",\"user_id\":\"563817305\",\"user_login\":\"playlostark\",\"user_name\":\"PlayLostArk\",\"title\":\"Lost Ark Midnight Circus TTRPG\",\"description\":\"\",\"created_at\":\"2022-10-12T18:45:18Z\",\"published_at\":\"2022-10-12T18:45:18Z\",\"url\":\"https://www.twitch.tv/videos/1622391775\",\"thumbnail_url\":\"https://thumbnail.jpg\",\"viewable\":\"public\",\"view_count\":53318,\"language\":\"en\",\"type\":\"archive\",\"duration\":\"2m30s\",\"muted_segments\":null}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6M319\"}}\n"
	token = "qq61yk4102dgxaz29feuuugrq44xqp"
	userID = "user-id"
	limit = "3"
)

// TwitchAdapterTestSuite contains additional fields that are used across all tests
type TwitchAdapterTestSuite struct {
	suite.Suite
}

func TestTwitchAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(TwitchAdapterTestSuite))
}

func (s *TwitchAdapterTestSuite) TestGetVideosForUser_HappyPath() {
	// set up adapter
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	accessTokenRequest, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	accessTokenRequest.Header.Set("Content-Type", "application/json")

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(adapters.VideosURL, userID, limit), nil /* body */)
	s.Assert().NoError(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Client-Id", config.ClientID)

	accessTokenResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(testAccessTokenString)),
	}

	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(videosResponseString)),
	}

	returnedVideoData := []entities.VideoData{
		{
			ID: "1622391775",
			StreamID: "47293140269",
			UserID: "563817305",
			UserLogin: "playlostark",
			UserName: "PlayLostArk",
			Title: "Lost Ark Midnight Circus TTRPG",
			Description: "",
			CreatedAt: "2022-10-12T18:45:18Z",
			PublishedAt: "2022-10-12T18:45:18Z",
			URL: "https://www.twitch.tv/videos/1622391775",
			ThumbnailURL: "https://thumbnail.jpg",
			Viewable: "public",
			ViewCount: 53318,
			Language: "en",
			Type: "archive",
			Duration: "2m30s",
			MutedSegments: nil,
		},
	}

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(accessTokenResponse, nil).Once()
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil).Once()

	// Act
	videos, err := adapter.GetVideosForUser(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().NoError(err)
	s.Assert().Equal(returnedVideoData, videos)
}

func (s *TwitchAdapterTestSuite) TestGetVideosForUser_AccessTokenError() {
	// set up adapter
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	accessTokenRequest, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	accessTokenRequest.Header.Set("Content-Type", "application/json")

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(adapters.VideosURL, userID, limit), nil /* body */)
	s.Assert().NoError(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Client-Id", config.ClientID)

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("access token error"))

	// Act
	videos, err := adapter.GetVideosForUser(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().EqualError(err, "access token error")
	s.Assert().Nil(videos)
	s.Assert().Equal("sending request", observedLogs.All()[0].Message)
}

func (s *TwitchAdapterTestSuite) TestGetVideosForUser_RequestReturnsError() {
	// set up adapter
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	accessTokenRequest, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	accessTokenRequest.Header.Set("Content-Type", "application/json")

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(adapters.VideosURL, userID, limit), nil /* body */)
	s.Assert().NoError(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Client-Id", config.ClientID)

	accessTokenResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(testAccessTokenString)),
	}

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(accessTokenResponse, nil).Once()
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("unexpected error")).Once()

	// Act
	videos, err := adapter.GetVideosForUser(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().EqualError(err, "unexpected error")
	s.Assert().Nil(videos)
	s.Assertions.Equal("sending request", observedLogs.All()[0].Message)
}

func (s *TwitchAdapterTestSuite) TestGetVideosForUser_UnauthorizedError() {
	// set up adapter
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	accessTokenRequest, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	accessTokenRequest.Header.Set("Content-Type", "application/json")

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(adapters.VideosURL, userID, limit), nil /* body */)
	s.Assert().NoError(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Client-Id", config.ClientID)

	accessTokenResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(testAccessTokenString)),
	}

	response := &http.Response{
		Status:     "401 Unauthorized",
		StatusCode: 401,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(videosResponseString)),
	}

	expectedErr := entities.NewUnauthorizedError()

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(accessTokenResponse, nil).Once()
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil).Once()

	// Act
	videos, err := adapter.GetVideosForUser(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().Equal(expectedErr, err)
	s.Assert().Nil(videos)
}

func (s *TwitchAdapterTestSuite) TestGetVideosForUser_UnmarshallingError() {
	// set up adapter
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	accessTokenRequest, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	accessTokenRequest.Header.Set("Content-Type", "application/json")

	// set up request
	request, err := http.NewRequest("GET", fmt.Sprintf(adapters.VideosURL, userID, limit), nil /* body */)
	s.Assert().NoError(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Client-Id", config.ClientID)

	accessTokenResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(testAccessTokenString)),
	}

	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString("1234")),
	}

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(accessTokenResponse, nil).Once()
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil).Once()

	// Act
	videos, err := adapter.GetVideosForUser(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().EqualError(err, "json: cannot unmarshal number into Go value of type entities.Video")
	s.Assert().Nil(videos)
	s.Assert().Equal("unmarshalling response body", observedLogs.All()[0].Message)
}

func (s *TwitchAdapterTestSuite) TestObtainAccessToken_HappyPath() {
	// set up adapter
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	request, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	request.Header.Set("Content-Type", "application/json")

	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(testAccessTokenString)),
	}

	expectedAccessToken := &entities.AccessToken{
		Token:     "qq61yk4102dgxaz29feuuugrq44xqp",
		ExpiresIn: 4901701,
		TokenType: "bearer",
	}

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil)

	// Act
	accessToken, err := adapter.ObtainAccessToken(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().NoError(err)
	s.Assert().Equal(expectedAccessToken, accessToken)
}

func (s *TwitchAdapterTestSuite) TestObtainAccessToken_RequestReturnsError() {
	// set up adapter
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	request, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	request.Header.Set("Content-Type", "application/json")

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("unexpected failure"))

	// Act
	accessToken, err := adapter.ObtainAccessToken(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().EqualError(err, "unexpected failure")
	s.Assert().Nil(accessToken)
	s.Assert().Equal("sending request", observedLogs.All()[0].Message)
	s.Assert().Equal("error", observedLogs.All()[0].Context[0].Key)
}

func (s *TwitchAdapterTestSuite) TestObtainAccessToken_RequestReturnsBadRequest() {
	// set up adapter
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	request, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	request.Header.Set("Content-Type", "application/json")

	response := &http.Response{
		Status:     "400 BadRequest",
		StatusCode: 400,
	}

	expectedErr := entities.NewBadRequestError()

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil)

	// Act
	accessToken, err := adapter.ObtainAccessToken(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().Equal(expectedErr, err)
	s.Assert().Nil(accessToken)
}

func (s *TwitchAdapterTestSuite) TestObtainAccessToken_FailToUnmarshalBody() {
	// set up adapter
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}
	client := &mocks.HTTPClient{}
	adapter := adapters.NewTwitchAdapter(observedLogger, config, client)

	// Arrange
	body := fmt.Sprintf(adapters.TokenString, config.ClientID, config.ClientSecret)
	requestBody := []byte(body)
	request, err := http.NewRequest("POST", adapters.TokenURL, bytes.NewBuffer(requestBody))
	s.Assert().NoError(err)
	request.Header.Set("Content-Type", "application/json")

	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString("1234")),
	}

	client.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil)

	// Act
	accessToken, err := adapter.ObtainAccessToken(config.ClientID, config.ClientSecret)

	// Assert
	s.Assert().EqualError(err, "json: cannot unmarshal number into Go value of type entities.AccessToken")
	s.Assert().Equal("unmarshalling response body", observedLogs.All()[0].Message)
	s.Assert().Nil(accessToken)
}
