package psql

import (
	"fmt"
	"os"
	"testing"
)

func TestTx(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	Connect(connStr)
	tx, err := DB.Begin()
	fmt.Printf("%v \n %v \n %v", DB, tx, err)
}

func TestCreateAccount(t *testing.T) {
	// connStr := os.Getenv("DATABASE_URL")
	// Connect(connStr)

	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	email := "a1112"
	username := "user1"
	passwordHash := "hash1"
	country := "AA"
	ipaddress := "1.1.1.1"

	if err := CreateAccount(tx, email, username, passwordHash, country, ipaddress); err != nil {
		t.Error(err)
		tx.Rollback()
	}

	tx.Commit()

}
func TestGetAccount(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	// var a *Account
	a, _ := GetAccount(tx, "a12")
	_, err = GetAccount(tx, "a22")
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", a)

	tx.Commit()
}

func TestUpdateAccount(t *testing.T) {

}
