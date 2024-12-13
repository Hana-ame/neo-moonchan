package main

import (
	"net/http"
	"os"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// connect to database
	connStr := os.Getenv("DATABASE_URL")

	route := gin.Default()

	route.Any("/api/echo", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, connStr)
	})

	route.Run("127.24.7.29:8080")
}
