package db

import "github.com/df-mc/goleveldb/leveldb"

var db *leveldb.DB

func init() {
	var err error
	db, err = leveldb.OpenFile("velvet.db", nil)
	if err != nil {
		panic(err)
	}
}
