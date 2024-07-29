package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"testing"
)

func TestConnect(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")

	// 连接到 PostgreSQL 数据库
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	fmt.Println("Successfully connected to the database!")

	DB = db

	fmt.Printf("%v", DB)
}

func TestConnectFunc(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")

	Connect(connStr)

	fmt.Printf("%v", DB)
	tx, err := Begin()
	fmt.Printf("%v\n%v\n%v", DB, tx, err)
}

func TestBegin(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")

	// 连接到 PostgreSQL 数据库
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	fmt.Println("Successfully connected to the database!")

	DB = db

	tx, err := Begin()
	fmt.Printf("%v\n%v\n%v", DB, tx, err)
}
