// 2025年2月24日

package inbox

import (
	"crypto"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/Hana-ame/neo-moonchan/action"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func Inbox(verify bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		// username := c.Param("username")

		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// log.Println(err)
			c.Header("X-Error", err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		o := orderedmap.New()
		if err := json.Unmarshal(b, &o); err != nil {
			// log.Println(err)
			c.Header("X-Error", err.Error())
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// for debug
		// log.Println(name, "|", c.Request.Host, "|", c.GetHeader("Host")) //   | fedi.moonchan.xyz |
		log.Println(string(b))
		h, _ := json.Marshal(c.Request.Header)
		log.Println(string(h))
		// j, _ := json.Marshal(o)
		// log.Println(string(j)) // same

		if verify {
			// 处理 httpsig
			err = tools.VerifyGin(c, b, retrieve)
			if err != nil {
				// delete 的时候无视找不到。
				v, ok := tools.Extract[string](o, "type")
				// if ok == nil && v == "Delete" && errors.Is(err, tools.ErrKeyNotExists) { // 会有tombstone的
				if ok == nil && v == "Delete" {
					c.AbortWithStatus(http.StatusOK)
					return
				} else {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}

		// 处理，直接写这里好了。
		// fmt.Println(string(b), c.Param("username"), c.Request.Host)
		if err := handle(o, c.Param("username"), c.Request.Host); err != nil {
			log.Println(err)
			c.Header("X-Error", err.Error())
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func retrieve(id string) (publicKey crypto.PublicKey, err error) {
	err = db.Exec(func(tx *sql.Tx) error {

		user, err := psql.ReadUser(tx, id)
		if err != nil {
			user, err = action.FetchUser(id)
			if err != nil {
				return err
			}
			psql.SaveUser(tx, id, user)
		}
		// tools.SaveStructToJsonFile(user, "user.json")
		// fetch在Tombstone的时候不返回err，所以这边返回not found.
		pkpem, err := tools.Extract[string](user, "publicKey", "publicKeyPem")
		if err != nil {
			return err
		}

		publicKey, err = tools.ParsePublicKey([]byte(pkpem))
		if err != nil {
			return err
		}

		return tx.Commit()
	})
	return
}

// TODO:
// 目前只是存起来。不过说起来是不是需要更多处理之类的。根据处理结果更新内部关系
func handle(o *orderedmap.OrderedMap, username, host string) error {
	typ, err := tools.Extract[string](o, "type")
	if err != nil {
		return err
	}
	switch typ {
	case "Follow":
		err = handleFollow(o, username, host)
	case "Block":
		err = handleBlock(o, username, host)
	case "Delete":
		err = handleDelete(o, username, host)
	case "Undo":
		err = handleUndo(o, username, host)
	default:
		err = fmt.Errorf("not supported")
	}
	return err
}

// 确实需要解析一下json-ld, 有空写吧。
func handleDelete(o *orderedmap.OrderedMap, username, host string) error {
	err := db.Exec(func(tx *sql.Tx) error {
		id, err := tools.Extract[string](o, "id")
		if err != nil {
			return err
		}
		o.Set("status", "received")
		err = psql.CreateActivity(tx, id, o)
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

// 逻辑有问题的，在AP这边应该仅仅做一个记录
// 不行，应该 还是需要做了
// 确实需要解析一下json-ld, 有空写吧。
func handleFollow(o *orderedmap.OrderedMap, username, host string) error {
	err := db.Exec(func(tx *sql.Tx) error {
		id, err := tools.Extract[string](o, "id")
		if err != nil {
			return err
		}
		// actor, err := tools.Extract[string](o, "actor")
		// if err != nil {
		// 	return err
		// }
		// object, err := tools.Extract[string](o, "object")
		// if err != nil {
		// 	return err
		// }
		o.Set("status", "received")
		err = psql.CreateActivity(tx, id, o)
		if err != nil {
			return err
		}

		// 检查是否blocked，之后检查？

		return tx.Commit()
	})

	if err != nil {
		return err
	}

	return nil
}

// 确实需要解析一下json-ld, 有空写吧。
func handleBlock(o *orderedmap.OrderedMap, username, host string) error {
	err := db.Exec(func(tx *sql.Tx) error {
		id, err := tools.Extract[string](o, "id")
		if err != nil {
			return err
		}
		o.Set("status", "received")
		err = psql.CreateActivity(tx, id, o)
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

// 确实需要解析一下json-ld, 有空写吧。
func handleUndo(o *orderedmap.OrderedMap, username, host string) error {
	err := db.Exec(func(tx *sql.Tx) error {
		id, err := tools.Extract[string](o, "id")
		if err != nil {
			return err
		}
		// o.Set("status", "received")
		err = psql.CreateActivity(tx, id, o)
		if err != nil {
			return err
		}

		objectID, err := tools.Extract[string](o, "object", "id")
		if err != nil {
			return err
		}
		object, err := psql.ReadActivity(tx, objectID)
		if err != nil {
			return err
		}
		object.Set("status", "undo")
		err = psql.SaveActivity(tx, objectID, object)
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
