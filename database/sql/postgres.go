package sql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresConfig struct {
	User   string
	Passwd string
	Host   string
	Port   int
	DBName string

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

// func (cfg *PostgresConfig) WithPort() *PostgresConfig {
// 	if cfg.Port == 0 {
// 		cfg.Port = 5432
// 	}
// 	return cfg
// }
//
// func (cfg *PostgresConfig) WithHost() *PostgresConfig {
// }

func (cfg *PostgresConfig) WithPostgresDefault() *PostgresConfig {
	if cfg == nil {
		return nil
	}

	c := *cfg

	if cfg.Port == 0 {
		c.Port = 5432
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

// db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
func (c PostgresConfig) URI() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.Host, c.Port, c.User, c.DBName, c.Passwd)
}

func (cfg *PostgresConfig) String() string {
	b, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err.Error()
	}
	return bytes.NewBuffer(b).String()
}

func NewMyPostgres(cfg *PostgresConfig) (db *gorm.DB, err error) {
	c := cfg.WithPostgresDefault()

	// 返回一个连接池
	db, err = gorm.Open("postgres", c.URI())

	// 防止无线连接，出现 too many connections 错误
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime))
	// db.LogMode(false)

	// 在外层，应用控制
	// db.Debug().AutoMigrate(&PostgresConfig{})
	return
}
