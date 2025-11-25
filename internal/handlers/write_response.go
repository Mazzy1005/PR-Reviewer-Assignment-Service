package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	errorCode := models.NOT_FOUND
	message := "internal server error"

	switch e := err.(type) {
	case models.ValidationError:
		statusCode = http.StatusBadRequest
		errorCode = models.BAD_REQUEST
		message = e.Message

	case models.AppError:
		message = e.Message
		errorCode = e.Code

		switch e.Code {
		case models.PR_MERGED, models.NOT_ASSIGNED, models.PR_EXISTS, models.NO_CANDIDATE:
			statusCode = http.StatusConflict
		case models.TEAM_EXISTS:
			statusCode = http.StatusBadRequest
		case models.NOT_FOUND:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
	default:
		if err != nil {
			message = "unexpected error"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: struct {
			Code    models.ErrorResponseErrorCode `json:"code"`
			Message string                        `json:"message"`
		}{
			Code:    errorCode,
			Message: message,
		},
	})
}
