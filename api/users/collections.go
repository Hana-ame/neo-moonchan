package users

import (
	"fmt"
	"net/http"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/gin-gonic/gin"
)

// 顶置
// 已经固定初值
//
//	{
//	    "@context": [
//	        "https://www.w3.org/ns/activitystreams",
//	        {
//	            "ostatus": "http://ostatus.org#",
//	            "atomUri": "ostatus:atomUri",
//	            "inReplyToAtomUri": "ostatus:inReplyToAtomUri",
//	            "conversation": "ostatus:conversation",
//	            "sensitive": "as:sensitive",
//	            "toot": "http://joinmastodon.org/ns#",
//	            "votersCount": "toot:votersCount"
//	        }
//	    ],
//	    "id": "https://mstdn.jp/users/nanakananoka/collections/featured",
//	    "type": "OrderedCollection",
//	    "totalItems": 1,
//	    "orderedItems": [
//	        {
//	            "id": "https://mstdn.jp/users/nanakananoka/statuses/113306680741186613",
//	            "type": "Note",
//	            "summary": null,
//	            "inReplyTo": null,
//	            "published": "2024-10-14T16:19:00Z",
//	            "url": "https://mstdn.jp/@nanakananoka/113306680741186613",
//	            "attributedTo": "https://mstdn.jp/users/nanakananoka",
//	            "to": [
//	                "https://www.w3.org/ns/activitystreams#Public"
//	            ],
//	            "cc": [
//	                "https://mstdn.jp/users/nanakananoka/followers"
//	            ],
//	            "sensitive": false,
//	            "atomUri": "https://mstdn.jp/users/nanakananoka/statuses/113306680741186613",
//	            "inReplyToAtomUri": null,
//	            "conversation": "tag:mstdn.jp,2024-10-14:objectId=550975646:objectType=Conversation",
//	            "content": "<p>说起来我xp也就是20代后半被整得死去活来満身創痍深重黑眼圈肾虚音的美少女了吧</p>",
//	            "contentMap": {
//	                "en": "<p>说起来我xp也就是20代后半被整得死去活来満身創痍深重黑眼圈肾虚音的美少女了吧</p>"
//	            },
//	            "attachment": [],
//	            "tag": [],
//	            "replies": {
//	                "id": "https://mstdn.jp/users/nanakananoka/statuses/113306680741186613/replies",
//	                "type": "Collection",
//	                "first": {
//	                    "type": "CollectionPage",
//	                    "next": "https://mstdn.jp/users/nanakananoka/statuses/113306680741186613/replies?only_other_accounts=true&page=true",
//	                    "partOf": "https://mstdn.jp/users/nanakananoka/statuses/113306680741186613/replies",
//	                    "items": []
//	                }
//	            }
//	        }
//	    ]
//	}
func Featured(c *gin.Context) {
	err := fmt.Errorf("not supported")
	c.JSON(http.StatusMethodNotAllowed, err)
}

// 不知道是什么
// 已经固定初值
//
//	{
//	    "@context": "https://www.w3.org/ns/activitystreams",
//	    "id": "https://mstdn.jp/users/nanakananoka/collections/devices",
//	    "type": "Collection",
//	    "totalItems": 0,
//	    "items": []
//	}
func Devices(c *gin.Context) {
	host := c.Request.Host
	username := c.Param("username")

	count := 0
	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", "https://"+host+"/users/"+username+"/collections/devices"),
		orderedmap.NewPair("type", "Collection"),
		orderedmap.NewPair("totalItems", count),
		orderedmap.NewPair("items", []struct{}{}),
	})

	c.JSON(http.StatusOK, o)
}

// 好像是主页的tag，不是很常用吧。
// 已经固定初值
//
//	{
//	    "@context": "https://www.w3.org/ns/activitystreams",
//	    "id": "https://mstdn.jp/users/nanakananoka/collections/tags",
//	    "type": "Collection",
//	    "totalItems": 0,
//	    "items": []
//	}
func Tags(c *gin.Context) {
	host := c.Request.Host
	username := c.Param("username")

	count := 0
	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", "https://"+host+"/users/"+username+"/collections/tags"),
		orderedmap.NewPair("type", "Collection"),
		orderedmap.NewPair("totalItems", count),
		orderedmap.NewPair("items", []struct{}{}),
	})

	c.JSON(http.StatusOK, o)
}
