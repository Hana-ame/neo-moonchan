package psql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

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

func CreateAccount(tx *sql.Tx, email, username, passwordHash, country, ipAddress string) error {
	query := `
	INSERT INTO accounts (email, username, password_hash, country, ip_address, flag)
	VALUES ($1, $2, $3, $4, $5, $6)
`
	if _, err := tx.Exec(query, email, username, passwordHash, country, ipAddress, "Padding"); err != nil {
		return err
	}

	return nil
}

func GetAccount(tx *sql.Tx, email string) (*Account, error) {
	query := `
	SELECT email, username, password_hash, country, ip_address, flag, last_login, failed_attempts, created_at, updated_at
	FROM accounts
	WHERE email = $1
`
	row := tx.QueryRow(query, email)

	var account Account
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
		return &account, fmt.Errorf("in the process of login: %v", err)
	}

	// check password
	// if passwordHash != account.PasswordHash {
	// 	// if mismatch, update FailedAttempts.
	// 	if err := UpdateAccount(tx, email, passwordHash, passwordHash, account.Flag, account.FailedAttempts+1); err != nil {
	// 		return &account, err
	// 	}
	// 	return &account, fmt.Errorf("password mismatch")
	// }

	// if err := UpdateAccount(tx, email, passwordHash, passwordHash, account.Flag, 0); err != nil {
	// 	return &account, err
	// }
	return &account, nil
}

func UpdateAccount(tx *sql.Tx, email, newPasswordHash, flag string, failedAttempts int) error {
	// 更新账户信息
	query := `
		UPDATE accounts
		SET password_hash = $1, flag = $2, failed_attempts = $3, updated_at = CURRENT_TIMESTAMP
		WHERE email = $4
	`
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
