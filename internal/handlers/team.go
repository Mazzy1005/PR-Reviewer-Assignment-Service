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
		writeError(w, models.ValidationError{Message: "invalid json"})
		return
	}

	if err := h.teamService.AddTeam(r.Context(), &team); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"team": team})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, models.ValidationError{Message: "team_name is required"})
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"team": team})
}
