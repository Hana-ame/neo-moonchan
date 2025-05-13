package main

import (
	"net/http"
	"os"
	"path/filepath"

	// "github.com/Hana-ame/neo-moonchan/api"
	// "github.com/Hana-ame/neo-moonchan/psql_old"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/liblib"
	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	handler "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/api/accounts"
	"github.com/Hana-ame/neo-moonchan/api/inbox"
	"github.com/Hana-ame/neo-moonchan/api/users"
	"github.com/Hana-ame/neo-moonchan/api/webfinger"
	"github.com/Hana-ame/neo-moonchan/ehentai"
	"github.com/Hana-ame/neo-moonchan/nft"
	"github.com/Hana-ame/neo-moonchan/register"
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

	{
		dapp := route.Group("/dapp")
		dapp.GET("/debug", nft.DebugInfo)
		dapp.GET("/cov/2edge", nft.Conv2Edge)

		dapp.POST("/generate/webui/img2img/ultra", liblib.Img2Img)
		dapp.POST("/generate/webui/text2img/ultra", liblib.Text2Img)
		dapp.POST("/generate/webui/status", liblib.GetStatus)

		// user object
		dapp.POST("/register", nft.Register)
		dapp.POST("/login", nft.Login)
		dapp.GET("/user/:username", nft.GetUserProfile)
		dapp.PATCH("/user/:username", nft.SetUserProfile)
		// post object
		dapp.POST("/post/create", nft.CreatePost)

		dapp.GET("/post/:id", nft.GetPost)
		dapp.PATCH("/post/:id", nft.CreatePost)

		dapp.GET("post/:id/owner", nft.GetOwnerOfPost)
		dapp.POST("post/:id/owner", nft.ChangeOwnerOfPost)
		dapp.PATCH("post/:id/owner", nft.PatchOwnerOfPost)

		dapp.GET("/user/:username/posts", nft.GetPostsByUsername)
		dapp.GET("/user/:username/owned", nft.GetPostsByOwnershipWithJoin)

		// comment object
		dapp.GET("/post/:id/comment")
		dapp.POST("/post/:id/comment")

		// catagory
		dapp.GET("/explore", nft.GetNewPosts)
		dapp.GET("/explore/:tag", nft.GetPostsByTagWithJoin)
		dapp.GET("/collection/:id", nft.GetPostsByOwnershipWithJoin)
	}

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

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.GetHeader("X-Forwarded-For"))
	})
	api.Any("/echo", handler.Echo)

	// ehentai
	if tools.HasEnv("EX_COOKIE") {
		api.POST("/register/by-mail", register.Register)

		// ehentai
		route.GET("/archiver.php", ehentai.Archiver)
		route.POST("/archiver.php", ehentai.Download)
		route.HEAD("/archiver.php", ehentai.Archiver)
		route.GET("/bounce_login.php", handler.File("bounce_login.html"))
		route.POST("/bounce_login.php", ehentai.Login)
		route.HEAD("/bounce_login.php", handler.File("bounce_login.html"))
	} // ehentai

	// for static files
	// staticRoot := os.Getenv("STATIC_ROOT")
	// filelist := []string{
	// 	"asset-manifest.json",
	// 	"favicon.ico",
	// 	"google5f29119424eae036.html",
	// 	"index.html",
	// 	"logo192.png",
	// 	"manifest.json",
	// 	"robots.txt",
	// 	"sw.js",
	// }
	// for _, path := range filelist {
	// 	route.GET("/"+path, handler.File(filepath.Join(staticRoot, path)))
	// }
	// route.GET("/static/*any", func(c *gin.Context) {
	// 	path := filepath.Join(staticRoot, "static", (c.Param("any")))
	// 	handler.File(path)(c)
	// })
	// end // for static files

	// 1. 配置静态文件服务（对应$uri检查）
	// route.Static("/static", staticRoot+"/static") // 前端构建产物目录

	// TODO chan 自身的逻辑 path = chan
	{
		group := api.Group("/chan")
		group.POST("/accounts/register", accounts.Register)
		group.POST("/accounts/login", accounts.Login)
		group.POST("/accounts/google/login", accounts.Login)
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

	route.GET("/favicon.png", handler.RedirectTo(http.StatusTemporaryRedirect, "/favicon.ico"))

	staticRoot := os.Getenv("STATIC_ROOT")
	// 2. 处理未匹配路由（对应/index.html回退）
	route.NoRoute(func(c *gin.Context) {
		c.File("o.html")
		return
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
