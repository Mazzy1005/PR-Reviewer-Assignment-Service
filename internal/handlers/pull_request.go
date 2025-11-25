package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/service"
)

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func NewPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prService: prService}
}

func (h *PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestCreateJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.ValidationError{Message: "invalid json"})
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, pr)
}

func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestReassignJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.ValidationError{Message: "invalid json"})
		return
	}

	resp, err := h.prService.ReassignReviewer(r.Context(), req.PullRequestId, req.OldUserId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestMergeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.ValidationError{Message: "invalid json"})
		return
	}

	pr, err := h.prService.MergePR(r.Context(), req.PullRequestId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pr)
}
