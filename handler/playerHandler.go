package handler

import (
	"CricTail_Backend/database"
	"CricTail_Backend/database/dbHelper"
	"CricTail_Backend/models"
	"CricTail_Backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterUser(c *gin.Context) {
	//req will hold data sent from user in form of body
	var req models.RegisterUser
	//maps user entered data with req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request- must send these fields in body -full_name, mobile_number, password"})
		return
	}

	req.FullName = strings.TrimSpace(req.FullName)
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
	req.Password = strings.TrimSpace(req.Password)

	if req.FullName == "" || req.MobileNumber == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields required-full_name, mobile_number, password"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		return
	}

	// Checks if user already exists
	existingUserID, err := dbHelper.GetUserIDByMobileNumber(req.MobileNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get userID by mobile number"})
		return
	}
	if existingUserID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Mobile Number already exists"})
		return
	}

	// Hash password
	hashed, _ := utils.HashPassword(req.Password)

	var userID string
	txErr := database.Tx(func(tx *sqlx.Tx) error {

		userID, err = dbHelper.CreateUser(tx, req.FullName, req.MobileNumber, hashed)
		if err != nil {
			return err
		}

		err = dbHelper.CreatePlayerCareerStats(tx, userID)
		if err != nil {

			return err
		}

		return nil
	})

	if txErr != nil {
		if txErr.Error() == "invalid userid" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": txErr.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
		"user":    userID,
	})
}

func LoginUser(c *gin.Context) {
	var req models.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request- must send these fields in body-mobile_number and password",
		})
		return
	}

	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
	req.Password = strings.TrimSpace(req.Password)

	if req.MobileNumber == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	user, err := dbHelper.GetUserInfoByMobileNumber(req.MobileNumber)
	if err != nil {
		//fmt.Println("----------", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) //"User Does Not exist"
		return
	}

	IsPasswordMatched := utils.ComparePasswordHash(req.Password, user.Password)
	if !IsPasswordMatched {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	// 🔥 Creating session for the user
	sessionID, err := dbHelper.CreateSession(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	//fmt.Println("Login ----->", user.UserRole)
	token, err := utils.GenerateToken(user.UserID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func LogoutUser(c *gin.Context) {
	sessionID := c.GetString("session_id")

	isSessionActive := dbHelper.IsSessionActive(sessionID)
	if !isSessionActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid SessionID",
		})
		return
	}

	err := dbHelper.ArchiveSession(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid sessionID-unable to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func ForgotPassword(c *gin.Context) {
	var req models.ForgotPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request- must send these fields in body -mobile_number, otp, password"})
		return
	}

	req.Otp = strings.TrimSpace(req.Otp)
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
	req.Password = strings.TrimSpace(req.Password)

	if req.Otp == "" || req.MobileNumber == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required-otp, mobile_number, password"})
		return
	}

	if req.Otp != "8080" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		return
	}

	// Checks if user exists
	existingUserID, err := dbHelper.GetUserIDByMobileNumber(req.MobileNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get userID by mobile number"})
		return
	}
	if existingUserID == "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Mobile Number does not exists"})
		return
	}

	// Hash password
	hashed, _ := utils.HashPassword(req.Password)

	err = dbHelper.UpdatePassword(req.MobileNumber, hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password updated successfully",
	})
}
