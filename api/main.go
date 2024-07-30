package api

import (
	ToolsHandler "github.com/Hana-ame/neo-moonchan/Tools/gin_handler"
	"github.com/gin-gonic/gin"
)

func Main() error {
	r := gin.Default()
	// middlewares
	r.Use(CORSMiddleware())
	r.Use(headersMiddleware())
	r.Use(sessionMiddleware())
	// echo for test
	r.Any("/api/echo", ToolsHandler.Echo)
	// login
	apiv1 := r.Group("/api/v1")
	{
		// accounts
		apiv1.POST("/register", register)
		apiv1.POST("/login", login)
		apiv1.POST("/logout", logout)
		// sessions
		apiv1.GET("/sessions", getSessions)
		apiv1.DELETE("/sessions", deleteSessions)
		apiv1.DELETE("/session/:id", deleteSession)
	}

	err := r.Run("127.24.7.29:8080") // Default listens on :8080

	return err
}
