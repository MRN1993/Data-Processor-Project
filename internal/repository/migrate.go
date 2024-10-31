package repository

import (
    "database/sql"
    "log"
)

// Migrate creates necessary tables if they do not exist
func Migrate(db *sql.DB) {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        quota INTEGER
    );`

    requestTable := `
    CREATE TABLE IF NOT EXISTS requests (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        data TEXT,
        received_at DATETIME,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`

    if _, err := db.Exec(userTable); err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }

    if _, err := db.Exec(requestTable); err != nil {
        log.Fatalf("Failed to create requests table: %v", err)
    }
}
