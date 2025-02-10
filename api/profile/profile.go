package profile

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Host      string    `json:"host,omitempty"` //omitempty 表示为空时省略这个字段
	Bio       string    `json:"bio,omitempty"`
	AvatarUrl string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Get(c *gin.Context) {

}
