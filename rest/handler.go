package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/gin-gonic/gin"
)

func ViewHandler(c *gin.Context, twitchAdapter *adapters.TwitchAdapter) {
	channelID := c.Param("channelID")
	numberOfVideosQueryParam := c.Query("limit")

	numberOfVideos, err := strconv.Atoi(numberOfVideosQueryParam)
	if err != nil {
		log.Printf("invalid query param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"message": "received request with invalid parameters",
		})
		return
	}

	twitchAdapter.GetVideosForUser(channelID, numberOfVideos)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	  })
}