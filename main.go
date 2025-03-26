package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	handler "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/Hana-ame/neo-moonchan/api/inbox"
	"github.com/Hana-ame/neo-moonchan/api/users"
	"github.com/Hana-ame/neo-moonchan/api/webfinger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	myfetch.DefaultClientPool = myfetch.NewClientPool([]*http.Client{myfetch.NewProxyClient(os.Getenv("HTTPS_PROXY"))}) // 没用啊。
	myfetch.SetDefaultHeader(http.Header{"User-Agent": []string{"MyFetch/1.1.0"}})
	// fs := &handler.FileServer{Path: "/mnt/d/ytdl"}

	// connect to database
	// connStr := os.Getenv("DATABASE_URL")

	route := gin.Default()

	route.Use(gzip.Gzip(gzip.DefaultCompression))

	route.Use(middleware.CORSMiddleware())
	api := route.Group("/api")

	// message
	api.PUT("/message/:receiver", handler.SendMsg)
	api.GET("/message/:receiver", handler.ReceiveMsg)

	// files 250210
	api.GET("/files/upload", handler.File("upload.html"))
	api.PUT("/files/upload", handler.UploadFilePsql)
	api.GET("/files/:id/:fn", handler.DownloadFilePsql)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.GetHeader("X-Forwarded-For"))
	})
	api.Any("/echo", handler.Echo)

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
		route.GET("/"+path, handler.File(filepath.Join(staticRoot, path)))
	}
	// route.GET("/static/*any", func(c *gin.Context) {
	// 	path := filepath.Join(staticRoot, "static", (c.Param("any")))
	// 	handler.File(path)(c)
	// })
	// end // for static files

	// 1. 配置静态文件服务（对应$uri检查）
	route.Static("/static", staticRoot+"/static") // 前端构建产物目录

	// TODO chan 自身的逻辑
	{
		group := api.Group("/chan")
		group.POST("/accounts/register", accounts.Register)
		group.POST("/accounts/login", accounts.Login)
		group.POST("/accounts/update", accounts.Update)
	}

	// activitypub
	route.GET("/.well-known/webfinger", webfinger.Webfinger)

	route.POST("/inbox", inbox.Inbox(true))
	route.POST("/users/:username/inbox", inbox.Inbox(true))

	route.GET("/users/:username", users.Users)
	// mastodon 的这部分根本不支持，敢情只是抓个初态。
	route.GET("/users/:username/outbox", users.Outbox)
	route.GET("/users/:username/following", users.Following)
	route.GET("/users/:username/followers", users.Followers)
	route.GET("/users/:username/collections/featured", users.Featured)
	route.GET("/users/:username/collections/tags", users.Tags)

	// route.GET("/users/:username/inbox", )	 // mastodon没实现，这里实现了之后用
	// route.POST("/users/:username/outbox", )    // mastodon没实现，这里实现了之后用

	// 2. 处理未匹配路由（对应/index.html回退）
	route.NoRoute(func(c *gin.Context) {
		log.Println(c.Request.URL.Path) // debug
		c.File(staticRoot + "/index.html")
	})

	route.Run("127.24.7.29:8080")
}
