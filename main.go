package main

import (
	"net/http"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	file "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	message "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// fs := &file.FileServer{Path: "/mnt/d/ytdl"}

	// connect to database
	// connStr := os.Getenv("DATABASE_URL")

	route := gin.Default()

	route.Use(middleware.CORSMiddleware())

	// message
	route.PUT("/api/message/:receiver", message.SendMsg)
	route.GET("/api/message/:receiver", message.ReceiveMsg)

	// files 250210
	route.GET("/api/files/upload", file.File("upload.html"))
	route.PUT("/api/files/upload", file.UploadFilePsql)
	route.GET("/api/files/:id/:fn", file.DownloadFilePsql)

	route.Any("/api/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	{
		group := route.Group("/api/chan")
		group.POST("/accounts/register", accounts.Register)
		group.POST("/accounts/login", accounts.Login)
		group.POST("/accounts/update", accounts.Update)
	}

	route.Run("127.24.7.29:8080")
}
