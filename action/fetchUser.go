package action

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	db "github.com/Hana-ame/neo-moonchan/Tools/db/pq"
	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/Hana-ame/neo-moonchan/psql"
)

// [@]user@domain.com
// /.well-known/webfinger?resource=acct:nanakananoka@mstdn.jp
func FetchWebfinger(acct string) (id string, err error) {
	host := tools.NewSlice(strings.Split(acct, "@")...).Last().Result()
	username := tools.NewSlice(strings.Split(acct, "@")...).FirstUnequal("")
	_, webfinger, err := myfetch.FetchJSON(http.MethodGet, "https://"+host+"/.well-known/webfinger?resource=acct:"+username+"@"+host, nil, nil)
	if err != nil {
		return
	}
	links := webfinger.GetOrDefault("links", []any{}).([]any)
	for _, link := range links {
		if o, ok := link.(orderedmap.OrderedMap); ok {
			if o.GetOrDefault("rel", "").(string) == "self" {
				return o.GetOrDefault("href", "").(string), nil
			}
		}
	}
	return acct, fmt.Errorf("not found")
}

func FetchUser(id string) (user *orderedmap.OrderedMap, err error) {
	_, user, err = myfetch.FetchJSON(http.MethodGet, id, nil, nil)
	if err != nil {
		return
	}

	err = db.Exec(func(tx *sql.Tx) error {
		return psql.SaveUser(tx, id, user)
	})
	return
}
