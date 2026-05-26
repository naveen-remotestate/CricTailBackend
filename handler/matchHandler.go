package handler

import (
	"CricTail_Backend/database"
	"CricTail_Backend/database/dbHelper"
	"CricTail_Backend/models"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateMatch(c *gin.Context) {

	var req models.CreateMatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// validate

	req.TeamA.Name = strings.TrimSpace(req.TeamA.Name)
	req.TeamB.Name = strings.TrimSpace(req.TeamB.Name)

	if req.TeamA.Name == "" || req.TeamB.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "team names required",
		})
		return
	}

	if req.TeamA.Name == req.TeamB.Name {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "both teams cannot have same name",
		})
		return
	}

	if req.Overs <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "overs must be greater than 0",
		})
		return
	}

	if req.TossWinnerTeam != "A" && req.TossWinnerTeam != "B" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid toss_winner_team",
		})
		return
	}

	if req.TossDecision != "BAT" && req.TossDecision != "BOWL" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid toss_decision",
		})
		return
	}

	// player validations for captain and one same player in both teamss
	allPlayers := make(map[string]bool)
	samePlayerCount := 0

	validateTeam := func(team models.TeamInput) error {

		if len(team.Players) == 0 {
			return fmt.Errorf("team players required")
		}

		captainCount := 0

		for _, player := range team.Players {

			if player.UserID == "" {
				return fmt.Errorf("user_id required")
			}

			if allPlayers[player.UserID] {

				samePlayerCount++

				if samePlayerCount > 1 {
					return fmt.Errorf("only one same player allowed in both teams")
				}
			}

			allPlayers[player.UserID] = true

			if player.IsCaptain {
				captainCount++
			}
		}

		if captainCount != 1 {
			return fmt.Errorf("exactly one captain required")
		}

		return nil
	}
	if err := validateTeam(req.TeamA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validateTeam(req.TeamB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var matchID string

	txErr := database.Tx(func(tx *sqlx.Tx) error {

		// create team A

		teamAID, err := dbHelper.CreateTeam(
			tx,
			req.TeamA.Name,
			req.HostedBy,
		)
		if err != nil {
			return err
		}

		for _, player := range req.TeamA.Players {

			err = dbHelper.AddPlayerToTeam(
				tx,
				teamAID,
				player,
			)
			if err != nil {
				return err
			}
		}

		// create team b

		teamBID, err := dbHelper.CreateTeam(
			tx,
			req.TeamB.Name,
			req.HostedBy,
		)
		if err != nil {
			return err
		}

		for _, player := range req.TeamB.Players {

			err = dbHelper.AddPlayerToTeam(
				tx,
				teamBID,
				player,
			)
			if err != nil {
				return err
			}
		}

		// who bat first

		var battingFirstTeamID string
		var bowlingFirstTeamID string

		if req.TossWinnerTeam == "A" {

			if req.TossDecision == "BAT" {
				battingFirstTeamID = teamAID
				bowlingFirstTeamID = teamBID
			} else {
				battingFirstTeamID = teamBID
				bowlingFirstTeamID = teamAID
			}

		} else {

			if req.TossDecision == "BAT" {
				battingFirstTeamID = teamBID
				bowlingFirstTeamID = teamAID
			} else {
				battingFirstTeamID = teamAID
				bowlingFirstTeamID = teamBID
			}
		}

		tossWinnerID := teamAID
		if req.TossWinnerTeam == "B" {
			tossWinnerID = teamBID
		}

		matchID, err = dbHelper.CreateMatch(
			tx,
			teamAID,
			teamBID,
			tossWinnerID,
			req.TossDecision,
			battingFirstTeamID,
			req.Overs,
			req.HostedBy,
		)
		if err != nil {
			return err
		}

		//create innings table
		InningID, err := dbHelper.CreateInning(
			tx,
			matchID,
			"1",
			battingFirstTeamID,
			bowlingFirstTeamID,
		)
		if err != nil {
			return err
		}

		//create Live Match table
		err = dbHelper.CreateLiveMatch(
			tx,
			matchID,
			InningID,
			req.StrikerID,
			req.NonStrikerID,
			req.BowlerID,
		)
		if err != nil {
			return err
		}
		if battingFirstTeamID == teamAID {
			//for batting scorecard of team batting first
			for _, player := range req.TeamA.Players {

				err = dbHelper.CreateBattingScorecard(
					tx, InningID, player.UserID,
				)
				if err != nil {
					return err
				}
			}
			//for bowling scorecard of team bowling first
			for _, player := range req.TeamB.Players {

				err = dbHelper.CreateBowlingScorecard(
					tx, InningID, player.UserID,
				)
				if err != nil {
					return err
				}
			}

		} else {
			for _, player := range req.TeamB.Players {

				err = dbHelper.CreateBattingScorecard(
					tx, InningID, player.UserID,
				)
				if err != nil {
					return err
				}
			}

			//for bowling scorecard of team bowling first
			for _, player := range req.TeamA.Players {

				err = dbHelper.CreateBowlingScorecard(
					tx, InningID, player.UserID,
				)
				if err != nil {
					return err
				}
			}

		}

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": txErr.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "match created successfully",
		"match_id": matchID,
	})
}

func GetMatches(c *gin.Context) {
	matches, err := dbHelper.GetMatches()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(), //"failed to get matches"
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
	})
}

func GetMatchByID(c *gin.Context) {

	matchID := c.Param("matchID")

	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "matchID required",
		})
		return
	}

	match, err := dbHelper.GetMatchByID(matchID)
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "match not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match": match,
	})
}

func ScoreLiveMatch(c *gin.Context) {

}
