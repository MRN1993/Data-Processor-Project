package api

import (
    "encoding/json"
    "net/http"
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/domain/services"
)

type UserAPI struct {
    service *services.UserService
}

func NewUserAPI(service *services.UserService) *UserAPI {
    return &UserAPI{service: service}
}

func (api *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input data", http.StatusBadRequest)
        return
    }

    if err := api.service.RegisterUser(user.ID, user.Quota); err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
