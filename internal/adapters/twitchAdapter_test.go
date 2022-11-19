package adapters_test

import (
	"testing"

	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type TwitchAdapterTestSuite struct {
	suite.Suite
}

func TestTwitchAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(TwitchAdapterTestSuite))
}

func (s *TwitchAdapterTestSuite) TestObtainAccessToken() {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	config := entities.ConfigAuth{ClientID: "client-id", ClientSecret: "client-secret"}

	adapter := adapters.NewTwitchAdapter(observedLogger, config)

	accessToken, err := adapter.ObtainAccessToken(config.ClientID, config.ClientSecret)
	s.Assert().NoError(err)
	s.Assert().Equal(&entities.AccessToken{}, accessToken)
}