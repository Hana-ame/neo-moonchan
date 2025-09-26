package main

import (
	"os"
	"path/filepath"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	middleware "github.com/Hana-ame/neo-moonchan/Tools/my_gin_middleware"
	"github.com/Hana-ame/neo-moonchan/twitter"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.Use(middleware.CORSMiddleware())
	route.Use(middleware.ProxyMiddleware())

	twitter.AddToGroup(route.Group("/api/twitter"))

	staticRoot := os.Getenv("STATIC_ROOT")
	route.NoRoute(func(c *gin.Context) {
		// 防止被读到.env
		if strings.HasPrefix(tools.Or(strings.Split(c.Request.URL.Path, "/")...), ".") {
			c.File(filepath.Join(staticRoot, "index.html"))
		}

		filePath := filepath.Join(staticRoot, c.Request.URL.Path)
		fileInfo, err := os.Stat(filePath)
		if err == nil && !fileInfo.IsDir() {
			c.File(filePath)
			return
		}
		c.File(filepath.Join(staticRoot, "index.html"))
	})

	route.Run(os.Getenv("LISTEN_ADDR")) // listen and serve on

}
