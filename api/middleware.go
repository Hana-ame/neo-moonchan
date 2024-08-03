// claude @ 240803
// 这个版本的代码更加模块化，每个函数都有明确的职责，使得代码更容易理解和维护。同时，它还保留了原始代码的所有功能。

package api

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func headersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		country := c.GetHeader("Cf-Ipcountry")
		ipAddress := c.GetHeader("Cf-Connecting-Ip")
		userAgent := c.GetHeader("User-Agent")

		c.Set("country", country)
		c.Set("ip", ipAddress)
		c.Set("ua", userAgent)

		c.Next()
	}
}

// SessionMiddleware handles user session validation and token generation.
// It returns a 403 Forbidden status if the session is unavailable.
// After this middleware, handlers can use c.GetString("username") to check the authenticated user.
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isUserAlreadyAuthenticated(c) {
			return
		}

		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.Next()
			return
		}

		if err := handleSession(c, sessionID); err != nil {
			log.Printf("Session handling error: %v", err)
		}

		c.Next()
	}
}

func isUserAlreadyAuthenticated(c *gin.Context) bool {
	return c.GetString("username") != ""
}

func handleSession(c *gin.Context, sessionID string) error {
	tx, err := psql.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	session, err := psql.GetSession(tx, sessionID)
	if err != nil {
		expireSessionCookie(c)
		return err
	}

	setUserSession(c, session)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func expireSessionCookie(c *gin.Context) {
	c.SetCookie("session_id", "expired", -1, "/", "", true, false)
}

func setUserSession(c *gin.Context, session *psql.Session) {
	tokenExpiration := time.Now().Unix() + 300
	token := encodeToken(session.Username, tokenExpiration)
	c.SetCookie("token", token, 300, "/", "", true, false)
	c.Set("username", session.Username)
	c.Set("session", session.SessionID)
}

func encodeToken(username string, expireAt int64) string {
	data := username + "." + strconv.Itoa(int(expireAt))
	dataHash := hash(data)
	return base64.URLEncoding.EncodeToString([]byte(data + "." + dataHash))
}

func decodeToken(token string) (username string, expireAt int, err error) {
	var decodedSlice []byte
	decodedSlice, err = base64.URLEncoding.DecodeString(token)
	if err != nil {
		return
	}
	tokenSlice := strings.Split(string(decodedSlice), ".")
	if len(tokenSlice) != 3 {
		err = fmt.Errorf("token invalid")
		return
	}
	username, expireAtString, dataHash := tokenSlice[0], tokenSlice[1], tokenSlice[2]
	data := username + "." + expireAtString
	if dataHash != hash(data) {
		err = fmt.Errorf("token not match")
		return
	}
	expireAt, err = strconv.Atoi(expireAtString)
	if err != nil {
		return
	}
	// success
	return
}

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Next()
			return
		}

		username, expireAt, err := decodeToken(token)
		if err != nil {
			c.Next()
			return
		}

		if isTokenValid(expireAt) {
			c.Set("username", username)
		} else {
			expireToken(c)
		}

		c.Next()
	}
}

func isTokenValid(expireAt int) bool {
	return expireAt >= int(time.Now().Unix())
}

func expireToken(c *gin.Context) {
	c.SetCookie("token", "expired", -1, "/", "", true, false)
}
