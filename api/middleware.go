package api

import (
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func registerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		country := c.GetHeader("Cf-Ipcountry")
		ipAddress := c.GetHeader("Cf-Connecting-Ip")
		// session := c.GetHeader("M-Session")
		userAgent := c.GetHeader("User-Agent")

		c.Set("country", country)
		c.Set("ip", ipAddress)
		// c.Set("session", session)
		c.Set("ua", userAgent)
	}
}

func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			return
		}

		tx, err := psql.Begin()
		if err != nil {
			return
		}
		session, err := psql.GetSession(tx, sessionID)
		if err != nil {
			return
		}
		tx.Commit()

		c.Set("session", sessionID)
		c.Set("username", session.Username)
	}
}
