// gpt4o mini @ 240801
// passed.

package psql

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Status 表示 statuses 表的结构体
type Status struct {
	ID         int64  `gorm:"primaryKey;type:bigint" json:"id"`
	Username   string `gorm:"type:varchar(50);not null" json:"username"`
	Warning    string `gorm:"type:text;not null" json:"warning"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Visibility string `gorm:"type:varchar(10);not null;default:'public'" json:"visibility"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateStatus 插入一个新的状态到 statuses 表中
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

// GetStatus 根据 ID 获取状态信息
func GetStatus(tx *sql.Tx, id int64) (*Status, error) {
	query := `
	SELECT id, username, warning, content, visibility, created_at, updated_at
	FROM statuses
	WHERE id = $1
`
	row := tx.QueryRow(query, id)

	var status Status

	if err := row.Scan(
		&status.ID,
		&status.Username,
		&status.Warning,
		&status.Content,
		&status.Visibility,
		&status.CreatedAt,
		&status.UpdatedAt,
	); err != nil {
		return &status, fmt.Errorf("could not retrieve status: %v", err)
	}

	return &status, nil
}

// UpdateStatus 更新状态信息
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

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no status found with the provided ID")
	}

	return nil
}

// SoftDeleteStatus 标记状态为 "deleted" 而不是实际删除
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

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no status found with the provided ID")
	}

	return nil
}

// GetStatuses 根据多个 ID 获取状态记录
func GetStatuses(tx *sql.Tx, ids []int64) ([]*Status, error) {
	if len(ids) == 0 {
		return []*Status{}, nil
	}

	query := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id = ANY($1) 
	` // instead of IN

	rows, err := tx.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer rows.Close()

	statusesMap := make(map[int64]*Status, 25)
	for rows.Next() {
		var status Status
		if err := rows.Scan(
			&status.ID,
			&status.Username,
			&status.Warning,
			&status.Content,
			&status.Visibility,
			&status.CreatedAt,
			&status.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan status: %v", err)
		}
		statusesMap[status.ID] = &status
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	var statuses []*Status = make([]*Status, 0, len(statusesMap))
	// 排序结果以确保与 IDs 列表的顺序一致
	for _, k := range ids {
		if v, exists := statusesMap[k]; exists {
			statuses = append(statuses, v)
		}
	}

	return statuses, nil
}

func GetStatusesMaxID(tx *sql.Tx, maxID int64, limit int) ([]*Status, error) {
	query := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id < $1
		ORDER BY id DESC
		LIMIT $2
	`
	if limit <= 0 {
		limit = 25
	}

	rows, err := tx.Query(query, maxID, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer rows.Close()

	var statuses []*Status
	for rows.Next() {
		var status Status
		if err := rows.Scan(
			&status.ID,
			&status.Username,
			&status.Warning,
			&status.Content,
			&status.Visibility,
			&status.CreatedAt,
			&status.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan status: %v", err)
		}
		statuses = append(statuses, &status)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	return statuses, nil
}

func GetStatusesMinID(tx *sql.Tx, minID int64, limit int) ([]*Status, error) {
	query := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2
	`
	if limit <= 0 {
		limit = 25
	}

	rows, err := tx.Query(query, minID, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer rows.Close()

	var statuses []*Status
	for rows.Next() {
		var status Status
		if err := rows.Scan(
			&status.ID,
			&status.Username,
			&status.Warning,
			&status.Content,
			&status.Visibility,
			&status.CreatedAt,
			&status.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan status: %v", err)
		}
		statuses = append(statuses, &status)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	// reverse
	slices.Reverse(statuses)

	return statuses, nil
}
