// gpt4o mini @ 240801
// passed.
// claude.ai @ 240803

// Package psql provides functions for managing status records in a PostgreSQL database.

package psql

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Status represents the structure of the 'statuses' table.
type Status struct {
	ID         int64     `gorm:"primaryKey;type:bigint" json:"id"`
	Username   string    `gorm:"type:varchar(50);not null" json:"username"`
	Warning    string    `gorm:"type:text;not null" json:"warning"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Visibility string    `gorm:"type:varchar(10);not null;default:'public'" json:"visibility"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateStatus inserts a new status into the 'statuses' table.
func CreateStatus(tx *sql.Tx, id int64, username, warning, content, visibility string) error {
	query := `
	INSERT INTO statuses (id, username, warning, content, visibility)
	VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := tx.Exec(query, id, username, warning, content, visibility); err != nil {
		return fmt.Errorf("could not create status: %v", err)
	}
	return nil
}

// GetStatus retrieves a status by its ID.
func GetStatus(tx *sql.Tx, id int64) (*Status, error) {
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	WHERE id = $1
	`
	row := tx.QueryRow(query, id)
	var status Status
	if err := row.Scan(&status.ID, &status.Username, &status.Warning, &status.Content,
		&status.Visibility, &status.CreatedAt, &status.UpdatedAt); err != nil {
		return &status, fmt.Errorf("could not retrieve status: %v", err)
	}
	return &status, nil
}

// UpdateStatus updates an existing status.
func UpdateStatus(tx *sql.Tx, id int64, warning, content string) error {
	query := `
	UPDATE statuses
	SET warning = $1, content = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	`
	result, err := tx.Exec(query, warning, content, id)
	if err != nil {
		return fmt.Errorf("could not update status: %v", err)
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no status found with the provided ID")
	}
	return nil
}

// SoftDeleteStatus marks a status as 'deleted' instead of actually deleting it.
func SoftDeleteStatus(tx *sql.Tx, id int64) error {
	query := `
	UPDATE statuses
	SET visibility = 'deleted', updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	`
	result, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete status: %v", err)
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no status found with the provided ID")
	}
	return nil
}

// GetStatusesByIds retrieves multiple statuses by their IDs.
func GetStatusesByIds(tx *sql.Tx, ids []int64) ([]*Status, error) {
	if len(ids) == 0 {
		return []*Status{}, nil
	}
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	WHERE id = ANY($1)
	`
	rows, err := tx.Query(query, pq.Array(ids))
	if err != nil {
		return []*Status{}, fmt.Errorf("could not query statuses: %v", err)
	}
	defer rows.Close()

	statusesMap := make(map[int64]*Status, len(ids))
	for rows.Next() {
		var status Status
		if err := rows.Scan(&status.ID, &status.Username, &status.Warning, &status.Content,
			&status.Visibility, &status.CreatedAt, &status.UpdatedAt); err != nil {
			return []*Status{}, fmt.Errorf("could not scan status: %v", err)
		}
		statusesMap[status.ID] = &status
	}
	if err := rows.Err(); err != nil {
		return []*Status{}, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	statuses := make([]*Status, 0, len(statusesMap))
	for _, k := range ids {
		if v, exists := statusesMap[k]; exists {
			statuses = append(statuses, v)
		}
	}
	return statuses, nil
}

// GetStatusesByMaxID retrieves statuses with IDs less than maxID, ordered by ID descending.
func GetStatusesByMaxID(tx *sql.Tx, maxID int64, limit int) ([]*Status, error) {
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	WHERE id < $1
	ORDER BY id DESC
	LIMIT $2
	`
	limit = intInterval(limit, 1, 25)
	return queryStatuses(tx, query, maxID, limit)
}

// GetStatusesByMinID retrieves statuses with IDs greater than minID, ordered by ID ascending.
func GetStatusesByMinID(tx *sql.Tx, minID int64, limit int) ([]*Status, error) {
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	WHERE id > $1
	ORDER BY id ASC
	LIMIT $2
	`
	limit = intInterval(limit, 1, 25)
	statuses, err := queryStatuses(tx, query, minID, limit)
	if err != nil {
		return statuses, err
	}
	slices.Reverse(statuses)
	return statuses, nil
}

// GetLatestStatuses retrieves the most recent statuses, ordered by ID descending.
func GetLatestStatuses(tx *sql.Tx, limit int) ([]*Status, error) {
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	ORDER BY id DESC
	LIMIT $1
	`
	limit = intInterval(limit, 1, 25)
	return queryStatuses(tx, query, limit)
}

// queryStatuses is a helper function to execute status queries and scan results.
func queryStatuses(tx *sql.Tx, query string, args ...interface{}) ([]*Status, error) {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return []*Status{}, fmt.Errorf("could not query statuses: %v", err)
	}
	defer rows.Close()

	statuses := make([]*Status, 0, 25)
	for rows.Next() {
		var status Status
		if err := rows.Scan(&status.ID, &status.Username, &status.Warning, &status.Content,
			&status.Visibility, &status.CreatedAt, &status.UpdatedAt); err != nil {
			return statuses, fmt.Errorf("could not scan status: %v", err)
		}
		statuses = append(statuses, &status)
	}
	if err := rows.Err(); err != nil {
		return statuses, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}
	return statuses, nil
}
