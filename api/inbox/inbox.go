// 2025年2月24日

package inbox

import (
	"crypto"
	"database/sql"
	"encoding/json"
	"errors"
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
			log.Println(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		o := orderedmap.New()
		if err := json.Unmarshal(b, &o); err != nil {
			log.Println(err)
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
			err = tools.VerifyGin(c, b, retrieve)
			if err != nil {
				// delete 的时候无视找不到。
				v, ok := tools.Extract[string](o, "type")
				if ok == nil && v == "Delete" && errors.Is(err, tools.ErrKeyNotExists) {
					c.AbortWithStatus(http.StatusOK)
					return
				}
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		// 处理，直接写这里好了。
		fmt.Println(string(b), c.Param("username"), c.Request.Host)
		// if err := handler.Inbox(o, name, host, err); err != nil {
		// 	log.Println(err)
		// 	c.JSON(http.StatusUnauthorized, err.Error())
		// 	return
		// }

		c.AbortWithStatus(http.StatusOK)
		return
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
