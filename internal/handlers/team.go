package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		writeError(w, http.StatusBadRequest, models.TEAM_EXISTS, "invalid json")
		return
	}

	if err := h.teamService.AddTeam(r.Context(), &team); err != nil {
		writeError(w, http.StatusBadRequest, models.TEAM_EXISTS, "team_name already exists")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"team": team})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, models.NOT_FOUND, "team_name is required")
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		writeError(w, http.StatusNotFound, models.NOT_FOUND, "team not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"team": team})
}

// TODO: вынести в отдельный файл
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code models.ErrorResponseErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: struct {
			Code    models.ErrorResponseErrorCode `json:"code"`
			Message string                        `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}
