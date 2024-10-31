package main

import (
    "log"
    "net/http"
    "data-processor-project/internal/api"
    "data-processor-project/internal/domain/services"
    "go.uber.org/zap"
)

func main() {
    // Create a new logger
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("cannot initialize zap logger: %v", err)
    }
    defer logger.Sync() // Flushes buffer, if any

    userService := services.NewUserService()
    requestService := services.NewRequestService(logger)
    userAPI := api.NewUserAPI(userService)
    // Assuming you have a request API for handling requests
    requestAPI := api.NewRequestAPI(requestService)

    http.HandleFunc("/users", userAPI.CreateUser)
    http.HandleFunc("/requests", requestAPI.HandleRequest) // Assuming you have this endpoint

    logger.Info("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Fatal("Failed to start server", zap.Error(err))
    }
}