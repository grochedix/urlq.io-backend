package main

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PROD : set to true if production, false if development.
const PROD bool = false

// MigrationDone : set to false if migrations needs to be done.
var MigrationDone bool = true

// database : representing the db instance.
var database *gorm.DB

// DBerr : indicates any problem when connecting to the database.
var DBerr error

// connectDB : connect to the DB instance.
func connectDB() {
	if PROD {
		// connecting to PROD database.
	} else {
		if _, err := os.Stat("goAPI.db"); err == nil {
			if !PROD {
				fmt.Println("Connecting to the database...")
			}
			database, DBerr = gorm.Open(sqlite.Open("goAPI.db"), &gorm.Config{})
			if !PROD && DBerr == nil {
				fmt.Println("OK!")
			}
		} else if os.IsNotExist(err) {
			fmt.Println("Creating database.")
			database, DBerr = gorm.Open(sqlite.Open("goAPI.db"), &gorm.Config{})
			MigrationDone = false
			if DBerr != nil {
				fmt.Println(DBerr)
			}
		} else {
			DBerr = errors.New("a problem occured while trying to connect to the database, stopping the server")
		}
	}
	return
}

// migrate : create the models in the DB instance.
func migrate() error {
	fmt.Println("Starting migration...")
	err := database.AutoMigrate(&Link{}, &RequestIP{})
	return err
}
