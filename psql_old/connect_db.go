package psql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(connStr string) *sql.DB {
	if DB != nil {
		DB.Close()
	}

	// 连接到 PostgreSQL 数据库
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	// defer db.Close() // not today

	// Set maximum open connections
	db.SetMaxOpenConns(5)
	// Set maximum idle connections
	db.SetMaxIdleConns(2)
	// Set maximum connection lifetime
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	fmt.Println("Successfully connected to the database!")

	DB = db

	return db
}

func Begin() (*sql.Tx, error) {
	return DB.Begin()
}
