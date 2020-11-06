package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PROD should be set to true if production, false if devlopment.
const PROD bool = false

// MigrationDone indicates if migrations needs to be done.
var MigrationDone bool = true

// DB is the var for db
var DB *gorm.DB

// DBerr is the var indicating any problem connecting to the DB
var DBerr error

func connectDB() {
	if PROD {

	} else {
		// In development, we'll be using a sqlite db.
		DB, DBerr = gorm.Open(sqlite.Open("goAPI.db"), &gorm.Config{})
	}
	return
}

func migrate() error {
	err := DB.AutoMigrate(&Link{}, &RequestIP{})
	return err
}
