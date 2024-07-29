package psql

import (
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	SessionID string    `gorm:"primaryKey;type:varchar(255)"`
	Username  string    `gorm:"type:varchar(50);not null"`
	LoginTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Country   string    `gorm:"type:char(2);not null"`
	IPAddress string    `gorm:"type:varchar(45);not null"`
	UserAgent string    `gorm:"type:text;not null"`
}

func CreateSession(tx *sql.Tx, sessionID, username, country, ipAddress, userAgent string) error {
	query := `
	INSERT INTO sessions (session_id, username, login_time, country, ip_address, user_agent)
	VALUES ($1, $2, CURRENT_TIMESTAMP, $3, $4, $5)
`
	if _, err := tx.Exec(query, sessionID, username, country, ipAddress, userAgent); err != nil {
		return fmt.Errorf("could not create session: %v", err)
	}
	return nil
}

func GetSession(tx *sql.Tx, sessionID string) (*Session, error) {
	query := `
	SELECT session_id, username, login_time, country, ip_address, user_agent
	FROM sessions
	WHERE session_id = $1
`
	row := tx.QueryRow(query, sessionID)

	var session Session
	if err := row.Scan(
		&session.SessionID,
		&session.Username,
		&session.LoginTime,
		&session.Country,
		&session.IPAddress,
		&session.UserAgent,
	); err != nil {
		return nil, fmt.Errorf("could not retrieve session: %v", err)
	}

	return &session, nil
}

func GetSessionList(tx *sql.Tx, username string) ([]*Session, error) {
	query := `
	SELECT session_id, username, login_time, country, ip_address, user_agent
	FROM sessions
	WHERE username = $1
	LIMIT 20
`
	rows, err := tx.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("could not query sessions: %v", err)
	}
	defer rows.Close()

	sessions := make([]*Session, 0, 20)
	for rows.Next() {
		var session Session
		if err := rows.Scan(
			&session.SessionID,
			&session.Username,
			&session.LoginTime,
			&session.Country,
			&session.IPAddress,
			&session.UserAgent,
		); err != nil {
			return nil, fmt.Errorf("could not scan session: %v", err)
		}
		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return sessions, nil
}

func UpdateSession(tx *sql.Tx, sessionID, country, ipAddress, userAgent string) error {
	query := `
	UPDATE sessions
	SET country = $1, ip_address = $2, user_agent = $3, login_time = CURRENT_TIMESTAMP
	WHERE session_id = $4
`
	result, err := tx.Exec(query, country, ipAddress, userAgent, sessionID)
	if err != nil {
		return fmt.Errorf("could not update session: %v", err)
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no session found with the provided session ID")
	}

	return nil
}

func DeleteSession(tx *sql.Tx, sessionID string) error {
	query := `
	DELETE FROM sessions
	WHERE session_id = $1
`
	result, err := tx.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("could not delete session: %v", err)
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no session found with the provided session ID")
	}

	return nil
}
