package db

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

var sess db.Session

// init creates the database connection.
func init() {
	conn, err := mongo.ParseURL("mongodb://Hashim:9AHn2GahV2IXJWHTr80f6dozWEzKMiks3@127.0.0.1/velvet")
	if err != nil {
		panic(err)
	}
	if sess, err = mongo.Open(conn); err != nil {
		panic(err)
	}
}
