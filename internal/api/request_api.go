package api

import (
    "encoding/json"
    "net/http"
    "data-processor-project/internal/domain/services"
)

// API represents the HTTP API.
type API struct {
    RequestService *services.RequestService
}

// NewAPI creates a new API instance.
func NewRequestAPI(requestService *services.RequestService) *API {
    return &API{RequestService: requestService}
}

// HandleRequest handles incoming requests.
func (api *API) HandleRequest(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID string `json:"user_id"`
        Data   string `json:"data"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Call the HandleRequest method to process the request
    if err := api.RequestService.HandleRequest(req.UserID, req.Data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Request processed successfully")
}

