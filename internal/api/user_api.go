package api

import (
    "encoding/json"
    "net/http"
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


    var user struct {
        Quota                int `json:"quota"`
        MonthlyDataLimit     int `json:"monthly_data_limit"`
        RequestLimitPerMinute int `json:"request_limit_per_minute"`
    }

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        logs.Logger.Error("Invalid input data", zap.Error(err))
        http.Error(w, "Invalid input data", http.StatusBadRequest)
        return
    }


    if err := api.service.RegisterUser(user.Quota, user.MonthlyDataLimit, user.RequestLimitPerMinute); err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }


    logs.Logger.Info("User created successfully")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
