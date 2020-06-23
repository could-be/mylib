package sql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlite(t *testing.T) {
	path := "/Users/goods/go/metro/tmp/db/dmall/dmall.db"
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=rwc", path))
	if err != nil {
		t.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec("UPDATE custregisterstatus set status='submit' where status=?;", "new"); err != nil {
		tx.Rollback()
		return
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
