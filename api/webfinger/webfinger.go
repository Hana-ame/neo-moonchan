// 2025年2月24日

package webfinger

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	db "github.com/Hana-ame/neo-moonchan/Tools/db/pq"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

// /.well-known/webfinger?resource=acct:nanakananoka@mstdn.jp
// 404
// {
//     "subject": "acct:nanakananoka@mstdn.jp",
//     "aliases": [
//         "https://mstdn.jp/@nanakananoka",
//         "https://mstdn.jp/users/nanakananoka"
//     ],
//     "links": [
//         {
//             "rel": "http://webfinger.net/rel/profile-page",
//             "type": "text/html",
//             "href": "https://mstdn.jp/@nanakananoka"
//         },
//         {
//             "rel": "self",
//             "type": "application/activity+json",
//             "href": "https://mstdn.jp/users/nanakananoka"
//         },
//         {
//             "rel": "http://ostatus.org/schema/1.0/subscribe",
//             "template": "https://mstdn.jp/authorize_interaction?uri={uri}"
//         }
//     ]
// }

func Webfinger(c *gin.Context) {
	subject := c.Query("resource")
	acct := tools.NewSlice(strings.Split(subject, ":")...).Last().Result()
	username := tools.NewSlice(strings.Split(acct, "@")...).FirstUnequal("")
	host := tools.NewSlice(strings.Split(acct, "@")...).Last().Result()
	id := "https://" + host + "/users/" + username
	// 这里是写死的，因为不会通过其他方式查询。
	err := db.Exec(func(tx *sql.Tx) error {
		o, err := psql.ReadUser(tx, id)
		if err != nil {
			return err
		}

		if o.GetOrDefault("deleted", false).(bool) {
			return fmt.Errorf("deleted")
		}

		return nil
	})
	if err != nil {
		c.Header("X-Error", err.Error())
		c.JSON(http.StatusNotFound, err)
		return
	}

	o := tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("subject", subject),
		orderedmap.NewPair("aliases", tools.Slice[string]{
			"https://" + host + "/@" + username,
			"https://" + host + "/users/" + username,
		}),
		orderedmap.NewPair("links", tools.Slice[*orderedmap.OrderedMap]{
			tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "http://webfinger.net/rel/profile-page"),
				orderedmap.NewPair("type", "text/html"),
				orderedmap.NewPair("href", "https://"+host+"/@"+username),
			}),
			tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "self"),
				orderedmap.NewPair("type", "application/activity+json"),
				orderedmap.NewPair("href", "https://"+host+"/users/"+username),
			}),
			tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "http://ostatus.org/schema/1.0/subscribe"),
				orderedmap.NewPair("template", "https://"+host+"/authorize_interaction?uri={uri}"),
			}),
		}),
	})

	c.JSON(http.StatusOK, o)
}
