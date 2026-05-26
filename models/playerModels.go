package models

type PlayerStats struct {
	ID     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`

	BattingStyle string `db:"batting_style" json:"batting_style"`
	BowlingStyle string `db:"bowling_style" json:"bowling_style"`

	MatchesPlayed int `db:"matches_played" json:"matches_played"`

	InningsBatted int `db:"innings_batted" json:"innings_batted"`
	InningsBowled int `db:"innings_bowled" json:"innings_bowled"`

	TotalRuns       int `db:"total_runs" json:"total_runs"`
	HighestScore    int `db:"highest_score" json:"highest_score"`
	TotalBallsFaced int `db:"total_balls_faced" json:"total_balls_faced"`
	TotalFours      int `db:"total_fours" json:"total_fours"`
	TotalSixes      int `db:"total_sixes" json:"total_sixes"`

	TotalWickets      int `db:"total_wickets" json:"total_wickets"`
	TotalBallsBowled  int `db:"total_balls_bowled" json:"total_balls_bowled"`
	TotalRunsConceded int `db:"total_runs_conceded" json:"total_runs_conceded"`
	TotalMaidens      int `db:"total_maidens" json:"total_maidens"`

	Catches int `db:"catches" json:"catches"`
	RunOuts int `db:"run_outs" json:"run_outs"`

	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type Player struct {
	UserID       string `db:"user_id" json:"user_id"`
	FullName     string `db:"full_name" json:"full_name"`
	MobileNumber string `db:"mobile_number" json:"mobile_number"`
}

type UpdatePlayer struct {
	FullName     string `db:"full_name" json:"full_name"`
	BattingStyle string `db:"batting_style" json:"batting_style"`
	BowlingStyle string `db:"bowling_style" json:"bowling_style"`
}
