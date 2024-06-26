package db

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
	"runtime"
)

var sess db.Session

// init creates the database connection.
func init() {
	host := "velvetpractice.tk:26919"
	if runtime.GOOS == "linux" {
		host = "127.0.0.1:27017"
	}

	var err error
	if sess, err = mongo.Open(mongo.ConnectionURL{
		User:     "Hashim",
		Password: "9AHn2GahV2IXJWHTr80f6dozWEzKMiks3",
		Host:     host,
		Database: "velvet",
		Options: map[string]string{
			"authSource":    "admin",
			"authMechanism": "SCRAM-SHA-1",
		},
	}); err != nil {
		panic(err)
	}
}
