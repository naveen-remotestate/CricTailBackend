package middleware

import (
	"CricTail_Backend/database/dbHelper"
	"CricTail_Backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token in header"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			c.Abort()
			return
		}
		sessionID, _ := claims["session_id"].(string)
		//fmt.Println("---------------->", sessionID)

		userID, err := dbHelper.GetUserIDBySessionID(sessionID)
		if err != nil || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("session_id", sessionID)

		c.Next()
	}
}
