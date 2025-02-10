package accounts

import (
	"database/sql"
	"net/http"
	"os"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/debug"
	"github.com/gin-gonic/gin"
)

// POST /api/chan/accounts/login
func Login(c *gin.Context) {
	o, err := tools.ReaderToJSON(c.Request.Body)
	defer c.Request.Body.Close()
	if tools.AbortWithError(c, http.StatusBadRequest, err) {
		return
	}

	email := o.GetOrDefault("email", "").(string)
	password := tools.Hash(o.GetOrDefault("password", "").(string), os.Getenv("SALT"))

	if err := db.Exec(func(tx *sql.Tx) error {
		var retrievedEmail, retrievedPassword string
		err := tx.QueryRow(`SELECT email, password FROM accounts WHERE email = $1 AND password = $2 LIMIT 1`, email, password).Scan(&retrievedEmail, &retrievedPassword)

		if err != nil {
			return err
		}

		return tx.Commit()

	}); err != nil {
		switch {
		case err == sql.ErrNoRows:
			// 如果没有找到匹配的记录，则返回404 Not Found
			tools.AbortWithError(c, http.StatusNotFound, err)
			debug.I("Login", "Failed! email: ", email, " IP: ", c.GetHeader("X-Forwarded-For"))
			return
		default:
			// 对于其他类型的错误，返回500 Internal Server Error
			tools.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.SetCookie("pass", email+"|"+password, 3600*24*365*10, "/", "", true, true)
	c.Status(http.StatusNoContent)
	debug.I("Login", "Succeed! email: ", email, " IP: ", c.GetHeader("X-Forwarded-For"))

}
