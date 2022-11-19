package entities

// AccessToken stores a clientID with its corresponding token, used for requests to twitch API
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

// Video contains all data returned from GetVideos requests
type Video struct {
	Data       []VideoData `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// VideoData contains all data returned on a single video
type VideoData struct {
	ID            string         `json:"id"`
	StreamID      string         `json:"stream_id"`
	UserID        string         `json:"user_id"`
	UserLogin     string         `json:"user_login"`
	UserName      string         `json:"user_name"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	CreatedAt     string         `json:"created_at"`
	PublishedAt   string         `json:"published_at"`
	URL           string         `json:"url"`
	ThumbnailURL  string         `json:"thumbnail_url"`
	Viewable      string         `json:"viewable"`
	ViewCount     int            `json:"view_count"`
	Language      string         `json:"language"`
	Type          string         `json:"type"`
	Duration      string         `json:"duration"`
	MutedSegments []MutedSegment `json:"muted_segments"`
}

// MutedSegment a segment of a video that Twitch Audio Recognition muted
type MutedSegment struct {
	Duration int `json:"duration"`
	Offset   int `json:"offset"`
}

// Pagination contains information used to page through results
type Pagination struct {
	Cursor string `json:"cursor"`
}

// VideoAnalytics response struct for '/videos' endpoint
type VideoAnalytics struct {
	SumOfVideoViews       int
	AverageViewsPerVideo  float64
	SumOfVideoLengths     string
	AverageViewsPerMinute float64
	MostViewedVideo       MostViewedVideo
}

// MostViewedVideo struct representing the most viewed video returned by the Twitch API
type MostViewedVideo struct {
	Title string
	Views int
}
