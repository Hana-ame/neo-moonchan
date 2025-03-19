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
			"https://w3id.org/security/v1",
		}),
		orderedmap.NewPair("id", "https://mastodon.social/users/yaosese#delete"),
		orderedmap.NewPair("type", "Delete"),
		orderedmap.NewPair("actor", "https://mastodon.social/users/yaosese"),
		orderedmap.NewPair("to", tools.Slice[any]{
			"https://www.w3.org/ns/activitystreams#Public",
		}),
		orderedmap.NewPair("object", "https://mastodon.social/users/yaosese"),
		orderedmap.NewPair("signature", tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
			orderedmap.NewPair("type", "RsaSignature2017"),
			orderedmap.NewPair("creator", "https://mastodon.social/users/yaosese#main-key"),
			orderedmap.NewPair("created", "2025-02-24T11:38:34Z"),
			orderedmap.NewPair("signatureValue", "XeY520USOMnaLAjm8/crb8Ep4zr15KFs7MP7NmVzGJID+GeZoIJBwzlaA/FMcik7rRI8Z4186rZnPY6wqFBhRP9W50QdbT8MRYRSswL6c9EZsATYYjwBD5ZVqfrzLPez2PNZl9B/Ad1K7XktqUXGd0HdOMiGYTK0NdRgoYui5aEccNssZsLUMBfXpTcJaf3g8GZwJoERqpMt2RHva8+IOg9VevQk7t/PCbhrqISrbBXWdkV0bCUFBeboyQOPcTbUpN4SgiO9K+Bi6AQ2SdPPgTe97OESYTvPp1LMSz+R51CGmkA57PjpaLv4Zs4hwD4WXghO5+2zq5IFzeRgKvJAcA=="),
		})),
	})
	b, _ := json.Marshal(o)
	log.Println(string(b))
}
