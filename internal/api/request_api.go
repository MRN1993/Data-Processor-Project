package api

import (
    "encoding/json"
    "net/http"
    "data-processor-project/internal/domain/services"
)

type RequestAPI struct {
    service *services.RequestService
}

func NewRequestAPI(service *services.RequestService) *RequestAPI {
    return &RequestAPI{service: service}
}

type Request struct {
    ID     int    `json:"id"`
    UserID int    `json:"user_id"`
    Data   string `json:"data"`
}

// AddRequest godoc
// @Summary      Add a new request
// @Description  Adds a new request to be processed
// @Tags         requests
// @Accept       json
// @Produce      json
// @Param        request body Request true "Request data"
// @Success      201 {string} string "Request processed successfully"
// @Failure      400 {string} string "Invalid request payload"
// @Failure      500 {string} string "Internal server error"
// @Router       /requests [post]
func (api *RequestAPI) AddRequest(w http.ResponseWriter, r *http.Request) {
    
    var request Request
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := api.service.ProcessRequest(request.ID,request.UserID,request.Data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Request processed successfully"))
}