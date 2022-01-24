package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

// init initializes the database and creates the tables required.
func init() {
	var err error
	if db, err = sqlx.Connect("sqlite3", "velvet.db"); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS Players (
		XUID TEXT PRIMARY KEY,
		IGN TEXT UNIQUE NOT NULL COLLATE NOCASE,
		DeviceID TEXT NOT NULL,
		Kills INTEGER DEFAULT 0,
		Deaths INTEGER DEFAULT 0
	)`); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS Bans (
		XUID TEXT PRIMARY KEY,
		IGN TEXT UNIQUE NOT NULL COLLATE NOCASE,
		Mod TEXT NOT NULL,
		Reason TEXT NOT NULL,
		Expires INT NOT NULL
	)`); err != nil {
		panic(err)
	}
}
