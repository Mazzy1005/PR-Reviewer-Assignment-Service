package models

type ErrorResponseErrorCode string

const (
	TEAM_EXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
	PR_EXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PR_MERGED    ErrorResponseErrorCode = "PR_MERGED"
	NOT_ASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOT_FOUND    ErrorResponseErrorCode = "NOT_FOUND"
	BAD_REQUEST  ErrorResponseErrorCode = "BAD_REQUEST"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type AppError struct {
	Code    ErrorResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(code ErrorResponseErrorCode) AppError {
	var msg string

	switch code {
	case TEAM_EXISTS:
		msg = "team_name already exists"
	case PR_EXISTS:
		msg = "PR id already exists"
	case PR_MERGED:
		msg = "cannot reassign on merged PR"
	case NOT_ASSIGNED:
		msg = "reviewer is not assigned to this PR"
	case NO_CANDIDATE:
		msg = "no active replacement candidate in team"
	case NOT_FOUND:
		msg = "resource not found"
	default:
		msg = "unknown error"
	}
	return AppError{Code: code, Message: msg}
}

type ErrorResponse struct {
	Error AppError `json:"error"`
}
