package dbHelper

import (
	"CricTail_Backend/database"
	"CricTail_Backend/models"
)

func GetPlayerStatsByUserID(userID string) (models.PlayerStats, error) {
	var PlayerStats models.PlayerStats

	query := `
		SELECT 
    id,
    user_id,
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
