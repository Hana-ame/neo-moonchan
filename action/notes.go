package action

import (
	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

func Create(object *orderedmap.OrderedMap) error {
	id, err := tools.Extract[string](object, "id")
	if err != nil {
		return err
	}
	actor, err := tools.Extract[string](object, "actors")
	if err != nil {
		return err
	}
	published, err := tools.Extract[string](object, "published")
	if err != nil {
		return err
	}
	to, err := tools.Extract[[]any](object, "to")
	if err != nil {
		return err
	}
	cc, err := tools.Extract[[]any](object, "cc")
	if err != nil {
		return err
	}
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
		orderedmap.NewPair("id", id+"/activity"),
		orderedmap.NewPair("type", "Create"),
		orderedmap.NewPair("actor", actor),
		orderedmap.NewPair("published", published),
		orderedmap.NewPair("to", to),
		orderedmap.NewPair("cc", cc),
		orderedmap.NewPair("object", object),
	})

	// body, err := json.Marshal(o)
	// if err != nil {
	// 	return err
	// }

	// for _, endpoint := range tools.NewSlice(to..., cc...){

	// 	FetchWithSign(actor, http.MethodPost, endpoint, body)
	// }

	_ = o

	return nil
}
