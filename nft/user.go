package nft

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func setUserProfile(username string, profile map[string]any) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	// 确保在函数返回前处理事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()
	// 将profile map序列化为JSON
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to marshal profile: %v", err)
	}

	query := `
	INSERT INTO accounts (username, profile) 
	VALUES ($1, $2::jsonb)
	ON CONFLICT (username) 
	DO UPDATE SET profile = EXCLUDED.profile;
`

	// 执行查询
	_, err = tx.Exec(query, username, profileJSON)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update profile: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

// json map
func SetUserProfile(c *gin.Context) {
	username := c.Param("username")
	profile := make(map[string]any)
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := setUserProfile(username, profile)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated successfully"})
}

func getUserProfile(username string) (map[string]any, error) {
	query := `
	SELECT profile, deposit
	FROM accounts 
	WHERE username = $1;
`
	var profileJSON string
	var deposit string
	err := db.QueryRow(query, username).Scan(&profileJSON, &deposit)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %v", err)
	}

	var profile map[string]any
	err = json.Unmarshal([]byte(profileJSON), &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %v", err)
	}

	profile["deposit"] = deposit
	return profile, nil
}

func GetUserProfile(c *gin.Context) {
	username := c.Param("username")
	profile, err := getUserProfile(username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, profile)
}

// username string
func Register(c *gin.Context) {
	username, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}
	if err := c.Request.Body.Close(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to close request body"})
		return
	}
	setUserProfile(string(username), map[string]any{"test": "test"})
	c.SetCookie("username", string(username), 3600, "/", "", false, false)
	c.JSON(200, gin.H{"username": string(username)})
}

// username string
func Login(c *gin.Context) {
	username := c.GetHeader("Dapp-Username")
	if username == "" {
		username, err := c.Cookie("username")
		if err != nil || username == "" {
			c.JSON(401, gin.H{"error": err.Error(), "username": username})
			return
		}
	}

	if err := c.Request.Body.Close(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to close request body"})
		return
	}
	c.SetCookie("username", string(username), 3600, "/", "", false, false)
	if _, err := getUserProfile(string(username)); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	if err := c.Request.Body.Close(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"username": username})
}
