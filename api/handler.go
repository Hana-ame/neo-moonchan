package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	e, err := newExtractor(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = psql.CreateAccount(tx,
		e.Get("email"),
		e.Get("username"),
		hash(e.Get("password")),
		c.GetString("country"),
		c.GetString("ip"),
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.Status(http.StatusCreated)
}

func login(c *gin.Context) {
	e, err := newExtractor(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	account, err := psql.GetAccount(tx,
		e.Get("email"),
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// exit if password not patch
	if hash(e.Get("password")) != account.PasswordHash {
		if err := psql.UpdateAccount(tx,
			account.Email, account.PasswordHash, account.Flag, account.FailedAttempts+1); err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password mimatch"})
		return
	}
	// generate an id.
	sessionIDSlince := make([]byte, 32)
	tools.DefaultRandomReader.Read(sessionIDSlince)
	sessionID := string(sessionIDSlince)
	// create session
	if err := psql.CreateSession(tx,
		sessionID, account.Username, c.GetString("country"), c.GetString("ip"), c.GetString("ua")); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.SetCookie("session_id", sessionID, 365*24*3600, "/", "", true, false)
	c.Status(http.StatusCreated)
}

func logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := psql.DeleteSession(tx, sessionID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.SetCookie("session_id", sessionID, 0, "/", "", true, false)
	c.Status(http.StatusNoContent)
}

func deleteSession(c *gin.Context) {
	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := psql.DeleteSession(tx, c.Param("id")); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.Status(http.StatusAccepted)
}

type extractor struct {
	cache map[string]string
	c     *gin.Context
}

func newExtractor(c *gin.Context) (*extractor, error) {

	extractor := &extractor{
		cache: nil,
		c:     c,
	}

	if c.GetHeader("ContentType") == "application/json" {
		cache := make(map[string]string)
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&cache); err != nil {
			return extractor, fmt.Errorf("error encoding body while application/json %v", err)
		}
	}

	return extractor, nil
}

func (e *extractor) Get(key string) string {
	if e.cache == nil {
		return e.c.PostForm(key)
	} else {
		return e.cache[key]
	}
}
