package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

func SaveActivity(tx *sql.Tx, id string, o *orderedmap.OrderedMap) error {
	content, err := json.Marshal(o)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO activities (id, content)
        VALUES ($1, $2::jsonb)
        ON CONFLICT (id) DO UPDATE SET content = EXCLUDED.content;
    `, id, content)
	if err != nil {
		return err
	}

	return nil
}

func ReadActivity(tx *sql.Tx, id string) (*orderedmap.OrderedMap, error) {
	var content []byte
	// 执行查询
	err := tx.QueryRow(
		`SELECT content FROM activities WHERE id = $1`,
		id,
	).Scan(&content)

	if err != nil {
		return nil, err // 包含 sql.ErrNoRows 等错误
	}

	// 反序列化 JSON 到 OrderedMap
	o := orderedmap.New()
	if err := json.Unmarshal(content, o); err != nil {
		return nil, err
	}

	return o, nil
}

func DeleteActivity(tx *sql.Tx, id string) error {
	_, err := tx.Exec(
		`DELETE FROM activities WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// 根据 actor 和 object 查询符合条件的活动 ID 列表
func QueryByMap(tx *sql.Tx, params map[string]string) ([]*orderedmap.OrderedMap, error) {
	// 使用 JSONB 包含操作符 @> 进行高效查询
	const query = `
        SELECT content 
        FROM activities 
        WHERE content @> $1::jsonb
    `

	// 构造查询参数（确保字段顺序和大小写一致）
	// params := map[string]string{
	// 	"actor":  actor,
	// 	"object": object,
	// }
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(query, paramBytes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objectArray []*orderedmap.OrderedMap
	for rows.Next() {
		var content []byte
		if err := rows.Scan(&content); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		var o *orderedmap.OrderedMap
		err := json.Unmarshal(content, &o)
		if err != nil {
			return objectArray, err
		}
		objectArray = append(objectArray, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return objectArray, nil
}

// not reviewed
// 高级查询：带分页和状态过滤
func QueryActivities(tx *sql.Tx, actor, object, status string, limit, offset int) ([]*orderedmap.OrderedMap, error) {
	query := &strings.Builder{}
	args := []interface{}{}
	query.WriteString(`
        SELECT content 
        FROM activities 
        WHERE 1=1
    `)

	// 动态构建查询条件
	if actor != "" {
		query.WriteString(" AND content->>'actor' = $1")
		args = append(args, actor)
	}
	if object != "" {
		query.WriteString(" AND content->>'object' = $2")
		args = append(args, object)
	}
	if status != "" {
		query.WriteString(" AND content->>'status' = $3")
		args = append(args, status)
	}

	query.WriteString(" ORDER BY (content->>'timestamp')::timestamptz DESC")
	query.WriteString(" LIMIT $4 OFFSET $5")
	args = append(args, limit, offset)

	rows, err := tx.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*orderedmap.OrderedMap
	for rows.Next() {
		var content []byte
		if err := rows.Scan(&content); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		o := orderedmap.New()
		if err := json.Unmarshal(content, o); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}
		results = append(results, o)
	}

	return results, nil
}
