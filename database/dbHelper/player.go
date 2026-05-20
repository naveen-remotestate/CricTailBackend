package dbHelper

import (
	"CricTail_Backend/database"
	"CricTail_Backend/models"
	"fmt"
)

func GetPlayerStatsByUserID(userID string) (models.PlayerStats, error) {
	var PlayerStats models.PlayerStats

	query := `
		SELECT 
    id,
    user_id,
    batting_style,
    bowling_style,
    matches_played,
    innings_batted,
    innings_bowled,
    total_runs,
    highest_score,
    total_balls_faced,
    total_fours,
    total_sixes,
    total_wickets,
    total_balls_bowled,
    total_runs_conceded,
    total_maidens,
    catches,
    run_outs,
    updated_at
FROM player_career_stats
WHERE user_id = $1;
	`

	err := database.DB.Get(&PlayerStats, query, userID)
	if err != nil {
		return models.PlayerStats{}, err //if error then returning empty player stats
	}

	return PlayerStats, nil
}

func GetAllPlayers(search string) ([]models.Player, error) {
	players := make([]models.Player, 0)

	query := `SELECT user_id,mobile_number, full_name FROM users where is_active=TRUE AND ($1 ='' OR full_name ILIKE '%' || $1 || '%' OR mobile_number ILIKE  '%' || $1 || '%'  )`
	err := database.DB.Select(&players, query, search)
	return players, err
}

func UpdatePlayerProfile(UserID, FullName, BattingStyle, BowlingStyle string) error {

	query1 := `
		UPDATE users
		SET
			full_name = COALESCE(NULLIF($1, ''), full_name)
		WHERE user_id = $2
	`
	rows1, err := database.DB.Exec(
		query1,
		FullName,
		UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to change full name")
	}
	count1, err := rows1.RowsAffected()
	if count1 == 0 {
		return fmt.Errorf("invalid userID") //because even if toodo id is wrong the query will run succsessfully
	}

	query2 := `
		UPDATE player_career_stats
		SET
			batting_style = COALESCE(NULLIF($1, ''), batting_style),
			bowling_style = COALESCE(NULLIF($2, ''), bowling_style)
		WHERE user_id = $3
	`
	rows2, err := database.DB.Exec(
		query2,
		BattingStyle,
		BowlingStyle,
		UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to change batting or bowling style")
	}
	count2, err := rows2.RowsAffected()
	if count2 == 0 {
		return fmt.Errorf("invalid userID") //because even if toodo id is wrong the query will run succsessfully
	}

	return nil
}
