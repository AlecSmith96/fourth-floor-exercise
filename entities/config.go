package entities

type Config struct {
	Rest    ConfigRest
	Logging ConfigLogging
	Auth    ConfigAuth
}

type ConfigRest struct {
	Port string
}

type ConfigLogging struct {
	LogLevel string
	Encoding string
}

type ConfigAuth struct {
	ClientID     string
	ClientSecret string
}
