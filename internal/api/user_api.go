package api

import (
    "encoding/json"
    "net/http"
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/domain/services"
    "data-processor-project/internal/logs"
    
    "go.uber.org/zap"
)

type UserAPI struct {
    service *services.UserService
}

func NewUserAPI(service *services.UserService) *UserAPI {
    return &UserAPI{service: service}
}

func (api *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    logs.Logger.Info("Received request to create user")

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        logs.Logger.Error("Invalid input data", zap.Error(err))
        http.Error(w, "Invalid input data", http.StatusBadRequest)
        return
    }

    if err := api.service.RegisterUser(user.ID, user.Quota); err != nil {
        logs.Logger.Error("Failed to register user", zap.Error(err), zap.Int("userID", user.ID))
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    logs.Logger.Info("User created successfully", zap.Int("userID", user.ID))
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
