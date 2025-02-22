package main

import (
	"log"
	"net/http"
	"path/filepath"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	file "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	message "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/Hana-ame/neo-moonchan/api/inbox"
	"github.com/Hana-ame/neo-moonchan/api/webfinger"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// fs := &file.FileServer{Path: "/mnt/d/ytdl"}

	// connect to database
	// connStr := os.Getenv("DATABASE_URL")

	route := gin.Default()

	route.Use(middleware.CORSMiddleware())
	api := route.Group("/api")

	// message
	api.PUT("/message/:receiver", message.SendMsg)
	api.GET("/message/:receiver", message.ReceiveMsg)

	// files 250210
	api.GET("/files/upload", file.File("upload.html"))
	api.PUT("/files/upload", file.UploadFilePsql)
	api.GET("/files/:id/:fn", file.DownloadFilePsql)

	api.Any("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// for static files
	staticRoot := "/home/lumin/chat-room/build"
	filelist := []string{
		"asset-manifest.json",
		"favicon.ico",
		"google5f29119424eae036.html",
		"index.html",
		"logo192.png",
		"manifest.json",
		"robots.txt",
		"sw.js",
	}
	for _, path := range filelist {
		route.GET("/"+path, file.File(filepath.Join(staticRoot, path)))
	}
	// route.GET("/static/*any", func(c *gin.Context) {
	// 	path := filepath.Join(staticRoot, "static", (c.Param("any")))
	// 	file.File(path)(c)
	// })
	// end // for static files

	// 1. 配置静态文件服务（对应$uri检查）
	route.Static("/static", staticRoot+"/static") // 前端构建产物目录

	//
	{
		group := api.Group("/chan")
		group.POST("/accounts/register", accounts.Register)
		group.POST("/accounts/login", accounts.Login)
		group.POST("/accounts/update", accounts.Update)
	}

	// activitypub
	route.GET("/.well-known/webfinger", webfinger.Webfinger)

	route.POST("/inbox", inbox.Inbox)

	// 2. 处理未匹配路由（对应/index.html回退）
	route.NoRoute(func(c *gin.Context) {
		log.Println(c.Request.URL.String()) // debug
		c.File(staticRoot + "/index.html")
	})

	route.Run("127.24.7.29:8080")
}
