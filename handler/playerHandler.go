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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player-stats": playerStats,
	})
}
