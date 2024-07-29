package api

import "github.com/gin-gonic/gin"

func Main() {
	r := gin.Default()

	r.GET("/echo", func(c *gin.Context) {
		c.String(200, "Hello, world!")
	})
	apiv1auth := r.Group("/api/v1/auth")
	apiv1auth.Use(registerMiddleware())
	{
		apiv1auth.POST("/register", register)
		apiv1auth.POST("/login", login)
		apiv1auth.POST("/logout", logout)
		apiv1auth.DELETE("/sesson/:id", deleteSession)
	}

	r.Run("127.24.7.29:8080") // Default listens on :8080
}
