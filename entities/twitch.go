package entities

// AccessToken stores a clientID with its corresponding token, used for requests to twitch API
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int
	TokenType string
}
