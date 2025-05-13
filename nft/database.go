package nft

import (
	"database/sql"

	tools_db "github.com/Hana-ame/neo-moonchan/Tools/db"
	log "github.com/Hana-ame/neo-moonchan/Tools/debug"
	"github.com/gin-gonic/gin"
)

var db = func() *sql.DB {
	db, err := tools_db.ConnectPostgreSQL("localhost", 5432, "lumin", "lumin", "dapp")
	if err != nil {
		log.F("InitialDataBase", err.Error())
	}
	return db
}()

// DebugInfo 返回当前数据库连接数
func DebugInfo(c *gin.Context) {
	connections, err := getCurrentConnections(db)
	c.JSON(200, gin.H{
		"connections": connections,
		"error":       err,
	})
}
