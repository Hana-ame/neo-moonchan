// gpt4o @ 240801
// gpt4o @ 240803

package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq" // 导入 PostgreSQL 驱动

	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap" // 有序映射的自定义工具包
)

// User 表示 users 表的结构体
// 定义了用户表的结构，包括用户名、域名、显示名、头像URL、简介、字段、设置、标志、创建时间和更新时间等
type User struct {
	Username    string                   `gorm:"primaryKey;type:varchar(50);not null" json:"username"`
	Domain      string                   `gorm:"type:varchar(100)" json:"domain"`
	DisplayName string                   `gorm:"type:varchar(50)" json:"display_name"`
	AvatarURL   string                   `gorm:"type:varchar(255)" json:"avatar_url"`
	Bios        string                   `gorm:"type:text" json:"bios"`                           // 用户简介
	Fields      []*orderedmap.OrderedMap `gorm:"type:json;not null;default:'{}'" json:"fields"`   // 额外字段（JSON 格式）
	Settings    *orderedmap.OrderedMap   `gorm:"type:json;not null;default:'{}'" json:"settings"` // 用户设置（有序 JSON 格式）
	Flag        string                   `gorm:"type:varchar(50);not null;default:'active'" json:"flag"`
	CreatedAt   time.Time                `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time                `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateUser 插入一个新的用户到 users 表中
// 接受一个事务 tx 和多个参数，用于向 users 表插入一个新用户
func CreateUser(tx *sql.Tx, username, domain, displayName, avatarURL, bios, fieldsJSON string, settings *orderedmap.OrderedMap) error {

	// 将 settings 转换为 JSON 字符串
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("could not marshal settings: %v", err)
	}
	if settings == nil {
		settingsJSON = []byte("{}")
	}

	// 插入用户信息的 SQL 查询语句
	query := `
	INSERT INTO users (username, domain, display_name, avatar_url, bios, fields, settings, flag)
	VALUES ($1, $2, $3, $4, $5, $6, $7, 'active')
	`
	// 执行插入操作
	if _, err := tx.Exec(query, username, domain, displayName, avatarURL, bios, fieldsJSON, settingsJSON); err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}

	return nil
}

// GetUser 根据用户名获取用户信息
// 接受一个事务 tx 和用户名，返回对应的 User 结构体
func GetUser(tx *sql.Tx, username string) (*User, error) {
	// 获取用户信息的 SQL 查询语句
	query := `
	SELECT username, domain, display_name, avatar_url, settings, flag, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	row := tx.QueryRow(query, username)

	var user User
	var fieldsJSON []byte
	var settingsJSON []byte

	// 扫描数据库返回的结果并赋值给 User 结构体
	if err := row.Scan(
		&user.Username,
		&user.Domain, // 扫描 domain 字段
		&user.DisplayName,
		&user.AvatarURL,
		&user.Bios,
		&fieldsJSON,
		&settingsJSON,
		&user.Flag,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return &user, fmt.Errorf("could not retrieve user: %v", err)
	}

	// 将 fieldsJSON 转换为 []*orderedmap.OrderedMap
	if err := json.Unmarshal(fieldsJSON, &user.Fields); err != nil {
		return nil, fmt.Errorf("could not unmarshal fields: %v", err)
	}
	// 将 settingsJSON 转换为 *orderedmap.OrderedMap
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, fmt.Errorf("could not unmarshal settings: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
// 接受一个事务 tx 和多个参数，用于更新指定用户的信息
func UpdateUser(tx *sql.Tx, username, displayName, avatarURL, bios, fieldsJSON string, settings *orderedmap.OrderedMap, flag string) error {

	// 将 settings 转换为 JSON 字符串
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("could not marshal settings: %v", err)
	}
	if settings == nil {
		settingsJSON = []byte("{}")
	}

	// 更新用户信息的 SQL 查询语句
	query := `
	UPDATE users
	SET display_name = $1, avatar_url = $2, bios = $3, fields = $4, settings = $5, flag = $6, updated_at = CURRENT_TIMESTAMP
	WHERE username = $7
	`
	// 执行更新操作
	result, err := tx.Exec(query, displayName, avatarURL, bios, fieldsJSON, settingsJSON, flag, username)
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
// 接受一个事务 tx 和用户名，将用户标记为 'deleted'
func SoftDeleteUser(tx *sql.Tx, username string) error {
	// 标记用户删除的 SQL 查询语句
	query := `
	UPDATE users
	SET flag = 'deleted', updated_at = CURRENT_TIMESTAMP
	WHERE username = $1
	`
	// 执行删除操作
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
