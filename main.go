package main

import (
	"net/http"
	"os"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	file "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	message "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	fs := &file.FileServer{Path: "/mnt/d/ytdl"}

	// connect to database
	connStr := os.Getenv("DATABASE_URL")

	route := gin.Default()

	// message
	route.PUT("/api/message/:receiver", message.SendMsg)
	route.GET("/api/message/:receiver", message.ReceiveMsg)
	route.PUT("/api/file/upload", fs.Upload)
	route.GET("/api/file/*path", fs.Get)

	route.Any("/api/echo", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, connStr)
	})

	{
		group := route.Group("/api/chan")
		group.POST("/accounts/register", accounts.Register)
		group.POST("/accounts/login", accounts.Login)
		group.POST("/accounts/update", accounts.Update)
	}

	route.Run("127.24.7.29:8080")
}
