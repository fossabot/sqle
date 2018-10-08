package sqle

import "database/sql"

// Sqle provides methods to simply a few specific database operations.
type Sqle struct {
	db *sql.DB
}

// New intializes a new `sqle.(*Sqle)` instance and returns it.
func New(db *sql.DB) *Sqle {
	return &Sqle{db: db}
}
