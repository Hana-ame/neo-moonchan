package webfinger

import (
	"net/http"
	"os"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
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
	acct := tools.NewSlice(strings.Split(subject, ":")...).Last()
	username := tools.NewSlice(strings.Split(acct, "@")...).FirstUnequal("")

	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("subject", subject),
		orderedmap.NewPair("aliases", tools.Slice[string]{
			"https://" + os.Getenv("HOST") + "/@" + username,
			"https://" + os.Getenv("HOST") + "/users/" + username,
		}),
		orderedmap.NewPair("links", tools.Slice[*orderedmap.OrderedMap]{
			tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "http://webfinger.net/rel/profile-page"),
				orderedmap.NewPair("type", "text/html"),
				orderedmap.NewPair("href", "https://"+os.Getenv("HOST")+"/@"+username),
			}),
			tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "self"),
				orderedmap.NewPair("type", "application/activity+json"),
				orderedmap.NewPair("href", "https://"+os.Getenv("HOST")+"/users/"+username),
			}),
			tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("rel", "http://ostatus.org/schema/1.0/subscribe"),
				orderedmap.NewPair("template", "https://"+os.Getenv("HOST")+"/authorize_interaction?uri={uri}"),
			}),
		}),
	})

	c.JSON(http.StatusOK, o)
}
