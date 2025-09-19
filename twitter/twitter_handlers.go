package twitter

import (
	"net/http"
	"os"
	"path"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/gin-gonic/gin"
)

// :fn
func CreateMetaData(c *gin.Context) {
	username, ok := c.GetQuery("username")
	if !ok || username == "" {
		c.JSON(400, gin.H{
			"error": "username is required",
		})
		return
	}
	c.Set("username", username)

	_, err := curlMetaData(username)
	if tools.AbortWithError(c, 500, err) {
		return
	}

	c.AbortWithStatus(200)
}

func GetMetaData(c *gin.Context) {
	fn := c.Param("fn")
	if fn == "" {
		c.JSON(400, gin.H{
			"error": "username is required",
		})
		return
	}

	// 根据fn打开文件返回

	result, err := tools.FileToJSON(path.Join(os.Getenv("TWITTER_DIR"), fn))
	if tools.AbortWithError(c, 404, err) {
		return
	}

	c.JSON(200, result)
}

// :fn
func GetLists(c *gin.Context) {
	list, ok := c.GetQuery("list")
	after, _ := c.GetQuery("after")

	if ok {
		r, err := getList(list, after)
		if tools.AbortWithError(c, 500, err) {
			return
		}
		c.JSON(200, r)
		return
	}

	search, ok := c.GetQuery("search")
	if ok {
		by, _ := c.GetQuery("by")
		r, err := getSearch(by, search)
		if tools.AbortWithError(c, 500, err) {
			return
		}
		c.JSON(200, r)
		return
	}

	c.String(http.StatusNotImplemented, "not implemented")

	return
}
