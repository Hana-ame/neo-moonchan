package accounts

import (
	"database/sql"
	"net/http"
	"os"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	db "github.com/Hana-ame/neo-moonchan/Tools/db/pq"
	"github.com/Hana-ame/neo-moonchan/Tools/debug"
	"github.com/gin-gonic/gin"
)

// 用户名邮箱密码
// POST /api/chan/accounts/register
func Register(c *gin.Context) {
	o, err := tools.ReaderToJSON(c.Request.Body)
	defer c.Request.Body.Close()
	if tools.AbortWithError(c, http.StatusBadRequest, err) {
		return
	}

	email := o.GetOrDefault("email", "").(string)
	password := tools.Hash(o.GetOrDefault("password", "").(string), os.Getenv("SALT"))

	// id := tools.NewTimeStamp()

	if err := db.Exec(func(tx *sql.Tx) error {

		if _, err := tx.Exec(`INSERT INTO accounts (email, password) VALUES ($1, $2);`, email, password); err != nil {
			return err
		}

		return tx.Commit()

	}); err != nil {
		switch {
		default:
			tools.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.SetCookie("pass", email+"|"+password, 3600*24*365*10, "/", "", false, true)
	c.Status(http.StatusCreated)
	debug.I("Register", "email: ", email, " IP: ", c.GetHeader("X-Forwarded-For"))

}
