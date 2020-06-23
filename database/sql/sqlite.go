package sql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type SQLiteConfig struct {
	Path string
}

func (cfg *SQLiteConfig) WithSQLiteDefault() *SQLiteConfig {
	if cfg == nil {
		return &SQLiteConfig{
			Path: "./sqlite.db",
		}
	}

	return cfg
}

// db, err = gorm.Open("mysql", "metro:metro1234@10.252.6.139:3306/crm?charset=utf8mb4")
func (c SQLiteConfig) URI() string {
	return c.Path
}

func NewSQLite(cfg *SQLiteConfig) (db *gorm.DB, err error) {
	c := cfg.WithSQLiteDefault()
	// 返回一个连接池
	db, err = gorm.Open("sqlite3", c.URI())

	return
}
