package nft

import "database/sql"

func getCurrentConnections(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity").Scan(&count)
	return count, err
}
