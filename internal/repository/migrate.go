package repository

import (
    "database/sql"
    "log"
)

// Migrate creates necessary tables if they do not exist
func Migrate(db *sql.DB) {

    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        quota INTEGER,
        monthly_data_limit INTEGER,    -- محدودیت حجم داده در ماه
        request_limit_per_minute INTEGER, -- محدودیت تعداد درخواست در دقیقه
        used_data INTEGER DEFAULT 0,  -- میزان داده‌ی استفاده شده ماهانه
        request_count INTEGER DEFAULT 0, -- تعداد درخواست‌ها در دقیقه فعلی
        last_request_time DATETIME    -- زمان آخرین درخواست
    );`

    
    requestTable := `
    CREATE TABLE IF NOT EXISTS requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
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
