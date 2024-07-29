package psql

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	connStr := os.Getenv("DATABASE_URL")
	Connect(connStr)

	// Run the tests
	code := m.Run()

	// Teardown code after running tests
	DB.Close()

	// Exit with the result of m.Run()
	os.Exit(code)
}

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
	email := "a12"
	username := "uA1"
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
	a, _ := GetAccount(tx, "a1")
	_, err = GetAccount(tx, "a22")
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", a)

	tx.Commit()
}

func TestUpdateAccount(t *testing.T) {

}
