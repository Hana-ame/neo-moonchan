package nft

import (
	"database/sql"
	"errors"
	"fmt"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/gin-gonic/gin"
)

// addTagToPost associates a tag with a specific post.
// If the association already exists, it does nothing and returns no error.
func addTagToPost(postID string, tagName string) error {
	if tagName == "" {
		return fmt.Errorf("tag name cannot be empty") // Or return nil if empty tag is acceptable no-op
	}
	if tools.Atoi(postID, 0) <= 0 { // Basic validation for postID
		return fmt.Errorf("invalid post ID: %s", postID)
	}

	// The table is post_tags, and the columns are post_id and tag_name.
	// The primary key (or a unique constraint) should be on (post_id, tag_name)
	// for ON CONFLICT to work as intended.
	query := `
		INSERT INTO post_tags (post_id, tag_name)
		VALUES ($1, $2)
		ON CONFLICT (post_id, tag_name) DO NOTHING
	`
	// ON CONFLICT (post_id, tag_name) assumes a unique constraint or primary key on these two columns.
	// DO NOTHING means if the row (this specific post_id and tag_name combination) already exists,
	// the INSERT is skipped, and no error is raised.

	_, err := db.Exec(query, postID, tagName)
	if err != nil {
		// This would be an unexpected error, like a connection issue or syntax error,
		// not a duplicate key error because of DO NOTHING.
		return fmt.Errorf("failed to add tag '%s' to post %s: %w", tagName, postID, err)
	}

	return nil
}

// removeTagFromPost disassociates a tag from a specific post.
// If the association doesn't exist, it does nothing and returns no error.
func removeTagFromPost(postID string, tagName string) error {
	if tagName == "" {
		return fmt.Errorf("tag name cannot be empty for deletion")
	}
	if tools.Atoi(postID, 0) <= 0 {
		return fmt.Errorf("invalid post ID: %s", postID)
	}

	query := `
		DELETE FROM post_tags
		WHERE post_id = $1 AND tag_name = $2
	`

	_, err := db.Exec(query, postID, tagName)
	if err != nil {
		return fmt.Errorf("failed to remove tag '%s' from post %s: %w", tagName, postID, err)
	}

	// Optionally, you could check if any row was actually deleted if that matters.
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	//  log.Printf("Warning: Could not get rows affected for tag deletion: %v", err)
	// }
	// if rowsAffected == 0 {
	//  log.Printf("Tag '%s' was not associated with post %d, or already removed.", tagName, postID)
	// }

	return nil
}

func createOwnerOfPost(postID string, username string, price int64, onsale bool) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty") // Or return nil if empty tag is acceptable no-op
	}
	if tools.Atoi(postID, 0) <= 0 { // Basic validation for postID
		return fmt.Errorf("invalid post ID: %s", postID)
	}

	query := `
		INSERT INTO ownerships (post_id, username, price, onsale)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (post_id) DO NOTHING
	;`
	_, err := db.Exec(query, postID, username, price, onsale)
	if err != nil {
		return err
	}

	return nil
}
func getOwnerOfPost(postID string) (string, string, bool, error) {
	if tools.Atoi(postID, 0) <= 0 {
		return "", "", false, fmt.Errorf("invalid post ID: %s", postID)
	}

	query := `SELECT username, price, onsale FROM ownerships WHERE post_id = $1`
	var owner string
	var price string
	var onsale bool

	// 使用 QueryRow 获取单条记录
	err := db.QueryRow(query, postID).Scan(&owner, &price, &onsale)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return owner, price, onsale, nil // 无数据时返回空字符串+无错误[7,5](@ref)
		}
		return "", "", false, fmt.Errorf("database query failed: %w", err)
	}

	return owner, price, onsale, nil // 无数据时返回空字符串+无错误[7,5](@ref)
}

