package users

import (
	"database/sql"
	"fmt"
	"net/http"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/gin-gonic/gin"
)

//	{
//	    "@context": "https://www.w3.org/ns/activitystreams",
//	    "id": "https://mstdn.jp/users/nanakananoka/outbox",
//	    "type": "OrderedCollection",
//	    "totalItems": 245,
//	    "first": "https://mstdn.jp/users/nanakananoka/outbox?page=true",
//	    "last": "https://mstdn.jp/users/nanakananoka/outbox?min_id=0&page=true"
//	}
func Outbox(c *gin.Context) {
	host := c.Request.Host
	username := c.Param("username")

	page := c.Query("page")
	minID := c.Query("min_id")
	maxID := c.Query("max_id")
	if page == "true" {
		o, err := OutboxPage(minID, maxID)
		if err != nil {
			c.Header("X-Error", err.Error())
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusServiceUnavailable, o)
	}

	var count int
	db.Exec(func(tx *sql.Tx) error {

		// count
		count = 0

		return tx.Commit()
	})

	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", "https://"+host+"/users/"+username+"/outbox"),
		orderedmap.NewPair("type", "OrderedCollection"),
		orderedmap.NewPair("totalItems", count),
		orderedmap.NewPair("first", "https://"+host+"/users/"+username+"/outbox?page=true"),
		orderedmap.NewPair("firslastt", "https://"+host+"/users/"+username+"/outbox?min_id=0&page=true"),
	})

	c.JSON(http.StatusOK, o)
}

func OutboxPage(min, max string) (*orderedmap.OrderedMap, error) {
	err := fmt.Errorf("not supported")
	return nil, err
}
