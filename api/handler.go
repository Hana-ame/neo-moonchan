// gpt4o @ 240803

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hana-ame/neo-moonchan/Tools/randomreader"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

// register 处理用户注册请求
// 从请求中提取数据，创建用户账户和相关信息
func register(c *gin.Context) {
	e, err := newExtractor(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	// 创建用户账户
	err = psql.CreateAccount(tx,
		e.Get("email"),
		e.Get("username"),
		hash(e.Get("password")),
		c.GetString("country"),
		c.GetString("ip"),
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 创建用户信息
	err = psql.CreateUser(tx,
		e.Get("username"),
		os.Getenv("DOMAIN"),
		e.Get("username"),
		os.Getenv("DEFAULT_AVATAR"),
		"",
		"{}",
		nil,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// login 处理用户登录请求
// 验证用户凭据并创建会话
func login(c *gin.Context) {
	e, err := newExtractor(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	account, err := psql.GetAccount(tx, e.Get("email"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 验证密码
	if hash(e.Get("password")) != account.PasswordHash {
		if err := psql.UpdateAccount(tx, account.Email, account.PasswordHash, account.Flag, account.FailedAttempts+1); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		} else {
			if err := tx.Commit(); err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
		}
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("password mismatch"))
		return
	}

	// 清除失败尝试计数
	if account.FailedAttempts != 0 {
		psql.UpdateAccount(tx, account.Email, account.PasswordHash, account.Flag, 0)
	}

	// 生成会话 ID
	sessionID := generateSessionID()

	// 创建会话
	if err := psql.CreateSession(tx,
		sessionID, account.Username, c.GetString("country"), c.GetString("ip"), c.GetString("ua")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 设置 cookie
	c.SetCookie("session_id", sessionID, 20*365*24*60*60, "/", "", true, false)
	c.Status(http.StatusCreated)
}

// logout 处理用户登出请求
// 删除用户会话
func logout(c *gin.Context) {
	sessionID := c.GetString("session")
	if sessionID == "" {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if err := psql.DeleteSession(tx, sessionID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 删除 cookie
	c.SetCookie("session_id", "", -1, "/", "", true, false)
	c.Status(http.StatusNoContent)
}

// deleteSession 删除指定的会话
// 用户可以删除自己指定的会话
func deleteSession(c *gin.Context) {
	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if err := psql.DeleteSession(tx, c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteSessions 删除指定用户的所有会话
// 用户可以删除自己所有的会话
func deleteSessions(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unauthorized"))
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if err := psql.DeleteSessions(tx, username); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// getSessions 获取指定用户的所有会话
// 返回用户的所有会话信息
func getSessions(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unauthorized"))
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	sessions, err := psql.GetSessions(tx, username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// newExtractor 初始化并返回一个 extractor 实例
// 根据请求的 Content-Type 提取数据
func newExtractor(c *gin.Context) (*extractor, error) {
	extractor := &extractor{
		cache: nil,
		c:     c,
	}

	if c.ContentType() == "application/json" {
		extractor.cache = make(map[string]string)
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&extractor.cache); err != nil {
			return extractor, fmt.Errorf("error encoding body while application/json %v", err)
		}
	}

	return extractor, nil
}

// extractor 是一个从请求中提取数据的工具
type extractor struct {
	cache map[string]string
	c     *gin.Context
}

// Get 从 extractor 中获取指定键的值
// 优先从缓存中获取，否则从 POST 表单中获取
func (e *extractor) Get(key string) string {
	if e.cache == nil {
		return e.c.PostForm(key)
	} else {
		return e.cache[key]
	}
}

// generateSessionID 生成一个随机的会话 ID
func generateSessionID() string {
	sessionIDSlice := make([]byte, 32)
	randomreader.Read(sessionIDSlice)
	return string(sessionIDSlice)
}
