package db

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

var sess db.Session

// init creates the database connection.
func init() {
	var err error
	if sess, err = mongo.Open(mongo.ConnectionURL{
		User:     "myUserAdmin",
		Password: "9AHn2GahV2IXJWHTr80f6dozWEzKMiks3",
		Host:     "127.0.0.1:27017",
		Database: "velvet",
		Options: map[string]string{
			"authSource":    "admin",
			"authMechanism": "SCRAM-SHA-1",
		},
	}); err != nil {
		panic(err)
	}
}
