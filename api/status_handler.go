// gpt4o @ 240802

package api

import (
	"net/http"
	"strconv"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/psql"
	"github.com/gin-gonic/gin"
)

func createStatus(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// get data
	e, err := newExtractor(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	timestamp := tools.Now()
	if err := psql.CreateStatus(tx,
		timestamp,
		username,
		e.Get("warning"),
		e.Get("content"),
		e.Get("visibility"),
	); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := psql.CreateLink(tx,
		timestamp,
		username,
		timestamp,
		e.Get("visibility"),
	); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusCreated)
}

func getStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	status, err := psql.GetStatus(tx, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if status.Visibility == "deleted" {
		c.JSON(http.StatusNotFound, nil)
	}

	c.JSON(http.StatusOK, status)
}

func updateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// get data
	e, err := newExtractor(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	if err := psql.UpdateStatus(tx,
		id,
		e.Get("warning"),
		e.Get("content"),
	); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	if err := psql.SoftDeleteStatus(tx, id); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func getStatuses(c *gin.Context) {
	limitString := c.DefaultQuery("limit", "0")
	limit, _ := strconv.Atoi(limitString)

	tx, err := psql.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	rollbackWithError := func(err error) {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	username := c.Param("username")
	var statuses []*psql.Status

	switch {
	case c.Query("max_id") != "":
		maxID, err := strconv.ParseInt(c.Query("max_id"), 10, 64)
		if err != nil {
			rollbackWithError(err)
			return
		}
		statuses, err = psql.GetStatusesFromLinksMaxID(tx, maxID, username, limit)
		if err != nil {
			rollbackWithError(err)
			return
		}

	case c.Query("min_id") != "":
		minID, err := strconv.ParseInt(c.Query("min_id"), 10, 64)
		if err != nil {
			rollbackWithError(err)
			return
		}
		statuses, err = psql.GetStatusesFromLinksMinID(tx, minID, username, limit)
		if err != nil {
			rollbackWithError(err)
			return
		}

	default:
		statuses, err = psql.GetStatusesFromLinks(tx, username, limit)
		if err != nil {
			rollbackWithError(err)
			return
		}
	}

	c.JSON(http.StatusOK, statuses)
}
