// gpt4o @ 240802
// gpt4o @ 240803

// Summary of Changes:
// Error Handling: Errors are consistently handled, with a defer tx.Rollback() to ensure a rollback occurs in case of any error.
// Consistency: Functions now follow a similar structure: data extraction, transaction initiation, core logic, transaction commit, and response.
// Comments: Added comments above each function to describe its purpose clearly.
// Readability: Improved the readability by organizing the code and adding spaces where necessary.
package api

import (
	"net/http"
	"strconv"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	psql "github.com/Hana-ame/neo-moonchan/psql_old"
	"github.com/gin-gonic/gin"
)

// createStatus handles creating a new status for a user.
func createStatus(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Extract data from the request
	e, err := tools.NewExtractor(c)
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

	// Create status and link
	timestamp := tools.Now()
	if err := psql.CreateStatus(tx, timestamp, username, e.Get("warning"), e.Get("content"), e.Get("visibility")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := psql.CreateLink(tx, timestamp, username, timestamp, e.Get("visibility")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": timestamp})
}

// getStatus retrieves a specific status by its ID.
func getStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	status, err := psql.GetStatus(tx, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// If the status has been deleted, return 404
	if status.Visibility == "deleted" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, status)
}

// updateStatus updates an existing status based on its ID.
func updateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Extract data from the request
	e, err := tools.NewExtractor(c)
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

	// Update the status
	if err := psql.UpdateStatus(tx, id, e.Get("warning"), e.Get("content")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteStatus soft-deletes a specific status by its ID.
func deleteStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	// Soft delete the status
	if err := psql.SoftDeleteStatus(tx, id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// getUserStatuses retrieves a list of statuses for a specific user.
func getUserStatuses(c *gin.Context) {
	limitString := c.DefaultQuery("limit", "0")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	username := c.Param("username")
	var statuses []*psql.Status

	switch {
	case c.Query("max_id") != "":
		maxID, err := strconv.ParseInt(c.Query("max_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		statuses, err = psql.GetStatusesByUsernameFromLinksMaxID(tx, maxID, username, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case c.Query("min_id") != "":
		minID, err := strconv.ParseInt(c.Query("min_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		statuses, err = psql.GetStatusesByUsernameFromLinksMinID(tx, minID, username, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	default:
		statuses, err = psql.GetLatestStatusesByUsernameFromLinks(tx, username, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, statuses)
}

// getStatuses retrieves statuses based on an ID range (min_id, max_id) or returns the latest statuses.
func getStatuses(c *gin.Context) {
	limitString := c.DefaultQuery("limit", "0")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := psql.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	var statuses []*psql.Status

	switch {
	case c.Query("max_id") != "":
		maxID, err := strconv.ParseInt(c.Query("max_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		statuses, err = psql.GetStatusesByMaxID(tx, maxID, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case c.Query("min_id") != "":
		minID, err := strconv.ParseInt(c.Query("min_id"), 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		statuses, err = psql.GetStatusesByMinID(tx, minID, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	default:
		statuses, err = psql.GetLatestStatuses(tx, limit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, statuses)
}
