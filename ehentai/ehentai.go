package ehentai

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	handler "github.com/Hana-ame/neo-moonchan/Tools/my_gin_handler"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
)

const PAGE_MOONCHAN_XYZ = "https://page.moonchan.xyz/#"

func helper(c *gin.Context) (ehentai []byte, err error) {
	email, err := c.Cookie("email")
	if err != nil {
		return
	}
	passkey, err := c.Cookie("passkey")
	if err != nil {
		return
	}
	// 验证用户
	err = db.Exec(func(tx *sql.Tx) error {
		var password string
		err := tx.QueryRow("SELECT password, ehentai FROM accounts WHERE email = $1", email).Scan(&password, &ehentai)
		if err != nil {
			return err
		}
		if tools.Hash(email, password) != passkey {
			return fmt.Errorf("not match")
		}
		return tx.Commit()
	})
	if err != nil {
		return
	}
	return
}

func getCost(c *gin.Context) (int, error) {
	header := tools.NewHeader(nil)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
	header.Set("Cookie", os.Getenv("EX_COOKIE"))
	resp, err := myfetch.Fetch(http.MethodGet, "https://exhentai.org"+c.Request.URL.String(), header.Header, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return 0, err
	}
	v := tools.FindAll(doc, "//strong", tools.InnerText)
	// fmt.Println(v.String()) // 等于的
	// fmt.Println(v)
	if len(v) < 4 {
		return 0, fmt.Errorf("not expected: %s", v.String())
	}

	cost := tools.Atoi(v[0], 0)

	return cost, nil
}

// route.POST("/archiver.php", ehentai.Download)
func Download(c *gin.Context) {

	gid := c.Query("gid")
	token := c.Query("token")
	if gid == "" || token == "" {
		c.String(http.StatusBadRequest, "gid or token is empty")
		return
	}

	ehentai, err := helper(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/bounce_login.php?r="+c.Request.URL.String())
		return
	}
	o := orderedmap.New()
	err = json.Unmarshal(ehentai, &o)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	cost, err := getCost(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if cost > 0 && o.GetOrDefault("gp", float64(0)).(float64)+o.GetOrDefault("limit", float64(0)).(float64) < 0 {
		c.String(http.StatusForbidden, "limit: %.0f < gp: %.0f", o.GetOrDefault("limit", float64(0)).(float64), o.GetOrDefault("gp", float64(0)).(float64))
		return
	}

	o.Set("gp", o.GetOrDefault("gp", float64(0)).(float64)-float64(cost))
	byteArray, err := json.Marshal(o)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	email, _ := c.Cookie("email")
	err = db.Exec(func(tx *sql.Tx) error {
		query := `
        UPDATE accounts 
        SET ehentai = $1::jsonb
        WHERE email = $2
    `

		_, err := tx.Exec(query, byteArray, email)
		if err != nil {
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	header := tools.NewHeader(nil)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
	header.Set("Cookie", os.Getenv("EX_COOKIE"))
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := myfetch.Fetch(http.MethodPost, "https://exhentai.org"+c.Request.URL.String(), header.Header, strings.NewReader("dltype=org&dlcheck=Download+Original+Archive"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	respHeader := tools.NewHeader(resp.Header)
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, respHeader.ToMap())

}

// route.GET("/archiver.php", ehentai.Archiver)
func Archiver(c *gin.Context) {
	_, err := c.Cookie("passkey")
	if err != nil {
		c.Redirect(http.StatusFound, "/bounce_login.php?gid="+c.Query("gid")+"&token="+c.Query("token"))
		return
	}

	if http.MethodGet == c.Request.Method {
		handler.File("archiver.html")(c)
		return
	}

	gid := c.Query("gid")
	token := c.Query("token")
	if gid == "" || token == "" {
		// c.String(http.StatusBadRequest, "gid or token is empty")
		c.Header("X-Error", "gid or token is empty")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ehentai, err := helper(c)
	if err != nil {
		// c.Redirect(http.StatusFound, "/bounce_login.php?gid="+c.Query("gid")+"&token="+c.Query("token"))
		c.Header("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	o := orderedmap.New()
	err = json.Unmarshal(ehentai, &o)
	if err != nil {
		// c.String(http.StatusInternalServerError, "error")
		c.Header("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// fmt.Println(o)

	var f64_0 float64

	gp := o.GetOrDefault("gp", f64_0).(float64)
	credit := o.GetOrDefault("credit", f64_0).(float64)
	hath := o.GetOrDefault("hath", f64_0).(float64)
	limit := o.GetOrDefault("limit", f64_0).(float64)

	c.Header("X-Gp", strconv.Itoa(int(gp)))
	c.Header("X-Credit", strconv.Itoa(int(credit)))
	c.Header("X-Hath", strconv.Itoa(int(hath)))
	c.Header("X-Limit", strconv.Itoa(int(limit)))

	cost, err := getCost(c)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		c.Header("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("X-Cost", strconv.Itoa(cost))
	c.Header("Access-Control-Expose-Headers", "X-Gp, X-Credit, X-Hath, X-Cost, X-Limit")
	c.Header("Access-Control-Allow-Origin", "*") // 或指定域名

	c.AbortWithStatus(http.StatusOK)
}

// route.POST("/bounce_login.php", ehentai.Login)
func Login(c *gin.Context) {

	c.Request.ParseForm()
	// for k, v := range c.Request.PostForm {
	// 	fmt.Printf("Key: %s, Value: %v\n", k, v)
	// }
	email := c.Request.PostFormValue("email")
	password := tools.Hash(c.Request.PostFormValue("password") + os.Getenv("SALT"))
	fmt.Println(email, password)

	// 验证用户
	err := db.Exec(func(tx *sql.Tx) error {
		var pw string
		err := tx.QueryRow("SELECT password FROM accounts WHERE email = $1", email).Scan(&pw)
		if err != nil {
			// c.String(http.StatusUnauthorized, "Invalid email")
			return fmt.Errorf("invalid email")
		}
		if pw != password {
			// c.String(http.StatusUnauthorized, "wrong password")
			return fmt.Errorf("wrong password")
		}
		return tx.Commit()
	})
	if err != nil {
		// log.Println("Error checking password:", err)
		c.String(http.StatusUnauthorized, err.Error())
		// c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gid := c.Query("gid")
	token := c.Query("token")
	if gid == "" || token == "" {
		c.Redirect(http.StatusFound, "https://ex.moonchan.xyz/")
		return
	}

	c.SetCookie("email", email, 3600*24*365*10, "/", "", false, false)
	c.SetCookie("passkey", tools.Hash(email, password), 3600*24*365*10, "/", "", false, false)

	c.Redirect(http.StatusFound, "/archiver.php?gid="+gid+"&token="+token+"")
}
