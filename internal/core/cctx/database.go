package cctx

import (
	"gorm.io/gorm"

	"github.com/saveblush/reraw-api/internal/core/connection/sql"
)

// GetDatabase get connection database
func (c *Context) GetDatabase() *gorm.DB {
	return sql.Database
}
