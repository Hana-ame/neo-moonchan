package api

import (
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

// forbidden if give an unavaliable session
func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err == nil {
			tx, err := psql.Begin()
			if err == nil {
				session, err := psql.GetSession(tx, sessionID)
				if err != nil { // 其实应该设置一下是not found
					c.SetCookie("session_id", sessionID, -1, "/", "", true, false)
				} else {
					// all success
					c.Set("username", session.Username)
					c.Set("session", sessionID)
				}
				tx.Commit()
			}
		}

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
