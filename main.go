package main

import (
	"net/http"
	"os"
	"path/filepath"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	r2 "github.com/Hana-ame/neo-moonchan/Tools/cloudflare/R2"
	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	handler "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/Tools/openai"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/Hana-ame/neo-moonchan/api/inbox"
	"github.com/Hana-ame/neo-moonchan/api/users"
	"github.com/Hana-ame/neo-moonchan/api/webfinger"
	"github.com/Hana-ame/neo-moonchan/twitter"
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

	// for bili
	route.Any("/v1/chat/completions", openai.GinHandler)

	api := route.Group("/api")

	// message
	api.PUT("/message/:receiver", handler.SendMsg)
	api.GET("/message/:receiver", handler.ReceiveMsg)

	// files 250210
	if tools.HasEnv("UPLOAD_PATH") {
		api.GET("/files/upload", handler.File("upload.html"))
		api.PUT("/files/upload", handler.UploadFilePsql)
		api.GET("/files/:id/:fn", handler.DownloadFilePsql)
		api.GET("/files/list", handler.ListFilesPsql)
	} // files 250210

	if tools.HasEnv("R2_NAME") && tools.HasEnv("R2_ACCOUNT_ID") && tools.HasEnv("R2_ACCESS_KEY_ID") && tools.HasEnv("R2_ACCESS_KEY_SECRET") {
		b, err := r2.NewBucket(os.Getenv("R2_NAME"), os.Getenv("R2_ACCOUNT_ID"), os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_ACCESS_KEY_SECRET"))
		if err == nil {
			api.GET("/r2/upload", handler.File("upload_r2.html"))
			api.PUT("/r2/upload", b.UploadHandler())
			api.GET("/r2/:id/:fn", b.DownloadHandler("id"))
		}
	}

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.GetHeader("X-Forwarded-For"))
	})
	api.Any("/echo", handler.Echo)

	// 删除了ehentai的下载方式,没人用

	chanRouter(api.Group("/chan"))
	mastodonRouter(route)

	twitterRouter(api.Group("/twitter"))

	route.GET("/favicon.png", handler.RedirectTo(http.StatusFound, "/favicon.ico"))

	// 1. 配置静态文件服务（对应$uri检查）
	// route.Static("/static", staticRoot+"/static") // 前端构建产物目录

	staticRoot := os.Getenv("STATIC_ROOT")
	// 2. 处理未匹配路由（对应/index.html回退）
	route.NoRoute(func(c *gin.Context) {
		// debug
		if staticRoot == "" {
			c.File("o.html")
			return
		}

		// 获取请求路径（如 "/asset-manifest.json"）
		requestedPath := c.Request.URL.Path
		// 拼接完整文件路径
		filePath := filepath.Join(staticRoot, requestedPath)

		// 1. 检查文件是否存在且是普通文件（非目录）
		fileInfo, err := os.Stat(filePath)
		if err == nil && !fileInfo.IsDir() {
			// 文件存在且非目录，直接返回文件内容
			c.File(filePath)
			return
		}

		// 2. 若文件不存在或路径是目录，返回index.html
		c.File(filepath.Join(staticRoot, "index.html"))
	})

	route.Run(os.Getenv("LISTEN_ADDR")) // listen and serve on
	// ~/script/vps/ssh.sh -R 127.24.7.29:8080:127.24.7.29:8080
}

func chanRouter(group *gin.RouterGroup) {
	group.POST("/accounts/register", accounts.Register)
	group.POST("/accounts/login", accounts.Login)
	group.POST("/accounts/google/login", accounts.Login)
	group.POST("/accounts/update", accounts.Update)
}

func mastodonRouter(route *gin.Engine) {
	// 以下，activitypub
	// 不使用 /api 开头
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

}

func twitterRouter(g *gin.RouterGroup) {
	g.POST("/", twitter.CreateMetaData)
	g.GET("/:fn", twitter.GetMetaData)
	g.GET("/", twitter.GetLists)
}
