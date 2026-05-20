package handler

import (
	"CricTail_Backend/database/dbHelper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayerProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	playerStats, err := dbHelper.GetPlayerStatsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get player Stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player-stats": playerStats,
	})
}

func GetAllPlayers(c *gin.Context) {
	search := c.Query("search")

	players, err := dbHelper.GetAllPlayers(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get players"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player": players,
		"total":  len(players),
	})
}

func UpdatePlayerProfile(c *gin.Context) {

}
