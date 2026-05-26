package dbHelper

import (
	"CricTail_Backend/database"
	"CricTail_Backend/models"

	"github.com/jmoiron/sqlx"
)

func CreateTeam(
	tx *sqlx.Tx,
	teamName string,
	hostedBy string,
) (string, error) {

	query := `
		INSERT INTO teams (
			name,
			created_by
		)
		VALUES ($1, $2)
		RETURNING id
	`

	var teamID string

	err := tx.Get(
		&teamID,
		query,
		teamName,
		hostedBy,
	)
	if err != nil {
		return "", err
	}

	return teamID, nil
}

func AddPlayerToTeam(
	tx *sqlx.Tx,
	teamID string,
	player models.TeamPlayerInput,
) error {

	query := `
		INSERT INTO team_players (
			team_id,
			user_id,
			is_captain
		)
		VALUES ($1, $2, $3)
	`

	_, err := tx.Exec(
		query,
		teamID,
		player.UserID,
		player.IsCaptain,
	)

	return err
}

func CreateMatch(
	tx *sqlx.Tx,
	teamAID string,
	teamBID string,
	tossWinnerID string,
	tossDecision string,
	battingFirstTeamID string,
	overs int,
	hostedBy string,
) (string, error) {

	query := `
		INSERT INTO matches (
			team_a_id,
			team_b_id,
			toss_winner_team_id,
			toss_decision,
			batting_first_team_id,
			overs,
			hosted_by
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7
		)
		RETURNING id
	`

	var matchID string

	err := tx.Get(
		&matchID,
		query,
		teamAID,
		teamBID,
		tossWinnerID,
		tossDecision,
		battingFirstTeamID,
		overs,
		hostedBy,
	)

	if err != nil {
		return "", err
	}

	return matchID, nil
}

func CreateInning(tx *sqlx.Tx, matchID, inningNumber, battingFirstTeamID, bowlingFirstTeamID string) (string, error) {
	query := `
		INSERT INTO innings (
			match_id,
			innings_no,
		    batting_team_id,
		    bowling_team_id
		)
		VALUES (
			$1,$2,$3,$4
		)
		RETURNING id
	`

	var InningID string

	err := tx.Get(
		&InningID,
		query,
		matchID,
		inningNumber,
		battingFirstTeamID,
		bowlingFirstTeamID,
	)

	if err != nil {
		return "", err
	}

	return InningID, nil
}

func CreateBattingScorecard(tx *sqlx.Tx, inningID, userID string) error {
	query := `
		INSERT INTO batting_scorecards (
			innings_id,
			user_id
		)
		VALUES (
			$1,$2
		)`

	_, err := tx.Exec(
		query,
		inningID,
		userID,
	)

	if err != nil {
		return err
	}
	return nil

}

func CreateBowlingScorecard(tx *sqlx.Tx, inningID, userID string) error {
	query := `
		INSERT INTO bowling_scorecards (
			innings_id,
			user_id
		)
		VALUES (
			$1,$2
		)`

	_, err := tx.Exec(
		query,
		inningID,
		userID,
	)

	if err != nil {
		return err
	}
	return nil

}

func CreateLiveMatch(tx *sqlx.Tx, matchID, InningID, StrikerID, NonStrikerID, BowlerID string) error {
	query := `
		INSERT INTO live_match (
			match_id,
			innings_id,
			striker_id,
		    non_striker_id,
		    current_bowler_id
		)
		VALUES (
			$1,$2,$3,$4,$5
		)`

	_, err := tx.Exec(
		query,
		matchID,
		InningID,
		StrikerID,
		NonStrikerID,
		BowlerID,
	)

	if err != nil {
		return err
	}
	return nil
}

