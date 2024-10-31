package api

import (
    "encoding/json"
    "net/http"
    "data-processor-project/internal/domain/services"
)

// UserAPI handles user-related API requests.
type UserAPI struct {
    userService *services.UserService
}

// NewUserAPI creates a new UserAPI.
func NewUserAPI(userService *services.UserService) *UserAPI {
    return &UserAPI{userService: userService}
}

// CreateUser handles the creation of a new user.
func (api *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Quota int `json:"quota"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    user, err := api.userService.CreateUser(req.Quota)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
