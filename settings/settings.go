package settings

import (
	"errors"
	"fmt"
	"os"
	"urlq/globals"
	"urlq/links"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbName string = "urlq.db"

// ConnectDB : connect to the DB instance.
func ConnectDB() {
	if globals.Prod {
		// connecting to prod database.
	} else {
		if _, err := os.Stat(dbName); err == nil {
			if !globals.Prod {
				fmt.Println("Connecting to the database...")
			}
			globals.Database, globals.DBerr = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
			if !globals.Prod && globals.DBerr == nil {
				fmt.Println("OK!")
			}
		} else if os.IsNotExist(err) {
			fmt.Println("Creating database.")
			globals.Database, globals.DBerr = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
			globals.MigrationDone = false
			if globals.DBerr != nil {
				fmt.Println(globals.DBerr)
			}
		} else {
			globals.DBerr = errors.New("a problem occured while trying to connect to the database, stopping the server")
		}
	}
	return
}

// Migrate : create the models in the DB instance.
func Migrate() error {
	fmt.Println("Starting migration...")
	err := globals.Database.AutoMigrate(&links.Link{}, &RequestIP{})
	return err
}
