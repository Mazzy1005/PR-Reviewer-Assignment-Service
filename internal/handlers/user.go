package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req models.PostUsersSetIsActiveJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.ValidationError{Message: "invalid json"})
		return
	}

	user, err := h.userService.SetIsActive(r.Context(), req.UserId, req.IsActive)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, models.PostUsersSetIsActiveResponse{User: *user})
}

func (h *UserHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, models.ValidationError{Message: "user_id is required"})
		return
	}

	resp, err := h.userService.GetReviews(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
