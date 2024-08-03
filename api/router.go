package api

import (
	ToolsHandler "github.com/Hana-ame/neo-moonchan/Tools/gin_handler"
	"github.com/gin-gonic/gin"
)

func Main() error {
	r := gin.Default()
	// middlewares
	r.Use(TokenMiddleware()) // not tested.
	r.Use(headersMiddleware())
	r.Use(SessionMiddleware())
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
		apiv1status := apiv1.Group("/status")
		{
			apiv1status.POST("", createStatus)       // 创建状态
			apiv1status.GET("/:id", getStatus)       // 获取单个状态
			apiv1status.PUT("/:id", updateStatus)    // 更新状态
			apiv1status.DELETE("/:id", deleteStatus) // 删除状态
		}
		apiv1.GET("/statuses", getStatuses)
		apiv1.GET("/:username/statuses", getUserStatuses)
	}

	err := r.Run("127.24.7.29:8080") // Default listens on :8080

	return err
}
