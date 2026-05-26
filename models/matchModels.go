package models

import "time"

type TeamPlayerInput struct {
	UserID    string `json:"user_id"`
	IsCaptain bool   `json:"is_captain"`
}

type TeamInput struct {
	Name    string            `json:"name"`
	Players []TeamPlayerInput `json:"players"`
}

type CreateMatchRequest struct {
	TeamA TeamInput `json:"team_a"`
	TeamB TeamInput `json:"team_b"`

	Overs    int    `json:"overs"`
	HostedBy string `json:"hosted_by"`

	TossWinnerTeam string `json:"toss_winner_team"` // A or B
	TossDecision   string `json:"toss_decision"`    // BAT or BOWL

	//for innings table

	//for live match table
	StrikerID    string `json:"striker_id"`
	NonStrikerID string `json:"non_striker_id"`
	BowlerID     string `json:"current_bowler_id"`
}

type MatchResponse struct {
	// match
	MatchID string `db:"match_id" json:"match_id"`

	TossWinnerTeamID *string `db:"toss_winner_team_id" json:"toss_winner_team_id"`
	WinnerTeamID     *string `db:"winner_team_id" json:"winner_team_id"`

	TossDecision string `db:"toss_decision" json:"toss_decision"`
	HostID       string `db:"hosted_by" json:"hosted_by"`

	CurrentInningNo int `db:"current_innings_no" json:"current_innings_no"`

	Overs int `db:"overs" json:"overs"`

	StartTime *time.Time `db:"start_time" json:"start_time"`
	EndTime   *time.Time `db:"end_time" json:"end_time"`

	// Team A
	TeamAID   string `db:"team_a_id" json:"team_a_id"`
	TeamAName string `db:"team_a_name" json:"team_a_name"`

	// Team B
	TeamBID   string `db:"team_b_id" json:"team_b_id"`
	TeamBName string `db:"team_b_name" json:"team_b_name"`

	// Live Score
	CurrentTotalRuns    *int `db:"current_total_runs" json:"current_total_runs"`
	CurrentTotalWickets *int `db:"current_total_wickets" json:"current_total_wickets"`
	LegalBalls          *int `db:"legal_balls" json:"legal_balls"`

	PreviousInningsScore      *int `db:"previous_innings_score" json:"previous_innings_score"`
	PreviousInningsLegalBalls *int `db:"previous_innings_legal_balls" json:"previous_innings_legal_balls"`

	CurrentInningID *string `db:"current_inning_id" json:"current_inning_id"`

	// Striker details
	StrikerID    *string `db:"striker_id" json:"striker_id"`
	StrikerName  *string `db:"striker_name" json:"striker_name"`
	StrikerRuns  *int    `db:"striker_runs" json:"striker_runs"`
	StrikerBalls *int    `db:"striker_balls" json:"striker_balls"`

	// nonstriker details
	NonStrikerID    *string `db:"non_striker_id" json:"non_striker_id"`
	NonStrikerName  *string `db:"non_striker_name" json:"non_striker_name"`
	NonStrikerRuns  *int    `db:"non_striker_runs" json:"non_striker_runs"`
	NonStrikerBalls *int    `db:"non_striker_balls" json:"non_striker_balls"`

	// current bowler details
	BowlerID         *string `db:"bowler_id" json:"bowler_id"`
	BowlerName       *string `db:"bowler_name" json:"bowler_name"`
	BowlerRunsGiven  *int    `db:"bowler_runs_given" json:"bowler_runs_given"`
	BowlerLegalBalls *int    `db:"bowler_legal_balls" json:"bowler_legal_balls"`
	BowlerWickets    *int    `db:"bowler_wickets" json:"bowler_wickets"`
}
