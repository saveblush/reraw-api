package sql

import (
	"gorm.io/gorm"

	"github.com/saveblush/reraw-api/internal/core/utils/logger"
)

func Migration(db *gorm.DB) error {
	var sqls []string
	sqls = append(sqls, `
		CREATE TABLE IF NOT EXISTS users (
			pubkey varchar(64) NOT NULL PRIMARY KEY,
			created_at integer DEFAULT NULL,
			updated_at integer DEFAULT NULL,
			deleted_at integer DEFAULT NULL,
			name text DEFAULT NULL,
			lightning_url text DEFAULT NULL
		);
	`)

	// index users
	sqls = append(sqls, `CREATE INDEX IF NOT EXISTS idx_deleted_at ON users (deleted_at);`)
	sqls = append(sqls, "CREATE INDEX IF NOT EXISTS idx_name ON users USING gin (to_tsvector('simple', name));")

	for _, sql := range sqls {
		err := db.Exec(sql).Error
		if err != nil {
			logger.Log.Errorf("db migration error: %s", err)
			return err
		}
	}

	return nil
}
