package sql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

func (cfg *MySQLConfig) WithMySQLDefault() *MySQLConfig {
	if cfg == nil {
		return nil
	}

	c := *cfg

	if cfg.Port == 0 {
		c.Port = 3306
	}
	if cfg.Host == "" {
		c.Host = "127.0.0.1"
	}

	if cfg.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}

	if cfg.MaxOpenConns == 0 {
		c.MaxOpenConns = 80
	}

	if cfg.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 0
	}

	return &c
}

// db, err = gorm.Open("mysql", "metro:metro1234@10.252.6.139:3306/crm?charset=utf8mb4")
func (c MySQLConfig) URI() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4", c.User, c.Password, c.Host, c.Port, c.DBName)
}

func NewMySQL(cfg *MySQLConfig) (db *gorm.DB, err error) {
	c := cfg.WithMySQLDefault()
	// 返回一个连接池
	db, err = gorm.Open("mysql", c.URI())

	// 防止无线连接，出现 too many connections 错误
	// MySQL默认是151
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime))
	// db.LogMode(false)

	return
}
