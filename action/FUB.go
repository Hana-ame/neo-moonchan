package action

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/google/uuid"
)

//	{
//	    "@context": "https://www.w3.org/ns/activitystreams",
//	    "id": "https://mstdn.jp/19d9313f-8903-42e5-b5c1-75f32a0f029c",
//	    "type": "Follow",
//	    "actor": "https://mstdn.jp/users/nanakananoka",
//	    "object": "https://fedi.moonchan.xyz/users/nanakananoka"
//	}
func Follow(actor, object string) error {

	id := "https://" + os.Getenv("HOST") + "/" + uuid.NewString()

	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", id),
		orderedmap.NewPair("type", "Follow"),
		orderedmap.NewPair("actor", actor),
		orderedmap.NewPair("object", object),
		orderedmap.NewPair("status", "pending"),
	})
	body, err := json.Marshal(o)
	if err != nil {
		return err
	}

	var inbox string
	err = db.Exec(func(tx *sql.Tx) error {

		// 首先检查是不是有已经有的。

		userObject, err := psql.ReadUser(tx, object)
		if err != nil {
			userObject, err = FetchUser(object)
			if err != nil {
				return err
			}
			psql.SaveUser(tx, object, userObject)
		}

		//
		// func(link string) string { u, _ := url.Parse(link); return u.Hostname() }(object)).(string)
		u, err := url.Parse(object)
		if err != nil {
			return err
		}
		inbox = userObject.GetOrDefault("inbox", u.Host).(string)

		psql.SaveActivity(tx, id, o)

		return tx.Commit()
	})
	if err != nil {
		return err
	}

	resp, err := FetchWithSign(
		actor,
		http.MethodPost, inbox, body)
	if err != nil {
		return err
	}

	if resp.StatusCode%100 == 2 {
		return db.Exec(func(tx *sql.Tx) error {
			o, err := psql.ReadActivity(tx, id)
			if err != nil {
				return err
			}

			o.Set("status", "done")

			err = psql.SaveActivity(tx, id, o)
			if err != nil {
				return err
			}

			return tx.Commit()
		})
	} else {
		return fmt.Errorf("resp: %s", resp.Status)
	}
	// return nil
}

//	{
//	    "@context": "https://www.w3.org/ns/activitystreams",
//	    "id": "https://mstdn.jp/19d9313f-8903-42e5-b5c1-75f32a0f029c",
//	    "type": "Block",
//	    "actor": "https://mstdn.jp/users/nanakananoka",
//	    "object": "https://fedi.moonchan.xyz/users/nanakananoka"
//	}
func Block(actor, object string) error {

	id := "https://" + os.Getenv("HOST") + "/" + uuid.NewString()

	o := tools.OrderedMap(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", id),
		orderedmap.NewPair("type", "Block"),
		orderedmap.NewPair("actor", actor),
		orderedmap.NewPair("object", object),
		orderedmap.NewPair("status", "pending"),
	})
	body, err := json.Marshal(o)
	if err != nil {
		return err
	}

	var inbox string
	err = db.Exec(func(tx *sql.Tx) error {

		userObject, err := psql.ReadUser(tx, object)
		if err != nil {
			userObject, err = FetchUser(object)
			if err != nil {
				return err
			}
			psql.SaveUser(tx, object, userObject)
		}

		//
		// func(link string) string { u, _ := url.Parse(link); return u.Hostname() }(object)).(string)
		u, err := url.Parse(object)
		if err != nil {
			return err
		}
		inbox = userObject.GetOrDefault("inbox", u.Host).(string)

		psql.SaveActivity(tx, id, o)

		return tx.Commit()
	})
	if err != nil {
		return err
	}

	resp, err := FetchWithSign(
		actor,
		http.MethodPost, inbox, body)
	if err != nil {
		return err
	}

	if resp.StatusCode%100 == 2 {
		return db.Exec(func(tx *sql.Tx) error {
			o, err := psql.ReadActivity(tx, id)
			if err != nil {
				return err
			}

			o.Set("status", "done")

			err = psql.SaveActivity(tx, id, o)
			if err != nil {
				return err
			}

			return tx.Commit()
		})
	} else {
		return fmt.Errorf("resp: %s", resp.Status)
	}
	// return nil
}

func UndoFollow(id string) error