func patchOwnerOfPost(postID string, username string, price int64, onsaleTrue, onsaleFalse bool) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty") // Or return nil if empty tag is acceptable no-op
	}
	if tools.Atoi(postID, 0) <= 0 { // Basic validation for postID
		return fmt.Errorf("invalid post ID: %s", postID)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var owner string
	var oldPrice int64
	var oldOnsale bool
	err = tx.QueryRow(`
        SELECT username, price, onsale FROM ownerships 
        WHERE post_id = $1 FOR UPDATE`, postID).Scan(&owner, &oldPrice, &oldOnsale)
	if err != nil {
		return err
	}

	if owner != username {
		return fmt.Errorf("not yours")
	}

	price = tools.Or(price, oldPrice)
	onsale := onsaleTrue || (oldOnsale && (onsaleFalse == onsaleTrue))
	// 更新库存
	_, err = tx.Exec(`
        UPDATE ownerships SET price = $1, onsale = $2 
        WHERE post_id = $3`, price, onsale, postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 有风险的。
func changeOwnerOfPost(postID string, username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty") // Or return nil if empty tag is acceptable no-op
	}
	if tools.Atoi(postID, 0) <= 0 { // Basic validation for postID
		return fmt.Errorf("invalid post ID: %s", postID)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var owner string
	var price int64
	var onsale bool
	err = tx.QueryRow(`
        SELECT username, price, onsale FROM ownerships 
        WHERE post_id = $1 FOR UPDATE`, postID).Scan(&owner, &price, &onsale)
	if err != nil {
		return err
	}

	if !onsale {
		return fmt.Errorf("not onsale")
	}

	var mydeposit int64
	err = tx.QueryRow(`
	SELECT deposit FROM accounts 
	WHERE username = $1 FOR UPDATE`, username).Scan(&mydeposit)
	if err != nil {
		return err
	}

	var ownersdeposit int64
	err = tx.QueryRow(`
	SELECT deposit FROM accounts 
	WHERE username = $1 FOR UPDATE`, owner).Scan(&ownersdeposit)
	if err != nil {
		return err
	}

	mydeposit -= price
	ownersdeposit += price

	_, err = tx.Exec(`
        UPDATE accounts SET deposit = $1
        WHERE username = $2`, mydeposit, username)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
        UPDATE accounts SET deposit = $1
        WHERE username = $2`, ownersdeposit, owner)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	UPDATE ownerships SET username = $1, onsale = false
	WHERE post_id = $2`, username, postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetOwnerOfPost(c *gin.Context) {
	postID := c.Param("id")
	owner, price, onsale, err := getOwnerOfPost(postID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"owner": owner, "id": postID, "price": price, "onsale": onsale})
}

func PatchOwnerOfPost(c *gin.Context) {
	postID := c.Param("id")
	username := c.GetHeader("Dapp-Username")
	if username == "" {
		username, err := c.Cookie("username")
		if err != nil || username == "" {
			c.JSON(401, gin.H{"error": err.Error(), "username": username})
			return
		}
	}

	onsale := c.Query("onsale")
	price := tools.Atoi(c.Query("price"), 0)

	// data := make(map[string]string)
	// if err := c.ShouldBindJSON(&data); err != nil {
	// c.JSON(400, gin.H{"error": "Invalid input"})
	// return
	// }
	err := patchOwnerOfPost(postID, username, int64(price), onsale == "true", onsale == "false")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"owner": username, "id": postID, "price": price, "onsale": onsale})
}
func ChangeOwnerOfPost(c *gin.Context) {
	postID := c.Param("id")
	username := c.GetHeader("Dapp-Username")
	if username == "" {
		username, err := c.Cookie("username")
		if err != nil || username == "" {
			c.JSON(401, gin.H{"error": err.Error(), "username": username})
			return
		}
	}

	// data := make(map[string]string)
	// if err := c.ShouldBindJSON(&data); err != nil {
	// c.JSON(400, gin.H{"error": "Invalid input"})
	// return
	// }
	err := changeOwnerOfPost(postID, username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"owner": username, "id": postID, "onsale": false})
}