func GetMatches() ([]models.MatchResponse, error) {
	query := `
		SELECT
			--match
			m.id AS match_id,

			m.toss_winner_team_id,
			m.winner_team_id,

			m.toss_decision,
			m.hosted_by,

			m.current_innings_no,

			m.overs,

			m.start_time,
			m.end_time,

			-- team a
			ta.id AS team_a_id,
			ta.name AS team_a_name,

			-- TEAM B
			tb.id AS team_b_id,
			tb.name AS team_b_name,

			-- live score
			lm.total_runs AS current_total_runs,
			lm.total_wickets AS current_total_wickets,
			lm.legal_balls,

			-- current inning
			i.id AS current_inning_id,

			-- previous iinnng
			pi.total_runs AS previous_innings_score,
			pi.legal_balls AS previous_innings_legal_balls,

			-- striker
			s.user_id AS striker_id,
			s.full_name AS striker_name,

			sb.runs AS striker_runs,
			sb.balls_faced AS striker_balls,

			-- non striker
			ns.user_id AS non_striker_id,
			ns.full_name AS non_striker_name,

			nsb.runs AS non_striker_runs,
			nsb.balls_faced AS non_striker_balls,

			-- bowler
			b.user_id AS bowler_id,
			b.full_name AS bowler_name,

			bs.runs_conceded AS bowler_runs_given,
			bs.legal_balls AS bowler_legal_balls,
			bs.wickets AS bowler_wickets

		FROM matches m

		-- teams
		INNER JOIN teams ta
			ON ta.id = m.team_a_id

		INNER JOIN teams tb
			ON tb.id = m.team_b_id

		-- live match
		LEFT JOIN live_match lm
			ON lm.match_id = m.id

		-- current inning
		LEFT JOIN innings i
			ON i.match_id = m.id
			AND i.innings_no = m.current_innings_no

		-- previous inning
		LEFT JOIN innings pi
			ON pi.match_id = m.id
			AND pi.innings_no = m.current_innings_no - 1

		-- striker
		LEFT JOIN users s
			ON s.user_id = lm.striker_id

		LEFT JOIN batting_scorecards sb
			ON sb.innings_id = i.id
			AND sb.user_id = s.user_id

		-- non striker
		LEFT JOIN users ns
			ON ns.user_id = lm.non_striker_id

		LEFT JOIN batting_scorecards nsb
			ON nsb.innings_id = i.id
			AND nsb.user_id = ns.user_id

		-- bowler
		LEFT JOIN users b
			ON b.user_id = lm.current_bowler_id

		LEFT JOIN bowling_scorecards bs
			ON bs.innings_id = i.id
			AND bs.user_id = b.user_id

		ORDER BY m.created_at DESC
	`

	var matches []models.MatchResponse

	err := database.DB.Select(
		&matches,
		query,
	)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func GetMatchByID(matchID string) (*models.MatchResponse, error) {

	query := `
		SELECT

			-- MATCH
			m.id AS match_id,

			m.toss_winner_team_id,
			m.winner_team_id,

			m.toss_decision,
			m.hosted_by,

			m.current_innings_no,

			m.overs,

			m.start_time,
			m.end_time,

			-- TEAM A
			ta.id AS team_a_id,
			ta.name AS team_a_name,

			-- TEAM B
			tb.id AS team_b_id,
			tb.name AS team_b_name,

			-- LIVE SCORE
			COALESCE(lm.total_runs, 0) AS current_total_runs,
			COALESCE(lm.total_wickets, 0) AS current_total_wickets,
			COALESCE(lm.legal_balls, 0) AS legal_balls,

			-- CURRENT INNING
			i.id AS current_inning_id,

			-- PREVIOUS INNING
			COALESCE(pi.total_runs, 0) AS previous_innings_score,
			COALESCE(pi.legal_balls, 0) AS previous_innings_legal_balls,

			-- STRIKER
			s.user_id AS striker_id,
			s.full_name AS striker_name,

			COALESCE(sb.runs, 0) AS striker_runs,
			COALESCE(sb.balls_faced, 0) AS striker_balls,

			-- NON STRIKER
			ns.user_id AS non_striker_id,
			ns.full_name AS non_striker_name,

			COALESCE(nsb.runs, 0) AS non_striker_runs,
			COALESCE(nsb.balls_faced, 0) AS non_striker_balls,

			-- BOWLER
			b.user_id AS bowler_id,
			b.full_name AS bowler_name,

			COALESCE(bs.runs_conceded, 0) AS bowler_runs_given,
			COALESCE(bs.legal_balls, 0) AS bowler_legal_balls,
			COALESCE(bs.wickets, 0) AS bowler_wickets

		FROM matches m

		-- TEAMS
		INNER JOIN teams ta
			ON ta.id = m.team_a_id

		INNER JOIN teams tb
			ON tb.id = m.team_b_id

		-- LIVE MATCH
		LEFT JOIN live_match lm
			ON lm.match_id = m.id

		-- CURRENT INNING
		LEFT JOIN innings i
			ON i.match_id = m.id
			AND i.innings_no = m.current_innings_no

		-- PREVIOUS INNING
		LEFT JOIN innings pi
			ON pi.match_id = m.id
			AND pi.innings_no = m.current_innings_no - 1

		-- STRIKER
		LEFT JOIN users s
			ON s.user_id = lm.striker_id

		LEFT JOIN batting_scorecards sb
			ON sb.innings_id = i.id
			AND sb.user_id = s.user_id

		-- NON STRIKER
		LEFT JOIN users ns
			ON ns.user_id = lm.non_striker_id

		LEFT JOIN batting_scorecards nsb
			ON nsb.innings_id = i.id
			AND nsb.user_id = ns.user_id

		-- BOWLER
		LEFT JOIN users b
			ON b.user_id = lm.current_bowler_id

		LEFT JOIN bowling_scorecards bs
			ON bs.innings_id = i.id
			AND bs.user_id = b.user_id

		WHERE m.id = $1
	`

	var match models.MatchResponse

	err := database.DB.Get(
		&match,
		query,
		matchID,
	)
	if err != nil {
		return nil, err
	}

	return &match, nil
}
