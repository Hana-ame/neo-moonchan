package action

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	db "github.com/Hana-ame/neo-moonchan/Tools/db/pq"
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

	// new follow object
	o := tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", id),
		orderedmap.NewPair("type", "Follow"),
		orderedmap.NewPair("actor", actor),
		orderedmap.NewPair("object", object),
		// orderedmap.NewPair("status", "pending"),
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

		inbox, err = tools.ExtractInSequence[string](userObject,
			tools.NewSlice("inbox"),
			tools.NewSlice("endpoints", "sharedInbox"))
		if err != nil {
			return err
		}

		o.Set("status", "pending")
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
	defer resp.Body.Close()

	// b,_:=io.ReadAll(resp.Body)
	tools.WriteReaderToFile("body.txt", resp.Body)

	if slices.Contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
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

func UndoFollow(id string) error {
	return fmt.Errorf("wip")
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

	o := tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
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
		// u, err := url.Parse(object)
		// if err != nil {
		// 	return err
		// }
		inbox, err = tools.ExtractInSequence[string](userObject,
			tools.NewSlice("inbox"),
			tools.NewSlice("endpoints", "sharedInbox"))
		if err != nil {
			return err
		}

		err = psql.CreateActivity(tx, id, o)
		if err != nil {
			return err
		}

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

	// 被接受
	if slices.Contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
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
//		"id": "https://mstdn.jp/users/nanakananoda#blocks/1353228/undo",
//		"type": "Undo",
//		"actor": "https://mstdn.jp/users/nanakananoda",
//		"object": {
//		  "id": "https://mstdn.jp/f3ffe12e-e648-4a42-85fb-02267463b7e7",
//		  "type": "Block",
//		  "actor": "https://mstdn.jp/users/nanakananoda",
//		  "object": "https://mstdn.work.gd/users/nanakananoka"
//		},
//		"@context": "https://www.w3.org/ns/activitystreams"
//	  }
func UndoBlock(actor, object string) error {
	// var err error
	// var inbox string
	err := db.Exec(func(tx *sql.Tx) error {

		// userObject, err := psql.ReadUser(tx, object)
		// if err != nil {
		// 	userObject, err = FetchUser(object)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	psql.SaveUser(tx, object, userObject)
		// }

		// u, err := url.Parse(object)
		// if err != nil {
		// 	return err
		// }
		// inbox, err = tools.ExtractInSequence[string](userObject,
		// 	tools.NewSlice("inbox"),
		// 	tools.NewSlice("endpoints", "sharedInbox"))
		// if err != nil {
		// 	return err
		// }

		objects, err := psql.QueryActivitiesByMap(tx, map[string]string{
			"actor":  actor,
			"object": object,
			"status": "done",
		})
		if err != nil {
			return err
		}
		if len(objects) == 0 {
			return fmt.Errorf("length of objects equals 0")
		}
		object := objects[0] // 偷懒了

		err = Undo(actor, object)
		if err != nil {
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		return err
	}
	return nil
}

// 直接undo一个已有的object
func Undo(actor string, o *orderedmap.OrderedMap) error {
	undo, err := newUndoObject(actor, o)
	if err != nil {
		return err
	}
	id, err := tools.Extract[string](o, "id")
	if err != nil {
		return err
	}
	object, err := tools.Extract[string](o, "object") // user
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

		inbox, err = tools.ExtractInSequence[string](userObject,
			tools.NewSlice("inbox"),
			tools.NewSlice("endpoints", "sharedInbox"))
		if err != nil {
			return err
		}

		body, err := json.Marshal(undo)
		if err != nil {
			return err
		}

		resp, err := FetchWithSign(
			actor,
			http.MethodPost, inbox, body)
		if err != nil {
			return err
		}

		// 被接受
		if slices.Contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
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
	})
	if err != nil {
		return err
	}
	return nil
}

//	{
//		"id": "https://mstdn.jp/users/nanakananoda#blocks/1353228/undo",
//		"type": "Undo",
//		"actor": "https://mstdn.jp/users/nanakananoda",
//		"object": {
//		  "id": "https://mstdn.jp/f3ffe12e-e648-4a42-85fb-02267463b7e7",
//		  "type": "Block",
//		  "actor": "https://mstdn.jp/users/nanakananoda",
//		  "object": "https://mstdn.work.gd/users/nanakananoka"
//		},
//		"@context": "https://www.w3.org/ns/activitystreams"
//	}
//
// 产生 undo 用的 object
func newUndoObject(actor string, object *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	id, err := tools.Extract[string](object, "id")
	if err != nil {
		return object, err
	}
	id = id + "#undo"
	o := tools.OrderedMapFromKVArray(tools.Slice[*orderedmap.Pair]{
		orderedmap.NewPair("@context", "https://www.w3.org/ns/activitystreams"),
		orderedmap.NewPair("id", id),
		orderedmap.NewPair("type", "Undo"),
		orderedmap.NewPair("actor", actor),
		orderedmap.NewPair("object", object),
		orderedmap.NewPair("status", "pending"),
	})

	return o, nil
}
