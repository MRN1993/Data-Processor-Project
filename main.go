package main

import (
    "database/sql"
    "log"
    "net/http"

    "data-processor-project/internal/api"
    "data-processor-project/internal/repository"
    "data-processor-project/internal/domain/services"

    _ "github.com/mattn/go-sqlite3"
    "go.uber.org/zap"
)

func main() {
    db, err := sql.Open("sqlite3", "data_processor.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    repository.Migrate(db)

    // تنظیم logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    userService := services.NewUserService(db)

    // ایجاد RequestRepository و RequestService
    requestRepo := repository.NewSQLRequestRepository(db)
    requestService := services.NewRequestService(requestRepo, logger)

    // ایجاد APIها
    userAPI := api.NewUserAPI(userService)
    requestAPI := api.NewRequestAPI(requestService)

    http.HandleFunc("/users", userAPI.CreateUser)
    http.HandleFunc("/requests", requestAPI.HandleRequest)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
