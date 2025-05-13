package nft

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// 定义帖子结构体（需与数据库表结构匹配）
type Post struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	// Owner    string `json:"owner,omitempty"`
	MetaData string `json:"meta_data,omitempty"`
}

// create on id == 0
func createPost(id, username, url string, content string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %v", err)
	}
	// 确保在函数返回前处理事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	id = tools.Or(id, strconv.Itoa(int(tools.NewTimeStamp())))

	query := `
	INSERT INTO posts (id, username, url, content) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) 
	DO UPDATE SET 
	username = EXCLUDED.username,
	url = EXCLUDED.url,
	content = EXCLUDED.content;
`

	// 执行查询
	_, err = tx.Exec(query, id, username, url, content)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to update post: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("could not commit transaction: %v", err)
	}

	return id, nil
}

func createTags(id string, tags []string) error {
	// 开启事务确保原子性
	tx, err := db.Begin()
	if err != nil {
		// log.Printf("事务启动失败: %v", err)
		return err
	}
	defer tx.Rollback() // 确保未提交时自动回滚[7,8](@ref)

	// 预编译插入语句（PostgreSQL占位符格式）
	stmt, err := tx.Prepare(`
        INSERT INTO post_tags(post_id, tag_name) 
        VALUES($1, $2)
        ON CONFLICT (post_id, tag_name) DO NOTHING
    `)
	if err != nil {
		// log.Printf("SQL预编译失败: %v", err)
		return err
	}
	defer stmt.Close()

	// 遍历标签并插入
	for _, tag := range tags {
		_, err := stmt.Exec(id, tag)
		if err != nil {
			// log.Printf("标签插入失败[post_id=%d, tag='%s']: %v", id, tag, err)
			return err // 立即终止并回滚
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		// log.Printf("事务提交失败: %v", err)
		return err
	}
	return nil
}

// 根据ID查询单条帖子记录
func getPost(id int64) (*Post, error) {
	query := `
        SELECT id, username, url, content 
        FROM posts 
        WHERE id = $1
    `

	var post Post
	// 使用 QueryRow 直接获取单条记录（避免事务开销）
	err := db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Username,
		&post.URL,
		&post.Content,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post id %d not found", id)
		}
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	return &post, nil
}

// 用到了吗。没？
func getPosts(ids []int64) ([]*Post, error) {
	// Handle empty ids slice to avoid an unnecessary query or potential DB error with ANY('{}')
	if len(ids) == 0 {
		return []*Post{}, nil // Return empty slice, no error
	}

	query := `
		SELECT id, username, url, content 
		FROM posts 
		WHERE id = ANY($1)
	`
	// db.Query is the correct function for fetching multiple rows.
	// db.QueryRow is for a single row.
	// We must wrap the 'ids' slice with pq.Array() for the pq driver
	// to correctly interpret it as a PostgreSQL array.
	rows, err := db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err) // Use %w for error wrapping
	}
	defer rows.Close()

	// Pre-allocate slice with a capacity, can be a small optimization
	posts := make([]*Post, 0, len(ids))

	for rows.Next() {
		var post Post
		// Ensure your Post struct fields match the scanned columns.
		// If 'url' can be NULL in the database, use sql.NullString or similar.
		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.URL, // If URL is nullable, post.URL should be &sql.NullString
			&post.Content,
		)
		if err != nil {
			// This error occurs if a specific row scan fails (e.g., type mismatch)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		posts = append(posts, &post)
	}

	// Check for errors encountered during iteration (e.g., network issues)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	// It's idiomatic in Go to return (nil, error) or (data, nil).
	// If no posts are found for the given IDs, 'posts' will be an empty slice,
	// and 'err' will be nil, which is usually the desired behavior.
	return posts, nil
}

// 获得时间线上的
func getNewPosts(before, limit int64) ([]*Post, error) {
	if limit <= 0 {
		limit = 1
	}
	if limit >= 64 {
		limit = 64
	}
	var rows *sql.Rows
	var err error

	if before > 0 {
		query := `
		SELECT id, username, url, content 
		FROM posts 
		WHERE id < $1
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, before, limit)
	} else {
		query := `
		SELECT id, username, url, content 
		FROM posts 
		ORDER BY id DESC
		LIMIT $1;
	`
		rows, err = db.Query(query, limit)
	}
	// db.Query is the correct function for fetching multiple rows.
	// db.QueryRow is for a single row.
	// We must wrap the 'ids' slice with pq.Array() for the pq driver
	// to correctly interpret it as a PostgreSQL array.
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err) // Use %w for error wrapping
	}
	defer rows.Close()

	// Pre-allocate slice with a capacity, can be a small optimization
	posts := make([]*Post, 0, limit)

	for rows.Next() {
		var post Post
		// Ensure your Post struct fields match the scanned columns.
		// If 'url' can be NULL in the database, use sql.NullString or similar.
		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.URL, // If URL is nullable, post.URL should be &sql.NullString
			&post.Content,
		)
		if err != nil {
			// This error occurs if a specific row scan fails (e.g., type mismatch)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		posts = append(posts, &post)
	}

	// Check for errors encountered during iteration (e.g., network issues)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	// It's idiomatic in Go to return (nil, error) or (data, nil).
	// If no posts are found for the given IDs, 'posts' will be an empty slice,
	// and 'err' will be nil, which is usually the desired behavior.
	return posts, nil
}
func getPostsByUsername(username string, before, limit int64) ([]*Post, error) {

	if limit <= 0 {
		limit = 1
	}
	if limit >= 64 {
		limit = 64
	}
	var rows *sql.Rows
	var err error

	if before > 0 {
		query := `
		SELECT id, username, url, content 
		FROM posts 
		WHERE id < $1
		AND username = $2
		ORDER BY id DESC
		LIMIT $3;
	`
		rows, err = db.Query(query, before, username, limit)
	} else {
		query := `
		SELECT id, username, url, content 
		FROM posts 
		WHERE username = $1
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, username, limit)
	}

	if err != nil {
		// c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Username, &post.URL, &post.Content)
		if err != nil {
			// c.JSON(500, gin.H{"error": err.Error()})
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		// c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	return posts, nil
}

func getPostsByTagWithJoin(tag string, before, limit int64) ([]*Post, error) {

	if limit <= 0 {
		limit = 1
	}
	if limit >= 64 {
		limit = 64
	}
	var rows *sql.Rows
	var err error

	if before > 0 {
		query := `
		SELECT p.id, p.username, p.url, p.content
		FROM posts p
		JOIN post_tags pt ON p.id = pt.post_id
		WHERE pt.post_id < $1
		AND pt.tag_name = $2
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, before, tag, limit)
	} else {
		query := `
		SELECT p.id, p.username, p.url, p.content
		FROM posts p
		JOIN post_tags pt ON p.id = pt.post_id
		WHERE pt.tag_name = $1
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, tag, limit)
	}

	if err != nil {
		log.Printf("Error querying posts with join for tag '%s': %v", tag, err)
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.URL,
			&post.Content,
		)
		if err != nil {
			log.Printf("Error scanning post row with join for tag '%s': %v", tag, err)
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating post rows with join for tag '%s': %v", tag, err)
		return nil, err
	}

	return posts, nil
}

