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

func register(c *gin.Context) {
	// get data
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
	err = psql.CreateUser(tx,
		e.Get("username"),
		os.Getenv("DOMAIN"),
		e.Get("username"),
		os.Getenv("DEFAULT_AVATAR"),
		"",
		"{}",
		nil,
	)

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusCreated)
}

func login(c *gin.Context) {
	// get data
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
	passwordHash := hash(e.Get("password"))
	if passwordHash != account.PasswordHash {
		// add count of failed attemps
		if err := psql.UpdateAccount(tx,
			account.Email, account.PasswordHash, account.Flag, account.FailedAttempts+1); err != nil {
			tx.Rollback()
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("error on commit: %v", err.Error())
				tx.Rollback()
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "password mimatch"})
		return
	}
	// clear the failed counts.
	if account.FailedAttempts != 0 {
		psql.UpdateAccount(tx, account.Email, account.PasswordHash, account.Flag, 0)
	}
	// generate an id.
	sessionIDSlince := make([]byte, 32)
	randomreader.Read(sessionIDSlince)
	sessionID := string(sessionIDSlince)
	// create session
	if err := psql.CreateSession(tx,
		sessionID, account.Username, c.GetString("country"), c.GetString("ip"), c.GetString("ua")); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}

	c.SetCookie("session_id", sessionID, 20*365*24*60*60, "/", "", true, false)
	c.Status(http.StatusCreated)
}

func logout(c *gin.Context) {
	sessionID := c.GetString("session")
	if sessionID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_id not found"})
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
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}

	// delete cookie
	c.SetCookie("session_id", sessionID, -1, "/", "", true, false)
	c.Status(http.StatusNoContent)
}

// in fact you dont need to check the authority cuz who can do this opration must be themself.
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
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}

	c.Status(http.StatusNoContent)
}

// in fact you dont need to check the authority either cuz it will not cause anything if it is a invalid request.
func deleteSessions(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := psql.DeleteSessions(tx, username); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}

	c.Status(http.StatusNoContent)
}

// in fact you dont need to check the authority once again cuz it will not cause anything if it is a invalid request.
func getSessions(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sessions, err := psql.GetSessionList(tx, username)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}

	c.JSON(http.StatusOK, sessions)
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
	// contentType := c.GetHeader("Content-Type")
	// _ = contentType
	if c.ContentType() == "application/json" {
		extractor.cache = make(map[string]string)
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&extractor.cache); err != nil {
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
