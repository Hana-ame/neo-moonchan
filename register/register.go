package register

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/debug"
	"github.com/gin-gonic/gin"
)

// 用户名邮箱密码
// POST /api/mail/register/by_mail
func Register(c *gin.Context) {
	// b, _ := io.ReadAll((c.Request.Body))
	// debug.I("b", string(b))
	o, err := tools.ReaderToJSON(c.Request.Body)
	defer c.Request.Body.Close()
	if tools.AbortWithError(c, http.StatusBadRequest, err) {
		return
	}

	email := o.GetOrDefault("from", "").(string)
	for _, suffix := range []string{
		"protonmail.com",
		"proton.me",
	} {
		if strings.HasSuffix(email, suffix) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	password := tools.Hash(o.GetOrDefault("subject", "").(string), os.Getenv("SALT"))

	// id := tools.NewTimeStamp()

	if err := db.Exec(func(tx *sql.Tx) error {
		if _, err := tx.Exec(`
		INSERT INTO accounts (email, password, ehentai) 
		VALUES (LOWER($1), $2, $3::jsonb)
		ON CONFLICT (email) 
		DO UPDATE SET 
		password = EXCLUDED.password, 
		ehentai = jsonb_set(
			COALESCE(accounts.ehentai, '{}'::jsonb),
			'{limit}',
			to_jsonb(
				LEAST(
					COALESCE((accounts.ehentai->>'limit')::integer, 0) + 20000,
					60000
				)
			),
			true
		)
		;`,
			email, password, []byte(`{"ip":"`+c.GetHeader("X-Forwarded-For")+`","limit":20000}`)); err != nil {
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

	// for k, i := range c.Request.Header {
	// 	debug.I("Register", k, i)
	// }

	// c.SetCookie("pass", email+"|"+password, 3600*24*365*10, "/", "", false, true)
	c.Status(http.StatusOK)
	debug.I("Register", "email: ", email, " IP: ", c.GetHeader("X-Forwarded-For"))

}
