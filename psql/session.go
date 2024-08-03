// gpt4o @ 240803

// 文件内容概述：
// Session 结构体：定义了数据库表 sessions 的字段和结构，包括会话 ID、用户名、登录时间、国家、IP 地址和用户代理。
// CreateSession 函数：用于插入新会话到 sessions 表。
// GetSession 函数：根据会话 ID 从 sessions 表中获取会话信息。
// GetSessionList 函数：获取指定用户名最近的 20 个会话记录。
// UpdateSession 函数：更新指定会话的信息。
// DeleteSession 函数：删除指定会话。
// DeleteSessions 函数：删除指定用户的所有会话。
// 使用说明：
// 事务管理：所有数据库操作均依赖于传入的事务 (tx) 对象，确保操作的原子性和一致性。
// 错误处理：每个操作都会返回错误信息，以便调用者能够进行适当的错误处理。

package psql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // 导入 PostgreSQL 驱动
)

// Session 表示 sessions 表的结构体
// 定义了会话表的结构，包括会话 ID、用户名、登录时间、国家、IP 地址和用户代理等
type Session struct {
	SessionID string    `gorm:"primaryKey;type:varchar(255)" json:"session_id"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	LoginTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"login_time"`
	Country   string    `gorm:"type:char(2);not null" json:"country"`
	IPAddress string    `gorm:"type:varchar(45);not null" json:"ip_address"`
	UserAgent string    `gorm:"type:text;not null" json:"user_agent"`
}

// CreateSession 插入一个新的会话到 sessions 表中
// 接受一个事务 tx 和多个参数，用于向 sessions 表插入一个新会话
func CreateSession(tx *sql.Tx, sessionID, username, country, ipAddress, userAgent string) error {
	// 插入会话信息的 SQL 查询语句
	query := `
	INSERT INTO sessions (session_id, username, login_time, country, ip_address, user_agent)
	VALUES ($1, $2, CURRENT_TIMESTAMP, $3, $4, $5)
	`
	// 执行插入操作
	if _, err := tx.Exec(query, sessionID, username, country, ipAddress, userAgent); err != nil {
		return fmt.Errorf("could not create session: %v", err)
	}
	return nil
}

// GetSession 根据会话 ID 获取会话信息
// 接受一个事务 tx 和会话 ID，返回对应的 Session 结构体
func GetSession(tx *sql.Tx, sessionID string) (*Session, error) {
	// 获取会话信息的 SQL 查询语句
	query := `
	SELECT session_id, username, login_time, country, ip_address, user_agent
	FROM sessions
	WHERE session_id = $1
	`
	row := tx.QueryRow(query, sessionID)

	var session Session
	// 扫描数据库返回的结果并赋值给 Session 结构体
	if err := row.Scan(
		&session.SessionID,
		&session.Username,
		&session.LoginTime,
		&session.Country,
		&session.IPAddress,
		&session.UserAgent,
	); err != nil {
		return &session, fmt.Errorf("could not retrieve session: %v", err)
	}

	return &session, nil
}

// GetSessions 根据用户名获取最近的 20 个会话记录
// 接受一个事务 tx 和用户名，返回一个包含最多 20 个会话记录的 Session 结构体切片
func GetSessions(tx *sql.Tx, username string) ([]*Session, error) {
	// 获取会话列表的 SQL 查询语句
	query := `
	SELECT session_id, username, login_time, country, ip_address, user_agent
	FROM sessions
	WHERE username = $1
	ORDER BY login_time DESC
	LIMIT 20
	`
	rows, err := tx.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("could not query sessions: %v", err)
	}
	defer rows.Close()

	// 创建一个切片来存储会话记录
	sessions := make([]*Session, 0, 20)
	for rows.Next() {
		var session Session
		// 扫描每一行并添加到 sessions 切片中
		if err := rows.Scan(
			&session.SessionID,
			&session.Username,
			&session.LoginTime,
			&session.Country,
			&session.IPAddress,
			&session.UserAgent,
		); err != nil {
			return sessions, fmt.Errorf("could not scan session: %v", err)
		}
		sessions = append(sessions, &session)
	}

	// 检查 rows 是否有任何错误
	if err := rows.Err(); err != nil {
		return sessions, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return sessions, nil
}

// UpdateSession 更新会话信息
// 接受一个事务 tx 和多个参数，用于更新指定会话的信息
func UpdateSession(tx *sql.Tx, sessionID, country, ipAddress, userAgent string) error {
	// 更新会话信息的 SQL 查询语句
	query := `
	UPDATE sessions
	SET country = $1, ip_address = $2, user_agent = $3, login_time = CURRENT_TIMESTAMP
	WHERE session_id = $4
	`
	// 执行更新操作
	result, err := tx.Exec(query, country, ipAddress, userAgent, sessionID)
	if err != nil {
		return fmt.Errorf("could not update session: %v", err)
	}
	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no session found with the provided session ID")
	}

	return nil
}

// DeleteSession 删除指定的会话
// 接受一个事务 tx 和会话 ID，删除对应的会话记录
func DeleteSession(tx *sql.Tx, sessionID string) error {
	// 删除会话的 SQL 查询语句
	query := `
	DELETE FROM sessions
	WHERE session_id = $1
	`
	// 执行删除操作
	result, err := tx.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("could not delete session: %v", err)
	}
	// 检查是否有删除的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no session found with the provided session ID")
	}

	return nil
}

// DeleteSessions 删除指定用户的所有会话
// 接受一个事务 tx 和用户名，删除对应用户的所有会话记录
func DeleteSessions(tx *sql.Tx, username string) error {
	// 删除用户所有会话的 SQL 查询语句
	query := `
	DELETE FROM sessions
	WHERE username = $1
	`
	// 执行删除操作
	result, err := tx.Exec(query, username)
	if err != nil {
		return fmt.Errorf("could not delete sessions: %v", err)
	}
	// 检查是否有删除的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no sessions found with the provided username")
	}

	return nil
}
