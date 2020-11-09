package globals

import "gorm.io/gorm"

// Prod : set to true if production, false if development.
const Prod bool = false

// MigrationDone : set to false if migrations needs to be done.
var MigrationDone bool = true

// Database : representing the db instance.
var Database *gorm.DB

// DBerr : indicates any problem when connecting to the database.
var DBerr error
