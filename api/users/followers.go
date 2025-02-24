package users

import (
	"fmt"
	"net/http"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/gin-gonic/gin"
)

func Followers(c *gin.Context) {
	host := c.Request.Host
	username := c.Param("username")
	page := c.Query("page")
	pageNum := tools.Atoi(page, 0)
	if pageNum > 0 {
		o, err := FollowingPage(username, pageNum)
		if err != nil {
			c.Header("X-Error", err.Error())
			c.JSON(http.StatusMethodNotAllowed, err)
			return
		}
		c.JSON(http.StatusOK, o)
		return
	}

	count := 0
	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", "https://"+host+"/users/"+username+"/followers"),
		orderedmap.NewPair("type", "OrderedCollection"),
		orderedmap.NewPair("totalItems", count),
		// orderedmap.NewPair("first", "https://"+host+"/users/"+username+"/following?page=1"), // 隐藏时 forbidden
	})
	// {
	// 	"@context": "https://www.w3.org/ns/activitystreams",
	// 	"id": "https://mstdn.jp/users/nanakananoka/following",
	// 	"type": "OrderedCollection",
	// 	"totalItems": 2,
	// 	"first": "https://mstdn.jp/users/nanakananoka/following?page=1"
	// }

	c.JSON(http.StatusOK, o)
}

func FollowersPage(id string, page int) (*orderedmap.OrderedMap, error) {
	return nil, fmt.Errorf("not supported")
}
