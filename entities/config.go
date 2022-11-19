package entities

// Config top level config struct that gets populated by config adapter
type Config struct {
	Rest    ConfigRest
	Logging ConfigLogging
	Auth    ConfigAuth
}

// ConfigRest contains config points for http server and router
type ConfigRest struct {
	Port    string
	GinMode string
}

// ConfigLogging contains config points for the zap logger initialisation
type ConfigLogging struct {
	LogLevel string
	Encoding string
}

// ConfigAuth contains authentication information for the Twitch API
type ConfigAuth struct {
	ClientID     string
	ClientSecret string
}
