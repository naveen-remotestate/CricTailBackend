package models

//type PlayerStats struct {
//	ID          string     `db:"id" json:"id"`
//	UserID          string     `db:"user_id" json:"user_id"`
//	MatchesPlayed          string     `db:"matches_played" json:"matches_played"`
//	InningsBatted string `db:"innings_batted" json:"innings_batted"`
//	InningsBowled string `db:"innings_bowled" json:"innings_bowled"`
//	TotalRuns string `db:"total_runs" json:"total_runs"`
//
//
//}

type PlayerStats struct {
	ID     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`

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
