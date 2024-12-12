// gpt4o @ 240803

// 文件内容概述：
// Account 结构体：定义了数据库表 accounts 的字段和结构。
// CreateAccount 函数：用于插入新账户到 accounts 表。
// GetAccount 函数：根据电子邮件从 accounts 表中获取账户信息。
// UpdateAccount 函数：更新账户的密码哈希、状态标志和失败尝试次数。
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

// Account 表示 accounts 表的结构体
// 定义了账户表的结构，包括电子邮件、用户名、密码哈希、国家、IP 地址、标志、最后登录时间、失败尝试次数、创建时间和更新时间等
type Account struct {
	Email          string    `gorm:"primaryKey;type:varchar(255)" json:"email"`
	Username       string    `gorm:"unique;type:varchar(50);not null" json:"username"`
	PasswordHash   string    `gorm:"type:varchar(255);not null" json:"password_hash"`
	Country        string    `gorm:"type:char(2);not null" json:"country"`
	IPAddress      string    `gorm:"type:varchar(45);not null" json:"ip_address"`
	Flag           string    `gorm:"type:varchar(255);not null" json:"flag"`
	LastLogin      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_login"`
	FailedAttempts int       `gorm:"default:0" json:"failed_attempts"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateAccount 插入一个新的账户到 accounts 表中
// 接受一个事务 tx 和多个参数，用于向 accounts 表插入一个新账户
func CreateAccount(tx *sql.Tx, email, username, passwordHash, country, ipAddress string) error {
	// 插入账户信息的 SQL 查询语句
	query := `
	INSERT INTO accounts (email, username, password_hash, country, ip_address)
	VALUES ($1, $2, $3, $4, $5)
	`
	// 执行插入操作
	if _, err := tx.Exec(query, email, username, passwordHash, country, ipAddress); err != nil {
		return err
	}

	return nil
}

// GetAccount 根据电子邮件获取账户信息
// 接受一个事务 tx 和电子邮件地址，返回对应的 Account 结构体
func GetAccount(tx *sql.Tx, email string) (*Account, error) {
	// 获取账户信息的 SQL 查询语句
	query := `
	SELECT email, username, password_hash, country, ip_address, flag, last_login, failed_attempts, created_at, updated_at
	FROM accounts
	WHERE email = $1
	`
	row := tx.QueryRow(query, email)

	var account Account
	// 扫描数据库返回的结果并赋值给 Account 结构体
	if err := row.Scan(
		&account.Email,
		&account.Username,
		&account.PasswordHash,
		&account.Country,
		&account.IPAddress,
		&account.Flag,
		&account.LastLogin,
		&account.FailedAttempts,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return &account, fmt.Errorf("could not retrieve account: %v", err)
	}

	return &account, nil
}

// UpdateAccount 更新账户信息
// 接受一个事务 tx 和多个参数，用于更新指定账户的信息
func UpdateAccount(tx *sql.Tx, email, newPasswordHash, flag string, failedAttempts int) error {
	// 更新账户信息的 SQL 查询语句
	query := `
	UPDATE accounts
	SET password_hash = $1, flag = $2, failed_attempts = $3, updated_at = CURRENT_TIMESTAMP
	WHERE email = $4
	`
	// 执行更新操作
	result, err := tx.Exec(query, newPasswordHash, flag, failedAttempts, email)
	if err != nil {
		return fmt.Errorf("could not update account: %v", err)
	}

	// 检查是否有更新的行
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no account found with the provided email")
	}

	return nil
}
