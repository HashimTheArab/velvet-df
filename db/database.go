package db

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

var sess db.Session

// init initializes the database and creates the tables required.
func init() {
	var err error
	conn, _ := mongo.ParseURL("mongodb+srv://Hashim:9AHn2GahV2IXJWHTr80f6dozWEzKMiks3@practice.oeekd.mongodb.net/test")
	if sess, err = mongo.Open(conn); err != nil {
		panic(err)
	}
}
