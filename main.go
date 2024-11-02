package main

import (
    "database/sql"
    "log"
    "net/http"

    "data-processor-project/internal/api"
    "data-processor-project/internal/repository"
    "data-processor-project/internal/domain/services"
    "data-processor-project/internal/logs"

    _ "github.com/mattn/go-sqlite3"
)

func main() {

    // Initialize logger
    logs.InitLogger()
    defer logs.Sync() 

    db, err := sql.Open("sqlite3", "data_processor.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    repository.Migrate(db)

    userService := services.NewUserService(db)
    requestService := services.NewRequestService(db)


    userAPI := api.NewUserAPI(userService)
    requestAPI := api.NewRequestAPI(requestService)

    http.HandleFunc("/users", userAPI.CreateUser)
    http.HandleFunc("/requests", requestAPI.AddRequest)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
