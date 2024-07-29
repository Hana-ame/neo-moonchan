package api

import "github.com/gin-gonic/gin"

func registerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		country := c.GetHeader("Cf-Ipcountry")
		ipAddress := c.GetHeader("Cf-Connecting-Ip")
		session := c.GetHeader("M-Session")
		userAgent := c.GetHeader("User-Agent")

		c.Set("country", country)
		c.Set("ip", ipAddress)
		c.Set("session", session)
		c.Set("ua", userAgent)
	}

}
