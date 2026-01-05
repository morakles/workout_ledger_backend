package database

import (
	"database/sql"
	"sync"
)

var (
	mu sync.RWMutex
	db *sql.DB
)

// SetDB stores the database connection for the services package.
func SetDB(d *sql.DB) {
	mu.Lock()
	defer mu.Unlock()
	db = d
}

// DB returns the stored database connection \(or nil if not set\).
func DB() *sql.DB {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

func Close() error {
	mu.RLock()
	d := db
	mu.RUnlock()
	if d == nil {
		return nil
	}
	return d.Close()
}