func getPostsByOwnershipWithJoin(username string, before, limit int64) ([]*Post, error) {
	if limit <= 0 {
		limit = 1
	}
	if limit >= 64 {
		limit = 64
	}
	var rows *sql.Rows
	var err error

	if before > 0 {
		query := `
		SELECT p.id, p.username, p.url, p.content
		FROM posts p
		JOIN ownerships p1 ON p.id = p1.post_id
		WHERE p1.username = $1
		AND p.id < $2
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, username, before, limit)
	} else {
		query := `
		SELECT p.id, p.username, p.url, p.content
		FROM posts p
		JOIN ownerships p1 ON p.id = p1.post_id
		WHERE p1.username = $1
		ORDER BY id DESC
		LIMIT $2;
	`
		rows, err = db.Query(query, username, limit)
	}

	if err != nil {
		log.Printf("Error querying posts with join for tag '%s': %v", username, err)
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.URL,
			&post.Content,
		)
		if err != nil {
			log.Printf("Error scanning post row with join for tag '%s': %v", username, err)
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating post rows with join for tag '%s': %v", username, err)
		return nil, err
	}

	return posts, nil
}

// json Post
func CreatePost(c *gin.Context) {
	username := c.GetHeader("Dapp-Username")
	if username == "" {
		username, err := c.Cookie("username")
		if err != nil || username == "" {
			c.JSON(401, gin.H{"error": err.Error(), "username": username})
			return
		}
	}
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	post.Username = username
	// post.ID = int64(tools.Atoi(c.Param("id"), 0)) // update的时候是会带这	个参数的。
	post.ID = c.Param("id") // update的时候是会带这	个参数的。
	metaData := orderedmap.New()
	json.Unmarshal([]byte(post.MetaData), &metaData)

	id, err := createPost(post.ID, post.Username, post.URL, post.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	post.ID = id
	err = createOwnerOfPost(post.ID, post.Username, 0, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = createTags(post.ID, tools.Slice[any](metaData.GetOrDefault("tags", []any{}).([]any)).ToString())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, post)
}

func GetPost(c *gin.Context) {

	id := c.Param("id")
	postID := int64(tools.Atoi(id, 0))
	if postID == 0 {
		c.JSON(400, gin.H{"error": "Invalid post ID", "id": id})
		return
	}
	post, err := getPost(postID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// owner, _ := getOwnerOfPost(id)
	// post.Owner = owner

	c.JSON(200, post)
}

func GetNewPosts(c *gin.Context) {
	limit := tools.Atoi(c.Query("limit"), 20)
	before := tools.Atoi(c.Query("before"), 0)

	posts, err := getNewPosts(int64(before), int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Even if no posts are found, 'posts' will be an empty slice, which is fine.
	c.JSON(http.StatusOK, posts)
}

func GetPostsByUsername(c *gin.Context) {
	username := c.Param("username")
	limit := tools.Atoi(c.Query("limit"), 20)
	before := tools.Atoi(c.Query("before"), 0)

	if username == "" {
		c.JSON(400, gin.H{"error": "Invalid username"})
		return
	}

	posts, err := getPostsByUsername(username, int64(before), int64(limit))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, posts)
}

func GetPostsByTagWithJoin(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag parameter is required"})
		return
	}
	limit := tools.Atoi(c.Query("limit"), 20)
	before := tools.Atoi(c.Query("before"), 0)

	posts, err := getPostsByTagWithJoin(tag, int64(before), int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Even if no posts are found, 'posts' will be an empty slice, which is fine.
	c.JSON(http.StatusOK, posts)
}

func GetPostsByOwnershipWithJoin(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is required"})
		return
	}
	limit := tools.Atoi(c.Query("limit"), 20)
	before := tools.Atoi(c.Query("before"), 0)

	posts, err := getPostsByOwnershipWithJoin(username, int64(before), int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Even if no posts are found, 'posts' will be an empty slice, which is fine.
	c.JSON(http.StatusOK, posts)
}
