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

// return forbidden if give an unavaliable session
// after this, handler function could use c.GetString("username") to check
func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetString("username") != "" {
			return
		}
		sessionID, err := c.Cookie("session_id")
		if err == nil {
			tx, err := psql.Begin()
			if err == nil {
				session, err := psql.GetSession(tx, sessionID)
				if err != nil { // 其实应该设置一下是not found
					c.SetCookie("session_id", "expired", -1, "/", "", true, false)
				} else {
					// all success
					c.SetCookie("token", encodeToken(session.Username, time.Now().Unix()+300), -1, "/", "", true, false)
					c.Set("username", session.Username)
					c.Set("session", sessionID)
				}
				if err := tx.Commit(); err != nil {
					log.Printf("error on commit: %v", err.Error())
					tx.Rollback()
				}
			}
		}

		c.Next()
	}
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

func tokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err == nil {
			username, expireAt, err := decodeToken(token)
			if err == nil {
				if expireAt < int(time.Now().Unix()) {
					c.Set("username", username)
				} else {
					c.SetCookie("token", "expired", -1, "/", "", true, false)
				}
			}
		}
		c.Next()
	}
}
