package api

import (
	"net/http"

	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func headersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		country := c.GetHeader("Cf-Ipcountry")
		ipAddress := c.GetHeader("Cf-Connecting-Ip")
		// session := c.GetHeader("M-Session")
		userAgent := c.GetHeader("User-Agent")

		c.Set("country", country)
		c.Set("ip", ipAddress)
		// c.Set("session", session)
		c.Set("ua", userAgent)

		c.Next()
	}
}

// forbidden if give an unavaliable session
func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		tx, err := psql.Begin()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		session, err := psql.GetSession(tx, sessionID)
		if err != nil {
			c.SetCookie("session_id", sessionID, 0, "/", "", true, false)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		tx.Commit()

		c.Set("session", sessionID)
		c.Set("username", session.Username)

		c.Next()
	}
}

// CORSMiddleware 添加CORS头，允许跨域请求携带cookie
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000") // 允许特定的前端地址
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true") // 允许cookie传递

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
