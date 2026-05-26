package dbHelper

import (
	"CricTail_Backend/database"
	"CricTail_Backend/models"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

func GetUserIDBySessionID(sessionID string) (string, error) {
	var userID string

	query := `
		SELECT user_id
		FROM user_sessions
		WHERE id = $1 AND archived_at IS NULL
	`

	err := database.DB.Get(&userID, query, sessionID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetUserIDByMobileNumber(mobileNumber string) (string, error) {
	var user_id string

	query := `SELECT user_id FROM users WHERE mobile_number=$1 AND is_active=TRUE`

	err := database.DB.Get(&user_id, query, mobileNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return user_id, err
}

func GetUserInfoByMobileNumber(mobileNumber string) (models.UserInfo, error) {
	var user models.UserInfo

	query := `
		SELECT user_id, full_name, mobile_number, password_hash
		FROM users
		WHERE mobile_number=$1
	`

	err := database.DB.Get(&user, query, mobileNumber)
	if err != nil {
		return models.UserInfo{}, err //if error then returning empty user
	}

	return user, nil
}

func CreateUser(tx *sqlx.Tx, fullName, mobileNumber, password string) (string, error) {
	query := `
		INSERT INTO users (full_name, mobile_number, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id`

	var userID string
	err := tx.Get(&userID, query, fullName, mobileNumber, password)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func CreatePlayerCareerStats(tx *sqlx.Tx, userID string) error {
	query := `INSERT INTO player_career_stats (user_id)
			VALUES ($1)
`
	res, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("invalid userid")
	}
	return nil
}

func CreateSession(userID string) (string, error) {
	var sessionID string

	query := `
		INSERT INTO user_sessions (user_id)
		VALUES ($1)
		RETURNING id
	`

	err := database.DB.Get(&sessionID, query, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func ArchiveSession(sessionID string) error {
	query := `
		UPDATE user_sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	_, err := database.DB.Exec(query, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func IsSessionActive(sessionID string) bool {
	var archivedAt *time.Time
	query := `SELECT archived_at FROM user_sessions WHERE id=$1 AND archived_at IS NULL`
	err := database.DB.Get(&archivedAt, query, sessionID)
	if err != nil {
		return false
	}
	return true
}

func UpdatePassword(MobileNumber, password string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at=NOW()
		WHERE mobile_number = $2 AND is_active = TRUE
	`

	_, err := database.DB.Exec(query, password, MobileNumber)
	if err != nil {
		return err
	}

	return nil
}
