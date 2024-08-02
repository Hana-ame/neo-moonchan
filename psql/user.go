// gpt4o @ 240801

package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

// User 表示 users 表的结构体
type User struct {
	Username    string                 `gorm:"primaryKey;type:varchar(50);not null" json:"username"`
	Domain      string                 `gorm:"type:varchar(100)" json:"domain"` // 新增的 domain 字段
	DisplayName string                 `gorm:"type:varchar(50)" json:"display_name"`
	AvatarURL   string                 `gorm:"type:varchar(255)" json:"avatar_url"`
	Settings    *orderedmap.OrderedMap `gorm:"type:json;not null;default:'{}'" json:"settings"`
	Flag        string                 `gorm:"type:varchar(50);not null;default:'active'" json:"flag"`
	CreatedAt   time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateUser 插入一个新的用户到 users 表中
func CreateUser(tx *sql.Tx, username, domain, displayName, avatarURL string, settings *orderedmap.OrderedMap) error {
	// 将 settings 转换为 JSON 字符串
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("could not marshal settings: %v", err)
	}

	query := `
	INSERT INTO users (username, domain, display_name, avatar_url, settings, flag)
	VALUES ($1, $2, $3, $4, $5, 'active')
	`
	if _, err := tx.Exec(query, username, domain, displayName, avatarURL, settingsJSON); err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}

	return nil
}

// GetUser 根据用户名获取用户信息
func GetUser(tx *sql.Tx, username string) (*User, error) {
	query := `
	SELECT username, domain, display_name, avatar_url, settings, flag, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	row := tx.QueryRow(query, username)

	var user User
	var settingsJSON []byte

	if err := row.Scan(
		&user.Username,
		&user.Domain, // Scan domain field
		&user.DisplayName,
		&user.AvatarURL,
		&settingsJSON,
		&user.Flag,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return &user, fmt.Errorf("could not retrieve user: %v", err)
	}

	// 将 settingsJSON 转换为 *orderedmap.OrderedMap
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, fmt.Errorf("could not unmarshal settings: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(tx *sql.Tx, username, displayName, avatarURL string, settings *orderedmap.OrderedMap, flag string) error {
	// 将 settings 转换为 JSON 字符串
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("could not marshal settings: %v", err)
	}

	query := `
		UPDATE users
		SET display_name = $1, avatar_url = $2, settings = $3, flag = $4, updated_at = CURRENT_TIMESTAMP
		WHERE username = $5
	`
	result, err := tx.Exec(query, displayName, avatarURL, settingsJSON, flag, username)
	if err != nil {
		return fmt.Errorf("could not update user: %v", err)
	}

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no user found with the provided username")
	}

	return nil
}

// SoftDeleteUser 标记用户为 "deleted" 而不是实际删除
func SoftDeleteUser(tx *sql.Tx, username string) error {
	query := `
	UPDATE users
	SET flag = 'deleted', updated_at = CURRENT_TIMESTAMP
	WHERE username = $1
	`
	result, err := tx.Exec(query, username)
	if err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	// 检查是否有更新的行
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no user found with the provided username")
	}

	return nil
}
