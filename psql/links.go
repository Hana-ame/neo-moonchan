package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Link 表示 links 表的结构体
type Link struct {
	LinkID     int64  `gorm:"primaryKey;type:bigint" json:"link_id"`
	Username   string `gorm:"type:varchar(50);not null" json:"username"`
	StatusID   int64  `gorm:"type:bigint;not null" json:"status_id"`
	Visibility string `gorm:"type:varchar(10);not null;default:'public'" json:"visibility"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定 Link 结构体对应的数据库表名称
func (Link) TableName() string {
	return "links"
}

// CreateLink 插入一个新的链接到 links 表中
func CreateLink(tx *sql.Tx, linkID int64, username string, statusID int64, visibility string) error {
	query := `
	INSERT INTO links (link_id, username, status_id, visibility)
	VALUES ($1, $2, $3, $4)
`
	if _, err := tx.Exec(query, linkID, username, statusID, visibility); err != nil {
		return fmt.Errorf("could not create link: %v", err)
	}

	return nil
}

// GetLink 根据 LinkID 获取链接信息
func GetLink(tx *sql.Tx, linkID int64) (*Link, error) {
	query := `
	SELECT link_id, username, status_id, visibility, created_at
	FROM links
	WHERE link_id = $1
`
	row := tx.QueryRow(query, linkID)

	var link Link

	if err := row.Scan(
		&link.LinkID,
		&link.Username,
		&link.StatusID,
		&link.Visibility,
		&link.CreatedAt,
	); err != nil {
		return &link, fmt.Errorf("could not retrieve link: %v", err)
	}

	return &link, nil
}

// UpdateLink 更新链接信息，包括所有字段
func UpdateLink(tx *sql.Tx, linkID int64, username string, statusID int64, visibility string) error {
	query := `
		UPDATE links
		SET username = $1, status_id = $2, visibility = $3
		WHERE link_id = $4
	`
	result, err := tx.Exec(query, username, statusID, visibility, linkID)
	if err != nil {
		return fmt.Errorf("could not update link: %v", err)
	}

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no link found with the provided ID")
	}

	return nil
}

// SoftDeleteLink 标记链接为 "deleted" 而不是实际删除
func SoftDeleteLink(tx *sql.Tx, linkID int64) error {
	query := `
		UPDATE links
		SET visibility = 'deleted'
		WHERE link_id = $1
	`
	result, err := tx.Exec(query, linkID)
	if err != nil {
		return fmt.Errorf("could not delete link: %v", err)
	}

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no link found with the provided ID")
	}

	return nil
}

// GetLinks 根据多个 LinkID 获取链接记录
func GetLinks(tx *sql.Tx, linkIDs []int64) ([]*Link, error) {
	if len(linkIDs) == 0 {
		return []*Link{}, nil
	}

	query := `
		SELECT link_id, username, status_id, visibility, created_at
		FROM links
		WHERE link_id = ANY($1) 
	` // instead of IN

	rows, err := tx.Query(query, pq.Array(linkIDs))
	if err != nil {
		return nil, fmt.Errorf("could not query links: %v", err)
	}
	defer rows.Close()

	linksMap := make(map[int64]*Link, 25)
	for rows.Next() {
		var link Link
		if err := rows.Scan(
			&link.LinkID,
			&link.Username,
			&link.StatusID,
			&link.Visibility,
			&link.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan link: %v", err)
		}
		linksMap[link.LinkID] = &link
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	var links []*Link = make([]*Link, 0, len(linksMap))
	// 排序结果以确保与 LinkIDs 列表的顺序一致
	for _, k := range linkIDs {
		if v, exists := linksMap[k]; exists {
			links = append(links, v)
		}
	}

	return links, nil
}

// GetStatusesFromLinksMaxID 根据小于某个 ID 和用户名获取链接记录，并按 ID 倒序排序，返回对应的状态信息
func GetStatusesFromLinks(tx *sql.Tx, username string, limit int) ([]*Status, error) {
	// 第一步：查询符合条件的 status_id 列表
	queryLinks := `
		SELECT status_id
		FROM links
		WHERE username = $1 
		ORDER BY link_id DESC
		LIMIT $2
	`

	if limit <= 0 {
		limit = 25
	}

	rows, err := tx.Query(queryLinks, username, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query links: %v", err)
	}
	defer rows.Close()

	var statusIDs []int64
	for rows.Next() {
		var statusID int64
		if err := rows.Scan(&statusID); err != nil {
			return nil, fmt.Errorf("could not scan status_id: %v", err)
		}
		statusIDs = append(statusIDs, statusID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over links rows: %v", err)
	}

	if len(statusIDs) == 0 {
		// 没有符合条件的记录时，直接返回空切片
		return []*Status{}, nil
	}

	// 第二步：根据 status_id 列表查询状态记录
	queryStatuses := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id = ANY($1)
	`
	statusRows, err := tx.Query(queryStatuses, pq.Array(statusIDs))
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer statusRows.Close()

	// 使用映射表来存储 status_id 与 Status 之间的关系
	statusMap := make(map[int64]*Status, len(statusIDs))
	for statusRows.Next() {
		var status Status
		if err := statusRows.Scan(
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
		statusMap[status.ID] = &status
	}

	if err := statusRows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over statuses rows: %v", err)
	}

	// 第三步：按照初始的 status_id 顺序构建结果切片
	var statuses []*Status = make([]*Status, 0, len(statusIDs))
	for _, id := range statusIDs {
		if status, found := statusMap[id]; found {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}

// GetStatusesFromLinksMaxID 根据小于某个 ID 和用户名获取链接记录，并按 ID 倒序排序，返回对应的状态信息
func GetStatusesFromLinksMaxID(tx *sql.Tx, maxID int64, username string, limit int) ([]*Status, error) {
	// 第一步：查询符合条件的 status_id 列表
	queryLinks := `
		SELECT status_id
		FROM links
		WHERE username = $1 AND link_id < $2
		ORDER BY link_id DESC
		LIMIT $3
	`

	if limit <= 0 {
		limit = 25
	}

	rows, err := tx.Query(queryLinks, username, maxID, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query links: %v", err)
	}
	defer rows.Close()

	var statusIDs []int64
	for rows.Next() {
		var statusID int64
		if err := rows.Scan(&statusID); err != nil {
			return nil, fmt.Errorf("could not scan status_id: %v", err)
		}
		statusIDs = append(statusIDs, statusID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over links rows: %v", err)
	}

	if len(statusIDs) == 0 {
		// 没有符合条件的记录时，直接返回空切片
		return []*Status{}, nil
	}

	// 第二步：根据 status_id 列表查询状态记录
	queryStatuses := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id = ANY($1)
	`
	statusRows, err := tx.Query(queryStatuses, pq.Array(statusIDs))
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer statusRows.Close()

	// 使用映射表来存储 status_id 与 Status 之间的关系
	statusMap := make(map[int64]*Status, len(statusIDs))
	for statusRows.Next() {
		var status Status
		if err := statusRows.Scan(
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
		statusMap[status.ID] = &status
	}

	if err := statusRows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over statuses rows: %v", err)
	}

	// 第三步：按照初始的 status_id 顺序构建结果切片
	var statuses []*Status = make([]*Status, 0, len(statusIDs))
	for _, id := range statusIDs {
		if status, found := statusMap[id]; found {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}

// GetStatusesFromLinksMinID 根据大于某个 ID 和用户名获取链接记录，并按 ID 升序排序，返回对应的状态信息
func GetStatusesFromLinksMinID(tx *sql.Tx, minID int64, username string, limit int) ([]*Status, error) {
	// 第一步：查询符合条件的 status_id 列表
	queryLinks := `
			SELECT status_id
			FROM links
			WHERE username = $1 AND link_id > $2
			ORDER BY link_id ASC
			LIMIT $3
	`

	if limit <= 0 {
		limit = 25
	}

	rows, err := tx.Query(queryLinks, username, minID, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query links: %v", err)
	}
	defer rows.Close()

	var statusIDs []int64
	for rows.Next() {
		var statusID int64
		if err := rows.Scan(&statusID); err != nil {
			return nil, fmt.Errorf("could not scan status_id: %v", err)
		}
		statusIDs = append(statusIDs, statusID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over links rows: %v", err)
	}

	if len(statusIDs) == 0 {
		// 没有符合条件的记录时，直接返回空切片
		return []*Status{}, nil
	}

	// 第二步：根据 status_id 列表查询状态记录
	queryStatuses := `
		SELECT id, username, warning, content, visibility, created_at, updated_at
		FROM statuses
		WHERE id = ANY($1)
	`
	statusRows, err := tx.Query(queryStatuses, pq.Array(statusIDs))
	if err != nil {
		return nil, fmt.Errorf("could not query statuses: %v", err)
	}
	defer statusRows.Close()

	// 使用映射表来存储 status_id 与 Status 之间的关系
	statusMap := make(map[int64]*Status, len(statusIDs))
	for statusRows.Next() {
		var status Status
		if err := statusRows.Scan(
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
		statusMap[status.ID] = &status
	}

	if err := statusRows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over statuses rows: %v", err)
	}

	// 第三步：按照初始的 status_id 顺序构建结果切片
	var statuses []*Status = make([]*Status, 0, len(statusIDs))
	for _, id := range statusIDs {
		if status, found := statusMap[id]; found {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}
