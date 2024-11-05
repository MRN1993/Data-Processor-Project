package repository

import (
    "database/sql"
    "go.uber.org/zap"
    "data-processor-project/internal/logs"
)

func Migrate(db *sql.DB) {

    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        monthly_data_limit INTEGER NOT NULL,   
        request_limit_per_minute INTEGER NOT NULL,
        used_data INTEGER DEFAULT 0,  
        request_count INTEGER DEFAULT 0, 
        last_request_time DATETIME   
    );`

    requestTable := `
    CREATE TABLE IF NOT EXISTS requests (
        request_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        data STRING NOT NULL,
        received_at DATETIME,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`

    logs.Logger.Info("Starting migration of tables...")

    if _, err := db.Exec(userTable); err != nil {
        logs.Logger.Fatal("Failed to create users table", zap.Error(err))
    } else {
        logs.Logger.Info("Users table created or already exists")
    }

    if _, err := db.Exec(requestTable); err != nil {
        logs.Logger.Fatal("Failed to create requests table", zap.Error(err))
    } else {
        logs.Logger.Info("Requests table created or already exists")
    }

    logs.Logger.Info("Migration completed successfully")
}
