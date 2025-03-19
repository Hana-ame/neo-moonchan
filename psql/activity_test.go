package psql

import (
	"encoding/json"
	"log"
	"testing"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

func TestJSONEncoder(t *testing.T) {
	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", tools.Slice[any]{
			"https://www.w3.org/ns/activitystreams",
			tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("ostatus", "http://ostatus.org#"),
				orderedmap.NewPair("atomUri", "ostatus:atomUri"),
				orderedmap.NewPair("inReplyToAtomUri", "ostatus:inReplyToAtomUri"),
				orderedmap.NewPair("conversation", "ostatus:conversation"),
				orderedmap.NewPair("sensitive", "as:sensitive"),
				orderedmap.NewPair("toot", "http://joinmastodon.org/ns#"),
				orderedmap.NewPair("votersCount", "toot:votersCount"),
			}),
		}),
		orderedmap.NewPair("id", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855/activity"),
		orderedmap.NewPair("type", "Create"),
		orderedmap.NewPair("actor", "https://mstdn.jp/users/nanakananoka"),
		orderedmap.NewPair("published", "2025-02-24T14:31:02Z"),
		orderedmap.NewPair("to", tools.Slice[any]{
			"https://fedi.moonchan.xyz/users/nanakananoka",
		}),
		orderedmap.NewPair("cc", tools.Slice[any]{}),
		orderedmap.NewPair("object", tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
			orderedmap.NewPair("id", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855"),
			orderedmap.NewPair("type", "Note"),
			orderedmap.NewPair("summary", nil),
			orderedmap.NewPair("inReplyTo", nil),
			orderedmap.NewPair("published", "2025-02-24T14:31:02Z"),
			orderedmap.NewPair("url", "https://mstdn.jp/@nanakananoka/114059343442062855"),
			orderedmap.NewPair("attributedTo", "https://mstdn.jp/users/nanakananoka"),
			orderedmap.NewPair("to", tools.Slice[any]{
				"https://fedi.moonchan.xyz/users/nanakananoka",
			}),
			orderedmap.NewPair("cc", tools.Slice[any]{}),
			orderedmap.NewPair("sensitive", false),
			orderedmap.NewPair("atomUri", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855"),
			orderedmap.NewPair("inReplyToAtomUri", nil),
			orderedmap.NewPair("conversation", "tag:mstdn.jp,2025-02-24:objectId=592059555:objectType=Conversation"),
			orderedmap.NewPair("content", "<p><span class=\"h-card\"><a href=\"https://fedi.moonchan.xyz/@nanakananoka\" class=\"u-url mention\">@<span>nanakananoka@fedi.moonchan.xyz</span></a></span></p>"),
			orderedmap.NewPair("contentMap", tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("zh", "<p><span class=\"h-card\"><a href=\"https://fedi.moonchan.xyz/@nanakananoka\" class=\"u-url mention\">@<span>nanakananoka@fedi.moonchan.xyz</span></a></span></p>"),
			})),
			orderedmap.NewPair("attachment", tools.Slice[any]{}),
			orderedmap.NewPair("tag", tools.Slice[any]{
				tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
					orderedmap.NewPair("type", "Mention"),
					orderedmap.NewPair("href", "https://fedi.moonchan.xyz/users/nanakananoka"),
					orderedmap.NewPair("name", "@nanakananoka@fedi.moonchan.xyz"),
				}),
			}),
			orderedmap.NewPair("replies", tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
				orderedmap.NewPair("id", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855/replies"),
				orderedmap.NewPair("type", "Collection"),
				orderedmap.NewPair("first", tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
					orderedmap.NewPair("type", "CollectionPage"),
					orderedmap.NewPair("next", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855/replies?only_other_accounts=true&page=true"),
					orderedmap.NewPair("partOf", "https://mstdn.jp/users/nanakananoka/statuses/114059343442062855/replies"),
					orderedmap.NewPair("items", tools.Slice[any]{}),
				})),
			})),
		})),
	})
	b, _ := json.Marshal(o)
	log.Println(string(b))
}
