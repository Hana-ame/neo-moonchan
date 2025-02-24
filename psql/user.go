package psql

import (
	"database/sql"
	"encoding/json"

	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

// func SaveUser(ctx context.Context, conn *pgx.Conn, id string, user *orderedmap.OrderedMap) error {
// 	content, err := json.Marshal(user)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = conn.Exec(ctx, `
//         INSERT INTO users (id, content)
//         VALUES ($1, $2::jsonb)
//         ON CONFLICT (id) DO UPDATE SET content = EXCLUDED.content;
//     `, id, content)
// 	return err
// }

func SaveUser(tx *sql.Tx, id string, user *orderedmap.OrderedMap) error {
	content, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO users (id, content)
        VALUES ($1, $2::jsonb)
        ON CONFLICT (id) DO UPDATE SET content = EXCLUDED.content;
    `, id, content)
	if err != nil {
		return err
	}

	return nil
}

func ReadUser(tx *sql.Tx, id string) (*orderedmap.OrderedMap, error) {
	var content []byte
	// 执行查询
	err := tx.QueryRow(
		`SELECT content FROM users WHERE id = $1`,
		id,
	).Scan(&content)

	if err != nil {
		return nil, err // 包含 sql.ErrNoRows 等错误
	}

	// 反序列化 JSON 到 OrderedMap
	user := orderedmap.New()
	if err := json.Unmarshal(content, user); err != nil {
		return nil, err
	}

	return user, nil
}
